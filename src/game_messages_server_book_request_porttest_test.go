//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestGameMessageServerBookRequest(t *testing.T) {
	o := newReliableReportsOwner(t)
	u := &o.units[0]
	pl := u.UpdateDataPlayer().Player
	items, freeItems := alloc.Make([]server.Object{}, 3)
	t.Cleanup(freeItems)
	data, freeData := alloc.Make([]byte{}, 3*64)
	t.Cleanup(freeData)
	session, freeSession := alloc.New(server.TradeSession{})
	t.Cleanup(freeSession)
	stock, freeStock := alloc.New(server.TradeItem{})
	t.Cleanup(freeStock)
	stock.Object = &items[1]
	session.Stock = stock
	oldList, oldPending := o.s.Objs.List, o.s.Objs.Pending
	oldInventory, oldTrade := u.InvFirstItem, u.UpdateDataPlayer().Trade70
	t.Cleanup(func() {
		o.s.Objs.List, o.s.Objs.Pending = oldList, oldPending
		u.InvFirstItem, u.UpdateDataPlayer().Trade70 = oldInventory, oldTrade
	})
	guides := bookAwardWords(t, 0x587000, 70500, 41)
	for i := range guides {
		name, free := alloc.CString(fmt.Sprintf("server-book-guide-%d", i))
		t.Cleanup(free)
		guides[i] = uint32(uintptr(unsafe.Pointer(name)))
	}
	type row struct {
		Sources, Subtype int
		Code             uint16
		Immobile         bool
		State            legacy.PortTestReliableReportState
	}
	var rows []row
	for mask := 0; mask < 8; mask++ {
		for _, kind := range []byte{0, 1, 2, 255} {
			for _, code := range []uint16{0, 7, 0x802a, 0xffff, 0x8000} {
				for _, immobile := range []bool{false, true} {
					func() {
						restoreLookup := legacy.PortTestCreatureXferLookupOwner()
						defer restoreLookup()
						o.reset()
						clear(data)
						for i := range items {
							items[i] = server.Object{NetCode: 7, Extent: uint32(40 + i), ObjClass: object.ClassSimple}
							if immobile {
								items[i].ObjClass = object.ClassImmobile
							}
							items[i].UseData.SetPtr(unsafe.Pointer(&data[64*i]))
							if kind == 2 {
								name := fmt.Sprintf("server-book-guide-%d", i+1)
								if i == 2 {
									name = "unknown-guide"
								}
								copy(data[64*i:], name)
							} else {
								data[64*i] = []byte{0, 17, 255}[i]
							}
						}
						u.InvFirstItem, u.UpdateDataPlayer().Trade70 = nil, nil
						o.s.Objs.List, o.s.Objs.Pending = nil, nil
						if mask&1 != 0 {
							u.InvFirstItem = &items[0]
						}
						if mask&2 != 0 {
							u.UpdateDataPlayer().Trade70 = session
						}
						if mask&4 != 0 {
							o.s.Objs.List = &items[2]
						}
						before := bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(&items[0])), len(items)*int(unsafe.Sizeof(items[0]))))
						useBefore := bytes.Clone(data)
						msg := []byte{226, byte(code), byte(code >> 8), kind}
						input := bytes.Clone(msg)
						n := legacy.Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(pl.PlayerIndex(), msg, pl, u, u.UpdateData)
						if n != 4 || !bytes.Equal(msg, input) || !bytes.Equal(data, useBefore) || !bytes.Equal(before, unsafe.Slice((*byte)(unsafe.Pointer(&items[0])), len(before))) {
							t.Fatal("book request changed input/object storage or consumed length")
						}
						selected := -1
						if code == 7 || code == 0x802a && mask&4 != 0 {
							for i := 0; i < 3; i++ {
								if mask&(1<<i) != 0 {
									selected = i
									break
								}
							}
						}
						state := o.state()
						if selected < 0 {
							if len(state.Nodes) != 0 {
								t.Fatal("missing book produced a response")
							}
						} else {
							value := []byte{0, 17, 255}[selected]
							if kind == 2 {
								value = []byte{1, 2, 0}[selected]
							}
							wireCode := uint16(7)
							if immobile {
								wireCode = 0x8000 | uint16(40+selected)
							}
							want := []byte{226, 0, 0, value}
							binary.LittleEndian.PutUint16(want[1:], wireCode)
							if len(state.Nodes) != 1 || !bytes.Equal(state.Nodes[0].Data, want) || state.Nodes[0].To != byte(pl.PlayerIndex()) || state.Nodes[0].Priority != 1 || state.Nodes[0].Ordered != 0 || state.Nodes[0].Related != 0 {
								t.Fatalf("book response sources%d kind%d code%x immobile%v: %+v", mask, kind, code, immobile, state.Nodes)
							}
						}
						rows = append(rows, row{mask, int(kind), code, immobile, state})
					}()
				}
			}
		}
	}
	interactionCapture(t, "game-server-book-request", rows)
}
