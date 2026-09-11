//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME4_3.h"
extern uint32_t dword_5d4594_251572;
extern uint32_t dword_5d4594_2489436;
extern uint32_t dword_5d4594_3835356;
extern uint32_t dword_5d4594_3835360;
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

const portTestBorderRows = 64
const portTestBorderRowSize = 60
const portTestBorderBase = 28644

type PortTestBorderRow struct {
	Index int
	Raw   [portTestBorderRowSize]byte
}

type PortTestBorderState struct {
	Count, Flag, Primary, Secondary uint32
}

type PortTestBorderSpec struct {
	// Mode: 0=543FB0 lookup, 1=544020 select name, 2=544070 primary,
	// 3=5440A0 secondary. NilName is valid only for direct lookup.
	Mode    byte
	Name    []byte // fixture appends the C NUL without modifying it
	NilName bool
	Value   int32
}

type PortTestBorderResult struct {
	Return                         int
	State                          PortTestBorderState
	TableUnchanged, InputUnchanged bool
	GuardsUnchanged                bool
}

type PortTestBorderSnapshot struct {
	Results                       []PortTestBorderResult
	Before, AfterRestore          PortTestBorderState
	TableRestored, GuardsRestored bool
}

func portTestBorderState() PortTestBorderState {
	return PortTestBorderState{Count: uint32(C.dword_5d4594_251572), Flag: uint32(C.dword_5d4594_2489436), Primary: uint32(C.dword_5d4594_3835356), Secondary: uint32(C.dword_5d4594_3835360)}
}
func portTestBorderSetState(v PortTestBorderState) {
	C.dword_5d4594_251572 = C.uint32_t(v.Count)
	C.dword_5d4594_2489436 = C.uint32_t(v.Flag)
	C.dword_5d4594_3835356 = C.uint32_t(v.Primary)
	C.dword_5d4594_3835360 = C.uint32_t(v.Secondary)
}
func portTestBorderTable() []byte {
	return unsafe.Slice(memmap.PtrUint8(0x85B3FC, portTestBorderBase), portTestBorderRows*portTestBorderRowSize)
}

// PortTestBorderSelection invokes the native lookup and three live C ABI entries over a bounded 64-row
// physical table. Count remains raw: tests may supply signed-negative words,
// but deliberately do not call a positive count beyond this physical table.
func PortTestBorderSelection(initial PortTestBorderState, rows []PortTestBorderRow, specs []PortTestBorderSpec) (snap PortTestBorderSnapshot) {
	table := portTestBorderTable()
	left := unsafe.Slice(memmap.PtrUint8(0x85B3FC, portTestBorderBase-8), 8)
	right := unsafe.Slice(memmap.PtrUint8(0x85B3FC, uintptr(portTestBorderBase+len(table))), 8)
	oldTable, oldLeft, oldRight := append([]byte(nil), table...), append([]byte(nil), left...), append([]byte(nil), right...)
	snap.Before = portTestBorderState()
	defer func() {
		copy(table, oldTable)
		copy(left, oldLeft)
		copy(right, oldRight)
		portTestBorderSetState(snap.Before)
		snap.AfterRestore = portTestBorderState()
		snap.TableRestored = bytes.Equal(table, oldTable)
		snap.GuardsRestored = bytes.Equal(left, oldLeft) && bytes.Equal(right, oldRight)
	}()
	for i := range table {
		table[i] = 0
	}
	for i := range left {
		left[i] = byte(0x31 + i)
	}
	for i := range right {
		right[i] = byte(0xc1 + i)
	}
	for _, row := range rows {
		if row.Index < 0 || row.Index >= portTestBorderRows {
			panic("invalid border row")
		}
		copy(table[row.Index*portTestBorderRowSize:], row.Raw[:])
	}
	portTestBorderSetState(initial)
	configured := append([]byte(nil), table...)
	wantLeft, wantRight := append([]byte(nil), left...), append([]byte(nil), right...)
	for _, s := range specs {
		var ret C.int
		input := append(append([]byte(nil), s.Name...), 0)
		inputBefore := append([]byte(nil), input...)
		switch s.Mode {
		case 0:
			if s.NilName {
				ret = C.int(findBorderName(nil))
			} else {
				ret = C.int(findBorderName((*C.char)(unsafe.Pointer(unsafe.SliceData(input)))))
			}
		case 1:
			if s.NilName {
				panic("544020 null input faults in C")
			}
			ret = C.sub_544020((*C.char)(unsafe.Pointer(unsafe.SliceData(input))))
		case 2:
			ret = C.nox_xxx_tileCheckByte3_544070(C.int(s.Value))
		case 3:
			ret = C.nox_xxx_tileCheckByte4_5440A0(C.int(s.Value))
		default:
			panic("invalid border mode")
		}
		snap.Results = append(snap.Results, PortTestBorderResult{Return: int(ret), State: portTestBorderState(), TableUnchanged: bytes.Equal(table, configured), InputUnchanged: bytes.Equal(input, inputBefore), GuardsUnchanged: bytes.Equal(left, wantLeft) && bytes.Equal(right, wantRight)})
	}
	return snap
}

// PortTestBorderRow makes a raw 60-byte record with an exact C string at +0
// and the loader's uint16 border-variation count at +44.
func NewPortTestBorderRow(index int, name []byte, limit uint16) PortTestBorderRow {
	var r PortTestBorderRow
	r.Index = index
	copy(r.Raw[:], name)
	binary.LittleEndian.PutUint16(r.Raw[44:], limit)
	return r
}
