//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/player"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestBookAwardsItemUse(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(0))
	bookAwardStrings(t, o.s)
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestBookSpellOwner())
	configure([]server.PortTestSpellClassDef{{Index: 1, Flags: 0x1000000, Valid: true}, {Index: 2, Flags: 0x2000000, Valid: true}})
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	ids, resetRecords, records, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	p.Prot4636 = ids[0]
	p.Prot4640 = ids[1]
	table := bookAwardWords(t, 0x587000, 70500, 41)
	for i := range table {
		table[i] = bookAwardString(t, fmt.Sprintf("UseGuide%d", i))
	}
	bookAwardWords(t, 0x587000, 216292, 1)
	it, freeItem := alloc.New(server.Object{})
	t.Cleanup(freeItem)
	data, freeData := alloc.New([256]byte{})
	t.Cleanup(freeData)
	oldDeleted := o.s.Objs.DeletedList
	t.Cleanup(func() { o.s.Objs.DeletedList = oldDeleted; o.s.PortTestCombatAudioReset() })
	var rows []map[string]any
	for kind, op := range []string{"nox_xxx_useSpellReward_53F9E0", "nox_xxx_useAbilityReward_53FAE0", "sub_53F930"} {
		for _, gf := range []uint32{0, 2048, 4096, 6144} {
			flags.ResetGame()
			flags.SetGame(flags.GameFlag(gf))
			for _, class := range []int{-1, 0, 1, 2} {
				for _, id := range []uint32{0, 1, 2} {
					for _, old := range []uint32{0, 1, 3, 5} {
						o.reset()
						resetRecords()
						o.s.PortTestCombatAudioReset()
						o.s.Objs.DeletedList = nil
						clear(p.SpellLvl[:])
						clear(p.BeastScrollLvl[:])
						p.SpellLvl[id] = old
						p.BeastScrollLvl[id] = old
						p.Field4792 = 0
						u.ObjClass = object.ClassPlayer
						if class < 0 {
							u.ObjClass = object.ClassMonster
						} else {
							p.Info().SetPlayerClass(player.Class(class))
						}
						*it = server.Object{}
						clear(data[:])
						it.UseData.Ptr = unsafe.Pointer(data)
						if kind == 2 {
							copy(data[:], fmt.Sprintf("UseGuide%d", id))
							if id == 0 {
								copy(data[:], "UnknownGuide")
							}
						} else {
							data[0] = byte(id)
						}
						ret := bookAwardCall(op, uint32(uintptr(unsafe.Pointer(u))), uint32(uintptr(unsafe.Pointer(it))))
						wantRet := uint32(0)
						consumed := false
						audio := 0
						if class >= 0 {
							switch kind {
							case 0:
								allowed := (class == 1 || class == 2) && id != 0 && (id != 2 || class == 1)
								if allowed {
									wantRet = 1
									consumed = old != 5 && !(gf&6144 != 0 && old == 3)
									if consumed {
										audio = 226
									} else {
										audio = 925
									}
								} else {
									audio = 925
								}
							case 1:
								if class == 0 {
									wantRet = 1
									consumed = id != 0 && old == 0
									if !consumed {
										audio = 925
									}
								} else {
									audio = 925
								}
							case 2:
								if !(gf&4096 != 0 && class != 2) && old == 0 {
									wantRet = 1
									consumed = true
									if id != 0 {
										audio = 227
									}
								}
							}
						}
						destroyed := it.Flags().Has(object.FlagDestroyed)
						if ret != wantRet || destroyed != consumed || (o.s.Objs.DeletedList == it) != consumed {
							t.Fatal("item use", kind, gf, class, id, old, ret, wantRet, destroyed, consumed)
						}
						if consumed && (it.DeletedAt != 123 || it.DeletedNext != nil) {
							t.Fatal("deletion bookkeeping", it.DeletedAt, it.DeletedNext)
						}
						events := o.s.PortTestCombatAudioSnapshot()
						if audio == 0 && len(events) != 0 || audio != 0 && len(events) != 1 {
							t.Fatal("item sound count", kind, gf, class, id, old, events, audio)
						}
						if audio != 0 {
							ev := events[0]
							if int(ev.ID) != audio || ev.Obj != u {
								t.Fatal("item sound", ev, audio)
							}
							wantKind, wantCode := 0, uint32(0)
							if audio == 925 {
								wantKind = 2
								wantCode = u.NetCode
							}
							if ev.Kind != wantKind || ev.Code != wantCode {
								t.Fatal("sound destination", ev, wantKind, wantCode)
							}
						}
						wantLevel := old
						if consumed && id != 0 {
							switch kind {
							case 0:
								wantLevel = old + 1
								if gf&4096 != 0 && wantLevel > 3 {
									wantLevel = 3
								}
							case 1:
								wantLevel = 5
							case 2:
								wantLevel = 1
							}
						}
						level := p.SpellLvl[id]
						if kind == 2 {
							level = p.BeastScrollLvl[id]
						}
						if level != wantLevel {
							t.Fatal("item level", kind, gf, class, id, old, level, wantLevel)
						}
						message, private := "", false
						if class >= 0 {
							switch kind {
							case 0:
								if ret == 0 {
									message = "use.c:SpellRewardClassFail"
									private = true
								} else if !consumed {
									message = "Maximum level"
								}
							case 1:
								if class != 0 {
									message = "pickup.c:ObjectEquipClassFail"
									private = true
								} else if id == 0 {
									message = "Invalid ability"
								} else if old != 0 {
									message = "use.c:HadAbility"
									private = true
								}
							case 2:
								if gf&4096 != 0 && class != 2 {
									message = "pickup.c:ObjectEquipClassFail"
									private = true
								} else if old != 0 {
									message = "objcoll.c:AlreadyHaveGuide"
									private = true
								} else if id == 0 {
									message = "Invalid guide"
								}
							}
						}
						bookAwardCheckDirect(t, o, message, private)
						rows = append(rows, map[string]any{"kind": kind, "flags": gf, "class": class, "id": id, "old": old, "return": ret, "consumed": destroyed, "audio": audio, "level": level, "records": records(), "queue": o.state(), "direct": visibilityEffectsPackets(o.s)})
					}
				}
			}
		}
	}
	spellbookCapture(t, "book-awards-item-use", rows, "ce8035fbce4dee00ac94d7250a5ea74c2cd5d70d148b08e760890bded0798a7f")
}
