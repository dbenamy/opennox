//go:build porttest

package legacy

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"unsafe"
)

type portTestEffectsRegion struct {
	data, old, initial []byte
	omitSnapshot       bool
}
type PortTestEffectsEnvironment struct {
	words        []*uint32
	old, initial []uint32
	regions      []portTestEffectsRegion
}

// PortTestNewEffectsEnvironment owns all writable effect state and the bounded
// production tables normally loaded by game startup. Named C words and mapped
// historical addresses are saved separately because they are distinct storage.
func PortTestNewEffectsEnvironment() *PortTestEffectsEnvironment {
	e := &PortTestEffectsEnvironment{words: []*uint32{
		(*uint32)(unsafe.Pointer(&dword_587000_180476)),
		(*uint32)(unsafe.Pointer(&dword_587000_180480)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1304328)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1313532)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1313536)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1313540)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1313564)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1313692)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1313880)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1316408)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1316412)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1316436)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1316448)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1316452)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1316456)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1316472)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1316476)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1316484)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1316492)),
		(*uint32)(unsafe.Pointer(&nox_color_white_2523948)),
		(*uint32)(unsafe.Pointer(&nox_xxx_lightningSteps_587000_178216)),
	}}
	for _, p := range e.words {
		e.old = append(e.old, *p)
		e.initial = append(e.initial, 0)
	}
	e.initial[3] = 0xffff
	e.initial[4] = 0xffff
	e.initial[5] = 0xffff
	e.initial[6] = 0xffff
	e.initial[19] = 0xffff
	e.initial[20] = 8
	for _, r := range []struct {
		base, off uintptr
		size      int
	}{
		{0x5D4594, 1303540, 812}, {0x5D4594, 1313528, 180},
		{0x5D4594, 1313828, 2708}, {0x5D4594, 1316980, 20},
	} {
		data := unsafe.Slice((*byte)(memmap.PtrOff(r.base, r.off)), r.size)
		e.regions = append(e.regions, portTestEffectsRegion{data: data, old: append([]byte(nil), data...), initial: make([]byte, len(data))})
	}
	for _, r := range blobdata.PortTestClientEffectsTables() {
		data := unsafe.Slice((*byte)(memmap.PtrOff(r.Base, r.Offset)), len(r.Data))
		e.regions = append(e.regions, portTestEffectsRegion{data: data, old: append([]byte(nil), data...), initial: r.Data, omitSnapshot: r.Base == 0x587000 && r.Offset == 155956})
	}
	e.Reset()
	// Every base color starts white; color-specific cases can replace these words.
	for off := uintptr(1313528); off <= 1313592; off += 4 {
		*memmap.PtrUint32(0x5D4594, off) = 0xffff
	}
	copy(e.regions[1].initial, e.regions[1].data)
	return e
}
func (e *PortTestEffectsEnvironment) Reset() {
	for i, p := range e.words {
		*p = e.initial[i]
	}
	for _, r := range e.regions {
		copy(r.data, r.initial)
	}
}
func (e *PortTestEffectsEnvironment) Restore() {
	for i, p := range e.words {
		*p = e.old[i]
	}
	for _, r := range e.regions {
		copy(r.data, r.old)
	}
}

// Snapshot returns exact non-pointer scratch state. Ray-cache pointers require
// normalization by the drawable owner and are deliberately captured separately.
func (e *PortTestEffectsEnvironment) Snapshot() []uint32 {
	out := make([]uint32, 0)
	for _, p := range e.words {
		out = append(out, *p)
	}
	for _, r := range e.regions[1:] {
		if r.omitSnapshot {
			if !bytes.Equal(r.data, r.initial) {
				panic("readonly effect distance table changed")
			}
			continue
		}
		out = append(out, unsafe.Slice((*uint32)(unsafe.Pointer(&r.data[0])), len(r.data)/4)...)
	}
	return out
}

func (e *PortTestEffectsEnvironment) Named(name string) *uint32 {
	switch name {
	case "dword_587000_180476":
		return e.words[0]
	case "dword_587000_180480":
		return e.words[1]
	case "dword_5d4594_1304328":
		return e.words[2]
	case "dword_5d4594_1313532":
		return e.words[3]
	case "dword_5d4594_1313536":
		return e.words[4]
	case "dword_5d4594_1313540":
		return e.words[5]
	case "dword_5d4594_1313564":
		return e.words[6]
	case "dword_5d4594_1313692":
		return e.words[7]
	case "dword_5d4594_1313880":
		return e.words[8]
	case "dword_5d4594_1316408":
		return e.words[9]
	case "dword_5d4594_1316412":
		return e.words[10]
	case "dword_5d4594_1316436":
		return e.words[11]
	case "dword_5d4594_1316448":
		return e.words[12]
	case "dword_5d4594_1316452":
		return e.words[13]
	case "dword_5d4594_1316456":
		return e.words[14]
	case "dword_5d4594_1316472":
		return e.words[15]
	case "dword_5d4594_1316476":
		return e.words[16]
	case "dword_5d4594_1316484":
		return e.words[17]
	case "dword_5d4594_1316492":
		return e.words[18]
	case "nox_color_white_2523948":
		return e.words[19]
	case "nox_xxx_lightningSteps_587000_178216":
		return e.words[20]
	default:
		panic("unknown named effect word")
	}
}
