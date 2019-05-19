package main

import (
	"fmt"
	"net/http"
	"strconv"
)

//表示一个区块链浏览器
type BlockBrowser struct {
	//中止BrowserServer通道
	chanQuit chan bool

	//当前区块链
	blockChain *BlockChain
}

const (
	SERVER_ADDR = "localhost"
)
//处理Url为 http://SERVER_ADDR:8080/ 的请求，显示区块链浏览器首页信息
func (bb *BlockBrowser) handleIndex(response http.ResponseWriter, request *http.Request) {
	content := "<html><br>"
	content = content + "<b>&nbsp;&nbsp;&nbsp;&nbsp;MiniBC区块链区块浏览器</b><br><br>"
	content = content + "&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"shutdown\">关闭MiniBC</a>"
	content = content + "&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"getblocks\">查看区块链</a>&nbsp;&nbsp;&nbsp;&nbsp;"
	content = content + "</html>"

	response.Write([]byte(content))
}

//处理Url为 http://SERVER_ADDR:8080/shutdown 的请求，关闭区块链浏览器，退出系统
func (bb *BlockBrowser) handleShutdown(response http.ResponseWriter, request *http.Request) {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("=========================================================================")
	fmt.Println("MiniBC 区块链和远程管理结束，谢谢使用!")
	fmt.Println("=========================================================================")

	bb.chanQuit <- true
}

//处理Url为 http://SERVER_ADDR:8080/getblocks 的请求，打印所有区块链信息
func (bb *BlockBrowser) handleGetBlocks(response http.ResponseWriter, request *http.Request) {

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
		content = content + "当前区块哈希值：" + BytesToHex(block.GetHash()) + "<br>"
		content = content + "当前区块内容为：" +  string(block.Data) + "<br>"
		content = content + "前一区块哈希值：" +  BytesToHex(block.HashPrevBlock) + "<br>"
		content = content + "=============================================" + "<br>"
	}

	content = content + "</html>"
	response.Write([]byte(content))
}

//处理Url为 http://SERVER_ADDR:8080/generateblock 的请求，生成新区块
func (bb *BlockBrowser) handleGenerateBlock(response http.ResponseWriter, request *http.Request) {

	//获取区块链高度
	height := bb.blockChain.Iterator().GetCount()

	//创建新的区块
	bb.blockChain.AddBlock("MiniBC Block " + strconv.Itoa(height))

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
		content = content + "当前区块哈希值：" + BytesToHex(block.GetHash()) + "<br>"
		content = content + "当前区块内容为：" +  string(block.Data) + "<br>"
		content = content + "前一区块哈希值：" +  BytesToHex(block.HashPrevBlock) + "<br>"
		content = content + "=============================================" + "<br>"
	}

	content = content + "</html>"
	response.Write([]byte(content))
}


//启动区块链浏览器
func (bb *BlockBrowser) Start() {
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

//创建一个区块链浏览器BlockBrowser
func NewBlockBrowser(blockchain *BlockChain) *BlockBrowser {
	BlockBrowser := BlockBrowser{}
	BlockBrowser.chanQuit = make(chan bool)
	BlockBrowser.blockChain = blockchain
	return &BlockBrowser
}
