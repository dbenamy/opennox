#ifndef NOX_PORT_GAME3_3
#define NOX_PORT_GAME3_3

#include "defs.h"

double sub_4E3CA0();
int nox_game_getQuestStage_4E3CC0();
void nox_game_setQuestStage_4E3CD0(int a1);
int nox_xxx_player_4E3CE0();
int sub_4E3D50();
void sub_4E3DD0();
int sub_4E4F30(int a1);
int nox_xxx_playerResetImportantCtr_4E4F40(int a1);
int nox_xxx_netSendPacket1_4E5390(int a1, int a2, int a3, int a4, int a5);
int sub_4E55F0(unsigned char a1);
void nox_xxx_playerLeaveObserver_0_4E6AA0(nox_playerInfo* pl);
unsigned char* sub_4E8E50();
int sub_4E8E60();
bool nox_server_questMaybeWarp_4E8F60();
int sub_4E9010();
int sub_4E9A30(nox_object_t* a1, nox_object_t* a2);
void nox_xxx_pickupFlagCtf_4EA490(int a1, int a2);
int sub_4EA7A0(int a1);
short sub_4EA800(int a1, int a2);
int sub_4EB250(int a1);
void sub_4EB340(float* a1, int a2);
void sub_4EB3E0(int a1);
int sub_4EB9B0(int a1, int a2);
void nox_xxx_unitRemoveChild_4EC470(nox_object_t* a1);
int sub_4ECBD0(int a1);
int sub_4ECC00(char** a1);
int nox_xxx_netGetUnitByExtent_4ED020(int a1);
int nox_xxx_plrReadVals_4EEDC0(nox_object_t* a1, int a2);
int sub_4EF140(int a1);
double nox_xxx_calcBoltDamage_4EF1E0(int a1, int a2);
void sub_4EF410(int a1, unsigned char a2);
int sub_4EF6F0(int a1);
nox_object_t* nox_xxx_playerRespawnItem_4EF750(nox_object_t* a1, char* a2, int* a3, int a4, int a5);
char nox_xxx_playerMakeDefItems_4EF7D0(int a1, int a2, int a3);
int sub_4EFF10(int a1);
int sub_4F2590(int a1);
int sub_4F2C30(int a1);
int nox_xxx_spell_4F2E70(int a1);
int sub_4F2EF0(int a1);
int nox_xxx_inventoryServPlace_4F36F0(nox_object_t* a1p, nox_object_t* a2p, int a3, int a4);

#endif // NOX_PORT_GAME3_3
