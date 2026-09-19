//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestPlayerDeathCoopPendingLoad(t *testing.T) {
	type row struct {
		Frame, Pending, Delay, Restore, State uint32
		Loading, Pointer                      bool
	}
	var rows []row
	for _, frame := range []uint32{0, 123, 0xffffffff} {
		t.Run(fmt.Sprint(frame), func(t *testing.T) {
			o := newMatchRosterOwner(t)
			oldAbilities := noxServer.abilities
			noxServer.abilities.Init(noxServer)
			t.Cleanup(func() { noxServer.abilities = oldAbilities })
			words, restore := legacy.PortTestServerOrchestrationGlobals()
			t.Cleanup(restore)
			serverConfigOwnBytes(t, 0x5D4594, 1563076, 4)
			oldPending, oldDelay, oldPointer, oldLoading := dword_5d4594_1563092, dword_5d4594_1563088, dword_5d4594_1563084, dword_5d4594_1563080
			t.Cleanup(func() {
				dword_5d4594_1563092, dword_5d4594_1563088, dword_5d4594_1563084, dword_5d4594_1563080 = oldPending, oldDelay, oldPointer, oldLoading
			})
			u := &o.units[0]
			u.ObjFlags = 0
			dword_5d4594_1563092 = 99
			dword_5d4594_1563088 = 77
			dword_5d4594_1563084 = u.CObj()
			dword_5d4594_1563080 = true
			*words["restore-cleanup"] = 1
			*memmap.PtrUint32(0x5D4594, 1563076) = 42
			noxflags.ResetGame()
			noxflags.SetGame(noxflags.GameHost | noxflags.GameFlag(2048))
			objectXferSetWord(u.CObj(), 520, 0)
			objectXferSetWord(u.CObj(), 524, 0)
			objectXferSetWord(u.UpdateDataPlayer().Player.C(), 3600, 0)
			o.reset()
			o.s.SetFrame(frame)
			ccall.CallVoidPtr(server.PortTestPlayerDeathCallback(), u.CObj())
			r := row{frame, dword_5d4594_1563092, dword_5d4594_1563088, *words["restore-cleanup"], *memmap.PtrUint32(0x5D4594, 1563076), dword_5d4594_1563080, dword_5d4594_1563084 != nil}
			if r != (row{Frame: frame, Delay: frame}) {
				t.Fatalf("pending load after death %+v", r)
			}
			rows = append(rows, r)
		})
	}
	spellbookCapture(t, "player-death-coop-pending-load", rows, "1c6082164c09eb67926cb79bb4f6bf3b3d5296ea9fae26db2780c8906e43a2e8")
}
