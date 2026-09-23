//go:build porttest

package legacy

import (
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

//go:noinline
func audioStreamRetiredAddress() uint32 {
	// The number deliberately falls inside a freed Go span. It is an opaque ABI
	// word, never dereferenced; treating it as a managed pointer would be a bug.
	b := make([]byte, 32<<20)
	b[4096] = 1
	address := uint32(uintptr(unsafe.Pointer(&b[4096])))
	runtime.KeepAlive(b)
	return address
}

func TestAudioStreamRawAddressGC(t *testing.T) {
	if os.Getenv("OPENNOX_AUDIO_RAW_GC_CHILD") != "1" {
		command := exec.Command(os.Args[0], "-test.run=^TestAudioStreamRawAddressGC$", "-test.count=1")
		command.Env = append(os.Environ(), "OPENNOX_AUDIO_RAW_GC_CHILD=1")
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("raw audio address child: %v\n%s", err, output)
		}
		return
	}
	chunk, freeChunk := alloc.New(audioStreamChunk{})
	defer freeChunk()
	buffer, freeBuffer := alloc.New(audioStreamBuffer{})
	defer freeBuffer()
	voice, freeVoice := alloc.New(audioStreamVoice{})
	defer freeVoice()
	address := audioStreamRetiredAddress()
	debug.FreeOSMemory()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100; i++ {
			runtime.GC()
		}
	}()
	loops := 0
	for {
		audioStreamBufferInit(buffer)
		PortTestAudioStreamCall("sub_487D30", uint32(uintptr(unsafe.Pointer(chunk))), address, 19)
		audioStreamBufferAppend(buffer, chunk)
		audioStreamVoiceBind(voice, buffer)
		if *(*uint32)(unsafe.Add(unsafe.Pointer(chunk), 12)) != address || *(*uint32)(unsafe.Add(unsafe.Pointer(voice), 296)) != address || voice.Length != 19 {
			t.Fatal("raw address or length changed")
		}
		audioStreamVoiceBind(voice, nil)
		loops++
		select {
		case <-done:
			t.Logf("preserved opaque address through %d buffer/chunk/voice cycles during100 GCs", loops)
			return
		default:
		}
	}
}
