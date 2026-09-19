//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestServerBrowserLANRows(t *testing.T) {
	o := newListboxOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	*words["dword_587000_87408"] = 1
	blank := serverConfigOwnBytes(t, 0x5D4594, 815120, 1)
	blank[0] = 0
	initial := serverConfigOwnBytes(t, 0x587000, 91164, 2)
	clear(initial)
	var entries []strman.Entry
	for _, name := range []string{"Quest", "CTF", "Highlander", "KotR", "Flagball", "Chat", "Arena", "Open", "Full"} {
		entries = append(entries, strman.Entry{ID: strman.ID("noxworld.c:" + name), Vals: []strman.Variant{{Str: name}}})
	}
	for _, name := range []string{"private", "closed"} {
		entries = append(entries, strman.Entry{ID: strman.ID("Noxworld.wnd:" + name), Vals: []strman.Variant{{Str: name}}})
	}
	set, restoreStrings := o.c.srv.Server.PortTestMeterStrings(entries...)
	defer restoreStrings()
	set(0)
	var lists []*gui.Window
	for _, name := range []string{"nox_wol_wnd_gameList_815012", "dword_5d4594_815016", "dword_5d4594_815020", "dword_5d4594_815024", "dword_5d4594_815028", "dword_5d4594_815032"} {
		o.create(t, 8, 180, 100, nil)
		lists = append(lists, o.win)
		*words[name] = uint32(uintptr(o.win.C()))
		o.win = nil // Keep each real listbox until GUI cleanup.
	}
	raw, free := alloc.Make([]byte{}, 172)
	defer free()
	type row struct {
		Name         string
		Mode         uint16
		Status       byte
		Players, Max byte
		Ping         int32
		Text         [6][]string
	}
	var rows []row
	modes := []struct {
		flags uint16
		name  string
	}{{0, "Arena"}, {0x1000, "Quest"}, {0x20, "CTF"}, {0x80, "Chat"}, {0x400, "Highlander"}, {0x1050, "Quest"}}
	for _, name := range []string{"", "Tiny", "Name with space", "ABCDEFGHIJKLMNO"} {
		for _, mode := range modes {
			for _, status := range []byte{0, 0x10, 0x20, 0x30} {
				for _, counts := range [][2]byte{{0, 0}, {1, 8}, {8, 8}} {
					for _, ping := range []int32{0, 123, 9999} {
						for _, w := range lists {
							w.Func94(gui.AsWindowEvent(16399, 0, 0))
						}
						clear(raw)
						copy(raw[12:28], "127.0.0.1")
						binary.LittleEndian.PutUint16(raw[109:], 18590)
						copy(raw[120:135], name)
						raw[100] = status
						raw[103] = counts[0]
						raw[104] = counts[1]
						binary.LittleEndian.PutUint16(raw[163:], mode.flags)
						binary.LittleEndian.PutUint16(raw[165:], 7)
						binary.LittleEndian.PutUint32(raw[96:], uint32(ping))
						before := append([]byte(nil), raw...)
						legacy.PortTestServerBrowserRow(unsafe.Pointer(&raw[0]))
						if !reflect.DeepEqual(raw, before) {
							t.Fatal("LAN row mutated server input")
						}
						r := row{Name: name, Mode: mode.flags, Status: status, Players: counts[0], Max: counts[1], Ping: ping}
						for i, w := range lists {
							r.Text[i] = serverPanelsListNames(w)
							if len(r.Text[i]) != 1 {
								t.Fatal("column row count", i, r.Text[i])
							}
						}
						wantMode := mode.name
						if mode.flags&0x1000 != 0 {
							wantMode += " 7"
						}
						wantPing := fmt.Sprint(ping)
						if ping == 9999 {
							wantPing = "--"
						}
						wantStatus := "Open"
						if counts[0] >= counts[1] {
							wantStatus = "Full"
						}
						switch status {
						case 0x10:
							wantStatus = "closed"
						case 0x20:
							wantStatus = "private"
						case 0x30:
							wantStatus = "private+closed"
						}
						for col, want := range map[int]string{2: fmt.Sprintf("%d/%d", counts[0], counts[1]), 3: wantMode, 4: wantPing, 5: wantStatus} {
							if r.Text[col][0] != want {
								t.Fatalf("column %d got %q want %q", col, r.Text[col][0], want)
							}
						}
						wantName := name
						if wantName == "" {
							wantName = "127.0.0.1:18590"
						}
						if !strings.HasPrefix(wantName, r.Text[1][0]) {
							t.Fatal("name truncation changed prefix", wantName, r.Text[1][0])
						}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	spellbookCapture(t, "server-browser-lan-rows", rows, "80f3bb8361995c6540ce3b09d291438433592014853abe88a1082f4800113630")
}
func TestServerBrowserNameWidth(t *testing.T) {
	o := newEntryOwner(t)
	type row struct {
		Input    string
		Width    byte
		Output   string
		Identity bool
		Words    []uint16
	}
	var rows []row
	words, free := alloc.Make([]uint16{}, 160)
	defer free()
	for _, text := range []string{"", "A", "WideName", "WWWWWWWWWWWWWWW", strings.Repeat("abc", 30), "Name with spaces"} {
		for _, width := range []byte{20, 50, 100, 255} {
			for i := range words {
				words[i] = 0xbeef
			}
			alloc.StrCopy16(words[8:152], text)
			ptr := unsafe.Pointer(&words[8])
			same := legacy.PortTestServerBrowserTrimName(ptr, width) == uintptr(ptr)
			got := alloc.GoString16(&words[8])
			if !same || !strings.HasPrefix(text, got) {
				t.Fatal("name width", text, width, got, same)
			}
			for _, v := range words[:8] {
				if v != 0xbeef {
					t.Fatal("name prefix guard")
				}
			}
			if words[152] != 0xbeef {
				t.Fatal("name suffix guard")
			}
			if o.c.R2().GetStringSizeWrapped(nil, got, 0).X+5 > int(width) {
				t.Fatal("name still exceeds width", got, width)
			}
			rows = append(rows, row{text, width, got, same, append([]uint16(nil), words...)})
		}
	}
	spellbookCapture(t, "server-browser-name-width", rows, "840800b220f259e61e4d6f9c4e893fa43800bbbcbd8b6e04d4d97b0731a9821e")
}
