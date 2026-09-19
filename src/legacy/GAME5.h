#ifndef NOX_PORT_GAME5
#define NOX_PORT_GAME5

#include "defs.h"

int sub_545E60(nox_object_t* a1);
void sub_548600(nox_object_t* a1, float a2, float a3);
void sub_548830(int a1);
void sub_548860(int a1, short a2);
int nox_xxx_strikeOgre_549220(float a1);
int nox_xxx_strikeMonsterDefault_549380(float a1);
int nox_xxx_strikeScorpion_5495B0(float a1);
int nox_xxx_strikeVileZombie_549700(float a1);
int nox_xxx_strikeStoneGolem_5497E0(float a1);
int nox_xxx_strikeMechGolem_549960(float a1);
int nox_xxx_strikeWasp_549980(float a1);
int nox_xxx_strikeGhost_549A60(float a1);
int nox_xxx_strikeBomber_549BB0();
int nox_xxx_strikeSpider_549BC0(float a1);
int nox_xxx_strikeSpittingSpider_549CA0(float a1);
int sub_549D80(int a1);
int sub_549E00(int a1);
int sub_549E70(int a1);
int sub_549E90(int a1);
int sub_549FA0(int a1);
int nox_bomberDead_54A150(nox_object_t* a1);
int sub_54A250(int a1);
int nox_xxx_monsterDeadTroll_54A270(int a1);
int sub_54A310(int a1);
int sub_54A750(int a1);
int sub_54A7D0(int a1);
int sub_54A850(int a1);
int sub_54A890(int a1);
int sub_54A900(int a1);
int sub_54A950(int a1);
short nox_xxx_monsterAutoSpells_54C0C0(nox_object_t* a1p);
void nox_xxx_monsterCreateFn_54C480(nox_object_t* a1);
int nox_xxx_createWeapon_54C710(int a1);
uintptr_t sub_54C950(int a1);
int nox_xxx_createFnObelisk_54CA10(int a1);
void nox_xxx_createFnAnim_54CA50(int a1);
uint8_t* nox_xxx_createTrigger_54CA60(int a1);
uint32_t* nox_xxx_createMonsterGen_54CA90(int a1);
uint32_t* nox_xxx_createRewardMarker_54CAC0(int a1);
int nox_xxx_dieImpEgg_54CAE0(int a1);
void nox_xxx_diePolyp_54CB10(int a1);
void nox_xxx_diePotion_54CBB0(int a1);
void sub_54CBD0(int a1);
int nox_xxx_diePlayer_54D2B0(int a1);
void nox_xxx_dieGlyph_54DF30(nox_object_t* a1);
void nox_xxx_dieBarrel_54DFA0(int a1);
void nox_xxx_dieCreateObject_54E010(int a1);
short nox_xxx_dieSpawnObject_54E070(int a1);
void nox_xxx_dieMarker_54E460(int a1);
void nox_xxx_dieBoulder_54E4B0(int a1);
int nox_xxx_dieGameBall_54E620(int a1);
void nox_xxx_dieMonsterGen_54E630(int a1);
char nox_xxx_updateMonsterGenerator_54E930(uint32_t* a1);
void nox_xxx_updateHarpoon_54F380(nox_object_t* a1);
void nox_xxx_unitUpdateMover_54F740(int a1);
int nox_xxx_updateShootingTrap_54F9A0(int a1);
void nox_xxx_collideTrigger_54FCD0(int a1, int a2);
float* nox_xxx_createSpark_54FD80(float a1, float a2, int a3, int a4, float a5, float a6, float a7, int a8);

#endif // NOX_PORT_GAME5
