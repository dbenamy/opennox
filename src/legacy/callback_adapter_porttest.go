//go:build porttest

package legacy

/*
#include "defs.h"
#include <stdint.h>

static uintptr_t adapter_draw_viewport;
static uintptr_t adapter_draw_drawable;
static int adapter_draw_return;
static int adapter_draw_id;
static int adapter_draw_calls;

static int adapter_draw_a(nox_draw_viewport_t *vp, nox_drawable *dr) {
    adapter_draw_viewport = (uintptr_t)vp;
    adapter_draw_drawable = (uintptr_t)dr;
    adapter_draw_id = 1;
    adapter_draw_calls++;
    return adapter_draw_return;
}
static int adapter_draw_b(nox_draw_viewport_t *vp, nox_drawable *dr) {
    adapter_draw_viewport = (uintptr_t)vp;
    adapter_draw_drawable = (uintptr_t)dr;
    adapter_draw_id = 2;
    adapter_draw_calls++;
    return adapter_draw_return;
}
static void *adapter_draw_callback(int which) {
    return which == 1 ? (void *)adapter_draw_a : (void *)adapter_draw_b;
}
static void adapter_draw_reset(int result) {
    adapter_draw_viewport = 0;
    adapter_draw_drawable = 0;
    adapter_draw_return = result;
    adapter_draw_id = 0;
    adapter_draw_calls = 0;
}
static uintptr_t adapter_draw_value(int field) {
    switch (field) {
    case 0: return adapter_draw_viewport;
    case 1: return adapter_draw_drawable;
    case 2: return (uintptr_t)adapter_draw_id;
    case 3: return (uintptr_t)adapter_draw_calls;
    default: return 0;
    }
}

*/
import "C"

import "unsafe"

func PortTestAdapterDrawCallback(which int) unsafe.Pointer {
	return C.adapter_draw_callback(C.int(which))
}
func PortTestAdapterDrawReset(result int) { C.adapter_draw_reset(C.int(result)) }
func PortTestAdapterDrawValue(field int) uintptr {
	return uintptr(C.adapter_draw_value(C.int(field)))
}
