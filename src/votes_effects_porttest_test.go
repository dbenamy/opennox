//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestPortVotesThresholdEffects(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	t.Cleanup(noxflags.PortTestGameFlags(0))
	// Observe the real blocked-player list; disconnect and quest-removal requests
	// use existing service hooks because their implementations are outside scope.
	h := unsafe.Slice(memmap.PtrUint32(0x5D4594, 371500), 3)
	oldHead := append([]uint32(nil), h...)
	addr := uint32(uintptr(unsafe.Pointer(&h[0])))
	h[0], h[1], h[2] = addr, addr, addr
	t.Cleanup(func() {
		for legacy.PortTestServerConfigAdmission("blocked-first", 0, nil, nil) != nil {
			legacy.PortTestServerConfigAdmission("blocked-remove", 0, nil, nil)
		}
		copy(h, oldHead)
	})
	oldTicks := legacy.PlatformTicks
	legacy.PlatformTicks = func() uint64 { return 1000 }
	t.Cleanup(func() { legacy.PlatformTicks = oldTicks })
	var disconnected []int
	oldDisconnect := legacy.Nox_xxx_playerCallDisconnect_4DEAB0
	legacy.Nox_xxx_playerCallDisconnect_4DEAB0 = func(ind ntype.PlayerInd, reason int8) {
		if reason != 4 {
			t.Fatalf("disconnect reason %d", reason)
		}
		disconnected = append(disconnected, int(ind))
	}
	t.Cleanup(func() { legacy.Nox_xxx_playerCallDisconnect_4DEAB0 = oldDisconnect })
	var questRemoved []*server.Object
	oldRemove := legacy.Sub_4DCFB0
	legacy.Sub_4DCFB0 = func(u *server.Object) { questRemoved = append(questRemoved, u) }
	t.Cleanup(func() { legacy.Sub_4DCFB0 = oldRemove })
	target := &o.units[1]
	target.UpdateDataPlayer().Player.SetName("target")
	target.UpdateDataPlayer().Player.Field4792 = 1
	for _, kind := range []int{0, 1, 3} {
		p := legacy.PortTestVoteCreate(kind, o.units[0].CObj())
		*(*unsafe.Pointer)(unsafe.Add(p, 28)) = target.CObj()
		*(*byte)(unsafe.Add(p, 4)) = 1
		*(*byte)(unsafe.Add(p, 12)) = 1
		legacy.Nox_xxx_voteUptate_506F30()
		if len(legacy.PortTestVoteList()) != 0 {
			t.Fatalf("kind%d vote not removed", kind)
		}
	}
	if len(disconnected) != 2 || disconnected[0] != 7 || disconnected[1] != 7 || len(questRemoved) != 1 || questRemoved[0] != target {
		t.Fatal("vote action dispatch")
	}
	rows, _ := serverConfigListSnapshot(t, "blocked")
	if len(rows) != 3 {
		t.Fatalf("blocked records %d", len(rows))
	}
	for _, r := range rows {
		if r.Name != "target" || r.Expires != 901000 {
			t.Fatalf("blocked record %+v", r)
		}
	}
	spellbookCapture(t, "votes-effects", rows, "e4cfcb9483e31973688022c4942d0fb8e61e4c4f199cc66b029a6af59b913978")
}

func TestPortVotesQuestReset(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameModeQuest))
	oldServer := legacy.GetServer
	catalog := legacy.PortTestMapCatalogOpen(1)
	t.Cleanup(catalog.Close)
	observer := &serverOptionsApplyServer{Server: oldServer()}
	legacy.GetServer = func() legacy.Server { return observer }
	t.Cleanup(func() { legacy.GetServer = oldServer })
	state := catalog.State()
	state.Count = 1
	clear(state.Quest)
	copy(state.Quest[4:], "quest-start.map")
	catalog.RestoreState(state)
	stage := memmap.PtrUint32(0x587000, 202028)
	oldStage := *stage
	t.Cleanup(func() { *stage = oldStage })
	*stage = 17
	for i := range o.units {
		o.units[i].UpdateDataPlayer().Player.Field4792 = 1
		// Exercise the actual respawn entry's quest-transition guard. Respawn's
		// full equipment/positioning behavior has its own accumulated contracts.
		*(*uint32)(unsafe.Add(o.units[i].UpdateData, 548)) = 1
	}
	p := legacy.PortTestVoteCreate(2, o.units[0].CObj())
	*(*byte)(unsafe.Add(p, 4)) = 2
	*(*uint32)(unsafe.Add(p, 8)) = 0x80000082
	legacy.Nox_xxx_voteUptate_506F30()
	if len(legacy.PortTestVoteList()) != 1 || len(observer.maps) != 0 || *stage != 17 {
		t.Fatal("reset before all participants voted")
	}
	*(*byte)(unsafe.Add(p, 4)) = 3
	legacy.Nox_xxx_voteUptate_506F30()
	if len(legacy.PortTestVoteList()) != 0 || len(observer.maps) != 1 || observer.maps[0] != "quest-start.map" || *stage != 0 {
		t.Fatalf("quest reset effects: records=%d maps=%q stage=%d", len(legacy.PortTestVoteList()), observer.maps, *stage)
	}
	nodes := o.state().Nodes
	var notified []byte
	for _, n := range nodes {
		if len(n.Data) == 2 && n.Data[0] == 238 && n.Data[1] == 7 {
			notified = append(notified, n.To)
		}
	}
	if len(notified) != 3 || notified[0] != 31 || notified[1] != 7 || notified[2] != 1 {
		t.Fatalf("reset voter notifications %v", notified)
	}
	spellbookCapture(t, "votes-quest-reset", struct {
		Maps     []string
		Stage    uint32
		Notified []byte
	}{observer.maps, *stage, notified}, "5998070c2b4e57e0758d25eb8150230ab640e170add7b91655cdd14ef5f764cc")
}
