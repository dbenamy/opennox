//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "defs.h"
#include "GAME4.h"
#include "GAME4_1.h"
static float prefabFloat(uint32_t v) { float x; memcpy(&x,&v,4); return x; }
static uint64_t prefabDouble(double v) { uint64_t x; memcpy(&x,&v,8); return x; }
static uint32_t prefab_calls[512][2];
static unsigned int prefab_call_count;
static void prefabRecord(int obj, int data) {
 if (prefab_call_count >= 512) abort();
 prefab_calls[prefab_call_count][0] = (uint32_t)obj;
 prefab_calls[prefab_call_count++][1] = (uint32_t)data;
}
static void* prefabObserver(void) { prefab_call_count = 0; return (void*)prefabRecord; }
static uint32_t* prefabCalls(void) { return &prefab_calls[0][0]; }
static unsigned int prefabCallCount(void) { return prefab_call_count; }
static uint64_t prefabInvoke(int op, const uint32_t* v) { switch(op) {
case 0: nox_server_scriptExecuteFnForEachGroupObj_502670((unsigned char*)(uintptr_t)v[0], (int)v[1], (void (*)(int,int))(uintptr_t)v[2], (int)v[3]); return 0;
case 1: return (uint32_t)(nox_xxx_mapgenMakeScript_502790((FILE*)(uintptr_t)v[0], (char*)(uintptr_t)v[1]));
case 2: return (uint32_t)(sub_5029A0((char*)(uintptr_t)v[0]));
case 3: return (uint32_t)(sub_5029F0((int)v[0]));
case 4: return (uint32_t)(sub_502A20());
case 5: return (uint32_t)(sub_502A50((char*)(uintptr_t)v[0]));
case 6: return (uint32_t)(sub_502AB0((char*)(uintptr_t)v[0]));
case 7: return (uint32_t)(sub_502B10());
case 8: return (uint32_t)(sub_502D70((int)v[0]));
case 9: return (uint32_t)(uintptr_t)(sub_502DA0((char*)(uintptr_t)v[0]));
case 10: return (uint32_t)(uintptr_t)(sub_502DF0());
case 11: return (uint32_t)(uintptr_t)(sub_502E10((int)v[0]));
case 12: return prefabDouble(sub_502E70((int)v[0]));
case 13: return prefabDouble(sub_502EA0((int)v[0]));
case 14: return (uint32_t)(nox_xxx_mapgenSaveMap_503830((int)v[0]));
case 15: return (uint32_t)(sub_503B30((float2*)(uintptr_t)v[0]));
case 16: return (uint32_t)(sub_503EC0((int)v[0], (float*)(uintptr_t)v[1]));
case 17: return (uint32_t)(uintptr_t)(nox_xxx_tileAllocTileInCoordList_5040A0((int)v[0], (int)v[1], prefabFloat(v[2])));
case 18: return (uint32_t)(nox_xxx_tileInit_504150((int)v[0], (int)v[1]));
case 19: return (uint32_t)(uintptr_t)(sub_504290((char)v[0], (char)v[1]));
case 20: return (uint32_t)(uintptr_t)(nox_xxx_cliWallGet_5042F0((int)v[0], (int)v[1]));
case 21: return (uint32_t)(sub_504330((int)v[0], (int)v[1]));
case 22: return (uint32_t)(uintptr_t)(sub_5044B0((int)v[0], prefabFloat(v[1]), prefabFloat(v[2])));
case 23: return (uint32_t)(sub_504560((int)v[0], (int)v[1]));
case 24: return (uint32_t)(uintptr_t)(nox_xxx_unitAddToList_5048A0((int)v[0]));
case 25: return (uint32_t)(sub_504910((int)v[0], (int)v[1]));
case 26: return (uint32_t)(sub_504980());
case 27: return (uint32_t)(sub_5049C0((int)v[0]));
case 28: return (uint32_t)(uintptr_t)(sub_5049D0());
case 29: return (uint32_t)(sub_5049E0((int)v[0]));
case 30: return (uint32_t)(sub_504A10((int)v[0]));
case 31: return (uint32_t)(uintptr_t)(sub_505060());
case 32: return (uint32_t)(nox_server_mapRWMapIntro_505080());
case 33: return (uint32_t)(nox_server_mapRWGroupData_505C30());
case 34: return (uint32_t)(nox_server_mapRWWaypoints_506260((uint32_t*)(uintptr_t)v[0]));
case 35: sub_51D0E0(); return 0;
case 36: return (uint32_t)(sub_51D0F0((char)v[0]));
case 37: return (uint32_t)(uintptr_t)(sub_51D120((float*)(uintptr_t)v[0]));
case 38: return (uint32_t)(uintptr_t)(sub_51D1A0((float2*)(uintptr_t)v[0]));
case 39: return (uint32_t)(uintptr_t)(sub_51D3F0((float2*)(uintptr_t)v[0], (float2*)(uintptr_t)v[1]));
default: abort(); } }
*/
import "C"
import "unsafe"

func PortTestPrefabCall(op int, args [6]uint32) uint64 {
	return uint64(C.prefabInvoke(C.int(op), (*C.uint32_t)(unsafe.Pointer(&args[0]))))
}

func PortTestPrefabObserver() unsafe.Pointer { return C.prefabObserver() }
func PortTestPrefabCalls() [][2]uint32 {
	n := int(C.prefabCallCount())
	return append([][2]uint32(nil), unsafe.Slice((*[2]uint32)(unsafe.Pointer(C.prefabCalls())), n)...)
}
