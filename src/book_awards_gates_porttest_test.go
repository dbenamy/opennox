//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func bookAwardStrings(t *testing.T, s *server.Server) {
	t.Helper()
	var entries []strman.Entry
	for _, pair := range [][2]string{{"plyrspel.c:AwardSpellError", "Invalid spell"}, {"plyrspel.c:MaxSpellLevel", "Maximum level"}, {"PlyrGide.c:AwardGuideError", "Invalid guide"}, {"Ability.c:AwardAbilityError", "Invalid ability"}} {
		entries = append(entries, strman.Entry{ID: strman.ID(pair[0]), Vals: []strman.Variant{{Str: pair[1]}}})
	}
	configure, restore := s.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
}
func TestBookAwardsReports(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(0))
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	var rows []map[string]any
	raw := uint32(uintptr(unsafe.Pointer(u)))
	for _, player := range []bool{false, true} {
		for _, id := range []uint32{1, 5, 136} {
			for _, level := range []uint32{0, 5, 127, 128, 255, 256, 0xffffffff} {
				for _, notify := range []uint32{0, 1, 128, 255, 256, 0xffffffff} {
					for _, auto := range []uint32{0, 1, 2, 0xffffffff} {
						for kind, op := range []string{"nox_xxx_netSendSpellAward_4D7F90", "nox_xxx_netAbilityReport_4D8060", "nox_xxx_netReportGuideAward_4D8000"} {
							o.reset()
							u.ObjClass = 2
							if player {
								u.ObjClass = 4
							}
							p.SpellLvl[id] = level
							var ret uint32
							switch kind {
							case 0:
								ret = bookAwardCall(op, raw, id, notify, auto)
							case 1:
								ret = bookAwardCall(op, raw, id, auto)
							case 2:
								ret = bookAwardCall(op, raw, id, notify, auto)
							}
							state := o.state()
							normalized := ret
							if !player {
								if ret != raw || len(state.Nodes) != 0 {
									t.Fatal("nonplayer report", op, ret, len(state.Nodes))
								}
								normalized = 0
							} else {
								var want []byte
								switch kind {
								case 0:
									want = []byte{111, byte(id), byte(level), byte(notify)}
									if auto != 0 {
										want[3] |= 128
									}
								case 1:
									want = []byte{205, byte(id), byte(level)}
									if auto != 0 {
										want[2] |= 128
									}
								case 2:
									want = []byte{209, byte(id), byte(notify)}
									if auto != 0 {
										want[2] |= 128
									}
								}
								if len(state.Nodes) != 1 || !bytes.Equal(state.Nodes[0].Data, want) || state.Nodes[0].To != p.PlayerInd {
									t.Fatal("report", op, id, level, notify, auto, state.Nodes, want)
								}
							}
							rows = append(rows, map[string]any{"kind": kind, "player": player, "id": id, "level": level, "notify": notify, "auto": auto, "return": normalized, "queue": state, "direct": visibilityEffectsPackets(o.s)})
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "book-awards-reports", rows, "bee7c31f90aa97f5454f2be51ee6ea2629bd9f4de1fc0f7e65619e3b64febd13")
}
func TestBookAwardsSpellGates(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(0))
	bookAwardStrings(t, o.s)
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestBookSpellOwner())
	configure([]server.PortTestSpellClassDef{{Index: 1, Flags: 0x100, Valid: true}, {Index: 34, Flags: 0x100, Valid: true}, {Index: 136, Flags: 0x100, Valid: true}})
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	u.ObjClass = 4
	ids, resetRecords, records, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	raw := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
	binary.LittleEndian.PutUint32(raw[4636:], ids[0])
	binary.LittleEndian.PutUint32(raw[4640:], ids[1])
	var rows []map[string]any
	for _, gf := range []uint32{0, 2048, 4096, 6144} {
		flags.ResetGame()
		flags.SetGame(flags.GameFlag(gf))
		for _, id := range []uint32{1, 34, 136} {
			for _, old := range []uint32{0, 1, 2, 3, 4, 5, 6, 0xffffffff} {
				for _, override := range []uint32{0, 1, 3, 6, 0xffffffff} {
					o.reset()
					resetRecords()
					clear(p.SpellLvl[:])
					p.SpellLvl[id] = old
					before := p.SpellLvl
					expected := before
					ok := old != 5 && !(gf&6144 != 0 && old == 3) && !(gf&4096 != 0 && id == 34 && old != 0)
					wantRet := uint32(0)
					bits := uint32(0)
					if ok {
						wantRet = 1
						value := old + 1
						if value > 5 {
							value = 5
						}
						if gf&4096 != 0 && value > 3 {
							value = 3
						}
						if override != 0 {
							value = override
						}
						expected[id] = value
						if value != 0 {
							bits = 1 << (id & 31)
						}
					}
					got := bookAwardCall("nox_xxx_spellGrantToPlayer_4FB550", uint32(uintptr(unsafe.Pointer(u))), id, 0, 0, override)
					rec := records()
					if got != wantRet || p.SpellLvl != expected || rec != [3]uint32{bits, 0, 0x2468ace0 ^ bits} {
						t.Fatal("spell gates", gf, id, old, override, got, wantRet, p.SpellLvl[id], expected[id], rec)
					}
					state := o.state()
					if ok && (len(state.Nodes) != 1 || !bytes.Equal(state.Nodes[0].Data, []byte{111, byte(id), byte(expected[id]), 0})) {
						t.Fatal("spell award packet", state.Nodes)
					}
					message := ""
					if !ok {
						message = "Maximum level"
					}
					bookAwardCheckDirect(t, o, message, false)
					rows = append(rows, map[string]any{"flags": gf, "id": id, "old": old, "override": override, "return": got, "levels": p.SpellLvl, "records": rec, "queue": state, "direct": visibilityEffectsPackets(o.s)})
				}
			}
		}
	}
	for _, player := range []bool{false, true} {
		for _, id := range []uint32{0, 137, 0xffffffff} {
			o.reset()
			resetRecords()
			clear(p.SpellLvl[:])
			u.ObjClass = 2
			if player {
				u.ObjClass = 4
			}
			got := bookAwardCall("nox_xxx_spellGrantToPlayer_4FB550", uint32(uintptr(unsafe.Pointer(u))), id, 0, 0, 0)
			if got != 0 || p.SpellLvl != [137]uint32{} || records() != [3]uint32{0, 0, 0x2468ace0} {
				t.Fatal("invalid spell", player, id, got)
			}
			message := ""
			if player {
				message = "Invalid spell"
			}
			bookAwardCheckDirect(t, o, message, false)
			rows = append(rows, map[string]any{"player": player, "invalid": id, "return": got, "queue": o.state(), "direct": visibilityEffectsPackets(o.s)})
		}
	}
	spellbookCapture(t, "book-awards-spell-gates", rows, "805d5b23953e73421d91fee801c8d3159cbc1bbe2135a4ba1e41b6b16d63b841")
}
func TestBookAwardsAbilityGuideGates(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(0))
	bookAwardStrings(t, o.s)
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	ids, resetRecords, records, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	raw := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
	binary.LittleEndian.PutUint32(raw[4636:], ids[0])
	binary.LittleEndian.PutUint32(raw[4640:], ids[1])
	bookAwardWords(t, 0x587000, 216292, 1) // No guide families in this gate matrix.
	var rows []map[string]any
	for kind, op := range []string{"nox_xxx_abilityRewardServ_4FB9C0_ability", "nox_xxx_awardBeastGuide_4FAE80_magic_plyrgide"} {
		for _, player := range []bool{false, true} {
			for _, id := range []uint32{0, 1, 5, 6, 40, 41, 0xffffffff} {
				for _, old := range []uint32{0, 1, 3, 5, 0xffffffff} {
					o.reset()
					resetRecords()
					clear(p.SpellLvl[:])
					clear(p.BeastScrollLvl[:])
					u.ObjClass = 2
					if player {
						u.ObjClass = 4
					}
					level := uint32(5)
					limit := uint32(6)
					if kind == 1 {
						level = 1
						limit = 41
					}
					valid := id > 0 && id < limit
					if valid {
						if kind == 0 {
							p.SpellLvl[id] = old
						} else {
							p.BeastScrollLvl[id] = old
						}
					}
					expectedSpells, expectedGuides := p.SpellLvl, p.BeastScrollLvl
					wantRet := uint32(0)
					wantBits := [3]uint32{0, 0, 0x2468ace0}
					if player && valid && old == 0 {
						wantRet = 1
						wantBits[kind] = 1 << (id & 31)
						wantBits[2] ^= wantBits[kind]
						if kind == 0 {
							expectedSpells[id] = level
						} else {
							expectedGuides[id] = level
						}
					}
					got := bookAwardCall(op, uint32(uintptr(unsafe.Pointer(u))), id, 0)
					if got != wantRet || p.SpellLvl != expectedSpells || p.BeastScrollLvl != expectedGuides || records() != wantBits {
						t.Fatal("award gates", kind, player, id, old, got, wantRet, records(), wantBits)
					}
					state := o.state()
					if wantRet == 1 {
						want := []byte{205, byte(id), 5}
						if kind == 1 {
							want = []byte{209, byte(id), 0}
						}
						if len(state.Nodes) != 1 || !bytes.Equal(state.Nodes[0].Data, want) {
							t.Fatal("award packet", kind, state.Nodes, want)
						}
					}
					message, private := "", false
					if player {
						if !valid {
							message = "Invalid ability"
							if kind == 1 {
								message = "Invalid guide"
							}
						} else if kind == 0 && old != 0 {
							message = "use.c:HadAbility"
							private = true
						}
					}
					bookAwardCheckDirect(t, o, message, private)
					rows = append(rows, map[string]any{"kind": kind, "player": player, "id": id, "old": old, "return": got, "spells": p.SpellLvl, "guides": p.BeastScrollLvl, "records": records(), "queue": state, "direct": visibilityEffectsPackets(o.s)})
				}
			}
		}
	}
	spellbookCapture(t, "book-awards-ability-guide-gates", rows, "1a3f6f6bea2dfe9e7576058a8f2e9c77f0a3eeea3fd8b8940ee3bf4813a43aee")
}

// The original award paths have two independent queues: notifications are
// reliable, while private and translated rejection messages are direct sends.
func bookAwardCheckDirect(t *testing.T, o *reliableReportsOwner, text string, private bool) {
	t.Helper()
	var want []byte
	if text != "" {
		if private {
			want = append([]byte{169, 15, 0}, []byte(text)...)
			want = append(want, 0)
		} else {
			want = make([]byte, 12+len(text))
			want[0] = 168
			want[3] = 2
			want[8] = byte(len(text) + 1)
			copy(want[11:], text)
		}
	}
	to := int((*server.PlayerUpdateData)(o.units[0].UpdateData).Player.PlayerInd)
	for i, got := range visibilityEffectsPackets(o.s) {
		var expected []byte
		if i == to {
			expected = want
		}
		if !bytes.Equal(got, expected) {
			t.Fatal("direct award message", i, text, private, got, expected)
		}
	}
}
