//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestSummonRecordStorage(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	for mask := 0; mask < 16; mask++ {
		o.reset(t)
		first := -1
		free := -1
		for i := 0; i < 4; i++ {
			r := o.record(i)
			for j := range r {
				r[j] = 0xab000000 + uint32(100*i+j)
			}
			r[0] = uint32(11 + i%2)
			r[2] = uint32(mask >> i & 1)
			r[6] = uint32(i)
			if r[2] != 0 && first < 0 {
				first = i
			}
			if r[2] == 0 && free < 0 {
				free = i
			}
		}
		want := uint32(0)
		if first >= 0 {
			want = o.ptr(first)
		}
		got := o.call("sub_4C2D60")
		o.check(t, got == want, "first active record")
		o.check(t, o.call("sub_4C3260") == uint32(summonBool(mask != 0)), "any active record")
		for i := 0; i < 4; i++ {
			want = 0
			for j := i + 1; j < 4; j++ {
				if mask>>j&1 != 0 {
					want = o.ptr(j)
					break
				}
			}
			o.check(t, o.call("sub_4C2D90", o.ptr(i)) == want, "next active record")
		}
		for id := uint32(10); id <= 13; id++ {
			want = 0
			for i := 0; i < 4; i++ {
				if mask>>i&1 != 0 && o.record(i)[0] == id {
					want = o.ptr(i)
					break
				}
			}
			got = o.call("sub_4C31D0", id)
			o.check(t, got == want, "lookup first matching active code")
			rows = append(rows, o.snapshot(fmt.Sprintf("mask%d-id%d", mask, id), got))
		}
		got = o.call("sub_4C2F20")
		want = 0
		if free >= 0 {
			want = o.ptr(free)
			r := o.record(free)
			for i, v := range r {
				exp := uint32(0)
				if i == 2 {
					exp = 1
				}
				if i == 6 {
					exp = uint32(free)
				}
				o.check(t, v == exp, "allocation clears exactly one complete record")
			}
		}
		o.check(t, got == want, "first free record or full")
		rows = append(rows, o.snapshot(fmt.Sprintf("mask%d-allocate", mask), got))
		for i := 0; i < 4; i++ {
			before := append([]uint32(nil), o.record(i)...)
			ret := o.call("sub_4C3210", o.ptr(i))
			o.check(t, ret == o.ptr(i), "deactivation returns record")
			for j, v := range o.record(i) {
				if j == 2 {
					o.check(t, v == 0, "deactivation clears active flag")
				} else {
					o.check(t, v == before[j], "deactivation retains other fields")
				}
			}
			rows = append(rows, o.snapshot(fmt.Sprintf("mask%d-release%d", mask, i), ret))
		}
	}
	spellbookCapture(t, "summon-record-storage", rows, "843e0e1d9e78276b7458cf421404827623ee9bc7bac2adc6a01d485d390b508b")
}
func TestSummonGrid(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	pos, free := alloc.New([2]int32{})
	defer free()
	p := uint32(uintptr(unsafe.Pointer(pos)))
	for mask := 0; mask < 16; mask++ {
		for _, size := range []int{1, 2, 4} {
			w, h := 1, 1
			if size >= 2 {
				h = 2
			}
			if size == 4 {
				w = 2
			}
			for x := 0; x <= 2-w; x++ {
				for y := 0; y <= 2-h; y++ {
					o.reset(t)
					for i := 0; i < 4; i++ {
						if mask>>i&1 != 0 {
							*o.word(1321180 + 4*i) = o.ptr(i)
						}
					}
					*pos = [2]int32{int32(x), int32(y)}
					available := true
					for xx := x; xx < x+w; xx++ {
						for yy := y; yy < y+h; yy++ {
							if mask>>(yy+2*xx)&1 != 0 {
								available = false
							}
						}
					}
					ret := o.call("sub_4C30C0", p, uint32(size))
					o.check(t, ret == uint32(summonBool(available)), "rectangle availability")
					before := append([]uint32(nil), o.raw...)
					ret = o.call("sub_4C3030", p, uint32(size), o.ptr(0))
					o.check(t, ret == uint32(y+h), "paint returns bottom coordinate")
					for i := 0; i < 4; i++ {
						want := before[(1321180-1320988)/4+i]
						xx, yy := i/2, i%2
						if xx >= x && xx < x+w && yy >= y && yy < y+h {
							want = o.ptr(0)
						}
						o.check(t, *o.word(1321180 + 4*i) == want, "paint rectangle only")
					}
					for xx := -1; xx <= 2; xx++ {
						for yy := -1; yy <= 2; yy++ {
							*pos = [2]int32{int32(xx), int32(yy)}
							want := uint32(0)
							if xx >= 0 && xx < 2 && yy >= 0 && yy < 2 {
								want = *o.word(1321180 + 4*(yy+2*xx))
							}
							o.check(t, o.call("nox_xxx_wndSummonGet_4C2410", p) == want, "grid lookup bounds")
						}
					}
					o.check(t, o.call("nox_xxx_wndSummonGet_4C2410", 0) == 0, "nil grid lookup")
					rows = append(rows, o.snapshot(fmt.Sprintf("mask%d-size%d-%d,%d", mask, size, x, y), ret))
					ret = o.call("sub_4C2BF0")
					o.check(t, ret == uint32(uintptr(unsafe.Pointer(o.word(1321200)))), "clear returns original loop end address")
					for i := 0; i < 4; i++ {
						o.check(t, *o.word(1321180 + 4*i) == 0, "grid clear")
					}
				}
			}
		}
	}
	spellbookCapture(t, "summon-grid", rows, "e729538433320046a4757cbcea9aedf10c862abc55555c22db8d2729d4bd8388")
}
func TestSummonLayout(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	for code := 0; code < 256; code++ {
		o.reset(t)
		v := code
		total := 0
		for i := 0; i < 4; i++ {
			kind := v % 4
			v /= 4
			if kind == 0 {
				continue
			}
			size := []int{0, 1, 2, 4}[kind]
			r := o.record(i)
			r[0] = uint32(100 + i)
			r[2] = 1
			r[5] = uint32(size)
			r[6] = uint32(i)
			total += size
		}
		ret := o.call("sub_4C2F70")
		o.check(t, ret == 0, "layout iterator finishes at nil")
		if total <= 4 {
			occupied := 0
			for i := 0; i < 4; i++ {
				if *o.word(1321180 + 4*i) != 0 {
					occupied++
				}
			}
			o.check(t, occupied == total, "fitting records occupy expected cells")
		}
		rows = append(rows, o.snapshot(fmt.Sprintf("layout%d", code), ret))
	}
	spellbookCapture(t, "summon-layout", rows, "bb1c3ade62a36c87b876b18b796d0dd7204787581281f4682365abed3b88e5fe")
}
