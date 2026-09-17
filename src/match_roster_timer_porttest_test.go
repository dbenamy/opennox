//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
	"unsafe"
)

func TestMatchRosterTimerModes(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name                 string
		Return, Flags, Timer uint32
		Queue                legacy.PortTestReliableReportState
		Packets              [][]byte
		Sounds               [][3]uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "match-roster-timer", rows, "729852a8a0d84268d8ae5693807943696684392590796ff9e473988416fa6c70")
	}()
	for _, mode := range []uint32{0, 16, 32, 64, 96, 256, 1024} {
		for control := 0; control < 16; control++ {
			for _, host := range []uint32{0, 1} {
				label := fmt.Sprintf("mode%x/control%d/host%d", mode, control, host)
				t.Run(label, func(t *testing.T) {
					o.reset()
					o.s.PortTestCombatAudioReset()
					restore := noxflags.PortTestGameFlags(noxflags.GameFlag(mode | host))
					defer restore()
					for i := range unsafe.Slice(memmap.PtrUint8(0x5D4594, 3500), 6) {
						*memmap.PtrUint8(0x5D4594, 3500+uintptr(i)) = byte(bool2int(control&8 == 0))
					}
					noxServer.flag3592 = control&4 != 0
					enabled := control&1 != 0
					expired := control&2 != 0
					*memmap.PtrUint32(0x587000, 4660) = uint32(bool2int(enabled))
					deadline := uint64(math.MaxUint64)
					if expired {
						deadline = 0
					}
					*memmap.PtrUint64(0x5D4594, 3468) = deadline
					for i := range o.units {
						p := o.units[i].UpdateDataPlayer().Player
						p.Lessons = int32(i + 1)
						p.Field2140 = uint32(i + 1)
						p.Field3680 = 0
						objectXferSetWord(p.C(), 3632, math.Float32bits(float32(i)+12.75))
						objectXferSetWord(p.C(), 3636, math.Float32bits(-float32(i)-12.75))
					}
					rv := matchRosterCall("check-limit", nil, nil, nil)
					active := enabled && expired && (control&8 == 0 || noxServer.flag3592)
					expectedFlags := mode | host
					if active {
						if noxServer.flag3592 {
							expectedFlags |= 0x4000000
						} else if mode != 0 {
							expectedFlags |= 8
						}
					}
					timer := memmap.Uint32(0x587000, 4660)
					wantTimer := uint32(bool2int(enabled))
					if active {
						wantTimer = 0
					}
					if uint32(noxflags.GetGame()) != expectedFlags || timer != wantTimer || rv != uint32(bool2int(active && host != 0)) {
						t.Fatal("timer transition", rv, noxflags.GetGame(), timer)
					}
					packets := visibilityEffectsPackets(o.s)
					sounds := [][3]uint32{}
					for _, e := range o.s.PortTestCombatAudioSnapshot() {
						sounds = append(sounds, [3]uint32{uint32(e.ID), uint32(e.Kind), uint32(e.Code)})
					}
					if len(sounds) != 3*bool2int(active && noxServer.flag3592) {
						t.Fatal("timer sound count", sounds)
					}
					for i, u := range o.units {
						want := []byte(nil)
						if active && noxServer.flag3592 {
							want = []byte{154, byte(12 + i), 0, 0, 0}
							binary.LittleEndian.PutUint16(want[3:], uint16(int16(-12-i)))
							if sounds[i] != [3]uint32{582, 2, u.NetCode} {
								t.Fatal("timer sound", sounds[i])
							}
						}
						if !bytes.Equal(packets[u.UpdateDataPlayer().Player.PlayerInd], want) {
							t.Fatal("timer position message", i)
						}
					}
					rows = append(rows, row{label, rv, uint32(noxflags.GetGame()), timer, o.state(), packets, sounds})
				})
			}
		}
	}
}
