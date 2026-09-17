//go:build porttest

package opennox

import (
	"bytes"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

type matchRosterOwner struct {
	*questRuntimeOwner
	roster map[string]*uint32
}

func newMatchRosterOwner(t *testing.T) *matchRosterOwner {
	t.Helper()
	o := &matchRosterOwner{questRuntimeOwner: newQuestRuntimeOwner(t)}
	t.Cleanup(noxflags.PortTestGameFlags(0))
	play := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
	t.Cleanup(func() { noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0)); noxflags.SetGamePlay(play) })
	engine := noxflags.GetEngine()
	t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(engine) })
	noxflags.ResetEngine()
	t.Cleanup(o.s.PortTestObjectiveTypes(nil, nil, nil))
	t.Cleanup(o.s.PortTestMapDrawableTeamMessages())
	oldSend := o.s.NetSendPacketXxx
	o.s.NetSendPacketXxx = legacy.Nox_xxx_netSendPacket_4E5030
	t.Cleanup(func() { o.s.NetSendPacketXxx = oldSend })
	ball := legacy.PortTestMapDrawableTeamWord()
	oldBall := *ball
	*ball = 0
	t.Cleanup(func() { *ball = oldBall })
	var restore func()
	o.roster, restore = legacy.PortTestMatchRosterGlobals()
	t.Cleanup(restore)
	for _, r := range []struct {
		base, off uintptr
		size      int
	}{
		{0x5D4594, 1324, 16}, {0x5D4594, 3464, 4}, {0x5D4594, 3468, 8},
		{0x5D4594, 3488, 24}, {0x5D4594, 371380, 58}, {0x5D4594, 371516, 144},
		{0x5D4594, 1563272, 4}, {0x5D4594, 1563284, 4}, {0x5D4594, 1567736, 106},
		{0x5D4594, 371688, 4}, {0x587000, 4660, 4}, {0x587000, 4704, 24},
	} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r.base, r.off)), r.size)
		old := bytes.Clone(b)
		clear(b)
		t.Cleanup(func() { copy(b, old) })
	}
	for _, table := range blobdata.PortTestScoreboardTables() {
		if table.Base == 0x587000 && table.Offset == 4704 {
			copy(unsafe.Slice((*byte)(memmap.PtrOff(table.Base, table.Offset)), len(table.Data)), table.Data)
		}
	}
	flag := noxServer.flag3592
	t.Cleanup(func() { noxServer.flag3592 = flag })
	noxServer.flag3592 = false
	return o
}
func matchRosterCall(op string, u *server.Object, pl *server.Player, data unsafe.Pointer, args ...uint32) uint32 {
	var a [4]uint32
	copy(a[:], args)
	return legacy.PortTestMatchRoster(op, u, pl, data, a)
}
