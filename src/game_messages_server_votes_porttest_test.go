//go:build porttest

package opennox

import (
	"bytes"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerVoteWithdrawal(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	settings, restore := legacy.PortTestVoteSettings()
	t.Cleanup(restore)
	*settings[0], *settings[1] = 1, 1
	old := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(old)
	t.Cleanup(func() { noxflags.UnsetGamePlay(noxflags.GetGamePlay()); noxflags.SetGamePlay(old) })
	o.units[1].UpdateDataPlayer().Player.SetName("target")
	for i := range o.units {
		o.units[i].UpdateDataPlayer().Player.Field4792 = 1
	}
	u := &o.units[0]
	pl := u.UpdateDataPlayer().Player
	for _, quest := range []bool{false, true} {
		flags := noxflags.GameFlag(0)
		if quest {
			flags = noxflags.GameModeQuest
		}
		restore := noxflags.PortTestGameFlags(flags)
		for _, cast := range []byte{0, 1, 4} {
			data := make([]byte, 52)
			data[0], data[1] = 238, cast
			for i, c := range []byte("target") {
				data[2+i*2] = c
			}
			length, kind, withdraw := 52, uint32(cast), cast+2
			if cast == 0 && quest {
				kind = 3
			}
			if cast == 4 {
				length, kind, withdraw = 2, 2, 5
				data = data[:2]
			}
			send := func(action byte) {
				data[1] = action
				before := bytes.Clone(data)
				if n := legacy.Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(pl.PlayerIndex(), data, pl, u, nil); n != length || !bytes.Equal(data, before) {
					t.Fatal("vote length/input")
				}
			}
			send(cast)
			list := legacy.PortTestVoteList()
			if len(list) != 1 || *(*uint32)(list[0]) != kind {
				t.Fatalf("vote creation quest=%v kind=%d", quest, kind)
			}
			// A different participant withdrawing must leave the caster's vote intact.
			if kind == 2 {
				legacy.PortTestVoteCast(int(kind), unsafe.Pointer(&o.units[2]), nil, true)
			}
			if kind == 2 && len(legacy.PortTestVoteList()) != 1 {
				t.Fatal("nonvoter removed reset vote")
			}
			send(withdraw)
			if len(legacy.PortTestVoteList()) != 0 {
				t.Fatalf("vote withdrawal quest=%v kind=%d", quest, kind)
			}
			send(withdraw)
			if len(legacy.PortTestVoteList()) != 0 {
				t.Fatal("repeated withdrawal created vote")
			}
		}
		restore()
	}
}
