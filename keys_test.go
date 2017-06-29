// TODO Test that priv.params matches pub.dict.params.
// TODO Test that Get() is the same as GetIdx(), GetRows(), GetValue().
// TODO Test a map with just one item and look at its internal representation.
package keys

import (
	"bytes"
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

var badM = map[string]string{
	"cool": string(make([]byte, MaxRowBytes-TagBytes)),
}

func TestNewStore(t *testing.T) {
	pub, priv, err := NewStore(goodK, badM)
	if err == nil {
		t.Fatal("NewStore(goodK, badM) succeeds, expected error")
	}

	pub, priv, err = NewStore(badK, goodM)
	if err == nil {
		t.Fatal("NewStore(badK, goodM) succeeds, expected error")
	}

	pub, priv, err = NewStore(goodK, emptyM)
	if err == nil {
		t.Fatalf("NewStore(goodK, emptyM) succeeds, expected error")
	}

	pub, priv, err = NewStore(goodK, goodM)
	if err != nil {
		t.Fatalf("NewStore(goodK, goodM) fails: %s", err)
	}
	defer pub.Free()
	defer priv.Free()

	if bytes.Compare(priv.key, goodK) != 0 {
		t.Errorf("priv.Key = %q, expected %q", priv.key, goodK)
	}
}

func TestNewStoreGenerateKey(t *testing.T) {
	pub, priv, err := NewStoreGenerateKey(goodM)
	if err != nil {
		t.Fatalf("NewStoreGenerateKey(goodM) fails: %s", err)
	}
	defer pub.Free()
	defer priv.Free()
}

func TestGetParams(t *testing.T) {
	pub, priv, err := NewStoreGenerateKey(goodM)
	if err != nil {
		t.Fatalf("NewStoreGenerateKey() fails: %s", err)
	}
	defer pub.Free()
	defer priv.Free()

	params := pub.GetParams()
	if params == nil {
		t.Error("st.GetParams() = nil, expected success")
	}

	// TODO Do something like AssertEq(var, exp, got)?
	expectedTableLen := 9
	if params.TableLen != expectedTableLen {
		t.Errorf("params.TableLen = %d, expected %d", params.TableLen, expectedTableLen)
	}

	expectedMaxOutputBytes := 5
	if params.MaxOutputBytes != expectedMaxOutputBytes {
		t.Errorf("params.MaxOutputBytes = %d, expected %d", params.MaxOutputBytes, expectedMaxOutputBytes)
	}

	expectedRowBytes := 8
	if params.RowBytes != expectedRowBytes {
		t.Errorf("params.RowBytes = %d, expected %d", params.RowBytes, expectedRowBytes)
	}

	expectedSaltLen := SaltBytes
	if len(params.Salt) != expectedSaltLen {
		t.Errorf("len(params.Salt) = %d, expected %d", len(params.Salt), expectedSaltLen)
	}
}

func TestGet(t *testing.T) {
	pub, priv, err := NewStoreGenerateKey(goodM)
	if err != nil {
		t.Fatalf("NewStore() fails: %s", err)
	}
	defer pub.Free()
	defer priv.Free()

	badInput := "tragically"
	output, err := Get(pub, priv, badInput)
	if err == nil {
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
