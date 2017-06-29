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

func TestGetparams(t *testing.T) {
	st, err := NewStore(goodK, goodM)
	if err != nil {
		t.Fatalf("NewStore() fails: %s", err)
	}
	defer st.Free()

	params := st.GetParams()
	if params == nil {
		t.Error("st.GetParams() = nil, expected success")
	}

	// TODO Do something like AssertEq(var, exp, got)?
	expectedTableLen := 9
	if params.TableLen != expectedTableLen {
		t.Errorf("params.TableLen = %d, expected %d", params.TableLen, expectedTableLen)
	}

	expectedMaxOutputBytes := 5
	if params.MaxOutputBytes != expectedMaxOutputBytes {
		t.Errorf("params.MaxOutputBytes = %d, expected %d", params.MaxOutputBytes, expectedMaxOutputBytes)
	}

	expectedRowBytes := 8
	if params.RowBytes != expectedRowBytes {
		t.Errorf("params.RowBytes = %d, expected %d", params.RowBytes, expectedRowBytes)
	}

	expectedSaltLen := SaltBytes
	if len(params.Salt) != expectedSaltLen {
		t.Errorf("len(params.Salt) = %d, expected %d", len(params.Salt), expectedSaltLen)
	}
}

func TestGet(t *testing.T) {
	st, err := NewStore(goodK, goodM)
	if err != nil {
		t.Fatalf("NewStore() fails: %s", err)
	}
	defer st.Free()

	badInput := "tragically"
	output, err := st.Get(badInput)
	if err == nil {
		t.Error("st.Get(badInput) succeeded, expected error")
	}

	goodInput := "hip"
	expectedOutput := "pizza"
	output, err = st.Get(goodInput)
	if err != nil {
		t.Error("st.Get(goodInput) erred, expected success")
	} else if output != expectedOutput {
		t.Error("st.get(goodInput) = %q, expected %q", output, expectedOutput)
	}
}
