//go:build porttest

package legacy

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"unsafe"
)

// Install only the production callback table/name region, with callback identities.
func portTestCallbackTablesEnvironment() func() {
	dst := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 287096)), 640)
	old := bytes.Clone(dst)
	copy(dst, blobdata.PortTestCallbackTables())
	*memmap.PtrPtr(0x587000, 287096) = memmap.PtrOff(0x587000, 287352)
	*memmap.PtrPtr(0x587000, 287100) = monsterCallbackKey(monsterCallbackID_0)
	*memmap.PtrPtr(0x587000, 287104) = memmap.PtrOff(0x587000, 287364)
	*memmap.PtrPtr(0x587000, 287108) = monsterCallbackKey(monsterCallbackID_1)
	*memmap.PtrPtr(0x587000, 287112) = memmap.PtrOff(0x587000, 287380)
	*memmap.PtrPtr(0x587000, 287116) = monsterCallbackKey(monsterCallbackID_2)
	*memmap.PtrPtr(0x587000, 287120) = memmap.PtrOff(0x587000, 287400)
	*memmap.PtrPtr(0x587000, 287124) = monsterCallbackKey(monsterCallbackID_3)
	*memmap.PtrPtr(0x587000, 287128) = memmap.PtrOff(0x587000, 287420)
	*memmap.PtrPtr(0x587000, 287132) = monsterCallbackKey(monsterCallbackID_4)
	*memmap.PtrPtr(0x587000, 287136) = memmap.PtrOff(0x587000, 287436)
	*memmap.PtrPtr(0x587000, 287140) = monsterCallbackKey(monsterCallbackID_5)
	*memmap.PtrPtr(0x587000, 287144) = memmap.PtrOff(0x587000, 287448)
	*memmap.PtrPtr(0x587000, 287148) = monsterCallbackKey(monsterCallbackID_6)
	*memmap.PtrPtr(0x587000, 287152) = memmap.PtrOff(0x587000, 287464)
	*memmap.PtrPtr(0x587000, 287156) = monsterCallbackKey(monsterCallbackID_7)
	*memmap.PtrPtr(0x587000, 287160) = memmap.PtrOff(0x587000, 287488)
	*memmap.PtrPtr(0x587000, 287164) = monsterCallbackKey(monsterCallbackID_8)
	*memmap.PtrPtr(0x587000, 287168) = memmap.PtrOff(0x587000, 287500)
	*memmap.PtrPtr(0x587000, 287172) = monsterCallbackKey(monsterCallbackID_9)
	*memmap.PtrPtr(0x587000, 287176) = memmap.PtrOff(0x587000, 287516)
	*memmap.PtrPtr(0x587000, 287180) = monsterCallbackKey(monsterCallbackID_10)
	*memmap.PtrPtr(0x587000, 287192) = memmap.PtrOff(0x587000, 287532)
	*memmap.PtrPtr(0x587000, 287196) = monsterCallbackKey(monsterCallbackID_11)
	*memmap.PtrPtr(0x587000, 287200) = memmap.PtrOff(0x587000, 287548)
	*memmap.PtrPtr(0x587000, 287204) = monsterCallbackKey(monsterCallbackID_12)
	*memmap.PtrPtr(0x587000, 287208) = memmap.PtrOff(0x587000, 287560)
	*memmap.PtrPtr(0x587000, 287212) = monsterCallbackKey(monsterCallbackID_13)
	*memmap.PtrPtr(0x587000, 287216) = memmap.PtrOff(0x587000, 287568)
	*memmap.PtrPtr(0x587000, 287220) = monsterCallbackKey(monsterCallbackID_14)
	*memmap.PtrPtr(0x587000, 287224) = memmap.PtrOff(0x587000, 287584)
	*memmap.PtrPtr(0x587000, 287228) = monsterCallbackKey(monsterCallbackID_15)
	*memmap.PtrPtr(0x587000, 287232) = memmap.PtrOff(0x587000, 287596)
	*memmap.PtrPtr(0x587000, 287236) = monsterCallbackKey(monsterCallbackID_16)
	*memmap.PtrPtr(0x587000, 287240) = memmap.PtrOff(0x587000, 287608)
	*memmap.PtrPtr(0x587000, 287244) = monsterCallbackKey(monsterCallbackID_17)
	*memmap.PtrPtr(0x587000, 287248) = memmap.PtrOff(0x587000, 287620)
	*memmap.PtrPtr(0x587000, 287252) = monsterCallbackKey(monsterCallbackID_18)
	*memmap.PtrPtr(0x587000, 287256) = memmap.PtrOff(0x587000, 287636)
	*memmap.PtrPtr(0x587000, 287260) = monsterCallbackKey(monsterCallbackID_19)
	*memmap.PtrPtr(0x587000, 287264) = memmap.PtrOff(0x587000, 287656)
	*memmap.PtrPtr(0x587000, 287268) = monsterCallbackKey(monsterCallbackID_20)
	*memmap.PtrPtr(0x587000, 287280) = memmap.PtrOff(0x587000, 287668)
	*memmap.PtrPtr(0x587000, 287284) = monsterCallbackKey(monsterCallbackID_21)
	*memmap.PtrPtr(0x587000, 287288) = memmap.PtrOff(0x587000, 287680)
	*memmap.PtrPtr(0x587000, 287292) = monsterCallbackKey(monsterCallbackID_22)
	*memmap.PtrPtr(0x587000, 287296) = memmap.PtrOff(0x587000, 287688)
	*memmap.PtrPtr(0x587000, 287300) = monsterCallbackKey(monsterCallbackID_23)
	*memmap.PtrPtr(0x587000, 287304) = memmap.PtrOff(0x587000, 287704)
	*memmap.PtrPtr(0x587000, 287308) = monsterCallbackKey(monsterCallbackID_24)
	*memmap.PtrPtr(0x587000, 287312) = memmap.PtrOff(0x587000, 287720)
	*memmap.PtrPtr(0x587000, 287316) = monsterCallbackKey(monsterCallbackID_25)
	return func() { copy(dst, old) }
}
