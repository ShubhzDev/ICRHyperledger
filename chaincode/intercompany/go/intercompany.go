package main

import (
    "encoding/json" // Package for encoding and decoding JSON data
    "fmt"           // Package for formatted I/O
    "log"           // Package for logging
    "time"          // Package for handling time

    "github.com/hyperledger/fabric-contract-api-go/contractapi" // Hyperledger Fabric contract API package
)

// Transaction represents an intercompany transaction
type Transaction struct {
    ID              string    `json:"id"`              // Identifier of the transaction
    Company         string    `json:"company"`         // Company initiating the transaction
    Counterparty    string    `json:"counterparty"`    // Company receiving the transaction
    Amount          float64   `json:"amount"`          // Amount of money involved in the transaction
    TransactionType string    `json:"transactionType"` // Type of transaction
    Date            time.Time `json:"date"`            // Date of the transaction
    Reconciled      bool      `json:"reconciled"`      // Indicates whether the transaction is reconciled
}

// SmartContract provides functions for managing transactions
type SmartContract struct {
    contractapi.Contract // Hyperledger Fabric contract API
}

// AddTransaction adds a new transaction to the ledger
func (s *SmartContract) AddTransaction(ctx contractapi.TransactionContextInterface, id string, company string, counterparty string, amount float64, transactionType string, date string) error {
    // Parsing the transaction date
    transactionDate, err := time.Parse("2006-01-02", date)
    if err != nil {
        return fmt.Errorf("failed to parse date: %v", err)
    }

    // Creating a new transaction instance
    transaction := Transaction{
        ID:              id,
        Company:         company,
        Counterparty:    counterparty,
        Amount:          amount,
        TransactionType: transactionType,
        Date:            transactionDate,
        Reconciled:      false,
    }

    // Marshaling the transaction to JSON format
    transactionJSON, err := json.Marshal(transaction)
    if err != nil {
        return fmt.Errorf("failed to marshal transaction: %v", err)
    }

    // Storing the transaction in the ledger
    return ctx.GetStub().PutState(id, transactionJSON)
}

// QueryTransaction returns the transaction stored in the ledger with the given id
func (s *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface, id string) (*Transaction, error) {
    // Retrieving the transaction JSON from the ledger
    transactionJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if transactionJSON == nil {
        return nil, fmt.Errorf("the transaction %s does not exist", id)
    }

    // Unmarshaling the transaction JSON into a Transaction struct
    var transaction Transaction
    err = json.Unmarshal(transactionJSON, &transaction)
    if err != nil {
        return nil, err
    }

    return &transaction, nil
}

// MatchTransactions matches transactions based on predefined rules and marks them as reconciled
func (s *SmartContract) MatchTransactions(ctx contractapi.TransactionContextInterface, id1 string, id2 string) error {
    // Retrieving the transactions from the ledger
    transaction1, err := s.QueryTransaction(ctx, id1)
    if err != nil {
        return err
    }
    transaction2, err := s.QueryTransaction(ctx, id2)
    if err != nil {
        return err
    }

    // Checking if either transaction is already reconciled
    if transaction1.Reconciled || transaction2.Reconciled {
        return fmt.Errorf("one or both transactions are already reconciled")
    }

    // Checking if the counterparties match
    if transaction1.Company != transaction2.Counterparty || transaction2.Company != transaction1.Counterparty {
        return fmt.Errorf("counterparty mismatch")
    }

    // Checking if the amounts match (one is the negative of the other)
    if transaction1.Amount != -transaction2.Amount {
        return fmt.Errorf("amount mismatch")
    }

    // Checking if the transaction dates match
    if transaction1.Date.Year() != transaction2.Date.Year() || transaction1.Date.Month() != transaction2.Date.Month() {
        return fmt.Errorf("date mismatch")
    }

    // Marking both transactions as reconciled
    transaction1.Reconciled = true
    transaction2.Reconciled = true

    // Marshaling the updated transactions back to JSON
    transaction1JSON, err := json.Marshal(transaction1)
    if err != nil {
        return err
    }
    transaction2JSON, err := json.Marshal(transaction2)
    if err != nil {
        return err
    }

    // Updating the ledger with the reconciled transactions
    err = ctx.GetStub().PutState(id1, transaction1JSON)
    if err != nil {
        return err
    }
    return ctx.GetStub().PutState(id2, transaction2JSON)
}

// Entry point of the program
func main() {
    // Creating a new chaincode instance
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        log.Panicf("Error creating intercompany chaincode: %v", err)
    }

    // Starting the chaincode
    if err := chaincode.Start(); err != nil {
        log.Panicf("Error starting intercompany chaincode: %v", err)
    }
}
