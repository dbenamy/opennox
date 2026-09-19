//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesStatusRead(t *testing.T) {
	o := newReliableReportsOwner(t)
	defer flags.PortTestGameFlags(0)()
	u := &o.units[0]
	ud := u.UpdateDataPlayer()
	p := ud.Player
	p.NetCodeVal = u.NetCode
	hp, free := alloc.New(server.HealthData{})
	t.Cleanup(free)
	u.HealthData = hp
	raw := unsafe.Slice((*byte)(unsafe.Pointer(u)), int(unsafe.Sizeof(*u)))
	player := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
	data := unsafe.Slice((*byte)(unsafe.Pointer(ud)), int(unsafe.Sizeof(*ud)))
	saved := unsafe.Slice(memmap.PtrUint32(0x5D4594, 527696), 2)
	old := append([]uint32(nil), saved...)
	t.Cleanup(func() { copy(saved, old) })
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 2, 3, 0x8000, 0xffff} {
		for _, gf := range []flags.GameFlag{0, 2048, 2048 | 8192} {
			for _, present := range []byte{0, 1} {
				for _, v := range []uint16{0, 32768, 65535} {
					o.reset()
					flags.ResetGame()
					flags.SetGame(gf)
					*hp = server.HealthData{Cur: 17, Max: 31, Field16: 99}
					ud.ManaCur = 19
					ud.ManaPrev = 23
					ud.ManaMax = 37
					u.Poison540 = 0
					u.Field541 = 13
					u.Field542 = 77
					u.Experience = 19.5
					u.Field38 = 0x1234
					u.Direction1 = 7
					u.Direction2 = 9
					p.Field3680 = 0x8400
					saved[0] = 0xdeadbeef
					saved[1] = 0x12345678
					expectedU, expectedP, expectedD := bytes.Clone(raw), bytes.Clone(player), bytes.Clone(data)
					expectedHP := *hp
					expectedSaved := [2]uint32{saved[0], saved[1]}
					input := binary.LittleEndian.AppendUint16(nil, version)
					input = append(input, present)
					maxMana, curHP, curMana := v^0x1234, v^0x2468, v^0x4321
					for _, x := range []uint16{v, maxMana, curHP, curMana} {
						input = binary.LittleEndian.AppendUint16(input, x)
					}
					poison := byte(v)
					input = append(input, poison, 0xa6)
					input = binary.LittleEndian.AppendUint16(input, 0x5678)
					input = binary.LittleEndian.AppendUint32(input, 0x42f78000) // 123.75 experience
					if int16(version) >= 2 {
						input = binary.LittleEndian.AppendUint16(input, 0xabcd)
					}
					wantRet, wantPos := uint32(1), int64(len(input))
					body := true
					if int16(version) > 2 {
						wantRet = 0
						wantPos = 2
						body = false
					} else if present == 0 {
						wantPos = 3
						body = false
					} else if gf&2048 == 0 {
						wantRet = 0
						wantPos = 3
						body = false
					}
					if body {
						expectedHP.Cur = v
						expectedHP.Max = v
						if poison != 0 {
							expectedHP.Field16 = 123
						}
						expectedSaved = [2]uint32{uint32(curHP), uint32(curMana)}
						binary.LittleEndian.PutUint16(expectedD[4:], maxMana)
						binary.LittleEndian.PutUint16(expectedD[6:], 19)
						binary.LittleEndian.PutUint16(expectedD[8:], maxMana)
						binary.LittleEndian.PutUint32(expectedU[152:], 0xffffffff)
						expectedU[540] = poison
						expectedU[541] = 0xa6
						binary.LittleEndian.PutUint16(expectedU[542:], 0x5678)
						binary.LittleEndian.PutUint32(expectedU[28:], 0x42f78000)
						status := uint32(0x8000)
						if poison != 0 {
							status |= 1024
						}
						binary.LittleEndian.PutUint32(expectedP[3680:], status)
						if int16(version) >= 2 {
							binary.LittleEndian.PutUint16(expectedU[124:], 0xabcd)
							binary.LittleEndian.PutUint16(expectedU[126:], 0xabcd)
						}
					}
					input = append(input, 0xde, 0xad, 0xbe, 0xef)
					ret, got, pos := playerFileSection(t, "sub_41AA30", input, uint32(uintptr(unsafe.Pointer(u))), 0)
					if ret != wantRet || pos != wantPos || !bytes.Equal(got, input) || !bytes.Equal(raw, expectedU) || !bytes.Equal(player, expectedP) || !bytes.Equal(data, expectedD) || *hp != expectedHP || [2]uint32{saved[0], saved[1]} != expectedSaved {
						t.Fatal("status read", version, gf, present, v, ret, pos, wantPos)
					}
					state := o.state()
					wantNodes := 0
					if body {
						wantNodes = 1
					}
					if len(state.Nodes) != wantNodes {
						t.Fatal("status report count", len(state.Nodes), wantNodes)
					}
					if body && !bytes.Equal(state.Nodes[0].Data, []byte{110, 123, 0, 0, 0}) {
						t.Fatal("experience report", state.Nodes[0].Data)
					}
					rows = append(rows, map[string]any{"version": version, "flags": uint32(gf), "present": present, "value": v, "return": ret, "position": pos, "saved": expectedSaved, "health": *hp, "queue": state})
				}
			}
		}
	}
	spellbookCapture(t, "player-files-status-read", rows, "f44884ed7ec9ff2bce3e603d6b8668d2e599693083ca60cb0a8dac6e8d5dabde")
}
