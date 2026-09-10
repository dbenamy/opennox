//go:build porttest

package opennox

import (
	"math/rand"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestNetworkCodeBitsABI(t *testing.T) {
	for low := uint32(0); low < 65536; low++ {
		for _, upper := range []uint32{0, 0x00010000, 0x7fff0000, 0x80000000, 0xffff0000} {
			input := upper | low
			clear, flag := legacy.PortTestNetworkBits(input)
			wantClear := low % 32768
			wantFlag := low / 32768
			if clear != wantClear || flag != wantFlag {
				t.Fatalf("input=%08x clear=%x flag=%x", input, clear, flag)
			}
			if uint32(nox_xxx_netClearHighBit_578B30(uint16(input))) != wantClear || nox_xxx_netTestHighBit_578B70(uint16(input)) != (wantFlag != 0) {
				t.Fatal("native bit helper disagrees", input)
			}
		}
	}
}

func TestNetworkClientCodeABI(t *testing.T) {
	specs := []legacy.PortTestClientCode{{Nil: true}, {Nil: true, Code: 0xffffffff, Class: 0xffffffff}}
	classes := []uint32{0, 0x20400000, 0xffffffff}
	for bit := 0; bit < 32; bit++ {
		classes = append(classes, uint32(1)<<bit, ^(uint32(1) << bit))
	}
	codes := []uint32{0, 1, 2, 0x7ffe, 0x7fff, 0x8000, 0x8001, 0xffff, 0x10000, 0x10001, 0x17fff, 0x7fffffff, 0x80000000, 0xffffffff}
	for _, code := range codes {
		for _, class := range classes {
			specs = append(specs, legacy.PortTestClientCode{Code: code, Class: class})
		}
	}
	gen := rand.New(rand.NewSource(0x578b00))
	for i := 0; i < 2048; i++ {
		code := gen.Uint32()
		if i%2 == 0 {
			code %= 32768
		}
		specs = append(specs, legacy.PortTestClientCode{Code: code, Class: gen.Uint32()})
	}
	got := legacy.PortTestClientCodes(specs)
	if len(got) != len(specs) {
		t.Fatal("missing client snapshots")
	}
	for i, spec := range specs {
		want := uint32(0)
		if !spec.Nil && spec.Code < 32768 {
			want = spec.Code
			if (spec.Class>>22)%2 == 1 || (spec.Class>>29)%2 == 1 {
				want += 32768
			}
		}
		if got[i].Code != want || !got[i].Unchanged {
			t.Fatalf("case=%d spec=%+v got=%+v want=%x", i, spec, got[i], want)
		}
	}
}
