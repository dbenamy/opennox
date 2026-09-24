//go:build porttest

package legacy

/*
#include "GAME5_2.h"
*/
import "C"
import (
	"slices"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/protection"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func PortTestProtectionBit(index, enabled int32) uint32 {
	return protection.Bit(index, enabled)
}

type PortTestBitsetResult struct {
	Award                   int32
	BeforeValid, AfterValid int
	Value, Checksum         uint32
	Unchanged               bool
}

func PortTestBitset(id int32, key, value, sum uint32, index, enabled int32, values []int32, count, mode int) PortTestBitsetResult {
	oldHead, oldKey, oldSum := dword_5d4594_2516344, dword_5d4594_2516348, dword_5d4594_2516328
	defer func() {
		dword_5d4594_2516344, dword_5d4594_2516348, dword_5d4594_2516328 = oldHead, oldKey, oldSum
	}()
	r, free := alloc.New([4]uint32{uint32(id) ^ key, value, 0, 0})
	defer free()
	*r = [4]uint32{uint32(id) ^ key, value, 0, 0}
	if mode == 2 {
		r[0] ^= 1
	} // populated list, missing ID
	head := uint32(uintptr(unsafe.Pointer(r)))
	if mode == 0 {
		head = 0
	}
	dword_5d4594_2516344, dword_5d4594_2516348, dword_5d4594_2516328 = head, uint32(key), uint32(sum)
	var data []int32
	if len(values) != 0 {
		var release func()
		data, release = alloc.CloneSlice(values)
		defer release()
	}
	p := unsafe.Pointer(unsafe.SliceData(data))
	out := PortTestBitsetResult{BeforeValid: int(C.nox_xxx_playerApplyProtectionCRC_56FD50(C.int(id), p, C.int(count)))}
	out.Award = int32(C.nox_xxx_playerAwardSpellProtectionCRC_56FCE0(C.int(id), C.int(index), C.int(enabled)))
	out.AfterValid = int(C.nox_xxx_playerApplyProtectionCRC_56FD50(C.int(id), p, C.int(count)))
	out.Value, out.Checksum = r[1], uint32(dword_5d4594_2516328)
	wantID := uint32(id) ^ key
	if mode == 2 {
		wantID ^= 1
	}
	out.Unchanged = slices.Equal(data, values) && r[0] == wantID && r[2] == 0 && r[3] == 0 && dword_5d4594_2516344 == head && dword_5d4594_2516348 == uint32(key)
	return out
}
