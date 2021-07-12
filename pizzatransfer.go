/*
SPDX-License-Identifier: Apache-2.0
*/

package main

/*
./network.sh up createChannel
./network.sh deployCC -ccn pizzacc -ccp ../../pizza-cc/ -ccl go


*/
import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Order struct {
	Pizza    string  `json:"pizza"`
	ID       string  `json:"ID"`
	Price    float32 `json:"price"`
	Holder   string  `json:"holder"`
	Datetime string  `json:"datetime"`
	Address  string  `json:"Address"`
}

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	orders := []Order{
		{Pizza: "margherita", ID: "1", Price: 7.5, Holder: "Satoshi", Datetime: "1230969425", Address: "50.0010, -60.3425"},
		{Pizza: "diavolo", ID: "2", Price: 10.75, Holder: "Vitalik", Datetime: "1284618372", Address: "30.1928, 50.3817"},
		{Pizza: "melanzano", ID: "3", Price: 12, Holder: "Ignaz", Datetime: "1386453183", Address: "10.1837, 53.1823"},
		{Pizza: "quattro formaggi", ID: "4", Price: 11, Holder: "Xabi", Datetime: "4531234861", Address: "89.1738, 23.1382"},
	}

	for _, order := range orders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(order.ID, orderJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) CreateOrder(ctx contractapi.TransactionContextInterface, pizza string, ID string,
	price float32, holder string, datetime string, Address string) error {
	exists, err := s.OrderExists(ctx, ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the order %s already exists", ID)
	}

	order := Order{
		Pizza:    pizza,
		ID:       ID,
		Price:    price,
		Holder:   holder,
		Datetime: datetime,
		Address:  Address,
	}
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, orderJSON)
}

func (s *SmartContract) UpdateOrder(ctx contractapi.TransactionContextInterface, pizza string, id string,
	price float32, holder string, datetime string, Address string) error {
	exists, err := s.OrderExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the order %s does not exist", id)
	}

	// overwriting original asset with new asset
	order := Order{
		Pizza:    pizza,
		ID:       id,
		Price:    price,
		Holder:   holder,
		Datetime: datetime,
		Address:  Address,
	}
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, orderJSON)
}

func (s *SmartContract) ReadOrder(ctx contractapi.TransactionContextInterface, id string) (*Order, error) {
	orderJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("the order %s does not exist", id)
	}

	var order Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *SmartContract) DeleteOrder(ctx contractapi.TransactionContextInterface, ID string) error {
	exists, err := s.OrderExists(ctx, ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the order %s does not exist", ID)
	}

	return ctx.GetStub().DelState(ID)
}

func (s *SmartContract) OrderExists(ctx contractapi.TransactionContextInterface, ID string) (bool, error) {
	orderJSON, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return orderJSON != nil, nil
}

func (s *SmartContract) TransferOrder(ctx contractapi.TransactionContextInterface, id string, owner string) error {
	order, err := s.ReadOrder(ctx, id)
	if err != nil {
		return err
	}

	order.Holder = owner
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, orderJSON)
}

func (s *SmartContract) GetAllOrders(ctx contractapi.TransactionContextInterface) ([]*Order, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var orders []*Order
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var order Order
		err = json.Unmarshal(queryResponse.Value, &order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func main() {
	orderChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating order-transfer-basic chaincode: %v", err)
	}

	if err := orderChaincode.Start(); err != nil {
		log.Panicf("Error starting order-transfer-basic chaincode: %v", err)
	}
}
