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
	"time"
	"unsafe"
)

func TestGameStatisticsPlayerBytes(t *testing.T) {
	o := newMatchRosterOwner(t)
	words, restore := legacy.PortTestStatisticsGlobals()
	t.Cleanup(restore)
	data := serverConfigOwnBytes(t, 0x5D4594, 600124, 8192)
	pl := o.units[0].UpdateDataPlayer().Player
	type row struct {
		Flags        uint32
		Index, Value int32
		Op           string
		Return       uint32
		Data         []byte
	}
	var rows []row
	for _, flags := range []uint32{0, 4096, 8192, 12288, 8193} {
		for _, index := range []int32{-1, 0, 1, 31, 254} {
			for _, value := range []int32{-2147483648, -129, -128, -1, 0, 1, 127, 128, 255, 256, 2147483647} {
				for _, op := range []string{"completion", "participation"} {
					func() {
						defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
						for i := range data {
							data[i] = 0xa5
						}
						want := bytes.Clone(data)
						objectXferSetWord(pl.C(), 4648, uint32(index))
						got := legacy.PortTestStatisticsCall(op, pl.C(), nil, value)
						if flags&8192 != 0 && flags&4096 == 0 && index >= 0 {
							off := 21
							if op == "participation" {
								off = 28
							}
							want[int(index)*32+off] = byte(value)
						}
						if !bytes.Equal(data, want) || objectXferGetWord(pl.C(), 4648) != uint32(index) || *words["players"] != 0 {
							t.Fatal("player byte mutation", flags, index, value, op)
						}
						rows = append(rows, row{flags, index, value, op, got, bytes.Clone(data)})
					}()
				}
			}
		}
	}
	spellbookCapture(t, "game-statistics-player-bytes", rows, "97699f6cc76ee3ef05a579ba0dca337127aad4dea95e1c7369f1baa53ade9f87")
}
func TestGameStatisticsPlayerRegistration(t *testing.T) {
	o := newMatchRosterOwner(t)
	words, restore := legacy.PortTestStatisticsGlobals()
	t.Cleanup(restore)
	data := serverConfigOwnBytes(t, 0x5D4594, 600124, 8192)
	oldStreams := o.s.NetStr
	o.s.NetStr = &netstr.Streams{}
	t.Cleanup(func() { o.s.NetStr = oldStreams })
	t.Cleanup(o.s.NetStr.PortTestStatisticsPeers())
	oldIP := o.s.OwnIP
	o.s.OwnIP = netip.MustParseAddr("198.51.100.7")
	t.Cleanup(func() { o.s.OwnIP = oldIP })
	type row struct {
		Flags, Observer uint32
		Unit            int
		Graphics        bool
		Initial         uint32
		Count, Index    uint32
		Data            []byte
	}
	var rows []row
	for _, flags := range []uint32{0, 4096, 8192, 12288} {
		for _, observer := range []uint32{0, 1, 0x20, 0x21, 0xffffffff} {
			for unit := range o.units {
				for _, graphics := range []bool{false, true} {
					for _, initial := range []uint32{0, 1, 254} {
						t.Run(fmt.Sprintf("flags%x/observer%x/unit%d/noGraphics%v/index%d", flags, observer, unit, graphics, initial), func(t *testing.T) {
							defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
							noxflags.UnsetEngine(noxflags.EngineNoRendering)
							if graphics {
								noxflags.SetEngine(noxflags.EngineNoRendering)
							}
							clear(data)
							*words["players"] = initial
							pl := o.units[unit].UpdateDataPlayer().Player
							objectXferSetWord(pl.C(), 4648, 0xffffffff)
							pl.Field3680 = observer
							name := unsafe.Slice((*byte)(unsafe.Add(pl.C(), 2096)), 12)
							clear(name)
							copy(name, fmt.Sprintf("Name%d", unit))
							objectXferSetWord(pl.C(), 2068, 0x12345678+uint32(unit))
							*(*byte)(unsafe.Add(pl.C(), 2251)) = byte(unit + 1)
							before := uint32(time.Now().Unix())
							legacy.PortTestStatisticsCall("register", pl.C(), nil, 0)
							after := uint32(time.Now().Unix())
							active := flags&8192 != 0 && flags&4096 == 0
							wantCount, wantIndex := initial, uint32(0xffffffff)
							if active {
								wantCount++
								wantIndex = initial
								record := data[int(initial)*32 : int(initial+1)*32]
								if !bytes.Equal(record[:12], name) || binary.LittleEndian.Uint32(record[16:]) != 0x12345678+uint32(unit) || record[20] != byte(unit+1) || record[21] != 1 {
									t.Fatal("registration fields")
								}
								stamp := binary.LittleEndian.Uint32(record[24:])
								if stamp < before || stamp > after {
									t.Fatal("registration time outside call interval")
								}
								binary.LittleEndian.PutUint32(record[24:], 0) // Normalize only the checked wall-clock field.
								participating := observer&1 == 0 || observer&0x20 != 0
								if pl.PlayerInd == 31 && graphics {
									participating = false
								}
								wantParticipation := byte(0)
								if participating {
									wantParticipation = 1
								}
								if record[28] != wantParticipation {
									t.Fatal("participation policy")
								}
								ip := []byte{192, 0, 2, byte(int(pl.PlayerInd) + 4)}
								if pl.PlayerInd == 31 {
									ip = []byte{198, 51, 100, 7}
								}
								if binary.LittleEndian.Uint32(record[12:]) != binary.BigEndian.Uint32(ip) {
									t.Fatalf("IP field %x != %x", record[12:16], ip)
								}
								stable := bytes.Clone(data)
								legacy.PortTestStatisticsCall("register", pl.C(), nil, 0)
								if !bytes.Equal(stable, data) {
									t.Fatal("already registered player changed")
								}
							} else if !bytes.Equal(data, make([]byte, len(data))) {
								t.Fatal("disabled registration changed table")
							}
							if *words["players"] != wantCount || objectXferGetWord(pl.C(), 4648) != wantIndex {
								t.Fatal("registration index/count")
							}
							rows = append(rows, row{flags, observer, unit, graphics, initial, *words["players"], objectXferGetWord(pl.C(), 4648), bytes.Clone(data[:int(wantCount)*32])})
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "game-statistics-player-registration", rows, "85d729612d5893a0b63da8be711102e8b717137a6df578282bd7fe0995fdaf53")
}
