//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestPrefabRuntimeReplaceSecretWall(t *testing.T) {
	s := newObjectXferOwner(t)
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	prefabRuntimePaths(t, words)
	secretHead, restoreSecrets := legacy.PortTestPrefabSecretList()
	defer restoreSecrets()
	live := legacy.PortTestPrefabAllocationBalance(func() {
		dst := s.Walls.CreateAtGrid(image.Pt(52, 52))
		dst.Flags4 = 4
		dst.Data = legacy.PortTestPrefabRawAllocation(32)
		data := (*[8]uint32)(dst.Data)
		data[1], data[2], data[3] = 52, 52, uint32(uintptr(dst.C()))
		*secretHead = uint32(uintptr(dst.Data))
		legacy.PortTestPrefabCall(19, [6]uint32{50, 50})
		if ret := legacy.PortTestPrefabCall(21, [6]uint32{46, 46}); ret != 1 {
			t.Fatalf("placement returned %d", ret)
		}
		if dst.Flags4&4 != 0 || dst.Data != nil || *secretHead != 0 {
			t.Error("plain replacement retained secret flag, data or list entry")
		}
		*words["placed"] = 1
		noxServer.Nox_xxx_free503F40()
		noxServer.Nox_xxx_free503F40()
		legacy.PortTestPrefabClearSecrets()
	})
	if len(live) != 0 {
		t.Errorf("unreleased allocation sizes %v", live)
	}
}
