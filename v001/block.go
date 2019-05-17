package main

import (
	"bytes"
	"crypto/sha256"
	"time"
)

const (
	VERSION = 1
)

//Block表示一个区块
type Block struct {
	Version       int32  //协议版本号
	HashPrevBlock []byte //上一个区块的hash值，长度为32个字节
	Time          int32  //时间戳，从1970.01.01 00:00:00到当前时间的秒数
	Bits          int32  //工作量证明(POW)的难度
	Nonce         int32  //要找的符合POW要求的的随机数

	Data []byte //区块存储的内容，在虚拟币中用来存储交易信息
}

//获得当前区块的hash值
func (block *Block) GetHash() []byte {

	version := IntToByte(block.Version)
	time := IntToByte(block.Time)
	bits := IntToByte(block.Bits)
	nonce := IntToByte(block.Nonce)

	data := bytes.Join([][]byte{version, block.HashPrevBlock, time, bits, nonce, block.Data}, []byte{})
	hash := sha256.Sum256(data)
	return hash[:]
}

//生成一个新区块，需要当前区块存储的内容和前一区块的hash两个参数
func NewBlock(data string, prevHash []byte) *Block {
	//Bits和Nonce暂时不赋值，后面引入了挖矿机制后再解决这个问题。
	block := Block{
		Version:       VERSION,
		HashPrevBlock: prevHash,
		Time:          int32(time.Now().Unix()),
		Bits:          0,
		Nonce:         0,
		Data:          []byte(data),
	}
	return &block
}

//生成创世纪区块，所谓创世纪区块是区块链的第一个区块。它存储的内容可以是任意内容，比如Genesis Block。
// 由于它没有前一区块，所以不需要提供前一区块的hash值。
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
