// Copyright (c) 2017, Christopher Patton.
// All rights reserved.

// hadee_serv is a toy server implementing the StoreProvider RPC specified in
// store.proto. It services requests for only one user, whose identity and table
// are specified via the command line.
//
// Usage: hadee_serv user store.pub
package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/cjpatton/store"
	"github.com/cjpatton/store/pb"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// HadeeStoreProvider implements the StoreProvider RPC.
type HadeeStoreProvider struct {
	pubs   map[string](*store.PubStore)
	params map[string](*pb.Params)
}

// NewHadeeStoreProvider creates a new HadeeStoreProvider.
//
// NOTE Must be dstroyed with s.CleanUp().
func NewHadeeStoreProvider(user string, table *pb.Store) *HadeeStoreProvider {
	s := new(HadeeStoreProvider)
	s.pubs = make(map[string](*store.PubStore))
	s.params = make(map[string](*pb.Params))
	s.pubs[user] = store.NewPubStoreFromProto(table)
	s.params[user] = table.GetDict().GetParams()
	return s
}

// CleanUp frees memory allocated to each pub in pubs. This is necessary because
// the underlying data structure is implemented in C.
func (s *HadeeStoreProvider) CleanUp() {
	for _, pub := range s.pubs {
		if pub != nil {
			pub.Free()
		}
	}
}

func (s *HadeeStoreProvider) GetShare(ctx context.Context, in *pb.ShareRequest) (*pb.ShareReply, error) {
	log.Println("GetShare")
	if pub, ok := s.pubs[in.GetUserId()]; ok {
		if pubShare, err := pub.GetShare(int(in.GetX()), int(in.GetY())); err == nil {
			return &pb.ShareReply{Error: pb.StoreProviderError_OK, PubShare: pubShare}, nil
		} else if err == store.ErrorIdx {
			return &pb.ShareReply{Error: pb.StoreProviderError_INDEX}, nil
		} else if err == store.ItemNotFound {
			return &pb.ShareReply{Error: pb.StoreProviderError_ITEM_NOT_FOUND}, nil
		} else {
			return nil, err // Unexpected error!
		}
	}
	return &pb.ShareReply{Error: pb.StoreProviderError_BAD_USER}, nil
}

func (s *HadeeStoreProvider) GetParams(ctx context.Context, in *pb.ParamsRequest) (*pb.ParamsReply, error) {
	log.Println("GetParams")
	if params, ok := s.params[in.GetUserId()]; ok {
		return &pb.ParamsReply{Error: pb.StoreProviderError_OK, Params: params}, nil
	}
	return &pb.ParamsReply{Error: pb.StoreProviderError_BAD_USER}, nil
}

func main() {

	if len(os.Args) != 3 {
		log.Fatal("error: usage: hadee_server user store.pub")
	}
	user := os.Args[1]
	tableString, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	table := new(pb.Store)
	if err = proto.Unmarshal(tableString, table); err != nil {
		log.Fatal("failed to parse protobuf: ", err)
	}

	// Begin serving.
	storeProvider := NewHadeeStoreProvider(user, table)
	defer storeProvider.CleanUp()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Opened TCP socket on", port)
	s := grpc.NewServer()
	pb.RegisterStoreProviderServer(s, storeProvider)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
