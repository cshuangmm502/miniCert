package main

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type SDKInfo struct {
	SDK                    *fabsdk.FabricSDK

	ChannelID         string
	ChannelConfig     string
	OrgAdmin          string
	OrgName           []string
	OrdererOrgName    string
	ChaincodePath     string
	ChaincodeID       string
	ChaincodeSequence int
	UserName          string

	ordererClientContext   contextAPI.ClientProvider
	org1AdminClientContext contextAPI.ClientProvider
	org2AdminClientContext contextAPI.ClientProvider
	org1ResMgmt            *resmgmt.Client
	org2ResMgmt            *resmgmt.Client
	org1MspClient          *mspclient.Client
	org2MspClient          *mspclient.Client
	LedgerClient           *ledger.Client
	channelClient          *channel.Client

}
