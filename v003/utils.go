package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//将32位整型数字按大端模式转换成字节切片
func IntToByte(val int32) []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, val)
	return buffer.Bytes()
}

//将字节切片转换成十六进制格式
func BytesToHex(data []byte) string {
	dst := fmt.Sprintf("%x", data)
	return dst
}
