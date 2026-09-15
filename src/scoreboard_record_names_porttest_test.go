//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"golang.org/x/image/font/basicfont"
)

func TestScoreboardRecordNameBoundaries(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	for _, advance := range []int{7, 2} {
		face := *basicfont.Face7x13
		face.Advance = advance
		font, restore := o.c.Render().GetFonts().PortTestWindowFont(&face, "default")
		o.c.dataRefs[uint32(uintptr(font))] = 0xf0900000 + uint32(advance)
		for _, mode := range []int{2, 5} {
			for _, elimination := range []bool{false, true} {
				for _, name := range []string{"", strings.Repeat("i", 25), strings.Repeat("i", 26), strings.Repeat("i", 27), strings.Repeat("i", 28), strings.Repeat("😀", 14), "界é😀Name%25"} {
					o.resetRank(t)
					o.constructRank(t)
					p := &o.players[0]
					p.Active = 1
					p.PlayerInd = 0
					p.Field3680 = 0
					p.Lessons = 12
					p.Field2140 = 13
					p.Field2108 = 14
					p.SetName(name)
					*(*byte)(unsafe.Add(unsafe.Pointer(p), 2251)) = 1
					original := append([]byte(nil), unsafe.Slice((*byte)(unsafe.Pointer(&p.NameFinal[0])), 56)...)
					*o.rankWords["dword_5d4594_1090120"] = uint32(mode)
					if elimination {
						noxflags.SetGame(noxflags.GameModeElimination)
					}
					ret := legacy.PortTestScoreboard(11, 0, 0, 0)
					record := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1084132)), 80)
					score := uint32(12)
					if elimination {
						score = 13
					}
					if mode == 5 {
						score = 14
					}
					if binary.LittleEndian.Uint32(record[52:]) != 0xffffffff || record[56] != 1 || binary.LittleEndian.Uint32(record[60:]) != 100 || binary.LittleEndian.Uint32(record[64:]) != score {
						t.Fatal("name changed non-name record fields")
					}
					if !bytes.Equal(original, unsafe.Slice((*byte)(unsafe.Pointer(&p.NameFinal[0])), 56)) {
						t.Fatal("score collection changed the player name")
					}
					rows = append(rows, o.rankCapture(t, 11, ret))
					o.rankWindow().Show()
					o.c.GUI.Draw()
					rows = append(rows, o.rankCapture(t, 3, 1))
				}
			}
		}
		o.releaseRank()
		restore()
	}
	scoreboardCapture(t, "record-name-boundaries", rows, "cf2f40ea64bd6f25921ae884e755522c969157cb1f5efe8161b7dd7e57e73c40")
}
