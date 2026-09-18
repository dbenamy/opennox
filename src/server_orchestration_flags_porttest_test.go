//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestServerOrchestrationDropFlags(t *testing.T) {
	o := newMatchRosterOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"RoundFlag", "OtherItem", "Glyph", "Torch", "Lantern"}, nil, true, 0, 0))
	globals, restore := legacy.PortTestServerOrchestrationGlobals()
	t.Cleanup(restore)
	*globals["drop-table"] = 1
	serverConfigOwnBytes(t, 0x587000, 279432, 192)
	clear(unsafe.Slice(memmap.PtrUint8(0x587000, 279432), 192))
	for _, off := range []uintptr{1568252, 1568256, 1568244} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	oldPending, oldRNG := o.s.Objs.Pending, o.s.Rand.Logic
	t.Cleanup(func() { o.s.Objs.Pending = oldPending; o.s.Rand.Logic = oldRNG })
	type row struct {
		Mask, Seed, RNG int
		Held            [3][3]bool
		Pos             [3][3]types.Pointf
		Queue           []uint32
		Packets         [][]byte
	}
	var rows []row
	for mask := 0; mask < 8; mask++ {
		for _, seed := range []int{1, 23, 253} {
			t.Run(fmt.Sprintf("mask%d/seed%d", mask, seed), func(t *testing.T) {
				o.reset()
				o.s.Rand.Logic = prand.New(seed)
				o.s.Objs.Pending = nil
				defer noxflags.PortTestGameFlags(2048)()
				var items [3][3]*server.Object
				defer func() {
					o.s.Objs.Pending = nil
					for i := range o.units {
						o.units[i].InvFirstItem = nil
						for _, u := range items[i] {
							if u != nil {
								legacy.PortTestPlayerStateMinimap("unmark", u, 0, ^uint32(0))
								u.InvHolder = nil
								o.s.Objs.FreeObject(u)
							}
						}
					}
				}()
				drops := 0
				for i := range o.units {
					owner := &o.units[i]
					owner.InvFirstItem = nil
					for j, name := range []string{"RoundFlag", "OtherItem", "RoundFlag"} {
						u := o.s.NewObjectByTypeID(name)
						if u == nil {
							t.Fatal("item")
						}
						items[i][j] = u
						u.NetCode = uint32(200 + i*3 + j)
						u.ObjFlags = 0
						u.ObjClass = 0
						if j != 1 {
							u.ObjClass = 0x10000000
						}
						u.PosVec = types.Pointf{X: 7, Y: 9}
						if mask&(1<<j) != 0 {
							u.InvHolder = owner
							if owner.InvFirstItem == nil {
								owner.InvFirstItem = u
							} else {
								prev := owner.InvFirstItem
								for prev.InvNextItem != nil {
									prev = prev.InvNextItem
								}
								prev.InvNextItem = u
								u.Field125 = prev
							}
							if j != 1 {
								drops++
							}
						}
					}
				}
				if legacy.PortTestServerOrchestration("drop-flags", nil, 0) != 0 {
					t.Fatal("loop return")
				}
				r := row{Mask: mask, Seed: seed, RNG: o.s.Rand.Logic.Index(), Packets: visibilityEffectsPackets(o.s)}
				rng := prand.New(seed)
				for i := 0; i < drops; i++ {
					rng.FloatClamp(-float64(float32(math.Pi)), float64(float32(math.Pi)))
				}
				if r.RNG != rng.Index() {
					t.Fatal("placement random consumption", r.RNG, rng.Index())
				}
				for i := range items {
					for j, u := range items[i] {
						r.Held[i][j] = u.InvHolder == &o.units[i]
						r.Pos[i][j] = u.PosVec
						held := j == 1 && mask&2 != 0
						if r.Held[i][j] != held {
							t.Fatal("inventory membership", i, j, r.Held)
						}
						dropped := j != 1 && mask&(1<<j) != 0
						if dropped {
							dx, dy := float64(u.PosVec.X-o.units[i].PosVec.X), float64(u.PosVec.Y-o.units[i].PosVec.Y)
							if math.Abs(dx*dx+dy*dy-2500) > 0.02 {
								t.Fatal("drop radius", dx, dy)
							}
							if objectXferGetWord(u.UpdateData, 8) != o.s.Frame() {
								t.Fatal("flag drop frame")
							}
						} else if u.PosVec != (types.Pointf{X: 7, Y: 9}) {
							t.Fatal("untouched position")
						}
					}
				}
				for u := o.s.Objs.Pending; u != nil; u = u.ObjNext {
					if len(r.Queue) >= drops {
						t.Fatal("pending cycle")
					}
					r.Queue = append(r.Queue, u.NetCode)
				}
				if len(r.Queue) != drops {
					t.Fatal("pending count", len(r.Queue), drops)
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "server-orchestration-drop-flags", rows, "c00d68b59aa24f3682ea00d5e567353cb977228ce1c6deeb4c928595b3f6048c")
}

func TestServerOrchestrationRoundFlags(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"RoundCrown", "Hecubah", "Necromancer"}, nil, true, 0, 0))
	for _, off := range []uintptr{1569740, 1569744} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	oldRNG := o.s.Rand.Logic
	t.Cleanup(func() { o.s.Rand.Logic = oldRNG })
	t.Cleanup(o.s.PortTestCombatAudioReset)
	type row struct {
		Assignment, Seed, Occupied, Capacity int
		Holder                               [3]int
		RNG                                  int
		Packets                              [][]byte
		Buffs                                [3]uint32
	}
	var rows []row
	for _, assignment := range []int{0, 5, 21, 26} {
		for _, seed := range []int{1, 23, 253} {
			for _, occupied := range []int{0, 1, 7} {
				for _, capacity := range []int{0, 64} {
					t.Run(fmt.Sprintf("assignment%d/seed%d/occupied%d/cap%d", assignment, seed, occupied, capacity), func(t *testing.T) {
						o.reset()
						o.s.PortTestCombatAudioReset()
						defer noxflags.PortTestGameFlags(1)()
						o.s.Teams.Reset()
						o.s.Teams.ActiveCnt = 0
						o.s.Rand.Logic = prand.New(seed)
						var teams [3]*server.Team
						var crowns [3]*server.Object
						defer func() {
							for i := range o.units {
								o.units[i].InvFirstItem = nil
								o.units[i].Field129 = nil
							}
							for _, u := range crowns {
								if u != nil {
									u.ObjOwner = nil
									u.InvHolder = nil
									o.s.Objs.FreeObject(u)
								}
							}
						}()
						for i := range teams {
							teams[i] = o.s.Teams.Create(server.TeamID(i + 1))
							crowns[i] = o.s.NewObjectByTypeID("RoundCrown")
							if crowns[i] == nil {
								t.Fatal("crown")
							}
							crowns[i].ObjFlags = 0
							crowns[i].NetCode = uint32(200 + i)
							objectXferSetWord(teams[i].C(), 76, uint32(uintptr(crowns[i].CObj())))
							if occupied&(1<<i) != 0 {
								crowns[i].InvHolder = &o.units[0]
							}
						}
						var counts [3]int
						var memberTeam [3]int
						v := assignment
						for i := range o.units {
							u := &o.units[i]
							u.TeamVal = server.ObjectTeam{}
							u.UpdateDataPlayer().Player.NetCodeVal = u.NetCode
							u.ObjFlags = 0
							u.Buffs = 0
							u.InvFirstItem = nil
							u.Field129 = nil
							u.CarryCapacity = uint16(capacity)
							objectXferSetWord(u.UpdateData, 264, 0)
							index := v % 3
							v /= 3
							memberTeam[i] = index
							counts[index]++
							legacy.Nox_xxx_createAtImpl_4191D0(teams[index].ID(), u.TeamPtr(), 0, int(u.NetCode), 0)
						}
						o.s.NetList.ResetAll()
						// Four active player records include the unitless record at slot 3.
						rng := prand.New(seed)
						want := [3]int{-1, -1, -1}
						for i := range teams {
							if occupied&(1<<i) != 0 {
								want[i] = 0
								continue
							}
							if counts[i] == 0 {
								continue
							}
							for tries := 0; ; tries++ {
								if tries > 1000 {
									t.Fatal("fixture random selection")
								}
								index := rng.IntClamp(0, 3)
								unit := -1
								switch index {
								case 0:
									unit = 0
								case 2:
									unit = 1
								case 3:
									unit = 2
								}
								if unit >= 0 && memberTeam[unit] == i {
									if capacity != 0 {
										want[i] = unit
									}
									break
								}
							}
						}
						legacy.PortTestServerOrchestration("flags", nil, 0)
						r := row{Assignment: assignment, Seed: seed, Occupied: occupied, Capacity: capacity, RNG: o.s.Rand.Logic.Index(), Packets: visibilityEffectsPackets(o.s)}
						if r.RNG != rng.Index() {
							t.Fatal("selection RNG", r.RNG, rng.Index())
						}
						for i, u := range crowns {
							r.Holder[i] = -1
							for j := range o.units {
								if u.InvHolder == &o.units[j] {
									r.Holder[i] = j
								}
							}
							if r.Holder[i] != want[i] {
								t.Fatal("selected carrier", r.Holder, want)
							}
							if occupied&(1<<i) == 0 && want[i] >= 0 {
								if u.ObjOwner != &o.units[want[i]] || objectXferGetWord(o.units[want[i]].UpdateData, 264) != o.s.Frame() {
									t.Fatal("crown owner/timestamp")
								}
							}
						}
						for i := range o.units {
							r.Buffs[i] = o.units[i].Buffs
							expected := false
							for j := range crowns {
								if occupied&(1<<j) == 0 && want[j] == i {
									expected = true
								}
							}
							if (r.Buffs[i]&(1<<30) != 0) != expected {
								t.Fatal("crown buff", i, r.Buffs)
							}
						}
						rows = append(rows, r)
					})
				}
			}
		}
	}
	spellbookCapture(t, "server-orchestration-round-flags", rows, "77fbf4346cfe311cdfcb2bfb63345a00d61945d5e98aa2e84b9e8238dd1a0842")
}
