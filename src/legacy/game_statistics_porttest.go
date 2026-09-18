//go:build porttest

package legacy

/*
#include <stdlib.h>
#include <stdint.h>
extern uint32_t dword_5d4594_534808;
*/
import "C"
import "unsafe"

// These fixtures own the production Go records and report operations.
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
	head  statisticsRecords
	nodes []*statisticsRecord
}

func PortTestStatisticsRecordsOpen(kind uint16) *PortTestStatisticsRecords {
	return &PortTestStatisticsRecords{head: statisticsRecords{kind: kind}}
}
func (p *PortTestStatisticsRecords) Close() { p.head.head = nil; p.nodes = nil }
func (p *PortTestStatisticsRecords) Add(f PortTestStatisticsField, path int) {
	var r *statisticsRecord
	if path == 2 {
		r = p.head.add(uint16(f.Kind), f.Tag, f.Value, f.Data)
	} else {
		r = new(statisticsRecord)
		r.set(uint16(f.Kind), f.Tag, f.Value, f.Data)
		p.head.prepend(r)
	}
	p.nodes = append(p.nodes, r)
}
func (p *PortTestStatisticsRecords) Replace(index int, f PortTestStatisticsField) {
	if len(p.nodes) != 1 || index != 0 {
		panic("replace requires one node")
	}
	p.nodes[index].set(uint16(f.Kind), f.Tag, f.Value, f.Data)
}
func (p *PortTestStatisticsRecords) Snapshot() []PortTestStatisticsNode {
	out := make([]PortTestStatisticsNode, 0, len(p.nodes))
	for _, n := range p.nodes {
		r := PortTestStatisticsNode{Tag: n.tag, Kind: n.kind, Length: n.length, Data: append([]byte{}, n.data[:int(n.length)]...), Next: -1}
		for i, v := range p.nodes {
			if n.next == v {
				r.Next = i
			}
		}
		if n.next != nil && r.Next < 0 {
			panic("unowned statistics link")
		}
		out = append(out, r)
	}
	return out
}
func (p *PortTestStatisticsRecords) Serialize() []byte { return p.head.serialize() }
func (p *PortTestStatisticsRecords) Endian(index int, encode bool) uint16 {
	if encode {
		return p.nodes[index].encode()
	}
	return p.nodes[index].decode()
}

// The caller supplies actual report bytes. This owner installs independently
// allocated arrays at the legacy offsets and frees every allocation afterward.
func PortTestStatisticsReport(quest bool, mode int, raw []byte, count int, events []byte) []byte {
	if quest {
		panic("retired quest statistics fixture")
	}
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
	*(*int16)(unsafe.Add(report, 6)) = int16(count)
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
			case off == 620 || off == 624 || off == 632:
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
	return statisticsReport(report, mode, int16(len(events)/2))
}
func PortTestStatisticsRuns(data []byte) ([]byte, int32) {
	out := statisticsRuns(data)
	return out, int32(len(out))
}
func PortTestStatisticsRandom(seed int32, count int) ([]float64, int32) {
	out := make([]float64, count)
	for i := range out {
		out[i] = statisticsRandom(&seed)
	}
	return out, seed
}
func PortTestStatisticsRandomBytes(seed int32, count int) []byte {
	return statisticsRandomBytes(seed, count)
}

func PortTestStatisticsGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{

		"start":       &statisticsStart,
		"players":     &statisticsPlayers,
		"events":      &statisticsEvents,
		"array-count": &statisticsArrayCount,
		"random-a":    &statisticsRandomA,
		"random-b":    &statisticsRandomB,
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
	switch op {
	case "event":
		return statisticsEvent(a, b)
	case "participation":
		return statisticsParticipation(a, byte(arg))
	case "register":
		statisticsRegister(a)
	case "flush":
		statisticsFlush()
	case "match-arrays":
		statisticsMatchArrays(a, b, int(arg))
	case "event-array":
		return statisticsEventArray(a, b, int(arg))
	case "match-send":
		return statisticsSend(a, int(arg))
	default:
		panic(op)
	}
	return 0
}
func PortTestStatisticsFree(p unsafe.Pointer) { statisticsFree(p) }
func PortTestStatisticsEncoding(data []byte, outer bool) ([]byte, int32) {
	if outer {
		return statisticsEncode(data)
	}
	return statisticsFrame(data)
}
