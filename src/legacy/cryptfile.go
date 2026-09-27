package legacy

/*
#include <stdio.h>
*/
import "C"
import (
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

func nox_xxx_mapgenGetSomeFile_426A60() *C.FILE {
	return NewFileHandle(cryptfile.Global().File.File)
}
