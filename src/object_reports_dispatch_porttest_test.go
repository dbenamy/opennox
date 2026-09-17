//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestObjectReportsRecipientMasks(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	s.PortTestAIEmptyMap()
	t.Cleanup(s.Map.Free)
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	u, fu := alloc.New(server.Object{})
	t.Cleanup(fu)
	var rows []struct {
		Name                 string
		Return               int
		Known, Dirty, Friend uint32
		Update               []byte
		Reliable             []legacy.PortTestShopPacketResult
	}
	defer func() {
		spellbookCapture(t, "object-reports-recipient-masks", rows, "1e0d69068cf86f0e00a94a1b69cb9380c8b7c3f0eca7ab5a9337d12ee6a51e2c")
	}()
	for _, cls := range []uint32{0, 0x100000, 0x200001, 0x40000000} {
		for _, flags := range []uint32{0, 0x20, 0x800} {
			for _, known := range []bool{false, true} {
				for _, dirty := range []bool{false, true} {
					for _, stationary := range []bool{false, true} {
						for _, friend := range []bool{false, true} {
							name := fmt.Sprintf("class%x/flags%x/known%t/dirty%t/stationary%t/friend%t", cls, flags, known, dirty, stationary, friend)
							t.Run(name, func(t *testing.T) {
								a := &units[0]
								*(*types.Pointf)(unsafe.Add(unsafe.Pointer(a.UpdateDataPlayer().Player), 3632)) = types.Pointf{100, 100}
								s.NetList.ResetAll()
								reset()
								*u = server.Object{NetCode: 123, TypeInd: 200, ObjClass: object.Class(cls), ObjFlags: object.Flags(flags), PosVec: types.Pointf{100, 100}}
								const bit = uint32(2)
								if known {
									u.Field37 = bit
								}
								if dirty {
									u.Field38 = bit
								}
								if stationary {
									u.Field5 = 32
								}
								if friend {
									u.Field35 = bit
									u.Field36 = bit
								}
								rv := legacy.PortTestObjectReports(6, a, u, 1, 0, 0, nil)
								update := s.NetList.CopyPacketsA(1, netlist.Kind2)
								rr := snapshot()
								eligible := flags&0x20 == 0 && cls&0x40000000 == 0
								wantFriend := uint32(0)
								if friend && !eligible {
									wantFriend = bit
								}
								if u.Field35 != wantFriend {
									t.Fatal("friend dirty bit")
								}
								wantReports := 0
								if friend && eligible {
									wantReports = 1
								}
								if len(rr) != wantReports {
									t.Fatalf("reliable reports%d want%d", len(rr), wantReports)
								}
								if len(rr) == 1 && (rr[0].Recipient != 1 || len(rr[0].Data) != 3 || rr[0].Data[0] != 52) {
									t.Fatal("friend report")
								}
								shouldSend := eligible && (dirty || cls&1 != 0) && flags&0x800 == 0 && !(known && stationary) && cls != 0
								if shouldSend {
									wantLen := 9
									if cls&0x200000 != 0 {
										wantLen = 11
									}
									if len(update) != wantLen || u.Field37&bit == 0 || u.Field38&bit != 0 {
										t.Fatalf("send state packet%x known%x dirty%x", update, u.Field37, u.Field38)
									}
								} else if len(update) != 0 {
									t.Fatalf("unexpected update%x", update)
								}
								rows = append(rows, struct {
									Name                 string
									Return               int
									Known, Dirty, Friend uint32
									Update               []byte
									Reliable             []legacy.PortTestShopPacketResult
								}{name, rv, u.Field37, u.Field38, u.Field35, update, rr})
							})
						}
					}
				}
			}
		}
	}
}

func TestObjectReportsDirtyPropagation(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	s.PortTestAIEmptyMap()
	t.Cleanup(s.Map.Free)
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	u, fu := alloc.New(server.Object{})
	t.Cleanup(fu)
	*(*types.Pointf)(unsafe.Add(unsafe.Pointer(units[0].UpdateDataPlayer().Player), 3632)) = types.Pointf{100, 100}
	var rows []struct {
		Name                string
		Return              int
		Flags, Known, Dirty uint32
		Reports             []legacy.PortTestShopPacketResult
	}
	defer func() {
		spellbookCapture(t, "object-reports-dirty-propagation", rows, "e7f8bbc7be839bc587968f2ef3da1e55b91ea20714df49bfc1aa209fbfed3826")
	}()
	for _, mask := range []uint32{0, 1, 4, 8, 16, 0x40, 0x80, 0xdd} {
		for _, known := range []bool{false, true} {
			for _, stationary := range []bool{false, true} {
				name := fmt.Sprintf("mask%x/known%t/stationary%t", mask, known, stationary)
				t.Run(name, func(t *testing.T) {
					s.NetList.ResetAll()
					reset()
					*u = server.Object{TypeInd: 200, NetCode: 123, ObjClass: object.ClassSimple, PosVec: types.Pointf{100, 100}, Field38: 2}
					u.Field140[1] = mask
					if known {
						u.Field37 = 2
					}
					if stationary {
						u.Field5 = 32
					}
					rv := legacy.PortTestObjectReports(6, &units[0], u, 1, 0, 0, nil)
					expected := mask
					if !known {
						expected |= (mask &^ 0xcd) << 16
					}
					if u.Field140[1] != expected {
						t.Fatalf("flags%x want%x", u.Field140[1], expected)
					}
					if !(known && stationary) && (u.Field37 != 2 || u.Field38 != 0) {
						t.Fatal("recipient bit transition")
					}
					rows = append(rows, struct {
						Name                string
						Return              int
						Flags, Known, Dirty uint32
						Reports             []legacy.PortTestShopPacketResult
					}{name, rv, u.Field140[1], u.Field37, u.Field38, snapshot()})
				})
			}
		}
	}
}
