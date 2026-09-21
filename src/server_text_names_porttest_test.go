//go:build porttest

package opennox

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf16"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestServerTextPlayerNames(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(1))
	oldList := o.s.Objs.List
	t.Cleanup(func() { o.s.Objs.List = oldList })
	names := [][]uint16{{}, {'A'}, {'A', 'l', 'i', 'c', 'e'}, {'a', 'l', 'i', 'c', 'e'}, {0xc5, 'n', 'g'}, {0x100, 'x'}, {0xffff, 0x8001}, {'a', 0, 'b'}, utf16.Encode([]rune("Rune😀")), utf16.Encode([]rune(strings.Repeat("Z", 27))), utf16.Encode([]rune(strings.Repeat("Y", 28))), utf16.Encode([]rune(strings.Repeat("Q", 29)))}
	record, free := alloc.Make([]byte{}, 80)
	t.Cleanup(free)
	result, freeResult := alloc.New([3]uint32{})
	t.Cleanup(freeResult)
	type row struct {
		Style, Active, Op int
		Missing           bool
		Query             string
		Return, Team      int
		Record            []byte
	}
	var rows []row
	for style, name := range names {
		var canonical []byte
		for _, ch := range name[:min(len(name), 28)] {
			if ch == 0 {
				break
			}
			canonical = append(canonical, byte(ch))
		}
		if i := bytes.IndexByte(canonical, 0); i >= 0 {
			canonical = canonical[:i]
		}
		queries := []string{string(canonical), strings.ToLower(string(canonical)), string(canonical) + "!", "", "Other"}
		for active := 0; active < 8; active++ {
			for _, missing := range []bool{false, true} {
				for i := 0; i < 32; i++ {
					o.s.Players.ByIndRaw(ntype.PlayerInd(i)).Active = 0
				}
				for i := range o.units {
					u := &o.units[i]
					p := u.UpdateDataPlayer().Player
					p.Active = byte((active >> i) & 1)
					p.NetCodeVal = u.NetCode
					if missing {
						p.NetCodeVal += 10000
					}
					clear(p.NameFinal[:])
					if i < 2 {
						copy(p.NameFinal[:], name)
					} else {
						copy(p.NameFinal[:], []uint16{'O', 't', 'h', 'e', 'r'})
					}
					*(*byte)(unsafe.Add(p.C(), 2251)) = []byte{0, 2, 255}[i]
					u.TeamPtr().ID = server.TeamID([]byte{0, 17, 255}[i])
					u.ObjNext = nil
					if i+1 < len(o.units) {
						u.ObjNext = &o.units[i+1]
					}
				}
				o.s.Objs.List = &o.units[0]
				if missing {
					o.s.Objs.List = nil
				}
				for _, query := range queries {
					first := -1
					for i := range o.units {
						if active&(1<<i) == 0 {
							continue
						}
						wantName := string(canonical)
						if i == 2 {
							wantName = "Other"
						}
						if wantName == query {
							first = i
							break
						}
					}
					if missing {
						first = -1
					}
					for op := 0; op < 2; op++ {
						for i := range record {
							record[i] = 0xa5
						}
						copy(record[2:], query)
						record[2+len(query)] = 0
						before := bytes.Clone(record)
						*result = [3]uint32{0x12345678, 0x76543210, 0x87654321}
						restoreLookup := legacy.PortTestCreatureXferLookupOwner()
						got := 0
						if op == 0 {
							got = legacy.PlayerInfoStructParser_0(unsafe.Pointer(&record[0]))
						} else {
							got = legacy.PlayerInfoStructParser_1(unsafe.Pointer(&record[0]), (*int32)(unsafe.Pointer(&result[1])))
						}
						restoreLookup()
						want := bytes.Clone(before)
						wantRet, wantTeam, teamID := 0, uint32(0x76543210), -1
						if first >= 0 {
							wantRet = 1
							want[0] = []byte{0, 2, 255}[first]
							want[1] = []byte{0, 17, 255}[first]
							if op == 1 {
								wantTeam = uint32(uintptr(o.units[first].TeamPtr().C()))
								teamID = first
							}
						}
						if got != wantRet || !bytes.Equal(record, want) || result[1] != wantTeam || result[0] != 0x12345678 || result[2] != 0x87654321 {
							t.Fatalf("player-name lookup style=%d active=%d missing=%v query=%q op=%d: return=%d/%d team=%x/%x record=%x/%x", style, active, missing, query, op, got, wantRet, result[1], wantTeam, record, want)
						}
						rows = append(rows, row{style, active, op, missing, query, got, teamID, bytes.Clone(record)})
					}
				}
			}
		}
	}
	// Invalid record/output sentinels return without touching caller-owned words.
	for _, ptr := range []unsafe.Pointer{nil, unsafe.Pointer(uintptr(0xfffffffe))} {
		result[1] = 0x76543210
		if legacy.PlayerInfoStructParser_0(ptr) != 0 || legacy.PlayerInfoStructParser_1(ptr, (*int32)(unsafe.Pointer(&result[1]))) != 0 || result[1] != 0x76543210 {
			t.Fatal("name lookup sentinel")
		}
	}
	if legacy.PlayerInfoStructParser_1(unsafe.Pointer(&record[0]), nil) != 0 {
		t.Fatal("nil name lookup output")
	}
	interactionCapture(t, "server-text-player-names", rows)
	t.Logf("%d name lookup cases", len(rows))
}

func TestServerTextFlagIndex(t *testing.T) {
	type row struct {
		Index int
		Value uint32
	}
	var rows []row
	// Negative exponents below -31 would divide by zero in C and are not valid
	// settings inputs. Cover its defined domain before replacing the helper.
	for i := -31; i <= 31; i++ {
		want := uint32(0)
		if i >= 0 {
			want = uint32(1) << uint(i)
		}
		got := legacy.GetFlagValueFromFlagIndex(i)
		if got != want {
			t.Fatal("flag index", i, got, want)
		}
		rows = append(rows, row{i, got})
	}
	interactionCapture(t, "server-text-flag-index", rows)
}
