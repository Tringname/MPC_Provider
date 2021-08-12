package main

import (
	"fmt"
	"strconv"
)

//bin表示转化后的位数
func convertToBin(n int, bin int) string {
	var b string
	switch {
	case n == 0:
		for i := 0; i < bin; i++ {
			b += "0"
		}
	case n > 0:
		//strcov.Itoa 将 1 转为 "1" , string(1)直接转为assic码
		for ; n > 0; n /= 2 {
			b = strconv.Itoa(n%2) + b
		}
		//加0
		j := bin - len(b)
		for i := 0; i < j; i++ {
			b = "0" + b
		}
	case n < 0:
		n = n * -1
		// fmt.Println("变为整数：",n)
		s := convertToBin(n, bin)
		// fmt.Println("bin:",s)
		//取反
		for i := 0; i < len(s); i++ {
			if s[i:i+1] == "1" {
				b += "0"
			} else {
				b += "1"
			}
		}
		// fmt.Println("~bin :",b)
		//转化为整形，之后加1 这里必须要64，否则在转化过程中可能会超出范围
		n, err := strconv.ParseInt(b, 2, 64)
		if err != nil {
			fmt.Println(err)
		}
		//转为bin
		//+1
		b = convertToBin(int(n+1), bin)
	}
	return b
}

func main() {
	fmt.Println(
		convertToBin(5, 11),
		TypeOf(convertToBin(5, 11)),
		//101
		// 	convertToBin(13, 8), //1101
		// 	convertToBin(11111, 8),
		// 	convertToBin(0, 8),
		// 	convertToBin(1, 8),
		// 	convertToBin(-5, 8),
		// 	convertToBin(-11111, 8),
	)

}
