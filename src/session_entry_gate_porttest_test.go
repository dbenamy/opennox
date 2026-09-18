//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestSessionEntrySaveGate(t *testing.T) {
	o := newMatchRosterOwner(t)
	frame := serverConfigOwnBytes(t, 0x5D4594, 1563068, 4)
	cursor := serverConfigOwnBytes(t, 0x5D4594, 1096672, 4)
	book, restoreBook := legacy.PortTestBookWords()
	defer restoreBook()
	oldLoading, oldPending, oldDelay := dword_5d4594_1563080, dword_5d4594_1563092, dword_5d4594_1563088
	oldBusy := legacy.Get_dword_5d4594_251744()
	defer func() {
		dword_5d4594_1563080 = oldLoading
		dword_5d4594_1563092 = oldPending
		dword_5d4594_1563088 = oldDelay
		legacy.Set_dword_5d4594_251744(oldBusy)
	}()
	pl := o.s.Players.ByInd(31)
	u := pl.PlayerUnit
	type row struct {
		Flags, Now, Last uint32
		Block            string
		Result           int
	}
	var rows []row
	for _, flags := range []uint32{0, 2048, 2048 | 8192} {
		for _, times := range [][2]uint32{{29, 0}, {30, 0}, {31, 0}, {0, 0xffffffe2}, {1, 0xffffffe2}, {0, 1}, {0xffffffff, 0}} {
			for _, block := range []string{"none", "inactive", "unitless", "loading", "pending", "busy", "cursor", "object-4000", "object-dead", "object-no-update", "trading", "book", "book-other"} {
				deferFlags := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
				o.reset()
				o.s.SetFrame(times[0])
				binary.LittleEndian.PutUint32(frame, times[1])
				clear(cursor)
				pl.Active = 1
				pl.PlayerUnit = u
				u.ObjFlags = 0
				objectXferSetWord(u.UpdateData, 284, 0)
				dword_5d4594_1563080 = false
				dword_5d4594_1563092 = 0
				dword_5d4594_1563088 = 0
				legacy.Set_dword_5d4594_251744(0)
				*book["dword_5d4594_1047520"] = 0
				switch block {
				case "inactive":
					pl.Active = 0
				case "unitless":
					pl.PlayerUnit = nil
				case "loading":
					dword_5d4594_1563080 = true
				case "pending":
					dword_5d4594_1563092 = 1
				case "busy":
					legacy.Set_dword_5d4594_251744(1)
				case "cursor":
					binary.LittleEndian.PutUint32(cursor, 2)
				case "object-4000":
					u.ObjFlags = 0x4000
				case "object-dead":
					u.ObjFlags = object.FlagDead
				case "object-no-update":
					u.ObjFlags = object.FlagNoUpdate
				case "trading":
					objectXferSetWord(u.UpdateData, 284, 1)
				case "book":
					*book["dword_5d4594_1047520"] = 1
				case "book-other":
					*book["dword_5d4594_1047520"] = 2
				}
				got := legacy.Nox_xxx_game_4DCCB0()
				want := 0
				if flags&2048 == 0 || times[0]-times[1] >= 30 && (block == "none" || block == "book-other") {
					want = 1
				}
				if got != want {
					t.Errorf("flags%x times%v block%s got%d want%d", flags, times, block, got, want)
				}
				rows = append(rows, row{flags, times[0], times[1], block, got})
				pl.Active = 1
				pl.PlayerUnit = u
				deferFlags()
			}
		}
	}
	spellbookCapture(t, "session-entry-save-gate", rows, "4f82efc7edc9b59d88eacec88918cdf3455576b6b269d4ea4067afcc017b208f")
}
