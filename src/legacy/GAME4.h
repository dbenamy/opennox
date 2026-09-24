#ifndef NOX_PORT_GAME4
#define NOX_PORT_GAME4

#include "defs.h"

int nox_xxx_XFerSpellReward_4F5F30(int* a1);
int nox_xxx_XFerAbilityReward_4F6240(int* a1);
int nox_xxx_XFerFieldGuide_4F6390(int* a1);
int nox_xxx_XFerWeapon_4F64A0(int a1);
int nox_xxx_XFerArmor_4F6860(int a1);
int nox_xxx_XFerAmmo_4F6B20(int* a1);
int nox_xxx_XFerTeam_4F6D20(int* a1);
int nox_xxx_XFerGold_4F6EC0(int a1);
int nox_xxx_XFerObelisk_4F6F60(int* a1);
int nox_xxx_XFerToxicCloud_4F70A0(int a1);
int nox_xxx_XFerMonsterGen_4F7130(int* a1);
int nox_xxx_XFerRewardMarker_4F74D0(int* a1);
void nox_xxx_mapFindPlayerStart_4F7AB0(float2* a1, nox_object_t* a2p);
int nox_xxx_weaponGetStaminaByType_4F7E80(int a1);
void nox_xxx_updatePlayer_4F8100(nox_object_t* a1);
void nox_xxx_animPlayerGetFrameRange_4F9F90(int a1, int* a2, int* a3);
int nox_xxx_unitGetStrength_4F9FD0(int a1);
int nox_xxx_playerSetState_4FA020(nox_object_t* a1, int a2);
int sub_4FA280(int a1);
int nox_common_mapPlrActionToStateId_4FA2B0(nox_object_t* a1);
int nox_xxx_checkInversionEffect_4FA4F0(int a1, int a2);
uint32_t* nox_xxx_playerAddGold_4FA590(int a1, int a2);
uint32_t* nox_xxx_playerSubGold_4FA5D0(int a1, unsigned int a2);
int nox_xxx_playerGetGold_4FA6B0(int a1);
int nox_object_getGold_4FA6D0(nox_object_t* a1);
char nox_xxx_mobMorphFromPlayer_4FAAC0(uint32_t* a1);
char nox_xxx_mobMorphToPlayer_4FAAF0(uint32_t* a1);
int nox_xxx_updatePlayerMonsterBot_4FAB20(uint32_t* a1);
int nox_xxx_netSendRewardNotify_4FAD50(int a1, int a2, int a3, char a4);
void sub_4FADD0(int a1, char* a2, char a3);
int sub_4FB050(int a1, int a2, int* a3);
void nox_xxx_playerExecuteAbil_4FBB70(nox_object_t* a1, int a2);
void sub_4FC0B0(nox_object_t* a1, int a2);
void nox_xxx_playerCancelAbils_4FC180(nox_object_t* a1);
int nox_common_playerIsAbilityActive_4FC250(nox_object_t* a1, int a2);
void sub_4FC300(nox_object_t* a1, int a2);
int nox_xxx_probablyWarcryCheck_4FC3E0(nox_object_t* a1, int a2);
void sub_4FC440(nox_object_t* a1, int a2);
void sub_4FC670(int a1);
unsigned short sub_4FD030(int a1, short a2);
void nox_xxx_teleportAllPixies_4FD090(nox_object_t* a1);
int nox_xxx_castSpellByUser_4FDD20(int a1, nox_object_t* a2, void* a3);
void nox_xxx_collide_4FDF90(int a1, int a2);
int nox_xxx_spellGetPhoneme_4FE1C0(int a1, char a2);
int sub_4FEA70(int a1, float2* a2);
int nox_xxx_playerCancelSpells_4FEAE0(nox_object_t* a1);
void nox_xxx_spellCancelDurSpell_4FEB10(int a1, nox_object_t* a2);
char* nox_xxx_netStartDurationRaySpell_4FF130(int a1);
int sub_4FF2D0(int a1, int a2);
int nox_xxx_testUnitBuffs_4FF350(nox_object_t* unit, char buff);
void nox_xxx_buffApplyTo_4FF380(nox_object_t* unit, int buff, short dur, char power);
int nox_xxx_unitGetBuffTimer_4FF550(nox_object_t* unit, int buff);
char nox_xxx_buffGetPower_4FF570(nox_object_t* unit, int buff);
void nox_xxx_unitClearBuffs_4FF580(nox_object_t* unit);
int nox_xxx_spellBuffOff_4FF5B0(nox_object_t* a1, int a2);
int nox_xxx_creatureIsMonitored_500CC0(nox_object_t* a1, nox_object_t* a2);
bool nox_xxx_checkSummonedCreaturesLimit_500D70(nox_object_t* a1, int a2);
int nox_xxx_summonStart_500DA0(int a1);
int nox_xxx_summonFinish_5010D0(int a1);
void nox_xxx_summonCancel_5011C0(int a1);
int nox_xxx_charmCreature1_5011F0(int* a1);
int nox_xxx_charmCreatureFinish_5013E0(int* a1);
int nox_xxx_charmCreature2_501690(int a1);
int nox_xxx_getSevenDwords3_501940(int a1);
void nox_xxx_aud_501960(int a1, nox_object_t* a2, int a3, int a4);
void nox_xxx_audCreate_501A30(int a1, float2* a2, int a3, int a4);
void nox_xxx_gameSetAudioFadeoutMb_501AC0(int a1);
uint32_t* nox_xxx_tileAllocTileInCoordList_5040A0(int a1, int a2, float a3);
uint32_t* sub_504290(char a1, char a2);

void nox_server_scriptExecuteFnForEachGroupObj_502670(unsigned char* groupPtr, int expectedType, void* a3,
													  int a4);

#endif // NOX_PORT_GAME4
