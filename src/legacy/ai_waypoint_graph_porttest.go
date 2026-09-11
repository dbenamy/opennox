//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME5.h"
extern uint32_t dword_5d4594_2490504;
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// PortTestWaypointGraphNode describes one C-backed waypoint. Parent and Next
// are initial traversal-state links, not graph edges; -1 means nil.
type PortTestWaypointGraphNode struct {
	Flags, Epoch uint32
	Flags2       byte
	Edges        []int
	Parent, Next int
}

// PortTestWaypointGraphSpec invokes the original C builder once. Start and End
// may be -1 to exercise its safe ineligible-endpoint path. Capacity is bounded
// to 0..32; the fixture always supplies Capacity+1 writable output words.
type PortTestWaypointGraphSpec struct {
	Nodes          []PortTestWaypointGraphNode
	Start, End     int
	Capacity       int
	Epoch, OneShot uint32
	Scratch        []uint32 // logical scratch slots at blob+2489476, up to 256
}

type PortTestWaypointGraphChange struct {
	Node   int
	Offset int
	Value  int32 // node pointers are normalized to one-based node IDs
}

type PortTestWaypointGraphState struct {
	Epoch        uint32
	Parent, Next int32 // zero=nil, one-based node ID, otherwise raw int32 bits
}

type PortTestWaypointGraphResult struct {
	Return             int
	Epoch, OneShot     uint32
	OutputWords        []int32 // includes the Capacity-th compatibility write
	ScratchWords       []int32
	NodeState          []PortTestWaypointGraphState
	NodeChanges        []PortTestWaypointGraphChange
	OnlyTraversalState bool
	GuardsOK           bool
}

const (
	portTestWaypointScratchOff = uintptr(2489476)
	portTestWaypointScratchLen = 256
	portTestWaypointOneShotOff = uintptr(2490500)
	portTestWaypointGuardWords = 2
)

// PortTestWaypointGraph runs arbitrary small graphs through the original C
// builder. It restores the standalone epoch, logical 256-word scratch, and
// one-shot flag after all operations. The C builder uses offsets 504 (epoch),
// 508 (parent), and 512 (next-frontier); all other waypoint words are checked.
func PortTestWaypointGraph(specs []PortTestWaypointGraphSpec) (out []PortTestWaypointGraphResult, restored bool) {
	oldEngine := noxflags.GetEngine()
	noxflags.UnsetEngine(noxflags.EngineShowAI)
	defer func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) }()
	scratch := unsafe.Slice(memmap.PtrUint32(0x5D4594, portTestWaypointScratchOff), portTestWaypointScratchLen)
	oneShot := memmap.PtrUint32(0x5D4594, portTestWaypointOneShotOff)
	oldScratch := append([]uint32(nil), scratch...)
	oldOneShot, oldEpoch := *oneShot, uint32(C.dword_5d4594_2490504)
	defer func() {
		copy(scratch, oldScratch)
		*oneShot = oldOneShot
		C.dword_5d4594_2490504 = C.uint32_t(oldEpoch)
		restored = bytes.Equal(bytesOfU32(scratch), bytesOfU32(oldScratch)) &&
			*oneShot == oldOneShot && uint32(C.dword_5d4594_2490504) == oldEpoch
	}()

	for _, spec := range specs {
		out = append(out, portTestWaypointGraphOne(spec, scratch, oneShot))
	}
	return out, false
}

