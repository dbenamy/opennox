#ifndef NOX_PORT_GAME5_2
#define NOX_PORT_GAME5_2

#include "defs.h"

void nullsub_22();
void nullsub_36();
void nullsub_29(void);
int sub_56F250();
int nox_xxx_protectionCreateInt_56F400(int a1);
int sub_56F4F0(int* a1);
uint32_t* sub_56F590(int a1);
int nox_xxx_protectData_56F5C0();
uint32_t sub_56F780(int a1, int a2);
uint32_t nox_xxx_playerResetProtectionCRC_56F7D0(int a1, int a2);
uint32_t sub_56F820(int a1, unsigned char a2);
uint32_t nox_xxx_protectPlayerHPMana_56F870(int a1, unsigned short a2);
uint32_t sub_56F8C0(int a1, float a2);
uint32_t sub_56F920(int a1, int a2);
uint32_t nox_xxx_protectMana_56F9E0(int a1, short a2);
uint32_t sub_56FA40(int a1, float a2);
int nox_xxx_protectionStringCRCLen_56FAE0(int* a1, unsigned int a2);
int sub_56FB00(int* a1, unsigned int a2, int a3);
int nox_xxx_protect_56FBF0(int a1, nox_object_t* item);
int nox_xxx_protect_56FC50(int a1, nox_object_t* object);
int nox_xxx_playerAwardSpellProtectionCRC_56FCE0(int a1, int a2, int a3);
int nox_xxx_playerApplyProtectionCRC_56FD50(int a1, void* a2, int a3);
unsigned int nox_xxx_netGetUnitCodeCli_578B00(int a1);
nox_waypoint_t* nox_xxx_waypointGetList_579860();
int nox_xxx_waypointNext_579870(int a1);
int sub_5798A0(int a1);
unsigned int nox_xxx_waypoint_5798C0();
char* sub_579A30();
int sub_579CA0();
uint32_t* sub_579E70();
char* sub_57A1B0(short a1);
char sub_57A1E0(int* a1, char* a2, int* a3, char a4, short a5);
int sub_57A950(char* a1);
int sub_57A9F0(char* a1, char* a2);
char sub_57AAA0(char* a1, char* a2, int* a3);
int nox_xxx_playerCheckSpellClass_57AEA0(int a1, int a2);
int sub_57B190(unsigned short a1, unsigned short a2);
int nox_xxx_client_57B400(int a1);
int nox_xxx_collideReflect_57B810(float* a1, int a2);
int nox_xxx_map_57B850(float2* a1, float* a2, float2* a3);
int sub_57B920(void* a1);
char nox_xxx_cliGenerateAlias_57B9A0(int a1, int a2, int a3, unsigned int a4);
int sub_57BA10(int a1, short a2, short a3, int a4);
void* nox_server_getFirstMapGroup_57C080();
int nox_xxx_mathPointOnTheLine_57C8A0(float4* a1, float2* a2, float2* a3);

#endif // NOX_PORT_GAME5_2
