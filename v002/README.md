## MiniBC区块链V002 - KV数据库的实现和区块链数据的持久化 #    

####  KV数据库

虽然我们已经创建了一条非常简单的区块链，但是当程序关闭后，内存中的区块数据却没有保存下来。这使得我们无法重复使用一个区块链，也无法与其他人分享，我们需要将它存储在硬盘中。我们接下来的任务就是实现一个极小的KV数据库，用来保存区块链数据。像比特币bitcoin使用了leveldb数据库，也有些golang开发的区块链采用了BoltDB，它们都属于单机KV数据库。


KV数据库，也就是key/value数据库，这种数据库没有关系型数据库系统RDMBS（比如MySQL，Oracle，PostgreSQL等）的表（table）、行（row）、列（column）等概念。数据均以键/值（key/value）的方式进行存储，类似Go语言中的Map数据结构，只不过Map是存放在内存中，而KV数据库是存放在硬盘文件中。我们实现的简易KV数据库，也引入了类似表（实际上一般称为桶）的概念，将一组类似的KV键值数据进行分组。要想获得一个value，你需要知道对应的table和key。这有点像二维的map[string][string]结构，第一维是表(table），第二维是Key。其实我们的kv数据库确实就是按照二维的map实现的，我们实现KV数据库就是解决二维Map如何存储到硬盘文件中。


#### 数据编解码encoding/gob

为了让某个数据结构能够在网络上传输或者保存到文件里，我们必须对这个数据结构进行编码，同时要保证编码的数据能够被解码，还原为原始的数据结构。当然，已经有许多可用的编码方式了：JSON，XML，Google 的 protocol buffers等等。在golang中，我们又多了一种，这就是由encoding/gob包提供的编解码方式。 gob是Golang包自带的一个数据结构序列化的编码/解码工具。 编码使用Encoder，解码使用Decoder。在关闭数据库的时候，我们使用Encoder对map数据结构进行编码，保存到文件中；在打开数据库的时候，再使用Decoder将序列化的数据转换成map数据结构。
  
关于encoding/gob的使用大家可自行揣摩，encoding包下的xml、json用于网络数据交换、rpc等场景，使用比较频繁，大家可以多练习一下。
encoding/gob的使用比较简单，比如编码：

		//创建一个字节缓冲区buffer，作为参数初始化一个新的编码器encoder，
		//然后就可以对数据机构struct进行编码Encode(struct)，返回二进制字节流[]byte
		buffer := new(bytes.Buffer)
		encoder := gob.NewEncoder(buffer)
		encoder.Encode(Database.data)
  
我们要实现的数据库，保存数据到文件的时候需要用到encoding/gob编码，从文件读取数据的时候需要用到encoding/gob解码。


####  实现数据库Database

Database 实现了一个简单的kv数据库

	type Database struct {
		data map[string]map[string][]byte
	}

比如：我要在人员表User中存放一些人的信息，比如有张三、李四和王五，那么可以这样存放数据：

	map["User"]["张三"] = 张三的信息，比如姓名+年龄+性别等  
	map["User"]["李四"] = 李四的信息，比如姓名+年龄+性别等
	map["User"]["王五"] = 王五的信息，比如姓名+年龄+性别等
	
我们把一个人的数据信息，比如姓名、年龄和性别等合成为一个value，存放到Map中，通过表名Table和键Key就能找到。


Database的操作也非常简单，只有两个操作方法Get和Set，分别用于读取数据和存放数据。

从Database中读取数据：

	//从表table的记录里读取主键为key的记录的值
	func (Database *Database) Get(table, key string) []byte {
		if row, ok := Database.data[table]; ok {
			if val, ok := row[key]; ok {
				return val
			}
		}

		//如果键值不存在，返回空数据
		return []byte{}
	}

向Database中存放数据：

	//向表table添加一条记录，主键为key, 值为val
	func (Database *Database) Set(table, key string, val []byte) {
		Database.data[table] = map[string][]byte{key: val}
	}

为了简化数据库的创建过程，提供了数据库的创建函数: 

	func NewDatabase() *Database {
		Database := new(Database)
		Database.open()
		return Database
	}

数据库使用完毕后，关闭数据库。这个过程中，需要将内存中的Database.data这个二维map转为二进制数据，写入硬盘文件中。

	func (Database *Database) Close() {

		//将Database.data数据序列化为二进制数据，等待写入数据库文件
		buffer := new(bytes.Buffer)
		encoder := gob.NewEncoder(buffer)
		err := encoder.Encode(Database.data)
		if err != nil {
			log.Panic(err)
		}

		//创建数据库文件，如果文件已经存在，那么直接覆盖掉
		f, err := os.Create(DATABASE_FILE)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()

		//将序列化数据写入文件
		_, err = f.Write(buffer.Bytes())
		if err != nil {
			log.Panic(err)
		}

		//并调用同步，把数据从缓冲区立即写盘
		f.Sync()
	}
  
这样，我们的数据库就可以开始工作了~~~~


####  用Database读取和存放区块Block

实现KV Database的主要目标就是存放区块数据。我们的Database只能保存[]byte数据，为了能够将Block结构存入数据库，就需要先对Block结构进行编码，因此我们为Block增加序列化方法Serialize：

	//区块序列化，也就是将区块结构的内部数据转换为可以存储的字节流的格式
	func (block *Block) Serialize() []byte {
		buffer := new(bytes.Buffer)
		encoder := gob.NewEncoder(buffer)
		err := encoder.Encode(block)
		if err != nil {
			log.Panic(err)
		}

		return buff.Bytes()
	}

同样的，为了能够读取数据库中已被序列化的数据，就需要对数据进行解码，还原为Block结构，因此我们为Block增加了反序列化方法Deserialize：

	//区块反序列化，也就是将字节流转换为含有内部数据的区块结构，这个过程跟Serialize正好相反
	func Deserialize(bytesBlock []byte) (*Block) {
		decoder := gob.NewDecoder(bytes.NewReader(bytesBlock))
		var block Block
		err := decoder.Decode(&block)
		if err != nil {
			log.Panic(err)
		}

		return &block
	}


####  改造BlockChain
前面我们讨论过，区块链从技术角度来讲，就是一个链表结构，最初我们为了简单起见，把区块链定义为了一个切片：

	type BlockChain struct {
		Blocks []*Block
	}

现在我们将利用KV Database来实现链表结构，因为BlockChain的定义修改为：

	type BlockChain struct {
		DB *Database
	}

首先，我们在数据库中，为每一个区块Block建立一个键值对，其中Key存放该区块Block的hash值，Value存放区块Block数据;
其次，我们单独存放最后一个区块的hash值，建立一个键值对，其中Key为固定名称“last”，Value存放最后一个区块的hash值。
  
遍历区块链的顺序是由后向前，流程如下：
先取出最后一个区块的Hash值，通过Hash值就能够从数据库中取得它存放在数据库中的Value，Value经过解码后还原为Block结构;
再从Block中取得前一个区块的哈希值HashPrevBlock，继而从数据库获得前一个区块；
重复这个过程，直至创世纪区块，我们就能得到区块链上所有区块的数据。


为了遍历区块链，实现以上的流程，我们专门定义了一个迭代子结构BlockchainIterator：

	type BlockchainIterator struct {
		//当前区块的hash
		hashCurrent []byte

		//区块链数据库
		DB *Database
	}

迭代子BlockchainIterator只要一个方法Next：

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


遍历区块链的代码如下：

		iterator := blockchain.Iterator()
		for {
			block := iterator.Next()
			if block == nil {
				break
			}
		}


#### 运行
  

让我们看看MiniBC能否正常工作：

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

运行后输出：  

	前一区块哈希值： 3a57d1664c26d4d1cedd5ccb430143dd620f93f5f198ec39f1c3492461cb4eb7
	当前区块内容为： Mini block 03
	当前区块哈希值： d8472bf7b9004fa9959404e1e769ffebab2fc94f5a47fb58212d85f7530cb7af
	=============================================
	前一区块哈希值： 9fb890edd668a65f484dab2756797f15a58ba5bf64a3b6d0a1360b6aba530214
	当前区块内容为： Mini block 02
	当前区块哈希值： 3a57d1664c26d4d1cedd5ccb430143dd620f93f5f198ec39f1c3492461cb4eb7
	=============================================
	前一区块哈希值： 462442c457bed5f8d04a24bc4b95e84cee8026e0338a0a2ea5b92e06a2243ef0
	当前区块内容为： Mini block 01
	当前区块哈希值： 9fb890edd668a65f484dab2756797f15a58ba5bf64a3b6d0a1360b6aba530214
	=============================================
	前一区块哈希值： 
	当前区块内容为： Genesis Block
	当前区块哈希值： 462442c457bed5f8d04a24bc4b95e84cee8026e0338a0a2ea5b92e06a2243ef0
	=============================================


好像运行的还不错。

#### 我们实现的和将要实现的  

我们已经构建了一个非常简单的区块链，还实现了一个简易的KV数据库，区块链数据已经可以保存到数据库。下一节，我们的目标是构建一个可人机交互的区块链以及区块链浏览器。  


#### 交流和疑问  

大家可以从v001开始逐步深入的研究学习，每个源码文件均有完整的说明和注释。如果有疑问、建议和不明白的问题，尽可以与我联系。

MiniBC区块链交流可以加入qq群：<font size=5><b>777804802</b></font>，go开发者乐园
我的微信：<font size=5><b>bkra50 </b></font>  
github addr: https://github.com/wangshizebin/minibc

另外，招募一群志同道合者，共同维护项目，共同推进项目，为开源世界贡献力量，欢迎勾搭。

Let's go~~~~~~~~~~~~~~~~~~~~~~~~~~~~

