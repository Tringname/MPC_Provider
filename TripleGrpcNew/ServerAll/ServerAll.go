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
var bitmx sync.RWMutex

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

var IntDataLIst = [100][2][]*pb.Triple{
	{{{A: 4, B: 79, C: 4460}, {A: 5147683243598419542, B: 4822789147564851964, C: 2033364834455673516}},
		{{A: 77, B: 8, C: 2587}, {A: 2817246957421286492, B: 988019566165373957, C: 9000756907300605446}}},
}
var BitDataList = [100][2][]*pb.BitTriple{
	{{{A: 0, B: 0, C: 1}, {A: 1, B: 1, C: 0}},
		{{A: 1, B: 0, C: 1}, {A: 0, B: 0, C: 1}}},
}

//读写函数
func IntResult(Group int, Role int, Nlen int) []*pb.Triple {
	mx.Lock()
	if !IsContain(Group, Role, Nlen) {
		DataA, DataB := BeaverTriple(Nlen)
		IntDataLIst[Group][0] = append(IntDataLIst[Group][0], DataA...)
		IntDataLIst[Group][1] = append(IntDataLIst[Group][1], DataB...)
	}
	res := GetValue(Group, Role, Nlen)
	mx.Unlock()
	return res
}

//判断原数据中是否存在该队列数据，如存在进入读函数，如不存在进入写函数
func whichlock(Group int, Role int, Nlen int) []*pb.Triple {
	var resultn []*pb.Triple
	if Group >= 0 && Role < 3 && Nlen > 0 {
		resultn = IntResult(Group, Role, Nlen)
	} else {
		fmt.Println("参数输入出错", time.Now(), "Role=", Role, "Nlen=", Nlen)
	}

	return resultn
}

//返回Map中该数据
func GetValue(Group int, Role int, Nlen int) []*pb.Triple {
	if len(IntDataLIst[Group][Role]) >= Nlen {
		result := IntDataLIst[Group][Role][0:Nlen]
		IntDataLIst[Group][Role] = IntDataLIst[Group][Role][Nlen:]
		return result
	} else {
		remain := IntDataLIst[Group][Role][0:len(IntDataLIst[Group][Role])]
		IntDataLIst[Group][Role] = []*pb.Triple{}
		return remain
	}
}

////判断Map中是否存在该数据
func IsContain(Group int, Role int, Nlen int) bool {
	if len(IntDataLIst) >= Group {
		if len(IntDataLIst[Group][Role]) >= Nlen {
			return true
		}
	} else {
		initdata := [2][]*pb.Triple{}
		IntDataLIst[Group] = initdata
		return false
	}
	return false
}

//三元组中已知a,随机生成a1,a2
func SecretCom(S int64) (int64, int64) {
	rand.Seed(time.Now().UnixNano())
	s1 := rand.Int63()
	s2 := S - s1
	return s1, s2
}

//三元组生成函数
func BeaverTriple(Nlen int) ([]*pb.Triple, []*pb.Triple) { //三元组生成函数
	var listA []*pb.Triple
	var listB []*pb.Triple
	for i := 1; i <= Nlen; i++ {
		rand.Seed(time.Now().UnixNano())
		var a = rand.Int63()
		var b = rand.Int63()
		var c = a * b
		a1, a2 := SecretCom(a)
		b1, b2 := SecretCom(b)
		c1, c2 := SecretCom(c)
		Atuple := []*pb.Triple{{A: a1, B: b1, C: c1}}
		listA = append(listA, Atuple...)
		Btuple := []*pb.Triple{{A: a2, B: b2, C: c2}}
		listB = append(listB, Btuple...)
	}
	return listA, listB
}

