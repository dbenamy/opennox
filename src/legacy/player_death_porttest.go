//go:build porttest

package legacy

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
)

func PortTestPlayerDeath(op string, victim, killer, assist *server.Object, info *server.Player) {
	switch op {
	case "death":
		playerDeath(victim)
	case "arena":
		playerDeathArena(victim, killer, assist, info)
	case "elimination":
		playerDeathElimination(victim, killer)
	case "kotr":
		playerDeathKotr(victim, killer)
	case "notify":
		playerDeathNotify(victim)
	default:
		panic(op)
	}
}
func PortTestPlayerCorpseCache()                               { playerCorpseCache() }
func PortTestPlayerCorpseCreate(pos types.Pointf, angle int32) { playerCorpseCreate(pos, angle) }

// Own the real lookup cache so assist resolution cannot retain fixture objects.
func PortTestPlayerDeathLookupOwner() func() {
	oldState, oldInit := netCodeCacheState, netCodeCacheNeedInit
	netCodeCacheState, netCodeCacheNeedInit = netCodeCacheStorage{}, 1
	return func() { netCodeCacheState, netCodeCacheNeedInit = oldState, oldInit }
}
