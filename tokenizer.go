package tokenizers

// TODO packaging: how do we build the rust lib for distribution?

/*
#cgo LDFLAGS: ./lib/libtokenizers.a -ldl -lstdc++
#include <stdlib.h>
#include "./lib/tokenizers.h"
*/
import "C"

// NOTE: There should be NO space between the comments and the `import "C"` line.
import (
	"io"
	"unsafe"
)

type Tokenizer struct {
	tokenizer unsafe.Pointer
}

var _ io.Closer = (*Tokenizer)(nil)

func FromFile(path string) (*Tokenizer, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	tokenizer, err := C.from_file(cPath)
	if err != nil {
		return nil, err
	}
	return &Tokenizer{tokenizer: tokenizer}, nil
}

func (t *Tokenizer) Close() error {
	C.free_tokenizer(t.tokenizer)
	t.tokenizer = nil
	return nil
}

func (t *Tokenizer) Encode(str string, addSpecialTokens bool) []uint32 {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	var len C.uint
	res := C.encode(t.tokenizer, cStr, &len, C.bool(addSpecialTokens))
	if len > 0 {
		// can't dealloc nil
		defer C.free(unsafe.Pointer(res))
	}
	slice := unsafe.Slice(res, len)

	tokenIDs := make([]uint32, len)
	for i, v := range slice {
		tokenIDs[i] = uint32(v)
	}
	return tokenIDs
}

func (t *Tokenizer) Decode(tokenIDs []uint32, skipSpecialTokens bool) string {
	len := C.uint(len(tokenIDs))
	res := C.decode(t.tokenizer, (*C.uint)(unsafe.Pointer(&tokenIDs[0])), len, C.bool(skipSpecialTokens))
	defer C.free(unsafe.Pointer(res))
	return C.GoString(res)
}

func (t *Tokenizer) VocabSize() uint32 {
	return uint32(C.vocab_size(t.tokenizer))
}
