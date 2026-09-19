//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"testing"
	"unsafe"
)

func TestPlayerFilesEnchantmentReadGates(t *testing.T) {
	o := newReliableReportsOwner(t)
	defer flags.PortTestGameFlags(0)()
	t.Cleanup(o.s.PortTestPlayerFileAbilities())
	oldDur := noxServer.spells.duration
	noxServer.spells.duration.Init(noxServer)
	t.Cleanup(func() { noxServer.spells.duration = oldDur })
	o.s.Spells.Dur.Init()
	t.Cleanup(o.s.Spells.Dur.Free)
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	p.NetCodeVal = u.NetCode
	classPtr := (*byte)(unsafe.Add(unsafe.Pointer(p), 2251))
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 3, 4, 5, 6, 0x8000, 0xffff} {
		for _, class := range []byte{0, 1} {
			for _, gf := range []flags.GameFlag{0, 2048, 8192, 2048 | 8192} {
				for _, present := range []byte{0, 1} {
					o.reset()
					flags.ResetGame()
					flags.SetGame(gf)
					*classPtr = class
					ad := o.s.Abils.GetFor(u)
					ad.ExecList = nil
					for i := range ad.Cooldowns {
						ad.Cooldowns[i] = 100 + i
					}
					wantCooldowns := ad.Cooldowns
					input := binary.LittleEndian.AppendUint16(nil, version)
					input = append(input, present)
					if present != 0 {
						input = append(input, 0)
					} // Empty enchantment list.
					tail := class == 0 && (version == 4 || int16(version) >= 5 && present != 0)
					if tail {
						input = append(input, 0, 0)
						input = binary.LittleEndian.AppendUint32(input, 0x12345678)
						input = append(input, make([]byte, 20)...)
					}
					wantRet, wantPos := uint32(1), int64(len(input))
					if int16(version) > 5 {
						wantRet = 0
						wantPos = 2
					} else if present != 0 && gf&2048 == 0 {
						wantRet = 0
						wantPos = 3
					} else if tail {
						for i := 1; i < 6; i++ {
							wantCooldowns[i] = 0
						}
					}
					input = append(input, 0xde, 0xad, 0xbe, 0xef)
					before := bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(u)), int(unsafe.Sizeof(*u))))
					ret, got, pos := playerFileSection(t, "nox_xxx_guiEnchantment_41B9C0", input, uint32(uintptr(unsafe.Pointer(u))), 0)
					if ret != wantRet || pos != wantPos || !bytes.Equal(got, input) || ad.Cooldowns != wantCooldowns || !bytes.Equal(before, unsafe.Slice((*byte)(unsafe.Pointer(u)), len(before))) {
						t.Fatal("enchant gates", version, class, gf, present, ret, pos, wantPos, ad.Cooldowns, wantCooldowns)
					}
					rows = append(rows, map[string]any{"version": version, "class": class, "flags": uint32(gf), "present": present, "return": ret, "position": pos, "cooldowns": ad.Cooldowns})
				}
			}
		}
	}
	spellbookCapture(t, "player-files-enchantment-read-gates", rows, "64886e8b29df9b824e3aec9d986b10e54f198a6df21ae3004ec5fe900a498b82")
}
