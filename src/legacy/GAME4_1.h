#ifndef NOX_PORT_GAME4_1
#define NOX_PORT_GAME4_1

#include "defs.h"

#define nox_xxx_monsterPushAction_50A260(obj, a2) nox_xxx_monsterPushAction_50A260_impl(obj, a2, __FILE__, __LINE__)
void* nox_xxx_monsterPushAction_50A260_impl(nox_object_t* a1p, int a2, char* file, int line);
void nox_xxx_monsterClearActionStack_50A3A0(nox_object_t* a1);
void sub_50E140(int a1);
int sub_51D2C0(int a1, int a2);
int sub_51D300(int a1, int a2, char a3);
int nox_xxx_tileGetDefByName_51D4D0(char* a1);
int nox_xxx_tileCheckImage_51D540(int a1);
int nox_xxx_tileCheckImageVari_51D570(int a1);
int nox_xxx_tile_51D5C0(int a1);

#endif // NOX_PORT_GAME4_1
