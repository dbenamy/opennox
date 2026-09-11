//go:build porttest

package server

import "sync/atomic"

// PortTestAttachAI registers an isolated server for C-owned fixture objects.
func PortTestAttachAI(s *Server, objects ...*Object) func() {
	h := atomic.AddUintptr(&serverLast, 1)
	servers.Store(h, s)
	old := make([]uintptr, len(objects))
	for i, obj := range objects {
		old[i], obj.serverHandle = obj.serverHandle, h
	}
	return func() {
		for i, obj := range objects {
			obj.serverHandle = old[i]
		}
		servers.Delete(h)
	}
}

// PortTestAIEmptyMap supplies empty spatial/wall indices for actual waypoint
// lookup and ray traversal, without loading map assets or allocating wall objects.
func (s *Server) PortTestAIEmptyMap() {
	if s.Map.bucketByPos == nil {
		s.Map.Init()
	} else {
		for _, row := range s.Map.bucketByPos {
			clear(row)
		}
	}
	if s.Walls.byPos == nil {
		s.Walls.byPos = make([]*Wall, wallsPerBucket*WallGridSize)
	}
}
