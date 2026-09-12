//go:build porttest

package server

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"unsafe"
)

// PortTestSpellEffectTypes preserves real type names and supplies complete
// monster update storage for summon factories, using the fixture's live template.
func (s *Server) PortTestSpellEffectTypes(names, monsters []string, template *Object) func() {
	saved := make(map[*ObjectType]ObjectType)
	var frees []func()
	save := func(t *ObjectType) {
		if _, ok := saved[t]; !ok {
			saved[t] = *t
		}
	}
	for _, name := range names {
		if t := s.Types.ByID(name); t != nil {
			save(t)
			t.id = name
		}
	}
	for _, name := range monsters {
		t := s.Types.ByID(name)
		if t == nil {
			panic("missing spell effect monster type: " + name)
		}
		save(t)
		t.class = object.ClassMonster
		t.subclass = 0
		t.Shape.Kind = ShapeKindCenter
		t.Mass = 10
		data, free := alloc.Make([]byte{}, 2216)
		frees = append(frees, free)
		copy(data[:2200], unsafe.Slice((*byte)(template.UpdateData), 2200))
		for i := 2200; i < 2216; i++ {
			data[i] = 0xa5
		}
		t.UpdateData = unsafe.Pointer(&data[0])
		t.UpdateDataSize = 2216
		hp, freeHP := alloc.New(HealthData{})
		*hp = HealthData{Cur: 100, Max: 100}
		frees = append(frees, freeHP)
		t.health = hp
		t.id = name
		s.Types.byID[strings.ToLower(name)] = t
	}
	return func() {
		for t, v := range saved {
			*t = v
		}
		for _, f := range frees {
			f()
		}
	}
}
