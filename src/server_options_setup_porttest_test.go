//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestServerOptionsSetup(t *testing.T) {
	type widget struct {
		ID             uint
		Flags, Checked uint32
	}
	type row struct {
		Flags                uint32
		Mode                 uint16
		Gameplay             uint
		Score, Time, Caption string
		Controls             []widget
		Settings             []byte
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 128, 129} {
		for _, mode := range []uint16{0x100, 0x20, 0x1000, 0x8020} {
			for _, gameplay := range []uint{0, 1, 2, 4, 7} {
				t.Run(fmt.Sprintf("%x-%x-%d", flags, mode, gameplay), func(t *testing.T) {
					catalog := legacy.PortTestMapCatalogOpen(23)
					t.Cleanup(catalog.Close)
					o := newServerOptionsOwner(t)
					serverOptionsRulesOwner(t, o)
					defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
					old := noxflags.GetGamePlay()
					noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
					noxflags.SetGamePlay(noxflags.GameplayFlag(gameplay))
					t.Cleanup(func() { noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0)); noxflags.SetGamePlay(old) })
					legacy.PortTestServerOptionsModeName(0)
					binary.LittleEndian.PutUint16(o.settings[52:], mode)
					binary.LittleEndian.PutUint16(o.settings[54:], 65535)
					o.settings[56] = 255
					copy(unsafe.Slice(memmap.PtrUint8(0x5D4594, 371438), 58), o.settings)
					o.call("setup", 0, "")
					r := row{Flags: flags, Mode: mode, Gameplay: gameplay, Score: o.entry(10134), Time: o.entry(10135), Caption: o.options.ChildByID(10119).DrawData().Text(), Settings: append([]byte(nil), o.settings...)}
					wantScore, wantTime := "65535", "255"
					if flags&1 != 0 {
						index := map[uint16]int{0x100: 0, 0x20: 2, 0x1000: 5}[mode&0x17f0]
						wantScore = fmt.Sprint(101 + 11*index)
						wantTime = fmt.Sprint(13 + 7*index)
					}
					if r.Score != wantScore || r.Time != wantTime || r.Caption != legacy.PortTestServerOptionsModeName(mode) {
						t.Fatalf("setup limits/title %+v want %s/%s", r, wantScore, wantTime)
					}
					for _, id := range []uint{10119, 10122, 10134, 10135, 10141, 10153, 10161, 10199, 10330, 10331, 10332, 10333} {
						w := o.options.ChildByID(id)
						r.Controls = append(r.Controls, widget{id, uint32(w.Flags), w.DrawData().Field0 & 4})
					}
					for _, x := range []struct {
						id   uint
						flag uint
					}{{10331, 2}, {10333, 1}} {
						if (o.options.ChildByID(x.id).DrawData().Field0&4 != 0) != (gameplay&x.flag != 0) {
							t.Fatal("gameplay checkbox")
						}
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "server-options-setup", rows, "19dd81e9ff1afb90aa04eab3a0a00911ded93aae1ab1681aeeaa4cdeb6773c22")
}
