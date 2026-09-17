//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestQuestRuntimeReconnectCleanup(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	defer noxflags.PortTestGameFlags(noxflags.GameModeQuest)()
	type row struct {
		Name      string
		Stats     [3][11]uint32
		Relations [3][224]byte
		Queue     legacy.PortTestReliableReportState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-reconnect", rows, "22693964f97bcfe27fd30a796b1b478eabc8afe3233d7a07d0034ac921f1b29a")
	}()
	for who := 0; who < 3; who++ {
		for _, stage := range []uint32{0, 65535, 0xffffffff} {
			for _, participation := range []uint32{0, 1, 2} {
				name := fmt.Sprintf("who%d/stage%x/participation%d", who, stage, participation)
				t.Run(name, func(t *testing.T) {
					o.reset()
					*o.quest["202028"] = stage
					u := &o.units[who]
					slot := int(u.UpdateDataPlayer().Player.PlayerInd)
					var want row
					for i := range o.units {
						pl := o.units[i].UpdateDataPlayer().Player
						objectXferSetWord(pl.C(), 4792, participation)
						for j := 0; j < 11; j++ {
							objectXferSetWord(pl.C(), 4652+4*j, uint32(10+i+j))
						}
						want.Stats[i] = questRuntimeStats(pl.C())
						if i == who {
							want.Stats[i] = [11]uint32{}
							want.Stats[i][9] = stage
							want.Stats[i][10] = 63
						}
						b := unsafe.Slice((*byte)(unsafe.Add(o.units[i].UpdateData, 324)), 224)
						for j := range b {
							b[j] = byte(1 + j)
						}
						copy(want.Relations[i][:], b)
						clear(want.Relations[i][4*slot : 4*slot+4])
						for k := 0; k < 3; k++ {
							want.Relations[i][128+32*k+slot] = 0
						}
					}
					if rv := questRuntimeCall("sub_4D79C0", u); rv != 0 {
						t.Fatal("reconnect return", rv)
					}
					var got row
					got.Name = name
					got.Queue = o.state()
					for i := range o.units {
						got.Stats[i] = questRuntimeStats(o.units[i].UpdateDataPlayer().Player.C())
						copy(got.Relations[i][:], unsafe.Slice((*byte)(unsafe.Add(o.units[i].UpdateData, 324)), 224))
					}
					if got.Stats != want.Stats || got.Relations != want.Relations {
						t.Fatal("reconnect cleanup scope")
					}
					if len(got.Queue.Nodes) != 1 {
						t.Fatal("reconnect notification count", len(got.Queue.Nodes))
					}
					data := got.Queue.Nodes[0].Data
					expected := []byte{240, 1, 0, 0}
					binary.LittleEndian.PutUint16(expected[2:], uint16(u.NetCode))
					if !bytes.Equal(data, expected) {
						t.Fatalf("reconnect payload%x want%x", data, expected)
					}
					rows = append(rows, got)
				})
			}
		}
	}
}
