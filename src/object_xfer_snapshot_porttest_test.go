//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

type objectXferState struct {
	Object [193]uint32
	Name   string
	Data   map[string][]byte
}
type objectXferCaptureRow struct {
	Case  string
	State objectXferState
}

// Capture every byte in the legacy object ABI and all type-owned data. Only
// declared pointer slots are reduced to nil/non-nil; scalar values stay intact.
// Relationship identity/order is checked independently by inventory/reference
// contracts. No allocation addresses or Go-only server handles become goldens.
func objectXferSnapshot(u *server.Object) objectXferState {
	s := objectXferState{Object: *(*[193]uint32)(unsafe.Pointer(u)), Name: u.ID(), Data: map[string][]byte{}}
	for _, off := range []int{0, 260, 264, 268, 276, 280, 284, 292, 296, 300, 308, 312, 316, 324, 328, 332, 444, 448, 452, 476, 480, 492, 496, 500, 504, 508, 512, 516, 520, 556, 688, 692, 696, 700, 704, 708, 712, 716, 720, 724, 728, 732, 736, 744, 748, 756, 760} {
		if s.Object[off/4] != 0 {
			s.Object[off/4] = 1
		}
	}
	typ := u.Server().Types.ByInd(int(u.TypeInd))
	copyData := func(name string, p unsafe.Pointer, size uintptr) {
		if p != nil && size != 0 {
			s.Data[name] = append([]byte(nil), unsafe.Slice((*byte)(p), int(size))...)
		}
	}
	copyData("init", u.InitData, typ.InitDataSize)
	copyData("collide", u.CollideData, typ.CollideDataSize)
	copyData("use", u.UseData.Ptr, typ.UseDataSize)
	copyData("update", u.UpdateData, typ.UpdateDataSize)
	copyData("script", u.Field189, 2572)
	return s
}

func objectXferCaptureCase(t *testing.T, rows *[]objectXferCaptureRow, u *server.Object) {
	*rows = append(*rows, objectXferCaptureRow{Case: t.Name(), State: objectXferSnapshot(u)})
}
