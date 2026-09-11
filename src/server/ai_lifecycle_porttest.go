//go:build porttest

package server

import (
	"strings"

	"github.com/opennox/libs/balance"
	"github.com/opennox/libs/object"
)

// PortTestLifecycleTypeIDs identifies the six synthetic types installed by
// PortTestLifecycleTypes. The IDs intentionally follow the combat fixture's
// projectile at type 1 on a fresh server.
type PortTestLifecycleTypeIDs struct {
	Zombie, VileZombie, ReleasedSoul       int
	ScorchSmall, ScorchMedium, ScorchLarge int
}

// PortTestLifecycleTypes installs the minimal real ObjectTypes used by the
// lifecycle fixture. It requires PortTestCombatProjectileType to have already
// installed its type at index 1 on a fresh server. All six are simple objects:
// creation therefore follows the normal C-backed allocator path without asset
// init, update data, or type callbacks. Cleanup restores the exact prior type
// tables and their cached Zombie/VileZombie IDs; it does not free the existing
// allocator or any object owned by its caller.
func (s *Server) PortTestLifecycleTypes() (PortTestLifecycleTypeIDs, func()) {
	if len(s.Types.byInd) != 2 || s.Types.byInd[1] == nil {
		panic("PortTestLifecycleTypes requires PortTestCombatProjectileType on a fresh server")
	}
	old := s.Types
	byInd := append([]*ObjectType(nil), old.byInd...)
	byID := make(map[string]*ObjectType, len(old.byID)+6)
	for k, v := range old.byID {
		byID[k] = v
	}
	s.Types.byInd = append(byInd, make([]*ObjectType, 6)...)
	s.Types.byID = byID
	s.Types.fast.zombie = 0
	s.Types.fast.zombieVile = 0

	ids := PortTestLifecycleTypeIDs{
		Zombie: 2, VileZombie: 3, ReleasedSoul: 4,
		ScorchSmall: 5, ScorchMedium: 6, ScorchLarge: 7,
	}
	for _, v := range []struct {
		id  string
		ind int
	}{
		{"Zombie", ids.Zombie},
		{"VileZombie", ids.VileZombie},
		{"ReleasedSoul", ids.ReleasedSoul},
		// These are the actual C scorch-init names in GAME4_3.c/blob_587000.
		{"ScorchMarkFloorSmallA", ids.ScorchSmall},
		{"ScorchMarkFloorMediumA", ids.ScorchMedium},
		{"ScorchMarkFloorLargeB", ids.ScorchLarge},
	} {
		name := strings.ToLower(v.id)
		t := &ObjectType{
			s: &s.Types, ind: uint16(v.ind), ind2: uint16(v.ind), id: name,
			class: object.ClassSimple, allowed: true,
		}
		s.Types.byInd[v.ind] = t
		s.Types.byID[name] = t
	}
	return ids, func() { s.Types = old }
}

// PortTestZombieDeadDuration overlays the live balance file with exactly the
// two values queried by the C lifecycle owner. The overlay is a real
// balance.File and falls back to the prior file for every other key. Cleanup
// restores the original balance file.
func (s *Server) PortTestZombieDeadDuration(min, max float64) func() {
	old := s.Balance.file
	s.Balance.file = &balance.File{
		Global: balance.Config{
			"zombiedeadduration": balance.Array{min, max},
		},
		Tags:   make(map[balance.Tag]balance.Config),
		Parent: old,
	}
	return func() { s.Balance.file = old }
}
