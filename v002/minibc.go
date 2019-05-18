package main

import (
	"fmt"
)
func main() {

	//新建一条区块链，如果区块链数据库中已经有了数据，将会读取数据库加载进来
	//如果尚未创建数据库，或者数据库为空，那么会自动生成一个创世纪区块
	bc := NewBlockChain()
	defer bc.DB.Close()

	//如果区块链中只有1个创世纪区块，我们就再添加3个区块。
	if bc.Iterator().GetCount() == 1 {
		bc.AddBlock("Mini block 01")
		bc.AddBlock("Mini block 02")
		bc.AddBlock("Mini block 03")
	}

	//区块链中应该有4个区块：1个创世纪区块，还有3个添加的区块
	iterator := bc.Iterator()
	for {
		block := iterator.Next()
		if block == nil {
			break
		}
		fmt.Println("前一区块哈希值：", BytesToHex(block.HashPrevBlock))
		fmt.Println("当前区块内容为：", string(block.Data))
		fmt.Println("当前区块哈希值：", BytesToHex(block.GetHash()))
		fmt.Println("=============================================")
	}
}
