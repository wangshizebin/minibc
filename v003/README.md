## MiniBC区块链V003 - 区块链浏览器与人机交互 #    

####  工作目标

区块链浏览器是用户通过浏览器的方式查看区块链的所有信息。我们本节的目标就是实现这一功能。
我们不仅提供查看功能，还提供了了管理功能。在真实的区块链项目中，一般会提供多种交互方式，比如提供客户端命令行方式、websocket方式等等，最终由区块链server提供统一的rpc服务。我们目前先实现浏览器方式，以后会逐步扩充。

#### Http Server

在golang中写一个http server非常简单，三两行代码就可以实现。

	http.HandleFunc("/", handleIndex)
	http.ListenAndServe(":8080", nil)

我们的区块链浏览器就是使用了golang标准包里的http服务完成的。使用起来非常简单，调用NewBlockBrower创建一个区块链浏览器对象，然后调用Start方法就可以工作了：

	//启动区块链浏览器，您可以通过浏览器 http://SERVER_ADDR:8080 访问
	NewBlockBrower(bc).Start()


####  区块链浏览器

区块链浏览器的结构比较简单，一个通知退出的channel和当前区块链对象BlockChain：  

	type BlockBrower struct {
		//中止BrowserServer通道
		chanQuit chan bool

		//当前区块链
		blockChain *BlockChain
	}

启动区块链浏览器方法Start：
目前区块链服务器共能处理四种Url："/"、"/shutdown"、"/getblocks"和"/generateblock"，监听8080端口。
如果MiniBC运行后，如果通过本机器打开浏览器，可直接输入：http://localhost：8080 

	func (bb *BlockBrower) Start() {
		fmt.Println("=========================================================================")
		fmt.Println("MiniBC 区块链浏览器已经启动，请通过浏览器http://" + SERVER_ADDR + ":8080访问....")
		fmt.Println("=========================================================================")

		http.HandleFunc("/", bb.handleIndex)
		http.HandleFunc("/shutdown", bb.handleShutdown)
		http.HandleFunc("/getblocks", bb.handleGetBlocks)
		http.HandleFunc("/generateblock", bb.handleGenerateBlock)

		go http.ListenAndServe(":8080", nil)

		//只有接到退出通知，才能结束
		select {
			case <-bb.chanQuit:
		}
	}

实现了4个Url处理函数：

处理Url为 http://SERVER_ADDR:8080/ 的请求，显示区块链浏览器首页信息
	func (bb *BlockBrower) handleIndex(response http.ResponseWriter, request *http.Request) {
		content := "<html><br>"
		content = content + "<b>&nbsp;&nbsp;&nbsp;&nbsp;MiniBC区块链区块浏览器</b><br><br>"
		content = content + "&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"shutdown\">关闭MiniBC</a>"
		content = content + "&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"getblocks\">查看区块链</a>&nbsp;&nbsp;&nbsp;&nbsp;"
		content = content + "</html>"

		response.Write([]byte(content))
	}

处理Url为 http://SERVER_ADDR:8080/shutdown 的请求，关闭区块链浏览器，退出系统
	func (bb *BlockBrower) handleShutdown(response http.ResponseWriter, request *http.Request) {
		fmt.Println("")
		fmt.Println("")
		fmt.Println("=========================================================================")
		fmt.Println("MiniBC 区块链和远程管理结束，谢谢使用!")
		fmt.Println("=========================================================================")

		bb.chanQuit <- true
	}

处理Url为 http://SERVER_ADDR:8080/getblocks 的请求，打印所有区块链信息
	func (bb *BlockBrower) handleGetBlocks(response http.ResponseWriter, request *http.Request) {

		//获取区块链高度
		blockHeight := strconv.Itoa(bb.blockChain.Iterator().GetCount())

		content := "<html><br>"
		content = content + "<a href=\"/\">返回首页</a>&nbsp;&nbsp;&nbsp;&nbsp;"
		content = content + "<a href=\"generateblock\">生成新区块</a>&nbsp;&nbsp;&nbsp;&nbsp;<br>"
		content = content + "<br><b>&nbsp;&nbsp;&nbsp;&nbsp;当前MiniBC区块链高度：" + blockHeight + "</b><br><br>"

		//遍历区块链，打印每一个区块的详细信息
		iterator := bb.blockChain.Iterator()
		for {
			block := iterator.Next()
			if block == nil {
				break
			}
			content = content + "当前区块哈希值：0x" + BytesToHex(block.GetHash()) + "<br>"
			content = content + "当前区块内容为：" +  string(block.Data) + "<br>"
			content = content + "前一区块哈希值：0x" +  BytesToHex(block.HashPrevBlock) + "<br>"
			content = content + "=============================================" + "<br>"
		}

		content = content + "</html>"
		response.Write([]byte(content))
	}

//处理Url为 http://SERVER_ADDR:8080/generateblock 的请求，生成新区块
	func (bb *BlockBrower) handleGenerateBlock(response http.ResponseWriter, request *http.Request) {

		//获取区块链高度
		height := bb.blockChain.Iterator().GetCount()

		//创建新的区块
		bb.blockChain.AddBlock("Mini block " + strconv.Itoa(height))

		blockHeight := strconv.Itoa(bb.blockChain.Iterator().GetCount())

		content := "<html><br>"
		content = content + "<a href=\"/\">返回首页</a>&nbsp;&nbsp;&nbsp;&nbsp;"
		content = content + "<a href=\"generateblock\">生成新区块</a>&nbsp;&nbsp;&nbsp;&nbsp;<br>"
		content = content + "<br><b>&nbsp;&nbsp;&nbsp;&nbsp;当前MiniBC区块链高度：" + blockHeight + "</b><br><br>"

		//遍历区块链
		iterator := bb.blockChain.Iterator()
		for {
			block := iterator.Next()
			if block == nil {
				break
			}
			content = content + "当前区块哈希值：0x" + BytesToHex(block.GetHash()) + "<br>"
			content = content + "当前区块内容为：" +  string(block.Data) + "<br>"
			content = content + "前一区块哈希值：0x" +  BytesToHex(block.HashPrevBlock) + "<br>"
			content = content + "=============================================" + "<br>"
		}

		content = content + "</html>"
		response.Write([]byte(content))
	}


运行后输出：


![Image text](https://github.com/wangshizebin/minibc/blob/master/v003/index.png)
![Image text](https://github.com/wangshizebin/minibc/blob/master/v003/viewblock.png)


#### 我们实现的和将要实现的  

我们已经构建了一个非常简单的区块链，实现了一个简易的KV数据库和区块链浏览器，下一步来实现工作量证明pow共识算法。  


#### 交流和疑问  

大家可以从v001开始逐步深入的研究学习，每个源码文件均有完整的说明和注释。如果有疑问、建议和不明白的问题，尽可以与我联系。

MiniBC区块链交流可以加入qq群：<font size=5><b>777804802</b></font>，go开发者乐园
我的微信：<font size=5><b>bkra50 </b></font>  
github addr: https://github.com/wangshizebin/minibc

另外，招募一群志同道合者，共同维护项目，共同推进项目，为开源世界贡献力量，欢迎勾搭。

Let's go~~~~~~~~~~~~~~~~~~~~~~~~~~~~

