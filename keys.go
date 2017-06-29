// TODO Clean up errors
// TODO Add docstrings to functions.
package keys

/* #cgo LDFLAGS: -lstruct -lcrypto
#include <struct/const.h>
#include <struct/dict.h>
#include "string.h"

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
	"crypto/rand"
	"errors"
	"unsafe"
)

// TODO Make these parameters instead of global constants?
const SaltBytes = 8
const TagBytes = 2
const MaxRowBytes = C.HASH_BYTES
const KeyBytes = C.HMAC_KEY_BYTES

type StoreParams struct {
	TableLen       int
	MaxOutputBytes int
	RowBytes       int
	Salt           []byte
}

// TODO GetRows, Getparms (refactor)
type PubStore struct {
	dict *C.cdict_t
}

// TODO GetKey, GetParams (refactor), GetIdx, GetValue
type PrivStore struct {
	key     []byte
	tinyCtx *C.tiny_ctx
	params  *C.dict_params_t
}

func NewStore(K []byte, M map[string]string) (*PubStore, *PrivStore, error) {

	// Check that K is the right length.
	if len(K) != KeyBytes {
		return nil, nil, errors.New("bad inputBytes")
	}

	pub := new(PubStore)
	priv := new(PrivStore)
	priv.key = K

	// Copy input/output pairs into C land.
	itemCt := C.int(len(M))
	inputs := C.new_str_list(itemCt)
	inputBytes := C.new_int_list(itemCt)
	outputs := C.new_str_list(itemCt)
	outputBytes := C.new_int_list(itemCt)
	defer C.free_str_list(inputs, itemCt)
	defer C.free_str_list(outputs, itemCt)
	defer C.free_int_list(inputBytes)
	defer C.free_int_list(outputBytes)

	maxOutputueBytes := 0
	i := C.int(0)
	for input, output := range M {
		if len(output) > maxOutputueBytes {
			maxOutputueBytes = len(output)
		}
		// NOTE C.CString() copies all the bytes of its input, even if it
		// encounters a null byte.
		C.set_str_list(inputs, i, C.CString(input))
		C.set_int_list(inputBytes, i, C.int(len(input)))
		C.set_str_list(outputs, i, C.CString(output))
		C.set_int_list(outputBytes, i, C.int(len(output)))
		i++
	}

	tableLen := C.dict_compute_table_length(C.int(len(M)))
	dict := C.dict_new(
		tableLen,
		C.int(maxOutputueBytes),
		C.int(TagBytes),
		C.int(SaltBytes))
	if dict == nil {
		return nil, nil, errors.New("dict_new")
	}
	defer C.dict_free(dict)

	priv.tinyCtx = C.tinyprf_new(tableLen)
	if priv.tinyCtx == nil {
		return nil, nil, errors.New("tinyprf_new")
	}

	cK := C.CString(string(K))
	defer C.free(unsafe.Pointer(cK))
	errNo := C.tinyprf_init(priv.tinyCtx, cK)
	if errNo != C.OK {
		priv.Free()
		return nil, nil, errors.New("tinyprf_init")
	}

	errNo = C.dict_create(
		dict, priv.tinyCtx, inputs, inputBytes, outputs, outputBytes, itemCt)
	if errNo != C.OK {
		priv.Free()
		return nil, nil, errors.New("dict_create")
	}

	pub.dict = C.dict_compress(dict)
	if pub.dict == nil {
		priv.Free()
		return nil, nil, errors.New("dict_compress")
	}

	// Copy parameters to priv.
	priv.params = new(C.dict_params_t)
	priv.params.table_length = pub.dict.params.table_length
	priv.params.max_value_bytes = pub.dict.params.max_value_bytes
	priv.params.tag_bytes = pub.dict.params.tag_bytes
	priv.params.row_bytes = pub.dict.params.row_bytes
	priv.params.salt_bytes = pub.dict.params.salt_bytes
	priv.params.salt = (*C.char)(C.malloc(C.size_t(pub.dict.params.salt_bytes + 1)))
	C.memcpy(unsafe.Pointer(priv.params.salt),
		unsafe.Pointer(pub.dict.params.salt),
		C.size_t(priv.params.salt_bytes))

	return pub, priv, nil
}

func NewStoreGenerateKey(M map[string]string) (*PubStore, *PrivStore, error) {
	K := make([]byte, KeyBytes)
	_, err := rand.Read(K)
	if err != nil {
		return nil, nil, errors.New("rand.Read")
	}
	return NewStore(K, M)
}

func (pub *PubStore) GetParams() *StoreParams {
	params := new(StoreParams)
	params.TableLen = int(pub.dict.params.table_length)
	params.MaxOutputBytes = int(pub.dict.params.max_value_bytes)
	params.RowBytes = int(pub.dict.params.row_bytes)
	params.Salt = C.GoBytes(
		unsafe.Pointer(pub.dict.params.salt), pub.dict.params.salt_bytes)
	return params
}

func (pub *PubStore) Free() {
	C.cdict_free(pub.dict)
}

func (priv *PrivStore) Free() {
	for i := 0; i < len(priv.key); i++ {
		priv.key[i] = 0
	}
	C.free(unsafe.Pointer(priv.params.salt))
	C.tinyprf_free(priv.tinyCtx)
}

// TODO Refactor.
func Get(pub *PubStore, priv *PrivStore, input string) (string, error) {
	cInput := C.CString(input)
	// FIXME Better way to do the following?
	cOutput := C.CString(string(make([]byte, pub.dict.params.max_value_bytes)))
	cOutputBytes := C.int(0)
	defer C.free(unsafe.Pointer(cInput))
	defer C.free(unsafe.Pointer(cOutput))
	errNo := C.cdict_get(
		pub.dict, priv.tinyCtx, cInput, C.int(len(input)), cOutput, &cOutputBytes)
	if errNo == C.ERR_DICT_BAD_KEY {
		return "", errors.New("item not found")
	} else if errNo != C.OK {
		return "", errors.New("cdict_get")
	}
	return C.GoStringN(cOutput, cOutputBytes), nil
}
