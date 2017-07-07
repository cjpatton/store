// Copyright (c) 2017, Christopher Patton
// All rights reserved.

package store

import "fmt"

func ExampleDeriveKeyFromPassword() {
	password := []byte("A really secure password")
	salt := []byte("Optional salt, useful in many applications")
	K := DeriveKeyFromPassword(password, salt)
	fmt.Println(len(K))
	// Output: 32
}

func ExampleGenerateDictKey() {
	K := GenerateKey()
	fmt.Println(len(K))
	// Output: 32
}

func ExampleNewDict() {
	K := GenerateDictKey()
	M := map[string]string{"Out": "of this world!"}

	pub, priv, err := NewDict(K, M)
	if err != nil {
		fmt.Println("NewDict() error:", err)
		return
	}
	defer pub.Free()
	defer priv.Free()

	x, y, err := priv.GetIdx("Out")
	if err != nil {
		fmt.Println("priv.GetIdx() error:", err)
		return
	}

	pubShare, err := pub.GetShare(x, y)
	if err != nil {
		fmt.Println("pub.GetShare() error:", err)
		return
	}

	out, err := priv.GetValue("Out", pubShare)
	if err != nil {
		fmt.Println("priv.GetValue() error:", err)
		return
	}

	fmt.Println(out)
	// Output: of this world!
}

func ExampleGet() {
	K := GenerateDictKey()
	M := map[string]string{"Out": "of this world!"}

	pub, priv, err := NewDict(K, M)
	if err != nil {
		fmt.Println("NewDict() error:", err)
		return
	}
	defer pub.Free()
	defer priv.Free()

	out, err := Get(pub, priv, "Out")
	if err != nil {
		fmt.Println("Get() error:", err)
		return
	}
	fmt.Println(out)

	out, err = Get(pub, priv, "Evil input")
	if err != nil {
		fmt.Println("Get() error:", err)
		return
	}
	// Output:
	// of this world!
	// Get() error: item not found
}

func ExampleNewPubDictFromTable() {
	K := GenerateDictKey()
	M := map[string]string{"Out": "of this world!"}

	pub, priv, err := NewDict(K, M)
	if err != nil {
		fmt.Println("NewDict() error:", err)
		return
	}
	defer pub.Free()
	defer priv.Free()

	pubFromTable := NewPubDictFromProto(pub.GetProto())
	defer pubFromTable.Free()

	fmt.Println(pub.String() == pubFromTable.String())
	// Output: true
}

func ExampleNewPrivDict() {
	K := GenerateDictKey()
	M := map[string]string{"Out": "of this world!"}

	pub, priv, err := NewDict(K, M)
	if err != nil {
		fmt.Println("NewDict() error:", err)
		return
	}
	defer pub.Free()
	defer priv.Free()

	privFromKeyAndPrivParams, err := NewPrivDict(K, priv.GetParams())
	if err != nil {
		fmt.Println("NewPrivDict() error:", err)
	}
	defer privFromKeyAndPrivParams.Free()

	privFromKeyAndPubParams, err := NewPrivDict(K, pub.GetProto().GetParams())
	if err != nil {
		fmt.Println("NewPrivDict() error:", err)
	}
	defer privFromKeyAndPubParams.Free()
}
