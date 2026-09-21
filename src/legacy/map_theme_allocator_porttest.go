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
static _Thread_local int grid_active, grid_fail, grid_count, grid_freed, grid_valid;
static _Thread_local void* grid_ptr[130];
static _Thread_local size_t grid_size[130];
static _Thread_local int grid_live[130];
void worldGridAllocObserve(int fail) {
 grid_active = 1; grid_fail = fail; grid_count = 0; grid_freed = 0; grid_valid = 1;
}
void worldGridAllocStop(void) { grid_active = 0; grid_fail = 0; }
int worldGridAllocStat(int index) {
 if (index == -1) return grid_count;
 if (index == -2) return grid_freed;
 if (index == -3) return grid_valid;
 return grid_size[index];
}
int worldGridAllocContains(void* p) {
 for (int i = 0; i < grid_count; i++) if (grid_ptr[i] == p && grid_live[i]) return 1;
 return 0;
}
// Resource teardown records raw addresses without calling Go inside free.
static _Thread_local uintptr_t* resource_free_events;
static _Thread_local int resource_free_capacity, resource_free_count;
void resourceFreeObserve(uintptr_t* events, int capacity) {
 resource_free_events = events; resource_free_capacity = capacity; resource_free_count = 0;
}
int resourceFreeStop(void) {
 resource_free_events = NULL;
 return resource_free_count;
}
static uint32_t theme_observe_time;
void themeTestObserve(int active, uint32_t epoch) {
 theme_observe_time = epoch;
 theme_observe_active = active;
}
void* __wrap_calloc(size_t n, size_t size) {
 if (grid_active && grid_fail > 0 && n == 128 && (size == 4 || size == 44) && --grid_fail == 0) return NULL;
 void* p = __real_calloc(n, size);
 if (grid_active && p) {
  if (grid_count >= 130) grid_valid = 0;
  else { grid_ptr[grid_count] = p; grid_size[grid_count] = n * size; grid_live[grid_count++] = 1; }
 }
 if (theme_observe_active && p) themeTestAllocated(p, n * size);
 return p;
}
void __wrap_free(void* p) {
 if (resource_free_events && p) {
  if (resource_free_count < resource_free_capacity) resource_free_events[resource_free_count] = (uintptr_t)p;
  resource_free_count++;
 }
 if (grid_active && p) {
  int found = 0;
  for (int i = 0; i < grid_count; i++) if (grid_ptr[i] == p && grid_live[i]) {
   grid_live[i] = 0; grid_freed++; found = 1; break;
  }
  if (!found) grid_valid = 0;
 }
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
