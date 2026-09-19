//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func TestServerBrowserDetails(t *testing.T) {
	o := newListboxOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	o.create(t, 8, 600, 3000, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Count = 256 })
	*words["dword_5d4594_815004"] = uint32(uintptr(o.win.C()))
	var entries []strman.Entry
	for _, s := range []string{"Name", "Ping", "GameType", "Quest", "Arena", "CTF", "Stage", "Map", "Individual", "Clan", "Ladder", "Occupancy", "Resolution", "DisabledSpells", "DisabledWeapons", "DisabledArmor", "None"} {
		entries = append(entries, strman.Entry{ID: strman.ID("noxworld.c:" + s), Vals: []strman.Variant{{Str: s}}})
	}
	set, restoreStrings := o.c.srv.PortTestMeterStrings(entries...)
	defer restoreStrings()
	set(0)
	for _, off := range []uintptr{89396, 89464, 89520, 89580, 89636, 89788, 89860, 89916, 90024, 90132} {
		clear(serverConfigOwnBytes(t, 0x587000, off, 2))
	}
	defs := []server.PortTestSpellClassDef{{Index: 1, Flags: 0x1000000, Valid: true}, {Index: 31, Flags: 0x2000000, Valid: true}, {Index: 32, Flags: 0x4000000, Valid: true}, {Index: 63, Flags: 0x7000000, Valid: true}, {Index: 127, Flags: 0x1000000, Valid: true}, {Index: 136, Flags: 0x2000000, Valid: true}, {Index: 137, Flags: 0x4000000, Valid: true}, {Index: 10, Flags: 0, Valid: true}, {Index: 11, Flags: 0x1000000, Valid: false}}
	saved := o.c.srv.Spells
	defer func() { o.c.srv.Spells = saved }()
	o.c.srv.Spells = server.PortTestSpellClassServer(defs).Spells
	for _, d := range defs {
		o.c.srv.Spells.DefByInd(spell.ID(d.Index)).Title = fmt.Sprintf("Spell %d", d.Index)
	}
	for _, table := range []struct {
		off   uintptr
		bits  []int
		label string
	}{{33392, []int{0, 7, 8, 26}, "Weapon"}, {35496, []int{0, 7, 25}, "Armor"}} {
		b := serverConfigOwnBytes(t, 0x587000, table.off, (len(table.bits)+1)*12)
		clear(b)
		for i, bit := range table.bits {
			p, free := alloc.CString16(fmt.Sprintf("%s %d", table.label, bit))
			defer free()
			binary.LittleEndian.PutUint32(b[i*12:], uint32(uintptr(unsafe.Pointer(p))))
			binary.LittleEndian.PutUint32(b[i*12+8:], 1<<bit)
		}
	}
	raw, free := alloc.Make([]byte{}, 172)
	defer free()
	type row struct {
		Mode   uint16
		Mask   byte
		Name   string
		Ping   int32
		Text   []string
		Colors []uint32
	}
	var rows []row
	for _, mode := range []uint16{0, 0x20, 0x1000, 0x2000, 0x3000, 0x6000, 0xa000, 0xe000} {
		for _, mask := range []byte{0, 0xff, 0xaa, 0x55} {
			for _, name := range []string{"", "Example"} {
				for _, ping := range []int32{-1, 123, 9999} {
					clear(raw)
					copy(raw[12:28], "127.0.0.1")
					copy(raw[111:120], "Arena")
					copy(raw[120:135], name)
					binary.LittleEndian.PutUint16(raw[109:], 18590)
					binary.LittleEndian.PutUint32(raw[96:], uint32(ping))
					raw[103] = 3
					raw[104] = 8
					binary.LittleEndian.PutUint16(raw[163:], mode)
					binary.LittleEndian.PutUint16(raw[165:], 65535)
					for i := 135; i < 163; i++ {
						raw[i] = mask
					}
					before := append([]byte(nil), raw...)
					legacy.PortTestServerBrowserDetails(unsafe.Pointer(&raw[0]))
					if !reflect.DeepEqual(raw, before) {
						t.Fatal("details mutated record")
					}
					text := serverPanelsListNames(o.win)
					title := name
					if title == "" {
						title = "127.0.0.1:18590"
					}
					pingText := fmt.Sprint(ping)
					if ping == 9999 {
						pingText = "--"
					}
					modeText := "Arena"
					if mode&0x20 != 0 {
						modeText = "CTF"
					}
					if mode&0x1000 != 0 {
						modeText = "Quest"
					}
					want := []string{"Name", title, "", "Ping", pingText, "", "GameType", modeText}
					if mode&0x1000 != 0 {
						want = append(want, "", "Stage", "65535")
					}
					want = append(want, "", "Map", "Arena", "")
					if mode&0xc000 != 0 {
						kind := "Clan"
						if mode&0x4000 != 0 {
							kind = "Individual"
						}
						want = append(want, kind, "Ladder")
					}
					want = append(want, "", "Occupancy", "3/8")
					if mode&0x2000 != 0 {
						want = append(want, "", "Resolution", get_video_mode_string(0), "", "DisabledSpells")
						n := len(want)
						for _, id := range []int{1, 31, 32, 63, 127, 136} {
							if raw[135+id/8]&(1<<uint(id%8)) == 0 {
								want = append(want, fmt.Sprintf("Spell %d", id))
							}
						}
						if len(want) == n {
							want = append(want, "None")
						}
						for _, table := range []struct {
							title, label string
							bits         []int
						}{{"DisabledWeapons", "Weapon", []int{0, 7, 8, 26}}, {"DisabledArmor", "Armor", []int{0, 7, 25}}} {
							want = append(want, "", table.title)
							n = len(want)
							for _, bit := range table.bits {
								if mask&(1<<uint(bit%8)) == 0 {
									want = append(want, fmt.Sprintf("%s %d", table.label, bit))
								}
							}
							if len(want) == n {
								want = append(want, "None")
							}
						}
					}
					normalized := append([]string(nil), text...)
					for i := range normalized {
						normalized[i] = strings.TrimSuffix(normalized[i], "\n")
					}
					if !reflect.DeepEqual(normalized, want) {
						t.Fatalf("details mode %x mask %x\ngot %#v\nwant %#v", mode, mask, normalized, want)
					}
					r := row{Mode: mode, Mask: mask, Name: name, Ping: ping, Text: text}
					rows = append(rows, r)
				}
			}
		}
	}
	spellbookCapture(t, "server-browser-details", rows, "7e50692b5c088e67d7f8e04e5988ea28bef7a028f72126efdafeeaebdc6ed741")
}
