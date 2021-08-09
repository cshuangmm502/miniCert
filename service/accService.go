package service

import "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

func (t *ServiceSetup) QueryCert() (string,error){
	eventID := "eventIssueCert"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "getState", Args: [][]byte{[]byte("issuedState"), []byte(eventID)}}
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

//func (t *ServiceSetup) RevokedCert(state Accumulator) (string,error){
//
//}