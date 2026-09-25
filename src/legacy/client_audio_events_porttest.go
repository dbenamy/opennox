//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/client/audio/ail"
	"unsafe"
)

var portTestAudioEventNames = []string{
	"sub_43DA80",
	"sub_43DAD0",
	"sub_43DB20",
	"sub_43DB30",
	"sub_43DB40",
	"sub_43DB60",
	"sub_43DBA0",
	"sub_43DC00",
	"sub_43DC10",
	"sub_43DC30",
	"sub_43EDB0",
	"sub_43EE00",
	"sub_43F0E0",
	"sub_44D960",
	"sub_44D970",
	"sub_44D990",
	"sub_450750",
	"sub_450760",
	"sub_451850",
	"sub_451920",
	"sub_451970",
	"sub_4519C0",
	"sub_451BE0",
	"sub_451CA0",
	"sub_451F30",
	"sub_451F90",
	"sub_451FE0",
	"sub_452010",
	"sub_452050",
	"sub_452120",
	"sub_452190",
	"sub_4521A0",
	"sub_4521F0",
	"sub_452230",
	"nox_xxx_draw_452300",
	"sub_4523D0",
	"sub_452410",
	"sub_452490",
	"sub_452510",
	"sub_452690",
	"sub_4526D0",
	"sub_4526F0",
	"sub_452810",
	"nox_xxx_clientPlaySoundSpecial_452D80",
	"sub_452DC0",
	"sub_452E10",
	"sub_452E90",
	"sub_452EB0",
	"sub_452EE0",
	"sub_452F10",
	"sub_452F50",
	"sub_452F80",
	"sub_452FA0",
	"sub_452FE0",
	"sub_453050",
	"nox_xxx____setargv_9_453060",
	"sub_453070",
	"sub_451CF0",
	"sub_451DC0",
	"sub_451E80",
	"sub_452580",
	"sub_452770",
}
var portTestAudioEventGlobals = []string{
	"dword_587000_122848",
	"dword_587000_126996",
	"dword_587000_127004",
	"dword_587000_93156",
	"dword_5d4594_1045420",
	"dword_5d4594_1045424",
	"dword_5d4594_1045428",
	"dword_5d4594_1045432",
	"dword_5d4594_1045436",
	"dword_5d4594_816368",
	"dword_5d4594_816372",
	"dword_5d4594_816376",
	"dword_5d4594_831092",
}

func portTestAudioEventGlobal(i int) *uint32 {
	switch i {
	case 0:
		return (*uint32)(unsafe.Pointer(&dword_587000_122848))
	case 1:
		return (*uint32)(unsafe.Pointer(&dword_587000_126996))
	case 3:
		return (*uint32)(unsafe.Pointer(&dword_587000_93156))
	case 4:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1045420))
	case 6:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1045428))
	case 7:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_1045432))
	case 9:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_816368))
	case 10:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_816372))
	case 11:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_816376))
	case 12:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_831092))
	case 2:
		return (*uint32)(unsafe.Pointer(&legacyGlobals.dword_587000_127004))
	case 5, 8:
		return nil
	default:
		return nil
	}
}

