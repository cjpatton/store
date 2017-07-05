// Copyright (c) 2017, Christopher Patton.
// All rights reserved.

// hadee_gen generates a sample store from a password. It outputs a file called
// store.pub.
package main

import (
	"io/ioutil"
	"log"

	"github.com/cjpatton/store"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/ssh/terminal"
)

var credentials = map[string]string{
	"pwned.org":                                                                                                                "someemail@kewl.org;L33tHaXor:25",
	"russia.us":                                                                                                                "vladimir@putin;Putin-on-the-ritz",
	"dell.com":                                                                                                                 "del@funkyhomosapients.com;This is quite funky",
	"happytreefriends.com":                                                                                                     "frank@zapa.gov;280jklsdjf89wfdsfjkl1234234!!:D",
	"hadi.com":                                                                                                                 "me@hadi.com;1947",
	"wikipedia.com":                                                                                                            "me@hadi.com;1947",
	"wikipedia.org":                                                                                                            "me@hadi.com;1947",
	"wikimedia.com":                                                                                                            "you@hadi.com;hadi",
	"ufl.edu":                                                                                                                  "hadi@ufl.edu;gatorgatorgatorgatorgatorgator",
	"ucadvis.edu":                                                                                                              "hadi@hotmail.com;1947hadi",
	"microsoft.com":                                                                                                            "billy@gates.com;Oh, billy oh billy boy! Where are you?",
	"www.doglovers.com":                                                                                                        "benjirodriguez@netzero.com;ruff",
	"cise.ufl.edu":                                                                                                             "hacker97@hotmail.com;When I heard the learn'd astronomber",
	"google.com":                                                                                                               "hadi@hotmail.com;hadi12345",
	"hotmail.com":                                                                                                              "hadi@hotmail.com;1947hadi",
	"yahoo.com":                                                                                                                "steve.bannion@netzero.gov;123456789abcdef",
	"myspace.com":                                                                                                              "pinkcatgirl@hotmail.com;OMG Bradly Cooper is so cute!",
	"linkedin.com":                                                                                                             "securityexpert@yahoo.com;s3kr1te3xp3rt!",
	"zombo.com":                                                                                                                "hadi@hotmail.com;1;2;3;4",
	"caves.org":                                                                                                                "splunker1990@netzero.com;cavesrkewl",
	"facebook.com":                                                                                                             "roberoppenheimer@morenukes.com;Duke Nuke 'em",
	"centerforants.org":                                                                                                        "derek@bluesteel.com;BalooSteal",
	"centerforkidswhocantreadgoodandwannalearntodootherstuffgoodtoo.org":                                                       "derek@bluesteel.com;LeTigre",
	"Actually the input string can be any length ... even this really long strong works, as long as the output is < 60 bytes.": ";",
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

	if err := ioutil.WriteFile("store.pub", pubString, 0644); err != nil {
		log.Fatalln("Writing table fails:", err)
	}

	log.Println("Wrote store.pub.")
}
