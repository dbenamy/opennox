//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestBookAwardsSpellFamilies(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(0))
	bookAwardStrings(t, o.s)
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestBookSpellOwner())
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	u.ObjClass = 4
	ids, resetRecords, records, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	*(*uint32)(unsafe.Add(unsafe.Pointer(p), 4636)) = ids[0]
	var rows []map[string]any
	for _, family := range []struct{ base, child uint32 }{{0x1000, 0x2000}, {0x4000, 0x8000}, {0x10000, 0x20000}, {0x15000, 0x2000}} {
		configure([]server.PortTestSpellClassDef{{Index: 1, Flags: family.base, Valid: true}, {Index: 2, Flags: family.child, Valid: true}, {Index: 3, Flags: family.child, Valid: true}, {Index: 4, Flags: family.child, Valid: false}, {Index: 5, Flags: 0x100, Valid: true}})
		for _, gf := range []uint32{0, 2048, 4096} {
			flags.ResetGame()
			flags.SetGame(flags.GameFlag(gf))
			for _, old := range []uint32{0, 2, 3, 4, 5, 0xffffffff} {
				for _, override := range []uint32{0, 1, 3, 6, 0xffffffff} {
					o.reset()
					resetRecords()
					clear(p.SpellLvl[:])
					p.SpellLvl[1] = 2
					p.SpellLvl[2] = old
					p.SpellLvl[3] = old
					p.SpellLvl[4] = 17
					p.SpellLvl[5] = 23
					expected := p.SpellLvl
					expected[1] = 3
					if override != 0 {
						expected[1] = override
					}
					v := old + 1
					if override != 0 {
						v = override
					}
					if v > 5 {
						v = 5
					}
					expected[2] = v
					expected[3] = v
					// Historical family handling caps and folds the original spell ID, not
					// each family member. Preserve this observation; review separately.
					if gf&4096 != 0 && expected[1] > 3 {
						expected[1] = 3
					}
					got := bookAwardCall("nox_xxx_spellGrantToPlayer_4FB550", uint32(uintptr(unsafe.Pointer(u))), 1, 0, 0, override)
					if got != 1 || p.SpellLvl != expected || records() != [3]uint32{2, 0, 0x2468ace2} {
						t.Fatal("family", family, gf, old, override, got, p.SpellLvl[:6], expected[:6], records())
					}
					state := o.state()
					if len(state.Nodes) != 1 || !bytes.Equal(state.Nodes[0].Data, []byte{111, 1, byte(expected[1]), 0}) {
						t.Fatal("family packet", state.Nodes)
					}
					rows = append(rows, map[string]any{"baseFlag": family.base, "flags": gf, "old": old, "override": override, "return": got, "levels": p.SpellLvl, "records": records(), "queue": state, "direct": visibilityEffectsPackets(o.s)})
				}
			}
		}
	}
	spellbookCapture(t, "book-awards-spell-families", rows, "59287c41d176682d7f0889b2b43d795d25ab5f931349f3371e38f0941a803ac0")
}
func TestBookAwardsNotifications(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(0))
	bookAwardStrings(t, o.s)
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestBookSpellOwner())
	configure([]server.PortTestSpellClassDef{{Index: 1, Flags: 0x100, Valid: true}})
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	u.ObjClass = 4
	ids, resetRecords, records, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	raw := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
	binary.LittleEndian.PutUint32(raw[4636:], ids[0])
	binary.LittleEndian.PutUint32(raw[4640:], ids[1])
	family1, f1 := alloc.New([4]uint32{})
	t.Cleanup(f1)
	*family1 = [4]uint32{24, 7, 8, 0}
	family2, f2 := alloc.New([4]uint32{})
	t.Cleanup(f2)
	*family2 = [4]uint32{24, 8, 9, 0}
	table := bookAwardWords(t, 0x587000, 216292, 3)
	table[0] = uint32(uintptr(unsafe.Pointer(family1)))
	table[1] = uint32(uintptr(unsafe.Pointer(family2)))
	t.Cleanup(func() { o.s.Players.SetXxx(u, 0); o.s.PortTestCombatAudioReset() })
	var rows []map[string]any
	for kind, op := range []string{"nox_xxx_spellGrantToPlayer_4FB550", "nox_xxx_abilityRewardServ_4FB9C0_ability", "nox_xxx_awardBeastGuide_4FAE80_magic_plyrgide"} {
		for _, gf := range []uint32{0, 2048, 4096, 6144} {
			flags.ResetGame()
			flags.SetGame(flags.GameFlag(gf))
			for _, notify := range []uint32{0, 1, 2} {
				for _, loading := range []int32{0, 1} {
					for _, questNotify := range []uint32{0, 1} {
						o.reset()
						resetRecords()
						o.s.PortTestCombatAudioReset()
						clear(p.SpellLvl[:])
						clear(p.BeastScrollLvl[:])
						p.Field4792 = questNotify
						o.s.Players.SetXxx(u, loading)
						id := uint32(1)
						if kind == 2 {
							id = 24
						}
						got := bookAwardCall(op, uint32(uintptr(unsafe.Pointer(u))), id, notify, 0, 0)
						if got != 1 {
							t.Fatal("notify return", kind, gf, notify, got)
						}
						notifications := false
						others := false
						audioCount := 0
						switch kind {
						case 0:
							notifications = notify != 0 && (gf&4096 == 0 || questNotify != 0)
							others = notifications && loading == 0
							if notify != 0 {
								audioCount = 1
							}
						case 1:
							notifications = gf&4096 != 0
							others = notifications && loading == 0
						case 2:
							notifications = notify != 0
							others = notifications
							if notify != 0 {
								audioCount = 1
							}
						}
						wantPackets := map[byte][][]byte{}
						award := []byte{111, 1, 1, byte(notify)}
						awardKind := byte(0)
						if kind == 1 {
							award = []byte{205, 1, 5}
							if notify != 0 {
								award[2] |= 128
							}
							awardKind = 2
						}
						if kind == 2 {
							award = []byte{209, 24, byte(notify)}
							awardKind = 1
						}
						wantPackets[p.PlayerInd] = append(wantPackets[p.PlayerInd], award)
						msg := []byte{240, 30 + awardKind, byte(id), byte(u.NetCode), byte(u.NetCode >> 8)}
						if notifications {
							wantPackets[p.PlayerInd] = append(wantPackets[p.PlayerInd], msg)
						}
						if others {
							for i := range o.units {
								v := &o.units[i]
								if v != u {
									ind := v.UpdateDataPlayer().Player.PlayerInd
									wantPackets[ind] = append(wantPackets[ind], msg)
								}
							}
						}
						state := o.state()
						for _, node := range state.Nodes {
							found := false
							list := wantPackets[node.To]
							for i, data := range list {
								if bytes.Equal(data, node.Data) {
									wantPackets[node.To] = append(list[:i], list[i+1:]...)
									found = true
									break
								}
							}
							if !found {
								t.Fatal("unexpected notification", kind, gf, notify, loading, questNotify, node.To, node.Data)
							}
						}
						for ind, left := range wantPackets {
							if len(left) != 0 {
								t.Fatal("missing notification", kind, gf, notify, loading, questNotify, ind, left)
							}
						}
						events := o.s.PortTestCombatAudioSnapshot()
						if len(events) != audioCount {
							t.Fatal("audio count", kind, notify, len(events), audioCount)
						}
						for _, ev := range events {
							want := 226
							if kind == 2 {
								want = 227
							}
							if int(ev.ID) != want || ev.Obj != u || ev.Kind != 0 || ev.Code != 0 {
								t.Fatal("award sound", kind, ev)
							}
						}
						bits := [3]uint32{2, 0, 0x2468ace2}
						if kind == 2 {
							bits = [3]uint32{0, (1 << 24) | (1 << 7) | (1 << 8) | (1 << 9), 0}
							bits[2] = 0x2468ace0 ^ bits[1]
						}
						if records() != bits {
							t.Fatal("notify records", kind, records(), bits)
						}
						rows = append(rows, map[string]any{"kind": kind, "flags": gf, "notify": notify, "loading": loading, "questNotify": questNotify, "return": got, "spells": p.SpellLvl, "guides": p.BeastScrollLvl, "records": records(), "audio": audioCount, "queue": state, "direct": visibilityEffectsPackets(o.s)})
					}
				}
			}
		}
	}
	spellbookCapture(t, "book-awards-notifications", rows, "4bf2890ec127bf24d27a049029c865a2273d1d1da644d692903a5ffd80a8f22f")
}
