package legacy

/*
#include "defs.h"
#include "GAME1_1.h"
#include "GAME4_3.h"
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func creatureXferArgument(entry unsafe.Pointer, delta uint32) uint32 {
	r := objectXferStream{cryptfile.Global()}
	id := (*uint32)(entry)
	var name [256]byte
	initial := GoString(C.sub_534650(C.int(*id)))
	copy(name[:], initial)
	n := r.byte(byte(len(initial)))
	r.raw(unsafe.Pointer(&name[0]), int(n))
	name[n] = 0
	*id = uint32(C.nox_xxx_actionByName_534670((*C.char)(unsafe.Pointer(&name[0]))))
	count := r.byte(*memmap.PtrUint8(0x587000, 255604+uintptr(16*(*id))))
	for j := uint32(0); j < uint32(count); j++ {
		p := unsafe.Add(entry, 4+8*j)
		kind := *memmap.PtrUint32(0x587000, 255608+uintptr(16*(*id)+4*j))
		switch kind {
		case 0:
			r.raw(p, 8)
		case 1, 2:
			if r.read() {
				r.raw(p, 4)
			} else {
				ref := *(*unsafe.Pointer)(p)
				if ref == nil {
					r.word(0)
				} else if kind == 1 {
					r.raw(unsafe.Add(ref, 44), 4)
				} else {
					r.raw(ref, 4)
				}
			}
		case 3, 4, 6:
			r.raw(p, 4)
		case 5:
			r.raw(p, 4)
			if r.read() {
				creatureXferAdjust((*uint32)(p), delta)
			}
		case 7:
			r.raw(p, 1)
		default:
			return kind
		}
	}
	if _, err := r.cf.ReadWrite(unsafe.Slice((*byte)(unsafe.Add(entry, 20)), 4)); err != nil {
		return 0
	}
	return 1
}

func creatureXferBuffs(u *server.Object) int {
	r := objectXferStream{cryptfile.Global()}
	version := int16(r.short(2))
	if version <= 0 || version > 2 {
		return 0
	}
	count := r.byte(byte(C.sub_424CB0(C.int(uintptr(u.CObj())))))
	if r.read() {
		for i := 0; i < int(count); i++ {
			var name [256]byte
			n := r.byte(0)
			r.raw(unsafe.Pointer(&name[0]), int(n))
			id, ok := server.ParseEnchant(alloc.GoString(&name[0]))
			if !ok {
				return 0
			}
			power := r.byte(0)
			timer := r.word(0)
			arg := server.SpellAcceptArg{Obj: u, Pos: u.PosVec}
			GetServer().Nox_xxx_spellAccept4FD400(id.Spell(), u, u, u, &arg, int(power))
			u.BuffsDur[id] = uint16(timer)
			if id == 26 && version >= 2 {
				extra := r.word(0)
				if duration := spellLifeFindDuration(51, u); duration != nil {
					*(*uint32)(unsafe.Add(unsafe.Pointer(duration), 72)) = extra
				}
			}
		}
		return 1
	}
	for i := int(C.sub_424D00()); i != -1; i = int(C.sub_424D20(C.int(i))) {
		id := server.EnchantID(i)
		if !u.HasEnchant(id) {
			continue
		}
		name := id.String()
		n := r.byte(byte(len(name)))
		r.cf.ReadWrite([]byte(name[:int(n)]))
		r.byte(byte(u.EnchantPower(id)))
		r.word(uint32(u.EnchantDur(id)))
		if id == 26 {
			extra := uint32(100)
			if duration := spellLifeFindDuration(51, u); duration != nil {
				extra = *(*uint32)(unsafe.Add(unsafe.Pointer(duration), 72))
			}
			r.word(extra)
		}
	}
	return 1
}
