//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
	"unsafe"
)

func TestPlayerFilesAttributesUnitRead(t *testing.T) {
	o := newMatchRosterOwner(t)
	defer flags.PortTestGameFlags(0)()
	meter := legacy.PortTestNewMeterEnvironment()
	t.Cleanup(meter.Restore)
	*meter.NamedWord("nox_color_white_2523948") = 0x7fff7fff
	*o.quest["202028"] = 17
	u := &o.units[0]
	ud := u.UpdateDataPlayer()
	p := ud.Player
	p.NetCodeVal = u.NetCode
	raw := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
	info := raw[2185:2282]
	livesPtr := (*uint32)(unsafe.Add(unsafe.Pointer(ud), 320))
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 2, 3, 4, 5, 0xffff} {
		for _, lives := range []uint32{0, 5, 6, 0xffffffff} {
			for _, name := range []string{"Alex Ω", strings.Repeat("X", 24)} {
				o.reset()
				flags.ResetGame()
				*livesPtr = 77
				playerFileInfoReset(info, "Before")
				fields := bytes.Clone(info)
				for i := 50; i < 89; i++ {
					fields[i] = byte(i*9 + 3)
				}
				fields[66] = 1
				for off := 4652; off <= 4696; off += 4 {
					binary.LittleEndian.PutUint32(raw[off:], 0xaaaaaaaa)
				}
				input := playerFileAttributes(version, 2, name, fields)
				nameBytes := playerFileName(name)
				lifeAt := 2 + 1 + len(nameBytes) + 33 + 1
				if int16(version) >= 2 {
					lifeAt += 5
				}
				if int16(version) >= 5 {
					lifeAt += 4
				}
				if int16(version) >= 3 {
					binary.LittleEndian.PutUint32(input[lifeAt:], lives)
				}
				if int16(version) >= 4 {
					binary.LittleEndian.PutUint32(input[len(input)-4:], 0x89abcdef)
				}
				wantRet, wantPos := uint32(1), int64(len(input))
				wantLives := uint32(77)
				if int16(version) >= 3 {
					wantLives = lives
					if lives > 5 {
						wantRet = 0
						wantPos = int64(lifeAt + 4)
					}
				}
				expected := bytes.Clone(info)
				copy(expected, nameBytes)
				binary.LittleEndian.PutUint16(expected[len(nameBytes):], 0)
				copy(expected[50:83], fields[50:83])
				if int16(version) >= 2 {
					copy(expected[83:88], fields[83:88])
				}
				expected[88] = fields[88]
				input = append(input, 0xde, 0xad, 0xbe, 0xef)
				ret, got, pos := playerFileSection(t, "sub_41A590", input, uint32(uintptr(unsafe.Pointer(u))), uint32(uintptr(unsafe.Pointer(&info[0]))))
				if ret != wantRet || pos != wantPos || !bytes.Equal(got, input) || !bytes.Equal(info, expected) || *livesPtr != wantLives || !bytes.Equal(raw[4704:4704+len(nameBytes)], nameBytes) || binary.LittleEndian.Uint16(raw[4704+len(nameBytes):]) != 0 {
					t.Fatal("unit attributes", version, lives, len(nameBytes), ret, pos, wantPos, *livesPtr, wantLives)
				}
				var colors []uint32
				for i, off := range []int{2296, 2292, 2300, 2304, 2308} {
					rgb := fields[68+3*i : 71+3*i]
					c := uint32(rgb[0]&0xf8)<<7 | uint32(rgb[1]&0xf8)<<2 | uint32(rgb[2]&0xf8)>>3
					c |= c << 16
					gotColor := binary.LittleEndian.Uint32(raw[off:])
					if gotColor != c {
						t.Fatal("attribute color", i, gotColor, c)
					}
					colors = append(colors, gotColor)
				}
				if binary.LittleEndian.Uint32(raw[2312:]) != 0x7fff7fff {
					t.Fatal("attribute white")
				}
				for off := 4652; off <= 4696; off += 4 {
					want := uint32(0xaaaaaaaa)
					if wantRet == 1 {
						want = 0
						if off == 4688 {
							want = 17
						}
						if off == 4692 {
							want = 63
						}
						if off == 4696 {
							want = 0xaaaaaaaa
							if int16(version) >= 4 {
								want = 0x89abcdef
							}
						}
					}
					if binary.LittleEndian.Uint32(raw[off:]) != want {
						t.Fatal("attribute quest reset", version, lives, off)
					}
				}
				state := o.state()
				wantMessages := 0
				if wantRet == 1 && int16(version) >= 4 {
					wantMessages = 1
				}
				if len(state.Nodes) != wantMessages {
					t.Fatal("attribute report count", version, lives, state)
				}
				if wantMessages != 0 && !bytes.Equal(state.Nodes[0].Data, []byte{240, 29, 0xef, 0xcd}) {
					t.Fatal("attribute stage report", state)
				}
				rows = append(rows, map[string]any{"version": version, "lives": lives, "name": name, "return": ret, "position": pos, "info": bytes.Clone(info), "stored_lives": *livesPtr, "colors": colors, "quest": bytes.Clone(raw[4652:4700]), "queue": state})
			}
		}
	}
	spellbookCapture(t, "player-files-attributes-unit-read", rows, "")
}