//读写函数
func BitResult(Group int, Role int, Nlen int) []*pb.BitTriple {
	startT2 := time.Now()
	bitmx.Lock()
	if !BitIsContain(Group, Role, Nlen) {
		DataA, DataB := BitBeaverTriple(Nlen)
		BitDataList[Group][0] = append(BitDataList[Group][0], DataA...)
		BitDataList[Group][1] = append(BitDataList[Group][1], DataB...)
	}

	tc2append := time.Since(startT2)
	res := BitGetValue(Group, Role, Nlen)
	tcGet := time.Since(startT2)
	bitmx.Unlock()
	tc2 := time.Since(startT2) //计算耗时
	//fmt.Printf("readGogene耗时 = %v\n", tc2gene)
	fmt.Printf("三元组添加到变量时间 = %v\n", tc2append)
	fmt.Printf("三元组查询时间 = %v\n", tcGet)
	fmt.Printf("三元组返回时间 = %v\n", tc2)
	return res
}

//判断原数据中是否存在该队列数据，如存在进入读函数，如不存在进入写函数
func Bitwhichlock(Group int, Role int, Nlen int) []*pb.BitTriple {
	var resultn []*pb.BitTriple
	if Group >= 0 && Role < 3 && Nlen > 0 {
		resultn = BitResult(Group, Role, Nlen)
	} else {
		fmt.Println("参数输入出错", time.Now(), "Role=", Role, "Nlen=", Nlen)
	}

	return resultn
}

//返回Map中该数据
func BitGetValue(Group int, Role int, Nlen int) []*pb.BitTriple {
	if len(BitDataList[Group][Role]) >= Nlen {
		result := BitDataList[Group][Role][0:Nlen]
		BitDataList[Group][Role] = BitDataList[Group][Role][Nlen:]
		return result
	} else {
		remain := BitDataList[Group][Role][0:len(BitDataList[Group][Role])]
		BitDataList[Group][Role] = []*pb.BitTriple{}
		return remain
	}
}

////判断Map中是否存在该数据
func BitIsContain(Group int, Role int, Nlen int) bool {
	if len(BitDataList) >= Group {
		if len(BitDataList[Group][Role]) >= Nlen {
			return true
		}
	} else {
		initdata := [2][]*pb.BitTriple{}
		BitDataList[Group] = initdata
		return false
	}
	return false
}

//位三元组生成
func BitBeaverTriple(Nlen int) ([]*pb.BitTriple, []*pb.BitTriple) { //三元组生成函数
	startT3 := time.Now()
	var listA []*pb.BitTriple
	var listB []*pb.BitTriple
	for i := 1; i <= Nlen; i++ {
		rand.Seed(time.Now().UnixNano())

		a := rand.Intn(2)
		b := rand.Intn(2)
		c := a & b
		a1, a2 := BitSecretCom(a)
		b1, b2 := BitSecretCom(b)
		c1, c2 := BitSecretCom(c)
		Atuple := []*pb.BitTriple{{A: a1, B: b1, C: c1}}
		listA = append(listA, Atuple...)
		Btuple := []*pb.BitTriple{{A: a2, B: b2, C: c2}}
		listB = append(listB, Btuple...)

	}
	tcBe := time.Since(startT3) //计算耗时
	fmt.Printf("BeaverTriple生成耗时 = %v\n", tcBe)
	return listA, listB
}
func BitSecretCom(S int) (uint32, uint32) {
	s1 := rand.Intn(2)
	s2 := s1 ^ S
	return uint32(s1), uint32(s2)
}

//GetIntTriple
func (s *GrpcserverNew) GetIntTriple(ctx context.Context, in *pb.GetTripleRequest) (*pb.TripleList, error) {

	res1 := whichlock(int(in.Group), int(in.Role), int(in.Nlen))
	Atnew := &pb.TripleList{Triple: res1}

	return Atnew, nil
}

//GetBitTriple
func (s *GrpcserverNew) GetBitTriple(ctx context.Context, in *pb.GetTripleRequest) (*pb.BitTripleList, error) {
	res2 := Bitwhichlock(int(in.Group), int(in.Role), int(in.Nlen))
	Atnew := &pb.BitTripleList{Bittriple: res2}
	return Atnew, nil
}
func (s *GrpcserverNew) GetBytesTriple(ctx context.Context, in *pb.GetBytesTripleRequest) (*pb.BytesTripleList, error) {
	re := &pb.BytesTripleList{}
	return re, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGetTripleServer(s, &GrpcserverNew{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
