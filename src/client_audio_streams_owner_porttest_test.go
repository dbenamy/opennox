//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/timer"
	"testing"
	"unsafe"
)

type audioStreamCallback struct{ Op, Ref int }
type audioStreamOwner struct {
	root, descriptor, voiceAPI []uint32
	callbacks                  []unsafe.Pointer
	refs                       map[uint32]int
	next, slots, voices        int
	failOp, failAt             int
	callbackCounts             [14]int
	events                     []audioStreamCallback
	ticks                      uint64
	clockCalls                 int
}

func newAudioStreamOwner(t *testing.T, slots, voices int) *audioStreamOwner {
	t.Helper()
	o := &audioStreamOwner{refs: make(map[uint32]int), slots: slots, voices: voices, failOp: -1, ticks: 100}
	var free func()
	o.root, free = alloc.Make([]uint32{}, 32)
	t.Cleanup(free)
	o.descriptor, free = alloc.Make([]uint32{}, 10)
	t.Cleanup(free)
	o.voiceAPI, free = alloc.Make([]uint32{}, 10)
	t.Cleanup(free)
	clear(serverConfigOwnBytes(t, 0x5D4594, 1193332, 4))
	clear(serverConfigOwnBytes(t, 0x5D4594, 1193340, 4))
	oldLegacyClock := legacy.PlatformTicks
	legacy.PlatformTicks = func() uint64 { o.clockCalls++; return o.ticks }
	t.Cleanup(func() { legacy.PlatformTicks = oldLegacyClock })
	oldClock := timer.PlatformTicks
	timer.PlatformTicks = func() uint64 { o.clockCalls++; return o.ticks }
	t.Cleanup(func() { timer.PlatformTicks = oldClock })
	oldEnabled := dword_5d4594_1193336
	t.Cleanup(func() { dword_5d4594_1193336 = oldEnabled })
	var restore func()
	o.callbacks, restore = legacy.PortTestAudioStreamGlobalOwner(unsafe.Pointer(&o.root[0]), o.callback)
	t.Cleanup(restore)
	sub_486F30()
	for i, op := range []int{0, 1, 2, 3} {
		o.descriptor[5+i] = audioStreamPointer(o.callbacks[op])
	}
	o.descriptor[9] = audioStreamPointer(unsafe.Pointer(&o.voiceAPI[0]))
	for index, op := range map[int]int{1: 4, 2: 5, 3: 6, 4: 7, 8: 8, 9: 9} {
		o.voiceAPI[index] = audioStreamPointer(o.callbacks[op])
	}
	t.Cleanup(func() { o.call("sub_4875F0"); o.call("sub_4870A0") })
	return o
}
func (o *audioStreamOwner) call(name string, args ...uint32) uint32 {
	return legacy.PortTestAudioStreamCall(name, args...)
}
func (o *audioStreamOwner) callback(op int, p unsafe.Pointer) int {
	address := audioStreamPointer(p)
	if op == 0 || op == 2 || op == 4 {
		o.next++
		o.refs[address] = o.next
	}
	ref := o.refs[address]
	if ref == 0 {
		panic("audio callback received an unowned object")
	}
	o.events = append(o.events, audioStreamCallback{op, ref})
	o.callbackCounts[op]++
	if op == 0 {
		audioStreamWords(address, 22)[5] = uint32(o.slots)
	}
	if op == 2 {
		audioStreamWords(address, 66)[49] = uint32(o.voices)
	}
	if op == o.failOp && (o.failAt == 0 || o.callbackCounts[op] == o.failAt) {
		return -12345
	}
	return 0
}
func audioStreamWords(address uint32, n int) []uint32 {
	return unsafe.Slice((*uint32)(unsafe.Pointer(uintptr(address))), n)
}
func (o *audioStreamOwner) device() uint32 {
	return o.call("sub_486FA0", audioStreamPointer(unsafe.Pointer(&o.descriptor[0])))
}
func (o *audioStreamOwner) count() uint32 { return *memmap.PtrUint32(0x5D4594, 1193332) }
