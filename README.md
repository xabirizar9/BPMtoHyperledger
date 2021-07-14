# BPMtoHyperledger

In order to run this code, you will need to have installed the hyperledger fabric test network. For this, run 

curl -sSL https://bit.ly/2ysbOFE | bash -s

in the desired directory.


This should create a new directory inside the current one named "fabric-samples". On the folder where you save "fabric-samples", create a new folder where you will save all the chaincode related files (this repository).

Your setup should look like this, for example:

<your-directory>/fabric-samples/  
<your-directory>/pizza-cc

Then, go to 

fabric-samples/test-network/



# HOW TO RUN THE NETWORK
**All of these commands must be run inside fabric-samples/test-network/**  

./network.sh up createChannel -c pizzachannel  
./network.sh deployCC -c pizzachannel -ccn pizzacc -ccp ../../pizza-cc/ -ccl go

export PATH=${PWD}/../bin:$PATH  
export FABRIC_CFG_PATH=$PWD/../config/

# Environment variables for Org1

export CORE_PEER_TLS_ENABLED=true  

export CORE_PEER_LOCALMSPID="Org1MSP"

exportCORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt  

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp  

export CORE_PEER_ADDRESS=localhost:7051  

# INITIALIZE LEDGER
 
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C pizzachannel -n pizzacc --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C pizzachannel -n pizzacc --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'


# RUN THIS COMMAND TO SEE EVERY ORDER AND ITS STATUS
**If the pizza has been delivered, the variable "holder" stores the name of the customer, as it is the end state.**  

peer chaincode query -C pizzachannel -n pizzacc -c '{"Args":["GetAllOrders"]}'

# CREATE NEW PIZZA ORDER

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C pizzachannel -n pizzacc --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"CreateOrder","Args":["salsiccia", "9", "10", "Anuel AA", "1230969425", "52.0010, -63.3425"]}'

# RUN THIS COMMAND TO TRANSFER THE ORDER TO THE NEXT PEER

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C pizzachannel -n pizzacc --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"TransferOrder","Args":["9","Delivery guy"]}'

# SEE THE CHANGE BY QUERYING AGAIN THE LEDGER

peer chaincode query -C pizzachannel -n pizzacc -c '{"Args":["GetAllOrders"]}'
