package legacy

/*
#include "GAME2_2.h"
#include "GAME3_1.h"
*/

import "C"

import "unsafe"

//export sub_4873C0
func sub_4873C0(a0 C.int) C.int {
	return C.int(audioStreamContextTick((*audioStreamContext)(unsafe.Pointer(uintptr(uint32(a0))))))
}

//export sub_487810
func sub_487810(a0 C.int, a1 C.int) *C.int {
	return (*C.int)(unsafe.Pointer(audioStreamVoiceSelect((*audioStreamContext)(unsafe.Pointer(uintptr(uint32(a0)))), int32(a1))))
}

//export sub_4BD280
func sub_4BD280(a0 C.int, a1 C.int) *C.uint {
	return (*C.uint)(unsafe.Pointer(audioStreamPoolNew(int32(a0), int32(a1))))
}

//export sub_4BD2D0
func sub_4BD2D0(a0 unsafe.Pointer) { audioStreamPoolFree(unsafe.Pointer(a0)) }

//export sub_4BD2E0
func sub_4BD2E0(a0 **C.uint) *C.uint {
	return (*C.uint)(unsafe.Pointer(audioStreamPoolPop((*unsafe.Pointer)(unsafe.Pointer(a0)))))
}

//export sub_4BD300
func sub_4BD300(a0 *C.uint, a1 C.int) C.int {
	return C.int(uintptr(unsafe.Pointer(audioStreamPoolPush((*unsafe.Pointer)(unsafe.Pointer(a0)), unsafe.Pointer(uintptr(uint32(a1)))))))
}

//export sub_4BD340
func sub_4BD340(a0 C.int, a1 C.int, a2 C.int, a3 C.int) *C.uint {
	return (*C.uint)(unsafe.Pointer(audioStreamCacheNew((*audioStreamCatalog)(unsafe.Pointer(uintptr(uint32(a0)))), int32(a1), int32(a2), int32(a3))))
}

//export sub_4BD3C0
func sub_4BD3C0(a0 unsafe.Pointer) { audioStreamCacheFree((*audioStreamCache)(unsafe.Pointer(a0))) }

//export sub_4BD470
func sub_4BD470(a0 **C.uint, a1 C.int) *C.uint {
	return (*C.uint)(unsafe.Pointer(audioStreamCacheLoad((*audioStreamCache)(unsafe.Pointer(a0)), int32(a1))))
}

//export sub_4BD650
func sub_4BD650(a0 C.int) C.int {
	return C.int(uintptr(unsafe.Pointer(audioStreamCacheRef((*audioStreamCacheEntry)(unsafe.Pointer(uintptr(uint32(a0))))))))
}

//export sub_4BD660
func sub_4BD660(a0 C.int) C.int {
	return C.int(audioStreamCacheUnref((*audioStreamCacheEntry)(unsafe.Pointer(uintptr(uint32(a0))))))
}

//export sub_4BD710
func sub_4BD710(a0 C.int) C.int {
	return C.int(uintptr(unsafe.Pointer(audioStreamCacheBuffer((*audioStreamCacheEntry)(unsafe.Pointer(uintptr(uint32(a0))))))))
}

//export sub_4BD8C0
func sub_4BD8C0(a0 C.int) C.int {
	return C.int(audioStreamVoiceData((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a0))))))
}

//export sub_4BD940
func sub_4BD940(a0 C.int) C.int {
	return C.int(audioStreamVoiceLoop((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a0))))))
}

//export sub_4BD9B0
func sub_4BD9B0(a0 *C.uint) C.int {
	return C.int(audioStreamVoiceEnd((*audioStreamVoice)(unsafe.Pointer(a0))))
}

//export sub_4BDA80
func sub_4BDA80(a0 C.int) C.int {
	return C.int(audioStreamVoiceStop((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a0))))))
}

//export sub_4BDB20
func sub_4BDB20(a0 C.int) C.int {
	return C.int(uintptr(unsafe.Pointer(audioStreamVoiceReserve((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a0))))))))
}

//export sub_4BDB30
func sub_4BDB30(a0 C.int) C.int {
	return C.int(uintptr(unsafe.Pointer(audioStreamVoiceUnreserve((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a0))))))))
}

//export sub_4BDB40
func sub_4BDB40(a0 C.int) C.int {
	return C.int(audioStreamVoiceStart((*audioStreamVoice)(unsafe.Pointer(uintptr(uint32(a0))))))
}

//export sub_4BDB90
func sub_4BDB90(a0 *C.uint, a1 *C.uint) {
	audioStreamVoiceBind((*audioStreamVoice)(unsafe.Pointer(a0)), (*audioStreamBuffer)(unsafe.Pointer(a1)))
}
