//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestCombatOverlayFriends(t *testing.T) {
	type record struct {
		Name                   string
		Added, Removed, Reused []uint32
	}
	var rows []record
	for _, remove := range []uint32{0, 1, 2, 3, 65535, 0xffffffff} {
		for _, duplicate := range []bool{false, true} {
			name := fmt.Sprintf("remove=%x/duplicate=%t", remove, duplicate)
			t.Run(name, func(t *testing.T) {
				o := newCombatOverlayOwner(t)
				codes := []uint32{0, 1, 2, 3, 65535, 0xffffffff}
				if duplicate {
					codes = append(codes, 2)
				}
				var want []uint32
				for _, code := range codes {
					if legacy.PortTestCombatFriendAdd(code) == nil || !legacy.PortTestCombatFriendHas(code) {
						t.Fatal("add/query")
					}
					want = append([]uint32{code}, want...)
				}
				r := record{Name: name, Added: o.friendCodes(t)}
				if !reflect.DeepEqual(r.Added, want) {
					t.Fatal("friend insertion order/duplicates")
				}
				legacy.PortTestCombatFriendRemove(remove)
				for i, code := range want {
					if code == remove {
						want = append(want[:i], want[i+1:]...)
						break
					}
				}
				r.Removed = o.friendCodes(t)
				if !reflect.DeepEqual(r.Removed, want) {
					t.Fatal("remove first matching friend")
				}
				present := false
				for _, code := range want {
					present = present || code == remove
				}
				if legacy.PortTestCombatFriendHas(remove) != present {
					t.Fatal("membership after removal")
				}
				legacy.PortTestCombatFriendClear()
				if *o.words["friends"] != 0 {
					t.Fatal("clear head")
				}
				if legacy.PortTestCombatFriendAdd(1234) == nil {
					t.Fatal("reuse cleared pool")
				}
				r.Reused = o.friendCodes(t)
				if !reflect.DeepEqual(r.Reused, []uint32{1234}) {
					t.Fatal("reuse list")
				}
				legacy.PortTestCombatFriendDestroy()
				if *o.words["friends"] != 0 || *o.pools["friend"] != nil {
					t.Fatal("destroy globals")
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "combat-overlay-friends", rows, "ff06a0e526210dc4dfd16c3778c42d2899b497719f443822e079c277892befb3")
}
func TestCombatOverlayFriendCapacity(t *testing.T) {
	o := newCombatOverlayOwner(t)
	for i := 0; i < 128; i++ {
		if legacy.PortTestCombatFriendAdd(uint32(i)) == nil {
			t.Fatal("early allocation failure")
		}
	}
	before := o.friendCodes(t)
	if legacy.PortTestCombatFriendAdd(128) != nil || !reflect.DeepEqual(before, o.friendCodes(t)) {
		t.Fatal("friend capacity")
	}
	legacy.PortTestCombatFriendRemove(64)
	if legacy.PortTestCombatFriendAdd(128) == nil {
		t.Fatal("freed slot not reused")
	}
	after := o.friendCodes(t)
	if len(after) != 128 || after[0] != 128 || legacy.PortTestCombatFriendHas(64) {
		t.Fatal("reuse identity/order")
	}
	spellbookCapture(t, "combat-overlay-friend-capacity", [][]uint32{before, after}, "e7c16f023a28f10c9788ba3bcdd21196a7396e09f302e591ef4da038003523b5")
}
