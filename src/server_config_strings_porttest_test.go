//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestServerConfigNameStorage(t *testing.T) {
	o := newServerOptionsOwner(t)
	storage := serverConfigOwnBytes(t, 0x5D4594, 1324, 20)
	type row struct {
		Input   string
		Nil     bool
		Result  int
		Bytes   []byte
		Updated uint32
	}
	var rows []row
	for _, tc := range []struct {
		input, before, want string
		nilInput, changed   bool
	}{
		{"Alpha", "", "Alpha", false, true}, {"alpha", "Alpha", "Alpha", false, false},
		{"Beta", "Alpha", "Beta", false, true}, {"", "Alpha", "", false, true},
		{"", "", "", false, false}, {"", "Alpha", "", true, true}, {"", "", "", true, true},
		{"1234567890123456", "", "123456789012345", false, true},
		{"123456789012345XYZ", "123456789012345", "123456789012345", false, false},
		{"name with space", "", "name with space", false, true},
	} {
		for i := range storage {
			storage[i] = 0xa5
		}
		clear(storage[:16])
		copy(storage, tc.before)
		before := append([]byte(nil), storage...)
		*o.optionWords["settings-updated"] = 0
		var p unsafe.Pointer
		if !tc.nilInput {
			name, free := alloc.CString(tc.input)
			defer free()
			p = unsafe.Pointer(name)
		}
		result := legacy.PortTestServerConfigPointer("name-set", 0, p)
		category := 0
		if result != nil {
			if result != unsafe.Pointer(&storage[0]) {
				t.Fatal("name return is not owned storage")
			}
			category = 1
		}
		expected := append([]byte(nil), before...)
		if tc.nilInput {
			expected[0] = 0
		} else if tc.changed {
			clear(expected[:16])
			copy(expected[:15], tc.want)
		}
		if !bytes.Equal(storage, expected) || (*o.optionWords["settings-updated"] != 0) != tc.changed {
			t.Fatal("name mutation", tc, storage, expected)
		}
		if category != bool2int(tc.changed && !tc.nilInput) {
			t.Fatal("name pointer/scalar return", tc, category)
		}
		if legacy.PortTestServerConfigPointer("name-get", 0, nil) != unsafe.Pointer(&storage[0]) || alloc.GoString(&storage[0]) != tc.want {
			t.Fatal("name getter")
		}
		rows = append(rows, row{tc.input, tc.nilInput, category, append([]byte(nil), storage...), *o.optionWords["settings-updated"]})
	}
	spellbookCapture(t, "server-config-name-storage", rows, "54a3e66efa076844e6b7bb1fba078cd8855998edefc1e19e3d614ee3e675b404")
}

func TestServerConfigPasswordStorage(t *testing.T) {
	newServerOptionsOwner(t)
	storage := serverConfigOwnBytes(t, 0x5D4594, 3540, 44)
	type row struct {
		Input []uint16
		Bytes []byte
	}
	var rows []row
	for _, text := range [][]uint16{{0}, {'a', 0}, {'a', ' ', 'b', 0}, {0xd800, 0xdc00, 0}, {0xd800, 'X', 0}, {0xffff, 0}, {'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 0}} {
		for i := range storage {
			storage[i] = 0xa5
		}
		p, free := alloc.Make([]uint16{}, len(text))
		copy(p, text)
		result := legacy.PortTestServerConfigPointer("password-set", 0, unsafe.Pointer(&p[0]))
		free()
		if result != unsafe.Pointer(&storage[0]) || legacy.PortTestServerConfigPointer("password-get", 0, nil) != result {
			t.Fatal("password return")
		}
		expected := bytes.Repeat([]byte{0xa5}, len(storage))
		for i, v := range text {
			expected[2*i] = byte(v)
			expected[2*i+1] = byte(v >> 8)
		}
		if !bytes.Equal(storage, expected) {
			t.Fatal("UTF16 copy width", text, storage, expected)
		}
		rows = append(rows, row{append([]uint16(nil), text...), append([]byte(nil), storage...)})
	}
	if unsafe.Pointer(&storage[0]) != memmap.PtrOff(0x5D4594, 3540) {
		t.Fatal("password owner")
	}
	spellbookCapture(t, "server-config-password-storage", rows, "bdab2e327c4efa191d44103ae181d18dd760024f559970a048ed9573749968d0")
}
