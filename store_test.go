package store

import (
	"testing"
)

func TestNewStore(t *testing.T) {
	pub, priv, err := NewStore(GenerateKey(), goodM)
	if err != nil {
		t.Fatal("NewStore() fails:", err)
	}
	defer pub.Free()
	defer priv.Free()

	t.Log(pub.String())
}

// Test priv.GetIdx, pub.GetShare, and priv.GetValue().
func TestStoreGetIdxRowValue(t *testing.T) {
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
		pubShare, err := pub.GetShare(x, y)
		if err != nil {
			t.Errorf("pub.GetShare(%d, %d) fails: %s", x, y, err)
		}
		if pubShare != nil {
			out, err := priv.GetOutput(in, pubShare)
			if err != nil {
				t.Errorf("priv.GetOutput(%q, %q) fails: %s", in, pubShare, err)
			} else if out != val {
				t.Errorf("out = %q, expected %q", out, val)
			}
		}
	}
}

func TestNewPubStoreFromProto(t *testing.T) {
	pub1, priv1, err := NewStore(GenerateKey(), goodM)
	if err != nil {
		t.Fatalf("NewStore() fails: %s", err)
	}
	defer pub1.Free()
	defer priv1.Free()

	pub2 := NewPubStoreFromProto(pub1.GetProto())
	if pub2 != nil {
		defer pub2.Free()
	}

	AssertStringEqError(t, "pub2.ToString()",
		pub2.GetProto().String(), pub1.GetProto().String())
}

func TestNewPrivStore(t *testing.T) {
	K := GenerateKey()
	pub1, priv1, err := NewStore(K, goodM)
	if err != nil {
		t.Fatalf("NewStore() fails: %s", err)
	}
	defer pub1.Free()
	defer priv1.Free()

	priv2, err := NewPrivStore(K, priv1.GetParams())
	if err != nil {
		t.Fatalf("NewPrivStore() fails: %s", err)
	}
	defer priv2.Free()

	testInput := "cool"
	x1, y1, err := priv1.GetIdx(testInput)
	if err != nil {
		t.Errorf("priv1.GetIdx() fails: %s", err)
	}

	x2, y2, err := priv2.GetIdx(testInput)
	if err != nil {
		t.Errorf("priv2.GetIdx() fails: %s", err)
	}

	AssertIntEqError(t, "x1", x1, x2)
	AssertIntEqError(t, "y1", y1, y2)

	nonce := []byte("123456789abc")
	output1 := priv1.aead.Seal(nil, nonce, []byte(testInput), nil)
	output2 := priv2.aead.Seal(nil, nonce, []byte(testInput), nil)
	AssertStringEqError(t, "output1", string(output1), string(output2))
}
