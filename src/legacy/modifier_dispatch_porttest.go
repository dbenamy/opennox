//go:build porttest

package legacy

/*
#include <stdint.h>
static uintptr_t portModifierArgs[6];
static int32_t portModifierResult;
static int portModifierCalls;
static int portModifierKind;
static void portModifierReset(int32_t value) {
    for (int i=0;i<6;i++) portModifierArgs[i]=0;
    portModifierResult=value;portModifierCalls=0;portModifierKind=0;
}
static void portModifierRecord3(void* a,void* b,void* c) {
    portModifierArgs[0]=(uintptr_t)a;portModifierArgs[1]=(uintptr_t)b;
    portModifierArgs[2]=(uintptr_t)c;portModifierCalls++;
}
static int portModifierInt3(void* a,void* b,void* c) {
    portModifierRecord3(a,b,c);portModifierKind=1;return portModifierResult;
}
static void portModifierVoid3(void* a,void* b,void* c) {
    portModifierRecord3(a,b,c);portModifierKind=2;
}
static void portModifierVoid5(void* a,void* b,void* c,void* d,void* e) {
    portModifierRecord3(a,b,c);portModifierKind=3;
    portModifierArgs[3]=(uintptr_t)d;portModifierArgs[4]=(uintptr_t)e;
    if(e) *(uint32_t*)e=~(uint32_t)portModifierResult;
}
static void portModifierVoid6(void* a,void* b,void* c,void* d,void* e,void* f) {
    portModifierRecord3(a,b,c);portModifierKind=4;
    portModifierArgs[3]=(uintptr_t)d;portModifierArgs[4]=(uintptr_t)e;portModifierArgs[5]=(uintptr_t)f;
    if(f) *(uint32_t*)f=~(uint32_t)portModifierResult;
}
static void* portModifierObserver(int kind) {
    switch(kind) {case 1:return portModifierInt3;case 2:return portModifierVoid3;
    case 3:return portModifierVoid5;case 4:return portModifierVoid6;default:return 0;}
}
static uintptr_t portModifierArgument(int i) {return portModifierArgs[i];}
static int portModifierCount(void) {return portModifierCalls;}
static int portModifierLastKind(void) {return portModifierKind;}
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

func PortTestModifierObserverReset(value int32) [4]unsafe.Pointer {
	C.portModifierReset(C.int32_t(value))
	return [4]unsafe.Pointer{C.portModifierObserver(1), C.portModifierObserver(2), C.portModifierObserver(3), C.portModifierObserver(4)}
}
func PortTestModifierObserverSnapshot() ([6]uintptr, int, int) {
	var args [6]uintptr
	for i := range args {
		args[i] = uintptr(C.portModifierArgument(C.int(i)))
	}
	return args, int(C.portModifierCount()), int(C.portModifierLastKind())
}
func PortTestModifierCall3Result(key unsafe.Pointer, m *server.ModifierEff, a, b *server.Object) int32 {
	return server.CallModifierEffect3Result(key, m, a, b)
}
func PortTestModifierCall3Discard(key unsafe.Pointer, m *server.ModifierEff, a, b *server.Object) {
	server.CallModifierEffect3Discard(key, m, a, b)
}
func PortTestModifierCall5(key unsafe.Pointer, m *server.ModifierEff, a, b, c *server.Object, data unsafe.Pointer) {
	server.CallModifierEffect5(key, m, a, b, c, data)
}
func PortTestModifierCall6(key unsafe.Pointer, m *server.ModifierEff, a, b, c, d *server.Object, data unsafe.Pointer) {
	server.CallModifierEffect6(key, m, a, b, c, d, data)
}
