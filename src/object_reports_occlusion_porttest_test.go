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

func TestObjectReportsOcclusion(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	s.PortTestAIEmptyMap()
	t.Cleanup(s.Map.Free)
	walls, unchanged, freeWalls := s.PortTestPathWalls()
	t.Cleanup(freeWalls)
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	node, fn := alloc.New(server.MinimapItem{})
	t.Cleanup(fn)
	u, fu := alloc.New(server.Object{})
	t.Cleanup(fu)
	pl := units[0].UpdateDataPlayer().Player
	*(*types.Pointf)(unsafe.Add(unsafe.Pointer(pl), 3632)) = types.Pointf{100, 100}
	var rows []struct {
		Name         string
		Return       int
		Known, Dirty uint32
		Update       []byte
		Reliable     []legacy.PortTestShopPacketResult
	}
	defer func() {
		spellbookCapture(t, "object-reports-occlusion", rows, "cd755d3090a9e7f650572fc6b96a2664090e75c350b1b84a8785814edcd91599")
	}()
	for _, mode := range []int{0, 1, 2} {
		for _, cls := range []uint32{2, 0x100000, 0x200001, 0x20100000} {
			for _, known := range []bool{false, true} {
				for _, tracked := range []bool{false, true} {
					name := fmt.Sprintf("wall%d/class%x/known%t/tracked%t", mode, cls, known, tracked)
					t.Run(name, func(t *testing.T) {
						walls(mode)
						s.Walls.DefByInd(0).Flags32 |= 1 // Report rays (flags 69) require this visibility bit.
						reset()
						s.NetList.ResetAll()
						*u = server.Object{TypeInd: 200, NetCode: 123, ObjClass: object.Class(cls), PosVec: types.Pointf{200, 100}, Field38: 2}
						if known {
							u.Field37 = 2
						}
						pl.Field4580 = nil
						if tracked {
							*node = server.MinimapItem{Field4: u}
							node.Field8 = node
							node.Field12 = node
							pl.Field4580 = node
						}
						rv := legacy.PortTestObjectReports(6, &units[0], u, 1, 0, 0, nil)
						got := s.NetList.CopyPacketsA(1, netlist.Kind2)
						rr := snapshot()
						blocked := mode == 1 && !tracked && cls&0x20000000 == 0
						if blocked {
							if len(got) != 0 {
								t.Fatal("occluded update")
							}
							want := 0
							if known {
								want = 1
							}
							if len(rr) != want {
								t.Fatalf("occlusion reports%d want%d", len(rr), want)
							}
							if known {
								op := byte(51)
								if cls&6 != 0 {
									op = 50
								}
								if len(rr[0].Data) != 3 || rr[0].Data[0] != op || u.Field37 != 0 || u.Field38 != 2 {
									t.Fatalf("lost sight state %+v", rr)
								}
							}
						} else if len(rr) != 0 {
							t.Fatal("unexpected lost sight")
						}
						if !unchanged() {
							t.Fatal("wall owner mutated")
						}
						rows = append(rows, struct {
							Name         string
							Return       int
							Known, Dirty uint32
							Update       []byte
							Reliable     []legacy.PortTestShopPacketResult
						}{name, rv, u.Field37, u.Field38, got, rr})
					})
				}
			}
		}
	}
}
