//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func questRuntimeEightPlayers(t *testing.T, o *questRuntimeOwner) []server.Object {
	units, freeU := alloc.Make([]server.Object{}, 8)
	data, freeD := alloc.Make([]server.PlayerUpdateData{}, 8)
	t.Cleanup(freeU)
	t.Cleanup(freeD)
	for i := 0; i < 32; i++ {
		p := o.s.Players.ByIndRaw(ntype.PlayerInd(i))
		p.Active = 0
		p.PlayerUnit = nil
	}
	for i := range units {
		slot := i
		if i == 7 {
			slot = 31
		}
		pl := o.s.Players.ByIndRaw(ntype.PlayerInd(slot))
		pl.PlayerInd = byte(slot)
		pl.Active = 1
		pl.PlayerUnit = &units[i]
		units[i].ObjClass = object.ClassPlayer
		units[i].NetCode = uint32(2001 + i)
		units[i].UpdateData = unsafe.Pointer(&data[i])
		data[i].Player = pl
		o.objects[units[i].CObj()] = units[i].NetCode
	}
	t.Cleanup(func() {
		for i := 0; i < 32; i++ {
			p := o.s.Players.ByIndRaw(ntype.PlayerInd(i))
			p.Active = 0
			p.PlayerUnit = nil
		}
	})
	return units
}
func TestQuestRuntimeScoreAggregation(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	units := questRuntimeEightPlayers(t, o)
	defer noxflags.PortTestGameFlags(0)()
	type row struct {
		Name  string
		Score uint64
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-score-aggregation", rows, "1e1f246a5595101f84b9bc725321eb221b1b4faf77dde65a195674043010ed2b")
	}()
	for _, mask := range []uint32{0, 1, 3, 127, 255} {
		for _, participation := range []uint32{1, 2} {
			for _, seed := range []uint32{0, 10, 10000000, 0xffffffff} {
				for i := range units {
					pl := units[i].UpdateDataPlayer().Player
					p := uint32(0)
					if mask&(1<<uint(i)) != 0 {
						p = participation
					}
					objectXferSetWord(pl.C(), 4792, p)
					objectXferSetWord(pl.C(), 4668, seed+uint32(i))
					objectXferSetWord(pl.C(), 4672, uint32(i*3))
					objectXferSetWord(pl.C(), 4664, uint32(i*7))
					objectXferSetWord(pl.C(), 4688, 1+uint32(i%3))
				}
				for i := range units {
					pl := units[i].UpdateDataPlayer().Player
					name := fmt.Sprintf("mask%x/participation%d/seed%d/slot%d", mask, participation, seed, pl.PlayerInd)
					t.Run(name, func(t *testing.T) {
						before := questRuntimeStats(pl.C())
						got := questRuntimeCall("sub_4D6540", nil, uint32(pl.PlayerInd))
						if got > 999999999 {
							t.Fatal("score cap", got)
						}
						if mask&(1<<uint(i)) == 0 && got != 0 {
							t.Fatal("nonparticipant score", got)
						}
						if participation == 1 && bits.OnesCount32(mask) == 1 && mask&(1<<uint(i)) != 0 {
							want := questRuntimeCall("sub_4D66E0", nil, seed+uint32(i), uint32(i*3), uint32(i*7), 1+uint32(i%3))
							if want > 999999999 {
								want = 999999999
							}
							if got != want {
								t.Fatalf("single participant%d want%d", got, want)
							}
						}
						if questRuntimeStats(pl.C()) != before {
							t.Fatal("scoring changed statistics")
						}
						rows = append(rows, row{name, got})
					})
				}
			}
		}
	}
	if got := questRuntimeCall("sub_4D6540", nil, 30); got != 0 {
		t.Fatal("missing player score", got)
	}
}
func TestQuestRuntimeScoreboard(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	units := questRuntimeEightPlayers(t, o)
	defer noxflags.PortTestGameFlags(0)()
	type row struct {
		Name   string
		Return uint64
		Queue  legacy.PortTestReliableReportState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-scoreboard", rows, "13534766ab7c5b64eb3bf59bf2fc997c3ce4b7cc402009331efc832153e75ca2")
	}()
	for _, mask := range []uint32{0, 1, 3, 63, 127, 255} {
		for _, part := range []uint32{1, 2} {
			for _, stage := range []uint32{0, 65535, 65536, 0xffffffff} {
				name := fmt.Sprintf("mask%x/participation%d/stage%x", mask, part, stage)
				t.Run(name, func(t *testing.T) {
					o.reset()
					*o.quest["1556132"] = 0x12345678
					var selected []int
					for i := range units {
						pl := units[i].UpdateDataPlayer().Player
						p := uint32(0)
						if mask&(1<<uint(i)) != 0 {
							p = part
						}
						objectXferSetWord(pl.C(), 4792, p)
						objectXferSetWord(pl.C(), 4688, stage)
						for field, off := range []int{4664, 4668, 4672, 4680} {
							objectXferSetWord(pl.C(), off, uint32(65534+i+100*field))
						}
						if p == 1 && len(selected) < 6 {
							selected = append(selected, i)
						}
					}
					rv := questRuntimeCall("sub_4D6770", nil, 1)
					st := o.state()
					if rv != 1 || len(st.Nodes) != 1 {
						t.Fatal("scoreboard enqueue", rv, len(st.Nodes))
					}
					data := st.Nodes[0].Data
					if len(data) != 90 || data[0] != 240 || data[1] != 12 || binary.LittleEndian.Uint16(data[2:]) != 0x5678 {
						t.Fatal("scoreboard header")
					}
					wantStage := uint16(0)
					if len(selected) > 0 {
						wantStage = uint16(stage)
					}
					if binary.LittleEndian.Uint16(data[4:]) != wantStage {
						t.Fatal("recipient stage")
					}
					for n, i := range selected {
						at := 6 + 14*n
						if binary.LittleEndian.Uint16(data[at:]) != uint16(units[i].NetCode) {
							t.Fatal("scoreboard order")
						}
						for j := 1; j <= 4; j++ {
							if binary.LittleEndian.Uint16(data[at+2*j:]) != uint16(65534+i+[]int{0, 100, 200, 300, 0}[j]) {
								t.Fatal("stat word narrowing")
							}
						}
						pl := units[i].UpdateDataPlayer().Player
						score := questRuntimeCall("sub_4D6540", nil, uint32(pl.PlayerInd))
						if binary.LittleEndian.Uint32(data[at+10:]) != uint32(score) {
							t.Fatal("serialized score")
						}
					}
					for _, v := range data[6+14*len(selected):] {
						if v != 0 {
							t.Fatal("scoreboard padding or seventh participant")
						}
					}
					rows = append(rows, row{name, rv, st})
				})
			}
		}
	}
}
func TestQuestRuntimeParticipantCapacity(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	units := questRuntimeEightPlayers(t, o)
	old := noxflags.GetEngine()
	t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(old) })
	type row struct {
		Name        string
		Count, Room uint64
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-participant-capacity", rows, "81c4f5099a1cf9ce257d4bcc64fb1975e699623a6b3fdc11772936c46e3b7785")
	}()
	for count := 0; count <= 8; count++ {
		for _, part := range []uint32{0, 1, 2} {
			for _, host := range []bool{false, true} {
				for _, headless := range []bool{false, true} {
					name := fmt.Sprintf("count%d/part%d/host%t/headless%t", count, part, host, headless)
					t.Run(name, func(t *testing.T) {
						flags := noxflags.GameFlag(0)
						if host {
							flags = 1
						}
						restore := noxflags.PortTestGameFlags(flags)
						defer restore()
						noxflags.UnsetEngine(noxflags.EngineNoRendering)
						if headless {
							noxflags.SetEngine(noxflags.EngineNoRendering)
						}
						for i := range units {
							p := uint32(0)
							if i < count {
								p = part
							}
							objectXferSetWord(units[i].UpdateDataPlayer().Player.C(), 4792, p)
						}
						n := count
						if count == 8 && host && headless {
							n--
						}
						wc := uint64(0)
						if part == 1 {
							wc = uint64(n)
						}
						wr := uint64(1)
						if part != 0 && n >= 6 {
							wr = 0
						}
						c := questRuntimeCall("nox_xxx_player_4E3CE0", nil)
						r := questRuntimeCall("sub_4E4100", nil)
						if c != wc || r != wr {
							t.Fatalf("count%d want%d room%d want%d", c, wc, r, wr)
						}
						rows = append(rows, row{name, c, r})
					})
				}
			}
		}
	}
}
