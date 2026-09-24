package legacy

/*
#include <stdint.h>
*/
import "C"
import (
	"unsafe"
)

//export nox_server_getFirstMapGroup_57C080
func nox_server_getFirstMapGroup_57C080() unsafe.Pointer {
	return GetServer().S().MapGroups.GetFirstMapGroup().C()
}

func Get_dword_5d4594_3835312() int {
	return int(*prefabGlobal(prefabInstance))
}
