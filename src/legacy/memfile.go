package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/internal/binfile"
)

var _ = [1]struct{}{}[16-unsafe.Sizeof(binfile.MemFile{})]

func asMemfileP(p unsafe.Pointer) *binfile.MemFile {
	return (*binfile.MemFile)(p)
}
