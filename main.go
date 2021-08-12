package main
import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/pkg/errors"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"miniCert/service"
	"miniCert/utils"
	"time"

	//"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"os"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	lcpackager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	mb "github.com/hyperledger/fabric-protos-go/msp"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
)


const configFile = "./conf/sdkconf.yaml"

var (
	// peers = []string{"peer0.org1.hauturier.com", "peer1.org1.hauturier.com", "peer0.org2.hauturier.com", "peer1.org2.hauturier.com"}
	peersorg1 = []string{"peer0.org1.hauturier.com", "peer1.org1.hauturier.com"}
	peersorg2 = []string{"peer0.org2.hauturier.com", "peer1.org2.hauturier.com"}
	peer1	         = "peer0.org1.hauturier.com"
	peer2		 = "peer0.org2.hauturier.com"
	orderers = []string{"orderer1.hauturier.com","orderer2.hauturier.com","orderer3.hauturier.com"}
)


func main(){
	sdkInfo := &SDKInfo{
		ChannelID:       "mychannel",
		ChannelConfig:   os.Getenv("GOPATH")+"/src/github.com/hauturier.com/miniCert/basenet/channel-artifacts/channel.tx",
		OrgAdmin:        "Admin",
		OrgName:         []string{"Org1","Org2"},
		OrdererOrgName:  "orderer1.hauturier.com",
		ChaincodePath:   os.Getenv("GOPATH")+"/src/github.com/hauturier.com/miniCert/basenet/chaincode/minicc",
		ChaincodeID:     "sacc",
		ChaincodeSequence:	1,
		UserName:        "User1",
	}
	sdkInfo.setSDK()

	err :=sdkInfo.initContext()
	if err != nil{
		fmt.Println(err)
	}

	err = createChannel(sdkInfo)
	if err != nil{
		fmt.Println(err)
	}
	//packageID: sacc:3b352ec80e8e7d3ff6df07b556cc54f6a1827abbbae70f376e74ff9d42994ecf
	//新生命周期测试
	label, ccPkg := packageCC()
	packageID := lcpackager.ComputePackageID(label, ccPkg)
	fmt.Println(packageID)
	installCC(label, ccPkg, sdkInfo)
	approveCC(packageID, sdkInfo)
	//queryApprovedCC(sdkInfo)
	commitCC(sdkInfo)
	////
	initCC(sdkInfo)
	serviceSetup := service.ServiceSetup{
		ChaincodeID: "sacc",
		Client:      sdkInfo.channelClient,
	}
	time.Sleep(10 * time.Second)

	msg,err := serviceSetup.ValidateCert()
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("查询最新发布状态成功 ")
	}
	fmt.Println(string(msg))
	//启动web服务
	//app := controller.Application{Setup:&serviceSetup}
	//web.WebStart(app)
	//启动web服务
}


func (sdkInfo *SDKInfo)setSDK(){
	sdk,err := fabsdk.New(config.FromFile(configFile))
	if err!=nil {
		fmt.Println("Fabric SDK实例化失败")
	}
	fmt.Println("Fabric SDK实例化成功")
	sdkInfo.SDK = sdk
}

