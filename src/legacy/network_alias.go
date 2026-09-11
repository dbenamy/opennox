package legacy

import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

func resetNetworkAliases(table *[255]server.PlayerNetData) {
	*table = [255]server.PlayerNetData{}
}

func selectNetworkAlias(table *[255]server.PlayerNetData, key1, key2 int32, frame uint32) byte {
	start := int(byte(key1))
	if start == 0 || start == 255 {
		start = 1
	}
	slot := start
	for {
		rec := &table[slot]
		// Stored keys are zero-extended for comparison with the original full ints.
		if (int32(rec.Field0) == key1 && int32(rec.Field2) == key2) || rec.Frame4 < frame {
			return byte(slot)
		}
		slot++
		if slot == 255 {
			slot = 1
		}
		if slot == start {
			return 255
		}
	}
}

//export sub_57B920
func sub_57B920(ptr unsafe.Pointer) C.int {
	resetNetworkAliases((*[255]server.PlayerNetData)(ptr))
	return 0
}

//export nox_xxx_cliGenerateAlias_57B9A0
func nox_xxx_cliGenerateAlias_57B9A0(ptr, key1, key2 C.int, frame C.uint) C.char {
	table := (*[255]server.PlayerNetData)(unsafe.Pointer(uintptr(uint32(ptr))))
	return C.char(selectNetworkAlias(table, int32(key1), int32(key2), uint32(frame)))
}

//export sub_57BA10
func sub_57BA10(ptr C.int, key1, key2 C.short, frame C.int) C.int {
	rec := (*server.PlayerNetData)(unsafe.Pointer(uintptr(uint32(ptr))))
	*rec = server.PlayerNetData{Field0: uint16(key1), Field2: uint16(key2), Frame4: uint32(frame)}
	return ptr
}
