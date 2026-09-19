//go:build porttest

package legacy

/*
#cgo CFLAGS: -DNOX_PORT_TEST_CLIENT_SOUND
*/
import "C"

var portTestClientSoundObserver func(id, volume int)

// PortTestClientSoundObserver observes requests without replacing the actual
// client sound function. Its ordinary lookup/allocation/playback path still runs.
func PortTestClientSoundObserver(fn func(id, volume int)) func() {
	old := portTestClientSoundObserver
	portTestClientSoundObserver = fn
	return func() { portTestClientSoundObserver = old }
}

//export nox_porttest_client_sound
func nox_porttest_client_sound(id, volume C.int) {
	if fn := portTestClientSoundObserver; fn != nil {
		fn(int(id), int(volume))
	}
}

func init() {
	audioEventPlayObserver = func(id, volume int) {
		if f := portTestClientSoundObserver; f != nil {
			f(id, volume)
		}
	}
}
