//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"math"
	"testing"
	"unsafe"
)

func TestResourceDefinitionsLight(t *testing.T) {
	bias, scale := legacy.PortTestResourceLightConstants()
	ob, os := *bias, *scale
	t.Cleanup(func() { *bias = ob; *scale = os })
	p := memmap.PtrFloat64(0x581450, 9560)
	old := *p
	t.Cleanup(func() { *p = old })
	*p = 1.0 / 360
	*bias = math.Float64bits(.5)
	*scale = math.Float64bits(65536)
	typ, free := alloc.New(client.ObjectType{})
	defer free()
	raw := unsafe.Slice((*byte)(typ.C()), 128)
	var rows []map[string]any
	for _, kind := range []string{"direction", "penumbra"} {
		for deg := -2; deg <= 361; deg++ {
			for _, suffix := range []string{"", "tail"} {
				input := fmt.Sprintf(" \t%d%s", deg, suffix)
				s, freeS := alloc.CString(input)
				for i := range raw {
					raw[i] = 0xa5
				}
				ok := legacy.PortTestResourceClientParser(kind, typ, nil, unsafe.Pointer(s))
				freeS()
				limit, off := 360, 24
				if kind == "penumbra" {
					limit, off = 180, 26
				}
				want := bytes.Repeat([]byte{0xa5}, 128)
				valid := deg >= 0 && deg < limit
				if valid {
					binary.LittleEndian.PutUint16(want[off:], uint16((uint64(deg)*65536+180)/360))
					if kind == "direction" {
						resourcePut32(want, 16, 0)
					}
				}
				if ok != valid || !bytes.Equal(raw, want) {
					t.Fatalf("%s %d: %v %x != %x", kind, deg, ok, raw[16:28], want[16:28])
				}
				rows = append(rows, map[string]any{"kind": kind, "input": input, "ok": ok, "data": append([]byte(nil), raw...)})
			}
		}
		for _, input := range []string{"", "bad", "+", "--1"} {
			s, freeS := alloc.CString(input)
			for i := range raw {
				raw[i] = 0xa5
			}
			ok := legacy.PortTestResourceClientParser(kind, typ, nil, unsafe.Pointer(s))
			freeS()
			if ok || !bytes.Equal(raw, bytes.Repeat([]byte{0xa5}, 128)) {
				t.Fatal("invalid light", input)
			}
			rows = append(rows, map[string]any{"kind": kind, "input": input, "ok": ok, "data": append([]byte(nil), raw...)})
		}
	}
	spellbookCapture(t, "resource-definitions-light", rows, "9bc95a1eda3d846bcdfabe3f5ac8580695a1d2234d72a5b4aa53bc2e0732d0ed")
}
func TestResourceDefinitionsClientUpdate(t *testing.T) {
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 175072), 8)
	old := append([]uint32(nil), table...)
	defer copy(table, old)
	for i, name := range []string{"UpdateA", "UpdateB", "UpdateA"} {
		p, free := alloc.CString(name)
		defer free()
		table[2*i] = uint32(uintptr(unsafe.Pointer(p)))
		table[2*i+1] = uint32(101 + i)
	}
	table[6] = 0
	table[7] = 0xdeadbeef
	typ, free := alloc.New(client.ObjectType{})
	defer free()
	raw := unsafe.Slice((*byte)(typ.C()), 128)
	var rows []map[string]any
	for _, input := range []string{"UpdateA", "UpdateB", "UpdateA ignored", " \t\r\nUpdateB  ", "updatea", "missing"} {
		for i := range raw {
			raw[i] = 0xa5
		}
		p, free := alloc.CString(input)
		ok := legacy.PortTestResourceClientParser("update", typ, nil, unsafe.Pointer(p))
		scratch := append([]byte(nil), unsafe.Slice(p, len(input)+1)...)
		free()
		value := uint32(0)
		switch input {
		case "UpdateA", "UpdateA ignored":
			value = 101
		case "UpdateB", " \t\r\nUpdateB  ":
			value = 102
		}
		want := bytes.Repeat([]byte{0xa5}, 128)
		if value != 0 {
			resourcePut32(want, 100, value)
		}
		if ok != (value != 0) || !bytes.Equal(raw, want) {
			t.Fatal("client update", input)
		}
		rows = append(rows, map[string]any{"input": input, "ok": ok, "data": append([]byte(nil), raw...), "scratch": scratch})
	}
	// Empty table accepts empty input without dereferencing strtok's nil result.
	table[0] = 0
	p, freeS := alloc.CString("")
	defer freeS()
	if legacy.PortTestResourceClientParser("update", typ, nil, unsafe.Pointer(p)) {
		t.Fatal("empty update table")
	}
	spellbookCapture(t, "resource-definitions-client-update", rows, "3bc939f3a3f9577cdc472555eaef96e9c68cebc8f5ad4408d6fba387818996c8")
}
