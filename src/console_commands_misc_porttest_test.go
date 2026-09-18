//go:build porttest

package opennox

import (
	"encoding/binary"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

func TestConsoleCommandsInformation(t *testing.T) {
	o := newConsoleCommandOwner(t)
	var sounds [][2]int
	t.Cleanup(legacy.PortTestClientSoundObserver(func(id, volume int) { sounds = append(sounds, [2]int{id, volume}) }))
	before := legacy.Get_dword_5d4594_811904()
	t.Cleanup(func() {
		if legacy.Get_dword_5d4594_811904() != before {
			o.call(t, "show info", false)
		}
	})
	for i := 0; i < 2; i++ {
		if !o.call(t, "show info", false) {
			t.Fatal("info return")
		}
		want := before
		if i == 0 {
			want = 1 - before
		}
		if legacy.Get_dword_5d4594_811904() != want {
			t.Fatal("info toggle")
		}
	}
	if !reflect.DeepEqual(sounds, [][2]int{{921, 100}, {921, 100}}) {
		t.Fatal(sounds)
	}
	oldRank := legacy.Sub_4703F0
	t.Cleanup(func() { legacy.Sub_4703F0 = oldRank })
	calls := 0
	legacy.Sub_4703F0 = func() { calls++ }
	for _, flags := range []uint32{0, 8192} {
		restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
		if !o.call(t, "show rank", false) {
			t.Fatal("rank return")
		}
		restore()
	}
	if calls != 1 {
		t.Fatal("rank dispatch", calls)
	}
	// The memory command calls the real accounting routine; seed a live entry and
	// verify its table output, while keeping the global record list fixture-owned.
	table := serverConfigOwnBytes(t, 0x5D4594, 252284, 86096)
	clear(table)
	record, free := alloc.Make([]byte{}, 64)
	t.Cleanup(free)
	binary.LittleEndian.PutUint32(record[16:], 37)
	copy(record[20:], "fixture")
	head := serverConfigOwnBytes(t, 0x5D4594, 338304, 4)
	binary.LittleEndian.PutUint32(head, uint32(uintptr(unsafe.Pointer(&record[0]))))
	if !o.call(t, "show mem", false) || binary.LittleEndian.Uint32(table[80:]) != 37 || string(table[:7]) != "fixture" {
		t.Fatal("memory command")
	}
	seq := serverConfigOwnBytes(t, 0x5D4594, 1197340, 12)
	addr := uint32(uintptr(unsafe.Pointer(&seq[0])))
	binary.LittleEndian.PutUint32(seq, addr)
	binary.LittleEndian.PutUint32(seq[4:], addr)
	binary.LittleEndian.PutUint32(seq[8:], addr)
	if !o.call(t, "show seq", false) {
		t.Fatal("sequence command")
	}
}
func TestConsoleCommandsAbilityRoute(t *testing.T) {
	o := newConsoleCommandOwner(t)
	units, configure, _, free := o.c.srv.S().PortTestEscortPlayers()
	t.Cleanup(free)
	configure(3)
	var got []*server.Object
	old := legacy.Nox_xxx_playerCancelAbils_4FC180
	legacy.Nox_xxx_playerCancelAbils_4FC180 = func(u *server.Object) { got = append(got, u) }
	t.Cleanup(func() { legacy.Nox_xxx_playerCancelAbils_4FC180 = old })
	for _, flags := range []uint32{0, 8192} {
		restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
		if !o.call(t, "cheat ability", false) {
			t.Fatal("ability return")
		}
		restore()
	}
	if !reflect.DeepEqual(got, []*server.Object{&units[0], &units[1], &units[2]}) {
		t.Fatal("ability service route")
	}
}
