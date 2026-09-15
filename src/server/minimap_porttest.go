//go:build porttest

package server

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// PortTestMinimapWalls installs a fresh actual wall owner. Rendering uses the
// normal creation, row index, positional lookup and cleanup implementations.
func (s *Server) PortTestMinimapWalls() func() {
	old := s.Walls
	s.Walls = serverWalls{}
	if s.Walls.Init() == 0 {
		panic("minimap wall owner allocation")
	}
	s.Walls.defsCnt = 3
	for i, name := range []string{"MinimapWall", "InvisibleWallSet", "InvisibleBlockingWallSet"} {
		copy(s.Walls.defs[i].Field0[:], name)
	}
	s.Walls.fast.magicSysUse = -1
	s.Walls.fast.invis = -1
	s.Walls.fast.invisBlock = -1
	return func() { s.Walls.Free(); s.Walls = old }
}

// PortTestMinimapDebug owns the actual path-point buffer and object pool/list
// consumed by the optional debug overlay. Configuration supplies input data only.
func (s *Server) PortTestMinimapDebug() (func([]types.Pointf, int), func()) {
	oldObjects := s.Objs
	oldPoints, oldCount := s.AI.Paths.points, s.AI.Paths.pointsCnt
	s.Objs = serverObjects{}
	s.Objs.init(s.handle)
	if !s.Objs.Init(5) {
		panic("minimap debug object pool")
	}
	points, freePoints := alloc.Make([]types.Pointf{}, 8)
	s.AI.Paths.points = points
	s.AI.Paths.pointsCnt = 0
	typ := ObjectType{ind: 1, ind2: 1, class: object.ClassSimple, flags: object.FlagNoCollide}
	units := make([]*Object, 5)
	for i := range units {
		units[i] = s.Objs.NewObject(&typ)
		units[i].PosVec = types.Ptf(float32(215+i*7), float32(222+i*3))
		if i > 0 {
			units[i-1].ObjNext = units[i]
			units[i].ObjPrev = units[i-1]
		}
	}
	configure := func(path []types.Pointf, monsters int) {
		s.AI.Paths.ResetPoints()
		for _, p := range path {
			if !s.AI.Paths.appendPoint(p) {
				panic("minimap debug path capacity")
			}
		}
		for i, u := range units {
			u.ObjClass = object.ClassSimple
			if (i == 1 && monsters > 0) || (i == 3 && monsters > 1) || (i == 4 && monsters > 2) {
				u.ObjClass = object.ClassMonster
			}
		}
		if monsters < 0 {
			s.Objs.SetObjects(nil)
		} else {
			s.Objs.SetObjects(units[0])
		}
	}
	configure(nil, -1)
	return configure, func() {
		s.Objs.SetObjects(nil)
		for _, u := range units {
			s.Objs.FreeObject(u)
		}
		s.Objs.FreeObjects()
		s.Objs = oldObjects
		s.AI.Paths.points, s.AI.Paths.pointsCnt = oldPoints, oldCount
		freePoints()
	}
}

// The lightweight rendering server has teams but no loaded team definitions.
// Supply normal palette input to its real GetTeamColor lookup and restore it.
func (s *Server) PortTestMinimapTeamColors() func() {
	old := s.Teams.defs
	s.Teams.defs = map[TeamColor]*TeamDef{
		TeamNone: {Color: nox_color_white_2523948},
		TeamRed:  {Color: nox_color_red_2589776},
		TeamBlue: {Color: nox_color_blue_2650684},
	}
	return func() { s.Teams.defs = old }
}
