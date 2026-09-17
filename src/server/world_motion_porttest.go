//go:build porttest

package server

import "github.com/opennox/noxscript/ns/asm"

// Install two small map-script functions in the actual interpreter. The returned
// indices reject/accept trigger contact; script dispatch and return handling stay live.
func (s *NoxScriptVM) PortTestWorldMotionPredicates() (int32, int32, func()) {
	old := s.vm
	s.vm.funcs = make([]ScriptFunc, 4)
	for i := 0; i < 4; i++ {
		s.vm.funcs[i] = ScriptFunc{FuncDef: asm.FuncDef{Name: "motion-predicate", Return: 1, Code: []uint32{uint32(asm.OpPushInt), uint32(i & 1), uint32(asm.OpReturn)}}}
	}
	s.vm.stack = nil
	s.vm.callbacks = nil
	return 2, 3, func() { s.vm = old }
}
