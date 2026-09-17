//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestCollisionCoreGateCircle(t *testing.T) {
	w := newWorldCollisionOwner(t)
	worldGeometryTables(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	reset, hits, restore := legacy.PortTestWorldGeometryHitOwner()
	t.Cleanup(restore)
	resetQueues, queues, restore := legacy.PortTestGeometryQueues()
	t.Cleanup(restore)
	a := newObjectXferSimple(t, w.s)
	b := &w.units[0]
	savedB := *b
	t.Cleanup(func() { *b = savedB })
	raw := unsafe.Slice((*byte)(w.record(t, 80)), 80)
	ud := unsafe.Pointer(&raw[8])
	oldUD := a.UpdateData
	t.Cleanup(func() { a.UpdateData = oldUD; a.Update = nil; w.s.Objs.RemoveFromUpdatable(a) })
	ids := map[unsafe.Pointer]uint32{a.CObj(): 1001, b.CObj(): 1002, ud: 2001}
	type row struct {
		Angle, Mode, Variant   int
		Offset                 [2]float32
		A, B                   worldGeometryObjectState
		Data                   []byte
		Hits                   [][5]uint32
		Queues                 [3]uint32
		Updatable, LastMessage uint32
		Packets                []byte
	}
	var rows []row
	contacts, turns, messages := 0, 0, 0
	for _, angle := range []int{0, 1, 31, 32, 63, 64, 96, 127, 128, 160, 192, 224, 255} {
		for _, offset := range [][2]float32{{0, 0}, {5, 0}, {-5, 0}, {0, 5}, {0, -5}, {15, 15}, {-15, -15}, {50, 50}} {
			for mode := 0; mode < 2; mode++ {
				for variant := 0; variant < 8; variant++ {
					reset()
					resetQueues()
					w.s.NetList.ResetAll()
					w.s.Objs.RemoveFromUpdatable(a)
					w.s.SetFrame(123)
					worldGeometryResetObject(a, 1001, 100, 100, true)
					worldGeometryResetObject(b, 1002, 100+offset[0], 100+offset[1], false)
					b.ObjClass = object.ClassPlayer
					b.VelVec = types.Pointf{3, -4}
					a.UpdateData = ud
					a.Update = legacy.PortTestGeometryDoorUpdate()
					a.ObjFlags = 4
					a.Field115 = 0
					a.Field116 = 0
					a.TeamVal = server.ObjectTeam{}
					b.TeamVal = server.ObjectTeam{}
					a.ObjOwner = nil
					a.Field34 = 0
					clear(raw)
					for i := 0; i < 8; i++ {
						raw[i] = 0xa5
						raw[72+i] = 0x5a
					}
					binary.LittleEndian.PutUint16(raw[48:], uint16(angle))
					switch variant {
					case 1:
						raw[9] = 1
					case 2:
						a.ObjOwner = b
					case 3:
						a.ObjOwner = &w.units[1]
					case 4:
						b.ObjFlags = 0x8000000
					case 5:
						b.ObjFlags = 0x8000008
					case 6:
						a.TeamVal.ID = 1
						b.TeamVal.ID = 2
					case 7:
						a.TeamVal.ID = 1
						b.TeamVal.ID = 1
					}
					legacy.PortTestCollisionCore("gate-circle", a, b, nil, int32(mode))
					h := hits(ids)
					q := queues(ids)
					if len(h) > 0 {
						contacts++
						if binary.LittleEndian.Uint32(raw[52:]) != 123 {
							t.Fatal("gate contact timestamp")
						}
					}
					if q[0] != 0 {
						turns++
						if a.IsUpdatable != 1 {
							t.Fatal("turning gate absent from update list")
						}
					}
					packets := w.s.NetList.CopyPacketsA(b.ControllingPlayer().PlayerIndex(), netlist.Kind1)
					if len(packets) > 0 {
						messages++
					}
					// Repeat a locked contact in the same frame: notification is throttled.
					if variant == 6 && len(h) > 0 {
						legacy.PortTestCollisionCore("gate-circle", a, b, nil, int32(mode))
						if len(w.s.NetList.CopyPacketsA(b.ControllingPlayer().PlayerIndex(), netlist.Kind1)) != 0 {
							t.Fatal("gate message throttle")
						}
					}
					for i := 0; i < 8; i++ {
						if raw[i] != 0xa5 || raw[72+i] != 0x5a {
							t.Fatal("gate data guards")
						}
					}
					data := append([]byte(nil), raw[8:72]...)
					binary.LittleEndian.PutUint32(data[36:], 0) // owned next pointer; queue identity captured separately
					rows = append(rows, row{angle, mode, variant, offset, worldGeometryState(a), worldGeometryState(b), data, h, q, a.IsUpdatable, a.Field34, packets})
				}
			}
		}
	}
	if contacts == 0 || turns == 0 || messages == 0 {
		t.Fatal("gate effects not exercised", contacts, turns, messages)
	}
	spellbookCapture(t, "collision-core-gate-circle", rows, "5fec35c65f256144b5b41a59d6d4788463e8dbe64a56c55369953746d002b69d")
}
