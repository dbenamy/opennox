//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

type matchRosterMessageRow struct {
	Name   string
	Return uint32
	Queue  legacy.PortTestReliableReportState
}

func TestMatchRosterGUISettingsMessage(t *testing.T) {
	o := newMatchRosterOwner(t)
	var rows []matchRosterMessageRow
	defer func() {
		spellbookCapture(t, "match-roster-gui-settings", rows, "b15fecbfaa412679b3c7c0c6d031f502f48df5672f9aad585570346f819e7104")
	}()
	for _, mode := range []uint32{0, 1, 127, 128, 255, 256, 0xffffffff} {
		for _, to := range []uint32{1, 7, 31, 159, 255} {
			for _, seed := range []byte{0, 1, 255} {
				name := fmt.Sprintf("mode%x/to%d/seed%d", mode, to, seed)
				t.Run(name, func(t *testing.T) {
					o.reset()
					data := make([]byte, 66)
					for i := range data {
						data[i] = seed + byte(i*31)
					}
					saved := bytes.Clone(data)
					rv := matchRosterCall("gui-settings", nil, nil, unsafe.Pointer(&data[4]), mode, to)
					state := o.state()
					want := append([]byte{177, byte(mode)}, data[4:62]...)
					if rv != 1 || len(state.Nodes) != 1 || !bytes.Equal(state.Nodes[0].Data, want) {
						t.Fatal("settings message", rv, state)
					}
					if !bytes.Equal(data, saved) {
						t.Fatal("settings input changed")
					}
					rows = append(rows, matchRosterMessageRow{name, rv, state})
				})
			}
		}
	}
}
func TestMatchRosterSimpleObjectMessage(t *testing.T) {
	o := newMatchRosterOwner(t)
	u := &o.units[0]
	var rows []matchRosterMessageRow
	defer func() {
		spellbookCapture(t, "match-roster-simple-object", rows, "5ac1e01b3e3b297734c0927c89e3f98046f96c40ce6b763a2f7a9c4f63449d87")
	}()
	bits := []uint32{0, 0x80000000, 0x3fc00000, 0xbfc00000, 0x477fffff, 0x47800000, 0x4effffff, 0x4f000000, 0x7f800000, 0x7fc12345}
	for _, x := range bits {
		for _, y := range bits {
			for _, code := range []uint32{0, 0x12345678, 0xffffffff} {
				name := fmt.Sprintf("x%x/y%x/code%x", x, y, code)
				t.Run(name, func(t *testing.T) {
					o.reset()
					u.PosVec = types.Pointf{X: math.Float32frombits(x), Y: math.Float32frombits(y)}
					u.NetCode = code
					u.TypeInd = 0x1234
					rv := matchRosterCall("simple-object", u, nil, nil, 1)
					state := o.state()
					want := []byte{47, 0, 0, 0x34, 0x12, 0, 0, 0, 0}
					binary.LittleEndian.PutUint16(want[1:], uint16(o.s.GetUnitNetCode(u)))
					binary.LittleEndian.PutUint16(want[5:], uint16(floatIntExpected(x)))
					binary.LittleEndian.PutUint16(want[7:], uint16(floatIntExpected(y)))
					if rv != 1 || len(state.Nodes) != 1 || !bytes.Equal(state.Nodes[0].Data, want) {
						t.Fatal("object message", rv, state)
					}
					rows = append(rows, matchRosterMessageRow{name, rv, state})
				})
			}
		}
	}
}
func TestMatchRosterWallMessages(t *testing.T) {
	o := newMatchRosterOwner(t)
	var rows []matchRosterMessageRow
	defer func() {
		spellbookCapture(t, "match-roster-wall-messages", rows, "e593205eb1d45d26ebed9d2b05abe2dc42a1f27ec0d113925b626234188a1a25")
	}()
	for _, op := range []string{"wall-open", "wall-close", "wall-destroy"} {
		for _, id := range []uint16{0, 1, 32767, 32768, 65535} {
			for _, to := range []int{1, 31, 32} {
				name := fmt.Sprintf("%s/id%d/to%d", op, id, to)
				t.Run(name, func(t *testing.T) {
					o.reset()
					wall := server.Wall{Field10: id}
					rv := uint32(0)
					count := 3
					opcode := byte(59)
					switch op {
					case "wall-open":
						rv = matchRosterCall(op, nil, nil, unsafe.Pointer(&wall))
					case "wall-close":
						opcode = 60
						rv = matchRosterCall(op, nil, nil, unsafe.Pointer(&wall))
					case "wall-destroy":
						opcode = 58
						o.s.Nox_xxx_wallSendDestroyed_4DF0A0(&wall, to)
						if to != 32 {
							count = 1
						}
					}
					state := o.state()
					want := []byte{opcode, byte(id), byte(id >> 8)}
					if rv != 0 || len(state.Nodes) != count {
						t.Fatal("wall recipients", rv, len(state.Nodes), count)
					}
					for _, n := range state.Nodes {
						if !bytes.Equal(n.Data, want) {
							t.Fatal("wall payload", n.Data, want)
						}
					}
					if wall.Field10 != id {
						t.Fatal("wall changed")
					}
					rows = append(rows, matchRosterMessageRow{name, rv, state})
				})
			}
		}
	}
}
func TestMatchRosterFlagState(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name   string
		Return uint32
		Record [6]byte
		Queue  legacy.PortTestReliableReportState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "match-roster-flag-state", rows, "c4d4bd49637339ccd49e82486085fe1a6f9fba79c5c6aa59b027d2feb6f1df6c")
	}()
	all := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1567736), 106)
	for _, index := range []byte{0, 1, 2, 8, 16} {
		for _, a := range []uint32{0, 1, 127, 128, 255, 256, 0xffffffff} {
			for _, b := range []uint32{0, 1, 127, 128, 255, 256, 0xffffffff} {
				for _, value := range []uint32{0, 65535, 65536, 0xffffffff} {
					name := fmt.Sprintf("index%d/a%x/b%x/value%x", index, a, b, value)
					t.Run(name, func(t *testing.T) {
						o.reset()
						for i := range all {
							all[i] = 0xa5
						}
						rv := matchRosterCall("flag-state", nil, nil, nil, uint32(index), a, b, value)
						base, entry := legacy.PortTestMatchRosterFlagPointers(index)
						if base != unsafe.Pointer(&all[0]) || entry != unsafe.Pointer(&all[4+6*int(index)]) {
							t.Fatal("flag record address")
						}
						var got [6]byte
						copy(got[:], all[4+6*int(index):])
						want := [6]byte{index, byte(b), byte(a), 0xa5, byte(value), byte(value >> 8)}
						if got != want {
							t.Fatal("flag record", got, want)
						}
						for i, v := range all {
							if (i < 4+6*int(index) || i >= 10+6*int(index)) && v != 0xa5 {
								t.Fatal("adjacent flag state changed", i)
							}
						}
						state := o.state()
						payload := []byte{216, byte(a), index, byte(b), byte(value), byte(value >> 8)}
						if rv != 1 || len(state.Nodes) != 1 || !bytes.Equal(state.Nodes[0].Data, payload) {
							t.Fatal("flag report", rv, state)
						}
						rows = append(rows, row{name, rv, got, state})
					})
				}
			}
		}
	}
}
