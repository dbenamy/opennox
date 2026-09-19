//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesInventoryReadGates(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestInventoryEnvironment(false, nil, nil))
	t.Cleanup(o.s.PortTestInventoryDisplayBalance())
	oldStats := o.s.Players.Stats
	t.Cleanup(func() { o.s.Players.Stats = oldStats })
	o.s.Players.Stats.Base = server.ClassStats{Health: 25, Mana: 15, Speed: 1500, Strength: 10}
	o.s.Players.Stats.Wizard = server.ClassStats{Health: 80, Mana: 150, Speed: 3500, Strength: 20}
	carry := memmap.PtrFloat64(0x581450, 10216)
	oldCarry := *carry
	*carry = 1
	t.Cleanup(func() { *carry = oldCarry })
	configure, restore := o.s.PortTestWorldCollisionBalance()
	t.Cleanup(restore)
	configure(map[string]float64{"ForceOfNatureStaffLimit": 3})
	limitTypes := serverConfigOwnBytes(t, 0x5D4594, 1568356, 52)
	for i := 0; i < 13; i++ {
		binary.LittleEndian.PutUint32(limitTypes[i*4:], uint32(o.s.Types.IndByID("PortGold")))
	}
	serverConfigOwnBytes(t, 0x5D4594, 527704, 4)
	_, _, restore = legacy.PortTestMeterInventory()
	t.Cleanup(restore)
	_, restore = legacy.PortTestUIInventoryWords()
	t.Cleanup(restore)
	u := &o.units[0]
	ud := u.UpdateDataPlayer()
	p := ud.Player
	hp, free := alloc.New(server.HealthData{})
	t.Cleanup(free)
	u.HealthData = hp
	*(*byte)(unsafe.Add(unsafe.Pointer(p), 2251)) = 1
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 2, 3, 4, 0x8000, 0xffff} {
		for _, gf := range []flags.GameFlag{0, 2048, 4096, 8192} {
			for _, which := range []string{"absent", "empty", "unknown-type", "quest-count-limit"} {
				o.reset()
				flags.ResetGame()
				flags.SetGame(gf)
				*hp = server.HealthData{Cur: 17, Max: 25}
				ud.ManaCur = 11
				ud.ManaMax = 15
				u.Experience = 100
				u.InvFirstItem = nil
				p.GoldVal = 321
				field := (*byte)(unsafe.Add(u.UpdateData, 244))
				*field = 0xa7
				var input mapDrawableStream
				input.u16(version)
				present := byte(1)
				if which == "absent" {
					present = 0
				}
				input.u8(present)
				input.u32(500)
				count := uint32(0)
				if which == "unknown-type" {
					count = 1
				}
				if which == "quest-count-limit" {
					count = 2561
				}
				input.u32(count)
				if count != 0 {
					name := "PortMissingInventoryType"
					input.u8(byte(len(name)))
					input.WriteString(name)
				} else {
					input.u8(0)
					input.u32(0)
					if int16(version) >= 2 {
						input.u32(0)
					}
					if int16(version) >= 3 {
						input.u8(0x5c)
					}
				}
				wantRet := uint32(1)
				wantPos := input.Len()
				wantGold := uint32(321)
				wantField := byte(0xa7)
				if int16(version) > 3 {
					wantRet = 0
					wantPos = 2
				} else if present == 0 {
					wantPos = 3
				} else if gf != 2048 && gf != 4096 {
					wantRet = 0
					wantPos = 3
				} else {
					wantGold = 500
					if count != 0 {
						wantRet = 0
						if gf == 4096 && count > 2560 {
							wantPos = 11
						}
					} else {
						wantField = 0
						if gf != 4096 && int16(version) >= 3 {
							wantField = 0x5c
						}
					}
				}
				input.Write([]byte{0xde, 0xad, 0xbe, 0xef})
				ret, got, pos := playerFileSection(t, "sub_41AC30", input.Bytes(), uint32(uintptr(unsafe.Pointer(u))), 0)
				if ret != wantRet || pos != int64(wantPos) || !bytes.Equal(got, input.Bytes()) || p.GoldVal != wantGold || *field != wantField || u.InvFirstItem != nil {
					t.Fatal("inventory read", version, gf, which, ret, wantRet, pos, wantPos, p.GoldVal, wantGold, *field, wantField)
				}
				state := o.state()
				loaded := 0
				for _, node := range state.Nodes {
					if bytes.Equal(node.Data, []byte{113}) {
						loaded++
					}
				}
				if loaded != int(wantRet) {
					t.Fatal("inventory loaded report", version, gf, which, state)
				}
				rows = append(rows, map[string]any{"version": version, "flags": uint32(gf), "case": which, "return": ret, "position": pos, "gold": p.GoldVal, "field": *field, "queue": state})
			}
		}
	}
	spellbookCapture(t, "player-files-inventory-read-gates", rows, "")
}
