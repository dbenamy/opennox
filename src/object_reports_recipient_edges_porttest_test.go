//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestObjectReportsRecipientEdges(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	s.PortTestAIEmptyMap()
	t.Cleanup(s.Map.Free)
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	u, fu := alloc.New(server.Object{})
	t.Cleanup(fu)
	parent, fp := alloc.New(server.Object{})
	t.Cleanup(fp)
	var rows []struct {
		Name                 string
		Return               int
		Known, Dirty, Friend uint32
		Update               []byte
		Reports              []legacy.PortTestShopPacketResult
	}
	defer func() {
		spellbookCapture(t, "object-reports-recipient-edges", rows, "bbf48241c3e8b8b6c921587c01c79bec057e7fa8b90fc6a32cf501fd2e0e0bca")
	}()
	for i, slot := range []int{1, 7, 31} {
		for _, owner := range []string{"none", "viewer", "chain", "other"} {
			for _, reveal := range []bool{false, true} {
				name := fmt.Sprintf("slot%d/owner%s/reveal%t", slot, owner, reveal)
				t.Run(name, func(t *testing.T) {
					flags := uint32(0x200000)
					if reveal {
						flags |= 512
					}
					t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
					a := &units[i]
					*(*types.Pointf)(unsafe.Add(unsafe.Pointer(a.UpdateDataPlayer().Player), 3632)) = types.Pointf{100, 100}
					reset()
					s.NetList.ResetAll()
					bit := uint32(1) << uint(slot)
					*u = server.Object{TypeInd: 200, NetCode: 123, ObjClass: object.ClassSimple, ObjFlags: 0x800, PosVec: types.Pointf{100, 100}, Field38: 0xffffffff, Field35: 0xffffffff, Field36: bit}
					*parent = server.Object{}
					switch owner {
					case "viewer":
						u.ObjOwner = a
					case "chain":
						u.ObjOwner = parent
						parent.ObjOwner = a
					case "other":
						u.ObjOwner = &units[(i+1)%3]
					}
					rv := legacy.PortTestObjectReports(6, a, u, slot, 0, 0, nil)
					got := s.NetList.CopyPacketsA(a.UpdateDataPlayer().Player.PlayerIndex(), netlist.Kind2)
					rr := snapshot()
					canSee := reveal || owner == "viewer" || owner == "chain"
					known, dirty := uint32(0), uint32(0xffffffff)
					size := 0
					if canSee {
						known = bit
						dirty &^= bit
						size = 9
					}
					if u.Field37 != known || u.Field38 != dirty || u.Field35 != ^bit || len(got) != size || len(rr) != 1 || rr[0].Recipient != byte(slot) {
						t.Fatalf("recipient state known%x dirty%x friend%x bytes%x reports%+v", u.Field37, u.Field38, u.Field35, got, rr)
					}
					rows = append(rows, struct {
						Name                 string
						Return               int
						Known, Dirty, Friend uint32
						Update               []byte
						Reports              []legacy.PortTestShopPacketResult
					}{name, rv, u.Field37, u.Field38, u.Field35, got, rr})
				})
			}
		}
	}
}
