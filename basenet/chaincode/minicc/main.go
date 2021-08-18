/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type miniCC struct {
}

const ISSUESTATE = "issuedState"
const REVOKESTATE = "revokedState"


func main()  {
	if err := shim.Start(new(miniCC)); err != nil {
		fmt.Printf("Error starting miniCC chaincode: %s", err)
	}
}

//'{"Args":["init","issuedState","accStateAsBytes"]}'
func (t *miniCC) Init(stub shim.ChaincodeStubInterface) peer.Response{
	_,args := stub.GetFunctionAndParameters()
	if len(args)!=2{
		return shim.Error("初始化参数错误")
	}
	fmt.Println(" 初始化ing")
	tag := args[0]
	accRecordString := args[1]
	accRecordAsBytes := []byte(accRecordString)
	err := stub.PutState(tag,accRecordAsBytes)
	if err!=nil{
		return shim.Error(err.Error())
	}
	record := &Record{}
	_ = json.Unmarshal(accRecordAsBytes,record)
	return shim.Success([]byte("accumulator state:"+record.A))
}

func (t *miniCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	if fn == "getState" {
		return t.queryState(stub, args)
	}else if fn == "setState" {
		return t.setState(stub,args)
	}

	// Return the result as success payload
	return shim.Error("指定的函数名称错误")
}

