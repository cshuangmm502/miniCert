package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
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

func (t *ServiceSetup) QueryIssuedAccState() (string,error){
	eventID := "eventIssuedAccQuery"
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

func (t *ServiceSetup) QueryRevokedAccState() (string,error){
	eventID := "eventRevokedAccQuery"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{
		ChaincodeID:     t.ChaincodeID,
		Fcn:             "getState",
		Args:            [][]byte{[]byte("revokedState")},
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

func (t *ServiceSetup) UpdateIssuedAccState() (string,error){
	eventID := "eventIssuedAccUpdate"
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
	fmt.Println(certName)
	return "",nil
}

//func (t *ServiceSetup) RevokedCert(state Accumulator) (string,error){
//
//}