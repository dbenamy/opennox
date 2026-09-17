//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

var itemXferKinds = []struct {
	name    string
	version uint16
}{
	{"SpellReward", 60}, {"AbilityReward", 61}, {"FieldGuide", 60},
	{"Weapon", 64}, {"Armor", 62}, {"Ammo", 60}, {"Team", 60}, {"Gold", 60},
	{"Obelisk", 61}, {"ToxicCloud", 61}, {"MonsterGenerator", 63}, {"RewardMarker", 63},
}

func newItemXferOwner(t *testing.T) *server.Server {
	s := newObjectXferOwner(t)
	t.Cleanup(s.PortTestItemXferTypes())
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 69736), 7)
	oldTable := append([]uint32(nil), table...)
	t.Cleanup(func() { copy(table, oldTable) })
	clear(table)
	for i, name := range server.AbilityNames {
		p, free := alloc.CString(name)
		t.Cleanup(free)
		table[i] = uint32(uintptr(unsafe.Pointer(p)))
	}
	for i := 1; i < s.Types.Count(); i++ {
		objectTypeCode16ByInd[i] = uint16(1000 + i)
	}
	return s
}
func newItemXferObject(t *testing.T, s *server.Server, name string) *server.Object {
	u := newObjectXferTyped(t, s, name)
	if name == "Weapon" || name == "Armor" {
		hp, free := alloc.New(server.HealthData{})
		*hp = server.HealthData{Cur: 80, Field2: 100, Max: 100}
		u.HealthData = hp
		t.Cleanup(func() { u.HealthData = nil; free() })
	}
	return u
}
func itemXferDefaultPayload(name string) []byte {
	var p mapDrawableStream
	str := func(s string) { p.u8(byte(len(s))); p.WriteString(s) }
	switch name {
	case "SpellReward":
		str(spell.ID(1).String())
	case "AbilityReward":
		str(server.AbilityBerserk.String())
	case "FieldGuide":
		str("PortGuide")
	case "Weapon", "Armor":
		p.Write(make([]byte, 4))
		p.u16(80)
		p.u32(0)
	case "Ammo":
		p.Write(make([]byte, 6))
	case "Team":
		p.Write(make([]byte, 4))
	case "Gold", "ToxicCloud":
		p.u32(0)
	case "Obelisk":
		p.u32(0)
		p.u8(0)
	case "MonsterGenerator":
		p.u8(3)
		p.Write(make([]byte, 3))
		p.u8(0)
		p.u8(0)
		p.u32(0)
		for i := 0; i < 4; i++ {
			p.u16(1)
			p.u32(0)
			p.u32(0)
		}
		p.u8(3)
		p.Write(make([]byte, 3))
		p.u8(3)
		p.Write(make([]byte, 3))
		p.u32(0)
	case "RewardMarker":
		p.Write(make([]byte, 39))
	default:
		panic(name)
	}
	return p.Bytes()
}
func TestItemXferDefaultRecords(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "default.bin")
	var rows []itemXferCaptureRow
	defer func() { spellbookCapture(t, "item-xfer-default", rows, "") }()
	for _, sp := range itemXferKinds {
		t.Run(sp.name, func(t *testing.T) {
			u := newItemXferObject(t, s, sp.name)
			u.ObjFlags = 0
			u.Extent = 0x12345678
			u.ScriptIDVal = 88
			u.PosVec = types.Pointf{X: 64.25, Y: 128.75}
			u.Field34 = 777
			switch sp.name {
			case "SpellReward", "AbilityReward":
				*(*byte)(u.UseData.Ptr) = 1
			case "FieldGuide":
				copy(unsafe.Slice((*byte)(u.UseData.Ptr), 64), "PortGuide")
			}
			var want mapDrawableStream
			want.u16(sp.version)
			want.u16(64)
			want.u32(u.Extent)
			want.u32(88)
			want.f32(64.25)
			want.f32(128.75)
			want.u8(0)
			want.Write(itemXferDefaultPayload(sp.name))
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			if err := u.CallXfer(nil); err != nil {
				t.Fatal(err)
			}
			writtenChecksum := cryptfile.Global().PortTestChecksum()
			if err := cryptfile.Close(); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want.Bytes()) {
				t.Fatalf("bytes=%x want=%x", got, want.Bytes())
			}
			if u.Field34 != 777 {
				t.Fatal("writer lifetime changed")
			}
			v := newItemXferObject(t, s, sp.name)
			v.Field34 = 999
			checksum := itemXferReadRecord(t, v, path)
			if v.Field34 != 999 || v.Extent != u.Extent || v.ScriptIDVal != 88 || v.PosVec != u.PosVec {
				t.Fatal("reader common state/lifetime")
			}
			itemXferCaptureCase(t, &rows, v, checksum, writtenChecksum)
		})
	}
}
func TestItemXferFutureVersions(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "future.bin")
	for _, sp := range itemXferKinds {
		for _, version := range []uint16{sp.version + 1, 32767} {
			t.Run(fmt.Sprintf("%s-%d", sp.name, version), func(t *testing.T) {
				u := newItemXferObject(t, s, sp.name)
				var p mapDrawableStream
				p.u16(version)
				p.u32(0x12345678)
				if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				before := *(*[193]uint32)(unsafe.Pointer(u))
				if u.CallXfer(nil) == nil {
					t.Fatal("future version accepted")
				}
				if before != *(*[193]uint32)(unsafe.Pointer(u)) {
					t.Fatal("rejected object changed")
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != 2 {
					t.Fatalf("position=%d err=%v", pos, err)
				}
			})
		}
	}
}
