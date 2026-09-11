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
