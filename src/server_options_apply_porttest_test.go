//go:build porttest

package opennox

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/ifs"
	"github.com/opennox/libs/spell"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

type serverOptionsApplyServer struct {
	legacy.Server
	maps []string
}

func (s *serverOptionsApplyServer) SwitchMap(name string) { s.maps = append(s.maps, name) }

func TestServerOptionsApply(t *testing.T) {
	type row struct {
		Selected, Current string
		Cycle             bool
		Mode              uint16
		Maps, Commands    []string
		Written           string
		Open              bool
		Settings          []byte
	}
	var rows []row
	for _, selected := range []string{"", "Arena", "Other"} {
		for _, cycle := range []bool{false, true} {
			for _, mode := range []uint16{0x100, 0x1000} {
				t.Run(fmt.Sprintf("%s-%t-%x", selected, cycle, mode), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					defer noxflags.PortTestGameFlags(1)()
					t.Cleanup(legacy.PortTestServerOptionsRuleState())
					t.Cleanup(handles.PortTestInit())
					rules, _ := server.PortTestRuleServerSetup()
					oldSpells := o.c.srv.Spells
					o.c.srv.Spells = rules.Spells
					t.Cleanup(func() { o.c.srv.Spells = oldSpells })
					dir := t.TempDir()
					oldDir, err := ifs.Workdir()
					if err != nil {
						t.Fatal(err)
					}
					if err := ifs.Chdir(dir); err != nil {
						t.Fatal(err)
					}
					t.Cleanup(func() {
						if err := ifs.Chdir(oldDir); err != nil {
							t.Error(err)
						}
					})
					if err := os.MkdirAll(filepath.Join(dir, "maps", "Arena"), 0700); err != nil {
						t.Fatal(err)
					}
					for _, region := range [][3]int{{0x85B3FC, 36, 12}, {0x587000, 131668, 5}} {
						buf := unsafe.Slice((*byte)(memmap.PtrOff(uintptr(region[0]), uintptr(region[1]))), region[2])
						old := append([]byte(nil), buf...)
						clear(buf)
						t.Cleanup(func() { copy(buf, old) })
					}
					copy(unsafe.Slice(memmap.PtrUint8(0x85B3FC, 36), 12), "Arena")
					copy(unsafe.Slice(memmap.PtrUint8(0x587000, 131668), 5), ".map\x00")
					copy(unsafe.Slice(memmap.PtrUint8(0x5D4594, 2598188), 80), "maps\\Arena\\Arena.map")
					copy(unsafe.Slice(memmap.PtrUint8(0x5D4594, 3452), 12), "Cycle")
					copy(o.settings, selected)
					copy(o.settings[9:], "Port server")
					binary.LittleEndian.PutUint16(o.settings[52:], mode)
					binary.LittleEndian.PutUint16(o.settings[54:], 65535)
					for i := 24; i < 52; i++ {
						o.settings[i] = 255
					}
					index := int(spell.SPELL_FIREBALL)
					o.settings[24+index/8] &^= 1 << uint(index%8)
					if cycle {
						o.settings[57] = 1
					}
					copy(unsafe.Slice(memmap.PtrUint8(0x5D4594, 371438), 58), o.settings)
					*memmap.PtrUint32(0x5D4594, 371688) = 1
					prevServer := legacy.GetServer
					adapter := &serverOptionsApplyServer{Server: prevServer()}
					legacy.GetServer = func() legacy.Server { return adapter }
					t.Cleanup(func() { legacy.GetServer = prevServer })
					var commands []string
					oldExec := legacy.ExecConsoleCmd
					legacy.ExecConsoleCmd = func(_ context.Context, cmd string) bool { commands = append(commands, cmd); return true }
					t.Cleanup(func() { legacy.ExecConsoleCmd = oldExec })
					o.call("apply", 0, "")
					effective := selected
					if cycle {
						effective = "Cycle"
					}
					wantMaps := []string(nil)
					if effective != "" && effective != "Arena" {
						wantMaps = []string{effective + ".map"}
					}
					if fmt.Sprint(adapter.maps) != fmt.Sprint(wantMaps) {
						t.Fatalf("map requests %q want %q", adapter.maps, wantMaps)
					}
					open := *o.optionWords["root"] != 0
					if open != (effective == "") {
						t.Fatal("apply close state")
					}
					written, _ := os.ReadFile(filepath.Join(dir, "maps", "Arena", "user.rul"))
					if effective == "Arena" && !strings.Contains(string(written), "SPELL_FIREBALL") {
						t.Fatalf("rule writer did not preserve selected spell mask: %q", written)
					}
					if effective != "" {
						if legacy.Nox_xxx_serverOptionsGetServername_40A4C0() != "Port server" {
							t.Fatal("server name not committed")
						}
						if o.settings[57] != 0 {
							t.Fatal("committed cycle marker retained")
						}
					}
					rows = append(rows, row{selected, "Arena", cycle, mode, adapter.maps, commands, string(written), open, append([]byte(nil), o.settings...)})
				})
			}
		}
	}
	spellbookCapture(t, "server-options-apply", rows, "9c746eb7460e13076dc8ad71920f07c79f19dcce71adf5367d9efa5f0a93797b")
}
