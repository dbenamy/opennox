//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionGauntletBookRemoval(t *testing.T) {
	q := newQuickbarOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldAbilities := q.c.srv.abilities.defs
	t.Cleanup(func() { q.c.srv.abilities.defs = oldAbilities })
	for id := 1; id <= 5; id++ {
		q.c.srv.abilities.defs[id] = AbilityDef{name: "Ability", field24: 1}
	}
	family := serverConfigOwnBytes(t, 0x587000, 132100, 32)
	copy(family, blobdata.PortTestBookGuideFamily())
	*memmap.PtrPtr(0x587000, 132124) = memmap.PtrOff(0x587000, 132100)
	guides := serverConfigOwnBytes(t, 0x5D4594, 740076, 41*28)
	clear(guides)
	for i := 1; i <= 40; i++ {
		binary.LittleEndian.PutUint32(guides[28*i+4:], 1)
	}
	type row struct {
		On, Kind, Rank int
		ID             uint16
		Return         int
		State          quickbarResult
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{17, 18, 19} {
			ids := []uint16{1, 3, 5}
			if kind == 19 {
				ids = []uint16{1, 24, 40}
			}
			for _, id := range ids {
				for _, rank := range []uint32{0, 1, 7} {
					setup := func() {
						q.reset(t)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						q.configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny) | 0x1000, Valid: true}, {Index: 2, Flags: uint32(things.SpellClassAny) | 0x2000, Valid: true}, {Index: 3, Flags: uint32(things.SpellClassAny), Valid: true}, {Index: 5, Flags: uint32(things.SpellClassAny), Valid: true}})
						p := &q.players[0]
						*(*byte)(unsafe.Add(p.C(), 2251)) = 2
						for _, r := range []struct{ off, n int }{{3696, 137}, {4244, 41}} {
							v := unsafe.Slice((*uint32)(unsafe.Add(p.C(), r.off)), r.n)
							for i := range v {
								v[i] = rank
							}
						}
						for i := 1; i <= 5; i++ {
							*memmap.PtrUint32(0x5D4594, 1047788+uintptr((i-1)*24+16)) = rank
						}
						for i := 0; i < 25; i++ {
							v := uint32(i%5 + 1)
							if kind == 19 {
								v = []uint32{1, 24, 7, 8, 25, 26, 36, 40}[i%8] + 74
							}
							q.bar[2*i], q.bar[2*i+1] = v, 0xa1000000+uint32(i)
						}
						if kind == 19 {
							*q.words["dword_5d4594_1046868"], *q.words["dword_5d4594_1046872"] = 1, 1
						}
					}
					op := map[int]string{17: "sub_45D320", 18: "nox_xxx_clientQuestDisableAbility_45D4A0", 19: "sub_45D400"}[kind]
					setup()
					if on != 0 {
						q.bookCall(op, uint32(id))
					}
					want := q.snapshot("quest removal", 0)
					setup()
					data := []byte{240, byte(kind), byte(id), byte(id >> 8)}
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
					got := q.snapshot("quest removal", 0)
					if n != 4 || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
						t.Fatal("quest knowledge removal", on, kind, id, rank, n)
					}
					rows = append(rows, row{on, kind, int(rank), id, n, got})
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-book-removal", rows)
}
