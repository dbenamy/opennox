#ifndef NOX_PORT_GAME4_3
#define NOX_PORT_GAME4_3

#include "defs.h"

void sub_532930(int a1, unsigned short a2, unsigned short a3);
int nox_xxx_mobActionToAnimation_533790(int a1);
void nox_xxx_mobCalcDir_533CC0(int a1, float* a2);
int nox_xxx_monsterHasShield_5342C0(int a1);
int nox_xxx_monsterCanCast_534300(nox_object_t* a1);
int nox_xxx_monsterIsMoveing_534320(int a1);
int sub_534340(int a1);
int nox_xxx_monsterCanAttackAtWill_534390(nox_object_t* a1);
int sub_5343C0(int a1);
int sub_534440(int a1);
int sub_5347C0(int a1);
int nox_xxx_mobGetMoveAttemptTime_534810(nox_object_t* a1);
void nox_xxx_monsterMimicCheckMorph_534950(nox_object_t* a1);
int nox_xxx_wallPreDestroy_534DA0(int* a1);
int nox_xxx_playerPreAttackEffects_538290(int a1, int a2, int a3, int a4);
int nox_xxx_playerTraceAttack_538330(int a1, int a2);
void sub_538510(int a1, int a2);
void sub_5386A0(int a3, int a2);
int nox_xxx_itemApplyAttackEffect_538840(int a1, int a2, int a3);
int nox_xxx_playerAttack_538960(nox_object_t* a1);
short nox_xxx_warcryStunMonsters_539B90(int a1, int a2);
int nox_xxx_shootBowCrossbow1_539BD0(int a1, int a2);
uint32_t* nox_xxx_shootBowCrossbow2_539D80(int a1, int a2, int a3, char* a4);
int nox_xxx_shootApplyEffects_539F40(int a1, int a2, int a3);
int sub_539FB0(uint32_t* a1);
int nox_xxx_playerTryReloadQuiver_539FF0(uint32_t* a1);
void nox_xxx_fnElevatorShaft_53B410(int a1, int a2);
void nox_xxx_elevatorAud_53B490(int a1, int a2);
void nox_xxx_elevatorFn_53B750(int a1, int a2);
void sub_53BD10(int a1, int a2);
void nox_xxx_fnPentagramTeleport_53C060(float* a1, int a2);
void sub_53C140(float* a1, int a2);
void sub_53C240(float* a1, int arg4);
int nox_xxx_rechargeItem_53C520(int a1, int a2);
int nox_xxx_getRechargeRate_53C940(uint32_t* a1);
void nox_xxx_waterBarrel_53CC30(float* a1, int a2);
void sub_53D170(int a1, int a2);
void nox_xxx_updateFlameCleanse_53D510(int a1);
void sub_53D8C0(int a1, int a2);
void nox_xxx_toxicCloudPoison_53D9D0(int a1, int a2);
uint32_t* nox_xxx_wandShot_53F480(int a1, int a2, int* a3, uint32_t* a4);
int sub_543E60(int a1, int a2);
int nox_xxx_mapGenEdge_543EB0(int a1, int a2);
int sub_544020(char* a1);
int nox_xxx_tileCheckByte3_544070(int a1);
int nox_xxx_tileCheckByte4_5440A0(int a1);
int nox_xxx_mobSearchEdible_544A00(nox_object_t* a1, float a2);
int sub_544AE0(int a1, float a2);

#endif // NOX_PORT_GAME4_3
