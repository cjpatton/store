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

var ww = "When I heard the learn'd astronomer,\nWhen the proofs, the figures, were ranged in columns before me,\nWhen I was shown the charts and diagrams, to add, divide, and measure them,\nWhen I sitting heard the astronomer where he lectured with much applause in the lecture-room,\nHow soon unaccountable I became tired and sick,\nTill rising and gliding out I wander'd off by myself,\nIn the mystical moist night-air, and from time to time,\nLook'd up in perfect silence at the stars."

// Source: http://www.ducksters.com/jokes/dogs.php
var M = map[string]string{
	"Why don't blind people like to sky dive?":       "Because it scares the dog!",
	"Why did the poor dog chase his own tail?":       "He was trying to make both ends meet!",
	"What dog keeps the best time?":                  "A watch dog!",
	"Why don't dogs make good dancers?":              "Because they have two left feet!",
	"What happens when it rains cats and dogs?":      "You can step in a poodle!",
	"Why are dogs like phones?":                      "Because they have collar IDs.",
	"What dog loves to take bubble baths?":           "A shampoodle!",
	"What did the dog say when he sat on sandpaper?": "Ruff!",
	"What do you call a dog that is left handed?":    "A south paw!",
	"What did one flea say to the other?":            "Should we walk or take a dog?",
	"What type of markets do dogs avoid?":            "Flea markets!",
	"What did the cowboy say when his dog ran away?": "Well, doggone!",
	"When I heard ...":                               ww,
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

	log.Println("The store:")
	log.Println("\n", pub.String())

	pubString, err := proto.Marshal(pub.GetProto())
	if err != nil {
		log.Fatalln("pub.GetProto().Marshal() fails:", err)
	}

	if err := ioutil.WriteFile("store.pub", pubString, 0644); err != nil {
		log.Fatalln("Writing table fails:", err)
	}

	log.Println("Wrote store.pub.")
}
