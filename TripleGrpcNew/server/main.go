package main

/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	pb "TripleGrpcNew/TripleGrpc"

	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

var mx sync.RWMutex

// server is used to implement helloworld.GreeterServer.
type GrpcserverNew struct {
	pb.UnimplementedGetTripleServer
}

type InputData struct {
	Role       int
	TripleList []*pb.Triple
}

// var DataList = []InputData{
// 	{Role:0,TripleList:{*pb.Triple{A: 4, B: 79, C: 4460}, {A: 5147683243598419542, B: 4822789147564851964, C: 2033364834455673516}}},
// 	{Role:1,TripleList:{{A: 77, B: 8, C: 2587}, {A: 2817246957421286492, B: 988019566165373957, C: 9000756907300605446}}}}

var DataListNew = [100][2][]*pb.Triple{
	{{{A: 4, B: 79, C: 4460}, {A: 5147683243598419542, B: 4822789147564851964, C: 2033364834455673516}},
		{{A: 77, B: 8, C: 2587}, {A: 2817246957421286492, B: 988019566165373957, C: 9000756907300605446}}},
}

//统计耗时函数
func timeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("time cost = %v\n", tc)
	}
}

//读函数
func readGo(Group int, Role int, Nlen int) []*pb.Triple {
	startT2 := time.Now()
	mx.Lock()
	if !IsContain(Group, Role, Nlen) {
		DataA, DataB := BeaverTriple(Nlen)
		DataListNew[Group][0] = append(DataListNew[Group][0], DataA...)
		DataListNew[Group][1] = append(DataListNew[Group][1], DataB...)
	}
	tc2append := time.Since(startT2)
	res := GetValue(Group, Role, Nlen)
	tcGet := time.Since(startT2)
	mx.Unlock()
	tc2 := time.Since(startT2) //计算耗时
	//fmt.Printf("readGogene耗时 = %v\n", tc2gene)
	fmt.Printf("三元组添加到变量时间 = %v\n", tc2append)
	fmt.Printf("三元组查询时间 = %v\n", tcGet)
	fmt.Printf("三元组返回时间 = %v\n", tc2)
	return res

}

// //写函数
// func writeGo(Group int, Role int, Nlen int) []*pb.Triple {
// 	mx.Lock()
// 	if !IsContain(Group, Role, Nlen) {
// 		DataA, DataB := BeaverTriple(Nlen)
// 		DataListNew[Group][0] = append(DataListNew[Group][0], DataA...)
// 		DataListNew[Group][1] = append(DataListNew[Group][1], DataB...)
// 	}
// 	res := GetValue(Group, Role, Nlen)
// 	mx.Unlock()
// 	return res
// }

//判断原数据中是否存在该队列数据，如存在进入读函数，如不存在进入写函数
func whichlock(Group int, Role int, Nlen int) []*pb.Triple {
	startT1 := time.Now()
	var resultn []*pb.Triple

	// lockres := IsContain(Group, Role, Nlen)
	// if lockres {
	// 	resultn = readGo(Group, Role, Nlen)
	// } else if !lockres {
	// 	resultn = writeGo(Group, Role, Nlen)
	// } else {
	// 	fmt.Println("参数输入出错", time.Now(), "Role=", Role, "Nlen=", Nlen)
	// }
	// return resultn
	if Group >= 0 && Role < 3 && Nlen > 0 {
		resultn = readGo(Group, Role, Nlen)
	} else {
		fmt.Println("参数输入出错", time.Now(), "Role=", Role, "Nlen=", Nlen)
	}
	tc1 := time.Since(startT1) //计算耗时
	fmt.Printf("whichclock = %v\n", tc1)
	return resultn
}

//返回Map中该数据
func GetValue(Group int, Role int, Nlen int) []*pb.Triple {
	if len(DataListNew[Group][Role]) >= Nlen {
		result := DataListNew[Group][Role][0:Nlen]
		DataListNew[Group][Role] = DataListNew[Group][Role][Nlen:]
		return result
	} else {
		remain := DataListNew[Group][Role][0:len(DataListNew[Group][Role])]
		DataListNew[Group][Role] = []*pb.Triple{}
		return remain
	}
}

////判断Map中是否存在该数据
func IsContain(Group int, Role int, Nlen int) bool {
	if len(DataListNew) >= Group {
		if len(DataListNew[Group][Role]) >= Nlen {
			return true
		}
	} else {
		initdata := [2][]*pb.Triple{}
		DataListNew[Group] = initdata
		return false
	}
	return false
}

//三元组生成函数
func BeaverTriple(Nlen int) ([]*pb.Triple, []*pb.Triple) { //三元组生成函数
	startT3 := time.Now()
	var listA []*pb.Triple
	var listB []*pb.Triple
	for i := 1; i <= Nlen; i++ {
		rand.Seed(time.Now().UnixNano())
		var a = rand.Int63n(100)
		var b = rand.Int63n(100)
		var c = a * b
		a1, a2 := SecretCom(a)
		b1, b2 := SecretCom(b)
		c1, c2 := SecretCom(c)
		Atuple := []*pb.Triple{{A: a1, B: b1, C: c1}}
		listA = append(listA, Atuple...)
		Btuple := []*pb.Triple{{A: a2, B: b2, C: c2}}
		listB = append(listB, Btuple...)
	}
	tcBe := time.Since(startT3) //计算耗时
	fmt.Printf("BeaverTriple生成耗时 = %v\n", tcBe)
	return listA, listB
}

//三元组中已知a,随机生成a1,a2
func SecretCom(S int64) (int64, int64) {
	rand.Seed(time.Now().UnixNano())
	s1 := rand.Int63n(100)
	s2 := S - s1
	return s1, s2
}

// func loop() {
// 	for i := 1; i < 10; i++ {
// 		go whichlock(0, i)
// 		go whichlock(1, i)
// 	}
// }

// SayHello implements helloworld.GreeterServer
func (s *GrpcserverNew) GetIntTriple(ctx context.Context, in *pb.GetTripleRequest) (*pb.TripleList, error) {

	res1 := whichlock(int(in.Group), int(in.Role), int(in.Nlen))
	Atnew := &pb.TripleList{Triple: res1}

	return Atnew, nil
}
func (s *GrpcserverNew) GetBitTriple(ctx context.Context, in *pb.GetTripleRequest) (*pb.BitTripleList, error) {
	re := &pb.BitTripleList{}
	return re, nil
}
func (s *GrpcserverNew) GetBytesTriple(ctx context.Context, in *pb.GetBytesTripleRequest) (*pb.BytesTripleList, error) {
	re := &pb.BytesTripleList{}
	return re, nil
}

func main() {
	startT := time.Now()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGetTripleServer(s, &GrpcserverNew{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	tcall := time.Since(startT) //计算耗时
	fmt.Printf("服务端响应总耗时 = %v\n", tcall)
}
