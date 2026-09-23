//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestDeathRegistryRawForwarding(t *testing.T) {
	u, freeU := alloc.New(server.Object{})
	defer freeU()
	target, freeTarget := alloc.New(server.Object{})
	defer freeTarget()
	for _, tc := range []struct {
		name                                                 string
		class                                                object.Class
		nilTarget, selfTarget, sameTeam, nilCallback, called bool
	}{
		{name: "player", class: 4, called: true},
		{name: "monster", class: 2, called: true},
		{name: "both", class: 6, called: true},
		{name: "nonunit", class: 8},
		{name: "nil target", nilTarget: true},
		{name: "self", class: 4, selfTarget: true},
		{name: "same team", class: 4, sameTeam: true},
		{name: "nil slot ineligible target", class: 8, nilCallback: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			*u = server.Object{ObjFlags: 0x40}
			*target = server.Object{ObjClass: tc.class}
			u.Death = legacy.PortTestDeathForwardReset()
			if tc.nilCallback {
				u.Death = nil
			}
			if tc.sameTeam {
				u.TeamVal.ID = 1
				target.TeamVal.ID = 1
			}
			arg := target
			if tc.nilTarget {
				arg = nil
			}
			if tc.selfTarget {
				u.ObjClass = tc.class
				arg = u
			}
			legacy.PortTestDeathProjectile(u, arg)
			ptr, count := legacy.PortTestDeathForwardSnapshot()
			if tc.called {
				if count != 1 || ptr != u.CObj() || u.ObjFlags != 0x8040 {
					t.Fatalf("count/pointer/flags %d/%p/%x", count, ptr, u.ObjFlags)
				}
			} else if count != 0 || ptr != nil || u.ObjFlags != 0x40 {
				t.Fatalf("unexpected dispatch/mutation %d/%p/%x", count, ptr, u.ObjFlags)
			}
		})
	}
}
