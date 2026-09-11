//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME4_1.h"
extern uint32_t nox_tile_def_cnt;
extern nox_tileDef_t nox_tile_defs_arr[176];
extern uint32_t dword_5d4594_3835348;
static void* portTestTileDefs(void) { return nox_tile_defs_arr; }
static uint32_t* portTestTileCount(void) { return &nox_tile_def_cnt; }
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

const portTestTileCount = 176

type PortTestTileNameSpec struct {
	Name []byte // no trailing NUL required
}
type PortTestTileNameResult struct {
	Return                                   int
	Selected, Variation, Flag, Count         uint32
	TableUnchanged, InputUnchanged, GuardsOK bool
}

type PortTestTileScalarResult struct {
	Return                           int
	Selected, Variation, Flag, Count uint32
	TableUnchanged, GuardsOK         bool
}

type PortTestTileVariationSpec struct {
	Selected      uint32 // caller must provide 0..175
	Width, Height byte
	Value         int32
}
type PortTestTileVariationResult struct {
	Return                           int
	Selected, Variation, Flag, Count uint32
	TableUnchanged, GuardsOK         bool
}

type portTestTileState struct{ selected, variation, flag, beforeSelected, afterFlag uint32 }

func portTestTilePtrs() ([]server.TileDef, *uint32, *uint32, *uint32, *uint32) {
	table := unsafe.Slice((*server.TileDef)(C.portTestTileDefs()), portTestTileCount)
	count := (*uint32)(unsafe.Pointer(C.portTestTileCount()))
	selected := memmap.PtrUint32(0x973F18, 35912)
	flag := memmap.PtrUint32(0x973F18, 35916)
	variation := (*uint32)(unsafe.Pointer(&C.dword_5d4594_3835348))
	return table, count, selected, variation, flag
}
func portTestTileGet(selected, variation, flag *uint32) portTestTileState {
	return portTestTileState{selected: *selected, variation: *variation, flag: *flag, beforeSelected: *memmap.PtrUint32(0x973F18, 35908), afterFlag: *memmap.PtrUint32(0x973F18, 35920)}
}
func portTestTileResult(ret C.int, st portTestTileState, same bool, count uint32) PortTestTileScalarResult {
	return PortTestTileScalarResult{Return: int(ret), Selected: st.selected, Variation: st.variation, Flag: st.flag, Count: count, TableUnchanged: same, GuardsOK: st.beforeSelected == 0xa5a5a5a5 && st.afterFlag == 0x5a5a5a5a}
}

// PortTestTileNames installs a full physical 176-entry table regardless of
// nox_tile_def_cnt, then invokes C name lookup. The fixed entries exercise
// duplicate-last and raw high-byte equality without depending on assets.
func PortTestTileNames(specs []PortTestTileNameSpec) (out []PortTestTileNameResult) {
	table, count, selected, variation, flag := portTestTilePtrs()
	oldTable := append([]byte(nil), tileBytes(table)...)
	oldCount := *count
	old := portTestTileGet(selected, variation, flag)
	defer func() {
		copy(tileBytes(table), oldTable)
		*count = oldCount
		*selected, *variation, *flag = old.selected, old.variation, old.flag
		*memmap.PtrUint32(0x973F18, 35908), *memmap.PtrUint32(0x973F18, 35920) = old.beforeSelected, old.afterFlag
	}()
	for i := range table {
		table[i] = server.TileDef{}
	}
	copy(table[3].NameBuf[:], []byte("Alpha\x00"))
	copy(table[17].NameBuf[:], []byte("aLPHa\x00")) // last wins
	copy(table[9].NameBuf[:], []byte("NONE\x00"))   // NONE must override this match
	copy(table[41].NameBuf[:], []byte{0x80, 'X', 0})
	copy(table[42].NameBuf[:], []byte{0x81, 'X', 0})
	copy(table[93].NameBuf[:], "1234567890123456789012345678901")
	*count = 1 // C still scans all 176 physical entries.
	*selected, *variation, *flag = 99, 0x11223344, 1
	*memmap.PtrUint32(0x973F18, 35908), *memmap.PtrUint32(0x973F18, 35920) = 0xa5a5a5a5, 0x5a5a5a5a
	before := append([]byte(nil), tileBytes(table)...)
	for _, spec := range specs {
		input := append(append([]byte(nil), spec.Name...), 0)
		inputBefore := append([]byte(nil), input...)
		ret := C.nox_xxx_tileGetDefByName_51D4D0((*C.char)(unsafe.Pointer(unsafe.SliceData(input))))
		st := portTestTileGet(selected, variation, flag)
		out = append(out, PortTestTileNameResult{Return: int(ret), Selected: st.selected, Variation: st.variation, Flag: st.flag, Count: *count, TableUnchanged: bytes.Equal(tileBytes(table), before), InputUnchanged: bytes.Equal(input, inputBefore), GuardsOK: st.beforeSelected == 0xa5a5a5a5 && st.afterFlag == 0x5a5a5a5a})
	}
	return out
}

