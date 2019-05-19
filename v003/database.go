package main

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

//Database 实现了一个简单的kv数据库
//数据库中可以创建多张表，但每张表的记录只能存放两个字段key和value
//数据库两个操作方法Get和Set,分别用于读取和存放数据
type Database struct {
	data map[string]map[string][]byte
}

//数据库文件名字
const (
	DATABASE_FILE = "./block.db"
)

//打开数据库，初始化数据
func (db *Database) open() {

	//初始化Database数据结构
	db.data = make(map[string]map[string][]byte)

	//如果数据库文件不存在，那么不需要读取数据，直接返回
	if IsNotExists(DATABASE_FILE) {
		return
	}

	//用只读方式打开数据库文件
	f, err := os.OpenFile(DATABASE_FILE, os.O_RDONLY, 0)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	//读取数据库文件的全部内容
	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Panic(err)
	}

	//解码数据库中序列化的数据，还原到Database.data中
	reader := bytes.NewReader(content)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&(db.data))
	if err != nil {
		log.Panic(err)
	}
}

//关闭数据库
func (db *Database) Close() {

	//将Database.data数据序列化为二进制数据，等待写入数据库文件
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(db.data)
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

//从表table的记录里读取主键为key的记录的值
func (db *Database) Get(table, key string) []byte {
	if row, ok := db.data[table]; ok {
		if val, ok := row[key]; ok {
			return val
		}
	}

	//如果键值不存在，返回空数据
	return []byte{}
}

//向表table添加一条记录，主键为key, 值为val
func (db *Database) Set(table, key string, val []byte) {
	if _, ok := db.data[table]; ok{
		db.data[table][key] = val
	}else {
		db.data[table] = map[string][]byte{key: val}
	}
}

//创建一个Database对象
func NewDatabase() *Database {
	Database := new(Database)
	Database.open()
	return Database
}

//判断文件是否不存在，如果不存在返回true
func IsNotExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return true
	}

	return false
}
