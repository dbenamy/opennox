//go:build porttest

package legacy

/*
#include "noxstring.h"
#include <stdarg.h>
static int portTestTextV(wchar2_t* dst, unsigned count, const wchar2_t* fmt, ...) {
 va_list args;
 va_start(args,fmt);
 int n=nox_vsnwprintf(dst,count,fmt,args);
 va_end(args);
 return n;
}
static int portTestTextFormat(int kind,wchar2_t* dst,unsigned count,const wchar2_t* fmt,
 unsigned a,unsigned b,double f,const wchar2_t* w1,const wchar2_t* w2,const char* s) {
 switch(kind) {
 case 0:return portTestTextV(dst,count,fmt);
 case 1:return portTestTextV(dst,count,fmt,a,b);
 case 2:return portTestTextV(dst,count,fmt,f);
 case 3:return portTestTextV(dst,count,fmt,w1,w2);
 case 4:return portTestTextV(dst,count,fmt,s);
 case 5:return portTestTextV(dst,count,fmt,w1,s,a,f);
 default:return -1;
 }
}
*/
import "C"

import (
	"math"
	"unsafe"
)

// All buffers are caller-owned, including explicit UTF-16 terminators. This
// adapter only supplies typed C varargs; the formatter remains production C.
func PortTestTextFormat(kind int, dst *uint16, count uint32, format *uint16, a, b uint32, bits uint64, w1, w2 *uint16, narrow *byte) int {
	return int(C.portTestTextFormat(C.int(kind), (*C.wchar2_t)(unsafe.Pointer(dst)), C.uint(count), (*C.wchar2_t)(unsafe.Pointer(format)), C.uint(a), C.uint(b), C.double(math.Float64frombits(bits)), (*C.wchar2_t)(unsafe.Pointer(w1)), (*C.wchar2_t)(unsafe.Pointer(w2)), (*C.char)(unsafe.Pointer(narrow))))
}
