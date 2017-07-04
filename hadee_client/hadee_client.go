// Copyright (c) 2017, Christopher Patton.
// All rights reserved.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cjpatton/store"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address  = "localhost:50051"
	greeting = "---- Hadee ----------------------------------------------------\n" +
		"Welcome to Hadee, an excellent source of password suggestions."
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
	c := store.NewStoreProviderClient(conn)

	// Get parameters for store associated to user.
	paramsReply, err := c.GetParams(context.Background(), &store.ParamsRequest{UserId: user})
	if err != nil {
		fmt.Println("ParamsRequest fails", err)
		return
	} else if paramsReply.GetError() == store.StoreProviderError_BAD_USER {
		fmt.Printf("ParamsRequest fails: user %q not found\n", user)
		return
	} else if paramsReply.GetError() != store.StoreProviderError_OK {
		fmt.Println("ParamsRequest fails:", paramsReply.GetError())
		return
	}

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
	fmt.Println("\nEnter the name of a website and we'll tell you your user name")
	fmt.Println("and password. Type \"quit\" to leave.")
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
			&store.ShareRequest{
				UserId: user,
				X:      int32(x),
				Y:      int32(y),
			},
		)
		if err != nil {
			fmt.Println("ShareRequest fails:", err)
			return
		} else if shareReply.GetError() != store.StoreProviderError_OK {
			fmt.Println("ShareRequest fails:", shareReply.GetError())
			return
		}

		out, err := priv.GetValue(in, shareReply.GetPubShare())
		if err == store.ItemNotFound {
			fmt.Println("Site does not exist or you entered the wrong master password.")
		} else if err != nil {
			fmt.Println("priv.GetValue() fails:", err)
			return
		} else {
			outs := strings.SplitN(out, ";", 2)
			if len(outs) != 2 {
				fmt.Printf("Huh ... the credentials were not properly formatted. Here's what we got: %q\n", out)
			} else {
				fmt.Println("User name:", outs[0])
				fmt.Println("Pqssword: ", outs[1])
			}
		}
	}
}
