#ifndef NOX_PORT_GAME4
#define NOX_PORT_GAME4

#include "defs.h"

void nox_xxx_mapFindPlayerStart_4F7AB0(float2* a1, nox_object_t* a2p);
int nox_xxx_weaponGetStaminaByType_4F7E80(int a1);
void nox_xxx_updatePlayer_4F8100(nox_object_t* a1);
int sub_4FA280(int a1);
int nox_common_mapPlrActionToStateId_4FA2B0(nox_object_t* a1);
int nox_xxx_checkInversionEffect_4FA4F0(int a1, int a2);
char nox_xxx_mobMorphFromPlayer_4FAAC0(uint32_t* a1);
char nox_xxx_mobMorphToPlayer_4FAAF0(uint32_t* a1);
int nox_xxx_updatePlayerMonsterBot_4FAB20(uint32_t* a1);
int nox_xxx_netSendRewardNotify_4FAD50(int a1, int a2, int a3, char a4);
void sub_4FADD0(int a1, char* a2, char a3);
int sub_4FB050(int a1, int a2, int* a3);
void nox_xxx_teleportAllPixies_4FD090(nox_object_t* a1);
int nox_xxx_summonStart_500DA0(int a1);
int nox_xxx_summonFinish_5010D0(int a1);
void nox_xxx_summonCancel_5011C0(int a1);
int nox_xxx_charmCreature1_5011F0(int* a1);
int nox_xxx_charmCreatureFinish_5013E0(int* a1);
int nox_xxx_charmCreature2_501690(int a1);
uint32_t* nox_xxx_tileAllocTileInCoordList_5040A0(int a1, int a2, float a3);
uint32_t* sub_504290(char a1, char a2);

void nox_server_scriptExecuteFnForEachGroupObj_502670(unsigned char* groupPtr, int expectedType, void* a3,
													  int a4);

#endif // NOX_PORT_GAME4
