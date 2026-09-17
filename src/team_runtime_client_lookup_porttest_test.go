//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestTeamRuntimeClientObjectLookup(t *testing.T) {
	o := newObjectDrawingOwner(t)
	defer noxflags.PortTestGameFlags(0)()
	oldClient, oldServer := noxClient, noxServer
	noxClient, noxServer = o.c.Client, o.c.srv
	t.Cleanup(func() { noxClient, noxServer = oldClient, oldServer })
	dynamic := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(10, 10))
	static := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(20, 20))
	static.ObjClass = object.Class(0x20000000)
	dynamic.ObjClass = object.Class(4)
	type row struct {
		Code           uint32
		Found, Missing bool
	}
	var rows []row
	for _, code := range []uint32{0, 1, 31, 0x7fff, 0x8000, 0xffff, 0x80000001, 0xffffffff} {
		t.Run(fmt.Sprintf("%x", code), func(t *testing.T) {
			dynamic.NetCode32 = code
			static.NetCode32 = code
			got := legacy.PortTestTeamRuntimeObject(int(code))
			if got != dynamic.TeamPtr() {
				t.Fatal("dynamic team lookup selected wrong drawable")
			}
			if nox_xxx_objGetTeamByNetCode_418C80(int(code)) != got {
				t.Fatal("existing Go object-team lookup differs")
			}
			dynamic.NetCode32 = code ^ 0x10000
			missing := legacy.PortTestTeamRuntimeObject(int(code))
			if missing != nil {
				t.Fatal("static drawable included in dynamic lookup")
			}
			rows = append(rows, row{code, got == dynamic.TeamPtr(), missing == nil})
		})
	}
	spellbookCapture(t, "team-runtime-client-object-lookup", rows, "2e3367a483507589057e9c447e4834b2242d725426c87c5488497c339a015ff8")
}
