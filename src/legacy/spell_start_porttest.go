//go:build porttest

package legacy

/*
#include "server__magic__spell__execdur.h"
#include "GAME4_3.h"
extern int nox_cheat_charmall;
*/
import "C"

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/strman"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestSpellStartSpec struct {
	Blocked, WarmCache, NilOwner, PlayerOwner bool
	Existing                                  int
	Radius                                    float32
}

func (p *portTestShopPools) spellStartTeleportContract(record unsafe.Pointer) uint32 {
	sp := p.sustainedSpec().Start
	if sp == nil {
		panic("missing spell-start contract")
	}
	caster := *(**server.Object)(unsafe.Add(record, 16))
	recipient := *(**server.Object)(unsafe.Add(record, 48))
	if recipient == nil {
		recipient = caster
	}
	caster.PosVec = types.Ptf(100, 100)
	recipient.PosVec = types.Ptf(100, 100)
	ray := [4]float32{100, 100, 200, 100}
	if spatialRay(&ray) == sp.Blocked {
		panic("spell-start fixture ray")
	}
	configure, restoreStrings := p.proxy.core.PortTestMeterStrings(strman.Entry{ID: "ExecDur.c:UnseenTarget", Vals: []strman.Variant{{Str: "Unseen destination"}}})
	configure(0)
	defer restoreStrings()
	beforeMessages := make([][]byte, 3)
	for i := range beforeMessages {
		beforeMessages[i] = p.proxy.core.NetList.CopyPacketsA(ntype.PlayerInd(i), 1)
	}
	audioBefore := len(p.proxy.core.PortTestCombatAudioSnapshot())
	frame := p.proxy.core.Frame()
	delay := *equipmentWord(record, 68)
	tile := p.spellLifeSpec().Effects.Tile
	if tile == nil {
		panic("spell-start fixture tile")
	}
	reject := *tile != 0 || sp.Blocked
	want := uint32(0)
	if reject {
		want = 1
	} else if noxflags.HasGame(noxflags.GameModeCoop) {
		level := int32(*equipmentWord(record, 8))
		delay = frame + uint32(int32(float32(p.proxy.core.Balance.FloatInd("TeleportDelay", int(level-1)))))
	} else {
		delay = frame + 1
	}
	got := uint32(C.sub_530A30_spell_execdur(C.int(uintptr(record))))
	if got != want || *equipmentWord(record, 68) != delay || *(**server.Object)(unsafe.Add(record, 48)) != recipient {
		panic(fmt.Sprintf("teleport-start return/deadline/recipient: got %d/%d want %d/%d", got, *equipmentWord(record, 68), want, delay))
	}
	var wantMessage []byte
	recipientIndex := -1
	if *tile != 0 && recipient.ObjClass&4 != 0 {
		recipientIndex = int(recipient.UpdateDataPlayer().Player.PlayerInd)
		wantMessage = append([]byte{168, 0, 0, 2, 0, 0, 0, 0, 19, 0, 0}, []byte("Unseen destination\x00")...)
	} else if *tile == 0 && sp.Blocked && caster.ObjClass&4 != 0 {
		recipientIndex = int(caster.UpdateDataPlayer().Player.PlayerInd)
		wantMessage = []byte{169, 0, 2, 0, 0, 0}
	}
	for i, before := range beforeMessages {
		want := append([]byte(nil), before...)
		if i == recipientIndex {
			want = append(want, wantMessage...)
		}
		if !bytes.Equal(p.proxy.core.NetList.CopyPacketsA(ntype.PlayerInd(i), 1), want) {
			panic("teleport-start notification contract")
		}
	}
	events := p.proxy.core.PortTestCombatAudioSnapshot()[audioBefore:]
	if *tile != 0 {
		if len(events) != 1 || events[0].ID != 231 || events[0].Obj != recipient {
			panic("teleport-start audio")
		}
	} else if len(events) != 0 {
		panic("teleport-start unexpected audio")
	}
	return got
}

