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
