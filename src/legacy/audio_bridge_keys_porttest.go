//go:build porttest

package legacy

/*
#include "GAME1_3.h"
*/
import "C"
import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

// Order follows the nine voice API words and four device descriptor callbacks.
func PortTestAudioBridgeKeys() [13]unsafe.Pointer {
	return [13]unsafe.Pointer{
		C.sub_43EC30, C.sub_43ECB0, C.sub_43ED00, C.sub_43EFD0,
		C.sub_43F010, C.sub_43F030, C.sub_43F050, C.sub_43F060, C.sub_43F0D0,
		C.sub_43E940, C.sub_43E9F0, C.sub_43EA20, C.sub_43EC10,
	}
}

// Own only the shipped audio table region, leaving the rest of blob initialization alone.
func PortTestAudioBridgeTableBindings() ([13]unsafe.Pointer, func()) {
	base := memmap.PtrOff(0x587000, 93952)
	data := unsafe.Slice((*byte)(base), 84)
	saved := bytes.Clone(data)
	keys := PortTestAudioBridgeKeys()
	offsets := [...]uintptr{4, 8, 12, 16, 20, 24, 28, 32, 36, 60, 64, 68, 72}
	for i, off := range offsets {
		*(*uint32)(unsafe.Add(base, off)) = uint32(uintptr(keys[i]))
	}
	*(*uint32)(unsafe.Add(base, 76)) = uint32(uintptr(base))
	*(*uint32)(unsafe.Add(base, 80)) = uint32(uintptr(unsafe.Add(base, 40)))
	api := (*audioStreamVoiceAPI)(base)
	desc := (*audioStreamDescriptor)(unsafe.Add(base, 40))
	out := [13]unsafe.Pointer{api.Init, api.Free, api.Start, api.Stop,
		unsafe.Pointer(uintptr(api.Fields20[0])), unsafe.Pointer(uintptr(api.Fields20[1])), unsafe.Pointer(uintptr(api.Fields20[2])),
		api.Update, api.Restart, desc.Init, desc.Free, desc.ContextInit, desc.ContextFree}
	if desc.VoiceAPI != api {
		copy(data, saved)
		panic("audio table API link mismatch")
	}
	return out, func() { copy(data, saved) }
}
