//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestGameStatisticsQuestArrays(t *testing.T) {
	o := newMatchRosterOwner(t)
	statisticsPeers(t, o)
	words, restore := legacy.PortTestStatisticsGlobals()
	t.Cleanup(restore)
	raw, free := alloc.Make([]byte{}, 640)
	defer free()
	record := unsafe.Pointer(&raw[0])
	defer legacy.PortTestStatisticsCall("quest-clear", record, nil, 0)
	type row struct {
		Mask, Count int
		Names       []string
		Words       [][9]uint32
		Classes     []byte
		Cleared     bool
	}
	var rows []row
	// The reserved local slot31 is not a remote quest participant. Its separate
	// match registration is covered with the actual host address in another test.
	for _, mask := range []int{0, 1, 2, 3, 2, 0, 3, 1, 0} {
		for pl := o.s.Players.First(); pl != nil; pl = o.s.Players.Next(pl) {
			objectXferSetWord(pl.C(), 4792, 0)
		}
		count := 0
		for i := range o.units {
			pl := o.units[i].UpdateDataPlayer().Player
			name := unsafe.Slice((*byte)(unsafe.Add(pl.C(), 2096)), 12)
			clear(name)
			copy(name, fmt.Sprintf("Quest%d", i))
			for off := 4652; off <= 4692; off += 4 {
				objectXferSetWord(pl.C(), off, 0)
			}
			objectXferSetWord(pl.C(), 4688, 1)
			*(*byte)(unsafe.Add(pl.C(), 2251)) = byte(i + 1)
			if i < 2 {
				objectXferSetWord(pl.C(), 4668, uint32(10*(i+1)))
				objectXferSetWord(pl.C(), 4660, uint32(17+i))
				objectXferSetWord(pl.C(), 4676, uint32(23+i))
				objectXferSetWord(pl.C(), 4684, uint32(31+i))
				if mask&(1<<i) != 0 {
					objectXferSetWord(pl.C(), 4792, 1)
					count++
				}
			}
		}
		legacy.PortTestStatisticsCall("quest-arrays", record, nil, 0)
		if int(binary.LittleEndian.Uint16(raw)) != count || *words["array-count"] != uint32(count) {
			t.Fatal("quest array count")
		}
		r := row{Mask: mask, Count: count}
		index := 0
		ptr := func(off int) unsafe.Pointer { return *(*unsafe.Pointer)(unsafe.Add(record, off)) }
		for player := 0; player < 2; player++ {
			if mask&(1<<player) == 0 {
				continue
			}
			name := (*byte)(*(*unsafe.Pointer)(unsafe.Add(ptr(536), 4*index)))
			text := alloc.GoString(name)
			if text != fmt.Sprintf("Quest%d", player) {
				t.Fatal("quest name/order", text)
			}
			var values [9]uint32
			for j := range values {
				values[j] = *(*uint32)(unsafe.Add(ptr(540+4*j), 4*index))
			}
			score := uint32(100 * (player + 1))
			if count == 2 {
				score = score * 3 / 2
			}
			// Both cost/high-stage fields intentionally read the same player word.
			expected := [9]uint32{0xc0000200 + uint32(int(o.units[player].UpdateDataPlayer().Player.PlayerInd)+4), 1, 1, 0, uint32(17 + player), uint32(10 * (player + 1)), 0, uint32(31 + player), score}
			if values != expected {
				t.Fatalf("quest array values %v != %v", values, expected)
			}
			class := *(*byte)(unsafe.Add(ptr(576), index))
			if class != byte(player+1) {
				t.Fatal("quest class")
			}
			r.Names = append(r.Names, text)
			r.Words = append(r.Words, values)
			r.Classes = append(r.Classes, class)
			index++
		}
		legacy.PortTestStatisticsCall("quest-clear", record, nil, 0)
		for off := 536; off <= 576; off += 4 {
			if ptr(off) != nil {
				t.Fatal("quest clear retained allocation", off)
			}
		}
		if binary.LittleEndian.Uint16(raw) != uint16(count) {
			t.Fatal("quest clear changed count")
		}
		legacy.PortTestStatisticsCall("quest-clear", record, nil, 0)
		r.Cleared = true
		rows = append(rows, r)
	}
	spellbookCapture(t, "game-statistics-quest-arrays", rows, "94de5329eb6efe52e9a8ded8457e4e476c1b165e6ad6c550db183c2a3b05b4cb")
}
