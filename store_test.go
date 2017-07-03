// Copyright (c) 2017, Christopher Patton.
// All rights reserved.
package store

import (
	"fmt"
	"testing"
)

var goodK = []byte("1234123412341234")
var badK = []byte("too short")

var goodM = map[string]string{
	"this":           "I",
	"is":             "like",
	"pretty\x00cool": "",
	"hip":            "pizza",
}
var emptyM = map[string]string{}
var oneM = map[string]string{
	"just": "one",
}
var badM = map[string]string{
	"cool": string(make([]byte, MaxRowBytes-TagBytes)),
}

// TODO Does testing have these sorts of functions already?
func AssertInt32EqError(t *testing.T, name string, got, exp int32) {
	if got != exp {
		t.Errorf("%s = %d, expected %d", name, got, exp)
	}
}

func AssertIntEqError(t *testing.T, name string, got, exp int) {
	if got != exp {
		t.Errorf("%s = %d, expected %d", name, got, exp)
	}
}

func AssertStringEqError(t *testing.T, name string, got, exp string) {
	if got != exp {
		t.Errorf("%s = %q, expected %q", name, got, exp)
	}
}

// Test that key generation and generation output a key of the correct length.
func TestKey(t *testing.T) {
	AssertIntEqError(t, "len(GenerateKey()))", len(GenerateKey()), KeyBytes)
	password := []byte("hadi")
	salt := []byte("1947")
	K := DeriveKeyFromPassword(password, salt)
	AssertIntEqError(t, "len(DeriveKeyFromPassword()))", len(K), KeyBytes)
}

// Test New() good and bad inputs.
func TestNew(t *testing.T) {

	// Test with map with a value that is too long.
	pub, priv, err := New(goodK, badM)
	if err == nil {
		t.Fatal("New(goodK, badM) succeeds, expected error")
	}

	// Test with map with no items.
	pub, priv, err = New(goodK, emptyM)
	if err == nil {
		t.Fatalf("New(goodK, emptyM) succeeds, expected error")
	}

	// Test with key that is not the right length.
	pub, priv, err = New(badK, goodM)
	if err == nil {
		t.Fatal("New(badK, goodM) succeeds, expected error")
	}

	// Test with good inputs.
	pub, priv, err = New(goodK, goodM)
	if err != nil {
		t.Fatalf("New(goodK, goodM) fails: %s", err)
	}
	t.Log("pub\n", pub.ToString())
	defer pub.Free()
	defer priv.Free()

	// Check that the parameters are the same.
	AssertInt32EqError(t, "priv.params.table_length", int32(priv.params.table_length), int32(pub.dict.params.table_length))
	AssertInt32EqError(t, "priv.params.max_value_bytes", int32(priv.params.max_value_bytes), int32(pub.dict.params.max_value_bytes))
	AssertInt32EqError(t, "priv.params.tag_bytes", int32(priv.params.tag_bytes), int32(pub.dict.params.tag_bytes))
	AssertInt32EqError(t, "priv.params.row_bytes", int32(priv.params.row_bytes), int32(pub.dict.params.row_bytes))
	AssertInt32EqError(t, "priv.params.salt_bytes", int32(priv.params.salt_bytes), int32(pub.dict.params.salt_bytes))
	AssertStringEqError(t, "priv.params.salt",
		cBytesToString(priv.params.salt, priv.params.salt_bytes),
		cBytesToString(pub.dict.params.salt, pub.dict.params.salt_bytes))

	pub1, priv1, err := New(goodK, oneM)
	if err != nil {
		t.Fatalf("New(goodK, goodM) fails: %s", err)
	}
	t.Log("pub1\n", pub1.ToString())
	defer pub1.Free()
	defer priv1.Free()
}

