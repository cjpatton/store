package keys

import "testing"

var goodK = "1234123412341234"
var badK = "too short"

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

	st, err := NewStore(goodK, badM)
	if err == nil {
		t.Fatal("NewStore(goodK, badM) succeeds, expected error")
	}

	st, err = NewStore(badK, goodM)
	if err == nil {
		t.Fatal("NewStore(badK, goodM) succeeds, expected error")
	}

	st, err = NewStore(goodK, emptyM)
	if err == nil {
		t.Fatalf("NewStore(goodK, emptyM) succeeds, expected error")
	}

	st, err = NewStore(goodK, goodM)
	if err != nil {
		t.Fatalf("NewStore(goodK, goodM) fails: %s", err)
	} else if st.key != goodK {
		t.Error("st.key = %q, expected %q", st.key, goodK)
	}
	defer st.Free()
}

func TestNewStoreGenerateKey(t *testing.T) {
	st, err := NewStoreGenerateKey(goodM)
	if err != nil {
		t.Fatalf("NewStoreGenerateKey(goodM) fails: %s", err)
	}
	defer st.Free()
}

func TestGet(t *testing.T) {
	st, err := NewStore(goodK, goodM)
	if err != nil {
		t.Fatalf("NewStore(goodK, goodM) fails: %s", err)
	}
	defer st.Free()

	badKey := "tragically"
	val, err := st.Get(badKey)
	if err == nil {
		t.Error("st.Get(badKey) succeeded, expected error")
	}

	goodKey := "hip"
	expectedVal := "pizza"
	val, err = st.Get(goodKey)
	if err != nil {
		t.Error("st.Get(goodKey) erred, expected success")
	} else if val != expectedVal {
		t.Error("st.get(goodKey) = %q, expected %q", val, expectedVal)
	}
}
