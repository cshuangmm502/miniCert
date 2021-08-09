/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type miniCC struct {
}

const ISSUESTATE = "issuedState"
const REVOKESTATE = "revokeState"

// main function starts up the chaincode in the container during instantiate
func main()  {
	if err := shim.Start(new(miniCC)); err != nil {
		fmt.Printf("Error starting miniCC chaincode: %s", err)
	}
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
//Init rsa accumulator state and the initial set
func (t *miniCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	accumulator := New()
	N := accumulator.GetN()
	set := make(map[string]int)
	A0 := accumulator.GetA0()
	accstate := &Record{
		set: set,
		A:  A0.String(),
		N:   N.String(),
	}
	//stateAsBytes,err := json.Marshal(accstate.A)
	//if err!= nil{
	//	return shim.Error(err.Error())
	//}
	err := stub.PutState(ISSUESTATE,[]byte(accstate.A))
	if err!=nil{
		return shim.Error(err.Error())
	}

	fmt.Printf("Init accumulator %s \n", accstate.A)

	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *miniCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error


	if fn == "getState" {
		return t.queryState(stub, args)
	}else if fn == "setState" {
		return t.setState(stub,args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

