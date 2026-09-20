//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionGauntletStart(t *testing.T) {
	q := newQuickbarOwner(t)
	q.ownFullQuickbar(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	pathField := serverConfigOwnBytes(t, 0x85B3FC, 10984, 1024)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	t.Cleanup(legacy.PortTestPlayerFileClientSections())
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) })
	rank, restoreRank := legacy.PortTestScoreboardWords()
	t.Cleanup(restoreRank)
	vote, restoreVote := legacy.PortTestVoteGUIOwner()
	t.Cleanup(restoreVote)
	dir := t.TempDir()
	existing := filepath.Join(dir, "empty.plr")
	missing := filepath.Join(dir, "missing.plr")
	f, err := cryptfile.OpenFile(existing, cryptfile.WriteOnly, 27)
	if err != nil {
		t.Fatal(err)
	}
	if err = f.WriteU32(0); err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
	// Observe the already-qualified file-transfer boundary. File loading, quickbar,
	// meter creation and vote reset below remain the actual production owners.
	var transfers []string
	oldTransfer := legacy.Sub_41CC00
	legacy.Sub_41CC00 = func(path string) { transfers = append(transfers, filepath.Base(path)) }
	t.Cleanup(func() { legacy.Sub_41CC00 = oldTransfer })
	type state struct {
		Bar              quickbarResult
		Meters           meterResult
		Choice, Previous uint32
		Transfers        []string
	}
	snapshot := func() state {
		return state{q.snapshot("quest start", 0), q.meterOwner.invoke(t, 0, 0, 0, nil, 0, 0, 0, 0), *vote["choice"], *vote["previousChoice"], append([]string(nil), transfers...)}
	}
	type row struct {
		On, Host, Quest, Class, File int
		State                        state
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			for quest := 0; quest < 2; quest++ {
				for class := 0; class < 3; class++ {
					for file := 0; file < 2; file++ {
						path := missing
						if file != 0 {
							path = existing
						}
						setup := func() {
							q.reset(t)
							q.named = nil
							// The weapon HUD is not installed at this entry boundary. Avoid seeding
							// unrelated fixture windows that the meter rebuild would destroy and retain.
							for _, i := range []int{4, 5, 6} {
								q.meters.Records[i].Window = nil
							}
							if class == 0 {
								q.meters.Records[1].Window = nil
								q.meters.Records[3].Window = nil
							}
							noxflags.ResetGame()
							noxflags.SetGame(noxflags.GameFlag(host))
							if quest != 0 {
								noxflags.SetGame(4096)
							}
							*(*byte)(unsafe.Add(q.players[0].C(), 2251)) = byte(class)
							binary.LittleEndian.PutUint32(connected, uint32(on))
							clear(pathField)
							copy(pathField, path)
							*vote["choice"], *vote["previousChoice"] = 17, 29
							*rank["dword_5d4594_1090048"] = 0
							*rank["dword_5d4594_1090120"] = 0
							for i := 0; i < 25; i++ {
								q.bar[2*i], q.bar[2*i+1] = uint32(i%5+1), 0xabcd0000+uint32(i)
							}
							transfers = nil
						}
						setup()
						if on != 0 {
							if host != 0 {
								q.call("sub_460380")
								q.call("nox_xxx_cliPrepareGameplay1_460E60")
								legacy.Sub_41CC00(path)
								legacy.PortTestPlayerFileCall("nox_xxx_plrLoad_41A480", uint32(uintptr(unsafe.Pointer(&pathField[0]))))
							}
							legacy.PortTestMeterCall(1, nil, 1, 0, 0, 0)
							legacy.PortTestVoteGUI("reset-choice", 0)
							legacy.PortTestMeterCall(27, nil, 0, 0, 0, 0)
							legacy.PortTestScoreboard(26, 0, 0, 0)
						}
						want := snapshot()
						setup()
						data := []byte{240, 0}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
						got := snapshot()
						if n != 2 || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
							interactionCapture(t, "quest-start-diff", []state{want, got})
							t.Fatal("quest start route", on, host, quest, class, file, n)
						}
						rows = append(rows, row{on, host, quest, class, file, got})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-start", rows)
}
