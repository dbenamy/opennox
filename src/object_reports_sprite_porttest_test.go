//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestObjectReportsSpriteBranches(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	caches := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2386972), 4)
	old := append([]uint32(nil), caches...)
	t.Cleanup(func() { copy(caches, old) })
	copy(caches, []uint32{101, 102, 103, 104})
	u, fu := alloc.New(server.Object{})
	t.Cleanup(fu)
	link, fl := alloc.New(server.Object{})
	t.Cleanup(fl)
	data, fd := alloc.Make([]byte{}, 256)
	t.Cleanup(fd)
	linked, fld := alloc.Make([]byte{}, 32)
	t.Cleanup(fld)
	link.UpdateData = unsafe.Pointer(&linked[0])
	var rows []struct {
		Name     string
		Return   int
		Bytes    []byte
		Reliable []legacy.PortTestShopPacketResult
	}
	defer func() {
		spellbookCapture(t, "object-reports-sprites", rows, "dbd4b63122abe4f7b206f3e2586a8a6221df307daabf324b41e44949ab6fcb24")
	}()
	for _, branch := range []string{"weapon", "teleport", "spike", "pressure", "elevator", "shaft-linked", "shaft-no-link", "shaft-no-data", "transporter", "fallback", "periodic-fallback"} {
		for value := 0; value < 256; value++ {
			name := fmt.Sprintf("%s/%d", branch, value)
			t.Run(name, func(t *testing.T) {
				s.NetList.ResetAll()
				reset()
				clear(data)
				clear(linked)
				*u = server.Object{NetCode: 123, TypeInd: 100, UpdateData: unsafe.Pointer(&data[0])}
				want := []byte{57, 123, 0, byte(value)}
				reliable := false
				switch branch {
				case "weapon":
					u.ObjClass = object.ClassImmobile
					u.ObjSubClass = object.SubClass(0x18)
					data[0] = byte(value)
					want[0] = 179
					want[1], want[2] = 0, 128
					reliable = true
				case "teleport":
					u.TypeInd = 101
					data[20] = byte(value)
				case "spike":
					u.TypeInd = 103
					u.ObjFlags = object.Flags(uint32(value) << 24)
					want[3] = byte(^value) & 1
				case "pressure":
					u.TypeInd = 102
					data[0] = byte(value)
					want[0] = 180
					want[3] = byte(value) & 1
				case "elevator":
					u.ObjClass = object.Class(0x4000)
					data[16] = byte(value)
					want[3] = byte(value) >> 2
				case "shaft-linked":
					u.ObjClass = object.Class(0x8000)
					*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 4)) = unsafe.Pointer(link)
					linked[16] = byte(value)
					want[3] = byte(value) >> 2
				case "shaft-no-link":
					u.ObjClass = object.Class(0x8000)
					want[3] = 0
				case "shaft-no-data":
					u.ObjClass = object.Class(0x8000)
					u.UpdateData = nil
					want[3] = 0
				case "transporter":
					u.ObjClass = object.Class(0x80)
					data[12] = byte(value)
					want[0] = 178
				case "fallback":
					want = nil
				case "periodic-fallback":
					u.TypeInd = 104
					want = nil
				}
				rv := legacy.PortTestObjectReports(4, &units[0], u, 1, 0, 0, nil)
				got := s.NetList.CopyPacketsA(1, netlist.Kind1)
				rr := snapshot()
				actual := got
				if reliable {
					if len(rr) != 1 || rr[0].Recipient != 1 || len(got) != 0 {
						t.Fatalf("reliable destination: %+v %x", rr, got)
					}
					actual = rr[0].Data
				} else if len(rr) != 0 {
					t.Fatal("unexpected reliable message")
				}
				if !bytes.Equal(actual, want) {
					t.Fatalf("got%x want%x", actual, want)
				}
				if !reliable {
					expected := 1
					if want == nil {
						expected = 0
					}
					if rv != expected {
						t.Fatalf("return%d want%d", rv, expected)
					}
				}
				rows = append(rows, struct {
					Name     string
					Return   int
					Bytes    []byte
					Reliable []legacy.PortTestShopPacketResult
				}{name, rv, got, rr})
			})
		}
	}
}
