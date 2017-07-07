// Copyright (c) 2017, Christopher Patton.
// All rights reserved.

package store

import (
	"crypto/rand"
	"fmt"
	"unsafe"

	"github.com/cjpatton/store/pb"
	"github.com/golang/protobuf/proto"
)

/*
// The next line gets things going on Mac:
#cgo CPPFLAGS: -I/usr/local/opt/openssl/include
#cgo LDFLAGS: -lstruct -lcrypto
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

int get_int_list(int *list, int idx) {
	return list[idx];
}

char *get_str_list(char **list, int idx) {
	return list[idx];
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

char *get_row_ptr(char *table, int row, int row_bytes) {
	return &table[row * row_bytes];
}

node_t *get_node(graph_t *graph, int idx) {
	return &graph->node[idx];
}

int get_edge(graph_t *graph, int i, int j) {
	return graph->node[i].adj_edge[j];
}
*/
import "C"

// Number of bytes to use for the salt, a random string used to construct the
// table. It is prepended to the input of each HMAC call.
const SaltBytes = 8

// Number of row bytes allocated for the tag.
const TagBytes = 3

// The maximum length of the row. In general, the length of the row depends on
// the length of the longest output in the map. HASH_BYTES is defined in
// c/const.h.
const MaxRowBytes = C.HASH_BYTES

// The maximum length of the outputs. 1 byte of each row is allocated for
// padding the output string.
const MaxOutputBytes = MaxRowBytes - TagBytes - 1

// Length of the HMAC key. HMAC_KEY_BYTES is defined in c/const.h.
const DictKeyBytes = C.HMAC_KEY_BYTES

// GenerateKey generates a fresh, random key and returns it.
func GenerateDictKey() []byte {
	K := make([]byte, DictKeyBytes)
	_, err := rand.Read(K)
	if err != nil {
		return nil
	}
	return K
}

type Error string

func (err Error) Error() string {
	return string(err)
}

// Returned by Get() and priv.GetValue() if the input was not found in the
// map.
const ItemNotFound = Error("item not found")

// Returned by pub.GetShare() in case x or y is not in the table index.
const ErrorIdx = Error("index out of range")

// cError propagates an error from the internal C code.
func cError(fn string, errNo C.int) Error {
	return Error(fmt.Sprintf("%s returns error %d", fn, errNo))
}

// The public representation of the map.
type PubDict struct {
	dict *C.dict_t
}

// The private state required for evaluation queries.
type PrivDict struct {
	tinyCtx    *C.tiny_ctx
	params     C.dict_params_t
	cZeroShare *C.char
}

type Graph [][]int32

// Storage of map[string]string for processing with the C code.
type cMap struct {
	itemCt, maxOutputBytes  C.int
	inputs, outputs         **C.char
	inputBytes, outputBytes *C.int

	// Indicates to free() whether or inputs and inputBytes are a shallow copy,
	// and hence should not be freed.
	freeInputs bool
}

// newCMap constructs a new *cMap from a Go map.
//
// This must be freed with cM.free().
func newCMap(M map[string]string) (cM *cMap) {
	cM = new(cMap)
	cM.freeInputs = true
	cM.itemCt = C.int(len(M))
	cM.inputs = C.new_str_list(cM.itemCt)
	cM.inputBytes = C.new_int_list(cM.itemCt)
	cM.outputs = C.new_str_list(cM.itemCt)
	cM.outputBytes = C.new_int_list(cM.itemCt)
	cM.maxOutputBytes = C.int(0)
	i := 0
	// NOTE Go does not guarantee that the map will be traversed in the same
	// order each time.
	for in, out := range M {
		if C.int(len(out)) > cM.maxOutputBytes {
			cM.maxOutputBytes = C.int(len(out))
		}
		// NOTE C.CString() copies all the bytes of its input, even if it
		// encounters a null byte.
		C.set_str_list(cM.inputs, C.int(i), C.CString(in))
		C.set_int_list(cM.inputBytes, C.int(i), C.int(len(in)))
		C.set_str_list(cM.outputs, C.int(i), C.CString(out))
		C.set_int_list(cM.outputBytes, C.int(i), C.int(len(out)))
		i++
	}
	return cM
}

// getInput returns a pointer to the i-th input, as well as its length.
func (cM *cMap) getInput(i int) (*C.char, C.int) {
	idx := C.int(i)
	return C.get_str_list(cM.inputs, idx), C.get_int_list(cM.inputBytes, idx)
}

