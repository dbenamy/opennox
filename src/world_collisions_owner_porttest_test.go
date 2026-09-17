//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

type worldCollisionOwner struct {
	*reliableReportsOwner
	globals  map[string]*uint32
	throttle *uint64
	ticks    uint64
	balance  func(map[string]float64)
}

func newWorldCollisionOwner(t *testing.T) *worldCollisionOwner {
	t.Helper()
	o := &worldCollisionOwner{reliableReportsOwner: newReliableReportsOwner(t), ticks: 10000}
	t.Cleanup(server.PortTestAttachAI(o.s, &o.units[0], &o.units[1], &o.units[2]))
	o.s.NoxScriptVM.Init(o.s)
	noxServer.noxScript.Init(noxServer)
	output := memmap.PtrUint32(0x5D4594, 1599076)
	oldOutput := *output
	t.Cleanup(func() { *output = oldOutput })

	var restore func()
	o.globals, o.throttle, restore = legacy.PortTestWorldCollisionGlobals()
	t.Cleanup(restore)
	oldClock := legacy.PlatformTicks
	legacy.PlatformTicks = func() uint64 { return o.ticks }
	t.Cleanup(func() { legacy.PlatformTicks = oldClock })
	o.balance, restore = o.s.PortTestWorldCollisionBalance()
	t.Cleanup(restore)
	o.balance(map[string]float64{"QuestExitTimerStart": 30, "MaxExtraLives": 5, "InversionRange": 100})
	stage := memmap.PtrUint32(0x587000, 202028)
	oldStage := *stage
	*stage = 0
	t.Cleanup(func() { *stage = oldStage })
	return o
}
func (o *worldCollisionOwner) record(t *testing.T, size int) unsafe.Pointer {
	t.Helper()
	b, free := alloc.Make([]byte{}, size)
	t.Cleanup(free)
	return unsafe.Pointer(&b[0])
}
