// TODO Test a map with just one item and look at its internal representation.
package keys

import (
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

// TODO Does testing have these sorts of functions already.
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

func TestKey(t *testing.T) {
	AssertIntEqError(t, "len(GenerateKey()))", len(GenerateKey()), KeyBytes)

	password := []byte("hadi")
	salt := []byte("1947")
	K := DeriveKeyFromPassword(password, salt)
	AssertIntEqError(t, "len(DeriveKeyFromPassword()))", len(K), KeyBytes)
}

func TestNewStore(t *testing.T) {

	// Test with map with a value that is too long.
	pub, priv, err := NewStore(goodK, badM)
	if err == nil {
		t.Fatal("NewStore(goodK, badM) succeeds, expected error")
	}

	// Test with map with no items.
	pub, priv, err = NewStore(goodK, emptyM)
	if err == nil {
		t.Fatalf("NewStore(goodK, emptyM) succeeds, expected error")
	}

	// Test with key that is not the right length.
	pub, priv, err = NewStore(badK, goodM)
	if err == nil {
		t.Fatal("NewStore(badK, goodM) succeeds, expected error")
	}

	// Test with good inputs.
	pub, priv, err = NewStore(goodK, goodM)
	if err != nil {
		t.Fatalf("NewStore(goodK, goodM) fails: %s", err)
	}
	t.Log("pub\n", pub.ToString())
	defer pub.Free()
	defer priv.Free()

	// Check that the parameters are the same.
	AssertIntEqError(t, "priv.params.table_length", int(priv.params.table_length), int(pub.dict.params.table_length))
	AssertIntEqError(t, "priv.params.max_value_bytes", int(priv.params.max_value_bytes), int(pub.dict.params.max_value_bytes))
	AssertIntEqError(t, "priv.params.tag_bytes", int(priv.params.tag_bytes), int(pub.dict.params.tag_bytes))
	AssertIntEqError(t, "priv.params.row_bytes", int(priv.params.row_bytes), int(pub.dict.params.row_bytes))
	AssertIntEqError(t, "priv.params.salt_bytes", int(priv.params.salt_bytes), int(pub.dict.params.salt_bytes))
	AssertStringEqError(t, "priv.params.salt",
		cBytesToString(priv.params.salt, priv.params.salt_bytes),
		cBytesToString(pub.dict.params.salt, pub.dict.params.salt_bytes))

	pub1, priv1, err := NewStore(goodK, oneM)
	if err != nil {
		t.Fatalf("NewStore(goodK, goodM) fails: %s", err)
	}
	t.Log("pub1\n", pub1.ToString())
	defer pub1.Free()
	defer priv1.Free()
}

func TestGetParams(t *testing.T) {
	pub, priv, err := NewStore(GenerateKey(), goodM)
	if err != nil {
		t.Fatalf("NewStore() fails: %s", err)
	}
	defer pub.Free()
	defer priv.Free()

	params := pub.GetParams()
	if params == nil {
		t.Error("pub.GetParams() = nil, expected success")
	}
	AssertIntEqError(t, "pub.GetParams(): params.TableLen", params.TableLen, 9)
	AssertIntEqError(t, "pub.GetParams(): params.MaxOutputBytes", params.MaxOutputBytes, 5)
	AssertIntEqError(t, "pub.GetParams(): params.TagBytes", params.TagBytes, 2)
	AssertIntEqError(t, "pub.GetParams(): params.RowBytes", params.RowBytes, 8)
	AssertIntEqError(t, "pub.GetParams(): len(params.Salt)", len(params.Salt), SaltBytes)

	params = priv.GetParams()
	if params == nil {
		t.Error("pub.GetParams() = nil, expected success")
	}
	AssertIntEqError(t, "priv.GetParams(): params.TableLen", params.TableLen, 9)
	AssertIntEqError(t, "priv.GetParams(): params.MaxOutputBytes", params.MaxOutputBytes, 5)
	AssertIntEqError(t, "priv.GetParams(): params.TagBytes", params.TagBytes, 2)
	AssertIntEqError(t, "priv.GetParams(): params.RowBytes", params.RowBytes, 8)
	AssertIntEqError(t, "priv.GetParams(): len(params.Salt)", len(params.Salt), SaltBytes)
}

func TestGet(t *testing.T) {
	pub, priv, err := NewStore(GenerateKey(), goodM)
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

func TestGetIdxRowsValue(t *testing.T) {
	pub, priv, err := NewStore(GenerateKey(), goodM)
	if err != nil {
		t.Fatalf("NewStore() fails: %s", err)
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
