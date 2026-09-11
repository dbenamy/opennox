//go:build porttest

package server

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"unsafe"

	"github.com/opennox/libs/balance"
	"github.com/opennox/libs/object"
)

// PortTestAICallbackTypes appends the simple allocator-backed types used by
// death/debris and generated-equipment callbacks. Existing main-fixture IDs
// 0..11, notably SmallToxicCloud at 11, remain untouched.
type PortTestAICallbackTypeIDs struct {
	Skull, DebrisA, DebrisB                                    int
	Sword, WoodenShield, SteelShield, StaffWooden, Bow, Quiver int
	OgreAxe, FanChakram                                        int
}

// PortTestAICallbackTypes requires PortTestMainTypes (12 entries). configure
// removes an explicitly disabled synthetic name from both real type lookup
// tables, so NewObjectByTypeID and NewObjectByTypeInd return nil. Omitted names
// remain enabled. Each configure restores prior synthetic entries before applying
// its map. It intentionally does not attach update data: the legacy fixture supplies
// that only for SmallToxicCloud where its callback requires MonsterUpdateData.
func (s *Server) PortTestAICallbackTypes() (ids PortTestAICallbackTypeIDs, configure func(map[string]bool), restore func()) {
	if len(s.Types.byInd) != 12 || s.Types.byInd[10] == nil || s.Types.byInd[11] == nil {
		panic("PortTestAICallbackTypes requires PortTestMainTypes on a fresh server")
	}
	// Debris callbacks can create more than the combat fixture's eight objects.
	s.Objs.FreeObjects()
	if !s.Objs.Init(64) {
		panic("callback allocator initialization failed")
	}
	old := s.Types
	byInd := append([]*ObjectType(nil), old.byInd...)
	byID := make(map[string]*ObjectType, len(old.byID)+11)
	for k, v := range old.byID {
		byID[k] = v
	}
	s.Types.byInd = append(byInd, make([]*ObjectType, 11)...)
	s.Types.byID = byID
	cloud := *old.byInd[11]
	cloudData, freeCloudData := alloc.New(uint32(0))
	*cloudData = 0
	cloud.UpdateDataSize = 4
	cloud.UpdateData = unsafe.Pointer(cloudData)
	s.Types.byInd[11] = &cloud
	s.Types.byID[cloud.id] = &cloud
	ids = PortTestAICallbackTypeIDs{Skull: 12, DebrisA: 13, DebrisB: 14, Sword: 15, WoodenShield: 16, SteelShield: 17, StaffWooden: 18, Bow: 19, Quiver: 20, OgreAxe: 21, FanChakram: 22}
	entries := []struct {
		name string
		ind  int
	}{
		{"Skull", ids.Skull}, {"PortTestDebrisA", ids.DebrisA}, {"PortTestDebrisB", ids.DebrisB},
		{"Sword", ids.Sword}, {"WoodenShield", ids.WoodenShield}, {"SteelShield", ids.SteelShield},
		{"StaffWooden", ids.StaffWooden}, {"Bow", ids.Bow}, {"Quiver", ids.Quiver}, {"OgreAxe", ids.OgreAxe}, {"FanChakram", ids.FanChakram},
	}
	created := make(map[string]*ObjectType, len(entries))
	for _, e := range entries {
		name := strings.ToLower(e.name)
		t := &ObjectType{s: &s.Types, ind: uint16(e.ind), ind2: uint16(e.ind), id: name, class: object.ClassSimple, allowed: true, Mass: 1}
		s.Types.byInd[e.ind], s.Types.byID[name], created[name] = t, t, t
	}
	configure = func(enabled map[string]bool) {
		// Reinstall every synthetic entry before applying this case; configurations
		// are independent and do not retain a prior case's disabled names.
		for name, t := range created {
			s.Types.byInd[int(t.ind)] = t
			s.Types.byID[name] = t
		}
		for name, disabled := range enabled {
			if !disabled { // map false explicitly means unavailable.
				if t := created[name]; t != nil {
					delete(s.Types.byID, name)
					s.Types.byInd[int(t.ind)] = nil
				}
			}
		}
	}
	return ids, configure, func() { s.Types = old; freeCloudData() }
}

// PortTestSmallToxicCloudLifetime overlays exactly the lifetime queried by
// 54A270. Repeated set calls mutate one file and never chain balance parents.
func (s *Server) PortTestSmallToxicCloudLifetime() (set func(float64), restore func()) {
	old := s.Balance.file
	overlay := &balance.File{Global: balance.Config{}, Tags: make(map[balance.Tag]balance.Config), Parent: old}
	s.Balance.file = overlay
	set = func(v float64) { overlay.Global["smalltoxiccloudlifetime"] = balance.Array{v} }
	return set, func() { s.Balance.file = old }
}
