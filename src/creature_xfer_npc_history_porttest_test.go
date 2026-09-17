//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

func creatureXferNPCHistorical(p *mapDrawableStream, v int, spellWords ...[]uint32) {
	p.u32(1)
	p.u32(0)
	objectXferScript(p, 11)
	objectXferScript(p, 22)
	p.u16(0x1234)
	objectXferScript(p, 33)
	if v >= 32 {
		for i := 0; i < 6; i++ {
			objectXferScript(p, uint32(40+i))
		}
		if v >= 50 {
			objectXferScript(p, 50)
		}
	}
	p.u32(0)
	for i := 0; i < 18; i++ {
		p.u8(byte(10 + i))
	}
	if v == 31 {
		p.u16(0)
		p.u8(0)
	}
	if v >= 32 {
		p.u8(7)
		if v < 49 {
			p.u16(0xa5a5)
		} else {
			p.u32(0x1234a5a5)
		}
		p.u32(0x3f400000)
		p.u32(0x3e800000)
		p.u32(0x3f000000)
		p.u32(0x3f800000)
		p.u16(55)
		p.u32(0x3fa00000)
		if v < 35 {
			p.u32(4)
		}
		itemXferRewardName(p, "PortNPC")
		creatureXferSpellWords(p, v >= 34, spellWords...)
		for i := 0; i < 5; i++ {
			if v < 47 {
				p.u8(byte(20 + i))
				p.u8(byte(30 + i))
			} else {
				p.u16(uint16(0x1200 + i))
				p.u16(uint16(0x3400 + i))
			}
			if v < 34 {
				p.u32(0)
			}
		}
		if v >= 33 {
			p.u32(0x3fc00000)
		}
		if v >= 34 {
			p.u32(0x12345678)
			p.u8(9)
			p.u32(0x33445566)
			p.u32(0x3fe00000)
			if v < 42 {
				p.u16(0)
			}
			p.u8(0)
			p.u8(0)
			p.u8(0)
		}
		if v >= 35 {
			itemXferRewardName(p, ai.ActionType(0).String())
		}
	}
	if v >= 41 {
		p.u16(4)
		p.u8(0)
	}
	if v >= 42 {
		p.u8(1)
	}
	if v >= 44 {
		p.u32(0x3f200000)
	}
	if v >= 45 {
		p.u32(0x12345678)
	}
	if v >= 46 {
		p.u32(0x180)
	}
	if v >= 48 {
		p.u16(42)
	}
	if v >= 51 {
		p.u32(0x11223344)
	}
	if v >= 52 {
		p.u8(0)
	}
	if v >= 61 {
		p.u16(2)
		p.u8(0)
	}
	if v >= 62 {
		p.u8(0)
	}
}

func TestCreatureXferNPCHistorical(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "npc-history.bin")
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "creature-xfer-npc-historical", rows, "1a3494a6e4b9316aba48108cd601de0e9473ea5d3404517cafe1b37a7bf5baff")
	}()
	for version := 1; version <= 62; version++ {
		t.Run(fmt.Sprintf("v%d", version), func(t *testing.T) {
			u := newCreatureXferObject(t, s, "NPC")
			p := new(mapDrawableStream)
			p.u16(uint16(version))
			objectXferEmptyBase(p, int16(version))
			creatureXferNPCHistorical(p, version)
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
				t.Fatalf("position=%d length=%d err=%v", pos, p.Len(), err)
			}
			health := uint16(80)
			if version >= 32 {
				health = 55
			}
			if version >= 48 {
				health = 42
			}
			max := uint16(100)
			if version >= 45 {
				max = 0x5678
			}
			if u.HealthData.Cur != health || u.HealthData.Max != max {
				t.Fatalf("health=%+v wantcur=%d max=%d", u.HealthData, health, max)
			}
			status := uint32(0xa5a5)
			if version >= 49 {
				status = 0x1234a5a5
			}
			if version >= 32 && (objectXferGetWord(u.UpdateData, 1440) != status || objectXferGetWord(u.UpdateData, 1308) != 0x3fa00000) {
				t.Fatal("status width/speed mirror")
			}
			for i, b := range unsafe.Slice((*byte)(unsafe.Add(u.UpdateData, 2076)), 18) {
				if b != byte(10+i) {
					t.Fatalf("color byte%d=%d", i, b)
				}
			}
			itemXferCaptureCase(t, &rows, u, cryptfile.Global().PortTestChecksum(), 0)
		})
	}
}
