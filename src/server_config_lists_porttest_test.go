//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"testing"
	"unsafe"
)

func serverConfigListOwner(t *testing.T) *serverOptionsOwner {
	t.Helper()
	o := newServerOptionsOwner(t)
	for _, off := range []uintptr{371364, 371500} {
		h := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, off)), 3)
		old := append([]uint32(nil), h...)
		t.Cleanup(func() { copy(h, old) })
		addr := uint32(uintptr(unsafe.Pointer(&h[0])))
		h[0], h[1], h[2] = addr, addr, addr
	}
	t.Cleanup(func() {
		for _, kind := range []string{"allowed", "blocked"} {
			for i := 0; i < 100; i++ {
				if legacy.PortTestServerConfigAdmission(kind+"-first", 0, nil, nil) == nil {
					break
				}
				legacy.PortTestServerConfigAdmission(kind+"-remove", 0, nil, nil)
			}
			if legacy.PortTestServerConfigAdmission(kind+"-first", 0, nil, nil) != nil {
				t.Error("list cleanup did not terminate", kind)
			}
		}
	})
	return o
}
func serverConfigListAdd(kind, name, address string, duration int) unsafe.Pointer {
	p, free := alloc.CString16(name)
	defer free()
	var q unsafe.Pointer
	if address != "" {
		a, free := alloc.CString(address)
		defer free()
		q = unsafe.Pointer(a)
	}
	return legacy.PortTestServerConfigAdmission(kind+"-add", duration, unsafe.Pointer(p), q)
}

type serverConfigListRow struct {
	Name, Address string
	Expires       uint64
}

func serverConfigListSnapshot(t *testing.T, kind string) (out []serverConfigListRow, nodes []unsafe.Pointer) {
	t.Helper()
	p := legacy.PortTestServerConfigAdmission(kind+"-first", 0, nil, nil)
	for p != nil {
		if len(nodes) >= 100 {
			t.Fatal("list cycle")
		}
		nodes = append(nodes, p)
		r := serverConfigListRow{Name: alloc.GoString16((*uint16)(unsafe.Add(p, 12)))}
		if kind == "blocked" {
			r.Expires = binary.LittleEndian.Uint64(unsafe.Slice((*byte)(unsafe.Add(p, 64)), 8))
			r.Address = alloc.GoString((*byte)(unsafe.Add(p, 72)))
		}
		out = append(out, r)
		p = legacy.PortTestServerConfigAdmission(kind+"-next", 0, p, nil)
	}
	return
}
func TestServerConfigAdmissionLists(t *testing.T) {
	type row struct {
		Kind          string
		Index         int
		Before, After []serverConfigListRow
		Return        int
	}
	var rows []row
	for _, kind := range []string{"allowed", "blocked"} {
		for _, index := range []int{-2147483648, -1, 0, 1, 2, 3, 2147483647} {
			t.Run(fmt.Sprintf("%s/%d", kind, index), func(t *testing.T) {
				serverConfigListOwner(t)
				if got := legacy.PortTestServerConfigAdmission(kind+"-remove", index, nil, nil); got != nil {
					t.Fatal("empty removal result")
				}
				for _, name := range []string{"Alice", "Bob Smith", "Carol"} {
					if serverConfigListAdd(kind, name, "127.0.0.1", 0) != nil {
						t.Fatal("closed-panel append return")
					}
				}
				before, nodes := serverConfigListSnapshot(t, kind)
				if len(before) != 3 || before[0].Name != "Alice" || before[1].Name != "Bob Smith" || before[2].Name != "Carol" {
					t.Fatal("append order", before)
				}
				result := legacy.PortTestServerConfigAdmission(kind+"-remove", index, nil, nil)
				category := -1
				for i, p := range nodes {
					if result == p {
						category = i
					}
				}
				if result != nil && category < 0 {
					t.Fatal("unknown removal pointer")
				}
				after, _ := serverConfigListSnapshot(t, kind)
				want := append([]serverConfigListRow(nil), before...)
				valid := index >= 0 && index < len(want)
				if valid {
					want = append(want[:index], want[index+1:]...)
				}
				if !reflect.DeepEqual(after, want) {
					t.Fatal("remove index", index, after, want)
				}
				expected := -1
				if kind == "allowed" && valid {
					expected = index
				}
				if category != expected {
					t.Fatal("removal return", category, expected)
				}
				rows = append(rows, row{kind, index, before, after, category})
			})
		}
	}
	spellbookCapture(t, "server-config-admission-lists", rows, "a9cdf7c1b2bb4b77319c8c7fdfb68f9c9eaf6e037a1196e0f5e028b6a2e96b45")
}
func TestServerConfigAdmissionDurations(t *testing.T) {
	serverConfigListOwner(t)
	old := legacy.PlatformTicks
	t.Cleanup(func() { legacy.PlatformTicks = old })
	type row struct {
		Ticks    uint64
		Duration int
		Expires  uint64
		Calls    int
	}
	var rows []row
	for _, tick := range []uint64{0, 1, 0xffffffff, 0x100000000, ^uint64(0)} {
		for _, duration := range []int{-2147483648, -1, 0, 1, 71582, 71583, 2147483647} {
			calls := 0
			legacy.PlatformTicks = func() uint64 { calls++; return tick }
			serverConfigListAdd("blocked", "Timed", "", duration)
			snapshot, _ := serverConfigListSnapshot(t, "blocked")
			want := uint64(0)
			wantCalls := 0
			if duration != 0 {
				want = uint64(uint32(duration)*60000 + uint32(tick))
				wantCalls = 1
			}
			if len(snapshot) != 1 || snapshot[0].Expires != want || calls != wantCalls {
				t.Fatal("duration width", tick, duration, snapshot, want, calls)
			}
			rows = append(rows, row{tick, duration, snapshot[0].Expires, calls})
			legacy.PortTestServerConfigAdmission("blocked-remove", 0, nil, nil)
		}
	}
	spellbookCapture(t, "server-config-admission-durations", rows, "3640ee6f5ba0000fb593bd1cdc1556eda44cc5f89dfeb6dfe725352001667840")
}
func TestServerConfigAdmissionExpiry(t *testing.T) {
	serverConfigListOwner(t)
	old := legacy.PlatformTicks
	t.Cleanup(func() { legacy.PlatformTicks = old })
	ticks := uint64(100)
	legacy.PlatformTicks = func() uint64 { return ticks }
	serverConfigListAdd("blocked", "First", "", 1)
	serverConfigListAdd("blocked", "Second", "", 1)
	ticks = 60100
	legacy.PortTestServerConfigAdmission("expire", 0, nil, nil)
	equal, _ := serverConfigListSnapshot(t, "blocked")
	if len(equal) != 2 {
		t.Fatal("expiry must be strictly later than deadline", equal)
	}
	ticks++
	legacy.PortTestServerConfigAdmission("expire", 0, nil, nil)
	after, _ := serverConfigListSnapshot(t, "blocked")
	if len(after) != 0 {
		t.Fatal("consecutive expired entries must both be removed", after)
	}
	spellbookCapture(t, "server-config-admission-expiry", []any{equal, after}, "308d61ca0b108a4d03e68ac44a4a94aa69615cf5b3c5633a0331ace9de7603d0")
}

