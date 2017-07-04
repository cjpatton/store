// Copyright (c) 2017, Christopher Patton.
// All rights reserved.
package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/cjpatton/store"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type HadeeStoreProvider struct {
	pubs   map[string](*store.PubStore)
	params map[string](*store.StoreParams)
}

func NewHadeeStoreProvider(user string, table *store.StoreTable) *HadeeStoreProvider {
	s := new(HadeeStoreProvider)
	s.pubs = make(map[string](*store.PubStore))
	s.params = make(map[string](*store.StoreParams))
	s.pubs[user] = store.NewPubStoreFromTable(table)
	s.params[user] = table.GetParams()
	return s
}

func (s *HadeeStoreProvider) CleanUp() {
	for _, pub := range s.pubs {
		if pub != nil {
			pub.Free()
		}
	}
}

func (s *HadeeStoreProvider) GetShare(ctx context.Context, in *store.ShareRequest) (*store.ShareReply, error) {
	log.Println("GetShare")
	if pub, ok := s.pubs[in.GetUserId()]; ok {
		if pubShare, err := pub.GetShare(int(in.GetX()), int(in.GetY())); err == nil {
			return &store.ShareReply{Error: store.StoreProviderError_OK, PubShare: pubShare}, nil
		} else if err == store.ErrorIdx {
			return &store.ShareReply{Error: store.StoreProviderError_INDEX}, nil
		} else {
			return nil, err // Unexpected error!
		}
	}
	return &store.ShareReply{Error: store.StoreProviderError_BAD_USER}, nil
}

func (s *HadeeStoreProvider) GetParams(ctx context.Context, in *store.ParamsRequest) (*store.ParamsReply, error) {
	log.Println("GetParams")
	if params, ok := s.params[in.GetUserId()]; ok {
		return &store.ParamsReply{Error: store.StoreProviderError_OK, Params: params}, nil
	}
	return &store.ParamsReply{Error: store.StoreProviderError_BAD_USER}, nil
}

func main() {

	if len(os.Args) != 3 {
		log.Fatal("error: usage: hadee_server user table.pub")
	}

	user := os.Args[1]

	tableString, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	table := new(store.StoreTable)
	if err = proto.Unmarshal(tableString, table); err != nil {
		log.Fatal("failed to parse protobuf: ", err)
	}

	storeProvider := NewHadeeStoreProvider(user, table)
	defer storeProvider.CleanUp()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Opened TCP socket on", port)
	s := grpc.NewServer()
	store.RegisterStoreProviderServer(s, storeProvider)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
