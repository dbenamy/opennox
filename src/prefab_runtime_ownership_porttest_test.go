//go:build porttest

package opennox

import (
	"fmt"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/wall"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestPrefabRuntimeSpecialWallTransfer(t *testing.T) {
	for _, flags := range []byte{4, 8, 12} {
		for _, partial := range []bool{false, true} {
			t.Run(fmt.Sprintf("flags%x/partial%t", flags, partial), func(t *testing.T) {
				s := newObjectXferOwner(t)
				words, restore := legacy.PortTestPrefabRuntimeGlobals()
				defer restore()
				prefabRuntimePaths(t, words)
				secretHead, restoreSecrets := legacy.PortTestPrefabSecretList()
				defer restoreSecrets()
				live := legacy.PortTestPrefabAllocationBalance(func() {
					if partial {
						legacy.PortTestPrefabCall(19, [6]uint32{255, 255})
					}
					raw := uint32(legacy.PortTestPrefabCall(19, [6]uint32{50, 50}))
					node := unsafe.Pointer(uintptr(raw))
					src := *(**server.Wall)(node)
					src.Flags4 = wall.Flags(flags)
					src.Field10 = 17
					var data *[8]uint32
					if flags&4 != 0 {
						src.Data = legacy.PortTestPrefabRawAllocation(32)
						data = (*[8]uint32)(src.Data)
						data[1], data[2], data[3] = 50, 50, uint32(uintptr(src.C()))
					}
					ret := legacy.PortTestPrefabCall(21, [6]uint32{46, 46})
					want := uint64(1)
					if partial {
						want = 0
					}
					if ret != want {
						t.Fatalf("placement returned%d want%d", ret, want)
					}
					dst := s.Walls.GetWallAtGrid(image.Pt(52, 52))
					if dst == nil {
						t.Fatal("no world wall")
					}
					if dst.Field10 != 17 {
						t.Errorf("wall identity was not transferred: %d", dst.Field10)
					}
					if flags&4 != 0 {
						if dst.Data != unsafe.Pointer(data) || *secretHead != uint32(uintptr(unsafe.Pointer(data))) {
							t.Error("secret data not linked to world owner")
						}
						if data[3] != uint32(uintptr(dst.C())) || data[1] != 52 || data[2] != 52 {
							t.Errorf("secret data retained cache reference/coordinates: %x", data[:4])
						}
						if src.Data != nil {
							t.Error("cache still owns transferred secret data")
						}
					}
					if flags&8 != 0 {
						p := s.Walls.FirstBreakable()
						if p == nil || p.Wall != dst {
							t.Error("breakable list does not reference world wall")
						}
					}
					if !partial {
						*words["placed"] = 1
					}
					noxServer.Nox_xxx_free503F40()
					noxServer.Nox_xxx_free503F40()
					if flags&4 != 0 && dst.Data != unsafe.Pointer(data) {
						t.Error("cache disposal changed world secret data")
					}
					legacy.PortTestPrefabClearSecrets()
					dst.Data = nil
					s.Walls.ClearBreakable()
				})
				if len(live) != 0 {
					t.Errorf("unreleased allocation sizes %v", live)
				}
			})
		}
	}
}
func TestPrefabRuntimeTileWallDisposal(t *testing.T) {
	for _, placed := range []uint32{0, 1} {
		t.Run(fmt.Sprintf("placed%d", placed), func(t *testing.T) {
			s := newObjectXferOwner(t)
			_ = s
			words, restore := legacy.PortTestPrefabRuntimeGlobals()
			defer restore()
			prefabRuntimePaths(t, words)
			live := legacy.PortTestPrefabAllocationBalance(func() {
				legacy.PortTestPrefabCacheNode(0)
				node := legacy.PortTestPrefabCacheNode(1)
				if placed == 0 {
					w := *(**server.Wall)(node)
					w.Flags4 = 4
					w.Data = legacy.PortTestPrefabRawAllocation(32)
				}
				*words["placed"] = placed
				noxServer.Nox_xxx_free503F40()
				noxServer.Nox_xxx_free503F40()
			})
			if len(live) != 0 {
				t.Errorf("unreleased allocation sizes %v", live)
			}
		})
	}
}
