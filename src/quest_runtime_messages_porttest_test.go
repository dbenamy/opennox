//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestQuestRuntimeStageMessages(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	defer noxflags.PortTestGameFlags(0)()
	mode := memmap.PtrUint32(0x5D4594, 2388656)
	old := *mode
	t.Cleanup(func() { *mode = old })
	nameBuf := unsafe.Slice(memmap.PtrUint8(0x973F18, 3838), 32)
	titleBuf := unsafe.Slice(memmap.PtrUint8(0x973F18, 3806), 32)
	type row struct {
		Name   string
		Return uint64
		Queue  legacy.PortTestReliableReportState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-stage-messages", rows, "fa3ac020a403d390b1fa1495e0df37e5fdcb6c620b06853b8299cc60f7b662e3")
	}()
	for _, op := range []string{"sub_4D6880", "nox_game_sendQuestStage_4D6960"} {
		for _, stage := range []uint32{0, 65535, 65536, 0xffffffff} {
			for _, flag := range []uint32{0, 1, 0xffffffff} {
				for _, extra := range []uint32{0, 1, 0xffffffff} {
					for _, name := range []string{"", "Quest.map", strings.Repeat("a", 31)} {
						label := fmt.Sprintf("%s/stage%x/flag%x/extra%x/name%q", op, stage, flag, extra, name)
						t.Run(label, func(t *testing.T) {
							o.reset()
							*o.quest["202028"] = stage
							*mode = flag
							clear(nameBuf)
							clear(titleBuf)
							copy(nameBuf, name)
							title := "Title:" + name
							if len(title) > 31 {
								title = title[:31]
							}
							copy(titleBuf, title)
							if legacy.PortTestQuestRuntimeString("sub_4D6940") != name || legacy.PortTestQuestRuntimeString("sub_4D6950") != title {
								t.Fatal("map string interface")
							}
							rv := questRuntimeCall(op, nil, 1, extra)
							st := o.state()
							if rv != 1 || len(st.Nodes) != 1 {
								t.Fatal("stage enqueue")
							}
							got := st.Nodes[0].Data
							want := make([]byte, 69)
							want[0] = 240
							want[1] = 14
							if op == "sub_4D6880" {
								want[1] = 13
								if extra != 0 {
									want[4] = 1
								}
							}
							if flag != 0 {
								want[4] |= 2
							}
							binary.LittleEndian.PutUint16(want[2:], uint16(stage))
							copy(want[5:37], name)
							copy(want[37:], title)
							if !bytes.Equal(got, want) {
								t.Fatalf("stage payload %x want%x", got, want)
							}
							rows = append(rows, row{label, rv, st})
						})
					}
				}
			}
		}
	}
}
func TestQuestRuntimeSmallMessages(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	u := &o.units[0]
	defer noxflags.PortTestGameFlags(0)()
	type row struct {
		Name   string
		Return uint64
		Queue  legacy.PortTestReliableReportState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-small-messages", rows, "eead2433a83a44b7de51b470c7ee280b5fc0f6e77f811ae273168c59ad16352f")
	}()
	for _, op := range []string{"sub_4D7280", "sub_4D7450", "sub_4D6A20"} {
		for _, to := range []uint32{1, 31, 255} {
			for _, value := range []uint32{0, 1, 255, 256, 65535, 65536, 0xffffffff} {
				name := fmt.Sprintf("%s/to%d/value%x", op, to, value)
				t.Run(name, func(t *testing.T) {
					o.reset()
					objectXferSetWord(u.CObj(), 40, value)
					rv := questRuntimeCall(op, u, to, value)
					st := o.state()
					if rv != 1 || len(st.Nodes) != 1 {
						t.Fatal("small message enqueue", rv, len(st.Nodes))
					}
					want := []byte{240, 29, byte(value), byte(value >> 8)}
					if op == "sub_4D7280" {
						want = []byte{240, 24, byte(value)}
					}
					if op == "sub_4D6A20" {
						want[1] = 15
					}
					if !bytes.Equal(st.Nodes[0].Data, want) {
						t.Fatalf("message %x want%x", st.Nodes[0].Data, want)
					}
					rows = append(rows, row{name, rv, st})
				})
			}
		}
	}
}
func TestQuestRuntimeSettingsFallback(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	slot := (*unsafe.Pointer)(memmap.PtrOff(0x5D4594, 1556152))
	fallback := memmap.PtrOff(0x5D4594, 371516)
	flags := (*byte)(unsafe.Add(fallback, 100))
	old := *flags
	t.Cleanup(func() { *flags = old })
	alternate := o.record(t, 120)
	type row struct {
		Alternate                  bool
		Before, Fallback, Override byte
		Identity                   int
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-settings-fallback", rows, "712fbc659deb80a828301818f6f8a4feeb9d45724ba06fb9ac3bc4d17f11c012")
	}()
	for _, useAlternate := range []bool{false, true} {
		for _, bits := range []byte{0, 1, 16, 17, 255} {
			t.Run(fmt.Sprintf("alternate%t/bits%x", useAlternate, bits), func(t *testing.T) {
				*slot = nil
				*flags = bits
				*(*byte)(unsafe.Add(alternate, 100)) = bits
				if useAlternate {
					*slot = alternate
				}
				got := legacy.PortTestQuestRuntimeSettings()
				identity := 1
				want := fallback
				wantFlags := bits &^ 16
				if useAlternate {
					identity = 2
					want = alternate
					wantFlags = bits
				}
				if got != want || *flags != wantFlags || *(*byte)(unsafe.Add(alternate, 100)) != bits {
					t.Fatal("settings fallback identity or mutation")
				}
				rows = append(rows, row{useAlternate, bits, *flags, *(*byte)(unsafe.Add(alternate, 100)), identity})
			})
		}
	}
}
