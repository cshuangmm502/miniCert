package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"math/big"
)

func create_all_membership_witness(A0 *big.Int,set map[string]int,N *big.Int)[]*big.Int{
	var primes []*big.Int
	for k,v := range set{
		prime := HashToPrimeWithNonce(k,v)
		primes=append(primes, prime)
		fmt.Println(k,prime)
	}

	fmt.Println(primes)
	return root_factor(A0,primes,N)
}

func root_factor(g *big.Int,primes []*big.Int,N *big.Int)[]*big.Int{
	n := len(primes)
	if n==1{
		var result []*big.Int
		result = append(result,g)
		return result
	}

	n_tag := n/2

	primes_L := primes[n_tag:n]
	product_L := calculate_product(primes_L)
	g_L := big.NewInt(1)
	g_L.Exp(g,product_L,N)

	primes_R := primes[0:n_tag]
	product_R := calculate_product(primes_R)
	g_R := big.NewInt(1)
	g_R.Exp(g,product_R,N)

	L := root_factor(g_L, primes_R,N)
	R := root_factor(g_R, primes_L,N)

	var result []*big.Int
	result = append(result, L...)
	result = append(result, R...)
	return result
}

func calculate_product(list []*big.Int)*big.Int{
	base := big.NewInt(1)
	for _,i := range list{
		base.Mul(base,i)
	}
	return base
}

//func(t *miniCC) issueCertificate(stub shim.ChaincodeStubInterface,args []string) peer.Response {
//	stateAsBytes,err :=stub.GetState(ISSUESTATE)
//	if err != nil{
//		return shim.Error(err.Error())
//	}
//
//}

//peer chaincode invoke -C myc -n mycc -c '{"Args":["queryState","issuedState"]}'
//peer chaincode invoke -C myc -n mycc -c '{"function":"queryState","Args":["revokeState"]}'
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

//peer chaincode invoke -C myc -n mycc -c '{"Args":["getState","issuedState"]}'
func(t *miniCC) setState(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 3 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}
	tag := args[0]
	err := stub.PutState(tag, []byte(args[1]))
	if err != nil{
		return shim.Error("fail to update the accumulator state ")
	}
	return shim.Success([]byte("success to update the state"))
}

