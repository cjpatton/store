package keys

// #cgo LDFLAGS: -lstruct -lcrypto
// #include <struct/const.h>
// #include <struct/dict.h>
import "C"
import (
	"errors"
)

// TODO Make these parameters, not global constants?
const SaltBytes = 16
const TagBytes = 2
const MaxRowBytes = C.HASH_BYTES

type StoreParams struct {
	tableLen      int
	maxValueBytes int
	rowBytes      int
	salt          []byte
}

type Store struct {
	tinyCtx *C.tiny_ctx
	dict    *C.cdict_t
}

func NewStore(K []byte, M map[string]string) (*Store, error) {

	if len(K) != C.HMAC_KEY_BYTES {
		return nil, errors.New("bad keyBytes")
	}

	st := new(Store)
	tableLen := int(C.dict_compute_table_length(C.int(len(M))))
	maxValueBytes := 0
	for _, val := range M {
		if len(val) > maxValueBytes {
			maxValueBytes = len(val)
		}
	}

	dict := C.dict_new(
		C.int(tableLen),
		C.int(maxValueBytes),
		C.int(TagBytes),
		C.int(SaltBytes))

	if dict == nil {
		return nil, errors.New("bad rowBytes")
	}

	defer C.dict_free(dict)
	return st, nil
}

func (*Store) Free() {

}

// TODO Remove this.
func ComputeTableLength(itemCt int) int {
	return int(C.dict_compute_table_length(C.int(itemCt)))
}
