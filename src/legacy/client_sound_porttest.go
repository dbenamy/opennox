//go:build porttest

package legacy

var portTestClientSoundObserver func(id, volume int)

// PortTestClientSoundObserver observes requests without replacing the actual
// client sound function. Its ordinary lookup/allocation/playback path still runs.
func PortTestClientSoundObserver(fn func(id, volume int)) func() {
	old := portTestClientSoundObserver
	portTestClientSoundObserver = fn
	return func() { portTestClientSoundObserver = old }
}

func init() {
	audioEventPlayObserver = func(id, volume int) {
		if f := portTestClientSoundObserver; f != nil {
			f(id, volume)
		}
	}
}
