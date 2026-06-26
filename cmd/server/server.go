// For now we'll do the tidwall/btree in memory implementation of single
// node server for persistant kv store and a client for interacting
// using tidwall/btree instead of google/btree, examples are simpler to understand
// TODO: replace of tidwall/btree with our own implementation

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/iips-oss/distributed-kv/protobuf"
	"github.com/tidwall/btree"
	"google.golang.org/grpc"
)

// INFO: https://en.wikipedia.org/wiki/Copy-on-write
// tidwall/btree library has a few useful funcitons,
// like bulk loading keys with Load() and Copy() copy-on-write

// reason to use btree instead of Hashmap, is that btree preserve the
// lexical order of keys which is efficent for bulk-match key retreval.

var (
	port  = flag.Int("port", 50051, "sever port")
	store btree.Map[string, string]
)

type server struct {
	pb.UnimplementedKvstoreServer
}

// gRPC functions
func (s *server) KvGet(_ context.Context, in *pb.OpKeyReq) *pb.OpGetRes {
	log.Printf("Received: %v", in.GetKey())
	value, err := getKey(store, in.GetKey())
	return &pb.OpGetRes{Value: value, Err: err}
}
func (s *server) KvSet(_ context.Context, in *pb.SetReq) *pb.OpRes {
	key := in.GetKey()
	value := in.GetValue()
	log.Printf("Received: %v %v", key, value)
	value, err := store.Set(key, value)
	return &pb.OpRes{Err: err}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterKvstoreServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Iterate over btree and print it in order
func printMap(kv btree.Map[string, string]) {
	kv.Scan(func(key string, value string) bool {
		fmt.Printf("%s %s\n", key, value)
		return true
	})
}

// Iterate over btree to fetch single key
// @TODO implment bulk fetch of keys
func getKey(kv btree.Map[string, string], k string) (string, bool) {
	var ret string
	kv.Scan(func(key string, value string) bool {
		if k == key {
			ret = value
		}
		return true
	})
	return ret, false
}
