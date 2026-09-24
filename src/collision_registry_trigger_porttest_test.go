//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestCollisionRegistryTrigger(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u, v := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	savedU, savedV := *u, *v
	data := collisionCoreGuarded(t, o, 56)
	t.Cleanup(func() { *u = savedU; *v = savedV })
	reject, accept, restore := o.s.NoxScriptVM.PortTestWorldMotionPredicates()
	t.Cleanup(restore)
	original := legacy.GetServer
	trace := &mapPolygonScriptTrace{mapPolygonScriptBase: original().NoxScriptC()}
	proxy := &mapPolygonServerTrace{Server: original(), script: trace}
	legacy.GetServer = func() legacy.Server { return proxy }
	t.Cleanup(func() { legacy.GetServer = original })
	ids := collisionCoreIDs(u, v)
	type condition struct {
		Name                      string
		Mass                      float32
		Class, Allow, Deny        uint32
		Team, AllowTeam, DenyTeam uint8
		Eligible                  bool
	}
	conditions := []condition{
		{"plain", 2, 4, 0, 0, 1, 0, 0, true}, {"zero-mass", 0, 4, 0, 0, 1, 0, 0, false}, {"negative-mass", -1, 4, 0, 0, 1, 0, 0, false}, {"nan-mass", float32(math.NaN()), 4, 0, 0, 1, 0, 0, false}, {"tiny-mass", 0.00001, 4, 0, 0, 1, 0, 0, true},
		{"allowed-class", 2, 4, 4, 0, 1, 0, 0, true}, {"wrong-class", 2, 4, 2, 0, 1, 0, 0, false}, {"denied-class", 2, 4, 0, 4, 1, 0, 0, false}, {"deny-wins", 2, 4, 4, 4, 1, 0, 0, false}, {"unrelated-deny", 2, 4, 0, 2, 1, 0, 0, true},
		{"allowed-team", 2, 4, 0, 0, 1, 1, 0, true}, {"wrong-team", 2, 4, 0, 0, 2, 1, 0, false}, {"denied-team", 2, 4, 0, 0, 1, 0, 1, false}, {"team-deny-wins", 2, 4, 0, 0, 1, 1, 1, false}, {"byte-team", 2, 4, 0, 0, 255, 255, 0, false}, {"last-positive-team", 2, 4, 0, 0, 127, 127, 0, true}, {"first-negative-team", 2, 4, 0, 0, 128, 128, 0, false}, {"byte-team-denied", 2, 4, 0, 0, 255, 0, 255, false},
	}
	type row struct {
		Name                   string
		Powered, Disabled, Nil bool
		Script                 int32
		Words                  [14]uint32
		Events                 []int
	}
	var rows []row
	for _, sp := range conditions {
		for _, powered := range []bool{false, true} {
			for _, disabled := range []bool{false, true} {
				for _, nilTarget := range []bool{false, true} {
					for _, script := range []int32{-1, reject, accept} {
						*u = savedU
						*v = savedV
						u.UpdateData = data
						u.NetCode = 1001
						v.NetCode = 1002
						u.ObjFlags = 0
						if powered {
							u.ObjFlags = 0x1000000
						}
						words := (*[14]uint32)(data)
						*words = [14]uint32{}
						words[0] = 0xabcdef80
						if disabled {
							words[2] = 5
						}
						cb := (*server.ScriptCallback)(unsafe.Add(data, 12))
						cb.Func = script
						words[11] = sp.Allow
						words[12] = sp.Deny
						*(*uint8)(unsafe.Add(data, 52)) = sp.AllowTeam
						*(*uint8)(unsafe.Add(data, 53)) = sp.DenyTeam
						v.Mass = sp.Mass
						v.ObjClass = object.Class(sp.Class)
						v.TeamVal.ID = server.TeamID(sp.Team)
						target := v
						if nilTarget {
							target = nil
						}
						trace.events = nil
						legacy.PortTestRegisteredCollision(u, target, nil, "TriggerCollide")
						admitted := powered && !disabled && !nilTarget && sp.Eligible
						called := admitted && script != -1
						fired := admitted && script != reject
						if (len(trace.events) == 1) != called || len(trace.events) > 1 || called && trace.events[0] != 1 {
							t.Fatal("trigger script admission", sp.Name, powered, disabled, nilTarget, script, trace.events)
						}
						if called && (o.s.NoxScriptVM.Caller() != v || o.s.NoxScriptVM.Trigger() != u) {
							t.Fatal("trigger script caller and source")
						}
						want := uint32(0xabcdef80)
						targetID := uint32(0)
						if fired {
							want |= 1
							targetID = 1002
						}
						got := *words
						got[1] = collisionCoreID(t, ids, got[1])
						if got[0] != want || got[1] != targetID {
							t.Fatal("trigger class/team/mass/script contract", sp.Name, powered, disabled, nilTarget, script, got[0], got[1], want, targetID)
						}
						rows = append(rows, row{sp.Name, powered, disabled, nilTarget, script, got, append([]int(nil), trace.events...)})
					}
				}
			}
		}
	}
	collisionRegistryCapture(t, "collision-registry-trigger", rows)
}
