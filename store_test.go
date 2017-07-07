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
