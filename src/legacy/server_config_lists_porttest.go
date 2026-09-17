//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "common/fs/nox_fs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func PortTestServerConfigAdmission(op string, index int, p, q unsafe.Pointer) unsafe.Pointer {
	switch op {
	case "allowed-add":
		return unsafe.Pointer(C.sub_4168A0((*C.wchar2_t)(p)))
	case "blocked-add":
		return unsafe.Pointer(C.sub_416770(C.int(index), (*C.wchar2_t)(p), (*C.char)(q)))
	case "allowed-remove":
		return unsafe.Pointer(C.sub_416860(C.int(index)))
	case "blocked-remove":
		C.sub_416820(C.int(index))
		return nil
	case "allowed-first":
		return unsafe.Pointer(C.sub_4168E0())
	case "allowed-next":
		return unsafe.Pointer(C.sub_4168F0((*C.int)(p)))
	case "blocked-first":
		return unsafe.Pointer(C.sub_416900())
	case "blocked-next":
		return unsafe.Pointer(C.sub_416910((*C.int)(p)))
	case "expire":
		C.sub_416720()
		return nil
	case "close":
		return unsafe.Pointer(C.sub_416950())
	default:
		panic(op)
	}
}
func PortTestServerConfigFile(op, path string) int {
	name, free := alloc.CString(path)
	defer free()
	switch op {
	case "read":
		return int(C.sub_4E41B0((*C.char)(unsafe.Pointer(name))))
	case "write":
		return int(uintptr(unsafe.Pointer(C.sub_4E43F0((*C.char)(unsafe.Pointer(name))))))
	case "allowed", "blocked":
		file := C.nox_fs_open_text((*C.char)(unsafe.Pointer(name)))
		if file == nil {
			return -1
		}
		defer C.nox_fs_close(file)
		if op == "allowed" {
			return int(C.sub_4E4390(file))
		}
		return int(C.sub_4E42C0(file))
	default:
		panic(op)
	}
}
func PortTestServerConfigRulePopulate(settings unsafe.Pointer) int {
	return int(uintptr(C.sub_4CED40((*C.char)(settings))))
}

func PortTestServerConfigAdmissionInit() int { return int(C.sub_416920()) }
