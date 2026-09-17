//go:build porttest

package opennox

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestServerConfigSubflags(t *testing.T) {
	o := newReliableReportsOwner(t)
	words, restore := legacy.PortTestServerOptionsWords()
	t.Cleanup(restore)
	globals, restore := legacy.PortTestMatchRosterGlobals()
	t.Cleanup(restore)
	type row struct {
		Op                                 string
		Game, Before, Mask, After, Updated uint32
		Result                             uint64
		Reports                            [][4]uint32
	}
	var rows []row
	for _, game := range []uint32{0, 1, 32, 33, 1024, 4096} {
		restoreFlags := noxflags.PortTestGameFlags(noxflags.GameFlag(game))
		for _, before := range []uint32{0, 2, 0x2000, 0x2002, 0x80000000, 0xffffffff} {
			for _, mask := range []uint32{0, 2, 0x2000, 0x2002, 0x80000000, 0xffffffff} {
				for _, op := range []string{"flags-set", "flags-add", "flags-remove", "flags-toggle", "flags-query"} {
					*globals["server-subflags"] = before
					*words["settings-updated"] = 0
					for i := range o.units {
						copy(unsafe.Slice((*uint32)(unsafe.Add(o.units[i].UpdateData, 248)), 4), []uint32{11, 12, 13, 0})
					}
					result := legacy.PortTestServerConfigScalar(op, int(mask), 0)
					after, notify, reset := before, false, false
					switch op {
					case "flags-set":
						after = mask
						notify = before != mask
					case "flags-add":
						after = before | mask
						notify = after != before
						reset = notify && game&1 != 0 && mask&0x2000 != 0
					case "flags-remove":
						after = before &^ mask
						notify = after != before
					case "flags-toggle":
						after = before ^ mask
						notify = true
						reset = game&1 != 0 && mask&0x2000 != 0 && after&0x2000 != 0
					case "flags-query":
						want := before&mask != 0 || mask == 0x2000 && game&1056 != 0
						if result != uint64(bool2int(want)) {
							t.Fatal("query", game, before, mask, result)
						}
					}
					if *globals["server-subflags"] != after || (*words["settings-updated"] != 0) != notify {
						t.Fatal("subflag mutation", game, op, before, mask, *globals["server-subflags"], after)
					}
					if got := legacy.PortTestServerConfigScalar("flags-get", 0, 0); got != uint64(int64(int32(after))) {
						t.Fatal("subflag getter", got, after)
					}
					r := row{Op: op, Game: game, Before: before, Mask: mask, After: after, Updated: *words["settings-updated"], Result: result}
					for i := range o.units {
						v := *(*[4]uint32)(unsafe.Add(o.units[i].UpdateData, 248))
						want := [4]uint32{11, 12, 13, 0}
						if reset {
							want = [4]uint32{0, 0, o.s.Frame(), 0}
						}
						if v != want {
							t.Fatal("report reset dispatch", game, op, before, mask, i, v, want)
						}
						r.Reports = append(r.Reports, v)
					}
					rows = append(rows, r)
				}
			}
		}
		restoreFlags()
	}
	spellbookCapture(t, "server-config-subflags", rows, "80094290ef413d5aebc686bcd5d192aa10fc53819dab70ff0f460edbecb0b6c2")
}
