package keys

import "testing"

func TestAreWeHavingFunYet(t *testing.T) {
	response := AreWeHavingFunYet()
	if response != "yeah!" {
		t.Errorf("AreWehavingFunYet() = %q, want %q", response, "yeah!")
	}
}
