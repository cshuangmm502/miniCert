package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"
	rd "math/rand"
)

const configFilePath = "./conf/"
const OutputFilePath = ""

func CertificateIssuance(){
	//解析根证书
	caFile, err := ioutil.ReadFile(configFilePath+"cacert.cert")
	if err != nil {
		return
	}
	caBlock, _ := pem.Decode(caFile)

	cert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return
	}

	//解析私钥
	keyFile, err := ioutil.ReadFile(configFilePath+"key.pem")
	if err != nil {
		return
	}
	keyBlock, _ := pem.Decode(keyFile)
	rootkey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		fmt.Println(err)
	}

	equiCer := &x509.Certificate{
		SerialNumber: big.NewInt(rd.Int63()), //证书序列号
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{"BJUT"},
			OrganizationalUnit: []string{"SE"},
			Province:           []string{"BeiJing"},
			CommonName:         "testuser",
			Locality:           []string{"BeiJing"},
		},
		NotBefore:             time.Now(),                  //证书有效期开始时间
		NotAfter:              time.Now().AddDate(1, 0, 0), //证书有效期结束时间
		//BasicConstraintsValid: true,                        //基本的有效性约束
		//IsCA:           false,                                                                  //是否是根证书
		//KeyUsage:       x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
		//EmailAddresses: []string{"test@test.com"},
		//IPAddresses:    []net.IP{net.ParseIP("0.0.0.0")},
	}

	//签发证书
	//生成公钥私钥对
	start := time.Now()
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	ca, err := x509.CreateCertificate(rand.Reader, equiCer, cert, &priKey.PublicKey, rootkey)
	if err != nil {
		return
	}

	//编码证书文件和私钥文件
	File1, err := os.Create("./conf/testCert.pem")
	defer File1.Close()
	if err != nil {
		fmt.Println(err)
	}
	caPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca,
	}
	pem.Encode(File1,caPem)

	File2, err := os.Create("testKey.key")
	defer File2.Close()
	if err != nil {
		fmt.Println(err)
	}
	buf := x509.MarshalPKCS1PrivateKey(priKey)
	keyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: buf,
	}
	pem.Encode(File2,keyPem)

	fmt.Println(time.Since(start))
}

func ExamplePkiDomain(Country []string,Organization []string,OrganizationalUnit []string,Province []string,CommonName string,Locality []string) pkix.Name{
	return pkix.Name{
		Country:            Country,
		Organization:       Organization,
		OrganizationalUnit: OrganizationalUnit,
		Locality:           Locality,
		Province:           Province,
		CommonName:         CommonName,
	}
}