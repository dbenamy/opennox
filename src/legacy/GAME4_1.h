#ifndef NOX_PORT_GAME4_1
#define NOX_PORT_GAME4_1

#include "defs.h"

int nox_xxx_monsterPopAction_50A160(nox_object_t* a1p);
#define nox_xxx_monsterPushAction_50A260(obj, a2) nox_xxx_monsterPushAction_50A260_impl(obj, a2, __FILE__, __LINE__)
void* nox_xxx_monsterPushAction_50A260_impl(nox_object_t* a1p, int a2, char* file, int line);
void nox_xxx_monsterClearActionStack_50A3A0(nox_object_t* a1);
void nox_xxx_unitUpdateMonster_50A5C0(nox_object_t* a1);
void sub_50E140(int a1);
int nox_xxx_shopGetItemCost_50E3D0(int a1, int a2, float a3);
void nox_xxx_shopCancelSession_510DC0(void* a1);
void sub_510E20(int a1);
signed int nox_xxx_updateSentryGlobe_510E60(int a1);
void nox_xxx_updateSprings_5113A0();
int sub_517590(float a1, float a2);
void nox_xxx_moveUpdateSpecial_517970(nox_object_t* unit);
void sub_51D0E0();
int sub_51D2C0(int a1, int a2);
int sub_51D300(int a1, int a2, char a3);
int nox_xxx_tileGetDefByName_51D4D0(char* a1);
int nox_xxx_tileCheckImage_51D540(int a1);
int nox_xxx_tileCheckImageVari_51D570(int a1);
int nox_xxx_tile_51D5C0(int a1);
int sub_51D8F0(float2* a1);

#endif // NOX_PORT_GAME4_1
