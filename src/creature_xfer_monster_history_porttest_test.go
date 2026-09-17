//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

// Independent historical record builder: no C writer or converted serializer is
// used to manufacture older layouts.
func creatureXferMonsterHistorical(p *mapDrawableStream, v int, spellWords ...[]uint32) int {
	p.u32(1)
	p.u32(0)
	script := func(flags uint32) {
		if v >= 3 {
			objectXferScript(p, flags)
		} else {
			p.u32(0)
			p.u32(flags)
			p.u32(0)
		}
	}
	script(11)
	script(22)
	p.u16(0x1234)
	script(33)
	if v >= 31 {
		for i := 0; i < 6; i++ {
			objectXferScript(p, uint32(40+i))
		}
		if v >= 52 {
			objectXferScript(p, 50)
		}
	}
	if v >= 11 {
		p.u32(0)
	}
	if v >= 31 {
		p.u8(7)
		if v < 51 {
			p.u16(0xa5a5)
		} else {
			p.u32(0x1234a5a5)
		}
		p.u32(0x3f400000)
		p.u32(0x3e800000)
		p.u32(0x3f000000)
		p.u32(0x3f800000)
		if v < 33 {
			p.u16(0)
		}
		p.u32(0x3fa00000)
		if v < 34 {
			p.u32(4)
		}
		itemXferRewardName(p, "PortMonster")
		creatureXferSpellWords(p, v >= 34, spellWords...)
		for i := 0; i < 5; i++ {
			if v < 46 {
				p.u8(byte(20 + i))
				p.u8(byte(30 + i))
			} else {
				p.u16(uint16(0x1200 + i))
				p.u16(uint16(0x3400 + i))
			}
			if v <= 32 {
				p.u32(0)
			}
		}
		if v >= 32 {
			p.u32(0x3fc00000)
		}
		if v >= 33 {
			p.u32(0x12345678)
			p.u32(0x3fe00000)
			if v < 42 {
				p.u16(0)
			}
			if v < 53 {
				p.u32(0)
				p.u32(0)
				p.u32(0)
			} else {
				p.u8(0)
				p.u8(0)
				p.u8(0)
			}
		}
		if v >= 34 {
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
	merchantPos := p.Len()
	if v >= 44 {
		p.u32(0x3f200000)
	}
	if v >= 45 {
		p.u32(0x180)
	}
	if v >= 49 {
		p.u16(42)
	}
	if v >= 51 {
		p.Write([]byte{0, 0, 0, 0})
	}
	if v >= 62 {
		p.u16(2)
		p.u8(0)
	}
	if v >= 64 {
		p.u8(0)
	}
	return merchantPos
}

func TestCreatureXferMonsterHistorical(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "monster-history.bin")
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "creature-xfer-monster-historical", rows, "bc3e7e0d2dcdcd12343c119c15c109f47fc75df5c0ee00cc8abf553d4eaf92a5")
	}()
	for version := 1; version <= 64; version++ {
		t.Run(fmt.Sprintf("v%d", version), func(t *testing.T) {
			u := newCreatureXferObject(t, s, "Monster")
			p := new(mapDrawableStream)
			p.u16(uint16(version))
			objectXferEmptyBase(p, int16(version))
			creatureXferMonsterHistorical(p, version)
			if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
				t.Fatal(err)
			}
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			defer func() {
				cf := cryptfile.Global()
				if cf != nil && cf.File.File.Handle != nil {
					legacy.Nox_fs_close((*legacy.FILE)(cf.File.File.Handle))
				}
				cryptfile.Close()
			}()
			if err := u.CallXfer(nil); err != nil {
				t.Fatal(err)
			}
			pos, err := cryptfile.Global().File.Seek(0, 1)
			if err != nil || pos != int64(p.Len()) {
				t.Fatalf("position=%d length=%d err=%v", pos, p.Len(), err)
			}
			health := uint16(80)
			if version >= 49 {
				health = 42
			}
			if u.HealthData.Cur != health {
				t.Fatalf("health=%d want=%d", u.HealthData.Cur, health)
			}
			if version >= 31 {
				status := uint32(0xa5a5)
				if version >= 51 {
					status = 0x1234a5a5
				}
				if objectXferGetWord(u.UpdateData, 1440) != status || objectXferGetWord(u.UpdateData, 1308) != 0x3fa00000 {
					t.Fatal("status width/speed mirror")
				}
			}
			itemXferCaptureCase(t, &rows, u, cryptfile.Global().PortTestChecksum(), 0)
		})
	}
}
