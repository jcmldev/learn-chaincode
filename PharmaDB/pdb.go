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
//	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*
Pharmaceutical database for validation of authenticity of sold drugs
*/



type PharmaAuthDB struct {
}


func main() {
	err := shim.Start(new(PharmaAuthDB))
	if err != nil {
		fmt.Printf("Error starting Humanity coins chaincode: %s", err)
	}
}

func (t *PharmaAuthDB) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running function: " + function)

	if function == "init" {													
		return t.Init(stub, "", args)
	} else if function == "create_record" {										
		return t.CreateRecord(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *PharmaAuthDB) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if function == "get_record" {
		return t.GetRecord(stub, args)
	}
	
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *PharmaAuthDB) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Init is running")
	
	return nil, nil
}

func (t *PharmaAuthDB) CreateRecord(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running CreateRecord()")

	var drug_name string
	var packaging_image string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}

	drug_name = args[0]
	packaging_image = args[1]
	
	err = stub.PutState(drug_name, []byte(packaging_image))
	if err != nil { return nil, err }
	
	return nil, nil
}

func (t *PharmaAuthDB) GetRecord(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running CreateRecord()")

	var drug_name string
	var packaging_image []byte
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}

	drug_name = args[0]
	
	
	packaging_image, err = stub.GetState(drug_name)
	if err != nil { return nil, errors.New("Failed to get state for drug: " + drug_name) }

	return packaging_image, nil
}

