//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestGameStatisticsEventFlush(t *testing.T) {
	type row struct {
		Before, Actor, Target                     int
		Result, Players, Events, Arrays, Sequence uint32
		Header, Rows                              []byte
	}
	var rows []row
	for _, before := range []int{253, 254, 255} {
		for _, actor := range []int{0, 2} {
			for _, target := range []int{-1, 1, 2} {
				t.Run(fmt.Sprintf("count%d/actor%d/target%d", before, actor, target), func(t *testing.T) {
					o := newMatchRosterOwner(t)
					statisticsPeers(t, o)
					words, _ := statisticsRandomOwner(t)
					statisticsOwnReportTags(t)
					t.Cleanup(noxflags.PortTestGameFlags(8192))
					clock, restore := legacy.PortTestStatisticsClock()
					t.Cleanup(restore)
					clock(1700000000)
					raw := serverConfigOwnBytes(t, 0x5D4594, 599476, 640)
					clear(raw)
					players := serverConfigOwnBytes(t, 0x5D4594, 600124, 8192)
					clear(players)
					clear(serverConfigOwnBytes(t, 0x5D4594, 608320, 131072))
					clear(serverConfigOwnBytes(t, 0x5D4594, 741300, 32))
					t.Cleanup(func() {
						names := *(*unsafe.Pointer)(unsafe.Pointer(&raw[608]))
						if names != nil {
							for i := 0; i < int(*words["array-count"]); i++ {
								legacy.PortTestStatisticsFree(*(*unsafe.Pointer)(unsafe.Add(names, 4*i)))
							}
						}
						for off := 608; off <= 636; off += 4 {
							legacy.PortTestStatisticsFree(*(*unsafe.Pointer)(unsafe.Pointer(&raw[off])))
						}
					})
					index := uint32(0)
					for pl := o.s.Players.First(); pl != nil; pl = o.s.Players.Next(pl) {
						objectXferSetWord(pl.C(), 4648, index)
						index++
						name := unsafe.Slice((*byte)(unsafe.Add(pl.C(), 2096)), 12)
						clear(name)
						copy(name, fmt.Sprintf("Peer%d", pl.PlayerInd))
						pl.Field3680 = 0
					}
					*words["players"] = uint32(before)
					a := o.units[actor].UpdateDataPlayer().Player
					objectXferSetWord(a.C(), 4648, 0xffffffff)
					var b unsafe.Pointer
					if target >= 0 {
						b = o.units[target].UpdateDataPlayer().Player.C()
					}
					result := legacy.PortTestStatisticsCall("event", a.C(), b, 0)
					expected := uint32(before + 1)
					if before >= 254 {
						expected = 4
						if *words["events"] != 0 || *words["array-count"] != uint32(before+1) || memmap.Uint32(0x5D4594, 741668) != 1 || binary.LittleEndian.Uint16(raw[6:]) != uint16(before+1) || result != 0 {
							t.Fatal("automatic flush threshold/state")
						}
					} else if *words["events"] != 1 || *words["array-count"] != 0 || memmap.Uint32(0x5D4594, 741668) != 0 || result != expected {
						t.Fatal("flushed before threshold")
					}
					if *words["players"] != expected {
						t.Fatal("automatic roster rebuild")
					}
					rows = append(rows, row{before, actor, target, result, *words["players"], *words["events"], *words["array-count"], memmap.Uint32(0x5D4594, 741668), bytes.Clone(raw[:608]), bytes.Clone(players[:32*expected])})
				})
			}
		}
	}
	spellbookCapture(t, "game-statistics-event-flush", rows, "8eea354d2a24357023e830948deb8544a9e0a3dc472b83147e17739a2564c458")
}
