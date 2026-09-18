//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"strings"
	"testing"
)

func TestMapMetadataTOCRead(t *testing.T) {
	o := newObjectRenderOwner(t, "PaintBook", "PaintObject", "PaintDoor")
	s := newObjectXferOwner(t)
	oldCount := objectTypeCode16ByInd_len
	t.Cleanup(func() { objectTypeCode16ByInd_len = oldCount })
	dir := t.TempDir()
	type entry struct {
		Name string
		Code uint16
	}
	groups := [][]entry{
		nil,
		{{"PaintObject", 1}, {"PaintBook", 65535}, {"PaintDoor", 32768}},
		{{"PaintBook", 7}, {"PaintBook", 0}, {"PaintBook", 19}},
		{{"unknown", 1}, {"", 65535}, {"PaintObject\x00ignored", 31}},
		{{strings.Repeat("z", 255), 42}, {"PaintDoor", 0}},
	}
	type record struct {
		Name  string
		IO    legacy.PortTestMapSectionWire
		Table []uint16
		Count int
	}
	var records []record
	for _, ver := range []uint16{0, 1, 2, 32767, 32768, 65535} {
		for _, flags := range []uint32{0, 1, 2, 3, 0x400000, 0x400002} {
			for gi, entries := range groups {
				name := fmt.Sprintf("v=%d/flags=%x/group=%d", ver, flags, gi)
				t.Run(name, func(t *testing.T) {
					noxflags.UnsetGame(^noxflags.GameFlag(0))
					noxflags.SetGame(noxflags.GameFlag(flags))
					objectTypeCode16ByInd = make([]uint16, 64)
					for i := range objectTypeCode16ByInd {
						objectTypeCode16ByInd[i] = uint16(i + 1000)
					}
					objectTypeCode16ByInd_len = 99
					want := append([]uint16(nil), objectTypeCode16ByInd...)
					count := 99
					wire := binary.LittleEndian.AppendUint16(nil, ver)
					wire = binary.LittleEndian.AppendUint16(wire, uint16(len(entries)))
					for _, e := range entries {
						wire = binary.LittleEndian.AppendUint16(wire, e.Code)
						wire = append(wire, byte(len(e.Name)))
						wire = append(wire, e.Name...)
					}
					ret := uint32(0)
					pos := int64(2)
					if int16(ver) <= 1 {
						ret = 1
						pos = int64(len(wire))
						clear(want)
						count = 0
						for _, e := range entries {
							name := strings.SplitN(e.Name, "\x00", 2)[0]
							id := s.Types.IndByID(name)
							if flags&2 != 0 && flags&1 == 0 {
								id = o.c.Things.IndByID(name)
							}
							want[uint16(id)] = e.Code
						}
					}
					result := legacy.PortTestMapMetadata(legacy.PortTestMapSectionIO{Function: "toc", Read: true, Data: wire}, dir)
					if result.Return != ret || result.Position != pos || !bytes.Equal(result.Data, wire) {
						t.Fatalf("IO got %+v want return %d position %d", result, ret, pos)
					}
					if !reflect.DeepEqual(objectTypeCode16ByInd, want) || objectTypeCode16ByInd_len != count {
						t.Fatalf("table %v count %d want %v count %d", objectTypeCode16ByInd, objectTypeCode16ByInd_len, want, count)
					}
					records = append(records, record{name, result, append([]uint16(nil), objectTypeCode16ByInd...), objectTypeCode16ByInd_len})
				})
			}
		}
	}
	if s.Types.IndByID("PaintObject") == o.c.Things.IndByID("PaintObject") {
		t.Fatal("fixture must distinguish server/client routing")
	}
	spellbookCapture(t, "map-metadata-toc-read", records, "45c642b26e92d366560160986e551b8f69b61eda2afc7ba2042424b650532e45")
}

