//go:build porttest

package legacy

/*
#include "GAME3_2.h"
#include "GAME3_3.h"
extern uint32_t dword_5d4594_1556128;
extern uint32_t dword_5d4594_1556144;
*/
import "C"

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

// Dispatch through retained C interfaces or directly into Go. Fixture pointer returns are
// normalized independently from scalar returns by the caller-selected operation.
func PortTestQuestRuntime(op string, u *server.Object, args [4]uint32) uint64 {
	ai := C.int(uintptr(unsafe.Pointer(u)))
	x := C.int(args[0])
	var rv uint32
	pointerResult := false
	switch op {
	case "sub_4D6000":
		rv = uint32(C.sub_4D6000(asObjectC(u)))
		pointerResult = true
	case "sub_4D60E0":
		rv = uint32(uintptr(unsafe.Pointer(C.sub_4D60E0(ai))))
		pointerResult = true
	case "sub_4D6130":
		rv = uint32(C.sub_4D6130(ai))
		pointerResult = true
	case "sub_4D6170":
		rv = questRuntimeIncrement(u, 4664, 4)
		pointerResult = true
	case "sub_4D61F0":
		rv = uint32(C.sub_4D61F0(ai))
		pointerResult = true
	case "sub_4D61B0":
		questRuntimeIncrement(u, 4668, 8)
	case "sub_4D60B0":
		rv = uint32(C.sub_4D60B0())
	case "nox_xxx_isQuest_4D6F50":
		rv = uint32(C.nox_xxx_isQuest_4D6F50())
	case "sub_4D6F70":
		rv = uint32(C.sub_4D6F70())
	case "sub_4D6FA0":
		rv = uint32(C.sub_4D6FA0())
	case "sub_4D7150":
		rv = uint32(questRuntimeObserverDeadline())
	case "sub_4D71F0":
		rv = uint32(questRuntimeSoulTimeout())
	case "sub_4D7300":
		rv = uint32(questRuntimeWord(1556132))
	case "sub_4D7430":
		rv = uint32(questRuntimeWord(1556116))
	case "sub_4D75E0":
		rv = uint32(C.sub_4D75E0())
	case "sub_4D76F0":
		rv = uint32(questRuntimeWord(1556124))
	case "sub_4D7A80":
		rv = uint32(questRuntimeDepartureTick())
	case "sub_4D7B40":
		rv = uint32(questRuntimeDepartureReset())
	case "nox_game_getQuestStage_4E3CC0":
		rv = uint32(C.nox_game_getQuestStage_4E3CC0())
	case "nox_xxx_player_4E3CE0":
		rv = uint32(C.nox_xxx_player_4E3CE0())
	case "sub_4E3D50":
		rv = uint32(C.sub_4E3D50())
	case "sub_4E4100":
		rv = uint32(bool2int(questRuntimeRoom()))
	case "sub_4D6540":
		rv = uint32(C.sub_4D6540(x))
	case "sub_4D6770":
		rv = uint32(questRuntimeScoreboard(int(x)))
	case "nox_game_sendQuestStage_4D6960":
		rv = uint32(questRuntimeStageMessage(int(x), 14, 0))
	case "nox_xxx_setQuest_4D6F60":
		rv = uint32(C.nox_xxx_setQuest_4D6F60(x))
	case "sub_4D6F80":
		rv = uint32(C.sub_4D6F80(x))
	case "nox_xxx_bookCreatureTest_4D70C0":
		rv = uint32(bool2int(questRuntimeBookAllowed(int(x), 37)))
	case "sub_4D7100":
		rv = uint32(bool2int(questRuntimeBookAllowed(int(x), 111)))
	case "sub_4D71E0":
		rv = uint32(questRuntimeSetSoulFrame(uint32(x)))
	case "sub_4D72D0":
		rv = uint32(C.sub_4D72D0(x))
	case "sub_4D7440":
		rv = uint32(questRuntimeSetWord(1556116, uint32(x)))
	case "sub_4D7520":
		rv = uint32(questRuntimeGateSet(uint32(x)))
	case "sub_4D75F0":
		rv = uint32(questRuntimeSetWord(1556108, uint32(x)))
	case "sub_4D76E0":
		rv = uint32(C.sub_4D76E0(x))
	case "sub_4D7A60":
		rv = uint32(C.sub_4D7A60(x))
	case "sub_4D66E0":
		rv = uint32(questRuntimeScore(args[0], args[1], args[2], args[3]))
	case "sub_4D6880":
		rv = uint32(questRuntimeStageMessage(int(x), 13, args[1]))
	case "sub_4D6A20":
		rv = uint32(questRuntimeCodeMessage(int(x), u))
	case "sub_4D7280":
		rv = uint32(questRuntimeSharedMessage(int(x), byte(args[1])))
	case "sub_4D7450":
		rv = uint32(C.sub_4D7450(x, C.short(args[1])))
	case "sub_4D79A0":
		rv = uint32(questRuntimeSlotMask(args[0]))
	case "sub_4D79C0":
		rv = uint32(questRuntimeReconnect(u))
	case "sub_4D7480":
		questRuntimeGateReturn(u)
	case "nox_server_checkWarpGate_4D7600":
		questRuntimeWarpTick()
	case "nox_game_setQuestStage_4E3CD0":
		C.nox_game_setQuestStage_4E3CD0(x)
	case "sub_4E3DD0":
		C.sub_4E3DD0() // live Go caller ignores decompiler return
	case "sub_4E3CB0":
		rv = uint32(questRuntimeSetDifficulty(math.Float32frombits(args[0])))
	case "sub_4E4080":
		questRuntimeCap(202032, "PlayerDamageCap", math.Float32frombits(args[0]))
	case "sub_4E40C0":
		questRuntimeCap(202036, "SystemHealthCap", math.Float32frombits(args[0]))
	case "sub_4E3CA0":
		return math.Float64bits(float64(C.sub_4E3CA0()))
	case "sub_4E40B0":
		return math.Float64bits(float64(questRuntimeFloat(202032)))
	case "sub_4E40F0":
		return math.Float64bits(float64(C.sub_4E40F0()))
	default:
		panic("unknown quest runtime operation: " + op)
	}
	if pointerResult && u != nil {
		if rv == uint32(ai) {
			return 1
		}
		if rv == uint32(uintptr(controlPlayer(u))) {
			return 2
		}
	}
	return uint64(rv)
}
func PortTestQuestRuntimeString(op string) string {
	switch op {
	case "sub_4D6940":
		return GoStringP(unsafe.Pointer(questRuntimeMapName(3838)))
	case "sub_4D6950":
		return GoStringP(unsafe.Pointer(questRuntimeMapName(3806)))
	default:
		panic(op)
	}
}
func PortTestQuestRuntimeGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"previousStage":    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1556128)),
		"observerDeadline": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1556144)),
	}
	for _, off := range []uintptr{1556104, 1556108, 1556116, 1556120, 1556124, 1556132, 1556160, 1556164, 1556300, 1563928, 1563932, 1563912, 1563908, 1563916, 1563920, 1563924} {
		words[fmt.Sprint(off)] = memmap.PtrUint32(0x5D4594, off)
	}
	for _, off := range []uintptr{202024, 202028, 202032, 202036} {
		words[fmt.Sprint(off)] = memmap.PtrUint32(0x587000, off)
	}
	saved := make(map[string]uint32)
	for k, p := range words {
		saved[k] = *p
		*p = 0
	}
	type region struct {
		base, off uintptr
		size      int
	}
	regions := []region{{0x5D4594, 1556172, 128}, {0x973F18, 3838, 32}, {0x973F18, 3806, 32}, {0x5D4594, 1556152, 4}}
	var buffers, old [][]byte
	for _, r := range regions {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r.base, r.off)), r.size)
		buffers = append(buffers, b)
		old = append(old, bytes.Clone(b))
		clear(b)
	}
	return words, func() {
		for k, p := range words {
			*p = saved[k]
		}
		for i, b := range buffers {
			copy(b, old[i])
		}
	}
}

// This char* is a settings record, not a terminated string.
func PortTestQuestRuntimeSettings() unsafe.Pointer { return questRuntimeSettings() }