func PortTestAudioEventCall(name string, args ...uint64) uint64 {
	if len(args) > 4 {
		panic("too many audio event arguments")
	}
	var a [4]uint64
	copy(a[:], args)
	switch name {
	case portTestAudioEventNames[2]:
		return uint64(int64(int32(*audioEventMusicCount)))
	case portTestAudioEventNames[3]:
		return uint64(int64(audioEventMusicSetCount(int32(a[0]))))
	case portTestAudioEventNames[4]:
		return uint64(uintptr(audioEventMusicSlotPointer(int32(a[0]))))
	case portTestAudioEventNames[16]:
		return uint64(audioEventByte())
	case portTestAudioEventNames[17]:
		return uint64(int64(audioEventSetByte(int8(uint8(a[0])))))
	case portTestAudioEventNames[40]:
		return uint64(int64(audioEventVoiceStopped((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a[0])))))))
	case portTestAudioEventNames[41]:
		return uint64(int64(audioEventVoiceEnded((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a[0])))))))
	case portTestAudioEventNames[61]:
		return uint64(int64(audioEventVoiceLoop((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a[0])))))))
	case "nox_xxx_draw_452300":
		return uint64(uintptr(unsafe.Pointer(audioEventNew((*audioEventMetadata)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_4523D0":
		return uint64(int64(int32(audioEventDelete((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_452E90":
		return uint64(int64(int32(audioEventHandleSet((*audioEventHandle)(unsafe.Pointer(uintptr(uint32(a[0])))), (*audioEvent)(unsafe.Pointer(uintptr(uint32(a[1]))))))))
	case "sub_452EB0":
		p := audioEventHandleGet((*audioEventHandle)(unsafe.Pointer(uintptr(uint32(a[0])))))
		return uint64(int64(int32(uintptr(unsafe.Pointer(p)))))
	case "sub_452EE0":
		return uint64(int64(audioEventSetVolume((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0])))), int32(a[1]))))
	case "sub_452F50":
		return uint64(int64(int32(audioEventFadeVolume((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0])))), int32(a[1])))))
	case "sub_452F80":
		return uint64(uintptr(audioEventSetPan((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0])))), int32(a[1]))))
	case "sub_452FE0":
		return uint64(int64(int32(audioEventFadePan((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0])))), int32(a[1])))))
	case "sub_43DA80":
		return uint64(int64(int32(audioEventMusicSave())))
	case "sub_43DAD0":
		audioEventMusicRestore()
		return 0
	case "nox_xxx_clientPlaySoundSpecial_452D80":
		audioEventPlay(int32(a[0]), int32(a[1]), 0, 0)
		return 0
	case "sub_452DC0":
		audioEventPlay(int32(a[0]), int32(a[1]), int32(a[2]), 1)
		return 0
	case "sub_452E10":
		audioEventPlay(int32(a[0]), int32(a[1]), int32(a[2]), 2)
		return 0
	case "sub_43DB60":
		return uint64(int64(int32(audioEventMusicEnter())))
	case "sub_43DBA0":
		audioEventMusicLeave()
		return 0
	case "sub_43DC00":
		audioEventMusicDisable()
		return 0
	case "sub_43DC10":
		return uint64(int64(int32(audioEventMusicEnable())))
	case "sub_43DC30":
		return uint64(int64(int32(audioEventMusicEnabled())))
	case "sub_43EDB0":
		audioEventSampleEnded(ail.Sample(unsafe.Pointer(uintptr(uint32(a[0])))))
		return 0
	case "sub_43EE00":
		return uint64(int64(int32(audioEventSampleRefill((*AudioSample)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_43F0E0":
		return uint64(int64(int32(audioEventSampleFormat((*audioStreamFormat)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_44D960":
		audioEventDialogDisable()
		return 0
	case "sub_44D970":
		return uint64(int64(int32(audioEventDialogEnable())))
	case "sub_44D990":
		return uint64(int64(int32(audioEventDialogEnabled())))
	case "sub_451850":
		return uint64(int64(int32(audioEventInit((*audioStreamContext)(unsafe.Pointer(uintptr(uint32(a[0])))), (*audioStreamCatalog)(unsafe.Pointer(uintptr(uint32(a[1]))))))))
	case "sub_451920":
		return uint64(int64(int32(audioEventDefaults((*audioEventMetadata)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_451970":
		audioEventFree()
		return 0
	case "sub_4519C0":
		audioEventUpdate()
		return 0
	case "sub_451BE0":
		return uint64(int64(int32(audioEventRank((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_451CA0":
		return uint64(int64(int32(uintptr(unsafe.Pointer(audioEventSelectFirst((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))))
	case "sub_451F30":
		return uint64(int64(int32(audioEventLoadSample((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0])))), int32(a[1])))))
	case "sub_451F90":
		return uint64(int64(int32(audioEventReleaseSamples((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_451FE0":
		return uint64(int64(int32(audioEventUnlink((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_452010":
		return uint64(int64(int32(audioEventBucketsClear())))
	case "sub_452050":
		audioEventBucketAdd((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0])))))
		return 0
	case "sub_452120":
		return uint64(audioEventEvict((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))
	case "sub_452190":
		audioEventBucketRemove((*audioEventMetadata)(unsafe.Pointer(uintptr(uint32(a[0])))))
		return 0
	case "sub_4521A0":
		return uint64(uintptr(unsafe.Pointer(audioEventBucketFirst(int32(a[0])))))
	case "sub_4521F0":
		return uint64(int64(int32(audioEventStopAll())))
	case "sub_452230":
		return uint64(audioEventReclaim())
	case "sub_452410":
		return uint64(int64(int32(audioEventDetach((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_452490":
		return uint64(int64(int32(audioEventStart((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_452510":
		audioEventStep((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0])))))
		return 0
	case "sub_452690":
		return uint64(audioEventSchedule((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0])))), uint64(a[1]), uint32(a[2])))
	case "sub_452810":
		return uint64(uintptr(unsafe.Pointer(audioEventChooseVoice(int32(a[0]), int8(int32(a[1]))))))
	case "sub_452F10":
		return uint64(int64(uint32(audioEventVolume((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0])))), int32(a[1])))))
	case "sub_452FA0":
		return uint64(int64(int32(audioEventPan(int32(a[0])))))
	case "sub_453050":
		*audioEventPlayback = 0
		return 0
	case "nox_xxx____setargv_9_453060":
		*audioEventPlayback = 1
		return 0
	case "sub_453070":
		return uint64(int64(int32(int32(*audioEventPlayback))))
	case "sub_451CF0":
		return uint64(int64(int32(uintptr(unsafe.Pointer(audioEventSelectNext((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))))
	case "sub_451DC0":
		return uint64(int64(int32(audioEventLoad((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_451E80":
		return uint64(int64(int32(audioEventReloadIndex((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	case "sub_452580":
		return uint64(int64(int32(audioEventReserve((*audioEvent)(unsafe.Pointer(uintptr(uint32(a[0]))))))))
	}
	for _, n := range portTestAudioEventNames {
		if name == n {
			return 0
		}
	}
	panic("unknown audio event operation: " + name)
}
func PortTestAudioEventGlobals() (map[string]*uint32, func()) {
	words := make(map[string]*uint32)
	saved := make(map[string]uint32)
	for i, n := range portTestAudioEventGlobals {
		p := portTestAudioEventGlobal(i)
		if n == "dword_5d4594_1045424" {
			p = audioEventCacheWord
		}
		if n == "dword_5d4594_1045436" {
			p = audioEventPoolWord
		}
		words[n] = p
		saved[n] = *p
		*p = 0
	}
	return words, func() {
		for n, p := range words {
			*p = saved[n]
		}
	}
}