// getInput returns a pointer to the i-th output, as well as its length.
func (cM *cMap) getOutput(i int) (*C.char, C.int) {
	idx := C.int(i)
	return C.get_str_list(cM.outputs, idx), C.get_int_list(cM.outputBytes, idx)
}

// free() frees memory allocated to cM.
func (cM *cMap) free() {
	if cM.freeInputs {
		C.free_str_list(cM.inputs, cM.itemCt)
		C.free_int_list(cM.inputBytes)
	}
	C.free_str_list(cM.outputs, cM.itemCt)
	C.free_int_list(cM.outputBytes)
}

// getNonceMap returns a *cMap mapping each input of cM to a fresh nonce.
//
// Must free with cN.free().
func (cM *cMap) getNonceMap(nonceBytes int) (cN *cMap) {
	cN = new(cMap)
	cN.itemCt = cM.itemCt
	cN.maxOutputBytes = C.int(nonceBytes)
	cN.inputs = cM.inputs
	cN.inputBytes = cM.inputBytes
	cN.freeInputs = false // The inputs are shallow copied.
	cN.outputs = C.new_str_list(cN.itemCt)
	cN.outputBytes = C.new_int_list(cN.itemCt)

	nonce := make([]byte, nonceBytes)
	for i := C.int(0); i < cN.itemCt; i++ {
		// TODO random initial nonce, then count
		_, err := rand.Read(nonce)
		if err != nil {
			cN.free()
			return nil
		}
		C.set_str_list(cN.outputs, C.int(i), C.CString(string(nonce)))
		C.set_int_list(cN.outputBytes, C.int(i), C.int(len(nonce)))
	}
	return cN
}

// New generates a new structure (pub, priv) for the map M and key K.
//
// You must call pub.Free() and priv.Free() before these variables go out
// of scope. These structures contain C types that were allocated on the heap
// and must be freed before losing a reference to them.
func NewDict(K []byte, M map[string]string) (*PubDict, *PrivDict, error) {

	if len(M) == 0 {
		return nil, nil, Error("yup")
	}

	cM := newCMap(M)
	defer cM.free()

	pub, priv, _, err := newDictAndGraph(K, cM)
	if err != nil {
		return nil, nil, err
	}
	return pub, priv, nil
}

// TODO Unpadded version!!
func newDictAndGraph(K []byte, cM *cMap) (*PubDict, *PrivDict, Graph, error) {

	pub := new(PubDict)

	// Allocate a new dictionary object.
	tableLen := C.dict_compute_table_length(cM.itemCt)
	pub.dict = C.dict_new(
		tableLen,
		cM.maxOutputBytes,
		C.int(TagBytes),
		C.int(SaltBytes))
	if pub.dict == nil {
		return nil, nil, nil, Error(fmt.Sprintf("maxOutputBytes > %d", MaxOutputBytes))
	}

	params := cParamsToParams(&pub.dict.params)

	// Create priv.
	//
	// NOTE dict.salt is not set, and so priv.params.salt is not set. It's
	// necessary to set it after calling C.dict_create().
	priv, err := NewPrivDict(K, params)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create the dictionary.
	var errNo C.int
	cGraph := C.dict_create_and_output_graph(
		pub.dict, priv.tinyCtx, cM.inputs, cM.inputBytes, cM.outputs, cM.outputBytes, cM.itemCt, &errNo)
	if errNo != C.OK {
		priv.Free()
		return nil, nil, nil, cError("dict_create_and_output_graph", errNo)
	}
	defer C.graph_free(cGraph)

	// Copy salt to priv.params.
	C.memcpy(unsafe.Pointer(priv.params.salt),
		unsafe.Pointer(pub.dict.params.salt),
		C.size_t(priv.params.salt_bytes))

	// Save adjcency list.
	graph := make([][]int32, int32(cGraph.node_ct))
	for i := C.int(0); i < cGraph.node_ct; i++ {
		cAdjCt := C.get_node(cGraph, i).adj_ct
		graph[i] = make([]int32, int32(cAdjCt))
		for j := C.int(0); j < cAdjCt; j++ {
			graph[i][j] = int32(C.get_edge(cGraph, i, j))
		}
	}

	return pub, priv, graph, nil
}

