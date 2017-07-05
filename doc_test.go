// Copyright (c) 2017, Christopher Patton
// All rights reserved.

package store

import "fmt"

func ExampleVDeriveKeyFromPassword() {
	password := []byte("A really secure password")
	salt := []byte("Optional salt, useful in many applications")
	K := DeriveKeyFromPassword(password, salt)
	fmt.Println(len(K))
	// Output: 16
}

func ExampleGenerateKey() {
	K := GenerateKey()
	fmt.Println(len(K))
	// Output: 16
}

func ExampleNew() {
	K := GenerateKey()
	M := map[string]string{"Out": "of this world!"}

	pub, priv, err := New(K, M)
	if err != nil {
		fmt.Println("New() error:", err)
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
	K := GenerateKey()
	M := map[string]string{"Out": "of this world!"}

	pub, priv, err := New(K, M)
	if err != nil {
		fmt.Println("New() error:", err)
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

func ExampleNewPubStoreFromTable() {
	K := GenerateKey()
	M := map[string]string{"Out": "of this world!"}

	pub, priv, err := New(K, M)
	if err != nil {
		fmt.Println("New() error:", err)
		return
	}
	defer pub.Free()
	defer priv.Free()

	pubFromTable := NewPubStoreFromTable(pub.GetTable())
	defer pubFromTable.Free()

	fmt.Println(pub.ToString() == pubFromTable.ToString())
	// Output: true
}

func ExampleNewPrivStore() {
	K := GenerateKey()
	M := map[string]string{"Out": "of this world!"}

	pub, priv, err := New(K, M)
	if err != nil {
		fmt.Println("New() error:", err)
		return
	}
	defer pub.Free()
	defer priv.Free()

	privFromKeyAndPrivParams, err := NewPrivStore(K, priv.GetParams())
	if err != nil {
		fmt.Println("NewPrivStore() error:", err)
	}
	defer privFromKeyAndPrivParams.Free()

	privFromKeyAndPubParams, err := NewPrivStore(K, pub.GetTable().GetParams())
	if err != nil {
		fmt.Println("NewPrivStore() error:", err)
	}
	defer privFromKeyAndPubParams.Free()
}
