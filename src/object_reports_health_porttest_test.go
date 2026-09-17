//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestObjectReportsHealthHistory(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	resetReliable, snapshotReliable, freeReliable := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(freeReliable)
	monster := newCreatureXferObject(t, s, "Monster")
	monster.UpdateDataMonster().AIStackInd = -1
	hp, free := alloc.New(server.HealthData{})
	t.Cleanup(free)
	var rows []struct {
		Name             string
		Return           int
		History          uint16
		Reliable, Update []byte
	}
	defer func() {
		spellbookCapture(t, "object-reports-health", rows, "d744a88ebb868b1017cdc1d2640bb4c0db2f50f402535e10e3e55fb12465c0f9")
	}()
	for _, kind := range []string{"generator", "monster", "player"} {
		for _, enabled := range []bool{false, true} {
			for _, hasHealth := range []bool{false, true} {
				for _, elapsed := range []uint32{0, 2, 3, 0xffffffff} {
					for _, cur := range []uint16{0, 1, 32767, 32768, 65535} {
						for _, prev := range []uint16{0, 1, 32767, 32768, 65535} {
							name := fmt.Sprintf("%s/enabled%t/health%t/elapsed%d/current%d/prev%d", kind, enabled, hasHealth, elapsed, cur, prev)
							t.Run(name, func(t *testing.T) {
								s.NetList.ResetAll()
								resetReliable()
								s.SetFrame(1)
								u := &units[1]
								op, off := 2, 96
								u.ObjClass = object.ClassMonsterGenerator
								u.NetCode = 123
								if kind == "monster" {
									u = monster
									op, off = 3, 412
								} else if kind == "player" {
									u.ObjClass = object.ClassPlayer
									op, off = 5, 12
								}
								u.NetCode = 123
								u.HealthData = nil
								if hasHealth {
									hp.Cur = cur
									u.HealthData = hp
								}
								u.Frame134 = 1 - elapsed
								history := (*uint16)(unsafe.Add(u.UpdateData, off+2))
								*history = prev
								flag := 0
								if enabled {
									flag = 1
								}
								rv := legacy.PortTestObjectReports(op, &units[0], u, 1, flag, flag, nil)
								ordinary := s.NetList.CopyPacketsA(1, netlist.Kind1)
								update := s.NetList.CopyPacketsA(1, netlist.Kind2)
								if kind == "player" && !enabled {
									update = ordinary
									ordinary = nil
								}
								wantPrev := prev
								var want []byte
								if hasHealth && elapsed > 2 && (kind == "generator" || enabled) {
									wantPrev = cur
									delta := int16(cur - prev)
									if delta < 0 {
										want = []byte{66, 123, 0, 0, 0}
										binary.LittleEndian.PutUint16(want[3:], uint16(delta))
									}
								}
								reliable := snapshotReliable()
								var delivered []byte
								if len(reliable) > 1 {
									t.Fatal("extra reliable report")
								}
								if len(reliable) == 1 {
									if reliable[0].Recipient != 1 {
										t.Fatal("wrong health recipient")
									}
									delivered = reliable[0].Data
								}
								if rv != 1 || *history != wantPrev || !bytes.Equal(delivered, want) || len(ordinary) != 0 {
									t.Fatalf("return%d history%d want%d message%x want%x", rv, *history, wantPrev, delivered, want)
								}
								size := 9
								if kind == "monster" {
									size = 11
								} else if kind == "player" {
									size = 12
								}
								if len(update) != size {
									t.Fatalf("update size %d want%d: %x", len(update), size, update)
								}
								rows = append(rows, struct {
									Name             string
									Return           int
									History          uint16
									Reliable, Update []byte
								}{name, rv, *history, delivered, update})
								u.HealthData = nil
							})
						}
					}
				}
			}
		}
	}
}
