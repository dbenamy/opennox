//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

type PortTestMeterRecord struct {
	Window                             *gui.Window
	Current, Maximum, Color, Alternate uint32
}
type PortTestMeterEnvironment struct {
	named      []*uint32
	saved      []uint32
	Regions    [][]byte
	previous   [][]byte
	Records    []PortTestMeterRecord
	oldRecords []PortTestMeterRecord
}

func PortTestNewMeterEnvironment() *PortTestMeterEnvironment {
	e := &PortTestMeterEnvironment{named: []*uint32{
		(*uint32)(unsafe.Pointer(&dword_5d4594_1090276)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1090280)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1090284)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1090292)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1090828)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1091364)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1096252)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1096256)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1096260)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1096264)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1096272)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1096276)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1096280)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1096284)),
		(*uint32)(unsafe.Pointer(&dword_5d4594_1096288)),
		(*uint32)(unsafe.Pointer(&dword_8531A0_2576)),
		(*uint32)(unsafe.Pointer(&nox_client_renderBubbles_80844)),
		(*uint32)(unsafe.Pointer(&nox_color_black_2650656)),
		(*uint32)(unsafe.Pointer(&nox_color_violet_2598268)),
		(*uint32)(unsafe.Pointer(&nox_color_white_2523948)),
		(*uint32)(unsafe.Pointer(&nox_color_yellow_2589772)),
		(*uint32)(unsafe.Pointer(&nox_gameDisableMapDraw_5d4594_2650672)),
		(*uint32)(unsafe.Pointer(&nox_player_netCode_85319C)),
		(*uint32)(unsafe.Pointer(&nox_win_width)),
		(*uint32)(unsafe.Pointer(&nox_win_height)),
	}}
	if unsafe.Sizeof(PortTestMeterRecord{}) != 20 {
		panic("meter record ABI")
	}
	for _, p := range e.named {
		e.saved = append(e.saved, *p)
	}
	e.Records = unsafe.Slice((*PortTestMeterRecord)(unsafe.Pointer(&uiMeterRecords[0])), 7)
	e.oldRecords = append([]PortTestMeterRecord(nil), e.Records...)
	// Complete meter/potion storage, charge-raster row state, and cursor text.
	// Mapped embedded potion drawables are not ordinary sprite-pool allocations.
	for _, r := range [][3]uintptr{{0x5D4594, 1090276, 6032}, {0x587000, 147904, 488}, {0x5D4594, 1096672, 516}} {
		p := unsafe.Slice((*byte)(memmap.PtrOff(r[0], r[1])), int(r[2]))
		e.Regions = append(e.Regions, p)
		e.previous = append(e.previous, append([]byte(nil), p...))
	}
	return e
}
func (e *PortTestMeterEnvironment) Restore() {
	for i, p := range e.named {
		*p = e.saved[i]
	}
	copy(e.Records, e.oldRecords)
	for i, p := range e.Regions {
		copy(p, e.previous[i])
	}
}
func (e *PortTestMeterEnvironment) Named() []uint32 {
	out := make([]uint32, len(e.named))
	for i, p := range e.named {
		out[i] = *p
	}
	return out
}
func (e *PortTestMeterEnvironment) SetNamed(i int, v uint32) { *e.named[i] = v }
func (e *PortTestMeterEnvironment) Reset() {
	for i, p := range e.Regions {
		copy(p, e.previous[i])
	}
	for i, p := range e.named {
		*p = e.saved[i]
	}
	clear(e.Records)
}

func (e *PortTestMeterEnvironment) NamedWord(name string) *uint32 {
	switch name {
	case "dword_5d4594_1090276":
		return e.named[0]
	case "dword_5d4594_1090280":
		return e.named[1]
	case "dword_5d4594_1090284":
		return e.named[2]
	case "dword_5d4594_1090292":
		return e.named[3]
	case "dword_5d4594_1090828":
		return e.named[4]
	case "dword_5d4594_1091364":
		return e.named[5]
	case "dword_5d4594_1096252":
		return e.named[6]
	case "dword_5d4594_1096256":
		return e.named[7]
	case "dword_5d4594_1096260":
		return e.named[8]
	case "dword_5d4594_1096264":
		return e.named[9]
	case "dword_5d4594_1096272":
		return e.named[10]
	case "dword_5d4594_1096276":
		return e.named[11]
	case "dword_5d4594_1096280":
		return e.named[12]
	case "dword_5d4594_1096284":
		return e.named[13]
	case "dword_5d4594_1096288":
		return e.named[14]
	case "dword_8531A0_2576":
		return e.named[15]
	case "nox_client_renderBubbles_80844":
		return e.named[16]
	case "nox_color_black_2650656":
		return e.named[17]
	case "nox_color_violet_2598268":
		return e.named[18]
	case "nox_color_white_2523948":
		return e.named[19]
	case "nox_color_yellow_2589772":
		return e.named[20]
	case "nox_gameDisableMapDraw_5d4594_2650672":
		return e.named[21]
	case "nox_player_netCode_85319C":
		return e.named[22]
	case "nox_win_width":
		return e.named[23]
	case "nox_win_height":
		return e.named[24]
	default:
		panic("unknown meter named word")
	}
}
