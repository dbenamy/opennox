//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unicode/utf16"
	"unsafe"
)

func TestCharacterCreationAppearance(t *testing.T) {
	o := newEntryOwner(t)
	colors := characterPalette(t)
	o.create(t, 0, 8, func(_ *gui.WindowData, d *gui.EntryFieldData) { d.Field_1040 = 25 })
	lang, restore := o.c.srv.PortTestMeterStrings(strman.Entry{ID: "SelColor.c:DefaultName", Vals: []strman.Variant{{Str: "DefaultHero"}}})
	t.Cleanup(restore)
	lang(0)
	var controls []*gui.Window
	var ptrs []unsafe.Pointer
	for i := 0; i < 14; i++ {
		w := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 1, 1, nil)
		controls = append(controls, w)
		ptrs = append(ptrs, w.C())
	}
	ptrs = append(ptrs, o.win.C())
	player, free := alloc.Calloc(1, 128)
	defer free()
	raw := unsafe.Slice((*byte)(player), 128)
	restoreOwner := legacy.PortTestCharacterAppearanceOwner(player, ptrs)
	defer restoreOwner()
	type row struct {
		Mask        int
		Class       byte
		Input, Name string
		Data        []byte
	}
	var rows []row
	for _, class := range []byte{0, 1, 2} {
		for mask := 0; mask < 16; mask++ {
			for _, tc := range []struct{ input, want string }{{"Hero", "Hero"}, {"  Hero  ", "Hero"}, {"", "DefaultHero"}, {" \t ", "DefaultHero"}, {" *A* ", "-A*"}, {"é英雄", "é英雄"}} {
				for i := range raw {
					raw[i] = 0xa5
				}
				raw[66] = class
				o.text(utf16.Encode([]rune(tc.input)))
				for i := 0; i < 5; i++ {
					*characterWord(controls[i], 32) = uint32(i%3) | uint32(3+5*i)<<16
				}
				for i := 5; i < 10; i++ {
					*characterWord(controls[i], 32) = uint32(0x100+17*i) << 16
				}
				for i := 0; i < 4; i++ {
					controls[10+i].Flags = 0
					if mask&(1<<i) != 0 {
						controls[10+i].Flags = 8
					}
				}
				if legacy.PortTestCharacterAppearance() != player {
					t.Fatal("return identity")
				}
				var units []uint16
				for i := 0; i < 26; i++ {
					v := binary.LittleEndian.Uint16(raw[2*i:])
					if v == 0 {
						break
					}
					units = append(units, v)
				}
				got := string(utf16.Decode(units))
				if got != tc.want {
					t.Fatalf("name %q => %q want %q", tc.input, got, tc.want)
				}
				skin := colors[3*3:][:3]
				if !bytes.Equal(raw[71:74], skin) {
					t.Fatal("skin")
				}
				for i, off := range []int{68, 74, 77, 80} {
					idx := 3 + 5*(i+1) + 32*((i+1)%3)
					want := skin
					if mask&(1<<i) != 0 {
						want = colors[3*idx:][:3]
					}
					if !bytes.Equal(raw[off:off+3], want) {
						t.Fatal("override", mask, i)
					}
				}
				for i := 0; i < 5; i++ {
					if raw[83+i] != byte(0x100+17*(5+i)) {
						t.Fatal("byte index", i)
					}
				}
				for i := 52; i < 68; i++ {
					want := byte(0xa5)
					if i == 66 {
						want = class
					}
					if raw[i] != want {
						t.Fatal("unrelated info field", i)
					}
				}
				if !bytes.Equal(raw[88:], bytes.Repeat([]byte{0xa5}, 40)) {
					t.Fatal("trailing fields")
				}
				rows = append(rows, row{mask, class, tc.input, got, append([]byte(nil), raw...)})
			}
		}
	}
	spellbookCapture(t, "character-creation-appearance", rows, "42a587e5a5f800883beeb06fb24a64cb8a8935a3d104455310de017c2d5fb010")
}
