//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"slices"
	"strings"
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInteractionCenteredMessageBoundaries(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	o.reset(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	rows := serverConfigOwnBytes(t, 0x5D4594, 823804, 1932)
	oldFrame, oldRate := o.c.srv.Frame(), o.c.srv.TickRate()
	defer func() { o.c.srv.SetFrame(oldFrame); o.c.srv.SetTickRate(oldRate) }()
	o.c.srv.SetFrame(0xfffffff0)
	o.c.srv.SetTickRate(30)
	type row struct {
		Length, Slot int
		Text         []uint16
		Expiry       uint32
	}
	var captured []row
	for _, length := range []int{0, 1, 255, 316, 317, 318, 319, 512, 4096} {
		input, free := alloc.Make([]uint16{}, length+1)
		alphabet := []uint16{'x', 0xff, 0x100, 0xd800, 0xdc00, 0xffff}
		for i := 0; i < length; i++ {
			input[i] = alphabet[i%len(alphabet)]
		}
		beforeInput := append([]uint16(nil), input...)
		for slot := 0; slot < 3; slot++ {
			for i := range rows {
				rows[i] = 0x5a
			}
			before := append([]byte(nil), rows...)
			*words["dword_5d4594_825736"] = uint32((slot + 2) % 3)
			o.console = nil
			interactionCall("nox_xxx_printCentered_445490", uintptr(unsafe.Pointer(&input[0])))
			if *words["dword_5d4594_825736"] != uint32(slot) {
				t.Fatal("message ring advance", slot)
			}
			want := append([]byte(nil), before...)
			n := length
			if n > 317 {
				n = 317
			}
			for i := 0; i < n; i++ {
				binary.LittleEndian.PutUint16(want[644*slot+2*i:], input[i])
			}
			binary.LittleEndian.PutUint16(want[644*slot+2*n:], 0)
			binary.LittleEndian.PutUint32(want[644*slot+636:], 134)
			want[644*slot+640] = 0
			if !bytes.Equal(rows, want) {
				t.Fatal("centered message row/timestamp/neighbor boundary", length, slot)
			}
			expected := "System: " + string(utf16.Decode(input[:length]))
			if len(o.console) != 1 || !strings.HasSuffix(o.console[0], expected) {
				t.Fatal("full console message", length, len(o.console))
			}
			if !slices.Equal(input, beforeInput) {
				t.Fatal("centered message changed input")
			}
			captured = append(captured, row{length, slot, append([]uint16(nil), input[:n]...), binary.LittleEndian.Uint32(rows[644*slot+636:])})
		}
		free()
	}
	old := append([]byte(nil), rows...)
	count := *words["dword_5d4594_825736"]
	o.console = nil
	interactionCall("nox_xxx_printCentered_445490", 0)
	if !bytes.Equal(rows, old) || *words["dword_5d4594_825736"] != count || len(o.console) != 0 {
		t.Fatal("nil centered message")
	}
	interactionCapture(t, "centered-boundaries", captured)
}

func TestClientInteractionChatMessageBoundaries(t *testing.T) {
	o := newMeterOwner(t)
	scoreboard, restore := legacy.PortTestScoreboardWords()
	defer restore()
	dr, freeDrawable := alloc.New(client.Drawable{})
	defer freeDrawable()
	oldList := o.c.Objs.List1
	defer func() { o.c.Objs.List1 = oldList }()
	dr.NetCode32 = 0x1234
	dr.PosVec = image.Pt(0x12345, -2)
	*scoreboard["nox_player_netCode_85319C"] = dr.NetCode32
	type row struct {
		Length, Spaces, Alphabet int
		Team                     uint32
		Player                   bool
		Return                   uint64
		Message                  []byte
	}
	var rows []row
	alphabets := [][]uint16{{'x', 0xff}, {0x100, 'y'}, {0xd800, 0xdc00, 0xffff}, {' ', 9, 13, 'z'}}
	for alpha, alphabet := range alphabets {
		for _, length := range []int{0, 1, 249, 250, 252, 253, 254, 255, 256, 257, 511, 512, 4096} {
			for _, spaces := range []int{0, 3} {
				input, free := alloc.Make([]uint16{}, spaces+length+1)
				for i := 0; i < spaces; i++ {
					input[i] = ' '
				}
				for i := 0; i < length; i++ {
					input[spaces+i] = alphabet[i%len(alphabet)]
				}
				before := append([]uint16(nil), input...)
				trim := 0
				for trim < len(input)-1 && input[trim] == ' ' {
					trim++
				}
				text := input[trim : len(input)-1]
				wide := false
				for _, u := range text {
					wide = wide || u > 255
				}
				for _, team := range []uint32{0, 1, 0xffffffff} {
					for _, player := range []bool{false, true} {
						o.c.srv.NetList.ResetAll()
						o.c.Objs.List1 = nil
						if player {
							o.c.Objs.List1 = dr
						}
						gotReturn := interactionCall("nox_xxx_cmdSayDo_46A4B0", uintptr(unsafe.Pointer(&input[0])), uintptr(team))
						var got []byte
						packets := 0
						o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { got = append(got, b...); packets++; return false })
						var want []byte
						wantReturn := uint64(len(input) - 1)
						if len(text) != 0 {
							width := 1
							flags := byte(2)
							if wide {
								width = 2
								flags = 4
							}
							if team != 0 {
								flags |= 1
							}
							count := byte(len(text) + 1)
							want = make([]byte, 11+width*int(count))
							want[0] = 168
							want[3] = flags
							want[8] = count
							binary.LittleEndian.PutUint16(want[1:], 0x1234)
							x, y := uint16(0xffff), uint16(0xffff)
							if player {
								x, y = 0x2345, 0xfffe
							}
							binary.LittleEndian.PutUint16(want[4:], x)
							binary.LittleEndian.PutUint16(want[6:], y)
							for i := 0; i < int(count) && i < len(text); i++ {
								if wide {
									binary.LittleEndian.PutUint16(want[11+2*i:], text[i])
								} else {
									want[11+i] = byte(text[i])
								}
							}
							wantReturn = 1
						}
						if !bytes.Equal(got, want) || gotReturn != wantReturn || packets > 1 {
							t.Fatalf("chat alpha%d length%d spaces%d team%x player%v: return%d want%d bytes%x want%x", alpha, length, spaces, team, player, gotReturn, wantReturn, got, want)
						}
						if !slices.Equal(input, before) {
							t.Fatal("chat input changed")
						}
						rows = append(rows, row{length, spaces, alpha, team, player, gotReturn, got})
					}
				}
				free()
			}
		}
	}
	interactionCapture(t, "chat-boundaries", rows)
}
