package legacy

/*
#include <stdint.h>
#include "compat.h"
#include "GAME2_2.h"
#include "client/audio/ail/compat_mss.h"
#include "client__io__win95__focus.h"

int sub_43F060(uint32_t* a1);







*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/client/audio/ail"
)

var (
	Sub_43EA20 func(a1 unsafe.Pointer) int
	Sub_43E9F0 func()
	Sub_43E940 func(a1 unsafe.Pointer) int
	Sub_43EFD0 func(a1 unsafe.Pointer) int
	Sub_43EC10 func() int
	Sub_43F130 func() ail.Driver
	Sub_43ED00 func(a1p unsafe.Pointer) int
	Sub_43D6A0 func()
	Sub_44D640 func()
	Sub_44D7E0 func(a1 int) int
	Sub_44D660 func(a1 string) bool
	Sub_43F060 func(a1p unsafe.Pointer) int
	Sub_43EC30 func(a1 unsafe.Pointer) int
	Sub_43ECB0 func(a1 unsafe.Pointer) int
)

var _ = [1]struct{}{}[32-unsafe.Sizeof(AudioSample{})]

type AudioSample struct {
	Dev    ail.Driver     // 0, 0
	Field1 unsafe.Pointer // 1, 4
	Smp    ail.Sample     // 2, 8
	Field3 uint32         // 3, 12
	Field4 unsafe.Pointer // 4, 16
	Data1  *byte          // 5, 20
	Data2  *byte          // 6, 24
	Flag7  uint32         // 7, 28
}

//export sub_43F050
func sub_43F050() int {
	return 0
}

//export sub_43F0D0
func sub_43F0D0() int {
	return 0
}

//export sub_43F030
func sub_43F030(a1 int) int {
	panic("abort")
}

//export AIL_set_stream_volume
func AIL_set_stream_volume(s C.HSTREAM, volume C.int32_t) {
	ail.Stream(unsafe.Pointer(s)).SetVolume(int(volume))
}

//export AIL_stream_position
func AIL_stream_position(s C.HSTREAM) C.int32_t {
	return C.int32_t(ail.Stream(unsafe.Pointer(s)).Position())
}

//export AIL_load_sample_buffer
func AIL_load_sample_buffer(s C.HSAMPLE, num C.uint32_t, buf unsafe.Pointer, sz C.uint32_t) {
	ail.Sample(unsafe.Pointer(s)).LoadBuffer(uint32(num), unsafe.Slice((*byte)(buf), int(sz)))
}

//export AIL_sample_buffer_ready
func AIL_sample_buffer_ready(s C.HSAMPLE) C.int32_t {
	return C.int32_t(ail.Sample(unsafe.Pointer(s)).BufferReady())
}

//export AIL_sample_user_data
func AIL_sample_user_data(s C.HSAMPLE) unsafe.Pointer {
	v := ail.Sample(unsafe.Pointer(s)).UserData()
	if v == nil {
		return nil
	}
	return unsafe.Pointer(v.(*AudioSample))
}

//export sub_43F010
func sub_43F010(a1 unsafe.Pointer) int {
	p := *(**AudioSample)(unsafe.Add(a1, 272))
	p.Smp.Stop()
	return 0
}

//export sub_43EA20
func sub_43EA20(a1 unsafe.Pointer) int {
	return Sub_43EA20(a1)
}

//export sub_43E9F0
func sub_43E9F0() {
	Sub_43E9F0()
}

//export sub_43E940
func sub_43E940(a1 unsafe.Pointer) int {
	return Sub_43E940(a1)
}

//export sub_43EFD0
func sub_43EFD0(a1 unsafe.Pointer) int {
	return Sub_43EFD0(a1)
}

//export sub_43EC10
func sub_43EC10() int {
	return Sub_43EC10()
}

//export sub_43F130
func sub_43F130() int {
	return int(Sub_43F130())
}

//export sub_43ED00
func sub_43ED00(a1p *C.uint32_t) int {
	return Sub_43ED00(unsafe.Pointer(a1p))
}

//export sub_44D640
func sub_44D640() {
	Sub_44D640()
}

//export sub_44D7E0
func sub_44D7E0(a1 int) int {
	return Sub_44D7E0(a1)
}

//export sub_43F060
func sub_43F060(a1p *C.uint32_t) int {
	return Sub_43F060(unsafe.Pointer(a1p))
}

//export sub_43EC30
func sub_43EC30(a1p unsafe.Pointer) int {
	return Sub_43EC30(unsafe.Pointer(a1p))
}

//export sub_43ECB0
func sub_43ECB0(a1p unsafe.Pointer) int {
	return Sub_43ECB0(unsafe.Pointer(a1p))
}

func Get_dword_587000_127004() unsafe.Pointer {
	return legacyGlobals.dword_587000_127004
}

func Sub_43F0E0(v unsafe.Pointer) int {
	return int(audioEventSampleFormat((*audioStreamFormat)(v)))
}

func Sub_43EE00(v *AudioSample) {
	audioEventSampleRefill(v)
}

func Sub_43EDB0(v ail.Sample) {
	audioEventSampleEnded(v)
}

func Sub_413890() string {
	return ""
}

func Nox_xxx_parseSoundSetBin_424170(path string) int {
	return resourceSoundLoad(path)
}

func Sub_43DC00() {
	audioEventMusicDisable()
}

func Sub_44D960() {
	audioEventDialogDisable()
}

func Sub_453050() {
	*audioEventPlayback = 0
}

func Get_dword_5d4594_816376() ail.Driver {
	return ail.Driver(dword_5d4594_816376)
}

func Set_dword_5d4594_816376(drv ail.Driver) {
	dword_5d4594_816376 = C.uint(drv)
}

func Sub_486640(a1 unsafe.Pointer, a2 int) int {
	return int(int32(audioStreamScaleVolume(a1, uint32(a2))))
}

func Get_dword_5d4594_805984() unsafe.Pointer {
	return legacyGlobals.dword_5d4594_805984
}

func Set_dword_5d4594_805984(v unsafe.Pointer) {
	legacyGlobals.dword_5d4594_805984 = v
}

func Set_dword_587000_81128(v unsafe.Pointer) {
	legacyGlobals.dword_587000_81128 = v
}

func Sub_451850(a1 unsafe.Pointer, a2 unsafe.Pointer) {
	audioEventInit((*audioStreamContext)(a1), (*audioStreamCatalog)(a2))
}

func Sub_486FA0(a1 int) {
	audioStreamDeviceRegister((*audioStreamDescriptor)(unsafe.Pointer(uintptr(a1))))
}

func Sub_487D00(a1 unsafe.Pointer) {
	audioStreamByteRate((*audioStreamFormat)(a1))
}

func Sub_487150(a1 int, a2 unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(audioStreamContextAcquire(int32(a1), (*audioStreamFormat)(a2)))
}

func Sub_487790(a1 unsafe.Pointer, a2 int) int {
	return int(audioStreamVoiceCreateMany((*audioStreamContext)(a1), int32(a2)))
}
