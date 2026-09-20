package legacy

/*
#include <stdint.h>
#include "GAME2.h"
#include "GAME1_3.h"
#include "client__audio__audevent.h"
*/
import "C"
import "unsafe"

func sub_43DA80() C.int { return C.int(audioEventMusicSave()) }

func sub_43DAD0() { audioEventMusicRestore() }

//export sub_43DB20
func sub_43DB20() C.int { return C.int(int32(*audioEventMusicCount)) }

//export sub_43DB30
func sub_43DB30(a0 C.int) C.int { return C.int(audioEventMusicSetCount(int32(a0))) }

//export sub_43DB40
func sub_43DB40(a0 C.int) *C.char {
	return (*C.char)(unsafe.Pointer(audioEventMusicSlotPointer(int32(a0))))
}

//export sub_450750
func sub_450750() C.uchar { return C.uchar(audioEventByte()) }

//export sub_450760
func sub_450760(a0 C.char) C.char { return C.char(audioEventSetByte(int8(int32(a0)))) }

//export nox_xxx_draw_452300
func nox_xxx_draw_452300(a0 *C.uint) *C.uint {
	return (*C.uint)(unsafe.Pointer(audioEventNew((*audioEventMetadata)(unsafe.Pointer(a0)))))
}

//export sub_4523D0
func sub_4523D0(a0 unsafe.Pointer) C.int {
	return C.int(audioEventDelete((*audioEvent)(unsafe.Pointer(a0))))
}

//export sub_4526D0
func sub_4526D0(a0 C.int) C.int {
	return C.int(audioEventVoiceStopped((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a0))))))
}

//export sub_4526F0
func sub_4526F0(a0 C.int) C.int {
	return C.int(audioEventVoiceEnded((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a0))))))
}

func nox_xxx_clientPlaySoundSpecial_452D80(a0 C.int, a1 C.int) {
	audioEventPlay(int32(a0), int32(a1), 0, 0)
}

func sub_452DC0(a0 C.int, a1 C.int, a2 C.int) { audioEventPlay(int32(a0), int32(a1), int32(a2), 1) }

func sub_452E10(a0 C.int, a1 C.int, a2 C.int) { audioEventPlay(int32(a0), int32(a1), int32(a2), 2) }

//export sub_452E90
func sub_452E90(a0 *C.uint, a1 C.int) C.int {
	return C.int(audioEventHandleSet((*audioEventHandle)(unsafe.Pointer(a0)), (*audioEvent)(unsafe.Pointer(uintptr(uint32(a1))))))
}

//export sub_452EB0
func sub_452EB0(a0 *C.int) C.int {
	return C.int(uintptr(unsafe.Pointer(audioEventHandleGet((*audioEventHandle)(unsafe.Pointer(a0))))))
}

//export sub_452EE0
func sub_452EE0(a0 C.int, a1 C.int) C.int {
	return C.int(audioEventSetVolume((*audioEvent)(unsafe.Pointer(uintptr(uint32(a0)))), int32(a1)))
}

//export sub_452F50
func sub_452F50(a0 C.int, a1 C.int) C.int {
	return C.int(audioEventFadeVolume((*audioEvent)(unsafe.Pointer(uintptr(uint32(a0)))), int32(a1)))
}

//export sub_452F80
func sub_452F80(a0 C.int, a1 C.int) *C.uint {
	return (*C.uint)(unsafe.Pointer(audioEventSetPan((*audioEvent)(unsafe.Pointer(uintptr(uint32(a0)))), int32(a1))))
}

//export sub_452FE0
func sub_452FE0(a0 C.int, a1 C.int) C.int {
	return C.int(audioEventFadePan((*audioEvent)(unsafe.Pointer(uintptr(uint32(a0)))), int32(a1)))
}

//export sub_452770
func sub_452770(a0 *C.uint) C.int {
	return C.int(audioEventVoiceLoop((*audioStreamVoice)(unsafe.Pointer(a0))))
}
