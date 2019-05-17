package main

import (
	"fmt"
)

func main() {
	//新建一条区块链，里面隐含着创建了一个创世纪区块(初始区块)
	bc := NewBlockChain()

	//添加3个区块
	bc.AddBlock("Mini block 01")
	bc.AddBlock("Mini block 02")
	bc.AddBlock("Mini block 03")

	//区块链中应该有4个区块：1个创世纪区块，还有3个添加的区块
	for _, block := range bc.Blocks {
		fmt.Println("当前区块存储的内容：", string(block.Data), "前一区块得hash值：", BytesToHex(block.HashPrevBlock))
	}
}
