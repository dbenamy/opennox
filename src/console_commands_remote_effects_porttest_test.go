//go:build porttest

package opennox

import (
	"bytes"
	"context"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
)

func consoleCommandRemoteOwner(t *testing.T) (*consoleCommandOwner, []server.Object) {
	o := newConsoleCommandOwner(t)
	units, configure, _, free := o.c.srv.S().PortTestEscortPlayers()
	t.Cleanup(free)
	configure(3)
	for i := range units {
		units[i].UpdateDataPlayer().Player.SetName(fmt.Sprintf("player%d", i))
	}
	serverConfigOwnBytes(t, 0x5D4594, 818228, 1024)
	return o, units
}
func TestConsoleCommandsRemoteScriptAndExecution(t *testing.T) {
	type row struct {
		Authorization, Action int
		Text                  string
		Called                bool
		Output                []consoleCommandLine
	}
	var rows []row
	for auth := 0; auth < 4; auth++ {
		for _, action := range []int{1, 2} {
			for _, text := range []string{"exe Target", "exe Missing"} {
				t.Run(fmt.Sprintf("%d/%d/%s", auth, action, text), func(t *testing.T) {
					o, units := consoleCommandRemoteOwner(t)
					pl := units[0].UpdateDataPlayer().Player
					oldCode := legacy.ClientPlayerNetCode()
					legacy.ClientSetPlayerNetCode(23)
					t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
					if auth == 1 {
						legacy.ClientSetPlayerNetCode(int(pl.PlayerInd))
					}
					if auth == 2 {
						t.Cleanup(noxflags.PortTestGameFlags(2048))
					}
					t.Cleanup(legacy.PortTestSessionEntryRosterOwner())
					if auth == 3 {
						legacy.PortTestSessionEntryRoster("add", int(pl.PlayerInd))
					}
					s := o.c.srv.S()
					oldVM := s.NoxScriptVM
					s.NoxScriptVM = server.NoxScriptVM{}
					s.NoxScriptVM.Init(s)
					t.Cleanup(func() { s.NoxScriptVM.Reset(); s.NoxScriptVM = oldVM })
					funcs := []prefabScriptFunction{{Name: "GLOBAL", Code: []uint32{72}}, {Name: "GLOBAL", Code: []uint32{72}}, {Name: "Target", Code: []uint32{72}}}
					if err := s.NoxScriptVM.ReadScript(bytes.NewReader(prefabScriptsEncode(nil, funcs))); err != nil {
						t.Fatal(err)
					}
					called := false
					oldExec := legacy.ExecConsoleCmd
					legacy.ExecConsoleCmd = func(_ context.Context, cmd string) bool {
						if cmd != text || legacy.PortTestConsoleSender() != pl.C() {
							t.Fatal("execution context")
						}
						called = true
						return true
					}
					t.Cleanup(func() { legacy.ExecConsoleCmd = oldExec })
					wide, freeText := alloc.CString16(text)
					defer freeText()
					if legacy.PortTestConsoleRemote(pl.C(), action, wide) != 1 || legacy.PortTestConsoleSender() != nil {
						t.Fatal("dispatch completion")
					}
					if action == 1 {
						called = s.NoxScriptVM.Caller() == &units[0] && s.NoxScriptVM.Trigger() == &units[0]
					}
					want := auth != 0 && (action == 2 || text == "exe Target")
					if called != want {
						t.Fatal("dispatch", called, want)
					}
					if auth == 0 && len(o.printer.lines) != 0 {
						t.Fatal("unauthorized output", o.printer.lines)
					}
					rows = append(rows, row{auth, action, text, called, o.printer.lines})
				})
			}
		}
	}
	spellbookCapture(t, "console-commands-remote-script", rows, "75cfd70bc551bcda8264b371272d410f05cef8e685f36915f6a48243bcd0713a")
}
func TestConsoleCommandsRemoteCamera(t *testing.T) {
	type row struct {
		Observer, Replay bool
		Text             string
		Target           int
	}
	var rows []row
	for _, observer := range []bool{false, true} {
		for _, replay := range []bool{false, true} {
			for _, text := range []string{"", "PLAYER1", "missing"} {
				t.Run(fmt.Sprintf("%t/%t/%s", observer, replay, text), func(t *testing.T) {
					_, units := consoleCommandRemoteOwner(t)
					p := units[0].UpdateDataPlayer().Player
					p.CameraFollowObj = &units[2]
					if observer {
						p.Field3680 = 1
					}
					oldEngine := noxflags.GetEngine()
					if replay {
						noxflags.SetEngine(noxflags.EngineReplayRead)
					} else {
						noxflags.UnsetEngine(noxflags.EngineReplayRead)
					}
					t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) })
					wide, free := alloc.CString16(text)
					defer free()
					if legacy.PortTestConsoleRemote(p.C(), 4, wide) != 1 || legacy.PortTestConsoleSender() != nil {
						t.Fatal("camera completion")
					}
					want := &units[2]
					index := 2
					if observer || replay {
						switch text {
						case "":
							want = nil
							index = -1
						case "PLAYER1":
							want = &units[1]
							index = 1
						}
					}
					if p.CameraFollowObj != want {
						t.Fatal("camera target")
					}
					rows = append(rows, row{observer, replay, text, index})
				})
			}
		}
	}
	spellbookCapture(t, "console-commands-remote-camera", rows, "8a273cb509faf709f3c51bfc4de5af18d60c1dad61fe607d9895b59911eab476")
}
func TestConsoleCommandsRemoteBroadcast(t *testing.T) {
	type row struct {
		Action  int
		Packets [][]byte
		Sounds  []int
	}
	var rows []row
	for _, action := range []int{3, 5} {
		t.Run(fmt.Sprint(action), func(t *testing.T) {
			o, units := consoleCommandRemoteOwner(t)
			s := o.c.srv.S()
			s.PortTestCombatAudioReset()
			t.Cleanup(s.PortTestCombatAudioReset)
			text, free := alloc.CString16("hi")
			defer free()
			if legacy.PortTestConsoleRemote(units[0].UpdateDataPlayer().Player.C(), action, text) != 1 {
				t.Fatal("broadcast return")
			}
			packets := make([][]byte, 32)
			for i := range packets {
				packets[i] = s.NetList.CopyPacketsA(ntype.PlayerInd(i), netlist.Kind1)
			}
			flag := byte(2)
			if action == 5 {
				flag |= 16
			}
			want := []byte{168, 0, 0, flag, 0, 0, 0, 0, 3, 0, 0, 'h', 'i', 0}
			for i, p := range packets {
				expected := []byte(nil)
				if i == 1 || i == 7 || i == 31 {
					expected = want
				}
				if !bytes.Equal(p, expected) {
					t.Fatal(i, p, expected)
				}
			}
			var sounds []int
			events := s.PortTestCombatAudioSnapshot()
			for i, e := range events {
				if e.ID != 902 || e.Obj != &units[i] || e.Kind != 0 || e.Code != 0 {
					t.Fatal("audio event", e)
				}
				sounds = append(sounds, int(e.ID))
			}
			count := 0
			if action == 5 {
				count = 3
			}
			if len(events) != count || legacy.PortTestConsoleSender() != nil {
				t.Fatal("broadcast completion")
			}
			rows = append(rows, row{action, packets, sounds})
		})
	}
	spellbookCapture(t, "console-commands-remote-broadcast", rows, "8f22016dfb993207ed9451162533808ef2acbb9eabe63e173c19ed8a48cb85f6")
}
func TestConsoleCommandsRemoteObserver(t *testing.T) {
	type row struct {
		Flags, Status   uint32
		Special, Accept bool
		Calls           [][2]int
		After           uint32
	}
	var rows []row
	for _, flags := range []uint32{0, 8, 4096} {
		for _, status := range []uint32{0, 1} {
			for _, special := range []bool{false, true} {
				for _, accept := range []bool{false, true} {
					t.Run(fmt.Sprintf("%x/%x/%t/%t", flags, status, special, accept), func(t *testing.T) {
						o, units := consoleCommandRemoteOwner(t)
						p := units[0].UpdateDataPlayer().Player
						p.Field3680 = status
						t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
						var calls [][2]int
						old := legacy.Nox_xxx_playerGoObserver_4E6860
						legacy.Nox_xxx_playerGoObserver_4E6860 = func(pl *server.Player, a, b int) int {
							if pl != p || legacy.PortTestConsoleSender() != p.C() {
								t.Fatal("observer context")
							}
							calls = append(calls, [2]int{a, b})
							if accept {
								return 1
							}
							return 0
						}
						t.Cleanup(func() { legacy.Nox_xxx_playerGoObserver_4E6860 = old })
						text := "observer"
						if special {
							text = "\uf00d\uf0ad"
						}
						wide, free := alloc.CString16(text)
						defer free()
						if legacy.PortTestConsoleRemote(p.C(), 0, wide) != 1 || legacy.PortTestConsoleSender() != nil {
							t.Fatal("observer completion")
						}
						var want [][2]int
						if flags == 0 && status == 0 {
							want = [][2]int{{bool2int(special), 0}}
						}
						if !reflect.DeepEqual(calls, want) {
							t.Fatal(calls, want)
						}
						if len(want) == 0 || !accept {
							if p.Field3680 != status || len(o.printer.lines) != 0 {
								t.Fatal("rejected observer mutated state")
							}
						}
						rows = append(rows, row{flags, status, special, accept, calls, p.Field3680})
					})
				}
			}
		}
	}
	spellbookCapture(t, "console-commands-remote-observer", rows, "41bc2c18a58c2013c5d41265675de841bcff39cf374d060f65f3d4942fa4984d")
}

// The remaining C quit-dialog caller supplies no command text for this action.
func TestConsoleCommandsQuitObserver(t *testing.T) {
	_, units := consoleCommandRemoteOwner(t)
	p := units[0].UpdateDataPlayer().Player
	p.Field3680 = 0
	t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameHost))
	old := legacy.Nox_xxx_playerGoObserver_4E6860
	t.Cleanup(func() { legacy.Nox_xxx_playerGoObserver_4E6860 = old })
	calls := 0
	legacy.Nox_xxx_playerGoObserver_4E6860 = func(pl *server.Player, force, extra int) int {
		calls++
		if pl != p || force != 0 || extra != 0 || legacy.PortTestConsoleSender() != p.C() {
			t.Fatal("quit dialog observer arguments/context changed")
		}
		return 0
	}
	if legacy.PortTestConsoleRemote(p.C(), 0, nil) != 1 || calls != 1 || legacy.PortTestConsoleSender() != nil {
		t.Fatal("quit dialog observer completion changed")
	}
}
