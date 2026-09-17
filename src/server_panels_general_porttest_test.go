//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestServerPanelsGeneralToggles(t *testing.T) {
	type row struct {
		Flags                uint32
		Initial              int
		ID                   uint
		Respawn, Motd, Cycle int
		Subflags, Updated    uint32
		Checked, Enabled     []bool
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 32, 1024, 4096} {
		for _, initial := range []int{0, 1} {
			for _, id := range []uint{10301, 10302, 10304, 10305, 10306} {
				t.Run(fmt.Sprintf("%x-%d-%d", flags, initial, id), func(t *testing.T) {
					catalog := legacy.PortTestMapCatalogOpen(29)
					t.Cleanup(catalog.Close)
					o := newServerOptionsOwner(t)
					o.installSubpanels(t, true)
					defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
					globals, restore := legacy.PortTestMatchRosterGlobals()
					t.Cleanup(restore)
					oldRespawn := *memmap.PtrUint32(0x5D4594, 3584)
					oldMotd := legacy.Get_nox_server_sendMotd_108752()
					t.Cleanup(func() { *memmap.PtrUint32(0x5D4594, 3584) = oldRespawn; legacy.Set_nox_server_sendMotd_108752(oldMotd) })
					legacy.Nox_xxx_ruleSetNoRespawn_40A5E0(initial)
					legacy.Set_nox_server_sendMotd_108752(initial)
					legacy.Sub_4D0D90(initial)
					*globals["server-subflags"] = 0x80000000 | uint32(initial)*0x2002
					raw := legacy.PortTestServerPanelsConstruct("general", o.options, unsafe.Pointer(&o.settings[0]))
					if raw == 0 {
						t.Fatal("general constructor")
					}
					w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw))))
					if legacy.PortTestServerPanelsConstruct("general", o.options, unsafe.Pointer(&o.settings[0])) != 0 {
						t.Fatal("duplicate general constructor")
					}
					r := row{Flags: flags, Initial: initial, ID: id}
					respawn := initial != 0 && flags&4096 == 0
					for _, wid := range []uint{10301, 10302, 10304, 10305, 10306} {
						c := w.ChildByID(wid)
						checked := c.DrawData().Field0&4 != 0
						wantChecked := initial != 0
						if wid == 10301 {
							wantChecked = respawn
						}
						if wid == 10306 && flags&1056 != 0 {
							wantChecked = true
						}
						if checked != wantChecked {
							t.Fatal("initial checkbox", wid, checked, wantChecked)
						}
						if wid == 10301 && c.Flags.IsEnabled() != (flags&1024 == 0 || !respawn) {
							t.Fatal("respawn disabled")
						}
						if wid == 10306 && c.Flags.IsEnabled() != (flags&1056 == 0) {
							t.Fatal("report disabled")
						}
						r.Checked = append(r.Checked, checked)
						r.Enabled = append(r.Enabled, c.Flags.IsEnabled())
					}
					*o.optionWords["settings-updated"] = 0
					child := w.ChildByID(id)
					if gui.EventRespInt(w.Func94(gui.AsWindowEvent(16391, uintptr(unsafe.Pointer(child)), 0))) != 0 {
						t.Fatal("general event")
					}
					wantRespawn, wantMotd, wantCycle := initial, initial, initial
					wantFlags := uint32(0x80000000) | uint32(initial)*0x2002
					wantUpdated := uint32(0)
					switch id {
					case 10301:
						wantRespawn = 1
						if respawn {
							wantRespawn = 0
						}
					case 10302:
						wantMotd ^= 1
					case 10304:
						wantCycle ^= 1
					case 10305:
						wantFlags ^= 2
						wantUpdated = 1
					case 10306:
						wantFlags ^= 0x2000
						wantUpdated = 1
					}
					r.Respawn = int(*memmap.PtrUint32(0x5D4594, 3584))
					r.Motd = legacy.Get_nox_server_sendMotd_108752()
					r.Cycle = legacy.Sub_4D0D70()
					r.Subflags = *globals["server-subflags"]
					r.Updated = *o.optionWords["settings-updated"]
					if r.Respawn != wantRespawn || r.Motd != wantMotd || r.Cycle != wantCycle || r.Subflags != wantFlags || r.Updated != wantUpdated {
						t.Fatalf("toggle got %+v want %d/%d/%d/%x/%d", r, wantRespawn, wantMotd, wantCycle, wantFlags, wantUpdated)
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "server-panels-general-toggles", rows, "2b730ea67b8bcf88e2bb70a07199ac0e27ba1aa91f6c41ae6afb2774096f8c30")
}
