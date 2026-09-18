//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestServerOrchestrationRestoreIntegration(t *testing.T) {
	o := newMatchRosterOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	_, restore := legacy.PortTestServerOrchestrationGlobals()
	t.Cleanup(restore)
	head, restorePending := legacy.PortTestMonsterPendingOwner()
	t.Cleanup(restorePending)
	legacy.PortTestMonsterPending("init", 0, 0)
	_, restoreCache := legacy.PortTestMonsterCacheInit()
	t.Cleanup(restoreCache)
	serverConfigOwnBytes(t, 0x5D4594, 2386620, 208)
	for _, off := range []uintptr{1563128, 1563132} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	serverConfigOwnBytes(t, 0x587000, 255604, 16)
	clear(unsafe.Slice(memmap.PtrUint8(0x587000, 255604), 16))
	*memmap.PtrUint32(0x587000, 255604) = 1
	*memmap.PtrUint32(0x587000, 255608) = 1
	t.Cleanup(o.s.PortTestRewardTypes([]string{"SaveGameLocation", "Glyph", "RestoreIntegration"}, nil, true, 0, 0))
	oldList, oldDeleted := o.s.Objs.List, o.s.Objs.DeletedList
	oldTicks, oldFrame, oldSolo := nox_gameTicks_371764, nox_gameFrame_371772, gameIsSwitchToSolo
	t.Cleanup(func() {
		o.s.Objs.List = oldList
		o.s.Objs.DeletedList = oldDeleted
		nox_gameTicks_371764 = oldTicks
		nox_gameFrame_371772 = oldFrame
		gameIsSwitchToSolo = oldSolo
	})
	host := &o.units[2]
	oldHost := *host
	t.Cleanup(func() { *host = oldHost })
	type row struct {
		Arg                         int32
		Protected                   bool
		Position                    types.Pointf
		HostID                      int
		MarkerIDs                   [2]int
		Deleted                     [2]bool
		AIResolved, PendingResolved bool
		ActivatorIDs                [][2]uint32
		ActivatorObjects            [][2]int
		Tracked                     [4]bool
		Packets                     [][]byte
	}
	var rows []row
	for _, arg := range []int32{0, 1, 2} {
		for _, protected := range []bool{false, true} {
			t.Run(fmt.Sprintf("arg%d/protected%t", arg, protected), func(t *testing.T) {
				o.reset()
				legacy.PortTestMonsterCache("reset", nil, 0)
				legacy.PortTestMonsterPending("clear", 0, 0)
				var units []*server.Object
				makeObj := func(name string) *server.Object {
					u := o.s.NewObjectByTypeID(name)
					if u == nil {
						t.Fatal(name)
					}
					u.ObjFlags = 0
					u.ScriptIDVal = 100 + len(units)
					u.NetCode = uint32(200 + len(units))
					units = append(units, u)
					return u
				}
				a, b, target, child, monster := makeObj("SaveGameLocation"), makeObj("SaveGameLocation"), makeObj("RestoreIntegration"), makeObj("RestoreIntegration"), makeObj("RestoreIntegration")
				a.PosVec = types.Pointf{X: 123, Y: 145}
				b.PosVec = types.Pointf{X: 201, Y: 233}
				var owned [4]*server.Object
				for i := range owned {
					u := makeObj("RestoreIntegration")
					owned[i] = u
					u.ObjClass = 2
					u.UpdateData = o.record(t, 2560)
					u.HealthData = (*server.HealthData)(o.record(t, int(unsafe.Sizeof(server.HealthData{}))))
					u.HealthData.Cur = 11
					u.HealthData.Max = 79
					if i&1 != 0 {
						*(*byte)(unsafe.Add(u.UpdateData, 1440)) = 0x80
					}
					if i&2 != 0 {
						u.ObjSubClass = 0x80
					}
					u.ObjOwner = host
					if i > 0 {
						owned[i-1].Field128 = u
					}
				}
				monster.ObjClass = 2
				monster.UpdateData = o.record(t, 2560)
				objectXferSetWord(monster.UpdateData, 556, uint32(target.ScriptIDVal))
				if protected {
					monster.ObjFlags = 0x80000000
				}
				a.ObjNext = b
				b.ObjNext = target
				target.ObjNext = child
				child.ObjNext = monster
				monster.ObjNext = nil
				o.s.Objs.List = a
				o.s.Objs.DeletedList = nil
				host.PosVec = types.Pointf{X: 70, Y: 80}
				host.ScriptIDVal = 77
				host.Field129 = owned[0]
				legacy.PortTestMonsterPending("add", int32(target.ScriptIDVal), int32(child.ScriptIDVal))
				observe, restoreActivators := o.s.PortTestOrchestrationActivators([][2]uint32{{uint32(target.ScriptIDVal), uint32(child.ScriptIDVal)}, {999, 0}})
				defer restoreActivators()
				defer func() {
					o.s.Objs.List = nil
					o.s.Objs.DeletedList = nil
					host.Field129 = nil
					for _, u := range units {
						legacy.PortTestPlayerStateMinimap("unmark", u, 0, ^uint32(0))
						u.InvFirstItem = nil
						o.s.Objs.FreeObject(u)
					}
				}()
				legacy.PortTestServerOrchestration("restore", nil, arg)
				r := row{Arg: arg, Protected: protected, Position: host.PosVec, HostID: host.ScriptIDVal, MarkerIDs: [2]int{a.ScriptIDVal, b.ScriptIDVal}, Deleted: [2]bool{a.ObjFlags.Has(object.FlagDestroyed), b.ObjFlags.Has(object.FlagDestroyed)}, AIResolved: objectXferGetWord(monster.UpdateData, 556) == uint32(uintptr(target.CObj())), PendingResolved: child.ObjOwner == target, Packets: visibilityEffectsPackets(o.s)}
				if arg == 1 {
					if r.Position != b.PosVec || r.HostID != 101 || r.MarkerIDs != [2]int{} || r.Deleted != [2]bool{true, true} || !r.PendingResolved || *head != 0 {
						t.Fatal("saved world restore", r)
					}
				} else if r.Position != (types.Pointf{X: 70, Y: 80}) || r.HostID != 77 || r.MarkerIDs != [2]int{100, 101} || r.Deleted != [2]bool{} || r.PendingResolved || *head == 0 {
					t.Fatal("restore bypass", r)
				}
				if r.AIResolved != (arg == 1 && !protected) {
					t.Fatal("AI postload dispatch", r.AIResolved)
				}
				refs, args := observe()
				r.ActivatorIDs = refs
				for _, a := range args {
					pair := [2]int{}
					if a.Trigger != nil {
						pair[0] = a.Trigger.ScriptIDVal
					}
					if a.Caller != nil {
						pair[1] = a.Caller.ScriptIDVal
					}
					r.ActivatorObjects = append(r.ActivatorObjects, pair)
					if a.Callback != 7 || a.Arg != 9 {
						t.Fatal("timer metadata")
					}
				}
				if arg == 1 {
					if fmt.Sprint(refs) != "[[0 0] [0 0]]" || fmt.Sprint(r.ActivatorObjects) != "[[102 103] [0 0]]" {
						t.Fatal("activator references", r)
					}
				} else if fmt.Sprint(refs) != "[[102 103] [999 0]]" || fmt.Sprint(r.ActivatorObjects) != "[[0 0] [0 0]]" {
					t.Fatal("activator bypass", r)
				}
				for i, u := range owned {
					r.Tracked[i] = legacy.PortTestPlayerStateMinimap("tracks", u, 31, 0) != 0
					if r.Tracked[i] != (i != 0) {
						t.Fatal("owned creature visibility", r.Tracked)
					}
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "server-orchestration-restore-integration", rows, "3425da6b38be9af303273132e17ec9a8a60444b6c2f4e96ead8abc3a0fde070a")
}