func (p *portTestShopPools) spellStartPixieContract(owner, origin *server.Object, record unsafe.Pointer) uint32 {
	sp := p.sustainedSpec().Start
	core := p.proxy.core
	typ := core.Types.ByID("Pixie")
	typeID := uint16(0x7000)
	if typ != nil {
		typeID = uint16(typ.Ind())
	}
	cache := memmap.PtrUint32(0x5D4594, 2489140)
	oldCache := *cache
	defer func() { *cache = oldCache }()
	*cache = 0
	if sp.WarmCache {
		*cache = uint32(typeID)
	}
	if sp.NilOwner {
		owner = nil
	} else if sp.PlayerOwner {
		owner = p.temporaryRef(100)
	}
	if owner != nil {
		oldOwned := owner.Field129
		defer func() { owner.Field129 = oldOwned }()
		// Exercise real owned-object traversal, including two ignored entries.
		owner.Field129 = nil
		for i := 0; i < sp.Existing+2; i++ {
			u := (*server.Object)(p.objectiveRegion(int(unsafe.Sizeof(server.Object{}))))
			u.TypeInd = typeID
			u.Field128 = owner.Field129
			owner.Field129 = u
			if i == sp.Existing {
				u.ObjFlags = object.FlagDestroyed
			}
			if i == sp.Existing+1 {
				u.TypeInd = typeID + 1
			}
		}
	} else if sp.Existing != 0 {
		panic("nil pixie owner cannot own objects")
	}
	count := owner.CountSubOfType(int(typeID))
	if count != sp.Existing {
		panic("pixie fixture owned count")
	}
	origin.PosVec = types.Ptf(100, 100)
	*(*float32)(unsafe.Add(origin.CObj(), 176)) = sp.Radius
	id, level := *equipmentWord(record, 4), *equipmentWord(record, 8)
	limit := int(int64(core.Balance.FloatInd("PixieCount", int(int32(level)-1))))
	attempts := limit - count
	if attempts < 0 {
		attempts = 0
	}
	// Independent expectations use shipped direction coefficients and the actual
	// wall query, without substituting a factory, RNG or target-selection result.
	rng := prand.New(core.Rand.Logic.Index())
	type expected struct {
		pos    types.Pointf
		dir    uint16
		expiry uint32
	}
	var want []expected
	radius := float32(sp.Radius + 4)
	for i := 0; i < attempts; i++ {
		dir := rng.Int(0, 255)
		pos := types.Ptf(float32(float64(radius)*float64(*memmap.PtrFloat32(0x587000, 194136+8*uintptr(dir)))+100), float32(float64(radius)*float64(*memmap.PtrFloat32(0x587000, 194140+8*uintptr(dir)))+100))
		if core.MapTraceRay(origin.PosVec, pos, 5) && typ != nil {
			expiry := core.Frame() + uint32(core.TickRate())*uint32(rng.Int(30, 90))
			want = append(want, expected{pos, uint16(dir), expiry})
		}
	}
	before := len(p.proxy.life.created)
	audioBefore := len(core.PortTestCombatAudioSnapshot())
	got := uint32(C.nox_xxx_castPixies_540440(C.int(id), 0, inventoryInt(owner), inventoryInt(origin), 0, C.int(level)))
	created := p.proxy.life.created[before:]
	if got != 1 || len(created) != len(want) || core.Rand.Logic.Index() != rng.Index() {
		panic(fmt.Sprintf("pixie count/RNG: got %d/%d/%d want 1/%d/%d", got, len(created), core.Rand.Logic.Index(), len(want), rng.Index()))
	}
	wantCache := uint32(typeID)
	if typ == nil && !sp.WarmCache {
		wantCache = 0
	}
	if *cache != wantCache {
		panic("pixie type cache")
	}
	p.proxy.trace = append(p.proxy.trace, 702, *cache, uint32(len(created)))
	for i, u := range created {
		e := want[i]
		if u.PosVec != e.pos || *(*uint16)(unsafe.Add(u.CObj(), 124)) != e.dir || *(*uint16)(unsafe.Add(u.CObj(), 126)) != e.dir || *equipmentWord(u.CObj(), 80) != 0 || *equipmentWord(u.CObj(), 84) != 0 || *equipmentWord(u.CObj(), 156) != math.Float32bits(100) || *equipmentWord(u.CObj(), 160) != math.Float32bits(100) {
			panic("pixie position/direction/motion contract")
		}
		ud := u.UpdateData
		if *(**server.Object)(ud) != owner || *equipmentWord(ud, 12) != id || *equipmentWord(ud, 20) != e.expiry || *equipmentWord(ud, 24) != core.Frame() {
			panic("pixie update-data contract")
		}
		for _, v := range unsafe.Slice((*byte)(ud), 80)[64:] {
			if v != 0xa5 {
				panic("pixie update guard")
			}
		}
		for _, v := range unsafe.Slice((*uint32)(ud), 16) {
			p.proxy.trace = append(p.proxy.trace, p.normalize(v))
		}
	}
	events := core.PortTestCombatAudioSnapshot()[audioBefore:]
	if count < limit {
		if len(events) != 1 || uint32(events[0].ID) != 100+id || events[0].Obj != origin {
			panic("pixie audio")
		}
	} else if len(events) != 0 {
		panic("pixie unexpected audio")
	}
	return got
}

// PortTestSpellCharmControl checks the actual command setter and restores its owner.
func PortTestSpellCharmControl(values []bool) []int32 {
	old := C.nox_cheat_charmall
	defer func() { C.nox_cheat_charmall = old }()
	var out []int32
	for _, v := range values {
		CheatCharmAll(v)
		out = append(out, int32(C.nox_cheat_charmall))
	}
	return out
}
