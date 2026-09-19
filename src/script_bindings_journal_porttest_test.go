//go:build porttest

package opennox

import (
	"strings"
	"testing"

	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/legacy"
)

func TestScriptBindingsJournal(t *testing.T) {
	o := newJournalOwner(t)
	s := o.c.srv
	s.Server.NoxScriptVM.Init(s.Server)
	s.noxScript.Init(s)
	serverConfigOwnBytes(t, 0x5D4594, 2386620, 208)
	_, restore := legacy.PortTestMonsterCacheInit()
	t.Cleanup(restore)
	type row struct {
		Builtin asm.Builtin
		Target  int
		Status  uint32
		Root    bool
		State   journalResult
	}
	var rows []row
	for _, fi := range []asm.Builtin{asm.BuiltinJournalDelete, asm.BuiltinJournalEdit} {
		for _, target := range []int{0, 101, 107, 131, 999} {
			for _, status := range []uint32{0, 1, 0x7fff, 0x8000, 0xffff, 0x1234abcd, 0xffffffff} {
				for _, root := range []bool{false, true} {
					for _, exists := range []bool{false, true} {
						o.resetJournal(t)
						legacy.PortTestMonsterCache("reset", nil, 0)
						for i, slot := range []int{1, 7, 31} {
							o.units[i].ScriptIDVal = 100 + slot
							legacy.PortTestMonsterCache("prepare", &o.units[i], 0)
							if exists {
								o.call(t, 0, slot, "Short", 2)
								o.call(t, 0, slot, "Short", 4)
							}
						}
						o.queueReset()
						const sentinel = 0x2468ace0
						s.Server.NoxScriptVM.PushU32(sentinel)
						s.Server.NoxScriptVM.PushU32(uint32(target))
						s.Server.NoxScriptVM.PushString("Short")
						if fi == asm.BuiltinJournalEdit {
							s.Server.NoxScriptVM.PushU32(status)
						}
						if root {
							if err := s.noxScript.callBuiltinNative(fi); err != nil {
								t.Fatal(err)
							}
						} else if r, ok := legacy.CallScriptBuiltin(fi); r != 0 || !ok {
							t.Fatal("journal builtin dispatch", r, ok)
						}
						if s.Server.NoxScriptVM.PopU32() != sentinel {
							t.Fatal("journal builtin changed surrounding stack")
						}
						for _, slot := range []int{1, 7, 31} {
							n := o.players[slot].Journal
							if !exists {
								if n != nil {
									t.Fatal("journal builtin created absent entry")
								}
								continue
							}
							selected := target == 0 || target == 100+slot
							if selected && fi == asm.BuiltinJournalDelete {
								if n == nil || n.Field3 != 2 || n.Next != nil {
									t.Fatal("delete did not remove exactly first match")
								}
							} else {
								want := uint16(4)
								if selected {
									want = uint16(status)
								}
								if n == nil || n.Field3 != want || n.Next == nil || n.Next.Field3 != 2 || n.Next.Next != nil {
									t.Fatal("edit/untargeted journal mismatch")
								}
							}
						}
						rows = append(rows, row{fi, target, status, root, o.snapshot(t, 0, 0, true)})
					}
				}
			}
		}
	}
	spellbookCapture(t, "script-bindings-journal", rows, "78f23980643922909107aae7ed48c777dde2ef28fec99af204e01600472f0a2b")
}

func TestScriptBindingsJournalNameBoundary(t *testing.T) {
	o := newJournalOwner(t)
	s := o.c.srv
	s.Server.NoxScriptVM.Init(s.Server)
	s.noxScript.Init(s)
	serverConfigOwnBytes(t, 0x5D4594, 2386620, 208)
	_, restore := legacy.PortTestMonsterCacheInit()
	t.Cleanup(restore)
	for _, fi := range []asm.Builtin{asm.BuiltinJournalDelete, asm.BuiltinJournalEdit} {
		for _, target := range []uint32{0, 101, 999} {
			for _, name := range []string{"Short", "Short\x00ignored", "", "\x00Short", "Missing", "invalid-index"} {
				for _, root := range []bool{false, true} {
					o.resetJournal(t)
					legacy.PortTestMonsterCache("reset", nil, 0)
					for i, slot := range []int{1, 7, 31} {
						o.units[i].ScriptIDVal = 100 + slot
						legacy.PortTestMonsterCache("prepare", &o.units[i], 0)
						o.call(t, 0, slot, "Short", 2)
						o.call(t, 0, slot, "", 4)
					}
					const sentinel = 0x2468ace0
					s.Server.NoxScriptVM.PushU32(sentinel)
					s.Server.NoxScriptVM.PushU32(target)
					if name == "invalid-index" {
						s.Server.NoxScriptVM.PushU32(0x7fffffff)
					} else {
						s.Server.NoxScriptVM.PushString(name)
					}
					if fi == asm.BuiltinJournalEdit {
						s.Server.NoxScriptVM.PushU32(0x12345678)
					}
					if root {
						if err := s.noxScript.callBuiltinNative(fi); err != nil {
							t.Fatal(err)
						}
					} else if r, ok := legacy.CallScriptBuiltin(fi); r != 0 || !ok {
						t.Fatal("journal builtin dispatch", r, ok)
					}
					if s.Server.NoxScriptVM.PopU32() != sentinel {
						t.Fatal("journal name changed surrounding stack")
					}
					lookup := strings.SplitN(name, "\x00", 2)[0]
					if name == "invalid-index" {
						lookup = ""
					}
					for _, slot := range []int{1, 7, 31} {
						got := map[string]uint16{}
						for n := o.players[slot].Journal; n != nil; n = n.Next {
							got[strings.TrimRight(string(n.EntryBuf[:]), "\x00")] = n.Field3
						}
						want := map[string]uint16{"": 4, "Short": 2}
						if target == 0 || target == uint32(100+slot) {
							if _, exists := want[lookup]; exists {
								if fi == asm.BuiltinJournalDelete {
									delete(want, lookup)
								} else {
									want[lookup] = 0x5678
								}
							}
						}
						if len(got) != len(want) {
							t.Fatal("journal C-string entry count", fi, target, name, root, slot)
						}
						for key, value := range want {
							if v, ok := got[key]; !ok || v != value {
								t.Fatal("journal C-string lookup/status", fi, target, name, root, slot)
							}
						}
					}
				}
			}
		}
	}
}
