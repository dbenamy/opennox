//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestObjectReportsSelfUpdate(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	u := &units[0]
	ud := u.UpdateDataPlayer()
	pl := ud.Player
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	var rows []struct {
		Name    string
		Return  int
		Packet  []byte
		Reports []legacy.PortTestShopPacketResult
		Words   [3]uint32
	}
	defer func() {
		spellbookCapture(t, "object-reports-self", rows, "28aab84b6198587e0e333f5ce1e6940eb4a30a7b457037b0a4ec6b6fef76fcec")
	}()
	for _, route := range []int{0, 1} {
		for _, changed := range []bool{false, true} {
			for _, full := range []bool{false, true} {
				if route == 1 && full {
					continue
				}
				name := fmt.Sprintf("route%d/changed%t/full%t", route, changed, full)
				t.Run(name, func(t *testing.T) {
					s.NetList.ResetAll()
					reset()
					*ud = server.PlayerUpdateData{Player: pl}
					s.SetFrame(123)
					*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 2164)) = 0x12345678
					*(*byte)(unsafe.Add(u.CObj(), 440)) = 19
					objectXferSetWord(u.UpdateData, 228, 0x3f000000)
					objectXferSetWord(u.UpdateData, 232, 0x3f000000)
					*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 2168)) = 0x12345678
					*(*byte)(unsafe.Add(unsafe.Pointer(pl), 2172)) = 19
					if changed {
						objectXferSetWord(u.UpdateData, 232, 0)
						*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 2168)) = 0
						*(*byte)(unsafe.Add(unsafe.Pointer(pl), 2172)) = 0
					}
					if full && !s.NetList.AddToMsgListCli(1, netlist.Kind1, bytes.Repeat([]byte{0xa5}, 2048)) {
						t.Fatal("fill")
					}
					rv := legacy.PortTestObjectReports(5, u, u, 1, route, 0, nil)
					kind := netlist.Kind1
					if route != 0 {
						kind = netlist.Kind2
					}
					got := s.NetList.CopyPacketsA(1, kind)
					rr := snapshot()
					want := 1
					size := 12
					if full {
						want = 0
						size = 2048
					}
					if rv != want || len(got) != size {
						t.Fatalf("return%d size%d", rv, len(got))
					}
					n := 0
					if changed {
						n = 3
					}
					if len(rr) != n {
						t.Fatalf("reports%d want%d", len(rr), n)
					}
					if changed {
						for i, op := range []byte{91, 74, 73} {
							if rr[i].Data[0] != op || rr[i].Recipient != 1 {
								t.Fatal("self report order/recipient")
							}
						}
					}
					state := [3]uint32{objectXferGetWord(u.UpdateData, 232), *(*uint32)(unsafe.Add(unsafe.Pointer(pl), 2168)), uint32(*(*byte)(unsafe.Add(unsafe.Pointer(pl), 2172)))}
					if state != [3]uint32{0x3f000000, 0x12345678, 19} {
						t.Fatal("self history not updated")
					}
					rows = append(rows, struct {
						Name    string
						Return  int
						Packet  []byte
						Reports []legacy.PortTestShopPacketResult
						Words   [3]uint32
					}{name, rv, got, rr, state})
				})
			}
		}
	}
}
