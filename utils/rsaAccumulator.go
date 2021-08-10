package utils

import (
	"fmt"
	"math/big"
)

const (
	RSA_KEY_SIZE = 12
	RSA_PRIME_SIZE = RSA_KEY_SIZE/2
)

type Pair struct {
	first *big.Int
	second *big.Int
}

func (p *Pair) GetFirst() *big.Int{
	return p.first
}

func (p *Pair) GetSecond() *big.Int {
	return p.second
}

func NewPair(bitLength int) *Pair{
	var p1 *big.Int
	var p2 *big.Int

	p1 = GenerateLargePrime(bitLength)
	p2 = GenerateLargePrime(bitLength)
	for p:=p1;p.String()==p2.String();p2=GenerateLargePrime(bitLength){
	}

	return &Pair{first:p1,second:p2}
}

type RSAAccumulator struct {
	data 	map[string]*big.Int			//["key":hashPrime]
	pair 	*Pair
	p		*big.Int
	q 		*big.Int
	n		*big.Int
	//random 	big.Int
	a0		*big.Int
	a		*big.Int
}

func (rsaObj *RSAAccumulator)GetP() *big.Int {
	return rsaObj.p
}

func (rsaObj *RSAAccumulator)GetQ() *big.Int {
	return rsaObj.q
}

func (rsaObj *RSAAccumulator)GetN() *big.Int {
	return rsaObj.n
}

func (rsaObj *RSAAccumulator)GetA() *big.Int {
	return rsaObj.a
}

func (rsaObj *RSAAccumulator)GetA0() *big.Int {
	return rsaObj.a0
}

func (rsaObj *RSAAccumulator)GetVal(bigInteger big.Int) *big.Int {
	return rsaObj.data[bigInteger.String()]
}

func (rsaObj *RSAAccumulator)AddMember(key string) *big.Int {
	_,ok := rsaObj.data[key]
	if ok{
		return rsaObj.a
	}
	hashPrime,_ :=HashToPrime(key)
	//fmt.Println(hashPrime)
	rsaObj.a.Exp(rsaObj.a,hashPrime,rsaObj.n)
	rsaObj.data[key]=hashPrime
	return rsaObj.a
}

func (rsaObj *RSAAccumulator)ProveMembership(key string) *big.Int {
	_,ok := rsaObj.data[key]
	if !ok{
		return nil
	}
	witness := rsaObj.iterateAndGetProductWithoutX(key)
	return witness.Exp(rsaObj.a0,witness,rsaObj.n)
}

func (rsaObj *RSAAccumulator)DeleteMember(bigInteger big.Int) *big.Int{
	return big.NewInt(0)
}

func (rsaObj *RSAAccumulator)VerifyMembership(key string,proof *big.Int) bool{
	hashPrime,_ := HashToPrime(key)
	return	doVerifyMembership(rsaObj.a,hashPrime,proof,rsaObj.n)
}


func doVerifyMembership(accumulatorState *big.Int,hashPrime *big.Int,proof *big.Int,n *big.Int) bool{
	result := big.NewInt(1)
	result.Exp(proof,hashPrime,n)
	fmt.Println("当前累加器状态",accumulatorState)
	fmt.Println("当前关键字hash",hashPrime)
	fmt.Println("当前关键字存在性证明",proof)
	fmt.Println("当前result",result)
	if result.Cmp(accumulatorState)==0{
		return true
	}
	return false
}

func (rsaObj *RSAAccumulator)iterateAndGetProductWithoutX(key string) *big.Int{
	result := big.NewInt(1)
	for k,v := range rsaObj.data{
		if k!=key{
			result.Mul(result,v)
		}
	}
	return result
}

func (rsaObj *RSAAccumulator)iterateAndGetProduct() *big.Int{
	result := big.NewInt(1)
	for _,v := range rsaObj.data{
		result.Mul(result,v)
	}
	return result
}

func (rsaObj *RSAAccumulator)getPair() *Pair {
	return rsaObj.pair
}


func New() *RSAAccumulator {
	data := make(map[string]*big.Int)
	pair := NewPair(RSA_PRIME_SIZE)
	var N = new(big.Int)
	N.Mul(pair.GetFirst(), pair.GetSecond())
	random := GenerateRandomNumber(*big.NewInt(0), *N)
	random2 := big.NewInt(0)
	random2.Set(random)
	return &RSAAccumulator{
		data: data,
		pair: pair,
		p:    pair.GetFirst(),
		q:    pair.GetSecond(),
		n:    N,
		a:    random,
		a0:   random2,
	}
}

type Record struct {
	A	string
	N 	string
}