package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

// Non-executable identities for the internal stream callbacks. Foreign callback
// addresses retain the call convention chosen by each original caller.
var audioStreamCallbackKeys [4]byte

func audioStreamKnownCallback(fn, arg unsafe.Pointer) (int32, bool) {
	switch fn {
	case unsafe.Pointer(&audioStreamCallbackKeys[0]):
		return audioStreamContextTick((*audioStreamContext)(arg)), true
	case unsafe.Pointer(&audioStreamCallbackKeys[1]):
		return audioStreamVoiceData((*audioStreamVoice)(arg)), true
	case unsafe.Pointer(&audioStreamCallbackKeys[2]):
		return audioStreamVoiceLoop((*audioStreamVoice)(arg)), true
	case unsafe.Pointer(&audioStreamCallbackKeys[3]):
		return audioStreamVoiceEnd((*audioStreamVoice)(arg)), true
	default:
		return 0, false
	}
}

func AudioStreamCallbackVoid(fn, arg unsafe.Pointer) {
	if _, ok := audioStreamKnownCallback(fn, arg); !ok {
		ccall.CallVoidPtr(fn, arg)
	}
}

func AudioStreamCallbackInt(fn, arg unsafe.Pointer) int32 {
	if result, ok := audioStreamKnownCallback(fn, arg); ok {
		return result
	}
	return int32(ccall.CallIntPtr(fn, arg))
}
