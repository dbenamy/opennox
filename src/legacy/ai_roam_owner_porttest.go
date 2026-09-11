//go:build porttest

package legacy

import (
	"math"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
)

type PortTestRoamOwnerSpec struct {
	Repeat                                int
	Frame, FPS, Aggression, Status, Buffs uint32
	Current, Fallback                     byte
	Register, Enemy, ExistingPath         bool
	PathMode                              byte
	X, Y, WX, WY                          uint32
}

type portTestRoamOwnerServer struct {
	portTestRandomServer
	combat           *portTestCombatState
	combatProjectile *server.Object
	pathEndpoints    [2]*server.Waypoint
	endpointCall     int
	generator        int
	mode             byte
	precheck         bool
	fallback         *server.Waypoint
	trace            []uint32
}

func (s *portTestRoamOwnerServer) Nox_xxx_creatureSetDetailedPath_50D220(u *server.Object, p *types.Pointf) {
	s.trace = append(s.trace, 1, math.Float32bits(p.X), math.Float32bits(p.Y))
	ud := u.UpdateDataMonster()
	ud.Field2 = 0
	ud.Field71 = 0
	switch s.mode {
	case 1:
		ud.Field71 = 2
	case 2:
		ud.Field2 = 1
		ud.Field67 = 1
	case 4:
		ud.Field2 = 2
		ud.Field67 = 0
		ud.Path[0] = u.PosVec
		ud.Path[1] = *p
	case 5:
		ud.Field2 = 1
		ud.Field67 = 0
		ud.Path[0] = *p
	case 6:
		ud.Field71 = 1
	case 3:
		ud.Field2 = 1
		ud.Field67 = 0
		ud.Path[0] = u.PosVec
	}
}
func (s *portTestRoamOwnerServer) Sub_50CB20(u *server.Object, p *types.Pointf) *server.Waypoint {
	s.trace = append(s.trace, 2, math.Float32bits(p.X), math.Float32bits(p.Y))
	if s.pathEndpoints[0] != nil || s.pathEndpoints[1] != nil {
		p := s.pathEndpoints[min(s.endpointCall, 1)]
		s.endpointCall++
		return p
	}
	return s.fallback
}

func (s *portTestRoamOwnerServer) Sub_50B810(u *server.Object, p *types.Pointf) bool {
	s.trace = append(s.trace, 3, math.Float32bits(p.X), math.Float32bits(p.Y))
	return s.precheck
}
