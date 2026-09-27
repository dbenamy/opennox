//go:build porttest

package legacy

import "unsafe"

var portTestAudioStreamNames = []string{
	"sub_486640",
	"sub_4866D0",
	"sub_486AA0",
	"sub_486B60",
	"sub_486DB0",
	"sub_486E00",
	"sub_486E30",
	"sub_486E90",
	"sub_486FA0",
	"sub_486FE0",
	"sub_487030",
	"sub_487050",
	"sub_487070",
	"sub_487090",
	"sub_4870A0",
	"sub_4870E0",
	"sub_487100",
	"sub_487150",
	"sub_4871C0",
	"sub_4872C0",
	"sub_487310",
	"sub_487360",
	"sub_4873C0",
	"sub_487590",
	"sub_4875B0",
	"sub_4875D0",
	"sub_4875F0",
	"sub_487680",
	"sub_4876A0",
	"sub_487750",
	"sub_487790",
	"sub_4877D0",
	"sub_4877F0",
	"sub_487810",
	"sub_487910",
	"sub_487970",
	"sub_487C30",
	"sub_487C50",
	"sub_487C80",
	"sub_487D00",
	"sub_487D30",
	"sub_487D60",
	"sub_4BD280",
	"sub_4BD2D0",
	"sub_4BD2E0",
	"sub_4BD300",
	"sub_4BD340",
	"sub_4BD3C0",
	"sub_4BD420",
	"sub_4BD470",
	"sub_4BD600",
	"sub_4BD650",
	"sub_4BD660",
	"sub_4BD680",
	"sub_4BD690",
	"sub_4BD710",
	"sub_4BD720",
	"sub_4BD7A0",
	"sub_4BD7C0",
	"sub_4BD840",
	"sub_4BD8C0",
	"sub_4BD940",
	"sub_4BD9B0",
	"sub_4BDA60",
	"sub_4BDA80",
	"sub_4BDB20",
	"sub_4BDB30",
	"sub_4BDB40",
	"sub_4BDB90",
	"sub_4BDC00",
}

