//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME5_2.h"
extern uint32_t dword_8531A0_2576;
extern int nox_cheat_allowall;
*/
import "C"

import (
	"bytes"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

const (
	portTestGlyphPlayerClassOff = 2251
)

type portTestGlyphClient struct {
	Client
	cli *client.Client
}

func (c *portTestGlyphClient) Cli() *client.Client { return c.cli }

// PortTestGlyphEligibilityCall controls one call. Client selects 57B400 and
// Item selects 57B450; Wrapper invokes legacy.Sub_57B450 for Item. A nil
// drawable is valid only for Item, matching its C short-circuit.
type PortTestGlyphEligibilityCall struct {
	Kind                  string
	Wrapper               bool
	Drawable              bool
	Type, Class, Subclass uint32
	Player                bool
	LocalPlayer           bool
	PlayerClass           byte
	Cheat                 bool
	ClassMask             byte
	LookupType            uint32 // Glyph ID returned by this call's bare client
}

type PortTestGlyphEligibilityTrace struct {
	Class, Subclass, Type uint32
}

type PortTestGlyphEligibilitySnapshot struct {
	Return            int
	ClientCache       uint32
	ItemCache         uint32
	LookupCalls       int
	ClassMaskTrace    []PortTestGlyphEligibilityTrace
	DrawableUntouched bool
	PlayerUntouched   bool
}

// PortTestGlyphEligibility supplies C-owned storage and a minimal real Glyph
// lookup, preserving and restoring all globals used by the predicates.
func PortTestGlyphEligibility(clientCacheInit, itemCacheInit uint32, calls []PortTestGlyphEligibilityCall) []PortTestGlyphEligibilitySnapshot {
	clientCache := &glyphClientType
	itemCache := &glyphItemType
	local := memmap.PtrUint32(0x852978, 8)
	oldClientCache, oldItemCache, oldLocal := *clientCache, *itemCache, *local
	*clientCache, *itemCache = clientCacheInit, itemCacheInit
	oldPlayer, oldCheat := C.dword_8531A0_2576, C.nox_cheat_allowall
	oldGet, oldMask := GetClient, Sub_57B370
	defer func() {
		*clientCache, *itemCache, *local = oldClientCache, oldItemCache, oldLocal
		C.dword_8531A0_2576, C.nox_cheat_allowall = oldPlayer, oldCheat
		GetClient, Sub_57B370 = oldGet, oldMask
	}()

	// The porttest helper initializes only Things.byID["glyph"]. A zero ID
	// leaves it absent, whose nil ObjectType Index returns zero.
	localDrawable, freeLocal := alloc.New(client.Drawable{})
	defer freeLocal()
	proxy := &portTestGlyphClient{}
	lookupCalls := 0
	GetClient = func() Client {
		lookupCalls++
		return proxy
	}
	var maskTrace []PortTestGlyphEligibilityTrace
	mask := byte(0)
	Sub_57B370 = func(cl object.Class, sub object.SubClass, typ int) byte {
		maskTrace = append(maskTrace, PortTestGlyphEligibilityTrace{uint32(cl), uint32(sub), uint32(typ)})
		return mask
	}

	out := make([]PortTestGlyphEligibilitySnapshot, 0, len(calls))
	for _, call := range calls {
		lookupClient, freeLookup := client.PortTestGlyphClient(call.LookupType)
		proxy.cli = lookupClient
		mask = call.ClassMask
		C.nox_cheat_allowall = C.int(0)
		if call.Cheat {
			C.nox_cheat_allowall = 1
		}
		var player unsafe.Pointer
		var freePlayer func()
		if call.Player {
			player, freePlayer = alloc.Malloc(unsafe.Sizeof(server.Player{}))
			for i := range unsafe.Slice((*byte)(player), unsafe.Sizeof(server.Player{})) {
				*(*byte)(unsafe.Add(player, i)) = 0xa7
			}
			*(*byte)(unsafe.Add(player, portTestGlyphPlayerClassOff)) = call.PlayerClass
			C.dword_8531A0_2576 = C.uint(uintptr(player))
		} else {
			C.dword_8531A0_2576 = 0
		}
		var drawable unsafe.Pointer
		var freeDrawable func()
		if call.Drawable {
			drawable, freeDrawable = alloc.Malloc(unsafe.Sizeof(client.Drawable{}))
			for i := range unsafe.Slice((*byte)(drawable), unsafe.Sizeof(client.Drawable{})) {
				*(*byte)(unsafe.Add(drawable, i)) = 0xa7
			}
			dr := (*client.Drawable)(drawable)
			dr.TypeIDVal = call.Type
			dr.ObjClass = object.Class(call.Class)
			dr.ObjSubClass = object.SubClass(call.Subclass)
		}
		if call.LocalPlayer {
			*local = uint32(uintptr(unsafe.Pointer(localDrawable)))
		} else {
			*local = 0
		}
		var playerBefore []byte
		if player != nil {
			playerBefore = append([]byte(nil), unsafe.Slice((*byte)(player), unsafe.Sizeof(server.Player{}))...)
		}
		var before []byte
		if drawable != nil {
			before = append([]byte(nil), unsafe.Slice((*byte)(drawable), unsafe.Sizeof(client.Drawable{}))...)
		}
		lookupStart, maskStart := lookupCalls, len(maskTrace)
		s := PortTestGlyphEligibilitySnapshot{}
		switch call.Kind {
		case "client":
			if drawable == nil && call.Player {
				panic("client call with player requires drawable")
			}
			s.Return = int(C.nox_xxx_client_57B400(C.int(uintptr(drawable))))
		case "item":
			if call.Wrapper {
				s.Return = Sub_57B450((*client.Drawable)(drawable))
			} else {
				s.Return = glyphItemAllowed((*client.Drawable)(drawable))
			}
		default:
			panic("unknown glyph eligibility call")
		}
		s.ClientCache, s.ItemCache = *clientCache, *itemCache
		s.LookupCalls = lookupCalls - lookupStart
		s.ClassMaskTrace = append([]PortTestGlyphEligibilityTrace(nil), maskTrace[maskStart:]...)
		if drawable == nil {
			s.DrawableUntouched = true
		} else {
			s.DrawableUntouched = bytes.Equal(before, unsafe.Slice((*byte)(drawable), unsafe.Sizeof(client.Drawable{})))
		}
		s.PlayerUntouched = player == nil || bytes.Equal(playerBefore, unsafe.Slice((*byte)(player), unsafe.Sizeof(server.Player{})))
		out = append(out, s)
		if freeDrawable != nil {
			freeDrawable()
		}
		if freePlayer != nil {
			freePlayer()
		}
		freeLookup()
	}
	return out
}
