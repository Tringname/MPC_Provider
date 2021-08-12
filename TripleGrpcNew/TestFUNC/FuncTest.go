package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
func main() {
	var a int
	a = 3
	res := IntToBytes(a)
	b := BytesToInt(res)
	fmt.Println(res)
	fmt.Println(b)
	// var a uint32
	// var b uint32

	// rand.Seed(time.Now().UnixNano())

	// a := rand.Intn(2)
	// b := rand.Intn(2)
	// c := a & b
	// a1, a2 := SecretCom(a)
	// b1, b2 := SecretCom(b)
	// c1, c2 := SecretCom(c)
	// fmt.Println(uint32(a1))
	// fmt.Println(uint32(a2))
	// fmt.Println(uint32(b1))
	// fmt.Println(uint32(b2))
	// fmt.Println(uint32(c1))
	// fmt.Println(uint32(c2))

	// v := uint32(500) //转换为无符号32位
	// buf := make([]byte, 4)

	// //将 v 写入 buf中
	// binary.BigEndian.PutUint32(buf, v)

	// //read
	// //从buf中读取并赋值给x
	// x := binary.BigEndian.Uint32(buf)

	// //res, err := GenerateRandomBytes(3)
	// fmt.Println("res", x)
	// var a int
	// p := &a
	// b := [10]int64{1}
	// s := "adsa"
	// bs := make([]byte, 10)

	// fmt.Println(binary.Size(a))  // -1
	// fmt.Println(binary.Size(p))  // -1
	// fmt.Println(binary.Size(b))  // 80
	// fmt.Println(binary.Size(s))  // -1
	// fmt.Println(binary.Size(bs)) // 10
	//00000011

}