func TestServerConfigAdmissionExpiryMixed(t *testing.T) {
	type row struct {
		Mask           int
		Equal, Expired []serverConfigListRow
	}
	var rows []row
	for mask := 0; mask < 32; mask++ {
		t.Run(fmt.Sprint(mask), func(t *testing.T) {
			serverConfigListOwner(t)
			old := legacy.PlatformTicks
			t.Cleanup(func() { legacy.PlatformTicks = old })
			tick := uint64(1)
			legacy.PlatformTicks = func() uint64 { return tick }
			var want []serverConfigListRow
			for i := 0; i < 5; i++ {
				duration := 0
				if mask&(1<<uint(i)) != 0 {
					duration = 1
				}
				name := fmt.Sprintf("Name%d", i)
				serverConfigListAdd("blocked", name, "", duration)
				if duration == 0 {
					want = append(want, serverConfigListRow{Name: name})
				}
			}
			tick = 60001
			legacy.PortTestServerConfigAdmission("expire", 0, nil, nil)
			equal, _ := serverConfigListSnapshot(t, "blocked")
			if len(equal) != 5 {
				t.Fatal("deadline equality", equal)
			}
			tick++
			legacy.PortTestServerConfigAdmission("expire", 0, nil, nil)
			after, _ := serverConfigListSnapshot(t, "blocked")
			if !reflect.DeepEqual(after, want) {
				t.Fatal("mixed expiry", mask, after, want)
			}
			legacy.PortTestServerConfigAdmission("expire", 0, nil, nil)
			again, _ := serverConfigListSnapshot(t, "blocked")
			if !reflect.DeepEqual(again, want) {
				t.Fatal("second pass", again, want)
			}
			rows = append(rows, row{mask, equal, after})
		})
	}
	spellbookCapture(t, "server-config-admission-expiry-mixed", rows, "157354b611f9a029054ae9765f5f3c40c2717a7d3f70d3647dcc6fb8b1331458")
}
