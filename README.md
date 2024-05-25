#Architecture Components

##Organizations:
- Org1: Represents subsidiary CompanyA.
- Org2: Represents subsidiary CompanyB.
- ParentOrg: Represents the parent company that maintains the reconciled ledger.

##Peers:
- Peer0.Org1: Peer node for Org1.
- Peer0.Org2: Peer node for Org2.
- Peer0.ParentOrg: Peer node for ParentOrg.

##Certificate Authorities (CAs):
- CA.Org1: Issues certificates for Org1.
- CA.Org2: Issues certificates for Org2.
- CA.ParentOrg: Issues certificates for ParentOrg.

##Orderer:
- Orderer: Manages the ordering of transactions across the network.

##Channel:
- mychannel: The communication channel that all peers are part of. This channel allows transactions between Org1, Org2, and ParentOrg.

##Transaction Flow

##Adding Transactions:
- Org1 and Org2 add transactions to the ledger using their respective peers.
- Each transaction contains details such as id, company, counterparty, amount, transactionType, and date.

##Matching Transactions:
- A function within the chaincode matches transactions between Org1 and Org2 based on predefined rules (e.g., complementary amounts, same fiscal period, counterparty company).
Matched transactions are marked as reconciled.

##Recording Reconciled Transactions:
- Reconciled transactions are recorded in the parent company's ledger.(Recommended)
- The chaincode ensures that these transactions are added to a separate ledger entry managed by ParentOrg.

##Querying Transactions:
- Users can query transactions by their unique identifier to check their status (reconciled or not).

##Diagram
Here's a high-level architecture diagram representing the updated structure:
- https://drive.google.com/file/d/1GnckqZzOznEyaql9ZMY-Uq4oSlkLop7o/view?usp=sharing

##Detailed Explanation

###Organizations and Peers:
- Org1 and Org2 represent two subsidiaries. Each has its own peer (Peer0.Org1 and Peer0.Org2), which hosts the chaincode and maintains the ledger for each organization.
ParentOrg represents the parent company. It has its own peer (Peer0.ParentOrg), which also hosts the chaincode and maintains a ledger for reconciled transactions.

###Certificate Authorities:
- Each organization has its own CA to manage identities and permissions within the network.
- CA.Org1, CA.Org2, and CA.ParentOrg issue certificates to their respective peers and users.

###Orderer:
- The single orderer node ensures that all transactions are ordered correctly and consistently across the network.

###Channel:
- All peers (from Org1, Org2, and ParentOrg) are part of a single channel (mychannel), allowing them to communicate and transact with each other.

###Transaction Matching and Reconciliation:
- The chaincode includes functions for adding transactions, matching them based on predefined rules, and marking matched transactions as reconciled.
- Once transactions between Org1 and Org2 are reconciled, they are recorded in a ledger entry managed by Peer0.ParentOrg.
- This ensures that the parent company has an up-to-date record of all reconciled transactions between its subsidiaries.

###Implementation Steps

###Setup Certificate Authorities:
- Deploy CA services for Org1, Org2, and ParentOrg to manage identities and issue certificates.

###Setup Peers and Orderer:
- Deploy and configure peers for each organization.
- Deploy and configure the orderer node.

###Create and Join Channel:
- Create a channel configuration transaction.
- Join peers from Org1, Org2, and ParentOrg to the channel.

###Deploy Chaincode:
- Package and install the chaincode on all peers.
- Approve and commit the chaincode definition for all organizations.

###Invoke Chaincode Functions:
- Use CLI or SDK to add transactions, match them, and record reconciled transactions in the parent company's ledger.

By following this architecture, you can effectively manage and reconcile intercompany transactions between subsidiaries, maintaining a clear and accurate ledger for the parent company.

How to Start Project:-
./teardown.sh
./generateArtifacts.sh
./start_network.sh

#Design of Network :-
Organizations - ParentOrg,Org1,Org2
Orderer - OrdererOrg
Channel - myChannel(all will have same channel)
Peer - org1.peer0 , org2.peer1 , parentOrg.peer0(nodes)
Ledger - L(all will have same ledger,can restrict access on chaincode)
Smart Contract - InterCompany