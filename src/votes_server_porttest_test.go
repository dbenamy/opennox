//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func voteRecordWords(p unsafe.Pointer, units map[unsafe.Pointer]uint32) [13]uint32 {
	w := *(*[13]uint32)(p)
	for _, off := range []int{4, 7} {
		ptr := *(*unsafe.Pointer)(unsafe.Add(p, off*4))
		if ptr == nil {
			continue
		}
		if off == 4 {
			ptr = unsafe.Add(ptr, -48)
		}
		n, ok := units[ptr]
		if !ok {
			panic("unknown vote object")
		}
		w[off] = n
	}
	w[11], w[12] = 0, 0
	return w
}

func TestPortVotesAdmissionMatrix(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	old := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(old)
	t.Cleanup(func() { noxflags.UnsetGamePlay(noxflags.GetGamePlay()); noxflags.SetGamePlay(old) })
	settings, restore := legacy.PortTestVoteSettings()
	t.Cleanup(restore)
	var rows []struct {
		Name    string
		Records [][13]uint32
	}
	defer func() {
		spellbookCapture(t, "votes-admission", rows, "40d456dc99f22a8690abce2cb475e27fd638142f58c0c64a82aa69d9247c1850")
	}()
	for i := range o.units {
		o.units[i].UpdateDataPlayer().Player.SetName(fmt.Sprintf("voter%d", i))
		o.units[i].TeamVal.ID = 1
	}
	for _, kind := range []int{0, 1, 2, 3, 9} {
		for _, enabled := range []uint32{0, 1, 5, 32, 33, 256} {
			for _, participation := range []uint32{0, 1, 2} {
				for _, target := range []int{0, 1, 2, 3} {
					for len(legacy.PortTestVoteList()) != 0 {
						legacy.PortTestVoteDelete(legacy.PortTestVoteList()[0])
					}
					o.reset()
					*memmap.PtrUint32(0x587000, 229980) = enabled
					*settings[0], *settings[1] = enabled, enabled
					for i := range o.units {
						o.units[i].UpdateDataPlayer().Player.Field4792 = participation
					}
					name := fmt.Sprintf("voter%d", target)
					data := make([]byte, 52)
					data[0] = 238
					switch kind {
					case 0, 3:
						data[1] = 0
					case 1:
						data[1] = 1
					case 2:
						data[1] = 4
					default:
						data[1] = 99
					}
					for i, b := range []byte(name) {
						data[2+i*2] = b
					}
					flags := noxflags.GameFlag(0)
					if kind == 3 {
						flags = noxflags.GameModeQuest
					}
					restoreFlags := noxflags.PortTestGameFlags(flags)
					caster := &o.units[0]
					n := legacy.Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(caster.UpdateDataPlayer().Player.PlayerIndex(), data, caster.UpdateDataPlayer().Player, caster, nil)
					restoreFlags()
					wantSize := 52
					if kind == 2 {
						wantSize = 2
					}
					if kind == 9 {
						wantSize = -1
					}
					if n != wantSize {
						t.Fatalf("decoder kind %d returned %d", kind, n)
					}
					want := enabled > 0 && (kind == 2 || target == 1) && kind != 9
					if kind < 2 {
						want = want && enabled <= 32
					}
					if kind == 2 || kind == 3 {
						want = want && participation != 0
					}
					list := legacy.PortTestVoteList()
					if (len(list) != 0) != want {
						t.Fatalf("kind%d enabled%d participation%d target%d: %d records", kind, enabled, participation, target, len(list))
					}
					row := struct {
						Name    string
						Records [][13]uint32
					}{Name: fmt.Sprintf("kind%d/enabled%d/participation%d/target%d", kind, enabled, participation, target)}
					for _, p := range list {
						row.Records = append(row.Records, voteRecordWords(p, o.objects))
					}
					rows = append(rows, row)
				}
			}
		}
	}
}

func TestPortVotesDeletedTargetTick(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	target := &o.units[1]
	target.ObjFlags |= object.FlagDestroyed
	for _, kind := range []int{0, 1, 3} {
		p := legacy.PortTestVoteCreate(kind, o.units[0].CObj())
		*(*unsafe.Pointer)(unsafe.Add(p, 28)) = target.CObj()
	}
	// Tick dispatch must follow saved successors through consecutive deletions.
	legacy.Nox_xxx_voteUptate_506F30()
	if len(legacy.PortTestVoteList()) != 0 {
		t.Fatal("deleted target retained")
	}
}

func TestPortVotesAlternateWithdrawal(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	t.Cleanup(noxflags.PortTestGameFlags(0))
	old := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(old)
	t.Cleanup(func() { noxflags.UnsetGamePlay(noxflags.GetGamePlay()); noxflags.SetGamePlay(old) })
	o.units[1].UpdateDataPlayer().Player.SetName("target")
	caster := &o.units[0]
	pl := caster.UpdateDataPlayer().Player
	msg := make([]byte, 52)
	msg[0] = 238
	for i, c := range []byte("target") {
		msg[2+i*2] = c
	}
	send := func(action byte) {
		msg[1] = action
		if n := legacy.Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(pl.PlayerIndex(), msg, pl, caster, nil); n != 52 {
			t.Fatal(n)
		}
	}
	send(0)
	send(1)
	if len(legacy.PortTestVoteList()) != 2 {
		t.Fatal("separate vote kinds")
	}
	send(3)
	list := legacy.PortTestVoteList()
	if len(list) != 1 || *(*uint32)(list[0]) != 0 {
		t.Fatal("withdrawing kind1 changed the wrong vote")
	}
	send(2)
	if len(legacy.PortTestVoteList()) != 0 {
		t.Fatal("normal withdrawal retained vote")
	}
}
