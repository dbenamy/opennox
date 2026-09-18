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

func TestGameStatisticsSessionLifecycle(t *testing.T) {
	type row struct {
		Flags, Clock                  uint32
		Step                          string
		Result                        uint32
		Globals                       map[string]uint32
		Match, Quest, Players, Events []byte
		Lengths                       [2]uint32
	}
	var rows []row
	for _, flags := range []uint32{0, 4096, 8192, 12288} {
		for _, start := range []uint32{0, 1, 1700000000, 0xffffffe0} {
			t.Run(fmt.Sprintf("flags%x/start%x", flags, start), func(t *testing.T) {
				o := newMatchRosterOwner(t)
				statisticsPeers(t, o)
				words, _ := statisticsRandomOwner(t)
				statisticsOwnReportTags(t)
				t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
				setClock, restore := legacy.PortTestStatisticsClock()
				t.Cleanup(restore)
				setClock(start)
				match := serverConfigOwnBytes(t, 0x5D4594, 599476, 640)
				clear(match)
				quest := serverConfigOwnBytes(t, 0x5D4594, 739396, 580)
				clear(quest)
				players := serverConfigOwnBytes(t, 0x5D4594, 600124, 8192)
				clear(players)
				events := serverConfigOwnBytes(t, 0x5D4594, 608320, 131072)
				clear(events)
				clear(serverConfigOwnBytes(t, 0x5D4594, 741300, 32))
				serverConfigOwnBytes(t, 0x587000, 60072, 4)
				*memmap.PtrUint32(0x587000, 60072) = 0x12345678
				slot := serverConfigOwnBytes(t, 0x5D4594, 371380, 124)
				clear(slot)
				copy(slot, "arena")
				copy(slot[9:], "Stats fixture")
				binary.LittleEndian.PutUint16(slot[52:], 0x100)
				settings := serverConfigOwnBytes(t, 0x5D4594, 371516, 144)
				clear(settings)
				binary.LittleEndian.PutUint32(settings[40:], 0x23456789)
				settings[103], settings[104] = 16, 32
				for pl := o.s.Players.First(); pl != nil; pl = o.s.Players.Next(pl) {
					objectXferSetWord(pl.C(), 4648, 0xffffffff)
					objectXferSetWord(pl.C(), 4792, 0)
					pl.Field3680 = 0
					name := unsafe.Slice((*byte)(unsafe.Add(pl.C(), 2096)), 12)
					clear(name)
					copy(name, fmt.Sprintf("Peer%d", pl.PlayerInd))
					for off := 4652; off <= 4692; off += 4 {
						objectXferSetWord(pl.C(), off, 0)
					}
					objectXferSetWord(pl.C(), 4688, 1)
				}
				for i := 0; i < 2; i++ {
					objectXferSetWord(o.units[i].UpdateDataPlayer().Player.C(), 4792, 1)
				}
				t.Cleanup(func() {
					if *(*unsafe.Pointer)(unsafe.Pointer(&quest[536])) != nil {
						legacy.PortTestStatisticsCall("quest-clear", unsafe.Pointer(&quest[0]), nil, 0)
					}
					names := *(*unsafe.Pointer)(unsafe.Pointer(&match[608]))
					if names != nil {
						for i := 0; i < int(*words["array-count"]); i++ {
							legacy.PortTestStatisticsFree(*(*unsafe.Pointer)(unsafe.Add(names, 4*i)))
						}
					}
					for off := 608; off <= 636; off += 4 {
						legacy.PortTestStatisticsFree(*(*unsafe.Pointer)(unsafe.Pointer(&match[off])))
					}
				})
				capture := func(step string, clock, result uint32) {
					globals := map[string]uint32{}
					for k, p := range words {
						globals[k] = *p
					}
					rows = append(rows, row{flags, clock, step, result, globals, bytes.Clone(match[:608]), bytes.Clone(quest[:536]), bytes.Clone(players[:32*int(*words["players"])]), bytes.Clone(events[:8]), [2]uint32{memmap.Uint32(0x5D4594, 741300), memmap.Uint32(0x5D4594, 741312)}})
				}
				legacy.PortTestStatisticsCall("initialize", nil, nil, 0)
				if *words["start"] != start {
					t.Fatal("session start time")
				}
				onlineMatch := flags&8192 != 0 && flags&4096 == 0
				if onlineMatch {
					count := uint32(1)
					host := o.s.Players.ByInd(31)
					if objectXferGetWord(host.C(), 4648) != 0 {
						t.Fatal("host must register first")
					}
					for pl := o.s.Players.First(); pl != nil; pl = o.s.Players.Next(pl) {
						if pl.PlayerInd == 31 {
							continue
						}
						if objectXferGetWord(pl.C(), 4648) != count {
							t.Fatal("initial player order")
						}
						count++
					}
					if *words["players"] != count || binary.LittleEndian.Uint16(match[6:]) != 0 {
						t.Fatal("initial roster/report count")
					}
				} else if *words["players"] != 0 {
					t.Fatal("inactive match registration")
				}
				capture("initialize", start, 0)
				setClock(start + 45)
				if onlineMatch {
					*words["events"] = 3
					copy(events, []byte{0, 1, 2, 255, 1, 0})
				}
				legacy.PortTestStatisticsCall("flush", nil, nil, 0)
				if onlineMatch {
					if *words["events"] != 0 || memmap.Uint32(0x5D4594, 741668) != 1 || binary.LittleEndian.Uint32(match[28:]) != 45 {
						t.Fatal("intermediate report reset/time")
					}
					if *words["players"] != 4 {
						t.Fatal("rebuilt player registration count")
					}
					for i := 0; i < 4; i++ {
						if binary.LittleEndian.Uint32(players[32*i+24:]) != start+45 {
							t.Fatal("re-registration clock")
						}
					}
				}
				capture("flush", start+45, 0)
				setClock(start + 100)
				result := legacy.PortTestStatisticsCall("finish", nil, nil, 0)
				if result != 1 {
					t.Fatal("report completion result")
				}
				if flags&4096 != 0 {
					if binary.LittleEndian.Uint32(quest[20:]) != 100 || binary.LittleEndian.Uint16(quest) != 2 || memmap.Uint32(0x5D4594, 741672) != 1 {
						t.Fatal("quest final report count/time/sequence")
					}
				} else if onlineMatch {
					if binary.LittleEndian.Uint32(match[28:]) != 100 || memmap.Uint32(0x5D4594, 741668) != 2 {
						t.Fatal("match final report time/sequence")
					}
				}
				capture("finish", start+100, result)
			})
		}
	}
	spellbookCapture(t, "game-statistics-session-lifecycle", rows, "25eb953c9a395a7c60cb3b0689ede5f1b97e2a87a7340dfebc7c74b43b827d8a")
}
