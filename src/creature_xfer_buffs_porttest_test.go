//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestCreatureXferBuffWriter(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "buffs.bin")
	type row struct {
		Case string
		Wire []byte
		CRC  uint32
	}
	var rows []row
	defer func() { spellbookCapture(t, "creature-xfer-buff-writer", rows, "") }()
	order := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 24, 25, 27, 29, 26, 23, 0}
	masks := []uint32{0, 0xffffffff, 0x55555555, 0xaaaaaaaa}
	for i := 0; i < 32; i++ {
		masks = append(masks, 1<<i)
	}
	for _, mask := range masks {
		for _, duration := range []bool{false, true} {
			t.Run(fmt.Sprintf("mask%08x-duration%v", mask, duration), func(t *testing.T) {
				u := newCreatureXferObject(t, s, "Monster")
				u.Buffs = mask
				for i := range u.BuffsDur {
					u.BuffsDur[i] = uint16(0xffe0 + i)
					u.BuffsPower[i] = uint8(128 + i)
				}
				d, free := alloc.New(server.DurSpell{})
				defer free()
				if duration {
					d.Spell = 51
					d.Target48 = u
					*(*uint32)(unsafe.Add(unsafe.Pointer(d), 72)) = 0xabcdef01
					old := s.Spells.Dur.List
					s.Spells.Dur.List = d
					defer func() { s.Spells.Dur.List = old }()
				}
				before := *u
				p := new(mapDrawableStream)
				p.u16(2)
				count := byte(0)
				for _, id := range order {
					if mask&(1<<id) != 0 {
						count++
					}
				}
				p.u8(count)
				for _, id := range order {
					if mask&(1<<id) == 0 {
						continue
					}
					itemXferRewardName(p, server.EnchantID(id).String())
					p.u8(u.BuffsPower[id])
					p.u32(uint32(u.BuffsDur[id]))
					if id == 26 {
						if duration {
							p.u32(0xabcdef01)
						} else {
							p.u32(100)
						}
					}
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
					t.Fatal(err)
				}
				if ret := legacy.PortTestCreatureXferHelper(2, u, nil, 0); ret != 1 {
					t.Fatalf("writer=%d", ret)
				}
				crc := cryptfile.Global().PortTestChecksum()
				cryptfile.Close()
				got, err := os.ReadFile(path)
				if err != nil || !bytes.Equal(got, p.Bytes()) || *u != before {
					t.Fatal("buff wire/state")
				}
				rows = append(rows, row{t.Name(), got, crc})
			})
		}
	}
}

func TestCreatureXferBuffReadGates(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "buff-gates.bin")
	for _, version := range []uint16{0, 1, 2, 3, 0x7fff, 0x8000, 0xffff} {
		for _, name := range []string{"", "unknown-buff"} {
			t.Run(fmt.Sprintf("v%d-name%d", version, len(name)), func(t *testing.T) {
				u := newCreatureXferObject(t, s, "Monster")
				u.Buffs = 0x12345678
				u.BuffsDur[1] = 999
				before := *u
				p := new(mapDrawableStream)
				p.u16(version)
				count := byte(0)
				if name != "" {
					count = 1
				}
				p.u8(count)
				if count != 0 {
					itemXferRewardName(p, name)
					p.u8(3)
					p.u32(123)
				}
				if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				ret := legacy.PortTestCreatureXferHelper(2, u, nil, 0)
				wantRet := uint32(0)
				wantPos := int64(2)
				if version == 1 || version == 2 {
					wantPos = 3
					if count == 0 {
						wantRet = 1
					} else {
						wantPos += int64(1 + len(name))
					}
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != wantPos || ret != wantRet || *u != before {
					t.Fatalf("ret=%d pos=%d want=%d/%d", ret, pos, wantRet, wantPos)
				}
			})
		}
	}
}
