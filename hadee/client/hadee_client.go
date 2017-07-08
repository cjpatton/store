// Copyright (c) 2017, Christopher Patton.
// All rights reserved.

// hadee_client is a toy client that makes RPC requests to hadee_server. The
// first request gets the parameters, then it prompts the user for actual
// requests.
//
// Usage: hadee_client user
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cjpatton/store"
	"github.com/cjpatton/store/pb"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address  = "localhost:50051"
	greeting = "---- Hadee ----------------------------------------------------\n" +
		"Welcome to Hadee, your super secret source of dog jokes."
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: hadee_client user")
		return
	}
	user := os.Args[1]

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewStoreProviderClient(conn)

	// Get parameters for store associated to user.
	paramsReply, err := c.GetParams(context.Background(), &pb.ParamsRequest{UserId: user})
	if err != nil {
		fmt.Println("ParamsRequest fails", err)
		return
	} else if paramsReply.GetError() == pb.StoreProviderError_BAD_USER {
		fmt.Printf("ParamsRequest fails: user %q not found\n", user)
		return
	} else if paramsReply.GetError() != pb.StoreProviderError_OK {
		fmt.Println("ParamsRequest fails:", paramsReply.GetError())
		return
	}

	// Get password from client and set up the private store context.
	fmt.Println(greeting)
	fmt.Print("Please enter the master password> ")
	password, err := terminal.ReadPassword(0) // os.Stdin
	if err != nil {
		fmt.Println("terminal.ReadPassword() fails:", err)
		return
	}
	key := store.DeriveKeyFromPassword(password, nil)
	priv, err := store.NewPrivStore(key, paramsReply.GetParams())
	if err != nil {
		fmt.Println("store.NewPrivStore() fails:", err)
		return
	}
	defer priv.Free()

	bio := bufio.NewReader(os.Stdin)
	fmt.Println("\nEnter an input and we'll give you the output. Type \"ls\" to")
	fmt.Println("see the list of inputs. Type \"quit\" to leave.")
	fmt.Println("---------------------------------------------------------------")

	for {
		fmt.Print("> ")
		bin, _, err := bio.ReadLine()
		if err != nil {
			fmt.Println("bio.ReadLine() fails:", err)
			return
		}
		in := string(bin)
		if in == "quit" {
			fmt.Println("Go gators!")
			break
		}

		x, y, err := priv.GetIdx(in)
		if err != nil {
			fmt.Println("priv.GetIdx(in) fails:", err)
			return
		}
		shareReply, err := c.GetShare(context.Background(),
			&pb.ShareRequest{
				UserId: user,
				X:      int32(x),
				Y:      int32(y),
			},
		)
		if err != nil {
			fmt.Println("ShareRequest fails:", err)
			return
		} else if shareReply.GetError() == pb.StoreProviderError_ITEM_NOT_FOUND {
			fmt.Println("Server says item not found. (Wrong master password?)")
			continue
		} else if shareReply.GetError() != pb.StoreProviderError_OK {
			fmt.Println("ShareRequest fails:", shareReply.GetError())
			return
		}

		out, err := priv.GetOutput(in, shareReply.GetPubShare())
		if err == store.ItemNotFound {
			fmt.Println("Item not found. (Wrong master password?)")
		} else if err != nil {
			fmt.Println("priv.GetValue() fails:", err)
			return
		} else {
			fmt.Println(out)
		}
	}
}
