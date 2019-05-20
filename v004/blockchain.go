package main

const (
	TABLE_BLOCKS = "blocks"
	BLOCK_LAST   = "last"
)

//BlockChain 表示区块链，每一条区块链中包含着多个区块
type BlockChain struct {
	DB *Database
}

type BlockchainIterator struct {
	//当前区块的hash
	hashCurrent []byte
	//区块链数据库
	DB *Database
}

//向区块链上增加一个区块
func (bc *BlockChain) AddBlock(data string) {

	//校验区块链上是否已经有了创世纪区块
	hashLast := bc.DB.Get(TABLE_BLOCKS, BLOCK_LAST)
	if len(hashLast) == 0 {
		//创建一个创世纪区块
		block := NewGenesisBlock()

		//取该区块的哈希值
		blockHash := block.GetHash()

		//将该区块的哈希值和序列化数据组成键值对存入数据库
		bc.DB.Set(TABLE_BLOCKS, string(blockHash), block.Serialize())

		//将最后一个区块的哈希值存入数据库，Key 标记为 "last"
		bc.DB.Set(TABLE_BLOCKS, BLOCK_LAST, blockHash)
		return
	}

	//取出当前区块链的最后一个区块
	val := bc.DB.Get(TABLE_BLOCKS, string(hashLast))
	prevBlock := Deserialize(val)

	//传入区块数据和最后一个区块的hash，建新的区块
	block := NewBlock(data, prevBlock.GetHash())

	//取该区块的哈希值
	blockHash := block.GetHash()

	//将该区块的哈希值和序列化数据组成键值对存入数据库
	bc.DB.Set(TABLE_BLOCKS, string(blockHash), block.Serialize())

	//将最后一个区块的哈希值存入数据库，Key 标记为 "last"
	bc.DB.Set(TABLE_BLOCKS, BLOCK_LAST, blockHash)
}

//向区块链上增加创世纪区块
func (bc *BlockChain) AddGenesisBlock() {

	//检查区块链上是否已经存在区块
	hashLast := bc.DB.Get(TABLE_BLOCKS, BLOCK_LAST)
	if len(hashLast) > 0 {
		//区块链上已经有区块了，不能添加创世纪区块
		return
	}

	//创建一个创世纪区块
	block := NewGenesisBlock()

	//取该区块的哈希值
	blockHash := block.GetHash()

	//将该区块的哈希值和序列化数据组成键值对存入数据库
	bc.DB.Set(TABLE_BLOCKS, string(blockHash), block.Serialize())

	//将最后一个区块的哈希值存入数据库，Key 标记为 "last"
	bc.DB.Set(TABLE_BLOCKS, BLOCK_LAST, blockHash)
}

// 获得遍历区块链的迭代子
func (bc *BlockChain) Iterator() *BlockchainIterator {
	hashLast := bc.DB.Get(TABLE_BLOCKS, BLOCK_LAST)
	if len(hashLast) == 0 {
		//数据库中没有最后区块的记录，说明还没生成区块，无需迭代
		return nil
	}

	//取出最后一个区块的hash，初始化迭代子，准备迭代
	it := &BlockchainIterator{
			hashCurrent:hashLast,
			DB: bc.DB,
	}

	return it
}

func (it *BlockchainIterator) Next() *Block {
	if it.hashCurrent == nil {
		//已经没有前一区块，到此结束
		return nil
	}

	//取出当前遍历到的区块
	val := it.DB.Get(TABLE_BLOCKS, string(it.hashCurrent))
	block := Deserialize(val)
	it.hashCurrent = block.HashPrevBlock

	return block
}

//统计区块链中区块的数量
func (it *BlockchainIterator) GetCount() int {
	var count = 0
	if it.hashCurrent == nil {
		return count
	}
	for {
		if it.Next() != nil {
			count++
			continue
		}
		break
	}
	return count
}

//新建一个区块链对象
func NewBlockChain() *BlockChain {
	blockChain := BlockChain{NewDatabase()}
	blockChain.AddGenesisBlock()
	return &blockChain
}
