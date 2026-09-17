//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestServerOptionsTopTabs(t *testing.T) {
	type row struct {
		Host                                               bool
		ID                                                 uint
		MainHidden, HeaderHidden, General, Access, Players bool
		Loads                                              []string
		Children                                           []uint
	}
	var rows []row
	for _, host := range []bool{false, true} {
		t.Run(fmt.Sprint(host), func(t *testing.T) {
			catalog := legacy.PortTestMapCatalogOpen(21)
			t.Cleanup(catalog.Close)
			o := newServerOptionsOwner(t)
			flags := noxflags.GameFlag(0)
			if host {
				flags = 1
			}
			defer noxflags.PortTestGameFlags(flags)()
			loads := o.installSubpanels(t)
			*o.optionWords["tabs2"] = uint32(uintptr(o.images[1].C()))
			*o.optionWords["tabs3"] = uint32(uintptr(o.images[2].C()))
			for _, id := range []uint{10159, 10160, 10159, 10161, 10160, 10159, 10162, 10160} {
				if legacy.PortTestServerOptionsEvent(o.options, o.options.ChildByID(id), 16391, 0) != 1 {
					t.Fatal("tab dispatch")
				}
				r := row{host, id, o.options.ChildByID(10150).Flags.IsHidden(), o.options.ChildByID(10196).Flags.IsHidden(), *o.optionWords["advanced-open"] != 0, *o.optionWords["access"] != 0, *o.optionWords["players"] != 0, append([]string(nil), (*loads)...), serverOptionsChildren(t, o.options)}
				if id == 10160 {
					if r.MainHidden || !r.HeaderHidden || r.General || r.Access || r.Players {
						t.Fatalf("main tab %+v", r)
					}
				}
				if id == 10159 {
					if !r.MainHidden || r.HeaderHidden == host || r.General != host || r.Players == host {
						t.Fatalf("options tab %+v", r)
					}
				}
				rows = append(rows, r)
				o.c.GUI.FreeDestroyed()
			}
		})
	}
	spellbookCapture(t, "server-options-top-tabs", rows, "5584a5229b99893dcc64e5d615d0b0a03cb3b207da360810f849ac930e829eb5")
}