func portTestWaypointGraphOne(spec PortTestWaypointGraphSpec, scratch []uint32, oneShot *uint32) PortTestWaypointGraphResult {
	if len(spec.Nodes) > 258 {
		panic("waypoint graph has more than 258 nodes")
	}
	if spec.Capacity < 0 || spec.Capacity > 32 {
		panic("waypoint graph capacity outside 0..32")
	}
	if len(spec.Scratch) > portTestWaypointScratchLen {
		panic("waypoint graph scratch exceeds 256 words")
	}

	// Each record has individual C-visible guards, so a bad waypoint offset
	// cannot be masked by a neighboring record's valid bytes.
	wpSize := int(unsafe.Sizeof(server.Waypoint{}))
	stride := wpSize + 2*portTestWaypointGuardWords*4
	raw, free := alloc.Make([]byte{}, len(spec.Nodes)*stride)
	defer free()
	ptrs := make([]*server.Waypoint, len(spec.Nodes))
	before := make([][]byte, len(spec.Nodes))
	for i := range ptrs {
		base := i * stride
		for j := 0; j < portTestWaypointGuardWords*4; j++ {
			raw[base+j] = byte(0x90 + j)
			raw[base+portTestWaypointGuardWords*4+wpSize+j] = byte(0xd0 + j)
		}
		ptrs[i] = (*server.Waypoint)(unsafe.Pointer(&raw[base+portTestWaypointGuardWords*4]))
	}
	nodeRef := func(ind int) *server.Waypoint {
		if ind == -1 {
			return nil
		}
		if ind < 0 || ind >= len(ptrs) {
			panic(fmt.Sprintf("invalid node reference %d", ind))
		}
		return ptrs[ind]
	}
	for i, n := range spec.Nodes {
		if len(n.Edges) > len(ptrs[i].Points) {
			panic("waypoint has more than 32 edges")
		}
		ptrs[i].Flags, ptrs[i].Flags2, ptrs[i].PointsCnt = n.Flags, n.Flags2, byte(len(n.Edges))
		for j, to := range n.Edges {
			ptrs[i].Points[j].Waypoint = nodeRef(to)
		}
		putWord(ptrs[i], 504, n.Epoch)
		putWord(ptrs[i], 508, ptrWord(nodeRef(n.Parent)))
		putWord(ptrs[i], 512, ptrWord(nodeRef(n.Next)))
		before[i] = bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(ptrs[i])), wpSize))
	}

	for i := range scratch {
		scratch[i] = 0xa5000000 + uint32(i)
	}
	copy(scratch, spec.Scratch)
	*oneShot = spec.OneShot
	C.dword_5d4594_2490504 = C.uint32_t(spec.Epoch)

	// C writes entry Capacity before checking its capacity predicate. Keep it
	// in-bounds and observable, surrounded by two word guards on each side.
	words := spec.Capacity + 1
	outRaw, freeOut := alloc.Make([]uint32{}, words+2*portTestWaypointGuardWords)
	defer freeOut()
	for i := range outRaw {
		outRaw[i] = 0xbeef0000 + uint32(i)
	}
	for i := 0; i < portTestWaypointGuardWords; i++ {
		outRaw[i] = 0x11110000 + uint32(i)
		outRaw[len(outRaw)-portTestWaypointGuardWords+i] = 0xeeee0000 + uint32(i)
	}
	output := outRaw[portTestWaypointGuardWords : portTestWaypointGuardWords+words]

	start, end := nodeRef(spec.Start), nodeRef(spec.End)
	ret := C.nox_xxx_BuildWaypointPath_547F70(
		(*C.uint32_t)(unsafe.Pointer(start)),
		C.int(uintptr(unsafe.Pointer(end))),
		(*C.uint32_t)(unsafe.Pointer(&output[0])),
		C.int(spec.Capacity),
	)

	norm := func(v uint32) int32 {
		if v == 0 {
			return 0
		}
		for i, p := range ptrs {
			if v == ptrWord(p) {
				return int32(i + 1)
			}
		}
		return int32(v)
	}
	r := PortTestWaypointGraphResult{
		Return:             int(ret),
		Epoch:              uint32(C.dword_5d4594_2490504),
		OneShot:            *oneShot,
		OnlyTraversalState: true,
		GuardsOK:           true,
		OutputWords:        make([]int32, len(output)),
		ScratchWords:       make([]int32, len(scratch)),
		NodeState:          make([]PortTestWaypointGraphState, len(ptrs)),
	}
	for i, v := range output {
		r.OutputWords[i] = norm(v)
	}
	for i, v := range scratch {
		r.ScratchWords[i] = norm(v)
	}
	for i, p := range ptrs {
		after := unsafe.Slice((*byte)(unsafe.Pointer(p)), wpSize)
		for off := 0; off < wpSize; off += 4 {
			beforeWord, afterWord := binary.LittleEndian.Uint32(before[i][off:]), binary.LittleEndian.Uint32(after[off:])
			if beforeWord == afterWord {
				continue
			}
			if off != 504 && off != 508 && off != 512 {
				r.OnlyTraversalState = false
			}
			r.NodeChanges = append(r.NodeChanges, PortTestWaypointGraphChange{Node: i, Offset: off, Value: norm(afterWord)})
		}
		r.NodeState[i] = PortTestWaypointGraphState{
			Epoch:  wordAt(p, 504),
			Parent: norm(wordAt(p, 508)),
			Next:   norm(wordAt(p, 512)),
		}
		base := i * stride
		for j := 0; j < portTestWaypointGuardWords*4; j++ {
			r.GuardsOK = r.GuardsOK && raw[base+j] == byte(0x90+j) && raw[base+portTestWaypointGuardWords*4+wpSize+j] == byte(0xd0+j)
		}
	}
	for i := 0; i < portTestWaypointGuardWords; i++ {
		r.GuardsOK = r.GuardsOK && outRaw[i] == 0x11110000+uint32(i) && outRaw[len(outRaw)-portTestWaypointGuardWords+i] == 0xeeee0000+uint32(i)
	}
	return r
}

func ptrWord(p *server.Waypoint) uint32 { return uint32(uintptr(unsafe.Pointer(p))) }
func putWord(p *server.Waypoint, off uintptr, v uint32) {
	*(*uint32)(unsafe.Add(unsafe.Pointer(p), off)) = v
}
func wordAt(p *server.Waypoint, off uintptr) uint32 {
	return *(*uint32)(unsafe.Add(unsafe.Pointer(p), off))
}
func bytesOfU32(v []uint32) []byte {
	if len(v) == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(&v[0])), len(v)*4)
}
