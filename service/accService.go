package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"math/big"
	"miniCert/utils"
)

func (t *ServiceSetup) ValidateCert() ([]byte,error){

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "getState", Args: [][]byte{[]byte("issuedState")}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	if respone.Payload==nil{
		return []byte("kong"),nil
	}
	return respone.Payload, nil
}

func (t *ServiceSetup) QueryIssuedAccState() ([]byte,error){
	//eventID := "eventIssuedAccQuery"
	//reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	//defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{
		ChaincodeID:     t.ChaincodeID,
		Fcn:             "queryState",
		Args:            [][]byte{[]byte("issuedState")},
	}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return []byte{0x00}, err
	}

	//err = eventResult(notifier, eventID)
	//if err != nil {
	//	return []byte{0x00}, err
	//}

	return respone.Payload, nil
}

func (t *ServiceSetup) QueryRevokedAccState() ([]byte,error){
	eventID := "eventRevokedAccQuery"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{
		ChaincodeID:     t.ChaincodeID,
		Fcn:             "queryState",
		Args:            [][]byte{[]byte("revokedState")},
	}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return []byte{0x00}, err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceSetup) UpdateIssuedAccState(newState []byte) (string,error){
	eventID := "eventIssuedAccUpdate"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{
		ChaincodeID:     t.ChaincodeID,
		Fcn:             "updateState",
		Args:            [][]byte{[]byte("issuedState"),newState,[]byte(eventID)},
	}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) UpdateRevokedAccState() (string,error){
	eventID := "eventRevokedAccUpdate"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{
		ChaincodeID:     t.ChaincodeID,
		Fcn:             "getState",
		Args:            [][]byte{[]byte("issuedState")},
	}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}


func (t* ServiceSetup) IssueCert() (string,error){
	serial++
	certName := utils.CreateCertificate(serial)
	certStr := utils.ReadCertFromFile(certName)
	primeOfCert,_ := utils.HashToPrime(certStr)
	currentState,err := t.QueryIssuedAccState()
	if err!=nil{
		fmt.Println(err.Error())
	}
	currentStateAsObj := &utils.Record{}
	json.Unmarshal([]byte(currentState),currentStateAsObj)
	currentA := new(big.Int)
	currentA,_ = currentA.SetString(currentStateAsObj.A,10)
	secretN := new(big.Int)
	secretN,_ = secretN.SetString(currentStateAsObj.N,10)
	newState := big.NewInt(0)
	fmt.Println("cuurentA:"+currentA.String())
	fmt.Println("prineOfCert:"+primeOfCert.String())
	fmt.Println("secretN:"+secretN.String())
	newState.Exp(currentA,primeOfCert,secretN)
	if err!=nil{
		return "",err
	}
	currentStateAsObj.A=newState.String()
	newStateAsBytes,_ := json.Marshal(currentStateAsObj)
	t.UpdateIssuedAccState(newStateAsBytes)

	fmt.Println(certName)
	return "",nil
}

//func (t *ServiceSetup) RevokedCert(state Accumulator) (string,error){
//
//}