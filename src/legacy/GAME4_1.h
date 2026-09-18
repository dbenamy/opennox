#ifndef NOX_PORT_GAME4_1
#define NOX_PORT_GAME4_1

#include "defs.h"

int nox_xxx_math_509EA0(int a1);
int nox_xxx_monsterPopAction_50A160(nox_object_t* a1p);
#define nox_xxx_monsterPushAction_50A260(obj, a2) nox_xxx_monsterPushAction_50A260_impl(obj, a2, __FILE__, __LINE__)
void* nox_xxx_monsterPushAction_50A260_impl(nox_object_t* a1p, int a2, char* file, int line);
void nox_xxx_monsterClearActionStack_50A3A0(nox_object_t* a1);
void nox_xxx_unitUpdateMonster_50A5C0(nox_object_t* a1);
int nox_xxx_mapTraceObstacles_50B580(nox_object_t* a1, float4* a3);
void sub_50E140(int a1);
int nox_xxx_shopGetItemCost_50E3D0(int a1, int a2, float a3);
void sub_50F3A0(uint32_t* a1);
void nox_xxx_shopExit_50F4C0(uint32_t* a1);
void nox_xxx_tradeAccept_50F5A0(int a1, int a2);
int nox_xxx_tradeP2PAddOfferMB_50FE20(int a1, int a2);
uint32_t* sub_5108D0(int a1, int a2, int a3);
uint32_t* sub_510AE0(int* a1, int a2, uint32_t* a3);
void sub_510D10(int* a1, int a2, int a3, unsigned int a4);
void nox_xxx_shopCancelSession_510DC0(void* a1);
int sub_510DE0(int a1, int a2);
void sub_510E20(int a1);
signed int nox_xxx_updateSentryGlobe_510E60(int a1);
void nox_xxx_updateSprings_5113A0();
int nox_xxx_unitSetDecayTime_511660(nox_object_t* a1, int a2);
int sub_515C80(int a1, uint8_t* a2);
void sub_516FC0();
int sub_517590(float a1, float a2);
void nox_xxx_moveUpdateSpecial_517970(nox_object_t* unit);
void nox_xxx_unitsGetInCircle_517F90(float2* a1, float a2, void* fnc, void* data);
void nox_xxx_getMissilesInCircle_518170(float2* a1, float a2, void* a3, nox_object_t* a4);
nox_waypoint_t* sub_518740(float2* a1, unsigned char a2);
int sub_51A500(int a1);
char* sub_51A550();
void nox_xxx_spawnHecubahQuest_51A5A0(int* a1);
void nox_xxx_spawnNecroQuest_51A7A0(int* a1);
int nox_xxx_getQuestStage_51A930();
int sub_51A940(int a1);
int sub_51A950();
void nox_xxx_playerResetControlBuffer_51AC30(int a1);
char sub_51B860(int a1);
void sub_51D0E0();
int sub_51D0F0(char a1);
int sub_51D100(int a1);
uint32_t* sub_51D120(float* a1);
float* sub_51D1A0(float2* a1);
int sub_51D2C0(int a1, int a2);
int sub_51D300(int a1, int a2, char a3);
float2* sub_51D3F0(float2* a1, float2* a2);
int nox_xxx_tileGetDefByName_51D4D0(char* a1);
int nox_xxx_tileCheckImage_51D540(int a1);
int nox_xxx_tileCheckImageVari_51D570(int a1);
int nox_xxx_tile_51D5C0(int a1);
int sub_51D8F0(float2* a1);
void sub_517B70(float2* a1, void* a2, void* a3);
void nox_xxx_getUnitsInRect_517C10(float4* a1, void* fnc, void* a3);
void nox_xxx_getUnitsInRectAdv_517ED0(float4* a1, void* a2, void* a3);

#endif // NOX_PORT_GAME4_1
