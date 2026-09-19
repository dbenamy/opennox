//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesEnchantmentWrite(t *testing.T) {
	o := newReliableReportsOwner(t)
	defer flags.PortTestGameFlags(0)()
	t.Cleanup(o.s.PortTestPlayerFileAbilities())
	order := []uint32{26, 0, 28, 5}
	t.Cleanup(legacy.PortTestPlayerFileEnchants(order))
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	p.NetCodeVal = u.NetCode
	classPtr := (*byte)(unsafe.Add(unsafe.Pointer(p), 2251))
	duration, free := alloc.New(server.DurSpell{})
	t.Cleanup(free)
	oldDur := o.s.Spells.Dur.List
	t.Cleanup(func() { o.s.Spells.Dur.List = oldDur })
	var rows []map[string]any
	for _, class := range []byte{0, 1, 2} {
		for _, gf := range []flags.GameFlag{0, 2048, 8192, 2048 | 8192} {
			for _, mask := range []uint32{0, 1 << 31, 1 << 26, 0xffffffff} {
				for _, active := range []bool{false, true} {
					o.reset()
					flags.ResetGame()
					flags.SetGame(gf)
					*classPtr = class
					u.Buffs = mask
					for i := range u.BuffsDur {
						u.BuffsDur[i] = uint16(i * 2001)
						u.BuffsPower[i] = byte(255 - i)
					}
					ad := o.s.Abils.GetFor(u)
					for i := range ad.Cooldowns {
						ad.Cooldowns[i] = i*101 - 200
					}
					ad.ExecList = nil
					o.s.Spells.Dur.List = nil
					first := server.ExecAbilityClass{Abil: 1, Active: 0, Frame: 177}
					fourth := server.ExecAbilityClass{Abil: 4, Active: 1, Frame: 17}
					if active {
						first.Next = &fourth
						fourth.Prev = &first
						ad.ExecList = &first
						*duration = server.DurSpell{Spell: 51, Target48: u, Field72: -123}
						o.s.Spells.Dur.List = duration
					}
					before := bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(u)), int(unsafe.Sizeof(*u))))
					present := byte(1)
					if gf&8192 != 0 {
						present = 0
					}
					want := []byte{5, 0, present}
					wantRet := uint32(1)
					if present != 0 {
						if gf&2048 == 0 {
							wantRet = 0
						} else {
							count := byte(0)
							var entries []byte
							for _, id := range order {
								if mask&(1<<id) == 0 {
									continue
								}
								count++
								name := server.EnchantID(id).String()
								entries = append(entries, byte(len(name)))
								entries = append(entries, name...)
								entries = binary.LittleEndian.AppendUint16(entries, u.BuffsDur[id])
								entries = append(entries, u.BuffsPower[id])
								if id == 26 {
									shield := uint32(100)
									if active {
										shield = 0xffffff85
									}
									entries = binary.LittleEndian.AppendUint32(entries, shield)
								}
							}
							want = append(want, count)
							want = append(want, entries...)
							if class == 0 {
								state := byte(0)
								remaining := uint32(0xffffffff)
								start := 1
								if active {
									state = 1
									remaining = 0xffffff96
									start = 2
								} // frame17 - frame123
								want = append(want, state, state)
								want = binary.LittleEndian.AppendUint32(want, remaining)
								for id := start; id < 6; id++ {
									want = binary.LittleEndian.AppendUint32(want, uint32(ad.Cooldowns[id]))
								}
							}
						}
					}
					ret, got, pos := playerFileSection(t, "nox_xxx_guiEnchantment_41B9C0", nil, uint32(uintptr(unsafe.Pointer(u))), 0)
					if ret != wantRet || pos != int64(len(want)) || !bytes.Equal(got, want) || !bytes.Equal(before, unsafe.Slice((*byte)(unsafe.Pointer(u)), len(before))) {
						t.Fatal("enchant write", class, gf, mask, active, ret, pos, len(want), got, want)
					}
					rows = append(rows, map[string]any{"class": class, "flags": uint32(gf), "buffs": mask, "active": active, "return": ret, "bytes": got, "position": pos})
				}
			}
		}
	}
	spellbookCapture(t, "player-files-enchantment-write", rows, "8029dca6811bc5d66d6111baa1f0f8f850a72360902fec4477ea8c3c7db7621e")
}