func PortTestAudioStreamCall(name string, args ...uint32) uint32 {
	if len(args) > 4 {
		panic("too many audio stream arguments")
	}
	var a [4]uint32
	copy(a[:], args)
	switch name {
	case "sub_487680":
		audioStreamContextDestroy((*audioStreamContext)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_487970":
		return uint32(audioStreamVoiceStopKind((*audioStreamContext)(unsafe.Pointer(uintptr(a[0]))), int32(a[1])))

	case "sub_486640":
		return uint32(audioStreamScaleVolume(unsafe.Pointer(uintptr(a[0])), uint32(a[1])))
	case "sub_4866D0":
		return uint32(uintptr(unsafe.Pointer(audioStreamCatalogEntry((*audioStreamCatalog)(unsafe.Pointer(uintptr(a[0]))), int32(a[1])))))
	case "sub_486AA0":
		return uint32(audioStreamReadFormat((*audioStreamCatalog)(unsafe.Pointer(uintptr(a[0]))), int32(a[1]), (*audioStreamFormat)(unsafe.Pointer(uintptr(a[2])))))
	case "sub_486B60":
		return uint32(audioStreamOpen((*audioStreamCatalog)(unsafe.Pointer(uintptr(a[0]))), int32(a[1])))
	case "sub_486DB0":
		return uint32(audioStreamRead((*audioStreamCatalog)(unsafe.Pointer(uintptr(a[0]))), unsafe.Pointer(uintptr(a[1])), int32(a[2])))
	case "sub_486E00":
		audioStreamClose((*audioStreamCatalog)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_486E30":
		return uint32(audioStreamVoiceLink((*audioStreamContext)(unsafe.Pointer(uintptr(a[0]))), (*audioStreamVoice)(unsafe.Pointer(uintptr(a[1])))))
	case "sub_486E90":
		return uint32(audioStreamVoiceUnlink((*audioStreamVoice)(unsafe.Pointer(uintptr(a[0])))))
	case "sub_486FA0":
		return uint32(uintptr(unsafe.Pointer(audioStreamDeviceRegister((*audioStreamDescriptor)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_486FE0":
		return uint32(uintptr(unsafe.Pointer(audioStreamDeviceNew((*audioStreamDescriptor)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_487030":
		audioStreamDeviceFree((*audioStreamDevice)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_487050":
		audioStreamDeviceLink((*audioStreamDevice)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_487070":
		audioStreamDeviceDestroy((*audioStreamDevice)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_487090":
		audioStreamDeviceUnlink((*audioStreamDevice)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_4870A0":
		audioStreamDeviceDestroyAll()
		return 0
	case "sub_4870E0":
		return uint32(uintptr(unsafe.Pointer(audioStreamDeviceFirst((**audioStreamDevice)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_487100":
		return uint32(uintptr(unsafe.Pointer(audioStreamDeviceNext((**audioStreamDevice)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_487150":
		return uint32(uintptr(unsafe.Pointer(audioStreamContextAcquire(int32(a[0]), (*audioStreamFormat)(unsafe.Pointer(uintptr(a[1])))))))
	case "sub_4871C0":
		return uint32(uintptr(unsafe.Pointer(audioStreamContextNew((*audioStreamDevice)(unsafe.Pointer(uintptr(a[0]))), int32(a[1]), (*audioStreamFormat)(unsafe.Pointer(uintptr(a[2])))))))
	case "sub_4872C0":
		audioStreamContextFree((*audioStreamContext)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_487310":
		return uint32(audioStreamContextLink((*audioStreamContext)(unsafe.Pointer(uintptr(a[0])))))
	case "sub_487360":
		return uint32(uintptr(unsafe.Pointer(audioStreamDeviceSlot(int32(a[0]), (**audioStreamDevice)(unsafe.Pointer(uintptr(a[1]))), (*int32)(unsafe.Pointer(uintptr(a[2])))))))
	case "sub_487590":
		return uint32(uintptr(unsafe.Pointer(audioStreamContextFormat((*audioStreamContext)(unsafe.Pointer(uintptr(a[0]))), (*audioStreamFormat)(unsafe.Pointer(uintptr(a[1])))))))
	case "sub_4875B0":
		return uint32(uintptr(unsafe.Pointer(audioStreamContextFirst((**audioStreamContext)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_4875D0":
		return uint32(uintptr(unsafe.Pointer(audioStreamContextNext((**audioStreamContext)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_4875F0":
		return uint32(audioStreamContextDestroyAll())
	case "sub_4876A0":
		return uint32(audioStreamContextUnlink((*audioStreamContext)(unsafe.Pointer(uintptr(a[0])))))
	case "sub_487750":
		return uint32(uintptr(unsafe.Pointer(audioStreamVoiceCreate((*audioStreamContext)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_487790":
		return uint32(audioStreamVoiceCreateMany((*audioStreamContext)(unsafe.Pointer(uintptr(a[0]))), int32(a[1])))
	case "sub_4877D0":
		return uint32(uintptr(unsafe.Pointer(audioStreamVoiceFirst((*audioStreamContext)(unsafe.Pointer(uintptr(a[0]))), (**audioStreamVoice)(unsafe.Pointer(uintptr(a[1])))))))
	case "sub_4877F0":
		return uint32(uintptr(unsafe.Pointer(audioStreamVoiceNext((**audioStreamVoice)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_487910":
		return uint32(audioStreamVoiceDeleteKind((*audioStreamContext)(unsafe.Pointer(uintptr(a[0]))), int32(a[1])))
	case "sub_487C30":
		audioStreamBufferInit((*audioStreamBuffer)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_487C50":
		return uint32(audioStreamBufferAppend((*audioStreamBuffer)(unsafe.Pointer(uintptr(a[0]))), (*audioStreamChunk)(unsafe.Pointer(uintptr(a[1])))))
	case "sub_487C80":
		return uint32(uintptr(unsafe.Pointer(audioStreamBufferFirst((*audioStreamBuffer)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_487D00":
		return uint32(audioStreamByteRate((*audioStreamFormat)(unsafe.Pointer(uintptr(a[0])))))
	case "sub_487D30":
		return uint32(uintptr(unsafe.Pointer(audioStreamChunkInit((*audioStreamChunk)(unsafe.Pointer(uintptr(a[0]))), a[1], uint32(a[2])))))
	case "sub_487D60":
		return uint32(uintptr(unsafe.Pointer(audioStreamChunkClear((*audioStreamChunk)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_4BD420":
		return uint32(uintptr(unsafe.Pointer(audioStreamCacheFind((*audioStreamCache)(unsafe.Pointer(uintptr(a[0]))), int32(a[1])))))
	case "sub_4BD600":
		return uint32(audioStreamCacheEvict((*audioStreamCache)(unsafe.Pointer(uintptr(a[0])))))
	case "sub_4BD680":
		return uint32(audioStreamCacheReferences((*audioStreamCacheEntry)(unsafe.Pointer(uintptr(a[0])))))
	case "sub_4BD690":
		return uint32(uintptr(unsafe.Pointer(audioStreamCacheDrop((*audioStreamCacheEntry)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_4BD720":
		return uint32(uintptr(unsafe.Pointer(audioStreamVoiceNew((*audioStreamContext)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_4BD7A0":
		audioStreamVoiceFree((*audioStreamVoice)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_4BD7C0":
		return uint32(uintptr(unsafe.Pointer(audioStreamVoiceInit((*audioStreamVoice)(unsafe.Pointer(uintptr(a[0])))))))
	case "sub_4BD840":
		audioStreamVoiceMix((*audioStreamVoice)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_4BDA60":
		audioStreamVoiceDestroy((*audioStreamVoice)(unsafe.Pointer(uintptr(a[0]))))
		return 0
	case "sub_4BDC00":
		return uint32(uintptr(unsafe.Pointer(audioStreamVoiceState(unsafe.Pointer(uintptr(a[0]))))))
	}
	for op, n := range portTestAudioStreamNames {
		if n == name {
			switch op {
			case 22:
				return uint32(audioStreamContextTick((*audioStreamContext)(unsafe.Pointer(uintptr(uint32(a[0]))))))
			case 33:
				return uint32(uintptr(unsafe.Pointer(audioStreamVoiceSelect((*audioStreamContext)(unsafe.Pointer(uintptr(uint32(a[0])))), int32(a[1])))))
			case 42:
				return uint32(uintptr(unsafe.Pointer(audioStreamPoolNew(int32(a[0]), int32(a[1])))))
			case 43:
				audioStreamPoolFree(unsafe.Pointer(uintptr(a[0])))
				return 0
			case 44:
				return uint32(uintptr(audioStreamPoolPop((*unsafe.Pointer)(unsafe.Pointer(uintptr(a[0]))))))
			case 45:
				return uint32(uintptr(audioStreamPoolPush((*unsafe.Pointer)(unsafe.Pointer(uintptr(a[0]))), unsafe.Pointer(uintptr(uint32(a[1]))))))
			case 46:
				return uint32(uintptr(unsafe.Pointer(audioStreamCacheNew((*audioStreamCatalog)(unsafe.Pointer(uintptr(uint32(a[0])))), int32(a[1]), int32(a[2]), int32(a[3])))))
			case 47:
				audioStreamCacheFree((*audioStreamCache)(unsafe.Pointer(uintptr(a[0]))))
				return 0
			case 49:
				return uint32(uintptr(unsafe.Pointer(audioStreamCacheLoad((*audioStreamCache)(unsafe.Pointer(uintptr(a[0]))), int32(a[1])))))
			case 51:
				return uint32(uintptr(unsafe.Pointer(audioStreamCacheRef((*audioStreamCacheEntry)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
			case 52:
				return uint32(audioStreamCacheUnref((*audioStreamCacheEntry)(unsafe.Pointer(uintptr(uint32(a[0]))))))
			case 55:
				return uint32(uintptr(unsafe.Pointer(audioStreamCacheBuffer((*audioStreamCacheEntry)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
			case 60:
				return uint32(audioStreamVoiceData((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a[0]))))))
			case 61:
				return uint32(audioStreamVoiceLoop((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a[0]))))))
			case 62:
				return uint32(audioStreamVoiceEnd((*audioStreamVoice)(unsafe.Pointer(uintptr(a[0])))))
			case 64:
				return uint32(audioStreamVoiceStop((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a[0]))))))
			case 65:
				return uint32(uintptr(unsafe.Pointer(audioStreamVoiceReserve((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
			case 66:
				return uint32(uintptr(unsafe.Pointer(audioStreamVoiceUnreserve((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
			case 67:
				return uint32(audioStreamVoiceStart((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a[0]))))))
			case 68:
				audioStreamVoiceBind((*audioStreamVoice)(unsafe.Pointer(uintptr(a[0]))), (*audioStreamBuffer)(unsafe.Pointer(uintptr(a[1]))))
				return 0
			default:
				return 0 // Preserve the original dispatcher default.
			}
		}
	}
	panic("unknown audio stream operation: " + name)
}

var portTestAudioStreamCallback func(int, unsafe.Pointer) int

var portTestAudioStreamKeys [14]byte

func init() {
	for op := range portTestAudioStreamKeys {
		audioStreamCallbacks[unsafe.Pointer(&portTestAudioStreamKeys[op])] = func(obj unsafe.Pointer) int32 {
			if portTestAudioStreamCallback == nil {
				panic("audio callback without an owner")
			}
			return int32(portTestAudioStreamCallback(op, obj))
		}
	}
}

func PortTestAudioStreamGlobalOwner(p unsafe.Pointer, fn func(int, unsafe.Pointer) int) ([]unsafe.Pointer, func()) {
	old, callback := audioStreamsRoot, portTestAudioStreamCallback
	audioStreamsRoot, portTestAudioStreamCallback = (*audioStreamSystem)(p), fn
	return []unsafe.Pointer{
		unsafe.Pointer(&portTestAudioStreamKeys[0]),
		unsafe.Pointer(&portTestAudioStreamKeys[1]),
		unsafe.Pointer(&portTestAudioStreamKeys[2]),
		unsafe.Pointer(&portTestAudioStreamKeys[3]),
		unsafe.Pointer(&portTestAudioStreamKeys[4]),
		unsafe.Pointer(&portTestAudioStreamKeys[5]),
		unsafe.Pointer(&portTestAudioStreamKeys[6]),
		unsafe.Pointer(&portTestAudioStreamKeys[7]),
		unsafe.Pointer(&portTestAudioStreamKeys[8]),
		unsafe.Pointer(&portTestAudioStreamKeys[9]),
		unsafe.Pointer(&portTestAudioStreamKeys[10]),
		unsafe.Pointer(&portTestAudioStreamKeys[11]),
		unsafe.Pointer(&portTestAudioStreamKeys[12]),
		unsafe.Pointer(&portTestAudioStreamKeys[13]),
	}, func() { audioStreamsRoot, portTestAudioStreamCallback = old, callback }
}
