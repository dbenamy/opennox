//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestMatchRosterSettingsMessage(t *testing.T) {
	o := newMatchRosterOwner(t)
	var rows []matchRosterMessageRow
	defer func() {
		spellbookCapture(t, "match-roster-settings", rows, "5c26e9c133d55a983f3d7a5bdd9786b1a81b6f57832a480710cf13735df53d6e")
	}()
	settings := unsafe.Slice(memmap.PtrUint8(0x5D4594, 371380), 58)
	nameBuf := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1324), 16)
	scores := unsafe.Slice(memmap.PtrUint16(0x5D4594, 3488), 6)
	times := unsafe.Slice(memmap.PtrUint8(0x5D4594, 3500), 6)
	modes := []uint16{256, 1024, 32, 16, 64, 4096, 0, 96, 0xffff}
	for modeIndex, mode := range modes {
		for _, seed := range []uint32{0, 1, 0x800001ff, 0xffffffff} {
			for timerCase := 0; timerCase < 6; timerCase++ {
				label := fmt.Sprintf("mode%x/seed%x/timer%d", mode, seed, timerCase)
				t.Run(label, func(t *testing.T) {
					o.reset()
					restore := noxflags.PortTestGameFlags(noxflags.GameFlag(seed))
					defer restore()
					for i := range settings {
						settings[i] = byte(seed + uint32(i*17))
					}
					binary.LittleEndian.PutUint16(settings[52:], mode)
					for i := range scores {
						scores[i] = uint16(seed) + uint16(i*31+7)
						times[i] = byte(i*19 + 3)
					}
					index := modeIndex
					if index >= 6 {
						index = 0
					}
					if timerCase == 0 || timerCase == 3 {
						times[index] = 0
					}
					noxServer.flag3592 = timerCase == 3
					*memmap.PtrUint32(0x587000, 4660) = uint32(timerCase & 1)
					deadline := uint32(12000)
					if timerCase == 1 {
						deadline = 1
					}
					if timerCase == 5 {
						deadline = 0xffffffff
					}
					*memmap.PtrUint64(0x5D4594, 3468) = uint64(deadline)
					*memmap.PtrUint32(0x5D4594, 3464) = seed
					*o.roster["server-subflags"] = seed ^ 0xa5a55a5a
					names := []string{"", "a", "abcdefghijklmno", "roster"}
					name := names[int(seed&3)]
					clear(nameBuf)
					copy(nameBuf, name)
					expectedSettings := bytes.Clone(settings)
					legacy.PortTestMatchRosterSendSettings()
					state := o.state()
					a := make([]byte, 20)
					a[0] = 175
					binary.LittleEndian.PutUint32(a[1:], o.s.Frame())
					binary.LittleEndian.PutUint32(a[5:], uint32(0x000f039a))
					binary.LittleEndian.PutUint32(a[9:], seed&0x7fff0)
					binary.LittleEndian.PutUint32(a[13:], seed^0xa5a55a5a)
					a[17] = byte(seed)
					a[18] = byte(scores[index])
					a[19] = times[index]
					b := make([]byte, 49)
					b[0] = 176
					copy(b[1:17], name)
					copy(b[17:45], settings[24:52])
					if timerCase&1 != 0 && (noxServer.flag3592 || times[index] != 0) {
						binary.LittleEndian.PutUint32(b[45:], deadline-uint32(o.ticks))
					}
					if len(state.Nodes) != 2 || !bytes.Equal(state.Nodes[1].Data, a) || !bytes.Equal(state.Nodes[0].Data, b) {
						t.Fatalf("settings payload count=%d; expected messages %x / %x", len(state.Nodes), a, b)
					}
					if !bytes.Equal(settings, expectedSettings) {
						t.Fatal("input settings changed")
					}
					rows = append(rows, matchRosterMessageRow{label, 0, state})
				})
			}
		}
	}
}

func TestMatchRosterPlayerIDs(t *testing.T) {
	o := newMatchRosterOwner(t)
	var rows []matchRosterMessageRow
	defer func() {
		spellbookCapture(t, "match-roster-player-ids", rows, "bd3f32ead3e649e55de5e1ad4789caefca4705a0f84a5bec3da1f2dec4428880")
	}()
	for _, requestSlot := range []int{1, 3, 7, 31} {
		pl := o.s.Players.ByInd(ntype.PlayerInd(requestSlot))
		for mask := 0; mask < 8; mask++ {
			for _, self := range []bool{false, true} {
				label := fmt.Sprintf("slot%d/mask%d/self%t", requestSlot, mask, self)
				t.Run(label, func(t *testing.T) {
					o.reset()
					pl.NetCodeVal = 0xffffffff
					if self {
						pl.NetCodeVal = o.units[1].NetCode
					}
					var want [][]byte
					for i := range o.units {
						u := &o.units[i]
						objectXferSetWord(u.UpdateData, 260, uint32(mask&(1<<i)))
						u.TypeInd = uint16(0x1234 + i)
						if mask&(1<<i) == 0 || u.NetCode == pl.NetCodeVal {
							continue
						}
						data := []byte{210, 0, 0, 0, 0, 1, 2}
						binary.LittleEndian.PutUint16(data[1:], uint16(o.s.GetUnitNetCode(u)))
						binary.LittleEndian.PutUint16(data[3:], u.TypeInd)
						want = append([][]byte{data}, want...)
					}
					rv := matchRosterCall("player-ids", nil, pl, nil)
					state := o.state()
					if rv != 0 || len(state.Nodes) != len(want) {
						t.Fatal("ID report count", rv, len(state.Nodes), len(want))
					}
					for i, data := range want {
						if !bytes.Equal(state.Nodes[i].Data, data) {
							t.Fatal("ID report payload", i)
						}
					}
					rows = append(rows, matchRosterMessageRow{label, rv, state})
				})
			}
		}
	}
}
