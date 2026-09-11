//go:build porttest

package server

import (
	"strings"
	"sync/atomic"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/sound"
)

// PortTestCombatAudioEvent is an event accepted by the normal deferred audio
// queue. Object events include the object's position at snapshot time.
type PortTestCombatAudioEvent struct {
	ID    sound.ID
	Obj   *Object
	Pos   types.Pointf
	Kind  int
	Code  uint32
	ByPos bool
}

// PortTestCombatAudioSnapshot returns a copy of events queued by Audio.EventObj
// and Audio.EventPos. It intentionally observes the default deferred queue;
// Audio.Init is neither needed nor wanted by combat fixtures.
func (s *Server) PortTestCombatAudioSnapshot() []PortTestCombatAudioEvent {
	out := make([]PortTestCombatAudioEvent, 0, len(s.Audio.delayedObj)+len(s.Audio.delayedPos))
	for _, it := range s.Audio.delayedObj {
		var pos types.Pointf
		if it.Obj != nil {
			pos = it.Obj.Pos()
		}
		out = append(out, PortTestCombatAudioEvent{
			ID: it.ID, Obj: it.Obj, Pos: pos, Kind: it.Kind, Code: it.Code,
		})
	}
	for _, it := range s.Audio.delayedPos {
		out = append(out, PortTestCombatAudioEvent{
			ID: it.ID, Pos: it.Pos, Kind: it.Kind, Code: it.Code, ByPos: true,
		})
	}
	return out
}

// PortTestCombatAudioReset drops queued audio events without delivering them.
func (s *Server) PortTestCombatAudioReset() {
	s.Audio.delayedObj = nil
	s.Audio.delayedPos = nil
}

// PortTestCombatProjectileType installs one real missile ObjectType and a
// C-backed object allocator on a fresh fixture server. Its cleanup frees every
// allocated projectile and restores the prior server state.
func (s *Server) PortTestCombatProjectileType(id string, speed float32) func() {
	if id == "" {
		panic("empty projectile type ID")
	}
	oldHandle, oldTypes, oldObjs := s.handle, s.Types, s.Objs
	h := atomic.AddUintptr(&serverLast, 1)
	servers.Store(h, s)
	s.handle = h
	s.Objs = serverObjects{}
	s.Objs.init(h)
	if !s.Objs.Init(8) {
		servers.Delete(h)
		s.handle, s.Types, s.Objs = oldHandle, oldTypes, oldObjs
		panic("projectile object allocator initialization failed")
	}
	s.Types = serverObjTypes{
		byInd: make([]*ObjectType, 2),
		byID:  make(map[string]*ObjectType),
	}
	name := strings.ToLower(id)
	typ := &ObjectType{
		s:         &s.Types,
		ind:       1,
		ind2:      1,
		id:        name,
		class:     object.ClassMissile,
		Speed:     speed,
		SpeedBase: speed,
	}
	s.Types.byInd[1] = typ
	s.Types.byID[name] = typ
	return func() {
		s.Objs.FreeObjects()
		s.handle, s.Types, s.Objs = oldHandle, oldTypes, oldObjs
		servers.Delete(h)
	}
}
