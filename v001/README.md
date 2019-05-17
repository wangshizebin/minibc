## MiniBC区块链V001 - 简单区块链的实现 #    

#### 区块

我们从“区块链”的“区块”部分开始。区块是区块链中最基本的数据结构，在区块链中，区块存储了有价值信息。例如，比特币区块存储了交易数据，除此之外，区块中还包含其它信息：版本号，当前时间和前一个块的哈希值等。我们把bitcoin的区块定义稍作简化，作为MiniBC的区块定义：  

	type Block struct {
		Version       int32  //协议版本号
		HashPrevBlock []byte //上一个区块的hash值，长度为32个字节
		Time          int32  //时间戳，从1970.01.01 00:00:00到当前时间的秒数
		Bits          int32  //工作量证明(POW)的难度
		Nonce         int32  //要找的符合POW要求的的随机数

		Data []byte //区块存储的内容，在虚拟币中用来存储交易信息
	}

其中：Version是版本号，类型是32位整数。Time是当前时间戳，也就是创建区块的时间，Bits和Nonce暂时不用考虑，后面在工作量证明(pow）部分再做详细讲解。Data是包含在区块中的实际最有价值信息，凡是区块链中要存储的内容全部包含在里面。PrevBlockHash存储前一个区块的哈希值（Hash）。在实际的区块链项目中，比如bitcoin（比特币），区块包含区块头和区块体两个部分，Data数据保存在去区块体，其它信息包含在区块头中。我们为了简单起见，将区块头和区块体合二为一。

那么我们如何计算哈希？哈希计算方法是区块链的一个非常重要的特点，正是这个特点使区块链变得安全。哈希计算比较耗时，像比特币挖矿的过程就是哈希计算的过程
，这也是人们购买计算能力更强大的GPU来挖掘比特币的原因，象比特大陆设计了专门的芯片和矿机以提高哈希计算的能力。关于引入哈希计算的目的，以后会详细讨论。  

现在，我们将获取区块中字段，把它们连接在一起，并在连接组合上计算SHA-256哈希值。我们在SetHash方法中这样做：  

	func (block *Block) GetHash() []byte {

		version := IntToByte(block.Version)
		time := IntToByte(block.Time)
		bits := IntToByte(block.Bits)
		nonce := IntToByte(block.Nonce)

		data := bytes.Join([][]byte{version, block.HashPrevBlock, time, bits, nonce, block.Data}, []byte{})
		hash := sha256.Sum256(data)
		return hash[:]
	}  

获得一个区块的哈市值非常简单，就是把区块头中的数据连接在一起，调用golang中"crypto/sha256"包函数，得到一个32字节(256位)长度的二进制数据，这就是该区块的哈希值。  
工作过程：首先把区块头一些整数项通过IntToByte方法加工为[]byte，然后调用bytes.Join方法将区块头的全部数据连接在一起，最后调用sha256.Sum256法法为输入的[]byte数据生成32字节的[]byte哈希值。

我们为了简化区块的创建过程，提供了区块创建函数:    

	func NewBlock(data string, prevHash []byte) *Block {
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

由于区块链中第一个区块没有前一个区块，所以区块结构中HashPrevBlock需要置为空数据[]byte{}。区块链中第一个区块通常称为“创世纪区块”，我们单独为创建创世纪区块提供了一个方法：  

	func NewGenesisBlock() *Block {
		return NewBlock("Genesis Block", []byte{})
	}

区块就这么简单！  

#### 区块链  

区块链本身的数据结构并不复杂，它是将数据块Block连接在一起的链表。所以我们首先要做的工作就是将区块链定义为一个链表，链表中每一个节点的内容是区块Block。链表对于软件开发人员来说，是最为熟悉的基本数据结构，通常的定义方式如下：

	type Object interface{} //节点中存储的数据的类型

	type Node struct {
		Data Object //节点存储的数据
		Next *Node //指向下一个节点的指针
	}

	type List struct {
		HeadNode *Node //头节点
	}

不熟悉链表的同学可以再回去补补链表的知识，包括链表的定义，遍历，我们的项目暂时还用不到，但后续的工作中会用到。为了简单起见，我们先简化区块链BlockChain的定义，将区块Block的组织方式暂定为切片（数组），切片既具备数据存储的能力，又具备顺序组织数据的能力。


	type Blockchain struct {
		blocks []*Block
	}

这是我们的第一个区块链！这也太简单了吧，实际上真正的区块链从总体上来看，也就是这么简单，但是要做到安全缜密，还需做大量的工作。  

现在可以随时在我们创建的区块链上添加一个区块： 
  
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


我们为了简化区块链的创建过程，提供了区块链的创建函数: 

	func NewBlockChain() *BlockChain {
		//预先创建一个创世纪区块
		blockChain := BlockChain{Blocks: []*Block{NewGenesisBlock()}}
		return &blockChain
	}

新的区块链要预先创建一个创世纪区块。

大功告成~~~~~~~~~~~~~~~~~~~~  


#### 运行

让我们看看MiniBC能否正常工作：    

	func main() {
		//新建一条区块链，里面隐含着创建了一个创世纪区块(初始区块)
		bc := NewBlockChain()

		//添加3个区块
		bc.AddBlock("Mini block 01")
		bc.AddBlock("Mini block 02")
		bc.AddBlock("Mini block 03")

		//区块链中应该有4个区块：1个创世纪区块，还有3个添加的区块
		for _, block := range bc.Blocks {
			fmt.Println("前一区块哈希值：", BytesToHex(block.HashPrevBlock))
			fmt.Println("当前区块内容为：", string(block.Data))
			fmt.Println("当前区块哈希值：", BytesToHex(block.GetHash()))
			fmt.Println("=============================================")
		}
	}  


运行后输出：  

	前一区块哈希值： 
	当前区块内容为： Genesis Block
	当前区块哈希值： 0e902e4108dfc0fb14b094ce3a130734bc5ef76f5aa58bc56a3635e031947e1e
	=============================================
	前一区块哈希值： 0e902e4108dfc0fb14b094ce3a130734bc5ef76f5aa58bc56a3635e031947e1e
	当前区块内容为： Mini block 01
	当前区块哈希值： 3a567c6151181dd88b6f1aa134e0b81151280c15d585df51898aa9fb5f37cd8b
	=============================================
	前一区块哈希值： 3a567c6151181dd88b6f1aa134e0b81151280c15d585df51898aa9fb5f37cd8b
	当前区块内容为： Mini block 02
	当前区块哈希值： c4145bc199da2d24d2b9d7390b8ae4b3cee3eb001b2dd57538b3142af8604fca
	=============================================
	前一区块哈希值： c4145bc199da2d24d2b9d7390b8ae4b3cee3eb001b2dd57538b3142af8604fca
	当前区块内容为： Mini block 03
	当前区块哈希值： 8cd906d8084fa2b2e83b4c3788d0b49dbd625e66268c3319c830282fcdaa5c05
	=============================================

好像运行的还不错。

#### 我们实现的和将要实现的  

我们已经构建了一个非常简单的区块链：区块链只是一个区块数组，但当前区块的HashPrevBlock已经指向前一个区块，可以利用这个字段实现区块的反向链接。当然实际的区块链要复杂很多。区块链的数据需要保存，我们通常存储数据，要么放在数据库里，要么放在文件中，我们MiniBC暂定为存放到数据库中。实际上bitcoin的区块数据存是放到文件里，其他信息存放到leveldb数据库中。下一节，我们的目标是构建一个简单的KV数据库以及持久化区块链。  


#### 交流和疑问
大家可以从v001开始逐步深入的研究学习，每个源码文件均有完整的说明和注释。如果有疑问、建议和不明白的问题，尽可以与我联系。

MiniBC区块链交流可以加入qq群：777804802，go开发者乐园 我的微信：bkra50 另外，招募一群志同道合者，共同维护项目，共同推进项目，为开源世界贡献力量，欢迎勾搭。

Let's go~~~~~~~~~~~~~~~~~~~~~~~~~~~~

