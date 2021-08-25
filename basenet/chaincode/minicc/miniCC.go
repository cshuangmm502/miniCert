package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

//func(t *miniCC) issueCertificate(stub shim.ChaincodeStubInterface,args []string) peer.Response {
//	stateAsBytes,err :=stub.GetState(ISSUESTATE)
//	if err != nil{
//		return shim.Error(err.Error())
//	}
//
//}

//peer chaincode invoke -C myc -n mycc -c '{"Args":["queryState","issuedState"]}'
func(t *miniCC) queryState(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting a key")
	}
	tag := args[0]
	recordAsBytes,err := stub.GetState(tag)
	if err != nil{
		return shim.Error(err.Error())
	}
	//record := &Record{}
	//json.Unmarshal(recordAsBytes,record)
	//state := record.A
	return shim.Success(recordAsBytes)
}

//peer chaincode invoke -C myc -n mycc -c '{"Args":["setState","issuedState","123456"]}'
func(t *miniCC) updateState(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 3 {
		return shim.Error("Incorrect arguments. Expecting a key,a value and a eventID")
	}
	tag := args[0]
	err := stub.PutState(tag, []byte(args[1]))
	if err != nil{
		return shim.Error("fail to update the accumulator state ")
	}

	err = stub.SetEvent(args[2],[]byte{})
	if err!= nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("success to update the state"))
}

