//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func itemXferHistoricalPayload(p *mapDrawableStream, name string, v int16) {
	str := func(s string) { p.u8(byte(len(s))); p.WriteString(s) }
	switch name {
	case "SpellReward":
		if v < 31 {
			p.u8(4)
			p.u8(1)
			p.u8(2)
			if v == 10 {
				p.u8(9)
			}
		} else if v < 41 {
			str(spell.ID(3).String())
			str(spell.ID(1).String())
			str(spell.ID(2).String())
		} else {
			str(spell.ID(2).String())
		}
	case "AbilityReward":
		str(server.AbilityWarcry.String())
	case "FieldGuide":
		str("HistoricalGuide")
	case "Weapon", "Armor":
		if v < 11 {
			return
		}
		p.Write(make([]byte, 4))
		if name == "Armor" && v >= 41 || name == "Weapon" && v >= 42 {
			p.u16(37)
		}
		if name == "Armor" && v == 61 || name == "Weapon" && v == 63 {
			p.u8(7)
		}
		if name == "Armor" && v >= 62 || name == "Weapon" && v >= 64 {
			p.u32(123)
		}
	case "Ammo":
		p.Write(make([]byte, 4))
		p.u8(7)
		p.u8(9)
	case "Team":
		p.Write(make([]byte, 4))
	case "Gold", "ToxicCloud":
		p.u32(0x12345678)
	case "Obelisk":
		if v >= 61 {
			p.u32(37)
			p.u8(1)
		}
	case "MonsterGenerator":
		p.u8(3)
		p.Write([]byte{1, 2, 3})
		p.u8(4)
		p.u8(5)
		p.u32(123)
		for i := 0; i < 4; i++ {
			objectXferScript(p, uint32(i+1))
		}
		p.u8(3)
		p.Write(make([]byte, 3))
		if v >= 62 {
			p.u8(3)
			p.Write([]byte{7, 8, 9})
		}
		if v >= 63 {
			p.u32(456)
		}
	case "RewardMarker":
		p.u32(123)
		p.u32(456)
		p.u16(0)
		p.u16(0)
		p.u16(0)
		for i := 0; i < 5; i++ {
			p.u32(uint32(100 + i))
		}
		if v >= 62 {
			p.u32(105)
		}
		if v >= 63 {
			p.u8(7)
		}
	default:
		panic(name)
	}
}
func TestItemXferHistoricalRecords(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "historical.bin")
	var rows []itemXferCaptureRow
	defer func() { spellbookCapture(t, "item-xfer-historical", rows, "") }()
	for _, sp := range itemXferKinds {
		for _, v := range []int16{-32768, -1, 0, 1, 2, 10, 11, 20, 21, 30, 31, 40, 41, 42, 60, 61, 62, 63, 64} {
			if v > int16(sp.version) {
				continue
			}
			if v <= 0 && (sp.name == "ToxicCloud" || sp.name == "MonsterGenerator" || sp.name == "RewardMarker") {
				continue
			}
			t.Run(fmt.Sprintf("%s-v%d", sp.name, v), func(t *testing.T) {
				u := newItemXferObject(t, s, sp.name)
				u.Field34 = 999
				var p mapDrawableStream
				p.u16(uint16(v))
				objectXferEmptyBase(&p, v)
				itemXferHistoricalPayload(&p, sp.name, v)
				if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				if err := u.CallXfer(nil); err != nil {
					t.Fatal(err)
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != int64(p.Len()) {
					t.Fatalf("position=%d want=%d err=%v", pos, p.Len(), err)
				}
				life := uint32(999)
				if v < 11 && (sp.name == "Weapon" || sp.name == "Armor") {
					life = 0
				}
				if u.Field34 != life {
					t.Fatalf("saved lifetime=%d want=%d", u.Field34, life)
				}
				switch sp.name {
				case "SpellReward", "AbilityReward":
					if *(*byte)(u.UseData.Ptr) != 2 {
						t.Fatal("reward identifier")
					}
				case "FieldGuide":
					if alloc.GoString((*byte)(unsafe.Pointer(u.UseData.Ptr))) != "HistoricalGuide" {
						t.Fatal("guide name")
					}
				case "Gold":
					if objectXferGetWord(u.InitData, 0) != 0x12345678 {
						t.Fatal("gold quantity")
					}
				case "ToxicCloud":
					if objectXferGetWord(u.UpdateData, 0) != 0x12345678 {
						t.Fatal("cloud state")
					}
				}
				itemXferCaptureCase(t, &rows, u)
			})
		}
	}
}