// NewPubDictFromProto creates a new *PubDict from a *pb.Dict.
//
// You must destroy with pub.Free().
func NewPubDictFromProto(table *pb.Dict) *PubDict {
	pub := new(PubDict)
	pub.dict = (*C.dict_t)(C.malloc(C.sizeof_dict_t))

	// Allocate memory for salt + 1 tweak byte and set the parameters.
	pub.dict.params.salt = (*C.char)(C.malloc(C.size_t(len(table.GetParams().Salt) + 1)))
	setCParamsFromParams(&pub.dict.params, table.GetParams())

	// Allocate memory for table + 1 zero row and copy the table.
	tableLen := C.int(table.GetParams().GetTableLen())
	rowBytes := C.int(table.GetParams().GetRowBytes())
	realTableLen := C.int(len(table.Table)) / rowBytes
	cBuf := C.CString(string(table.Table))
	defer C.free(unsafe.Pointer(cBuf))
	pub.dict.table = (*C.char)(C.malloc(C.size_t(tableLen * rowBytes)))
	C.memset(unsafe.Pointer(pub.dict.table), 0, C.size_t(tableLen*rowBytes))
	for i := 0; i < int(realTableLen); i++ {
		src := C.get_row_ptr(cBuf, C.int(i), rowBytes)
		dst := C.get_row_ptr(pub.dict.table, C.int(table.Idx[i]), rowBytes)
		C.memcpy(unsafe.Pointer(dst), unsafe.Pointer(src), C.size_t(rowBytes))
	}

	return pub
}

// Get queries input on the structure (pub, priv). The result is M[input] =
// output, where M is the map represented by (pub, priv).
func Get(pub *PubDict, priv *PrivDict, input string) (string, error) {
	cInput := C.CString(input)
	cOutput := C.CString(string(make([]byte, pub.dict.params.max_value_bytes)))
	cOutputBytes := C.int(0)
	defer C.free(unsafe.Pointer(cInput))
	defer C.free(unsafe.Pointer(cOutput))
	errNo := C.dict_get(
		pub.dict, priv.tinyCtx, cInput, C.int(len(input)), cOutput, &cOutputBytes)
	if errNo == C.ERR_DICT_BAD_KEY {
		return "", ItemNotFound
	} else if errNo != C.OK {
		return "", cError("cdict_get", errNo)
	}
	return C.GoStringN(cOutput, cOutputBytes), nil
}

// GetShare returns the bitwise-XOR of the x-th and y-th rows of the table.
func (pub *PubDict) GetShare(x, y int) ([]byte, error) {
	if x < 0 || x >= int(pub.dict.params.table_length) ||
		y < 0 || y >= int(pub.dict.params.table_length) {
		return nil, ErrorIdx
	}
	xRow := getRow(pub.dict.table, C.int(x), pub.dict.params.row_bytes)
	yRow := getRow(pub.dict.table, C.int(y), pub.dict.params.row_bytes)
	for i := 0; i < len(xRow); i++ {
		xRow[i] ^= yRow[i]
	}
	return xRow, nil
}

// ToString returns a string representation of the table.
func (pub *PubDict) ToString() string {
	return pub.GetProto().String()
}

// GetProto returns a *pb.Dict representation of the dictionary.
func (pub *PubDict) GetProto() *pb.Dict {
	cdict := C.dict_compress(pub.dict)
	defer C.cdict_free(cdict)
	rowBytes := int(pub.dict.params.row_bytes)
	tableLen := int(cdict.compressed_table_length)
	tableIdx := make([]int32, tableLen)
	for i := 0; i < tableLen; i++ {
		tableIdx[i] = int32(C.get_int_list(cdict.idx, C.int(i)))
	}
	return &pb.Dict{
		Params: cParamsToParams(&pub.dict.params),
		Table:  C.GoBytes(unsafe.Pointer(cdict.table), C.int(tableLen*rowBytes)),
		Idx:    tableIdx,
	}
}

// Free deallocates memory associated with the underlying C implementation of
// the data structure.
func (pub *PubDict) Free() {
	C.dict_free(pub.dict)
}

// NewPrivDict creates a new *PrivDict from a key and parameters.
//
// You must destroy this with priv.Free().
func NewPrivDict(K []byte, params *pb.Params) (*PrivDict, error) {
	priv := new(PrivDict)

	// Check that K is the right length.
	if len(K) != DictKeyBytes {
		return nil, Error(fmt.Sprintf("len(K) = %d, expected %d", len(K), DictKeyBytes))
	}

	// Create new tinyprf context.
	priv.tinyCtx = C.tinyprf_new(C.int(params.GetTableLen()))
	if priv.tinyCtx == nil {
		return nil, Error("tableLen < 2")
	}

	// Allocate memory for salt.
	priv.params.salt = (*C.char)(C.malloc(C.size_t(len(params.Salt) + 1)))

	// Initialize tinyprf.
	cK := C.CString(string(K))
	defer C.memset(unsafe.Pointer(cK), 0, C.size_t(DictKeyBytes))
	defer C.free(unsafe.Pointer(cK))
	errNo := C.tinyprf_init(priv.tinyCtx, cK)
	if errNo != C.OK {
		priv.Free()
		return nil, cError("tinyprf_init", errNo)
	}

	// Set parameters.
	setCParamsFromParams(&priv.params, params)

	// A 0-byte string used by GetValue().
	priv.cZeroShare = (*C.char)(C.malloc(C.size_t(priv.params.row_bytes)))
	C.memset(unsafe.Pointer(priv.cZeroShare), 0, C.size_t(priv.params.row_bytes))

	return priv, nil
}

