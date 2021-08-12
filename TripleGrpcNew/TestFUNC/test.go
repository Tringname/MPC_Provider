package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

func main() {
	var num int64 = 1596 // 64位整数，8个字节
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, num)
	if err != nil {
		log.Fatal()
	}
	bytes := buf.Bytes()
	fmt.Println(bytes) // [0 0 0 0 0 0 0 15]

	var decodingNum int64
	err = binary.Read(&buf, binary.BigEndian, &decodingNum)
	if err != nil {
		log.Fatal()
	}
	fmt.Println(decodingNum) // 15

	buf.Reset()
	err = binary.Write(&buf, binary.LittleEndian, num)
	if err != nil {
		log.Fatal()
	}
	bytes = buf.Bytes()
	fmt.Println(bytes) // [15 0 0 0 0 0 0 0]
	err = binary.Read(&buf, binary.LittleEndian, &decodingNum)
	if err != nil {
		log.Fatal()
	}
	fmt.Println(decodingNum) // 15

}
