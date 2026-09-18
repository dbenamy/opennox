//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameStatisticsMetadata(t *testing.T) {
	_ = newMatchRosterOwner(t)
	words, restore := legacy.PortTestStatisticsGlobals()
	t.Cleanup(restore)
	settings := serverConfigOwnBytes(t, 0x5D4594, 371516, 144)
	slot := serverConfigOwnBytes(t, 0x5D4594, 371380, 124)
	report := serverConfigOwnBytes(t, 0x5D4594, 599476, 640)
	serverConfigOwnBytes(t, 0x587000, 60072, 4)
	*memmap.PtrUint32(0x587000, 60072) = 0xfedcba98
	for i := range settings {
		settings[i] = byte(i*7 + 1)
	}
	for i := range slot {
		slot[i] = byte(i*13 + 3)
	}
	copy(slot[:9], []byte("LongMap!\x00"))
	copy(slot[9:24], []byte("A long game nam"))
	type row struct {
		Flags  uint32
		Mode   uint16
		Result uint32
		Report []byte
	}
	var rows []row
	for _, flags := range []uint32{0, 8192, 4096, 12288} {
		for _, mode := range []uint16{0, 0x10, 0x20, 0x40, 0x100, 0x400, 0x570, 0xc010, 0xffff} {
			func() {
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
				for i := range report {
					report[i] = 0xa5
				}
				*words["mode"] = 0x76543210
				binary.LittleEndian.PutUint16(slot[52:], mode)
				legacy.PortTestStatisticsCall("metadata", nil, nil, 0)
				enabled := flags&8192 != 0 && flags&4096 == 0
				if !enabled {
					if !bytes.Equal(report, bytes.Repeat([]byte{0xa5}, len(report))) || *words["mode"] != 0x76543210 {
						t.Fatal("disabled metadata changed")
					}
				} else {
					word := func(off int) uint32 { return binary.LittleEndian.Uint32(report[off:]) }
					if binary.LittleEndian.Uint16(report[6:]) != uint16(settings[103]) || word(8) != uint32(settings[104]) || word(12) != binary.LittleEndian.Uint32(settings[40:]) || word(16) != 0xfedcba98 {
						t.Fatal("session metadata widths")
					}
					expected := uint32(0x76543210)
					for i, mask := range []uint16{0x100, 0x20, 0x40, 0x10, 0x400} {
						if mode&mask != 0 {
							expected = uint32(i)
							break
						}
					}
					if *words["mode"] != expected {
						t.Fatal("game mode precedence")
					}
					if string(report[352:368]) != "A long game nam\x00" || string(report[96:105]) != "LongMap!\x00" {
						t.Fatal("metadata string limits")
					}
					if word(32) != uint32(binary.LittleEndian.Uint16(slot[54:])) || word(36) != uint32(slot[56]) || word(44) != uint32(settings[101]&15) || word(48) != uint32(settings[101]>>4) {
						t.Fatal("limits and nibble fields")
					}
					if word(52) != uint32(binary.LittleEndian.Uint16(settings[105:])) || word(56) != uint32(binary.LittleEndian.Uint16(settings[107:])) {
						t.Fatal("unaligned settings fields")
					}
					for _, off := range []int{608, 612, 616, 620, 636} {
						if word(off) != 0 {
							t.Fatal("report array reset", off)
						}
					}
				}
				rows = append(rows, row{flags, mode, *words["mode"], bytes.Clone(report)})
			}()
		}
	}
	spellbookCapture(t, "game-statistics-metadata", rows, "fcef574e7262b724ade96dc46653d78d26366e458550d16d85884fcc3696729a")
}
