//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestPlayerFilesServerLoadFraming(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestPlayerFileAbilities())
	oldAbilities := noxServer.abilities
	noxServer.abilities.Init(noxServer)
	t.Cleanup(func() { noxServer.abilities = oldAbilities })
	t.Cleanup(o.s.PortTestInventoryDisplayBalance())
	oldStats := o.s.Players.Stats
	t.Cleanup(func() { o.s.Players.Stats = oldStats })
	o.s.Players.Stats.Base = server.ClassStats{Health: 25, Mana: 15, Speed: 1500, Strength: 10}
	o.s.Players.Stats.Warrior = server.ClassStats{Health: 150, Mana: 60, Speed: 4000, Strength: 40}
	o.s.Players.Stats.Wizard = server.ClassStats{Health: 80, Mana: 150, Speed: 3500, Strength: 20}
	o.s.Players.Stats.Conjurer = server.ClassStats{Health: 110, Mana: 110, Speed: 3750, Strength: 30}
	carry := memmap.PtrFloat64(0x581450, 10216)
	oldCarry := *carry
	*carry = 1
	t.Cleanup(func() { *carry = oldCarry })
	saved := unsafe.Slice(memmap.PtrUint32(0x5D4594, 527696), 2)
	oldSaved := append([]uint32(nil), saved...)
	t.Cleanup(func() { copy(saved, oldSaved) })
	t.Cleanup(o.s.PortTestBookSpellOwner())
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	configure([]server.PortTestSpellClassDef{{Index: 1, Flags: 0x100, Valid: true}, {Index: 2, Flags: 0x100, Valid: true}})
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) })
	flags.ResetGame()
	u := &o.units[0]
	ud := u.UpdateDataPlayer()
	p := ud.Player
	p.NetCodeVal = u.NetCode
	hp, free := alloc.New(server.HealthData{})
	t.Cleanup(free)
	u.HealthData = hp
	*(*byte)(unsafe.Add(unsafe.Pointer(p), 2251)) = 1
	var rows []map[string]any
	for _, which := range []string{"empty", "unknown", "bad-attributes", "bad-status", "missing", "no-player", "no-unit"} {
		o.reset()
		flags.ResetGame()
		*hp = server.HealthData{Cur: 17, Max: 25}
		ud.ManaCur = 11
		ud.ManaPrev = 7
		ud.ManaMax = 15
		u.Experience = 100
		u.Buffs = 0
		p.PlayerUnit = u
		path := filepath.Join(t.TempDir(), "server.plr")
		if which != "missing" {
			f, err := cryptfile.OpenFile(path, cryptfile.WriteOnly, 27)
			if err != nil {
				t.Fatal(err)
			}
			section := func(id uint32, data []byte) {
				t.Helper()
				if err := f.WriteU32(id); err != nil {
					t.Fatal(err)
				}
				f.SectionStart()
				if _, err := f.Write(data); err != nil {
					t.Fatal(err)
				}
				f.SectionEnd()
			}
			if which == "unknown" {
				section(99, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9})
			}
			if which == "bad-attributes" {
				section(2, binary.LittleEndian.AppendUint16(nil, 6))
			}
			if which == "bad-status" {
				section(3, binary.LittleEndian.AppendUint16(nil, 3))
			}
			if err = f.WriteU32(0); err != nil {
				t.Fatal(err)
			}
			if err = f.Close(); err != nil {
				t.Fatal(err)
			}
		}
		index := uint32(p.PlayerInd)
		if which == "no-player" {
			index = 2
		}
		if which == "no-unit" {
			p.PlayerUnit = nil
		}
		var before []byte
		if which != "missing" {
			var err error
			before, err = os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
		}
		name, release := alloc.CString(path)
		ret := legacy.PortTestPlayerFileCall("nox_xxx_cliPlrInfoLoadFromFile_41A2E0", uint32(uintptr(unsafe.Pointer(name))), index)
		release()
		success := which == "empty" || which == "unknown"
		wantRet := uint32(0)
		if success {
			wantRet = 1
		}
		if ret != wantRet || cryptfile.Global() != nil || o.s.Players.CheckXxx(u) {
			t.Fatal("server load return/cleanup", which, ret, wantRet)
		}
		if success && (hp.Cur != 17 || ud.ManaCur != 11 || hp.Max != 25 || ud.ManaMax != 15 || ud.Field19_0 != 17) {
			t.Fatal("server load resource restoration", which, *hp, ud.ManaCur, ud.ManaMax, ud.Field19_0)
		}
		if which != "missing" {
			after, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatal("server load changed source", which, err)
			}
		}
		rows = append(rows, map[string]any{"case": which, "return": ret, "health": *hp, "mana": ud.ManaCur, "max_mana": ud.ManaMax, "queue": o.state()})
		p.PlayerUnit = u
	}
	spellbookCapture(t, "player-files-server-load-framing", rows, "26372fb3cc108c066e1cfaa4838157828ff623c62447f3db03b1a7a6a2732de7")
}
