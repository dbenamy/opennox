//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME1_1.h"
#include "GAME1_2.h"
void sub_426060();
extern uint32_t dword_5d4594_599496, dword_5d4594_600116, dword_5d4594_608316, dword_5d4594_739392;
extern uint32_t dword_5d4594_741332, dword_5d4594_741648, dword_5d4594_741652;
extern uint32_t dword_5d4594_534808;
#include "GAME5_2.h"
*/
import "C"
import "unsafe"

// These fixtures own real C records and invoke the original record operations.
// Pointer identities are observed separately from the serialized payload.
type PortTestStatisticsField struct {
	Kind  int
	Tag   string
	Value int32
	Data  []byte
}
type PortTestStatisticsNode struct {
	Tag          [4]byte
	Kind, Length uint16
	Data         []byte
	Next         int
}
type PortTestStatisticsRecords struct {
	head  unsafe.Pointer
	nodes []unsafe.Pointer
}

func PortTestStatisticsRecordsOpen(kind uint16) *PortTestStatisticsRecords {
	p := C.calloc(1, 12)
	*(*uint16)(unsafe.Add(p, 2)) = kind
	return &PortTestStatisticsRecords{head: p}
}
func (p *PortTestStatisticsRecords) Close() {
	C.sub_42C330((*C.uint32_t)(p.head))
	C.free(p.head)
	p.head = nil
	p.nodes = nil
}
func (p *PortTestStatisticsRecords) Add(f PortTestStatisticsField, path int) {
	tag := C.CString(f.Tag)
	defer C.free(unsafe.Pointer(tag))
	data := C.CBytes(append(append([]byte(nil), f.Data...), 0))
	defer C.free(data)
	head := (*C.uint32_t)(p.head)
	if path == 2 {
		switch f.Kind {
		case 2:
			C.sub_42BCE0(head, tag, C.char(f.Value))
		case 3:
			C.sub_42BD50(head, tag, C.char(f.Value))
		case 6:
			C.sub_42BDC0(head, tag, C.char(f.Value))
		case 7:
			C.sub_42BE30(head, tag, (*C.char)(data))
		case 20:
			C.sub_42BEA0(head, tag, data, C.ushort(len(f.Data)))
		default:
			panic(f.Kind)
		}
		p.nodes = append(p.nodes, *(*unsafe.Pointer)(unsafe.Add(p.head, 4)))
		return
	}
	node := C.calloc(1, 16)
	x := C.int(uintptr(node))
	if path == 0 {
		switch f.Kind {
		case 2:
			C.sub_42C7F0(x, tag, C.char(f.Value))
		case 3:
			C.sub_42C820(x, tag, C.char(f.Value))
		case 6:
			C.sub_42C8B0(x, tag, C.char(f.Value))
		case 7:
			C.sub_42C8E0(x, tag, (*C.char)(data))
		case 20:
			C.sub_42C910(x, tag, data, C.ushort(len(f.Data)))
		default:
			panic(f.Kind)
		}
	} else {
		statisticsFieldSet(node, f, tag, data)
	}
	C.sub_42C360(head, x)
	p.nodes = append(p.nodes, node)
}
func statisticsFieldSet(node unsafe.Pointer, f PortTestStatisticsField, tag *C.char, data unsafe.Pointer) {
	p := (*unsafe.Pointer)(node)
	switch f.Kind {
	case 2:
		C.sub_42C9A0(p, tag, C.char(f.Value))
	case 3:
		C.sub_42CA00(p, tag, C.char(f.Value))
	case 6:
		C.sub_42CB20(p, tag, C.char(f.Value))
	case 7:
		C.sub_42CB80(p, tag, (*C.char)(data))
	case 20:
		C.sub_42CBF0(p, tag, data, C.ushort(len(f.Data)))
	default:
		panic(f.Kind)
	}
}
func (p *PortTestStatisticsRecords) Replace(index int, f PortTestStatisticsField) {
	// Setters clear the link. Only use this operation on the sole owned node.
	if len(p.nodes) != 1 || index != 0 {
		panic("replace requires one node")
	}
	tag := C.CString(f.Tag)
	defer C.free(unsafe.Pointer(tag))
	data := C.CBytes(append(append([]byte(nil), f.Data...), 0))
	defer C.free(data)
	statisticsFieldSet(p.nodes[index], f, tag, data)
}
func (p *PortTestStatisticsRecords) Snapshot() []PortTestStatisticsNode {
	out := make([]PortTestStatisticsNode, 0, len(p.nodes))
	for _, n := range p.nodes {
		r := PortTestStatisticsNode{Tag: *(*[4]byte)(n), Kind: *(*uint16)(unsafe.Add(n, 4)), Length: *(*uint16)(unsafe.Add(n, 6)), Next: -1}
		r.Data = C.GoBytes(*(*unsafe.Pointer)(unsafe.Add(n, 8)), C.int(r.Length))
		next := *(*unsafe.Pointer)(unsafe.Add(n, 12))
		for i, v := range p.nodes {
			if next == v {
				r.Next = i
			}
		}
		if next != nil && r.Next < 0 {
			panic("unowned statistics link")
		}
		out = append(out, r)
	}
	return out
}
func (p *PortTestStatisticsRecords) Serialize() []byte {
	var n C.uint
	buf := C.sub_42C480((*C.uint32_t)(p.head), &n)
	defer C.free(unsafe.Pointer(buf))
	return C.GoBytes(unsafe.Pointer(buf), C.int(n))
}
func (p *PortTestStatisticsRecords) Endian(index int, encode bool) uint16 {
	n := p.nodes[index]
	if encode {
		return uint16(C.sub_42CC70(C.int(uintptr(n))))
	}
	return uint16(C.sub_42CCE0((*C.uint16_t)(n)))
}

