//go:build porttest

package legacy

import (
	"bytes"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"math"
	"unsafe"
)

type PortTestSpawnPlayer struct {
	Pos                 [2]uint32
	View                [2]uint16
	Joined, Busy, Flags uint32
}

// The existing lifecycle fixture owns the sparse player list. Save every
// affected byte, including inactive records, before configuring this case.
func portTestSpawnPlayers(proxy *portTestRoamOwnerServer, specs []PortTestSpawnPlayer) (snapshot func() [][]uint32, restore func()) {
	if len(specs) > len(proxy.life.players) {
		panic("too many spawn-policy players")
	}
	var buffers, before [][]byte
	for i := range proxy.life.players {
		p := &proxy.life.players[i]
		ud := p.UpdateDataPlayer()
		buffers = append(buffers, unsafe.Slice((*byte)(p.CObj()), int(unsafe.Sizeof(*p))), unsafe.Slice((*byte)(p.UpdateData), int(unsafe.Sizeof(*ud))), unsafe.Slice((*byte)(unsafe.Pointer(ud.Player)), int(unsafe.Sizeof(*ud.Player))))
	}
	for _, b := range buffers {
		before = append(before, bytes.Clone(b))
	}
	for i := range proxy.life.players {
		u := &proxy.life.players[i]
		ud := u.UpdateDataPlayer()
		info := ud.Player
		info.Active = 0
		info.PlayerUnit = nil
		proxy.life.ids[uint32(uintptr(unsafe.Pointer(info)))] = uint32(12000 + i)
		proxy.life.ids[uint32(uintptr(u.UpdateData))] = uint32(12010 + i)
		if i >= len(specs) {
			continue
		}
		sp := specs[i]
		info.Active = 1
		info.PlayerUnit = u
		u.ObjFlags = object.Flags(sp.Flags)
		u.PosVec = types.Pointf{X: math.Float32frombits(sp.Pos[0]), Y: math.Float32frombits(sp.Pos[1])}
		*(*uint16)(unsafe.Add(unsafe.Pointer(info), 10)) = sp.View[0]
		*(*uint16)(unsafe.Add(unsafe.Pointer(info), 12)) = sp.View[1]
		*(*uint32)(unsafe.Add(unsafe.Pointer(info), 4792)) = sp.Joined
		*(*uint32)(unsafe.Add(u.UpdateData, 312)) = sp.Busy
	}
	norm := func(v uint32) uint32 {
		if id, ok := proxy.life.ids[v]; ok {
			return id
		}
		return v
	}
	snapshot = func() (out [][]uint32) {
		for _, b := range buffers {
			out = append(out, policyWords(unsafe.Pointer(&b[0]), len(b)/4*4, norm))
		}
		return
	}
	restore = func() {
		for i, b := range buffers {
			copy(b, before[i])
		}
	}
	return
}
