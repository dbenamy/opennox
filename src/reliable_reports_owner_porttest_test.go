//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

type reliableReportsOwner struct {
	s       *server.Server
	units   []server.Object
	words   map[string]*uint32
	reset   func()
	objects map[unsafe.Pointer]uint32
}

func newReliableReportsOwner(t *testing.T) *reliableReportsOwner {
	t.Helper()
	s, units, _ := objectReportsPlayers(t)
	reset, _, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	words, restore := legacy.PortTestReliableReportGlobals()
	t.Cleanup(restore)
	oldReserved := netPlayerBufSize
	netPlayerBufSize = 0
	t.Cleanup(func() { netPlayerBufSize = oldReserved })
	o := &reliableReportsOwner{s: s, units: units, words: words, objects: map[unsafe.Pointer]uint32{}}
	for i := range units {
		units[i].NetCode = uint32(1001 + i)
		o.objects[units[i].CObj()] = units[i].NetCode
	}
	o.reset = func() {
		reset()
		s.NetList.ResetAll()
		clear(unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565124)), 384))
		*words["capacity"] = 0
		*words["mask"] = 0x80000082
		*words["rateMode"] = 0
		*words["rate"] = 1
		s.SetFrame(123)
		s.SetTickRate(30)
	}
	o.reset()
	return o
}
func (o *reliableReportsOwner) state() legacy.PortTestReliableReportState {
	return legacy.PortTestReliableReportSnapshot(o.objects)
}
