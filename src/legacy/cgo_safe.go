//go:build safe

package legacy

/*
#cgo CFLAGS: -g -O0
#cgo CFLAGS: -fsanitize=address
#cgo LDFLAGS: -fsanitize=address
#cgo CFLAGS: -Dmalloc=nox_malloc
#cgo CFLAGS: -Drealloc=nox_realloc
#cgo CFLAGS: -Dcalloc=nox_calloc
#cgo CFLAGS: -Dfree=nox_free
#cgo CFLAGS: -Dmemset=nox_memset
#cgo CFLAGS: -Dmemcpy=nox_memcpy
#cgo CFLAGS: -Dmemcmp=nox_memcmp
#cgo CFLAGS: -Dstrlen=nox_strlen
#cgo CFLAGS: -Dstrcpy=nox_strcpy
#cgo CFLAGS: -Dstrcat=nox_strcat
#cgo CFLAGS: -Dstrcmp=nox_strcmp

// Preserve const-qualified C arguments in the generated export declarations.
typedef const void * nox_safe_const_void_ptr;
typedef const char * nox_safe_const_char_ptr;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

const cgoSafe = true

func init() {
	memmap.SetRuntimeChecks(true)
}

//export nox_malloc
func nox_malloc(size C.uint) unsafe.Pointer {
	p, _ := alloc.Malloc(uintptr(size))
	return p
}

//export nox_realloc
func nox_realloc(ptr unsafe.Pointer, size C.uint) unsafe.Pointer {
	return alloc.Realloc(ptr, uintptr(size))
}

//export nox_calloc
func nox_calloc(num, size C.uint) unsafe.Pointer {
	p, _ := alloc.Calloc(int(num), uintptr(size))
	return p
}

//export nox_free
func nox_free(ptr unsafe.Pointer) {
	alloc.FreePtr(ptr)
}

//export nox_memset
func nox_memset(ptr unsafe.Pointer, v C.int, size C.uint) unsafe.Pointer {
	return alloc.Memset(ptr, byte(v), uintptr(size))
}

//export nox_memcpy
func nox_memcpy(dst unsafe.Pointer, src C.nox_safe_const_void_ptr, size C.uint) unsafe.Pointer {
	return alloc.Memcpy(dst, unsafe.Pointer(src), uintptr(size))
}

//export nox_memcmp
func nox_memcmp(ptr1, ptr2 C.nox_safe_const_void_ptr, size C.uint) C.int {
	return C.int(alloc.Memcmp(unsafe.Pointer(ptr1), unsafe.Pointer(ptr2), uintptr(size)))
}

//export nox_strlen
func nox_strlen(ptr C.nox_safe_const_char_ptr) C.uint {
	return C.uint(alloc.Strlen(unsafe.Pointer(ptr)))
}

//export nox_strcpy
func nox_strcpy(dst *C.char, src C.nox_safe_const_char_ptr) *C.char {
	return (*C.char)(alloc.Strcpy(unsafe.Pointer(dst), unsafe.Pointer(src)))
}

//export nox_strcat
func nox_strcat(dst *C.char, src C.nox_safe_const_char_ptr) *C.char {
	return (*C.char)(alloc.Strcat(unsafe.Pointer(dst), unsafe.Pointer(src)))
}

//export nox_strcmp
func nox_strcmp(str1, str2 C.nox_safe_const_char_ptr) C.int {
	return C.int(alloc.Strcmp(unsafe.Pointer(str1), unsafe.Pointer(str2)))
}