func tileBytes(v []server.TileDef) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(v))), len(v)*int(unsafe.Sizeof(server.TileDef{})))
}

// PortTestTileScalars exercises both image index validation and flag setter;
// mode false is image, true is flag. It retains unrelated selection state.
func PortTestTileScalars(values []int32, flagMode bool) (out []PortTestTileScalarResult) {
	table, count, selected, variation, flag := portTestTilePtrs()
	oldTable := append([]byte(nil), tileBytes(table)...)
	oldCount := *count
	old := portTestTileGet(selected, variation, flag)
	defer func() {
		copy(tileBytes(table), oldTable)
		*count = oldCount
		*selected, *variation, *flag = old.selected, old.variation, old.flag
		*memmap.PtrUint32(0x973F18, 35908), *memmap.PtrUint32(0x973F18, 35920) = old.beforeSelected, old.afterFlag
	}()
	before := append([]byte(nil), tileBytes(table)...)
	*count = 7
	*memmap.PtrUint32(0x973F18, 35908), *memmap.PtrUint32(0x973F18, 35920) = 0xa5a5a5a5, 0x5a5a5a5a
	for _, v := range values {
		*selected, *variation, *flag = 99, 0x11223344, 0xa5a5a5a5
		var ret C.int
		if flagMode {
			ret = C.nox_xxx_tile_51D5C0(C.int(v))
		} else {
			ret = C.nox_xxx_tileCheckImage_51D540(C.int(v))
		}
		out = append(out, portTestTileResult(ret, portTestTileGet(selected, variation, flag), bytes.Equal(tileBytes(table), before), *count))
	}
	return out
}

// PortTestTileVariations requires valid selected indices; 255/NONE and other
// invalid indices are deliberately excluded because original C indexes its
// 176-entry array without a guard.
func PortTestTileVariations(specs []PortTestTileVariationSpec) (out []PortTestTileVariationResult) {
	table, count, selected, variation, flag := portTestTilePtrs()
	oldTable := append([]byte(nil), tileBytes(table)...)
	oldCount := *count
	old := portTestTileGet(selected, variation, flag)
	defer func() {
		copy(tileBytes(table), oldTable)
		*count = oldCount
		*selected, *variation, *flag = old.selected, old.variation, old.flag
		*memmap.PtrUint32(0x973F18, 35908), *memmap.PtrUint32(0x973F18, 35920) = old.beforeSelected, old.afterFlag
	}()
	*count = 7
	*variation, *flag = 0x11223344, 1
	*memmap.PtrUint32(0x973F18, 35908), *memmap.PtrUint32(0x973F18, 35920) = 0xa5a5a5a5, 0x5a5a5a5a
	expected := append([]byte(nil), tileBytes(table)...)
	for _, s := range specs {
		if s.Selected >= portTestTileCount {
			panic("invalid selected tile")
		}
		table[s.Selected].Field52, table[s.Selected].Field53 = s.Width, s.Height
		*selected = s.Selected
		expected[int(s.Selected)*60+52], expected[int(s.Selected)*60+53] = s.Width, s.Height
		ret := C.nox_xxx_tileCheckImageVari_51D570(C.int(s.Value))
		st := portTestTileGet(selected, variation, flag)
		out = append(out, PortTestTileVariationResult{Return: int(ret), Selected: st.selected, Variation: st.variation, Flag: st.flag, Count: *count, TableUnchanged: bytes.Equal(tileBytes(table), expected), GuardsOK: st.beforeSelected == 0xa5a5a5a5 && st.afterFlag == 0x5a5a5a5a})
	}
	return out
}

// PortTestTileSelectionState permits independent before/after fixture checks.
func PortTestTileSelectionState() []byte {
	table, count, selected, variation, flag := portTestTilePtrs()
	out := append([]byte(nil), tileBytes(table)...)
	for _, p := range []*uint32{count, selected, variation, flag, memmap.PtrUint32(0x973F18, 35908), memmap.PtrUint32(0x973F18, 35920)} {
		out = binary.LittleEndian.AppendUint32(out, *p)
	}
	return out
}