// GetIdx computes the two indices of the table associated with input and
// returns them.
func (priv *PrivDict) GetIdx(input string) (int, int, error) {
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))
	var x, y C.int
	errNo := C.dict_compute_rows(
		priv.params, priv.tinyCtx, cInput, C.int(len(input)), &x, &y)
	if errNo != C.OK {
		return 0, 0, cError("dict_compute_rows", errNo)
	}
	return int(x), int(y), nil
}

// GetValue computes the output associated with the input and the table rows.
// TODO Rename GetOutput.
func (priv *PrivDict) GetValue(input string, pubShare []byte) (string, error) {
	cInput := C.CString(input)
	cOutput := C.CString(string(make([]byte, priv.params.max_value_bytes)))
	defer C.free(unsafe.Pointer(cInput))
	defer C.free(unsafe.Pointer(cOutput))
	cOutputBytes := C.int(0)

	cPubShare := C.CString(string(pubShare))
	defer C.free(unsafe.Pointer(cPubShare))

	errNo := C.dict_compute_value(priv.params, priv.tinyCtx, cInput,
		C.int(len(input)), cPubShare, priv.cZeroShare, cOutput, &cOutputBytes)

	if errNo == C.ERR_DICT_BAD_KEY {
		return "", ItemNotFound
	} else if errNo != C.OK {
		return "", cError("dict_compute_value", errNo)
	}
	return C.GoStringN(cOutput, cOutputBytes), nil
}

// GetParams returns the public parameters of the data structure.
func (priv *PrivDict) GetParams() *pb.Params {
	return cParamsToParams(&priv.params)
}

// Free deallocates moemory associated with the C implementation of the
// underlying data structure.
func (priv *PrivDict) Free() {
	C.free(unsafe.Pointer(priv.params.salt))
	C.free(unsafe.Pointer(priv.cZeroShare))
	C.tinyprf_free(priv.tinyCtx)
}

// cBytesToString maps a *C.char to a []byte.
func cBytesToString(str *C.char, bytes C.int) string {
	return C.GoStringN(str, bytes)
}

// cBytesToString maps a *C.char to a []byte.
func cBytesToBytes(str *C.char, bytes C.int) []byte {
	return C.GoBytes(unsafe.Pointer(str), bytes)
}

// cParamsToParams creates *Params from a *C.dict_params_t, making a
// deep copy of the salt.
//
// Called by pub.GetParams() and priv.GetParams().
func cParamsToParams(cParams *C.dict_params_t) *pb.Params {
	return &pb.Params{
		TableLen:       *proto.Int32(int32(cParams.table_length)),
		MaxOutputBytes: *proto.Int32(int32(cParams.max_value_bytes)),
		RowBytes:       *proto.Int32(int32(cParams.row_bytes)),
		TagBytes:       *proto.Int32(int32(cParams.tag_bytes)),
		Salt:           C.GoBytes(unsafe.Pointer(cParams.salt), cParams.salt_bytes),
	}
}

// setCParamsFromDictparams copies parameters to a *C.dict_params_t.
//
// Must call C.free(cParams.salt)
func setCParamsFromParams(cParams *C.dict_params_t, params *pb.Params) {
	cParams.table_length = C.int(params.GetTableLen())
	cParams.max_value_bytes = C.int(params.GetMaxOutputBytes())
	cParams.row_bytes = C.int(params.GetRowBytes())
	cParams.tag_bytes = C.int(params.GetTagBytes())
	cParams.salt_bytes = C.int(len(params.Salt))
	cBuf := C.CString(string(params.Salt))
	C.memcpy(unsafe.Pointer(cParams.salt),
		unsafe.Pointer(cBuf),
		C.size_t(cParams.salt_bytes))
}

// getRow returns a []byte corresponding to row in the table.
func getRow(table *C.char, idx, rowBytes C.int) []byte {
	rowPtr := C.get_row_ptr(table, idx, rowBytes)
	return C.GoBytes(unsafe.Pointer(rowPtr), rowBytes)
}