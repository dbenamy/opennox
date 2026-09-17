//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"math"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

type questRuntimeOwner struct {
	*worldCollisionOwner
	quest map[string]*uint32
}

func newQuestRuntimeOwner(t *testing.T) *questRuntimeOwner {
	o := &questRuntimeOwner{worldCollisionOwner: newWorldCollisionOwner(t)}
	var restore func()
	o.quest, restore = legacy.PortTestQuestRuntimeGlobals()
	t.Cleanup(restore)
	raw := unsafe.Slice(memmap.PtrUint8(0x581450, 10088), 12)
	old := bytes.Clone(raw)
	t.Cleanup(func() { copy(raw, old) })
	copy(raw, blobdata.PortTestQuestScoreConstant())
	return o
}
func questRuntimeCall(op string, u *server.Object, args ...uint32) uint64 {
	var a [4]uint32
	copy(a[:], args)
	return legacy.PortTestQuestRuntime(op, u, a)
}
func questRuntimeStats(p unsafe.Pointer) [11]uint32 {
	var out [11]uint32
	for i := range out {
		out[i] = objectXferGetWord(p, 4652+4*i)
	}
	return out
}
func TestQuestRuntimeStatistics(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	u := &o.units[0]
	pl := u.UpdateDataPlayer().Player
	type row struct {
		Name   string
		Return uint64
		Words  [11]uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-statistics", rows, "504dbd76e14d3c504490d10d4fc2b1b839165f7726aedd8dd02cc687dd4f3950")
	}()
	for _, op := range []string{"sub_4D6000", "sub_4D60E0", "sub_4D6130", "sub_4D6170", "sub_4D61B0", "sub_4D61F0"} {
		for _, flags := range []uint32{0, 0x20, 0x8000, 0x8020} {
			for _, participation := range []uint32{0, 1, 2} {
				for _, seed := range []uint32{0, 1, 63, 0x7fffffff, 0x80000000, 0xfffffffe, 0xffffffff} {
					name := fmt.Sprintf("%s/flags%x/participation%d/seed%x", op, flags, participation, seed)
					t.Run(name, func(t *testing.T) {
						u.ObjFlags = object.Flags(flags)
						objectXferSetWord(pl.C(), 4792, participation)
						for i := 0; i < 11; i++ {
							objectXferSetWord(pl.C(), 4652+4*i, seed)
						}
						*o.quest["202028"] = 0x80001234
						want := questRuntimeStats(pl.C())
						rvWant := uint64(2)
						if op == "sub_4D6000" {
							want = [11]uint32{}
							want[9] = 0x80001234
							want[10] = 63
						} else if flags&0x20 != 0 {
							rvWant = 1
						} else {
							switch op {
							case "sub_4D60E0":
								if participation == 1 {
									want[0]++
									want[10] |= 1
								}
							case "sub_4D6130":
								want[2]++
								want[10] |= 2
							case "sub_4D6170":
								want[3]++
								want[10] |= 4
							case "sub_4D61B0":
								want[4]++
								want[10] |= 8
							case "sub_4D61F0":
								want[5]++
								want[6]++
								want[10] |= 16
							}
						}
						if op == "sub_4D61B0" {
							rvWant = 0
						}
						rv := questRuntimeCall(op, u)
						got := questRuntimeStats(pl.C())
						if rv != rvWant || got != want {
							t.Fatalf("return=%d want%d words=%v want%v", rv, rvWant, got, want)
						}
						rows = append(rows, row{name, rv, got})
					})
				}
			}
		}
	}
	for _, op := range []string{"sub_4D6000", "sub_4D6130", "sub_4D6170", "sub_4D61B0", "sub_4D61F0"} {
		if got := questRuntimeCall(op, nil); got != 0 {
			t.Fatalf("%s nil return=%d", op, got)
		}
	}
}
func TestQuestRuntimeScalarState(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	type row struct {
		Name         string
		Return, Read uint64
		Stored       uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-scalar-state", rows, "0acd99e577807559d73fc0c72e190c6e1a62e078f8e3bc5da8f70bf99b304a34")
	}()
	pairs := [][3]string{{"nox_xxx_setQuest_4D6F60", "nox_xxx_isQuest_4D6F50", "1556160"}, {"sub_4D6F80", "sub_4D6F70", "1556164"}, {"sub_4D7440", "sub_4D7430", "1556116"}, {"sub_4D76E0", "sub_4D76F0", "1556124"}, {"nox_game_setQuestStage_4E3CD0", "nox_game_getQuestStage_4E3CC0", "202028"}}
	for _, p := range pairs {
		for _, v := range []uint32{0, 1, 2, 255, 65535, 65536, 0x7fffffff, 0x80000000, 0xffffffff} {
			name := fmt.Sprintf("%s/%x", p[0], v)
			t.Run(name, func(t *testing.T) {
				rv := questRuntimeCall(p[0], nil, v)
				rd := questRuntimeCall(p[1], nil)
				want := uint64(v)
				if p[0] == "nox_game_setQuestStage_4E3CD0" {
					want = 0
				}
				if rv != want || rd != uint64(v) || *o.quest[p[2]] != v {
					t.Fatal("state round trip", rv, rd, *o.quest[p[2]])
				}
				rows = append(rows, row{name, rv, rd, *o.quest[p[2]]})
			})
		}
	}
	for _, old := range []uint32{0, 0x12345678, 0xffffffff} {
		for _, v := range []uint32{0, 1, 0x80000000, 0xffffffff} {
			name := fmt.Sprintf("previous/%x/%x", old, v)
			t.Run(name, func(t *testing.T) {
				*o.quest["previousStage"] = old
				rv := questRuntimeCall("sub_4D72D0", nil, v)
				rd := questRuntimeCall("sub_4D7300", nil)
				if rv != uint64(old) || rd != uint64(old) || *o.quest["previousStage"] != v {
					t.Fatal("previous-stage exchange")
				}
				rows = append(rows, row{name, rv, rd, *o.quest["previousStage"]})
			})
		}
	}
	for _, v := range []uint32{0, 1, 0x7fffffff, 0xffffffff} {
		for _, p := range [][2]string{{"sub_4D71E0", "soulFrame"}, {"sub_4D75F0", "1556108"}} {
			dst := o.quest[p[1]]
			if p[1] == "soulFrame" {
				dst = o.globals[p[1]]
			}
			rv := questRuntimeCall(p[0], nil, v)
			if rv != uint64(v) || *dst != v {
				t.Fatal("frame store")
			}
			rows = append(rows, row{p[0] + fmt.Sprint(v), rv, 0, *dst})
		}
		*o.quest["1556104"] = v
		rv := questRuntimeCall("sub_4D6FA0", nil)
		if rv != uint64(v) {
			t.Fatal("state getter")
		}
		rows = append(rows, row{"get-state" + fmt.Sprint(v), rv, 0, v})
	}
}
func TestQuestRuntimeParticipants(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	oldEngine := noxflags.GetEngine()
	t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) })
	type row struct {
		Name        string
		Count, Room uint64
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-participants", rows, "3479d7c06b3287b245a091bc1707525ff00ae47d1e9139d91a95503a2f59d34f")
	}()
	for active := 0; active < 8; active++ {
		for _, participation := range []uint32{0, 1, 2, 0xffffffff} {
			for _, hosting := range []bool{false, true} {
				for _, headless := range []bool{false, true} {
					name := fmt.Sprintf("active%d/participation%x/hosting%t/headless%t", active, participation, hosting, headless)
					t.Run(name, func(t *testing.T) {
						flags := noxflags.GameFlag(0)
						if hosting {
							flags = 1
						}
						restore := noxflags.PortTestGameFlags(flags)
						defer restore()
						noxflags.UnsetEngine(noxflags.EngineNoRendering)
						if headless {
							noxflags.SetEngine(noxflags.EngineNoRendering)
						}
						want := uint64(0)
						for i := range o.units {
							u := &o.units[i]
							pl := u.UpdateDataPlayer().Player
							pl.Active = 0
							pl.PlayerUnit = nil
							if active&(1<<uint(i)) != 0 {
								pl.Active = 1
								pl.PlayerUnit = u
								if participation == 1 && !(i == 2 && hosting && headless) {
									want++
								}
							}
							objectXferSetWord(pl.C(), 4792, participation)
						}
						rv := questRuntimeCall("nox_xxx_player_4E3CE0", nil)
						room := questRuntimeCall("sub_4E4100", nil)
						if rv != want || room != 1 {
							t.Fatalf("count%d want%d room%d", rv, want, room)
						}
						rows = append(rows, row{name, rv, room})
					})
				}
			}
		}
	}
}
func TestQuestRuntimeBookEligibility(t *testing.T) {
	_ = newQuestRuntimeOwner(t)
	type row struct {
		Name   string
		Return uint64
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-book-eligibility", rows, "79d29d1b88236879e4a7bd414b77f3b5c351dd0ad6e88a2cf13cd1d59d84d1e8")
	}()
	for _, quest := range []bool{false, true} {
		for _, op := range []string{"nox_xxx_bookCreatureTest_4D70C0", "sub_4D7100"} {
			for _, id := range []uint32{0, 36, 37, 38, 39, 40, 41, 110, 111, 112, 113, 114, 115, math.MaxUint32} {
				name := fmt.Sprintf("%s/quest%t/id%d", op, quest, id)
				t.Run(name, func(t *testing.T) {
					flags := noxflags.GameFlag(0)
					if quest {
						flags = noxflags.GameModeQuest
					}
					restore := noxflags.PortTestGameFlags(flags)
					defer restore()
					want := uint64(1)
					if !quest && (op == "nox_xxx_bookCreatureTest_4D70C0" && id >= 37 && id <= 40 || op == "sub_4D7100" && id >= 111 && id <= 114) {
						want = 0
					}
					rv := questRuntimeCall(op, nil, id)
					if rv != want {
						t.Fatalf("got%d want%d", rv, want)
					}
					rows = append(rows, row{name, rv})
				})
			}
		}
	}
}
