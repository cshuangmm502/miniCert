package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

const PRIME_CERTAINTY  = 5

func GenerateLargePrime(bitLength int) *big.Int {
	prime,err := rand.Prime(rand.Reader,bitLength)
	if err != nil {
		fmt.Println("素数生成出错:",err)
	}
	return prime
}

func CheckPrime(p *big.Int)bool{
	return p.ProbablyPrime(PRIME_CERTAINTY)
}

func HashToPrime(input string) (*big.Int,*big.Int){
	count := 0
	x := HashToLength(input)
	//fmt.Println(x)
	for{
		if CheckPrime(x)==true {
			break
		}
		x.Add(x,big.NewInt(1))
		count++
	}
	return x,big.NewInt(int64(count));
}

func HashToPrimeWithNonce(input string,nonce int)*big.Int{
	val := big.NewInt(int64(nonce))
	return val.Add(val,HashToLength(input))
}

func HashToLength(x string) *big.Int {
	var randomHexString string
	randomHexString = x
	hash := sha256.New()
	hash.Write([]byte(randomHexString))
	hashResult := hash.Sum(nil)
	//32个字节表示
	//fmt.Println(hashResult)
	hashcode := hex.EncodeToString(hashResult)
	//64个十六进制数表示
	//fmt.Println(hashcode)
	n := new(big.Int)
	n,ok := n.SetString(hashcode,16)
	if !ok {
		fmt.Println("SetString: error")
		return big.NewInt(0)
	}
	//十进制表示
	//fmt.Println(n)
	return n
}

func GenerateRandomNumber(min big.Int,max big.Int) *big.Int{
	temp := big.NewInt(0)
	temp.Sub(&max,&min)
	temp.Add(temp,big.NewInt(1))
	result,_ := rand.Int(rand.Reader,temp)
	return result.Add(result,&min)
}

func Exgcd(a big.Int,b big.Int,x *big.Int,y *big.Int) big.Int {
	if b.Cmp(&a)==1{
		return Exgcd(b,a,y,x)
	}
	if b.Cmp(big.NewInt(0))==0{

		x.Set(big.NewInt(1))
		y.Set(big.NewInt(0))
		return a
	}
	var x1 = new(big.Int)
	var temp,temp1,temp2 big.Int
	temp.Mod(&a,&b)
	var d = Exgcd(b,temp,x1,x)
	temp1.Mod(&a,&b)
	temp2.Mul(&temp1,x)
	y.Sub(x1,&temp2)
	return d
}

func Bezoute_Coefficients(a big.Int,b big.Int)(big.Int,big.Int){
	var x = new(big.Int)
	var y = new(big.Int)
	_ = Exgcd(a,b,x,y)
	return *x,*y
}
