package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Transaction represents an intercompany transaction
type Transaction struct {
    ID              string    `json:"id"`
    Company         string    `json:"company"`
    Counterparty    string    `json:"counterparty"`
    Amount          float64   `json:"amount"`
    TransactionType string    `json:"transactionType"`
    Date            time.Time `json:"date"`
    Reconciled      bool      `json:"reconciled"`
}

// SmartContract provides functions for managing transactions
type SmartContract struct {
    contractapi.Contract
}

// AddTransaction adds a new transaction to the ledger
func (s *SmartContract) AddTransaction(ctx contractapi.TransactionContextInterface, id string, company string, counterparty string, amount float64, transactionType string, date string) error {
    transactionDate, err := time.Parse("2006-01-02", date)
    if err != nil {
        return fmt.Errorf("failed to parse date: %v", err)
    }

    transaction := Transaction{
        ID:              id,
        Company:         company,
        Counterparty:    counterparty,
        Amount:          amount,
        TransactionType: transactionType,
        Date:            transactionDate,
        Reconciled:      false,
    }

    transactionJSON, err := json.Marshal(transaction)
    if err != nil {
        return fmt.Errorf("failed to marshal transaction: %v", err)
    }

    return ctx.GetStub().PutState(id, transactionJSON)
}

// QueryTransaction returns the transaction stored in the ledger with the given id
func (s *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface, id string) (*Transaction, error) {
    transactionJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if transactionJSON == nil {
        return nil, fmt.Errorf("the transaction %s does not exist", id)
    }

    var transaction Transaction
    err = json.Unmarshal(transactionJSON, &transaction)
    if err != nil {
        return nil, err
    }

    return &transaction, nil
}

// MatchTransactions matches transactions based on predefined rules and marks them as reconciled
func (s *SmartContract) MatchTransactions(ctx contractapi.TransactionContextInterface, id1 string, id2 string) error {
    transaction1, err := s.QueryTransaction(ctx, id1)
    if err != nil {
        return err
    }
    transaction2, err := s.QueryTransaction(ctx, id2)
    if err != nil {
        return err
    }

    if transaction1.Reconciled || transaction2.Reconciled {
        return fmt.Errorf("one or both transactions are already reconciled")
    }

    if transaction1.Company != transaction2.Counterparty || transaction2.Company != transaction1.Counterparty {
        return fmt.Errorf("counterparty mismatch")
    }

    if transaction1.Amount != -transaction2.Amount {
        return fmt.Errorf("amount mismatch")
    }

    if transaction1.Date.Year() != transaction2.Date.Year() || transaction1.Date.Month() != transaction2.Date.Month() {
        return fmt.Errorf("date mismatch")
    }

    transaction1.Reconciled = true
    transaction2.Reconciled = true

    transaction1JSON, err := json.Marshal(transaction1)
    if err != nil {
        return err
    }
    transaction2JSON, err := json.Marshal(transaction2)
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(id1, transaction1JSON)
    if err != nil {
        return err
    }
    return ctx.GetStub().PutState(id2, transaction2JSON)
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        log.Panicf("Error creating intercompany chaincode: %v", err)
    }

    if err := chaincode.Start(); err != nil {
        log.Panicf("Error starting intercompany chaincode: %v", err)
    }
}
