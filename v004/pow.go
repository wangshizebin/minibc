package main

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

const (
	//默认的工作量难度目标
	DIFFICULTY_TARGET_BITS = 242
)

//Pow表示一个工作量算法结构
type Pow struct {
	//要进行工作量计算的区块
	block *Block

	//工作量难度目标
	bits int32
}

//指定难度目标bits和随机数nonce，计算区块的hash值
func (pow *Pow) getHash(bits, nonce int32) []byte {

	//将Block整型成员全部转换成[]byte
	version := IntToByte(pow.block.Version)
	time := IntToByte(pow.block.Time)
	bbits := IntToByte(bits)
	bnonce := IntToByte(nonce)

	//将Block全部成员连接成一个[]byte
	data := bytes.Join([][]byte{version, time, bbits, bnonce, pow.block.HashPrevBlock, pow.block.Data}, []byte{})

	//对data进行两次hash计算
	hash := sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])

	return hash[:]
}

//通过不断变换nonce值来计算区块的hash，找到符合给定难度目标的nonce值
func (pow *Pow) Run() {

	//将0x1左移pow.bits位生成一个256位的大整数target，这是真正的难度目标值
	target := big.NewInt(1)
	target.Lsh(target, uint(pow.bits))

	//通过不断变换nonce值来计算区块的hash，直至找到小于或等于target的nonce值
	var hashInt big.Int
	for nonce := int32(0); nonce < math.MaxInt32; nonce++ {
		hashInt.SetBytes(pow.getHash(pow.bits, nonce))
		if hashInt.Cmp(target) <= 0 {
			//找到合适的nonce值，将其写入block中
			pow.block.Nonce = nonce
			pow.block.Bits = pow.bits
			break
		}
	}
}

//创建一个工作量对象
func NewPow(block *Block) *Pow {
	pow := &Pow{}
	pow.block = block
	pow.bits = DIFFICULTY_TARGET_BITS
	return pow
}
