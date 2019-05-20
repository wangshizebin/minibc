package main

func main() {

	//新建一条区块链，如果区块链数据库中已经有了数据，将会读取数据库加载进来
	//如果尚未创建数据库，或者数据库为空，那么会自动生成一个创世纪区块
	bc := NewBlockChain()
	defer bc.DB.Close()

	//启动区块链浏览器，您可以通过浏览器 http://localhost:8080 访问
	NewBlockBrowser(bc).Start()
}