func TestMapMetadataTOCWrite(t *testing.T) {
	s := newObjectXferOwner(t)
	oldCount := objectTypeCode16ByInd_len
	t.Cleanup(func() { objectTypeCode16ByInd_len = oldCount })
	// Real factory objects, linked through each list consumed by the dictionary.
	// Distinct valid type indices need no behavior-specific update data here.
	objects := make([]*server.Object, 7)
	for i := range objects {
		objects[i] = newObjectXferSimple(t, s)
	}
	oldList, oldPending, oldMissile := s.Objs.List, s.Objs.Pending, s.Objs.MissileList
	t.Cleanup(func() {
		s.Objs.List, s.Objs.Pending, s.Objs.MissileList = oldList, oldPending, oldMissile
		for _, u := range objects {
			u.ObjNext = nil
			u.InvFirstItem = nil
			u.InvNextItem = nil
			u.TypeInd = 1
		}
	})
	dir := t.TempDir()
	type record struct {
		Name  string
		IO    legacy.PortTestMapSectionWire
		Table []uint16
		Count int
	}
	var records []record
	for _, flags := range []uint32{0, 1, 0x200000, 0x200001, 0x400000, 0x400001} {
		for pattern := 0; pattern < 4; pattern++ {
			name := fmt.Sprintf("flags=%x/pattern=%d", flags, pattern)
			t.Run(name, func(t *testing.T) {
				noxflags.UnsetGame(^noxflags.GameFlag(0))
				noxflags.SetGame(noxflags.GameFlag(flags))
				objectTypeCode16ByInd = make([]uint16, 64)
				for i := range objectTypeCode16ByInd {
					objectTypeCode16ByInd[i] = uint16(i + 1000)
				}
				objectTypeCode16ByInd_len = 99
				s.Objs.List, s.Objs.Pending, s.Objs.MissileList = nil, nil, nil
				for _, u := range objects {
					u.ObjNext = nil
					u.InvFirstItem = nil
					u.InvNextItem = nil
				}
				types := []uint16{3, 2, 4, 1, 3, 2, 4}
				if pattern == 2 {
					types = []uint16{1, 1, 1, 1, 1, 1, 1}
				}
				if pattern == 3 {
					types = []uint16{4, 3, 2, 1, 2, 3, 4}
				}
				for i, u := range objects {
					u.TypeInd = types[i]
				}
				if pattern != 0 {
					s.Objs.List = objects[0]
					objects[0].ObjNext = objects[3]
					objects[0].InvFirstItem = objects[1]
					objects[1].InvNextItem = objects[2]
					s.Objs.Pending = objects[4]
					objects[4].ObjNext = objects[5]
					s.Objs.MissileList = objects[6]
				}
				want := make([]uint16, 64)
				count := 0
				if pattern != 0 && flags&0x200001 != 0 {
					for _, id := range types {
						if want[id] == 0 {
							count++
							want[id] = uint16(count)
						}
					}
				}
				wire := []byte{1, 0}
				wire = binary.LittleEndian.AppendUint16(wire, uint16(count))
				for id, code := range want {
					if code == 0 {
						continue
					}
					name := s.Types.ByInd(id).ID()
					wire = binary.LittleEndian.AppendUint16(wire, code)
					wire = append(wire, byte(len(name)))
					wire = append(wire, name...)
				}
				result := legacy.PortTestMapMetadata(legacy.PortTestMapSectionIO{Function: "toc"}, dir)
				if result.Return != 1 || result.Position != int64(len(wire)) || !bytes.Equal(result.Data, wire) {
					t.Fatalf("IO got %+v want bytes %x", result, wire)
				}
				if !reflect.DeepEqual(objectTypeCode16ByInd, want) || objectTypeCode16ByInd_len != count {
					t.Fatalf("table %v count %d want %v count %d", objectTypeCode16ByInd, objectTypeCode16ByInd_len, want, count)
				}
				records = append(records, record{name, result, append([]uint16(nil), objectTypeCode16ByInd...), objectTypeCode16ByInd_len})
			})
		}
	}
	spellbookCapture(t, "map-metadata-toc-write", records, "227b5d9c2ad091c8af85f49989086adb822bf664600dcfd80a873a3723959181")
}
