//go:build porttest

package legacy

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func unitGameplayTimerContract(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestMonsterStateSpec) uint64 {
	off := 532
	if sp.Op == 29 {
		off = 520
	}
	word := (*uint32)(unsafe.Add(u.UpdateData, off))
	*word = sp.Deadline
	if !sp.AnimData {
		u.UpdateDataMonster().SoundSet122 = nil
	}
	oldClass := u.ObjClass
	defer func() { u.ObjClass = oldClass }()
	if sp.NilUnit {
		u.ObjClass = 0
	} // Valid non-monster record.
	before := append([]byte(nil), unsafe.Slice((*byte)(u.UpdateData), 2200)...)
	rng := proxy.core.Rand.Logic.Index()
	expectedRNG := *proxy.core.Rand.Logic
	expectedDelay := uint32(expectedRNG.IntClamp(int(2*sp.FPS), int(4*sp.FPS)))
	if sp.Op == 29 {
		unitDamageTimer(u)
		want := sp.Deadline
		if want == 0 {
			want = sp.Frame
		}
		if *word != want {
			panic("damage timer first write")
		}
		proxy.core.SetFrame(sp.Frame + 1)
		unitDamageTimer(u)
		if want == 0 {
			want = sp.Frame + 1
		}
		if *word != want {
			panic("damage timer repeated write")
		}
	} else {
		// The original char result is unused by every production caller.
		unitHurtSound(u)
		active := !sp.NilUnit && sp.Frame >= sp.Deadline
		events := proxy.core.PortTestCombatAudioSnapshot()
		wantAudio := 0
		if active {
			delta := *word - sp.Frame
			if delta < 2*sp.FPS || delta > 4*sp.FPS {
				panic(fmt.Sprintf("hurt interval %d", delta))
			}
			if proxy.core.Rand.Logic.Index() != expectedRNG.Index() || delta != expectedDelay {
				panic("hurt timer RNG value/consumption")
			}
			if sp.AnimData {
				wantAudio = 1
			}
		} else if *word != sp.Deadline || proxy.core.Rand.Logic.Index() != rng {
			panic("inactive hurt timer mutation")
		}
		if len(events) != wantAudio {
			panic(fmt.Sprintf("hurt audio count %d want %d", len(events), wantAudio))
		}
		if wantAudio != 0 && (events[0].ID != 302 || events[0].Obj != u || events[0].Kind != 0 || events[0].Code != 0) {
			panic("hurt audio fields")
		}
	}
	got := *word
	*word = sp.Deadline
	if !bytes.Equal(before, unsafe.Slice((*byte)(u.UpdateData), 2200)) {
		panic("timer changed unrelated update data")
	}
	*word = got
	return uint64(got)
}
