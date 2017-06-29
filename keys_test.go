package keys

import "testing"

func TestComputeTableLength(t *testing.T) {
	var itemCt int = 100
	tableLength := ComputeTableLength(itemCt)
	if tableLength != 209 {
		t.Errorf("ComputeTableLength(%d) = %d, want %q", itemCt, tableLength, "hella")
	}
}

func TestNewStore(t *testing.T) {

	goodM := map[string]string{
		"this":   "I",
		"is":     "like",
		"pretty": "",
		"hip":    "pizza",
	}

	badM := map[string]string{
		"cool": string(make([]byte, MaxRowBytes-TagBytes)),
	}

	goodK := []byte("1234123412341234")

	badK := []byte("too short")

	st, err := NewStore(goodK, badM)
	if err == nil {
		t.Fatal("CreateDict(goodK, badM) succeeds, expected error")
	}

	st, err = NewStore(badK, goodM)
	if err == nil {
		t.Fatal("CreateDict(badK, goodM) succeeds, expected error")
	}

	st, err = NewStore(goodK, goodM)
	if err != nil {
		t.Fatalf("CreateDict(goodK, goodM) fails: %q", err)
	}

	defer st.Free()
}
