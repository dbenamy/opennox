//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestConsoleCommandsAdmission(t *testing.T) {
	type row struct {
		Command        string
		Quest, Missing bool
		Slot           int
		Calls          []string
		Blocked        []serverConfigListRow
		Output         []consoleCommandLine
	}
	var rows []row
	for _, command := range []string{"ban", "kick"} {
		for _, quest := range []bool{false, true} {
			for _, missing := range []bool{false, true} {
				for _, slot := range []int{1, 31} {
					t.Run(fmt.Sprintf("%s/%t/%t/%d", command, quest, missing, slot), func(t *testing.T) {
						o := newConsoleCommandOwner(t)
						units, configure, _, free := o.c.srv.S().PortTestEscortPlayers()
						t.Cleanup(free)
						configure(3)
						p := units[0].UpdateDataPlayer().Player
						if slot == 31 {
							p = units[2].UpdateDataPlayer().Player
						}
						p.SetName("target")
						clear(p.SerialBuf[:])
						copy(p.SerialBuf[:], "192.0.2.1")
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
						flags := noxflags.GameFlag(0)
						if quest {
							flags = 4096
						}
						t.Cleanup(noxflags.PortTestGameFlags(flags))
						var calls []string
						oldDisconn, oldCall, oldQuest := legacy.Nox_xxx_playerDisconnByPlrID_4DEB00, legacy.Nox_xxx_playerCallDisconnect_4DEAB0, legacy.Sub_4DCFB0
						t.Cleanup(func() {
							legacy.Nox_xxx_playerDisconnByPlrID_4DEB00 = oldDisconn
							legacy.Nox_xxx_playerCallDisconnect_4DEAB0 = oldCall
							legacy.Sub_4DCFB0 = oldQuest
						})
						legacy.Nox_xxx_playerDisconnByPlrID_4DEB00 = func(i ntype.PlayerInd) { calls = append(calls, fmt.Sprintf("disconnect:%d", i)) }
						legacy.Nox_xxx_playerCallDisconnect_4DEAB0 = func(i ntype.PlayerInd, r int8) { calls = append(calls, fmt.Sprintf("request:%d:%d", i, r)) }
						legacy.Sub_4DCFB0 = func(u *server.Object) {
							if u != p.PlayerUnit {
								t.Fatal("quest removal target")
							}
							calls = append(calls, "quest")
						}
						name := "target"
						if missing {
							name = "missing"
						}
						if !o.call(t, command, false, name) {
							t.Fatal("admission result")
						}
						blocked, _ := serverConfigListSnapshot(t, "blocked")
						wantCalls := 0
						if !missing && slot != 31 {
							wantCalls = 1
						}
						if len(calls) != wantCalls {
							t.Fatal(calls, wantCalls)
						}
						if wantCalls == 1 {
							want := "quest"
							if !quest {
								if command == "ban" {
									want = "disconnect:1"
								} else {
									want = "request:1:4"
								}
							}
							if calls[0] != want {
								t.Fatal(calls, want)
							}
						}
						wantBlocked := 0
						if command == "ban" && (missing || slot != 31) {
							wantBlocked = 1
						}
						if len(blocked) != wantBlocked {
							t.Fatal(blocked, wantBlocked)
						}
						if wantBlocked == 1 {
							address := "192.0.2.1"
							if missing {
								address = ""
							}
							if blocked[0].Name != name || blocked[0].Address != address || blocked[0].Expires != 0 {
								t.Fatal(blocked)
							}
						}
						rows = append(rows, row{command, quest, missing, slot, calls, blocked, o.printer.lines})
					})
				}
			}
		}
	}
	spellbookCapture(t, "console-commands-admission", rows, "62622506e37204de18dabfc0de05cc7311ed96ae639ca435d9d78d72eaf5e596")
}
