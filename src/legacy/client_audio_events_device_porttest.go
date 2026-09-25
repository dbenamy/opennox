//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/client/audio/ail"
	"unsafe"
)

type portTestAudioEventDevice struct {
	sample ail.Sample
	user   *AudioSample
	ready  func() int
	load   func(uint32, []byte)
}

var portTestAudioEventDeviceOwner *portTestAudioEventDevice

// PortTestAudioEventDevice supplies the external device queue only. The actual
// native refill/end callback bodies and voice/buffer owners continue to execute.
func PortTestAudioEventDevice(sample ail.Sample, user *AudioSample, ready func() int, load func(uint32, []byte)) func() {
	old := portTestAudioEventDeviceOwner
	portTestAudioEventDeviceOwner = &portTestAudioEventDevice{sample, user, ready, load}
	return func() { portTestAudioEventDeviceOwner = old }
}

func nox_porttest_audio_event_user_data(h uintptr) unsafe.Pointer {
	sample := ail.Sample(h)
	if o := portTestAudioEventDeviceOwner; o != nil {
		if sample != o.sample {
			panic("unexpected audio device sample")
		}
		return unsafe.Pointer(o.user)
	}
	p := sample.UserData()
	if p == nil {
		return nil
	}
	return unsafe.Pointer(p.(*AudioSample))
}

func nox_porttest_audio_event_buffer_ready(h uintptr) int32 {
	sample := ail.Sample(h)
	if o := portTestAudioEventDeviceOwner; o != nil {
		if sample != o.sample {
			panic("unexpected audio device sample")
		}
		return int32(o.ready())
	}
	return int32(sample.BufferReady())
}

func nox_porttest_audio_event_load_buffer(h uintptr, n uint32, p unsafe.Pointer, size uint32) {
	sample := ail.Sample(h)
	data := unsafe.Slice((*byte)(p), int(size))
	if o := portTestAudioEventDeviceOwner; o != nil {
		if sample != o.sample {
			panic("unexpected audio device sample")
		}
		o.load(uint32(n), data)
		return
	}
	sample.LoadBuffer(uint32(n), data)
}

func init() {
	audioEventDeviceUser = func(s ail.Sample) *AudioSample {
		return (*AudioSample)(nox_porttest_audio_event_user_data(uintptr(s)))
	}
	audioEventDeviceReady = func(s ail.Sample) int { return int(nox_porttest_audio_event_buffer_ready(uintptr(s))) }
	audioEventDeviceLoad = func(s ail.Sample, index uint32, data []byte) {
		nox_porttest_audio_event_load_buffer(uintptr(s), index, unsafe.Pointer(unsafe.SliceData(data)), uint32(len(data)))
	}
}
