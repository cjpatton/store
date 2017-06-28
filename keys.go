package keys

// #cgo LDFLAGS: -lstruct -lcrypto
// #include <struct/dict.h>
import "C"

func ComputeTableLength(itemCt int) int {
	return int(C.dict_compute_table_length(C.int(itemCt)))
}
