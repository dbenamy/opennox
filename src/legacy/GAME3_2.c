#include "GAME3_2.h"

// Shared catalog head; catalog algorithms are already native.
nox_list_item_t nox_common_maplist = {0};

// C varargs are formatted here; Go owns message serialization and fanout.
extern int nox_gameplayTextLine(int unit, wchar2_t* text);
extern int nox_gameplayTextAll(char flags, wchar2_t* text);
int nox_xxx_netSendLineMessage_4D9EB0(int unit, wchar2_t* format, ...) {
 if (!unit || !(*(uint8_t*)(unit + 8) & 4)) { return unit; }
 wchar2_t text[256];
 va_list args;
 va_start(args, format);
 nox_vswprintf(text, format, args);
 va_end(args);
 return nox_gameplayTextLine(unit, text);
}
int nox_xxx_printToAll_4D9FD0(char flags, wchar2_t* format, ...) {
 wchar2_t text[256];
 va_list args;
 va_start(args, format);
 nox_vswprintf(text, format, args);
 va_end(args);
 return nox_gameplayTextAll(flags, text);
}
