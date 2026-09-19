//go:build porttest

package legacy

import (
	"fmt"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestUnitUpdateSpec struct {
	Frame, Spawn uint32
	Target       bool
	TargetState  byte
}

func (p *portTestShopPools) unitUpdateContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.UnitUpdate
	u, target := p.items[0].u, p.items[1].u
	callback, size := server.PortTestUnitGameplayRegistration("UndeadKillerUpdate", false)
	if size != 0 {
		panic("undead update registration data size")
	}
	p.identify(callback, 93064)
	u.Update = callback
	init := p.objectiveRegion(4)
	oldCollision := u.CollideData
	defer func() { u.CollideData = oldCollision }()
	u.CollideData = init
	if sp.Target {
		*(*unsafe.Pointer)(init) = target.CObj()
	}
	*(*byte)(unsafe.Add(target.CObj(), 88)) = sp.TargetState
	u.Field34 = sp.Spawn
	p.proxy.core.SetFrame(sp.Frame)
	before := len(p.proxy.trace)
	u.CallUpdate()
	want := 0
	if sp.Target && sp.TargetState&1 != 0 || sp.Frame-sp.Spawn > 70 {
		want = 1
	}
	trace := p.proxy.trace[before:]
	if len(trace) != 2*want || want != 0 && (trace[0] != 32 || trace[1] != p.proxy.life.ids[uint32(uintptr(u.CObj()))]) {
		panic(fmt.Sprintf("undead deletion %v want %d", trace, want))
	}
	return []uint32{uint32(want), u.Field34}
}
