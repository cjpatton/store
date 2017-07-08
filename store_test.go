package store

import (
	"testing"

	"encoding/binary"
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

	badIn := "Nooooooo"
	x, y, err := priv.GetIdx(badIn)
	if err != nil {
		t.Errorf("priv.GetIdx(%q) fails: %s", badIn, err)
	}
	pubShare, err := pub.GetShare(x, y)
	if err == nil { // GetShare might succeed ...
		out, err := priv.GetOutput(badIn, pubShare)
		if err == nil { // ... But this should fail!
			t.Errorf("priv.GetOutput(%q, %q) succeeds, expected failure: %s", out)
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

// Evaluate how frequently pub.GetShare() returns a ciphertext on an incorrect
// query. On this small test, It appears to do so about 10% of the time.
func TestBadGetShare(t *testing.T) {
	K := GenerateKey()
	M := map[string]string{
		"This":             "is pretty cool, I think.",
		"It":               "is basically a mechanism for searchable encryption.",
		"I wonder":         "if anyone will use it?",
		"If they don't":    "I will",
		"You can use it":   "to store a bunch of files in an immutable fashion.",
		"The graph":        "is kept around for two reason: (1) to make it so the sealed output can be fetched in one round of communication",
		"and (2)":          "So that the sealed outputs can be modified.",
		"This would mean":  "That the AEAD needs to be nonce misuse resistant, since we're using the same nonce to store a different message",
		"I'll think about": "all these things moving forward. This has",
		"been":             "fun!",
	}

	pub, priv, err := NewStore(K, M)
	if err != nil {
		t.Fatal("NewStore() fails:", err)
	}
	defer pub.Free()
	defer priv.Free()
	t.Log(pub.String())

	trials := 1000
	ct := 0
	badInput := make([]byte, 4)
	for trial := 0; trial < trials; trial++ {
		binary.LittleEndian.PutUint32(badInput, uint32(trial))
		x, y, err := priv.GetIdx(string(badInput))
		if err != nil {
			t.Error("priv.GetIdx() fails:", err)
			return
		}

		_, err = pub.GetShare(x, y)
		if err == nil {
			ct++
		} else if err != ItemNotFound {
			t.Error("pub.GetShare() fails:", err)
			return
		}
	}
	t.Logf("%d / %d", ct, trials)
}
