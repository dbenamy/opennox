//go:build porttest

package server

import (
	"bytes"
	"crypto/sha256"
	"image"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// The bounded wall pool uses normal production creation/deletion and indices.
// Every Wall starts at its tracked allocation base so normal frees remain valid.
type PortTestPaintOwners struct {
	S                  *Server
	blocks             [][]byte
	walls              []*Wall
	defs               [80]WallDef
	byPosFree, byYFree func()
	typeData           []unsafe.Pointer
	objects            []*Object
	objectData         []unsafe.Pointer
	handle             uintptr
	definitionBytes    []byte
	definitionHash     [32]byte
}
type PortTestPaintWallState struct {
	Head, Free  uint32
	ByPos, ByY  [][2]uint32
	Records     [][9]uint32
	Definitions [32]byte
	Intact      bool
}

func (s *Server) PortTestMapPaintingOwners(xfers [2]unsafe.Pointer) *PortTestPaintOwners {
	if s.Walls.byPos != nil || s.handle != 0 {
		panic("map painting needs fresh owner")
	}
	p := &PortTestPaintOwners{S: s}
	p.handle = atomic.AddUintptr(&serverLast, 1)
	s.handle = p.handle
	servers.Store(p.handle, s)
	s.Walls.byPos, p.byPosFree = alloc.Make([]*Wall{}, wallsPerBucket*WallGridSize)
	s.Walls.indexY, p.byYFree = alloc.Make([]*Wall{}, WallGridSize)
	for i := 0; i < 512; i++ {
		b, _ := alloc.Make([]byte{}, int(unsafe.Sizeof(Wall{}))+16)
		p.blocks = append(p.blocks, b)
		p.walls = append(p.walls, (*Wall)(unsafe.Pointer(&b[0])))
	}
	for i, name := range []string{"PaintWall", "PaintWallTwo", "InvisibleWallSet"} {
		d := &p.defs[i]
		copy(d.Field0[:], name)
		d.Health41 = 80
		d.Field749 = 3
		for a := range d.Variations12272 {
			for j := range d.Variations12272[a] {
				d.Variations12272[a][j] = 4
			}
		}
	}
	s.Types = serverObjTypes{byInd: make([]*ObjectType, 5), byID: map[string]*ObjectType{}}
	for i, name := range []string{"PaintObject", "PaintDoor", "PaintBook", "PaintMonster"} {
		id := strings.ToLower(name)
		typ := &ObjectType{s: &s.Types, ind: uint16(i + 1), ind2: uint16(i + 1), id: id, class: object.ClassSimple, allowed: true, Mass: 1}
		switch i {
		case 1:
			typ.class = object.ClassDoor
			typ.Xfer = xfers[0]
			typ.UpdateDataSize = 64
		case 2:
			typ.Xfer = xfers[1]
			typ.UseDataSize = 4
		case 3:
			typ.class = object.ClassMonster
			typ.UpdateDataSize = unsafe.Sizeof(MonsterUpdateData{})
		}
		if typ.UpdateDataSize != 0 {
			typ.UpdateData, _ = alloc.Malloc(uintptr(typ.UpdateDataSize))
			p.typeData = append(p.typeData, typ.UpdateData)
		}
		if typ.UseDataSize != 0 {
			ptr, _ := alloc.Malloc(typ.UseDataSize)
			typ.UseData.SetPtr(ptr)
			p.typeData = append(p.typeData, ptr)
		}
		s.Types.byInd[i+1] = typ
		s.Types.byID[id] = typ
	}
	s.Objs.init(p.handle)
	if !s.Objs.Init(64) {
		panic("map painting object pool")
	}
	s.SetFrame(123)
	p.Reset(12345, 3, 4)
	return p
}
func (p *PortTestPaintOwners) TrackObjects(objects ...*Object) {
	for _, u := range objects {
		if u == nil {
			continue
		}
		found := false
		for _, old := range p.objects {
			if old == u {
				found = true
				break
			}
		}
		if found {
			continue
		}
		if u.UseData.Ptr != nil {
			p.objectData = append(p.objectData, u.UseData.Ptr)
		}
		p.objects = append(p.objects, u)
		if u.UpdateData != nil {
			p.objectData = append(p.objectData, u.UpdateData)
		}
	}
}
func (p *PortTestPaintOwners) Objects() []*Object { return p.objects }
func (p *PortTestPaintOwners) Reset(seed int, cycle, variations byte) {
	for _, ptr := range p.objectData {
		alloc.FreePtr(ptr)
	}
	p.objectData = nil
	p.objects = nil
	p.S.Objs.FreeObjects()
	p.S.Objs = serverObjects{}
	p.S.Objs.init(p.handle)
	if !p.S.Objs.Init(64) {
		panic("map painting object pool reset")
	}
	clear(p.S.Walls.byPos)
	clear(p.S.Walls.indexY)
	p.S.Walls.head = nil
	p.S.Walls.freeList = nil
	for i := len(p.walls) - 1; i >= 0; i-- {
		clear(p.blocks[i])
		for j := int(unsafe.Sizeof(Wall{})); j < len(p.blocks[i]); j++ {
			p.blocks[i][j] = 0x5a
		}
		w := p.walls[i]
		w.Next20 = p.S.Walls.freeList
		p.S.Walls.freeList = w
	}
	p.S.Walls.defs = p.defs
	p.S.Walls.defsCnt = 3
	for i := 0; i < 3; i++ {
		d := &p.S.Walls.defs[i]
		d.Field749 = cycle
		for a := range d.Variations12272 {
			for j := range d.Variations12272[a] {
				d.Variations12272[a][j] = variations
			}
		}
	}
	p.S.Rand.Logic, p.S.Rand.Other = prand.New(seed), prand.New(seed+1)
}
func (p *PortTestPaintOwners) NormalizeWall(v uint32) uint32 {
	for i, w := range p.walls {
		start := uint32(uintptr(unsafe.Pointer(w)))
		if v >= start && v-start < uint32(unsafe.Sizeof(Wall{})) {
			return 0x20000000 + uint32(i)*64 + v - start
		}
	}
	return v
}
func (p *PortTestPaintOwners) WallSnapshot(normalize func(uint32) uint32) (out PortTestPaintWallState) {
	norm := func(w *Wall) uint32 { return normalize(uint32(uintptr(unsafe.Pointer(w)))) }
	out.Head = norm(p.S.Walls.head)
	out.Free = norm(p.S.Walls.freeList)
	out.Intact = true
	for i, w := range p.S.Walls.byPos {
		if w != nil {
			out.ByPos = append(out.ByPos, [2]uint32{uint32(i), norm(w)})
		}
	}
	for i, w := range p.S.Walls.indexY {
		if w != nil {
			out.ByY = append(out.ByY, [2]uint32{uint32(i), norm(w)})
		}
	}
	for i, w := range p.walls {
		var record [9]uint32
		for j, v := range *(*[9]uint32)(unsafe.Pointer(w)) {
			record[j] = normalize(v)
		}
		out.Records = append(out.Records, record)
		out.Intact = out.Intact && bytes.Equal(p.blocks[i][36:], bytes.Repeat([]byte{0x5a}, 16))
	}
	// Compare the complete storage so mutations outside active definitions are
	// still detected. Reuse only the digest of byte-identical state.
	definitions := unsafe.Slice((*byte)(unsafe.Pointer(&p.S.Walls.defs[0])), int(unsafe.Sizeof(p.S.Walls.defs)))
	if !bytes.Equal(p.definitionBytes, definitions) {
		p.definitionBytes = append(p.definitionBytes[:0], definitions...)
		p.definitionHash = sha256.Sum256(definitions)
	}
	out.Definitions = p.definitionHash
	return out
}
func (p *PortTestPaintOwners) CreateWall(x, y int) *Wall {
	return p.S.Walls.CreateAtGrid(image.Pt(x, y))
}
func (p *PortTestPaintOwners) NormalizeObject(u *Object, words []uint32) bool {
	if u.serverHandle != p.handle || u.objectHandle != 0 {
		return false
	}
	words[unsafe.Offsetof(u.serverHandle)/4] = 1
	return true
}
func (p *PortTestPaintOwners) ObjectState() [7]uint32 {
	return [7]uint32{uint32(p.S.Objs.Alive), uint32(p.S.Objs.Created), uint32(p.S.Objs.CreatedSimple), uint32(p.S.Objs.CreatedImmobile), uint32(p.S.Objs.lastScriptID), uint32(uintptr(unsafe.Pointer(p.S.Objs.List))), uint32(uintptr(unsafe.Pointer(p.S.Objs.Pending)))}
}
func (p *PortTestPaintOwners) Close() {
	for _, ptr := range p.objectData {
		alloc.FreePtr(ptr)
	}
	p.objectData = nil
	p.S.Objs.FreeObjects()
	for _, ptr := range p.typeData {
		alloc.FreePtr(ptr)
	}
	for _, b := range p.blocks {
		alloc.FreeSlice(b)
	}
	p.byYFree()
	p.byPosFree()
	servers.Delete(p.handle)
}