func createChannel(info *SDKInfo) error{
	fmt.Println("**********")
	org1MspClient := info.org1MspClient
	adminIdentity,err := org1MspClient.GetSigningIdentity(info.OrgAdmin)
	if err !=nil {
		return fmt.Errorf("failed to get the signed identity of the special ID: %v",err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID:info.ChannelID,ChannelConfigPath:info.ChannelConfig,SigningIdentities:[]msp.SigningIdentity{adminIdentity}}
	org1ResMgmt := info.org1ResMgmt
	_,err = org1ResMgmt.SaveChannel(req,resmgmt.WithRetry(retry.DefaultResMgmtOpts),resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err !=nil {
		return fmt.Errorf("failed to create application channel: %v",err)
	}
	fmt.Println("success to create application channel")

	err = info.org1ResMgmt.JoinChannel(info.ChannelID,resmgmt.WithRetry(retry.DefaultResMgmtOpts),resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err!=nil {
		return fmt.Errorf("failed to join channel%v",err)
	}
	fmt.Println("peers success to join channel")

	clientChannelContext := info.SDK.ChannelContext(info.ChannelID,fabsdk.WithUser(info.UserName),fabsdk.WithOrg("Org1"))
	channelClient,err := channel.New(clientChannelContext)
	if err != nil {
		return fmt.Errorf("创建应用通道客户端失败: %v", err)
	}
	info.channelClient = channelClient
	return nil
}

//2.0新生命周期形式
func packageCC() (string, []byte) {
	desc := &lcpackager.Descriptor{
		Path:  os.Getenv("GOPATH")+"/src/github.com/hauturier.com/miniCert/basenet/chaincode/minicc",
		Type:  pb.ChaincodeSpec_GOLANG,
		Label: "sacc",
	}
	ccPkg, err := lcpackager.NewCCPackage(desc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("1")
	return desc.Label, ccPkg
}

func installCC(label string, ccPkg []byte, sdkInfo *SDKInfo) {
	installCCReq := resmgmt.LifecycleInstallCCRequest{
		Label:   label,
		Package: ccPkg,
	}
	for _, peer := range peersorg1 {
		resp1, err := sdkInfo.org1ResMgmt.LifecycleInstallCC(installCCReq, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			fmt.Println("2 err: ", peer, err)
		}
		fmt.Println("2: ", peer, resp1)
	}
	//for _, peer := range peersorg2 {
	//	resp2, err := org2RMgmt.LifecycleInstallCC(installCCReq, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	//	if err != nil {
	//		fmt.Println("2 err: ", peer, err)
	//	}
	//	fmt.Println("2: ", peer, resp2)
	//}

	fmt.Println("2 结束")
}

func approveCC(packageID string,sdkInfo *SDKInfo) {
	org1peers,err := DiscoverLocalPeers(sdkInfo.org1AdminClientContext,2)
	if err != nil {
		fmt.Println("333333")
		fmt.Println("3 err",err)
	}
	//org2peers,err := DiscoverLocalPeers(sdkInfo.org2AdminClientContext,2)
	//if err != nil {
	//	fmt.Println("4 err",err)
	//}

	ccPolicy := policydsl.SignedByNOutOfGivenRole(1, mb.MSPRole_MEMBER, []string{"Org1MSP", "Org2MSP"})
	approveCCReq := resmgmt.LifecycleApproveCCRequest{
		Name:              sdkInfo.ChaincodeID,
		Version:           "0",
		PackageID:         packageID,
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}

	resq1, err := sdkInfo.org1ResMgmt.LifecycleApproveCC(sdkInfo.ChannelID, approveCCReq, resmgmt.WithTargets(org1peers...), resmgmt.WithOrdererEndpoint("orderer1.hauturier.com"), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println("322222")
		fmt.Println("3 err: ", peersorg1[0], err)
	}
	fmt.Println("3: ", peersorg1[0], resq1)

	//resq2, err := sdkInfo.Org2ResMgmt.LifecycleApproveCC(ChannelID, approveCCReq, resmgmt.WithTargets(org2peers...), resmgmt.WithOrdererEndpoint("orderer.hauturier.com"), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	//if err != nil {
	//	fmt.Println("3 err: ", peersorg2[0], err)
	//}
	//fmt.Println("3: ", peersorg2[0], resq2)

	//fmt.Println("3 结束")

	//ccPolicy := policydsl.SignedByNOutOfGivenRole(2, mb.MSPRole_MEMBER, []string{"Org1MSP", "Org2MSP"})
	//approveCCReq := resmgmt.LifecycleApproveCCRequest{
	//	Name:              CCID,
	//	Version:           "0",
	//	PackageID:         packageID,
	//	Sequence:          1,
	//	EndorsementPlugin: "escc",
	//	ValidationPlugin:  "vscc",
	//	SignaturePolicy:   ccPolicy,
	//	InitRequired:      true,
	//}
	//
	//_, err := sdkInfo.Org1ResMgmt.LifecycleApproveCC(ChannelID, approveCCReq, resmgmt.WithTargetEndpoints(peer1), resmgmt.WithOrdererEndpoint("orderer.hauturier.com"), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//resq, err := sdkInfo.Org2ResMgmt.LifecycleApproveCC(ChannelID, approveCCReq, resmgmt.WithTargetEndpoints(peer2), resmgmt.WithOrdererEndpoint("orderer.hauturier.com"), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(resq)
	//fmt.Printf("3")
}

func queryApprovedCC(sdkInfo *SDKInfo){
	org1peers,err := DiscoverLocalPeers(sdkInfo.org1AdminClientContext,2)
	if err != nil {
		fmt.Println("4 err",err)
	}
	//org2peers,err := DiscoverLocalPeers(sdkInfo.org2AdminClientContext,2)
	if err != nil {
		fmt.Println("4 err",err)
	}
	queryApprovedCCReq := resmgmt.LifecycleQueryApprovedCCRequest{
		Name:     sdkInfo.ChannelID,
		Sequence: 1,
	}

	for _, p := range org1peers {
		fmt.Println(p)
			resp, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
				func() (interface{}, error) {
					resp1, err := sdkInfo.org1ResMgmt.LifecycleQueryApprovedCC(sdkInfo.ChannelID, queryApprovedCCReq, resmgmt.WithTargets(p))
					if err != nil {
						return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("LifecycleQueryApprovedCC returned error: %v", err), nil)
					}
					return resp1, err
				},
			)
			if err != nil {
				fmt.Println("4 err",err)
			}
			fmt.Println("4 ",resp)
	}

	//for _, p := range org2peers {
	//	resp, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
	//		func() (interface{}, error) {
	//			resp1, err := sdkInfo.Org2ResMgmt.LifecycleQueryApprovedCC(ChannelID, queryApprovedCCReq, resmgmt.WithTargets(p), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	//			if err != nil {
	//				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("LifecycleQueryApprovedCC returned error: %v", err), nil)
	//			}
	//			return resp1, err
	//		},
	//	)
	//	if err != nil {
	//		fmt.Println("4 err",err)
	//	}
	//	fmt.Println("4 ",resp)
	//}
	//fmt.Println("4 结束")
}

func commitCC(sdkInfo *SDKInfo) {
	//ccPolicy := policydsl.SignedByNOutOfGivenRole(2, mb.MSPRole_MEMBER, []string{"Org1MSP", "Org2MSP"})
	//req := resmgmt.LifecycleCommitCCRequest{
	//	Name:              CCID,
	//	Version:           "0",
	//	Sequence:          Sequence,
	//	EndorsementPlugin: "escc",
	//	ValidationPlugin:  "vscc",
	//	SignaturePolicy:   ccPolicy,
	//	InitRequired:      true,
	//}
	//_, err := org1RMgmt.LifecycleCommitCC(ChannelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.hauturier.com"))
	//if err != nil {
	//	fmt.Println("5 org1RMgmt err: ", err)
	//	// req.Sequence = req.Sequence + 1
	//	// _, err = org1RMgmt.LifecycleCommitCC(channelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.hauturier.com"))
	//}
	//// fmt.Println("4 : ", resp1)
	//req.Sequence = Sequence + 1
	//_, err = org2RMgmt.LifecycleCommitCC(ChannelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.hauturier.com"))
	//if err != nil {
	//	fmt.Println("5 org2RMgmt err: ", err)
	//	// req.Sequence = req.Sequence + 1
	//	// _, err = org2RMgmt.LifecycleCommitCC(channelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.hauturier.com"))
	//}
	//// fmt.Println("4 : ", resp2)
	//fmt.Println("5 : 结束")

	ccPolicy := policydsl.SignedByNOutOfGivenRole(1, mb.MSPRole_MEMBER, []string{"Org1MSP", "Org2MSP"})
	req := resmgmt.LifecycleCommitCCRequest{
		Name:              sdkInfo.ChaincodeID,
		Version:           "0",
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}
	_, err := sdkInfo.org1ResMgmt.LifecycleCommitCC(sdkInfo.ChannelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer1.hauturier.com"))
	if err != nil {
		fmt.Println("5 org1RMgmt err: ", err)
	}
	fmt.Println("5 : 结束")
}

//'{"Args":["init","issuedState","accStateAsBytes"]}'
func initCC(sdkInfo *SDKInfo){
	accumulator := utils.New()
	accstate := &utils.Record{
		A: accumulator.GetA().String(),
		N: accumulator.GetN().String(),
	}
	accstateAsBytes,_ := json.Marshal(accstate)
	channelClient := sdkInfo.channelClient
	_,err := channelClient.Execute(channel.Request{ChaincodeID:sdkInfo.ChaincodeID,Fcn:"init",Args:[][]byte{[]byte("issuedState"), accstateAsBytes},IsInit:true},channel.WithRetry(retry.DefaultResMgmtOpts))
	if err!= nil {
		fmt.Println("链码sacc初始化累加器失败")
		fmt.Println(err.Error())
	}
	fmt.Println("初始化链码sacc成功")
}

func(sdkInfo *SDKInfo) initContext() error {
	//clientChannelContext := sdkInfo.SDK.ChannelContext(sdkInfo.ChannelID,fabsdk.WithUser(sdkInfo.UserName),fabsdk.WithOrg("Org1"))
	//channelClient,err := channel.New(clientChannelContext)
	//if err != nil {
	//	return fmt.Errorf("创建应用通道客户端失败: %v", err)
	//}
	//sdkInfo.channelClient = channelClient

	clientContext1 := sdkInfo.SDK.Context(fabsdk.WithUser(sdkInfo.OrgAdmin), fabsdk.WithOrg("Org1"))
	clientContext2 := sdkInfo.SDK.Context(fabsdk.WithUser(sdkInfo.OrgAdmin), fabsdk.WithOrg("Org2"))
	if clientContext1 == nil || clientContext2 == nil{
		return fmt.Errorf("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	}

	sdkInfo.org1AdminClientContext = clientContext1
	//sdkInfo.org2AdminClientContext = clientContext2

	org1RMgmt, err := resmgmt.New(clientContext1)
	if err != nil {
		fmt.Println("Failed to create org1RMgmt:", err)
	}
	//org2RMgmt, err := resmgmt.New(clientContext2)
	//if err != nil {
	//	fmt.Println("Failed to create org2RMgmt:", err)
	//}

	sdkInfo.org1ResMgmt = org1RMgmt
	//sdkInfo.Org2ResMgmt = org2RMgmt

	org1MspClient,err :=mspclient.New(sdkInfo.SDK.Context(),mspclient.WithOrg("Org1"))
	if err != nil {
		fmt.Println("Failed to create org1MspClient:", err)
	}
	//org2MspClient,err :=mspclient.New(sdkInfo.SDK.Context(),mspclient.WithOrg(Org2Name))
	//if err != nil {
	//	fmt.Println("Failed to create org2MspClient:", err)
	//}
	sdkInfo.org1MspClient = org1MspClient
	//sdkInfo.org2MspClient = org2MspClient
	fmt.Println("sdk 资源初始化成功")
	return nil
}

// DiscoverLocalPeers queries the local peers for the given MSP context and returns all of the peers. If
// the number of peers does not match the expected number then an error is returned.
func DiscoverLocalPeers(ctxProvider contextAPI.ClientProvider, expectedPeers int) ([]fabAPI.Peer, error) {
	ctx, err := contextImpl.NewLocal(ctxProvider)
	if err != nil {
		return nil, errors.Wrap(err, "error creating local context")
	}

	discoveredPeers, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			peers, serviceErr := ctx.LocalDiscoveryService().GetPeers()
			if serviceErr != nil {
				return nil, errors.Wrapf(serviceErr, "error getting peers for MSP [%s]", ctx.Identifier().MSPID)
			}
			if len(peers) < expectedPeers {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Expecting %d peers but got %d", expectedPeers, len(peers)), nil)
			}
			return peers, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return discoveredPeers.([]fabAPI.Peer), nil
}

