package main

//BlockChain 表示区块链，每一条区块链中包含着多个区块
type BlockChain struct {
	Blocks []*Block
}

//向区块链上增加一个区块
func (bc *BlockChain) AddBlock(data string) {
	if (bc.Blocks == nil || len(bc.Blocks) < 1) {
		return
	}

	//取出当前区块链的最后一个区块
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	//传入区块数据和最后一个区块的hash，建新的区块
	block := NewBlock(data, prevBlock.GetHash())
	bc.Blocks = append(bc.Blocks, block)
}

//新建一个区块链对象
func NewBlockChain() *BlockChain {
	//预先创建一个创世纪区块
	blockChain := BlockChain{Blocks: []*Block{NewGenesisBlock()}}

	return &blockChain
}
