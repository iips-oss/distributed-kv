// centralized server connecting gRPC client

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/iips-oss/distributed-kv/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("didn't connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKvstoreClient(conn)

	for { // maybe unsafe
		fmt.Printf("> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		cmd := strings.Split(line, " ")
		if len(cmd) > 3 || len(cmd) < 2 {
			continue
		}
		method := cmd[0]
		key := cmd[1]

		switch method {
		case "GET":
			fmt.Printf("GET reviceved \n")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.KvGet(ctx, &pb.OpKeyReq{Key: key})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", r.GetValue())
		case "SET":
			if len(cmd) != 3 {
				fmt.Printf("Error: two arguments required\n")
			}
			value := cmd[2]
			fmt.Printf("SET reviceved \n")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := c.KvSet(ctx, &pb.SetReq{Key: key, Value: value})
			if err != nil {
				log.Fatal(err)
			}
		case "DEL":
			fmt.Printf("DEL reviceved \n")
		case "QUIT":
			fmt.Printf("QUITING\n")
			break
		}
		// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		// defer cancel()
		// r, err := c.HiLol(ctx, &pb.HiReq{ClientId: *client_id})
	}
}
