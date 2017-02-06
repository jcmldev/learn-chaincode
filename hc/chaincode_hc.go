/*
Copyright IBM Corp 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*
Humanity coins 

The scenario consists of accounts holding coins. 
The ledger maintains state of the accounts. Coins can be transfered between pairs of accounts. 
An account can be closed with remaning balance transferred to a master account.  

Functions
1 - Open account - with given account name and default value of 200
2 - Add coins to account
3 - Transfer coins between two accounts
4 - Close account - transfers remaining coins into a master account  
5 - Enquire an account balance
*/



type HumanityCoinsChaincode struct {
}

const ACCOUNT_MASTER_NAME		= "master_account"
const ACCOUNT_DEFAULT_BALANCE	= 200

const FUNCTION_OPEN_ACCOUNT		= "open_account"	
const FUNCTION_ADD_COINS		= "add_coins"
const FUNCTION_TRANSFER_COINS	= "transfer_coins"
const FUNCTION_CLOSE_ACCOUNT	= "close_account"
const FUNCTION_ACCOUNT_BALANCE  = "account_balance"


func main() {
	err := shim.Start(new(HumanityCoinsChaincode))
	if err != nil {
		fmt.Printf("Error starting Humanity coins chaincode: %s", err)
	}
}


func IntToDBValue(value int) ([]byte) {
	return []byte(strconv.Itoa(value))
}

func DBValueToInt(value []byte) (int, error) {
	return strconv.Atoi(string(value))
}


func (t *HumanityCoinsChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Init is running")
	
	/*
	if len(args) >= 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}
	*/

	err := stub.PutState(ACCOUNT_MASTER_NAME, IntToDBValue(0))
	if err != nil { return nil, err }

	return nil, nil
}


func (t *HumanityCoinsChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running function: " + function)

	if function == "init" {													
		return t.Init(stub, "", args)
	} else if function == FUNCTION_OPEN_ACCOUNT {										
		return t.OpenAccount(stub, args)
	} else if function == FUNCTION_ADD_COINS {											
		return t.AddCoins(stub, args)
	} else if function == FUNCTION_TRANSFER_COINS {									
		return t.TransferCoins(stub, args)
	} else if function == FUNCTION_CLOSE_ACCOUNT {										
		return t.CloseAccount(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}


func (t *HumanityCoinsChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if function == FUNCTION_ACCOUNT_BALANCE {
		return t.GetAccountBalance(stub, args)
	}
	
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}


func (t *HumanityCoinsChaincode) OpenAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running OpenAccount()")

	var account_name string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. Name of the account")
	}

	account_name = args[0]
	
	err = stub.PutState(account_name, IntToDBValue(ACCOUNT_DEFAULT_BALANCE))
	if err != nil { return nil, err }
	
	return nil, nil
}


func (t *HumanityCoinsChaincode) AddCoins(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running AddCoins()")

	var account_name string
	var account_balance_arr []byte
	var account_balance, amount_to_add int
	var err error


	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. Name of the account and amount of coins to add")
	}

	account_name = args[0]
	amount_to_add, err = strconv.Atoi(args[1])
	if err != nil {return nil, err }
	
	account_balance_arr, err = t.GetAccountBalance(stub, []string{account_name})
	if err != nil { return nil, err }
	account_balance, err = DBValueToInt(account_balance_arr)	
	if err != nil { return nil, err }
	
	account_balance += amount_to_add
	
	err = stub.PutState(account_name, IntToDBValue(account_balance)) 
	if err != nil { return nil, err }
	
	return nil, nil
}


func (t *HumanityCoinsChaincode) TransferCoins(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running TransferCoins()")

	var account_from_name, account_to_name string
	var account_from_balance_arr, account_to_balance_arr []byte
	var account_from_balance, account_to_balance, amount_to_transfer int
	var err error


	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. Names of the from and to accounts and amount of coins to transfer")
	}

	account_from_name = args[0]
	account_to_name = args[1]
	amount_to_transfer, err = strconv.Atoi(args[2])
	if err != nil {return nil, err }
	
	account_from_balance_arr, err = t.GetAccountBalance(stub, []string{account_from_name})
	if err != nil {return nil, err }
	account_from_balance, err = DBValueToInt(account_from_balance_arr)	
	if err != nil {return nil, err }

	account_to_balance_arr, err = t.GetAccountBalance(stub, []string{account_to_name})
	if err != nil {return nil, err }
	account_to_balance, err = DBValueToInt(account_to_balance_arr)	
	if err != nil {return nil, err }
	
	if amount_to_transfer <= 1 {
		return nil, errors.New("The amount of coins to transfer must be higher than 0")
	}
	
	if account_from_balance < amount_to_transfer {
		return nil, errors.New("Account: " + account_from_name + " - does not have sufficient funds")		
	}
	
	account_from_balance -= amount_to_transfer
	account_to_balance += amount_to_transfer
		
	err = stub.PutState(account_from_name, IntToDBValue(account_from_balance)) 
	if err != nil {return nil, err }
	
	err = stub.PutState(account_to_name, IntToDBValue(account_to_balance)) 
	if err != nil {return nil, err }
	
	return nil, nil
}


func (t *HumanityCoinsChaincode) CloseAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running CloseAccount()")

	var account_name string
	var account_balance_arr []byte
	var account_balance int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. Name of the account to close")
	}

	account_name = args[0]
	
	account_balance_arr, err = t.GetAccountBalance(stub, []string{account_name})
	if err != nil { return nil, err }
	account_balance, err = DBValueToInt(account_balance_arr)	
	if err != nil { return nil, err }
	
	t.TransferCoins(stub, []string{account_name, ACCOUNT_MASTER_NAME, strconv.Itoa(account_balance)})
	
	err = stub.PutState(account_name, IntToDBValue(ACCOUNT_DEFAULT_BALANCE)) 
	if err != nil { return nil, err }
	
	return nil, nil
}


func (t *HumanityCoinsChaincode) GetAccountBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running GetAccountBalance()")
	
	var account_name string
	var account_balance []byte
	var err error


	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. Name of the account.")
	}

	account_name = args[0]
	
	account_balance, err = stub.GetState(account_name)
	if err != nil { return nil, errors.New("Failed to get state for account: " + account_name) }

	return account_balance, nil
}
