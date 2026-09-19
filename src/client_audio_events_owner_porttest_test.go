//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/legacy/timer"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

type audioEventsOwner struct {
	*audioStreamOwner
	words        map[string]*uint32
	rows         []byte
	ctx          uint32
	cat, entries []uint32
	s            *server.Server
	globalTimers *timer.TimerGroup
}

func newAudioEventsOwner(t *testing.T, voices int) *audioEventsOwner {
	t.Helper()
	t.Cleanup(handles.PortTestInit())
	o := new(audioEventsOwner)
	var restore func()
	o.words, restore = legacy.PortTestAudioEventGlobals()
	t.Cleanup(restore)
	o.rows, _, _, restore = legacy.PortTestClientAudioAssetsOwner()
	t.Cleanup(restore)
	clear(serverConfigOwnBytes(t, 0x5D4594, 839892, 736))
	clear(serverConfigOwnBytes(t, 0x5D4594, 1045228, 96))
	clear(serverConfigOwnBytes(t, 0x5D4594, 1045440, 12))
	clear(serverConfigOwnBytes(t, 0x587000, 127000, 4))
	o.audioStreamOwner = newAudioStreamOwner(t, 1, voices)
	if o.device() == 0 {
		t.Fatal("event device")
	}
	o.ctx = o.call("sub_487150", 0, 0)
	if o.ctx == 0 {
		t.Fatal("event context")
	}
	if o.call("sub_487790", o.ctx, uint32(voices)) != uint32(voices) {
		t.Fatal("event voices")
	}
	o.s = new(server.Server)
	o.s.Rand.Other = prand.New(31)
	o.s.Rand.Logic = prand.New(17)
	oldServer := legacy.GetServer
	legacy.GetServer = func() legacy.Server { return &audioEventRandomOwner{s: o.s} }
	t.Cleanup(func() { legacy.GetServer = oldServer })
	var free func()
	o.cat, free = alloc.Make([]uint32{}, 72)
	t.Cleanup(free)
	o.entries, free = alloc.Make([]uint32{}, 27)
	t.Cleanup(free)
	o.cat[0] = audioStreamPointer(unsafe.Pointer(&o.entries[0]))
	o.cat[1] = 3
	var bag []byte
	for i, n := range []int{8, 20, 32} {
		e := o.entries[9*i : 9*i+9]
		e[4], e[5], e[6], e[7] = uint32(len(bag)), uint32(n), 22050, 4
		for j := 0; j < n; j++ {
			bag = append(bag, byte(11+i*67+j*13))
		}
	}
	files, _ := prefabScriptsFiles(t, bag)
	o.cat[67] = audioStreamPointer(files[0])
	o.globalTimers, free = alloc.New(timer.TimerGroup{})
	t.Cleanup(free)
	o.globalTimers.Init()
	*o.words["dword_587000_127004"] = audioStreamPointer(unsafe.Pointer(o.globalTimers))
	if o.eventCall("sub_451850", uint64(o.ctx), audioEventPointer(unsafe.Pointer(&o.cat[0]))) != 1 {
		t.Fatal("event manager init")
	}
	*o.words["dword_587000_126996"] = 1
	t.Cleanup(func() { o.eventCall("sub_451970") })
	return o
}
func (o *audioEventsOwner) eventCall(name string, args ...uint64) uint64 {
	return legacy.PortTestAudioEventCall(name, args...)
}
func (o *audioEventsOwner) metadata(index int, flags uint32, samples int) uint32 {
	p := audioStreamPointer(unsafe.Pointer(&o.rows[index*200]))
	w := audioStreamWords(p, 50)
	w[0] = 1
	w[1] = flags
	w[48] = uint32(samples)
	for i := 0; i < samples; i++ {
		binary.LittleEndian.PutUint16(o.rows[index*200+128+2*i:], uint16(i%3))
	}
	return p
}
func (o *audioEventsOwner) create(meta uint32) uint32 {
	return uint32(o.eventCall("nox_xxx_draw_452300", uint64(meta)))
}
