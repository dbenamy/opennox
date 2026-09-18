//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netstr"
	"github.com/opennox/opennox/v1/legacy"
	"net/netip"
	"testing"
	"unsafe"
)

func statisticsPeers(t *testing.T, o *matchRosterOwner) {
	oldStreams := o.s.NetStr
	o.s.NetStr = &netstr.Streams{}
	t.Cleanup(func() { o.s.NetStr = oldStreams })
	t.Cleanup(o.s.NetStr.PortTestStatisticsPeers())
	oldIP, oldText := o.s.OwnIP, o.s.OwnIPStr
	o.s.OwnIP = netip.MustParseAddr("198.51.100.7")
	o.s.OwnIPStr = o.s.OwnIP.String()
	t.Cleanup(func() { o.s.OwnIP, o.s.OwnIPStr = oldIP, oldText })
}
func TestGameStatisticsEvents(t *testing.T) {
	o := newMatchRosterOwner(t)
	statisticsPeers(t, o)
	words, restore := legacy.PortTestStatisticsGlobals()
	t.Cleanup(restore)
	players := serverConfigOwnBytes(t, 0x5D4594, 600124, 8192)
	events := serverConfigOwnBytes(t, 0x5D4594, 608320, 131072)
	type row struct {
		Flags                 uint32
		Mask, Actor, Target   int
		Count, Events, Return uint32
		Indices               [3]uint32
		Players, Pairs        []byte
	}
	var rows []row
	for _, flags := range []uint32{0, 4096, 8192, 12288} {
		for mask := 0; mask < 8; mask++ {
			for actor := -1; actor < 3; actor++ {
				for target := -1; target < 3; target++ {
					func() {
						defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
						clear(players)
						clear(events)
						*words["players"], *words["events"] = 3, 0
						indices := [3]uint32{0, 1, 2}
						for i := range o.units {
							pl := o.units[i].UpdateDataPlayer().Player
							if mask&(1<<i) != 0 {
								indices[i] = 0xffffffff
							}
							objectXferSetWord(pl.C(), 4648, indices[i])
							objectXferSetWord(pl.C(), 2068, 0x12345678+uint32(i))
							name := unsafe.Slice((*byte)(unsafe.Add(pl.C(), 2096)), 12)
							clear(name)
							copy(name, fmt.Sprintf("Player%d", i))
							*(*byte)(unsafe.Add(pl.C(), 2251)) = byte(i + 1)
						}
						ptr := func(i int) unsafe.Pointer {
							if i < 0 {
								return nil
							}
							return o.units[i].UpdateDataPlayer().Player.C()
						}
						count, pairs := uint32(3), uint32(0)
						wantPair := [2]byte{}
						if flags&8192 != 0 && flags&4096 == 0 && actor >= 0 {
							if indices[actor] == 0xffffffff {
								indices[actor] = count
								count++
							}
							if target >= 0 && indices[target] == 0xffffffff {
								indices[target] = count
								count++
							}
							wantPair[0] = byte(indices[actor])
							wantPair[1] = 255
							if target >= 0 {
								wantPair[1] = byte(indices[target])
							}
							pairs = 1
						}
						got := legacy.PortTestStatisticsCall("event", ptr(actor), ptr(target), 0)
						if *words["players"] != count || *words["events"] != pairs || events[0] != wantPair[0] || events[1] != wantPair[1] {
							t.Fatal("event indices/count/pair", flags, mask, actor, target)
						}
						for i := range o.units {
							if objectXferGetWord(ptr(i), 4648) != indices[i] {
								t.Fatal("event player index")
							}
						}
						if pairs != 0 {
							for i := range o.units {
								if mask&(1<<i) == 0 || indices[i] == 0xffffffff {
									continue
								}
								r := players[indices[i]*32 : (indices[i]+1)*32]
								name := fmt.Sprintf("Player%d", i)
								ip := []byte{192, 0, 2, byte(int(o.units[i].UpdateDataPlayer().Player.PlayerInd) + 4)}
								if o.units[i].UpdateDataPlayer().Player.PlayerInd == 31 {
									ip = []byte{198, 51, 100, 7}
								}
								if string(r[:len(name)]) != name || r[len(name)] != 0 || binary.LittleEndian.Uint32(r[12:]) != binary.BigEndian.Uint32(ip) {
									t.Fatalf("event registration mask%d actor%d target%d: player%d row%d name/address %x", mask, actor, target, i, indices[i], r[:16])
								}
							}
						}
						rows = append(rows, row{flags, mask, actor, target, count, pairs, got, indices, bytes.Clone(players[:count*32]), bytes.Clone(events[:8])})
					}()
				}
			}
		}
	}
	spellbookCapture(t, "game-statistics-events", rows, "fc3ac4835fe589365efb0b7e2de29c472f41365bb0b7f09aa36fad554e236993")
}
