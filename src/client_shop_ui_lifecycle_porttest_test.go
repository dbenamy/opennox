//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
	"unsafe"
)

func TestClientShopUIStartClose(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, render := range []bool{false, true} {
		for _, name := range []string{"", "Merchant", strings.Repeat("界", 24)} {
			for _, greeting := range []string{"", "FixtureGreeting", "FixtureSilentGreeting"} {
				o.reset(t)
				o.construct(t)
				o.constructShop(t)
				nox_client_renderGUI_80828 = render
				var title uintptr
				if name != "" {
					title = txptr(unsafe.Pointer(alloc.InternCString16(name)))
				}
				key := txptr(unsafe.Pointer(alloc.InternCString(greeting)))
				rows = append(rows, o.shopCapture(t, len(rows), 37, title, key, uintptr(o.c.Things.TypeByID("Shopkeeper").Index())))
				if *o.windowWords["dword_5d4594_1098624"] != 1 || *o.windowWords["dword_5d4594_1098628"] != 1 || o.shopWindow().Flags.IsHidden() || nox_client_renderGUI_80828 {
					t.Fatal("shop did not enter open/intro state")
				}
				wantSound := ""
				if greeting == "FixtureGreeting" {
					wantSound = "shop-greeting.wav"
				}
				if legacy.Dialogs.FileToRead() != wantSound {
					t.Fatalf("dialogue queued %q want%q", legacy.Dialogs.FileToRead(), wantSound)
				}
				rows = append(rows, o.shopCapture(t, len(rows), 8))
				rows = append(rows, o.shopCapture(t, len(rows), 15))
				if *o.windowWords["dword_5d4594_1098624"] != 0 || *o.windowWords["dword_5d4594_1098628"] != 0 || !o.shopWindow().Flags.IsHidden() || nox_client_renderGUI_80828 != render || legacy.Dialogs.FileToRead() != "" {
					t.Fatal("shop close did not restore UI/clear dialogue")
				}
				rows = append(rows, o.shopCapture(t, len(rows), 15))
			}
		}
	}
	shopUICapture(t, "start-close", rows, "e8db70fd55728a8feb4f874adcd98afa3697fd9cf67d3837a734f9c0d2e9e9f3")
}
func TestClientShopUIShopkeeperPictures(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	names := []string{"Shopkeeper", "ShopkeeperWarriorsRealm", "ShopkeeperConjurerRealm", "ShopkeeperWizardRealm", "ShopkeeperLandOfTheDead", "ShopkeeperMagicShop", "ShopkeeperPurple", "ShopkeeperYellow", "RedApple"}
	images := []string{"ShopKeeperPic", "ShopKeeperWarriorPic", "ShopKeeperConjurerPic", "ShopKeeperWizardPic", "ShopKeeperLandOfTheDeadPic", "ShopKeeperMagicShopPic", "ShopKeeperPurplePic", "ShopKeeperBrownPic", "ShopKeeperPic"}
	for _, cached := range []bool{false, true} {
		for i, name := range names {
			o.reset(t)
			o.construct(t)
			o.constructShop(t)
			if cached {
				legacy.PortTestShopUI(14, uintptr(o.c.Things.TypeByID("Shopkeeper").Index()))
			}
			r := o.shopCapture(t, len(rows), 14, uintptr(o.c.Things.TypeByID(name).Index()))
			rows = append(rows, r)
			if o.loads[len(o.loads)-1] != images[i] {
				t.Fatalf("%s picture request %q want%q", name, o.loads[len(o.loads)-1], images[i])
			}
			if r.Return == 0 || r.Return != o.norm(uint32(txptr(unsafe.Pointer(o.shopWindow().ChildByID(3805).DrawData().BgImageHnd)))) {
				t.Fatal("shopkeeper picture not assigned to widget")
			}
		}
	}
	shopUICapture(t, "shopkeeper-pictures", rows, "2cd3e566fe0e5d20a306475425815ec8bfaf703b7f12def8222cced73794a9db")
}
func TestClientShopUIDestroyRecreate(t *testing.T) {
	o := newShopUIOwner(t)
	o.reset(t)
	o.construct(t)
	for cycle := 0; cycle < 100; cycle++ {
		o.constructShop(t)
		c := &legacy.PortTestShopUICells()[0]
		c.Drawable = o.item(t, "RedApple", 123)
		c.Count = 1
		dr := c.Drawable
		o.releaseShop()
		o.c.GUI.FreeDestroyed()
		if o.objects[o.identities[dr]].Live || o.shopWindow() != nil || *o.windowWords["dword_5d4594_1098624"] != 0 {
			t.Fatalf("cycle%d retained shop ownership", cycle)
		}
	}
}
