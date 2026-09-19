//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestServerBrowserStateWords(t *testing.T) {
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	type row struct {
		Value        int32
		Set          int32
		Reads        [4]uint32
		Error, State uint32
	}
	var rows []row
	for _, v := range []int32{0, 1, 2, 3, 6, 9, 10, -1, -2147483648, 2147483647} {
		for _, name := range []string{"dword_5d4594_815052", "nox_game_createOrJoin_815048", "dword_5d4594_815104", "dword_5d4594_815044"} {
			*words[name] = uint32(v)
		}
		r := row{Value: v, Set: legacy.PortTestServerBrowserStateSet(v)}
		for i := range r.Reads {
			r.Reads[i] = legacy.PortTestServerBrowserStateGet(i)
			if r.Reads[i] != uint32(v) {
				t.Fatal("state word", i, v, r.Reads[i])
			}
		}
		if r.Set != v {
			t.Fatal("setter return", v, r.Set)
		}
		legacy.PortTestServerBrowserErrorSet(v)
		r.Error = *words["nox_client_connError_814552"]
		r.State = *words["dword_5d4594_814548"]
		if r.Error != uint32(v) || r.State != 2 {
			t.Fatal("error setter", r)
		}
		rows = append(rows, r)
	}
	spellbookCapture(t, "server-browser-state-words", rows, "ba62d2f4c7d3587d4fc0d9bb0f92103864f1884b7a24deadd8acab48526916cf")
}
func TestServerBrowserSelectedEndpoint(t *testing.T) {
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	port := serverConfigOwnBytes(t, 0x5D4594, 814604, 4)
	record, free := alloc.Calloc(1, 172)
	defer free()
	raw := unsafe.Slice((*byte)(record), 172)
	type row struct {
		Gate, Port uint32
		Address    string
		Mode       uint16
		Got        [3]uint32
	}
	var rows []row
	addresses := []string{"0.0.0.0", "127.0.0.1", "255.255.255.255", "192.168.1.9", "127.1", "1", "0x7f000001", "0127.0.0.1", "1.2.3.999", "localhost", "", "1.2.3.4 "}
	for _, gate := range []uint32{0, 1, 2, 0xffffffff} {
		for _, address := range addresses {
			for _, value := range []uint32{0, 18590, 65535, 0x12345678, 0xffffffff} {
				*words["dword_5d4594_815056"] = gate
				*words["dword_5d4594_814624"] = uint32(uintptr(record))
				clear(raw[12:28])
				copy(raw[12:28], address)
				binary.LittleEndian.PutUint16(raw[163:], 0xf431)
				binary.LittleEndian.PutUint32(port, value)
				r := row{Gate: gate, Port: value, Address: address, Mode: 0xf431}
				for i := range r.Got {
					r.Got[i] = legacy.PortTestServerBrowserStateGet(4 + i)
				}
				if gate == 0 {
					if r.Got != [3]uint32{} {
						t.Fatal("inactive selection", r)
					}
				} else {
					if r.Got[1] != value || r.Got[2] != 0xf431 {
						t.Fatal("selected words", r)
					}
					if address == "127.0.0.1" && r.Got[0] != 0x0100007f {
						t.Fatal("address byte order", r)
					}
				}
				rows = append(rows, r)
			}
		}
	}
	// Inactive selection must not dereference stale selection storage.
	*words["dword_5d4594_815056"] = 0
	*words["dword_5d4594_814624"] = 1
	for i := 4; i <= 6; i++ {
		if legacy.PortTestServerBrowserStateGet(i) != 0 {
			t.Fatal("inactive getter", i)
		}
	}
	spellbookCapture(t, "server-browser-selected-endpoint", rows, "cd479701a8db191cf882b8216bfe888c763bbd09542bead8de93955e2ea2ecff")
}
func TestServerBrowserSortButtons(t *testing.T) {
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	sort := words["nox_wol_servers_sorting_166704"]
	type row struct {
		Before uint32
		ID     int
		After  uint32
	}
	var rows []row
	for _, before := range []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0xffffffff} {
		for _, id := range []int{0, 10046, 10047, 10048, 10049, 10050, 10051, 10052} {
			*sort = before
			legacy.PortTestServerBrowserSortClick(id)
			want := before
			if id >= 10047 && id <= 10051 {
				want = uint32(2 * (id - 10047))
				if before == want {
					want++
				}
			}
			if *sort != want {
				t.Fatal("sort toggle", before, id, *sort, want)
			}
			rows = append(rows, row{before, id, *sort})
		}
	}
	spellbookCapture(t, "server-browser-sort-buttons", rows, "120969dcd7bc15721779a06438dbf30618ee225b6cb324070b16ffe5074a62a1")
}
func TestServerBrowserAddressFormatting(t *testing.T) {
	out, free := alloc.Make([]byte{}, 128)
	defer free()
	type row struct {
		Address string
		Port    uint16
		Return  int
		Text    string
		Guard   byte
	}
	var rows []row
	for _, address := range []string{"", "127.0.0.1", "host-name", "255.255.255.255", "name with space"} {
		for _, port := range []uint16{0, 1, 18590, 32768, 65535} {
			for i := range out {
				out[i] = 0xad
			}
			ret := legacy.PortTestServerBrowserFormat(address, port, unsafe.Pointer(&out[0]))
			text := alloc.GoStringS(out)
			want := fmt.Sprintf("%s:%d", address, port)
			if text != want || ret != len(want) || out[len(want)+1] != 0xad {
				t.Fatal("endpoint format", text, want, ret)
			}
			rows = append(rows, row{address, port, ret, text, out[len(want)+1]})
		}
	}
	spellbookCapture(t, "server-browser-address-formatting", rows, "83fe9df285956bee14e1d664665b665222aee8fb741094abd9db75a8e647a623")
}