// The caller supplies actual report bytes. This owner installs independently
// allocated arrays at the legacy offsets and frees every allocation afterward.
func PortTestStatisticsReport(quest bool, mode int, raw []byte, count int, events []byte) []byte {
	report := C.calloc(1, 640)
	defer C.free(report)
	copy(unsafe.Slice((*byte)(report), 640), raw)
	var owned []unsafe.Pointer
	defer func() {
		for _, p := range owned {
			C.free(p)
		}
	}()
	keep := func(size int) unsafe.Pointer {
		if size == 0 {
			size = 1
		}
		p := C.calloc(1, C.size_t(size))
		owned = append(owned, p)
		return p
	}
	start, end := 608, 632
	if quest {
		start, end = 536, 576
		*(*int16)(report) = int16(count)
	} else {
		*(*int16)(unsafe.Add(report, 6)) = int16(count)
	}
	for off := start; off <= end; off += 4 {
		array := keep(4 * count)
		*(*unsafe.Pointer)(unsafe.Add(report, off)) = array
		for i := 0; i < count; i++ {
			switch {
			case off == start:
				text := []byte("Player-" + string(rune('A'+i)))
				str := keep(len(text) + 1)
				copy(unsafe.Slice((*byte)(str), len(text)), text)
				*(*unsafe.Pointer)(unsafe.Add(array, 4*i)) = str
			case !quest && (off == 620 || off == 624 || off == 632), quest && off == 576:
				*(*byte)(unsafe.Add(array, i)) = byte(i*71 + off)
			default:
				*(*uint32)(unsafe.Add(array, 4*i)) = uint32(0x12340000 + off + i*131)
			}
		}
	}
	if !quest {
		ev := keep(len(events))
		copy(unsafe.Slice((*byte)(ev), len(events)), events)
		*(*unsafe.Pointer)(unsafe.Add(report, 636)) = ev
	}
	var n C.uint
	var out *C.uint16_t
	if quest {
		out = C.sub_42B810((*C.short)(report), &n)
	} else {
		out = C.sub_42ADA0(C.int(uintptr(report)), C.int(mode), C.short(len(events)/2), &n)
	}
	defer C.free(unsafe.Pointer(out))
	return C.GoBytes(unsafe.Pointer(out), C.int(n))
}
func PortTestStatisticsRuns(data []byte) (out []byte, ret int32) {
	// The original routine reads the lookahead byte even at the final boundary.
	// Own a zero sentinel explicitly, without including it in the logical length.
	input := C.CBytes(append(append([]byte(nil), data...), 0))
	defer C.free(input)
	buf := C.calloc(1, C.size_t(2*len(data)+8))
	defer C.free(buf)
	n := C.int(len(data))
	ret = int32(C.sub_42A970((*C.uint8_t)(input), (*C.uint8_t)(buf), &n))
	return C.GoBytes(buf, n), ret
}
func PortTestStatisticsRandom(seed int32, count int) ([]float64, int32) {
	v := C.int(seed)
	out := make([]float64, count)
	for i := range out {
		out[i] = float64(C.sub_42AAA0(&v))
	}
	return out, int32(v)
}
func PortTestStatisticsRandomBytes(seed int32, count int) []byte {
	buf := C.calloc(1, C.size_t(count+1))
	defer C.free(buf)
	C.sub_42ABF0(C.int(uintptr(buf)), C.int(count), C.int(seed))
	return C.GoBytes(buf, C.int(count))
}

func PortTestStatisticsGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"mode":        (*uint32)(unsafe.Pointer(&C.dword_5d4594_599496)),
		"start":       (*uint32)(unsafe.Pointer(&C.dword_5d4594_600116)),
		"players":     (*uint32)(unsafe.Pointer(&C.dword_5d4594_608316)),
		"events":      (*uint32)(unsafe.Pointer(&C.dword_5d4594_739392)),
		"array-count": (*uint32)(unsafe.Pointer(&C.dword_5d4594_741332)),
		"random-a":    (*uint32)(unsafe.Pointer(&C.dword_5d4594_741648)),
		"random-b":    (*uint32)(unsafe.Pointer(&C.dword_5d4594_741652)),
		"services":    (*uint32)(unsafe.Pointer(&C.dword_5d4594_534808)),
	}
	saved := map[string]uint32{}
	for k, p := range words {
		saved[k] = *p
		*p = 0
	}
	return words, func() {
		for k, p := range words {
			*p = saved[k]
		}
	}
}
func PortTestStatisticsCall(op string, a, b unsafe.Pointer, arg int32) uint32 {
	x, y := C.int(uintptr(a)), C.int(uintptr(b))
	switch op {
	case "event":
		return uint32(uintptr(unsafe.Pointer(C.sub_425CA0(x, y))))
	case "completion":
		return uint32(C.sub_425E90(a, C.char(arg)))
	case "participation":
		return uint32(C.sub_425ED0(x, C.char(arg)))
	case "register":
		C.sub_425F10((*C.nox_playerInfo)(a))
	case "initialize":
		C.sub_426060()
	case "metadata":
		C.sub_426150()
	case "flush":
		C.nox_xxx_net_4263C0()
	case "finish":
		return uint32(C.sub_4264D0())
	case "address":
		C.sub_4282D0((*C.char)(a), C.int(arg))
	case "match-arrays":
		C.sub_4282F0(x, y, C.size_t(arg))
	case "event-array":
		return uint32(C.sub_428540(x, (*C.char)(b), C.int(arg)))
	case "quest-arrays":
		C.sub_4285C0((*C.short)(a))
	case "quest-clear":
		C.sub_4289D0((*unsafe.Pointer)(a))
	case "match-send":
		return uint32(C.sub_428810(x, C.int(arg)))
	case "quest-send":
		return uint32(C.sub_428890((*C.short)(a)))
	default:
		panic(op)
	}
	return 0
}

// PortTestStatisticsFree releases allocations made by the C report owners.
func PortTestStatisticsFree(p unsafe.Pointer) { C.free(p) }

func PortTestStatisticsEncoding(data []byte, outer bool) ([]byte, int32) {
	input := C.CBytes(append(append([]byte(nil), data...), 0))
	defer C.free(input)
	n := C.int(len(data))
	var buf unsafe.Pointer
	if outer {
		buf = unsafe.Pointer(C.sub_42A8B0((*C.uint8_t)(input), &n))
	} else {
		buf = unsafe.Pointer(C.sub_42AC50((*C.uint8_t)(input), (*C.size_t)(unsafe.Pointer(&n))))
	}
	defer C.free(buf)
	if buf == nil {
		return nil, int32(n)
	}
	return C.GoBytes(buf, n), int32(n)
}
