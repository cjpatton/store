package keys

import "testing"

func TestComputeTableLength(t *testing.T) {
	var itemCt int = 100
	tableLength := ComputeTableLength(itemCt)
	if tableLength != 209 {
		t.Errorf("ComputeTableLength(%d) = %d, want %q", itemCt, tableLength, "hella")
	}
}
