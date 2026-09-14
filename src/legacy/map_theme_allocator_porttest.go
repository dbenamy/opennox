//go:build porttest

package legacy

/*
#cgo LDFLAGS: -Wl,--wrap=calloc -Wl,--wrap=free -Wl,--wrap=time
#include <stdint.h>
#include <stdlib.h>
#include <time.h>
void* __real_calloc(size_t, size_t);
void __real_free(void*);
time_t __real_time(time_t*);
extern void themeTestAllocated(void*, size_t);
extern void themeTestReleased(void*);
static int theme_observe_active;
static uint32_t theme_observe_time;
void themeTestObserve(int active, uint32_t epoch) {
 theme_observe_time = epoch;
 theme_observe_active = active;
}
void* __wrap_calloc(size_t n, size_t size) {
 void* p = __real_calloc(n, size);
 if (theme_observe_active && p) themeTestAllocated(p, n * size);
 return p;
}
void __wrap_free(void* p) {
 if (theme_observe_active && p) themeTestReleased(p);
 __real_free(p);
}
time_t __wrap_time(time_t* out) {
 if (!theme_observe_active) return __real_time(out);
 time_t value = (time_t)theme_observe_time;
 if (out) *out = value;
 return value;
}
*/
import "C"

func themeObserve(active bool, epoch uint32) {
	C.themeTestObserve(C.int(bool2int(active)), C.uint32_t(epoch))
}
