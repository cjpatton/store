// TODO Clean up errors
package keys

/* #cgo LDFLAGS: -lstruct -lcrypto
#include <struct/const.h>
#include <struct/dict.h>

char **new_str_list(int len) {
	return calloc(sizeof(char *), len);
}

int *new_int_list(int len) {
	return calloc(sizeof(int), len);
}

void set_str_list(char **list, int idx, char *val) {
	list[idx] = val;
}

void set_int_list(int *list, int idx, int val) {
	list[idx] = val;
}

void free_str_list(char **list, int len) {
	int i;
	for (i = 0; i < len; i++) {
		if (list[i] != NULL) {
			free(list[i]);
		}
	}
	free(list);
}

void free_int_list(int *list) {
	free(list);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

// TODO Make these parameters instead of global constants?
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

func NewStore(K string, M map[string]string) (*Store, error) {

	// Check that K is the right length.
	if len(K) != C.HMAC_KEY_BYTES {
		return nil, errors.New("bad keyBytes")
	}

	st := new(Store)

	// Copy key/value pairs into C land.
	itemCt := C.int(len(M))
	keys := C.new_str_list(itemCt)
	keyBytes := C.new_int_list(itemCt)
	values := C.new_str_list(itemCt)
	valueBytes := C.new_int_list(itemCt)
	defer C.free_str_list(keys, itemCt)
	defer C.free_str_list(values, itemCt)
	defer C.free_int_list(keyBytes)
	defer C.free_int_list(valueBytes)

	maxValueBytes := 0
	i := C.int(0)
	for key, value := range M {
		if len(value) > maxValueBytes {
			maxValueBytes = len(value)
		}
		// NOTE C.CString() copies all the bytes of its input, even if it
		// encounters a null byte.
		C.set_str_list(keys, i, C.CString(key))
		C.set_int_list(keyBytes, i, C.int(len(key)))
		C.set_str_list(values, i, C.CString(value))
		C.set_int_list(valueBytes, i, C.int(len(value)))
		i++
	}

	tableLen := C.dict_compute_table_length(C.int(len(M)))
	dict := C.dict_new(
		tableLen,
		C.int(maxValueBytes),
		C.int(TagBytes),
		C.int(SaltBytes))
	if dict == nil {
		return nil, errors.New("dict_new")
	}
	defer C.dict_free(dict)

	st.tinyCtx = C.tinyprf_new(tableLen)
	if st.tinyCtx == nil {
		return nil, errors.New("tinyprf_new")
	}

	cK := C.CString(K)
	defer C.free(unsafe.Pointer(cK))
	errNo := C.tinyprf_init(st.tinyCtx, cK)
	if errNo != C.OK {
		st.Free()
		return nil, errors.New("tinyprf_init")
	}

	errNo = C.dict_create(
		dict, st.tinyCtx, keys, keyBytes, values, valueBytes, itemCt)
	if errNo != C.OK {
		st.Free()
		return nil, errors.New("dict_create")
	}

	st.dict = C.dict_compress(dict)
	if st.dict == nil {
		st.Free()
		return nil, errors.New("dict_compress")
	}

	return st, nil
}

func (st *Store) Get(key string) (string, error) {
	cKey := C.CString(key)
	// FIXME Better way to do the following?
	cVal := C.CString(string(make([]byte, st.dict.params.max_value_bytes)))
	cValBytes := C.int(0)
	defer C.free(unsafe.Pointer(cKey))
	defer C.free(unsafe.Pointer(cVal))
	errNo := C.cdict_get(
		st.dict, st.tinyCtx, cKey, C.int(len(key)), cVal, &cValBytes)
	if errNo == C.ERR_DICT_BAD_KEY {
		return "", errors.New("item not found")
	} else if errNo != C.OK {
		return "", errors.New("cdict_get")
	}
	return C.GoStringN(cVal, cValBytes), nil
}

func (st *Store) Free() {
	C.tinyprf_free(st.tinyCtx)
	C.cdict_free(st.dict)
}
