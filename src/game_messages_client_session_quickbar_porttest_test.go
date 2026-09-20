//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionQuickbar(t *testing.T) {
	q := newQuickbarOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Kind, Value, Return int
		ID, Frame               uint32
		State                   quickbarResult
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{206, 207, 223, 235} {
			ids := []uint32{0, 1, 2, 3, 4, 5, 6, 31, 255}
			if kind == 235 {
				ids = []uint32{0, 1, 2, 3, 4, 5, 6}
			}
			if kind == 223 {
				ids = []uint32{0, 1, 135, 136, 139, 140, 0x80000000, 0xffffffff}
			}
			for _, id := range ids {
				for _, value := range []int{0, 1, 127, 128, 255} {
					for _, frame := range []uint32{0, 0xffffffff} {
						setup := func() {
							q.reset(t)
							q.c.srv.SetFrame(frame)
							binary.LittleEndian.PutUint32(connected, uint32(on))
							for n := 1; n <= 5; n++ {
								off := uintptr(1047764 + 24*n)
								*memmap.PtrUint32(0x5D4594, off+8) = 2
								*memmap.PtrUint32(0x5D4594, off+12) = 0xa5a5a5a5
								*memmap.PtrUint32(0x5D4594, off+20) = 77
							}
						}
						setup()
						if on != 0 || kind == 223 {
							switch kind {
							case 206:
								q.call("sub_461090", id, uint32(value))
							case 207:
								q.call("sub_461120", id, uint32(value))
							case 223:
								q.call("sub_460EB0", id, uint32(value))
							case 235:
								q.call("sub_4610D0", id)
							}
						}
						want := q.snapshot("message", 0)
						setup()
						data := []byte{byte(kind), byte(id), byte(value)}
						if kind == 235 {
							data = data[:2]
						}
						if kind == 223 {
							data = binary.LittleEndian.AppendUint32([]byte{223}, id)
							data = append(data, byte(value))
						}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
						got := q.snapshot("message", 0)
						if n != len(data) || !bytes.Equal(data, input) || !reflect.DeepEqual(got, want) {
							t.Fatalf("quickbar on%d kind%d id%d value%d frame%d return%d", on, kind, id, value, frame, n)
						}
						for slot := 1; slot <= 5; slot++ {
							off := uintptr(1047764 + 24*slot)
							state, flags, stamp := uint32(2), uint32(0xa5a5a5a5), uint32(77)
							if on != 0 {
								if kind == 206 && id == uint32(slot) {
									state = uint32(value)
									stamp = 0
									if value == 0 {
										stamp = frame
									}
								}
								if kind == 207 && id == uint32(slot) {
									if value != 0 {
										flags |= 1 << id
									} else {
										flags &^= 1 << id
									}
								}
								if kind == 235 && (id == uint32(slot) || id == 6) {
									state = 1
									stamp = 0
								}
							}
							if memmap.Uint32(0x5D4594, off+8) != state || memmap.Uint32(0x5D4594, off+12) != flags || memmap.Uint32(0x5D4594, off+20) != stamp {
								t.Fatal("ability state/flag/cooldown isolation", on, kind, id, value, slot)
							}
						}
						if kind == 223 {
							for i := 0; i < 140; i++ {
								want := byte(0)
								if uint32(i) == id {
									want = byte(value)
								}
								if memmap.Uint8(0x5D4594, 1049544+uintptr(i)) != want {
									t.Fatal("spell timer bounds/byte value", id, value, i)
								}
							}
						}
						rows = append(rows, row{on, kind, value, n, id, frame, got})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-quickbar", rows)
}
