// Copyright (c) 2017, Christopher Patton.
// All rights reserved.
//
// hadee_gen generates a new table from a password. It outputs a file called
// table.pub.
package main

import (
	"io/ioutil"
	"log"

	"github.com/cjpatton/store"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/ssh/terminal"
)

var credentials = map[string]string{
	"google.com":   "hadi@hotmail.com;hadi12345",
	"hotmail.com":  "hadi@hotmail.com;1947hadi",
	"yahoo.com":    "steve.bannion@netzero.gov;123456789abcdef",
	"myspace.com":  "pinkcatgirl@hotmail.com;Bradly Cooper is so cute!",
	"linkedin.com": "securityexpert@yahoo.com;s3kr1te3xp3rt!",
	"zombo.com":    "hadi@hotmail.com;1;2;3;4",
}

func main() {
	log.Println("Please enter your super secret password:")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalln("terminal.ReadPassword(0) fails:", err)
	}

	key := store.DeriveKeyFromPassword(password, nil)
	pub, priv, err := store.New(key, credentials)
	if err != nil {
		log.Fatalln("store.New(key, credentials) fails:", err)
	}
	priv.Free()
	defer pub.Free()

	pubString, err := proto.Marshal(pub.GetTable())
	if err != nil {
		log.Fatalln("pub.GetTable().Marshal() fails:", err)
	}

	if err := ioutil.WriteFile("table.pub", pubString, 0644); err != nil {
		log.Fatalln("Writing table fails:", err)
	}

	log.Println("Wrote table.pub.")
}
