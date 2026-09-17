//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3_3.h"
extern uint32_t dword_5d4594_1565512;
extern uint32_t dword_5d4594_1565516;
extern uint32_t dword_5d4594_1565520;
extern uint32_t dword_5d4594_2649712;
extern unsigned int dword_5d4594_2650652;
*/
import "C"

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Entry dispatch and observation only; retired entries call production Go directly.
func PortTestReliableReports(op, to, index int, arg uint32, data []byte, related *server.Object, priority int, ordered bool) uint32 {
	var empty byte
	p := unsafe.Pointer(unsafe.SliceData(data))
	if p == nil {
		p = unsafe.Pointer(&empty)
	}
	node := uint32(C.dword_5d4594_1565512)
	if op == 7 || op == 15 {
		for i := 0; i < index && node != 0; i++ {
			node = *(*uint32)(unsafe.Pointer(uintptr(node) + 408))
		}
		if node == 0 {
			panic("missing queue node")
		}
	}
	obj := C.int(uintptr(unsafe.Pointer(related)))
	a, b, n := C.int(to), C.int(uintptr(p)), C.int(len(data))
	order := C.char(0)
	if ordered {
		order = 1
	}
	switch op {
	case 0:
		return uint32(reliableInit())
	case 1:
		return uint32(reliableRecalculate(to))
	case 2:
		return uint32(reliableResetSequences())
	case 3:
		return uint32(reliableResetRates())
	case 4:
		return uint32(C.sub_4E4F30(a))
	case 5:
		return uint32(C.nox_xxx_playerResetImportantCtr_4E4F40(a))
	case 6:
		return uint32(reliableTrim())
	case 7:
		reliableUnlink((*reliableMessage)(unsafe.Pointer(uintptr(node))))
	case 8:
		return uint32(reliableEnqueue(to, data, related, priority, byte(order)))
	case 9:
		return uint32(reliablePressure())
	case 10:
		return uint32(reliableRemoveSlowPlayer(to))
	case 11:
		return uint32(C.nox_xxx_netSendPacket1_4E5390(a, b, n, obj, C.int(priority)))
	case 12:
		return uint32(C.nox_xxx_netClientSend2_4E53C0(a, p, n, obj, C.int(priority)))
	case 13:
		return uint32(C.nox_xxx_netSendPacket0_4E5420(a, p, n, obj, C.int(priority)))
	case 14:
		return uint32(reliableCoalesce(to, data, related, priority))
	case 15:
		reliableAcknowledge(arg, (*reliableMessage)(unsafe.Pointer(uintptr(node))), to)
	case 16:
		return uint32(C.nox_net_importantACK_4E55A0(a, C.int(arg)))
	case 17:
		return uint32(C.sub_4E55F0(C.uchar(to)))
	case 18:
		return uint32(reliableAdapt(byte(to)))
	case 19:
		reliableDeliver(byte(to), int(arg))
	default:
		panic("reliable report operation")
	}
	return 0
}

func PortTestReliableReportGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{"capacity": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1565520)), "mask": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2649712)), "rateMode": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2650652)), "rate": memmap.PtrUint32(0x587000, 4728)}
	saved := map[string]uint32{}
	for k, p := range words {
		saved[k] = *p
	}
	var regions [][]byte
	var originals [][]byte
	for _, r := range [][2]uintptr{{1565124, 384}, {1564964, 160}} {
		p := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, r[0])), r[1])
		regions = append(regions, p)
		originals = append(originals, bytes.Clone(p))
		clear(p)
	}
	return words, func() {
		for k, p := range words {
			*p = saved[k]
		}
		for i, p := range regions {
			copy(p, originals[i])
		}
	}
}

type PortTestReliableReportNode struct {
	Frame                                    uint32
	LastSent                                 [32]uint32
	Countdown                                [32]byte
	Retries                                  byte
	Acknowledged, Sent, Recipients, Priority uint32
	Ordered                                  byte
	Sequence                                 [32]uint16
	To                                       byte
	Data                                     []byte
	Related                                  uint32
}
type PortTestReliableReportState struct {
	Pool     bool
	Capacity uint32
	Sequence [32]uint16
	Rates    []byte
	Nodes    []PortTestReliableReportNode
}

func PortTestReliableReportSnapshot(objects map[unsafe.Pointer]uint32) PortTestReliableReportState {
	r := PortTestReliableReportState{Pool: *memmap.PtrPtr(0x5D4594, 1565508) != nil, Capacity: uint32(C.dword_5d4594_1565520)}
	r.Sequence = *(*[32]uint16)(memmap.PtrOff(0x5D4594, 1565524))
	r.Rates = bytes.Clone(unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565124)), 384))
	previous := uint32(0)
	for q := uint32(C.dword_5d4594_1565512); q != 0; {
		p := unsafe.Pointer(uintptr(q))
		word := func(off uintptr) uint32 { return *(*uint32)(unsafe.Add(p, off)) }
		if word(412) != previous || len(r.Nodes) > 4096 {
			panic("queue list links")
		}
		size := *(*byte)(unsafe.Add(p, 401))
		if size > 150 {
			panic("queue payload length")
		}
		v := PortTestReliableReportNode{Frame: word(0), LastSent: *(*[32]uint32)(unsafe.Add(p, 4)), Countdown: *(*[32]byte)(unsafe.Add(p, 132)), Retries: *(*byte)(unsafe.Add(p, 164)), Acknowledged: word(168), Sent: word(172), Recipients: word(176), Priority: word(180), Ordered: *(*byte)(unsafe.Add(p, 184)), Sequence: *(*[32]uint16)(unsafe.Add(p, 186)), To: *(*byte)(unsafe.Add(p, 250)), Data: bytes.Clone(unsafe.Slice((*byte)(unsafe.Add(p, 251)), int(size)))}
		if obj := *(*unsafe.Pointer)(unsafe.Add(p, 404)); obj != nil {
			var ok bool
			v.Related, ok = objects[obj]
			if !ok {
				panic("unexpected queue object reference")
			}
		}
		r.Nodes = append(r.Nodes, v)
		previous = q
		q = word(408)
	}
	if previous != uint32(C.dword_5d4594_1565516) {
		panic("queue tail link")
	}
	return r
}
