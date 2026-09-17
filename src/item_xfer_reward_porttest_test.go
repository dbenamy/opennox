//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

var itemXferRewardGroups = []struct {
	offset, count int
	name          func(index int) string
}{
	{8, 137, func(i int) string { return spell.ID(i).String() }},
	{145, 6, func(i int) string { return server.Ability(i).String() }},
	{151, 41, func(i int) string { return fmt.Sprintf("PortGuide%02d", i) }},
}

func itemXferRewardOwner(t *testing.T) *server.Server {
	s := newItemXferOwner(t)
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 70500), 41)
	old := append([]uint32(nil), table...)
	for i := range table {
		p, free := alloc.CString(itemXferRewardGroups[2].name(i))
		t.Cleanup(free)
		table[i] = uint32(uintptr(unsafe.Pointer(p)))
	}
	t.Cleanup(func() { copy(table, old) })
	return s
}
func itemXferRewardName(p *mapDrawableStream, name string) {
	p.u8(byte(len(name)))
	p.WriteString(name)
}
func itemXferRewardTail(p *mapDrawableStream, version uint16) {
	for _, off := range []int{196, 192, 200, 204, 208} {
		p.u32(uint32(1000 + off))
	}
	if version >= 62 {
		p.u32(1212)
	}
	if version >= 63 {
		p.u8(17)
	}
}
func TestItemXferRewardMasks(t *testing.T) {
	s := itemXferRewardOwner(t)
	path := filepath.Join(t.TempDir(), "reward.bin")
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "item-xfer-reward-masks", rows, "db3294dda28c6f942b267731f0a320ccec573e3f3af4aadc93e6c893497b2ec4")
	}()
	run := func(label string, group, index int, value byte) {
		t.Run(label, func(t *testing.T) {
			u := newItemXferObject(t, s, "RewardMarker")
			objectXferSetCommon(u)
			u.Field34 = 777
			data := unsafe.Slice((*byte)(u.InitData), 220)
			objectXferSetWord(u.InitData, 0, 123)
			objectXferSetWord(u.InitData, 4, 456)
			for _, off := range []int{192, 196, 200, 204, 208, 212} {
				objectXferSetWord(u.InitData, off, uint32(1000+off))
			}
			data[216] = 17
			for g, spec := range itemXferRewardGroups {
				for i := 0; i < spec.count; i++ {
					if group < 0 && i > 0 || group == g && index == i {
						data[spec.offset+i] = value
					}
				}
			}
			want := objectXferCurrentRecord(63)
			want.u32(123)
			want.u32(456)
			for _, spec := range itemXferRewardGroups {
				count := uint16(0)
				for i := 0; i < spec.count; i++ {
					if data[spec.offset+i] == 1 {
						count++
					}
				}
				want.u16(count)
				for i := 0; i < spec.count; i++ {
					if data[spec.offset+i] != 0 {
						itemXferRewardName(want, spec.name(i))
					}
				}
			}
			itemXferRewardTail(want, 63)
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			if err := u.CallXfer(nil); err != nil {
				t.Fatal(err)
			}
			written := cryptfile.Global().PortTestChecksum()
			cryptfile.Close()
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want.Bytes()) {
				t.Fatal("reward wire format")
			}
			if u.Field34 != 777 {
				t.Fatal("writer lifetime")
			}
			if value != 1 || index == 0 {
				itemXferCaptureCase(t, &rows, u, 0, written)
				return
			}
			v := newItemXferObject(t, s, "RewardMarker")
			v.Field34 = 999
			checksum := itemXferReadRecord(t, v, path)
			if !bytes.Equal(unsafe.Slice((*byte)(v.InitData), 220), data) || v.Field34 != 999 {
				t.Fatal("reward state/lifetime")
			}
			itemXferCaptureCase(t, &rows, v, checksum, written)
		})
	}
	for g, spec := range itemXferRewardGroups {
		for i := 1; i < spec.count; i++ {
			run(fmt.Sprintf("group%d-id%d", g, i), g, i, 1)
		}
		for _, value := range []byte{1, 2, 255} {
			run(fmt.Sprintf("group%d-zero-id-value%d", g, value), g, 0, value)
			if value != 1 {
				run(fmt.Sprintf("group%d-noncanonical%d", g, value), g, 1, value)
			}
		}
	}
	run("all", -1, -1, 1)
}
func TestItemXferRewardReadFailures(t *testing.T) {
	s := itemXferRewardOwner(t)
	path := filepath.Join(t.TempDir(), "reward-failure.bin")
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "item-xfer-reward-failures", rows, "ec77c9e1785dcf9d5d4ceb9b20c86fa3dcea89f3bb0907fd9deb27e1859494f2")
	}()
	for group, spec := range itemXferRewardGroups {
		for _, first := range []bool{false, true} {
			for _, bad := range []string{"", "unknown", strings.Repeat("x", 127), strings.Repeat("x", 255), spec.name(0)} {
				t.Run(fmt.Sprintf("group%d-first%v-bad%q", group, first, bad), func(t *testing.T) {
					u := newItemXferObject(t, s, "RewardMarker")
					u.Field34 = 999
					data := unsafe.Slice((*byte)(u.InitData), 220)
					for _, g := range itemXferRewardGroups {
						data[g.offset+g.count-1] = 7
					}
					wantData := append([]byte(nil), data...)
					p := objectXferCurrentRecord(63)
					p.u32(123)
					p.u32(456)
					for g := 0; g <= group; g++ {
						n := uint16(1)
						if g == group && first {
							n = 2
						}
						p.u16(n)
						if g < group || first {
							itemXferRewardName(p, itemXferRewardGroups[g].name(1))
							wantData[itemXferRewardGroups[g].offset+1] = 1
						}
						if g == group {
							itemXferRewardName(p, bad)
						}
					}
					stop := p.Len()
					p.u32(0xdeadbeef)
					if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
						t.Fatal(err)
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
						t.Fatal(err)
					}
					defer cryptfile.Close()
					if err := u.CallXfer(nil); err == nil {
						t.Fatal("accepted invalid reward name")
					}
					pos, err := cryptfile.Global().File.Seek(0, 1)
					if err != nil || pos != int64(stop) {
						t.Fatalf("position=%d want=%d err=%v", pos, stop, err)
					}
					objectXferSetWord(unsafe.Pointer(&wantData[0]), 0, 123)
					objectXferSetWord(unsafe.Pointer(&wantData[0]), 4, 456)
					if !bytes.Equal(data, wantData) || u.Field34 != 0 {
						t.Fatal("partial reward mutation/lifetime")
					}
					itemXferCaptureCase(t, &rows, u)
				})
			}
		}
	}
}
func TestItemXferRewardReadMerge(t *testing.T) {
	s := itemXferRewardOwner(t)
	path := filepath.Join(t.TempDir(), "reward-merge.bin")
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "item-xfer-reward-merge", rows, "181f8397e0143fdfe072e8a5f89cd3b83a9dbc7c99011bfcd6678e0e862f6c19")
	}()
	for _, version := range []uint16{61, 62, 63} {
		t.Run(fmt.Sprintf("v%d", version), func(t *testing.T) {
			u := newItemXferObject(t, s, "RewardMarker")
			u.Field34 = 999
			data := unsafe.Slice((*byte)(u.InitData), 220)
			for _, g := range itemXferRewardGroups {
				data[g.offset+g.count-1] = 7
			}
			objectXferSetWord(u.InitData, 212, 999)
			data[216] = 99
			p := objectXferCurrentRecord(version)
			p.u32(123)
			p.u32(456)
			for _, g := range itemXferRewardGroups {
				p.u16(2)
				itemXferRewardName(p, g.name(1))
				itemXferRewardName(p, g.name(1))
			}
			itemXferRewardTail(p, version)
			if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
				t.Fatal(err)
			}
			checksum := itemXferReadRecord(t, u, path)
			for _, g := range itemXferRewardGroups {
				for i := 0; i < g.count; i++ {
					expected := byte(0)
					if i == 1 {
						expected = 1
					}
					if i == g.count-1 {
						expected = 7
					}
					if data[g.offset+i] != expected {
						t.Fatal("reward mask merge")
					}
				}
			}
			tail := uint32(999)
			if version >= 62 {
				tail = 1212
			}
			last := byte(99)
			if version >= 63 {
				last = 17
			}
			if objectXferGetWord(u.InitData, 212) != tail || data[216] != last || u.Field34 != 999 {
				t.Fatal("versioned reward tail/lifetime")
			}
			itemXferCaptureCase(t, &rows, u, checksum)
		})
	}
}
