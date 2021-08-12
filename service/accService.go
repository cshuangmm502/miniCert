package service

import "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

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

func (t * ServiceSetup) UpdateIssuedAccState(accState []byte) (string,error){
	eventID := "eventIssuedAccUpdate"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{
		ChaincodeID:     t.ChaincodeID,
		Fcn:             "setState",
		Args:            [][]byte{[]byte("issuedState"),accState,[]byte(eventID)},
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
//func (t* ServiceSetup) IssueCert() (string,error){
//
//}

//func (t *ServiceSetup) RevokedCert(state Accumulator) (string,error){
//
//}