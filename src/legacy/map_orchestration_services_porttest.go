//go:build porttest

package legacy

/*
#cgo LDFLAGS: -Wl,--wrap=nox_xxx_mapSaveMap_51E010
int __real_nox_xxx_mapSaveMap_51E010(char* path, int flags);
*/
import "C"

// Only map serialization is substituted by outer orchestration contracts.
// Generation, retries and filesystem operations continue through real services.
var orchestrationTestSave func(string, int) int

//export __wrap_nox_xxx_mapSaveMap_51E010
func __wrap_nox_xxx_mapSaveMap_51E010(path *C.char, flags C.int) C.int {
	if orchestrationTestSave != nil {
		return C.int(orchestrationTestSave(C.GoString(path), int(flags)))
	}
	return C.__real_nox_xxx_mapSaveMap_51E010(path, flags)
}
