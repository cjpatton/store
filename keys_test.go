package keys

import "testing"

func TestComputeTableLength(t *testing.T) {
	var itemCt int = 100
	tableLength := ComputeTableLength(itemCt)
	if tableLength != 209 {
		t.Errorf("ComputeTableLength(%d) = %d, want %q", itemCt, tableLength, "hella")
	}
}

func TestCreateDict(t *testing.T) {
	m := map[string]string{
		"this":   "I",
		"is":     "like",
		"pretty": "",
		"hip":    "pizza",
	}

	dict, err := NewDict(m)
	if err != nil {
		t.Fatalf("CreateDict(m) fails: %q", err)
	}

	defer dict.Free()
}
