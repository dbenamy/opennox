//go:build porttest

package legacy

import "github.com/opennox/opennox/v1/server"
import "unsafe"

func PortTestServerOrchestration(op string, u *server.Object, arg int32) uint32 {
	switch op {
	case "flags":
		orchestrationRoundFlags()
	case "players":
		orchestrationTransitionPlayers()
	case "restore":
		orchestrationRestore(arg)
	case "load-state":
		return uint32(orchestrationLoadState())
	case "difficulty":
		orchestrationDifficulty()
	case "drop-flags":
		orchestrationDropFlags()
	case "reset-player":
		return orchestrationResetPlayer(u)
	case "rewards":
		orchestrationRewards()
	case "walls":
		orchestrationWalls()
	case "timeout":
		return uint32(bool2int(orchestrationTimeout()))
	}
	return 0
}

func PortTestServerOrchestrationGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{"ankh-marker": (*uint32)(unsafe.Pointer(&dword_5d4594_1568280)), "selected-marker": (*uint32)(unsafe.Pointer(&dword_5d4594_1568288)), "drop-table": (*uint32)(unsafe.Pointer(&dword_5d4594_2488728)), "restore-cleanup": &orchestrationRestoreCleanup, "reward-marker": &orchestrationRewardMarker}
	saved := map[string]uint32{}
	for k, p := range words {
		saved[k] = *p
		*p = 0
	}
	return words, func() {
		for k, p := range words {
			*p = saved[k]
		}
	}
}

func PortTestServerOrchestrationChestInit() unsafe.Pointer { return lifecycleInitKey(initIDChest) }
