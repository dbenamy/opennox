//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// PortTestBookAwardShop owns a real attached trade session. The quest exit path
// caches it instead of releasing it through the general trade-session pool.
func PortTestBookAwardShop(a, b *server.Object) (func() [4]uint32, func()) {
	s, free := alloc.New(shopSession{})
	s.Active = 1
	s.Units = [2]*server.Object{a, b}
	da, db := a.UpdateDataPlayer(), b.UpdateDataPlayer()
	oldA, oldB := da.Trade70, db.Trade70
	ind := int(da.Player.PlayerInd)
	oldCache := shopCached()[ind]
	shopCached()[ind] = nil
	da.Trade70 = s
	db.Trade70 = s
	return func() [4]uint32 {
		return [4]uint32{s.Active, uint32(bool2int(da.Trade70 == s)), uint32(bool2int(db.Trade70 == s)), uint32(bool2int(shopCached()[ind] == s))}
	}, func() { da.Trade70 = oldA; db.Trade70 = oldB; shopCached()[ind] = oldCache; free() }
}
