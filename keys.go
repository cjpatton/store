package keys

// #cgo LDFLAGS: -lstruct -lcrypto
// #include <struct/dict.h>
import "C"
import (
	"errors"
)

// TODO Make these parameters, not global constants?
const SaltBytes = 16
const TagBytes = 2

type DictParams struct {
	tableLen      int
	maxValueBytes int
	rowBytes      int
	salt          []byte
}

type Dict struct {
	params DictParams
	cdict  *C.cdict_t
}

func NewDict(m map[string]string) (*Dict, error) {
	dict := new(Dict)
	dict.params.tableLen = int(C.dict_compute_table_length(C.int(len(m))))
	dict.params.maxValueBytes = 0
	for _, val := range m {
		if len(val) > dict.params.maxValueBytes {
			dict.params.maxValueBytes = len(val)
		}
	}

	x := C.dict_new(
		C.int(dict.params.tableLen),
		C.int(dict.params.maxValueBytes),
		C.int(TagBytes),
		C.int(SaltBytes))

	if x == nil {
		return nil, errors.New("rowBytes exceeds maximum")
	}
	defer C.dict_free(x)
	dict.params.rowBytes = int(x.params.row_bytes)
	return dict, nil
}

func (*Dict) Free() {

}

// TODO Remove this.
func ComputeTableLength(itemCt int) int {
	return int(C.dict_compute_table_length(C.int(itemCt)))
}
