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

// Package main implements a client for Greeter service.
package main

import (
	pb "TripleGrpcNew/TripleGrpc"
	"context"
	"fmt"
	"log"
	"time"
	"unsafe"

	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1:50052"
)

func main() {
	start := time.Now()
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//cost1 := time.Since(start)

	c := pb.NewGetTripleClient(conn)

	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//cost2 := time.Since(start)
	defer cancel()
	r, err := c.GetIntTriple(ctx, &pb.GetTripleRequest{Group: 7, Role: 1, Nlen: 10000})
	cost3 := time.Since(start)
	// fmt.Println("cost1=", cost1)
	// fmt.Println("cost2=", cost2)
	fmt.Println("服务请求总耗时=", cost3)
	log.Fatalf("over", unsafe.Sizeof(r), err)

}