// Test pub.GetParams() and priv.GetParams().
func TestGetParams(t *testing.T) {
	pub, priv, err := New(GenerateKey(), goodM)
	if err != nil {
		t.Fatalf("New() fails: %s", err)
	}
	defer pub.Free()
	defer priv.Free()

	params := pub.GetTable().GetParams()
	if params == nil {
		t.Error("pub.GetParams() = nil, expected success")
	}
	AssertInt32EqError(t, "pub.GetParams(): params.TableLen", params.GetTableLen(), 9)
	AssertInt32EqError(t, "pub.GetParams(): params.MaxOutputBytes", params.GetMaxOutputBytes(), 5)
	AssertInt32EqError(t, "pub.GetParams(): params.TagBytes", params.GetTagBytes(), 3)
	AssertInt32EqError(t, "pub.GetParams(): params.RowBytes", params.GetRowBytes(), 9)
	AssertIntEqError(t, "pub.GetParams(): len(params.Salt)", len(params.Salt), SaltBytes)

	params = priv.GetParams()
	if params == nil {
		t.Error("pub.GetParams() = nil, expected success")
	}
	AssertInt32EqError(t, "priv.GetParams(): params.TableLen", params.GetTableLen(), 9)
	AssertInt32EqError(t, "priv.GetParams(): params.MaxOutputBytes", params.GetMaxOutputBytes(), 5)
	AssertInt32EqError(t, "priv.GetParams(): params.TagBytes", params.GetTagBytes(), 3)
	AssertInt32EqError(t, "priv.GetParams(): params.RowBytes", params.GetRowBytes(), 9)
	AssertIntEqError(t, "priv.GetParams(): len(params.Salt)", len(params.Salt), SaltBytes)
}

// Test Get().
func TestGet(t *testing.T) {
	pub, priv, err := New(GenerateKey(), goodM)
	if err != nil {
		t.Fatalf("New() fails: %s", err)
	}
	defer pub.Free()
	defer priv.Free()

	badInput := "tragically"
	output, err := Get(pub, priv, badInput)
	if err != ItemNotFound {
		t.Error("st.Get(badInput) succeeded, expected error")
	}

	goodInput := "hip"
	expectedOutput := "pizza"
	output, err = Get(pub, priv, goodInput)
	if err != nil {
		t.Error("st.Get(goodInput) erred, expected success")
	} else if output != expectedOutput {
		t.Error("st.get(goodInput) = %q, expected %q", output, expectedOutput)
	}
}

// Test priv.GetIdx, pub.GetRows, and priv.GetValue().
func TestGetIdxRowsValue(t *testing.T) {
	pub, priv, err := New(GenerateKey(), goodM)
	if err != nil {
		t.Fatalf("New() fails: %s", err)
	}
	defer pub.Free()
	defer priv.Free()

	for in, val := range goodM {
		x, y, err := priv.GetIdx(in)
		if err != nil {
			t.Errorf("priv.GetIdx(%q) fails: %s", in, err)
		}
		X, err := pub.GetRow(x)
		if err != nil {
			t.Errorf("pub.GetRow(%d) fails: %s", x, err)
		}
		Y, err := pub.GetRow(y)
		if err != nil {
			t.Errorf("pub.GetRow(%d) fails: %s", y, err)
		}
		if X != nil && Y != nil {
			rows := [][]byte{X, Y}
			out, err := priv.GetValue(in, rows)
			if err != nil {
				t.Errorf("priv.GetValue(%q, %q) fails: %s", in, rows, err)
			} else if out != val {
				t.Errorf("out = %q, expected %q", out, val)
			}
		}
	}
}

// Test NewPubStoreFromTable().
func TestNewPubStoreFromTable(t *testing.T) {
	pub, priv, err := New(goodK, goodM)
	if err != nil {
		t.Fatalf("New(goodK, goodM) fails: %s", err)
	}
	defer pub.Free()
	defer priv.Free()

	pub2 := NewPubStoreFromTable(pub.GetTable())
	defer pub2.Free()

	AssertStringEqError(t, "pub2.ToString()", pub2.ToString(), pub.ToString())

	for in, val := range goodM {
		out2, err := Get(pub2, priv, in)
		if err != nil {
			t.Fatalf("Get(pub2, priv, %q) fails: %s", in, err)
		}
		AssertStringEqError(t, fmt.Sprintf("Get(pub2, priv, %q)", in), out2, val)
	}
}
