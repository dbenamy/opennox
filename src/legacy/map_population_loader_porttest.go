//go:build porttest

package legacy

/*
#cgo LDFLAGS: -Wl,--wrap=sub_502D70
int __real_sub_502D70(int index);
*/
import "C"

// The batch owns selection/placement, while file decoding remains a service.
// A fixture may supply an already-decoded cache and record loader requests.
// Every other caller, including integration tests, reaches the real loader.
var populationTestLoad func(int32) (int32, bool)

//export __wrap_sub_502D70
func __wrap_sub_502D70(index C.int) C.int {
	if populationTestLoad != nil {
		if result, ok := populationTestLoad(int32(index)); ok {
			return C.int(result)
		}
	}
	return C.__real_sub_502D70(index)
}
