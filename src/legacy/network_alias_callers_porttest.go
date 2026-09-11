//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
#include "GAME2_3.h"
*/
import "C"

import (
	"image"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

const (
	portTestAliasTableOff  = 1198020
	portTestAliasTableSize = 255 * 8
	portTestAliasGuardSize = 8
)

type portTestAliasClient struct {
	Client
	calls [][4]int
	last  *client.Drawable
	alloc []unsafe.Pointer
}

func (c *portTestAliasClient) Nox_xxx_spriteCreate_48E970(typeID int, code uint16, x, y int) *client.Drawable {
	p := C.calloc(1, C.size_t(unsafe.Sizeof(client.Drawable{})))
	if p == nil {
		panic("fixture drawable allocation failed")
	}
	dr := (*client.Drawable)(p)
	dr.ObjClass = object.Class(0x200000)
	dr.PosVec = image.Pt(x, y)
	c.calls = append(c.calls, [4]int{typeID, int(code), x, y})
	c.last = dr
	c.alloc = append(c.alloc, p)
	return dr
}

func (c *portTestAliasClient) free() {
	for _, p := range c.alloc {
		C.free(p)
	}
	c.alloc, c.last = nil, nil
}

// PortTestAliasCallerCall invokes one of the two legacy incoming-alias paths.
// Stream 1 is 494A60; Stream 2 is 494C30. InvalidPosition only applies to
// stream 2 and preserves the pre-sprite early-return probe. Recipient is always
// the host because Kind0's production netlist has only a host message queue.
type PortTestAliasCallerCall struct {
	Stream          int
	Key1, Key2      uint16
	Frame           uint32
	InvalidPosition bool
}

type PortTestAliasCallerSnapshot struct {
	ReturnRaw       uint32
	Position        [2]int32
	Camera          [2]int
	CameraCalls     int
	SpriteCalls     [][4]int // typeID, code, x, y
	SpriteFrame     uint32
	SpriteAnim      uint32
	SpriteDirection byte
	Table           []byte
	Guard           []byte
	Announcements   []byte
	AnnouncementCnt int
}

// PortTestAliasCallers saves and restores the fixed blob table and the eight
// bytes directly after it. In an original-C full-table run those guard bytes
// intentionally expose the historical slot-255 write; callers can then assert
// the fixed Go path preserves both the table and guard and emits no message.
func PortTestAliasCallers(initial []byte, calls []PortTestAliasCallerCall) []PortTestAliasCallerSnapshot {
	if len(initial) != portTestAliasTableSize {
		panic("unexpected alias table size")
	}
	raw := memmap.Slice(0x5D4594, portTestAliasTableOff)[:portTestAliasTableSize+portTestAliasGuardSize]
	saved := append([]byte(nil), raw...)
	defer copy(raw, saved)
	copy(raw[:portTestAliasTableSize], initial)
	for i := range raw[portTestAliasTableSize:] {
		raw[portTestAliasTableSize+i] = 0xD7
	}

	oldGet, oldClient, oldCamera := GetServer, GetClient, Nox_xxx_cliUpdateCameraPos_435600
	core := new(server.Server)
	core.NetList = netlist.New()
	core.NetList.Init()
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	proxy := new(portTestAliasClient)
	GetClient = func() Client { return proxy }
	var camera [2]int
	cameraCalls := 0
	Nox_xxx_cliUpdateCameraPos_435600 = func(x, y int) {
		camera, cameraCalls = [2]int{x, y}, cameraCalls+1
	}
	defer func() {
		proxy.free()
		core.NetList.Free()
		GetServer, GetClient, Nox_xxx_cliUpdateCameraPos_435600 = oldGet, oldClient, oldCamera
	}()

	out := make([]PortTestAliasCallerSnapshot, 0, len(calls))
	for _, call := range calls {
		core.SetFrame(call.Frame)
		camera, cameraCalls = [2]int{-999, -999}, 0
		proxy.free()
		proxy.calls = nil
		core.NetList.ResetByInd(server.HostPlayerIndex, netlist.Kind0)
		s := PortTestAliasCallerSnapshot{}
		switch call.Stream {
		case 1:
			// FF, keys, coordinates 321/654, status zero, animation seven.
			packet, freePacket := alloc.Make([]byte{0xFF, byte(call.Key1), byte(call.Key1 >> 8), byte(call.Key2), byte(call.Key2 >> 8), 321 & 0xff, 321 >> 8, 654 & 0xff, 654 >> 8, 0, 7}, 11)
			coords, freeCoords := alloc.Make([]uint32{}, 2)
			s.ReturnRaw = uint32(C.nox_xxx_netCliProcUpdateStream_494A60(
				(*C.uchar)(unsafe.Pointer(&packet[0])), C.int(server.HostPlayerIndex), (*C.uint)(unsafe.Pointer(&coords[0]))))
			s.Position = [2]int32{int32(coords[0]), int32(coords[1])}
			freeCoords()
			freePacket()
		case 2:
			// FF, keys, signed deltas +5/-7, then status seven. A negative
			// start X verifies the range return without creating a drawable.
			packet, freePacket := alloc.Make([]byte{0xFF, byte(call.Key1), byte(call.Key1 >> 8), byte(call.Key2), byte(call.Key2 >> 8), 5, 0xF9, 7}, 8)
			start := [2]int32{300, 400}
			if call.InvalidPosition {
				start = [2]int32{-10, 400}
			}
			coords, freeCoords := alloc.Make([]int32{start[0], start[1]}, 2)
			ret := C.nox_xxx_netCliUpdateStream2_494C30(
				(*C.uchar)(unsafe.Pointer(&packet[0])), C.int(server.HostPlayerIndex), (*C.int)(unsafe.Pointer(&coords[0])))
			s.ReturnRaw = uint32(uintptr(unsafe.Pointer(ret)))
			s.Position = [2]int32{coords[0], coords[1]}
			freeCoords()
			freePacket()
		default:
			panic("unknown alias caller stream")
		}
		s.Camera, s.CameraCalls = camera, cameraCalls
		s.SpriteCalls = append([][4]int(nil), proxy.calls...)
		if dr := proxy.last; dr != nil {
			s.SpriteFrame = dr.Field_72
			s.SpriteAnim = dr.AnimInd
			s.SpriteDirection = dr.AnimDir
		}
		s.Table = append([]byte(nil), raw[:portTestAliasTableSize]...)
		s.Guard = append([]byte(nil), raw[portTestAliasTableSize:]...)
		s.Announcements = core.NetList.CopyPacketsA(ntype.PlayerInd(server.HostPlayerIndex), netlist.Kind0)
		if len(s.Announcements) != 0 {
			s.AnnouncementCnt = len(s.Announcements) / 10
		}
		out = append(out, s)
	}
	return out
}
