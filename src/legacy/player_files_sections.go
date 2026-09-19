package legacy

/*
#include "GAME1_1.h"
#include "GAME3_3.h"
#include "server__magic__plyrgide.h"
#include "server__magic__plyrspel.h"
#include "server__ability__ability.h"
*/
import "C"
import (
	"github.com/opennox/libs/spell"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/music"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func playerFileStream() objectXferStream { return objectXferStream{cf: cryptfile.Global()} }
func playerFileName(r objectXferStream, name string) string {
	n := r.byte(byte(len(name)))
	if !r.read() {
		r.cf.ReadWrite([]byte(name)[:int(n)])
		return name
	}
	var buf [256]byte
	r.raw(unsafe.Pointer(&buf[0]), int(n))
	return alloc.GoStringS(buf[:])
}
func playerFilePresent(r objectXferStream, skip bool) bool {
	value := byte(1)
	if skip {
		value = 0
	}
	return r.byte(value) != 0
}
func playerFileInventoryCount() int {
	count := 0
	for row := 0; row < 20; row++ {
		for col := 0; col < 4; col++ {
			count += int(uiInventoryGrid()[row+21*col].Count)
		}
	}
	return count
}
func playerFileInventoryAllowed(u *server.Object) bool {
	cache := memmap.PtrUint32(0x5D4594, 527724)
	if *cache == 0 {
		*cache = uint32(GetServer().S().Types.IndByID("Glyph"))
	}
	return uint32(u.ObjClass)&0x40 == 0 && uint32(u.TypeInd) != *cache
}
func playerFileMusic() int {
	r := playerFileStream()
	version := int16(r.short(11))
	if version > 11 {
		return 0
	}
	if version >= 11 && !playerFilePresent(r, noxflags.HasGame(8192)) {
		return 1
	}
	current := music.MusicState{}
	count := uint32(0)
	if !r.read() {
		current = MusicModule.GetCurrentBlock()
		count = *audioEventMusicCount
	}
	transfer := func(s *music.MusicState) {
		s.D = r.word(s.D)
		s.Position = r.word(s.Position)
		s.MusicIdx = r.word(s.MusicIdx)
		s.Volume = r.word(s.Volume)
	}
	transfer(&current)
	count = r.word(count)
	for i := int32(0); i < int32(count); i++ {
		slot := audioEventMusicSlot(uint32(i))
		if r.read() {
			*slot = music.MusicState{}
		}
		transfer(slot)
	}
	if r.read() && !noxflags.HasGame(8192) {
		MusicModule.SetNextMusic(current)
		audioEventMusicSetCount(int32(count))
	}
	return 1
}
func playerFileJournal(u *server.Object) int {
	p := u.UpdateDataPlayer().Player
	r := playerFileStream()
	if int16(r.short(1)) > 1 {
		return 0
	}
	if !playerFilePresent(r, noxflags.HasGame(8192)) {
		return 1
	}
	if !noxflags.HasGame(2048) {
		return 0
	}
	count := uint16(0)
	for n := p.Journal; n != nil; n = n.Next {
		count++
	}
	count = r.short(count)
	if r.read() {
		journalRemoveMask(u, 0xffff)
		for i := 0; i < int(count); i++ {
			name := playerFileName(r, "")
			flags := r.short(0)
			journalUnitAdd(u, name, flags)
		}
	} else if count != 0 {
		n := p.Journal
		for n.Next != nil {
			n = n.Next
		}
		for ; n != nil; n = n.Prev {
			playerFileName(r, alloc.GoString(&n.EntryBuf[0]))
			r.raw(unsafe.Pointer(&n.Field3), 2)
		}
	}
	return 1
}
func playerFileGuides(u *server.Object) int {
	p := GetServer().S().Players.ByID(int(u.NetCode))
	if p == nil {
		return 0
	}
	r := playerFileStream()
	if int16(r.short(1)) > 1 {
		return 0
	}
	if !playerFilePresent(r, noxflags.HasGame(8192) && !noxflags.HasGame(4096)) {
		return 1
	}
	if !noxflags.HasGame(2048 | 4096) {
		return 0
	}
	known := unsafe.Slice((*uint32)(unsafe.Add(unsafe.Pointer(p), 4244)), 41)
	if r.read() {
		count := r.byte(0)
		if count > 41 {
			return 0
		}
		for i := 0; i < int(count); i++ {
			name := playerFileName(r, "")
			id := uint32(C.nox_xxx_guide_427010(internCStr(name)))
			if noxflags.HasGame(4096) && !questEligibilityBeast(id) {
				return 0
			}
			C.nox_xxx_awardBeastGuide_4FAE80_magic_plyrgide(inventoryInt(u), C.int(id), 0)
		}
	} else {
		count := byte(0)
		for i := 1; i <= 40; i++ {
			if known[i] != 0 {
				count++
			}
		}
		r.byte(count)
		for i := 1; i <= 40; i++ {
			if known[i] != 0 {
				playerFileName(r, GoString(C.nox_xxx_guideNameByN_427230(C.int(i))))
			}
		}
	}
	return 1
}
func playerFileSpells(u *server.Object) int {
	p := GetServer().S().Players.ByID(int(u.NetCode))
	if p == nil {
		return 0
	}
	r := playerFileStream()
	version := int16(r.short(3))
	if version > 3 {
		return 0
	}
	if !playerFilePresent(r, noxflags.HasGame(8192) && !noxflags.HasGame(4096)) {
		return 1
	}
	if !noxflags.HasGame(2048 | 4096) {
		return 0
	}
	class := *(*byte)(unsafe.Add(unsafe.Pointer(p), 2251))
	levels := unsafe.Slice((*uint32)(unsafe.Add(unsafe.Pointer(p), 3696)), 137)
	if r.read() {
		count := r.byte(0)
		if count > 137 {
			return 0
		}
		for i := 0; i < int(count); i++ {
			name := playerFileName(r, "")
			level := uint32(3)
			if version >= 2 {
				level = r.word(level)
			}
			if noxflags.HasGame(4096) && class != 0 && int32(level) > 3 {
				return 0
			}
			id := spell.ParseID(name)
			if id < 0 {
				id = 0
			}
			if noxflags.HasGame(4096) && class != 0 && !questEligibilitySpell(uint32(id)) {
				return 0
			}
			if version < 3 || class != 0 {
				Nox_xxx_spellGrantToPlayer_4FB550(u, id, 0, 0, int(int32(level)))
			} else {
				aid := C.nox_xxx_abilityNameToN_424D80(internCStr(name))
				C.nox_xxx_abilityRewardServ_4FB9C0_ability(inventoryInt(u), aid, 0)
			}
		}
	} else {
		limit := 136
		if class == 0 {
			limit = 5
		}
		count := byte(0)
		for i := 1; i <= limit; i++ {
			if levels[i] != 0 {
				count++
			}
		}
		r.byte(count)
		for i := 1; i <= limit; i++ {
			if levels[i] != 0 {
				name := spell.ID(i).String()
				if class == 0 {
					name = server.Ability(i).String()
				}
				playerFileName(r, name)
				r.raw(unsafe.Pointer(&levels[i]), 4)
			}
		}
	}
	return 1
}
