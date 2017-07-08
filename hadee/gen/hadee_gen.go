// Copyright (c) 2017, Christopher Patton.
// All rights reserved.

// hadee_gen generates a sample store from a password. It outputs a file called
// store.pub.
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cjpatton/store"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/ssh/terminal"
)

var M = map[string]string{
	"cool": "guy",
}

func main() {
	log.Println("Please enter your super secret password:")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalln("terminal.ReadPassword() fails:", err)
	}

	ls := make([]string, 0)
	for in, out := range M {
		if in == "ls" {
			log.Printf("warning: overwriting item (%q, %q) in input", in, out)
		}
		ls = append(ls, in)
	}
	lsStr := ""
	for i := 0; i < len(ls); i++ {
		lsStr += fmt.Sprintf("%s\n", ls[i])
	}
	M["ls"] = lsStr

	K := store.DeriveKeyFromPassword(password, nil)
	pub, priv, err := store.NewStore(K, M)
	if err != nil {
		log.Fatalln("store.New() fails:", err)
	}
	priv.Free()
	defer pub.Free()

	pubString, err := proto.Marshal(pub.GetProto())
	if err != nil {
		log.Fatalln("pub.GetProto().Marshal() fails:", err)
	}

	if err := ioutil.WriteFile("store.pub", pubString, 0644); err != nil {
		log.Fatalln("Writing table fails:", err)
	}

	log.Println("Wrote store.pub.")
}
