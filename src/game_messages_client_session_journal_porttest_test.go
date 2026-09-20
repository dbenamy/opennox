//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"strings"
	"testing"
)

func TestGameMessageClientSessionJournal(t *testing.T) {
	o := newJournalOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Present, Kind, Seed, Return int
		Name                            string
		Flags                           uint16
		State                           journalResult
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for _, kind := range []int{0, 1, 2, 3, 4, 255} {
				for seed := 0; seed < 2; seed++ {
					for _, name := range []string{"Short", "Wrapped", "Unicode", strings.Repeat("X", 63)} {
						for _, flags := range []uint16{0, 2, 8, 0x8000, 0xffff} {
							setup := func() {
								o.resetJournal(t)
								binary.LittleEndian.PutUint32(connected, uint32(on))
								o.call(t, 0, 1, "Short", 2)
								o.call(t, 0, 7, "Wrapped", 8)
								if seed != 0 {
									o.call(t, 0, 31, "Short", 2)
									o.call(t, 0, 31, "Wrapped", 8)
									o.call(t, 0, 31, "Short", 4)
								}
								if present == 0 {
									legacy.Set_dword_8531A0_2576(nil)
								}
							}
							setup()
							if on != 0 {
								if present != 0 {
									switch kind {
									case 1:
										o.call(t, 0, 31, name, flags)
									case 2:
										o.call(t, 2, 31, name, flags)
									case 3:
										o.call(t, 5, 31, name, flags)
									}
								}
								if kind == 1 || kind == 2 {
									legacy.PortTestJournal(9, 0, 0, 0)
								}
							}
							want := o.snapshot(t, 0, 0, false)
							setup()
							data := make([]byte, 68)
							data[0], data[1] = 213, byte(kind)
							copy(data[2:66], name)
							binary.LittleEndian.PutUint16(data[66:], flags)
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(213), data)
							if on != 0 && present != 0 && kind == 1 {
								o.remember(o.players[31].Journal)
							}
							got := o.snapshot(t, 0, 0, false)
							wantN := -1
							if kind >= 1 && kind <= 3 {
								wantN = 68
							}
							if n != wantN || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
								t.Fatal("journal message", on, present, kind, seed, name, flags, n)
							}
							rows = append(rows, row{on, present, kind, seed, n, name, flags, got})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-journal", rows)
}
