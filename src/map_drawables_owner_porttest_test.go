//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"testing"
	"unsafe"
)

func mapDrawableTables(t *testing.T) {
	t.Helper()
	for _, r := range blobdata.PortTestMapDrawableTables() {
		p := unsafe.Slice((*byte)(memmap.PtrOff(r.Base, r.Offset)), len(r.Data))
		old := append([]byte(nil), p...)
		copy(p, r.Data)
		t.Cleanup(func() { copy(p, old) })
	}
}
