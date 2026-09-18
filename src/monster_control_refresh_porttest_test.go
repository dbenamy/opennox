//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestMonsterControlRefresh(t *testing.T) {
	o := newCollisionCoreOwner(t)
	table := blobdata.PortTestMonsterActionTable()
	copy(serverConfigOwnBytes(t, 0x587000, 230388, len(table)), table)
	// The pointer slot for ESCORT is required: a zero-filled schema would hide removal.
	if binary.LittleEndian.Uint32(table[3*16:]) != 2 || binary.LittleEndian.Uint32(table[3*16+8:]) != 1 {
		t.Fatal("shipped escort schema")
	}
	a, b := newCreatureXferObject(t, o.s, "Monster"), newCreatureXferObject(t, o.s, "Monster")
	sa, sb := *a, *b
	ud := a.UpdateDataMonster()
	saved := *ud
	t.Cleanup(func() { *a = sa; *b = sb; *ud = saved })
	ids := collisionCoreIDs(a, b)
	type row struct {
		Action, Flags  uint32
		Stack          int8
		Mode           int
		Return, Target uint32
		Actions        [24][6]uint32
	}
	var rows []row
	for action := uint32(0); action < 72; action++ {
		for _, flags := range []object.Flags{4, 0x24, 0x8004, 0x8024} {
			for _, stack := range []int8{-1, 0, 1, 23} {
				for mode := 0; mode < 4; mode++ {
					*a = sa
					*b = sb
					*ud = server.MonsterUpdateData{}
					worldGeometryResetObject(a, 1001, 100, 100, false)
					worldGeometryResetObject(b, 1002, 133.25, 150.75, false)
					a.ObjClass = object.ClassMonster
					b.ObjClass = object.ClassMonster
					b.ObjFlags = flags
					ud.AIStackInd = stack
					target := uintptr(b.CObj())
					if mode == 0 {
						target = 0
					}
					ud.Field304 = uint32(target)
					if mode != 0 {
						ud.CurrentEnemy = b
					}
					if mode == 2 || mode == 3 {
						a.Buffs |= 1 << server.ENCHANT_BLINDED
					}
					for i := 0; i <= int(stack); i++ {
						ac := action
						if mode == 3 && i == 0 && stack > 0 {
							ac = 3
						}
						item := server.AIStackItem{Action: ac, Args: [4]uintptr{0x3f800000, 0x40000000, 0, 0xdeadbeef}, Field5: 0x12345678}
						count := binary.LittleEndian.Uint32(table[ac*16:])
						for j := uint32(0); j < count; j++ {
							if binary.LittleEndian.Uint32(table[ac*16+4+4*j:]) == 1 {
								item.Args[2*j] = target
							}
						}
						ud.AIStack[i] = item
					}
					rv := legacy.PortTestMonsterControl("refresh", a, b, nil, 0)
					if stack < 0 && int32(rv) != int32(stack) {
						t.Fatal("empty refresh return")
					}
					if stack >= 0 && rv != 0 {
						t.Fatal("nonempty refresh return")
					}
					if flags&0x8020 != 0 && ud.Field304 != 0 {
						t.Fatal("stale target retained")
					}
					if stack >= 0 && action == 3 && mode != 0 && flags&0x20 == 0 {
						h := ud.AIStack[stack]
						if uint32(h.Args[0]) != math.Float32bits(b.PosVec.X) || uint32(h.Args[1]) != math.Float32bits(b.PosVec.Y) {
							t.Fatal("escort position not refreshed")
						}
					}
					r := row{Action: action, Flags: uint32(flags), Stack: stack, Mode: mode, Return: rv, Target: collisionCoreID(t, ids, ud.Field304)}
					for i, s := range ud.AIStack {
						args := s.Args
						count := binary.LittleEndian.Uint32(table[s.Action*16:])
						for j := uint32(0); j < count; j++ {
							if binary.LittleEndian.Uint32(table[s.Action*16+4+4*j:]) == 1 {
								args[j*2] = uintptr(collisionCoreID(t, ids, uint32(args[j*2])))
							}
						}
						r.Actions[i] = [6]uint32{s.Action, uint32(args[0]), uint32(args[1]), uint32(args[2]), uint32(args[3]), s.Field5}
					}
					rows = append(rows, r)
				}
			}
		}
	}
	spellbookCapture(t, "monster-control-refresh", rows, "b27d9e78c0b2a54d9ee3aa8d8821999b1f0bcf5b0459516361812938ffa7046b")
}

func TestMonsterControlAnimation(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a := newCreatureXferObject(t, o.s, "Monster")
	sa := *a
	ud := a.UpdateDataMonster()
	saved := *ud
	t.Cleanup(func() { *a = sa; *ud = saved })
	data := unsafe.Slice((*byte)(collisionCoreGuarded(t, o, 256)), 256)
	type row struct {
		Frames, Delay, Loop, Frame, Tick, Done, Mode byte
		Result                                       [4]byte
	}
	var rows []row
	for _, frames := range []byte{0, 1, 2, 254, 255} {
		for _, delay := range []byte{0, 1, 254, 255} {
			for _, loop := range []byte{0, 1} {
				for _, frame := range []byte{0, 1, 253, 254, 255} {
					for _, tick := range []byte{0, 1, 253, 254, 255} {
						for _, done := range []byte{0, 1, 255} {
							for mode := byte(0); mode < 4; mode++ {
								*a = sa
								*ud = server.MonsterUpdateData{}
								a.ObjSubClass = 0
								clear(data)
								for i := 0; i < 16; i++ {
									data[16*i+9] = frames
									data[16*i+10] = delay
									binary.LittleEndian.PutUint32(data[16*i+12:], uint32(loop))
								}
								ud.Field119 = (*[16]server.MonsterAnim)(unsafe.Pointer(&data[0]))
								ud.AIStackInd = 0
								ud.AIStack[0].Action = 0
								ud.Field120_0 = 17
								ud.Field120_1 = frame
								ud.Field120_2 = tick
								ud.Field120_3 = done
								if mode == 1 {
									ud.Field119 = nil
								}
								if mode >= 2 {
									a.ObjSubClass = 0x10
									ud.AIStack[0].Action = uint32(14 + mode)
								}
								before := append([]byte(nil), data...)
								legacy.PortTestMonsterControl("animation", a, nil, nil, 0)
								got := [4]byte{ud.Field120_0, ud.Field120_1, ud.Field120_2, ud.Field120_3}
								if mode == 1 && got != [4]byte{17, frame, tick, done} {
									t.Fatal("nil animation changed state")
								}
								if mode >= 2 && got != [4]byte{17, frame, tick, 0} {
									t.Fatal("NPC attack animation reset")
								}
								if mode == 0 && done == 0 && frames == 0 && got != [4]byte{0, frame, tick, 1} {
									t.Fatal("empty animation completion")
								}
								if mode == 0 && done == 0 && frames > 0 && delay == 255 && (got[1] != frame || got[2] != tick+1 || got[3] != 0) {
									t.Fatal("maximum delay must never advance frame")
								}
								for i, v := range data {
									if v != before[i] {
										t.Fatal("animation table changed")
									}
								}
								rows = append(rows, row{frames, delay, loop, frame, tick, done, mode, got})
							}
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "monster-control-animation", rows, "3ea0a6d87f207730699c90b54a42c5b387aa8f4b5737dda8c09af58d83d610fa")
}
