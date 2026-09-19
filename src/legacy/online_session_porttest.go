//go:build porttest

package legacy

import "unsafe"

func PortTestOnlineSessionWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"failed":   (*uint32)(unsafe.Pointer(&onlineConnectionFailed)),
		"status":   (*uint32)(unsafe.Pointer(&onlineSessionStatus)),
		"pending":  (*uint32)(unsafe.Pointer(&onlineRetryPending)),
		"active":   (*uint32)(unsafe.Pointer(&onlineRetryActive)),
		"deadline": (*uint32)(unsafe.Pointer(&onlineRetryDeadline)),
		"origin":   (*uint32)(unsafe.Pointer(&onlineRetryOrigin)),
	}
	old := make(map[string]uint32)
	for k, p := range words {
		old[k] = *p
	}
	return words, func() {
		for k, p := range words {
			*p = old[k]
		}
	}
}
func PortTestOnlineSessionCall(op string, a, b uint32) uint32 {
	switch op {
	case "reset":
		return onlineSessionReset()
	case "start":
		onlineSessionStart()
		return 0
	case "attempt":
		return onlineSessionAttempt()
	case "retry":
		return onlineSessionRetry()
	case "marker":
		return onlineSessionMarker(a)
	case "briefing":
		return onlineSessionBriefing(a)
	case "map":
		return onlineSessionMap()
	case "channels", "users":
		return onlineSessionListCleanup()
	case "append":
		return onlineSessionEvent()
	case "clear", "player", "players":
		return 0
	case "status":
		return onlineSessionStatus
	case "login":
		return onlineSessionLogin()
	}
	panic(op)
}
