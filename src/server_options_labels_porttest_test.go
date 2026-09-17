//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestServerOptionsLabels(t *testing.T) {
	for _, flags := range []uint32{0x100, 0x101, 0x20, 0x21, 0x80, 0x81} {
		t.Run(fmt.Sprintf("%x", flags), func(t *testing.T) {
			o := newServerOptionsOwner(t)
			defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
			copy(o.settings, "Arena")
			o.call("labels", 0, "")
			wantMode := "Choose mode"
			if flags&128 == 0 {
				wantMode = "Mode: Arena battle"
				if flags&0x20 != 0 {
					wantMode = "Mode: Capture flag"
				}
			}
			action := "Options"
			if flags&1 != 0 {
				action = "Start game"
			}
			for id, want := range map[uint]string{10121: "Settings for Arena", 10118: wantMode, 10117: action} {
				got := alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(o.event(id, 16386, 0, 0)))))
				if got != want {
					t.Errorf("label %d: %q want %q", id, got, want)
				}
			}
		})
	}
}
