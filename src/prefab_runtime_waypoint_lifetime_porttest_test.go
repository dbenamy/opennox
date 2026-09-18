//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestPrefabRuntimeWaypointLifetime(t *testing.T) {
	_, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	old := noxServer
	core := &server.Server{}
	noxServer = &Server{Server: core}
	defer func() { noxServer = old }()
	defer noxflags.PortTestGameFlags(0)()
	node := legacy.PortTestPrefabCacheNode(2)
	if node == nil {
		t.Fatal("prefab waypoint allocation")
	}
	defer legacy.MapPrefabFreeNode(node)
	wp := *(**server.Waypoint)(node)
	if legacy.PortTestPrefabCall(23, [6]uint32{10, 20}) != 1 || core.WPs.Pending != wp {
		t.Fatal("waypoint not transferred to actual pending owner")
	}
	if wp.PosVec.X != 12 || wp.PosVec.Y != 23 {
		t.Fatal("waypoint placement offset")
	}
	core.Nox_xxx_waypoint_5799C0()
	if core.WPs.List != wp || core.WPs.Pending != nil {
		t.Fatal("waypoint not transferred to actual world owner")
	}
	// The regular game teardown must accept the allocation used by the prefab path.
	core.Nox_xxx_waypointDeleteAll_579DD0()
	if core.WPs.List != nil {
		t.Fatal("waypoint survived teardown")
	}
	// The cache wrapper is separate storage and is released by its own owner.
}

func TestPrefabRuntimeWaypointCacheDisposal(t *testing.T) {
	for _, placed := range []bool{false, true} {
		t.Run(fmt.Sprintf("placed%t", placed), func(t *testing.T) {
			words, restore := legacy.PortTestPrefabRuntimeGlobals()
			defer restore()
			old := noxServer
			core := &server.Server{}
			noxServer = &Server{Server: core}
			defer func() { noxServer = old }()
			defer noxflags.PortTestGameFlags(0)()
			path, freePath := alloc.Make([]byte{}, 2048)
			defer freePath()
			alternate, freeAlt := alloc.Make([]byte{}, 2048)
			defer freeAlt()
			*words["path"] = uint32(uintptr(unsafe.Pointer(&path[0])))
			*words["alternate"] = uint32(uintptr(unsafe.Pointer(&alternate[0])))
			live := legacy.PortTestPrefabAllocationBalance(func() {
				if legacy.PortTestPrefabCacheNode(2) == nil {
					t.Fatal("cache waypoint allocation")
				}
				if placed {
					if legacy.PortTestPrefabCall(23, [6]uint32{10, 20}) != 1 {
						t.Fatal("placement failed")
					}
					core.Nox_xxx_waypoint_5799C0()
					*words["placed"] = 1
				}
				noxServer.Nox_xxx_free503F40()
				if placed {
					if core.WPs.List == nil || core.WPs.List.Index != 17 {
						t.Fatal("cache cleanup destroyed transferred waypoint")
					}
					core.Nox_xxx_waypointDeleteAll_579DD0()
				}
			})
			if len(live) != 0 {
				t.Fatalf("unreleased allocation sizes %v", live)
			}
			if *words["waypoints"] != 0 || core.WPs.List != nil || core.WPs.Pending != nil {
				t.Fatal("stale waypoint ownership")
			}
		})
	}
}
