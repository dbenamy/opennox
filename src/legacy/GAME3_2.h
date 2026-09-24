#ifndef NOX_PORT_GAME3_2
#define NOX_PORT_GAME3_2

#include "defs.h"

int nox_xxx_updDrawDBall_4CDF80(int a1, int a2);
int sub_4CE0A0(int a1, int a2);
int nox_xxx_updDrawCloud_4CE1D0(int a1, int a2);
int sub_4CE340(int a1, int a2);
int sub_4CE360(int a1, int a2);
int nox_xxx_updDrawColorlight_4CE390(uint32_t* a1, int a2);
int sub_4D6000(nox_object_t* a1);
int sub_4D60B0();
uint32_t* sub_4D60E0(int a1);
int sub_4D6130(int a1);
int sub_4D61F0(int a1);
unsigned int sub_4D6540(int a1);
int nox_xxx_isQuest_4D6F50();
int nox_xxx_setQuest_4D6F60(int a1);
int sub_4D6F70();
int sub_4D6F80(int a1);
int sub_4D6FA0();
int sub_4D72D0(int a1);
int sub_4D7450(int a1, short a2);
int sub_4D75E0();
int sub_4D76E0(int a1);
int sub_4D7A60(int a1);
int nox_xxx_netPrintLineToAll_4DA390(const char* a1);
nox_object_t* nox_get_and_zero_server_objects_4DA3C0(void);
void nox_set_server_objects_4DA3E0(nox_object_t* p);
nox_object_t* nox_xxx_getFirstUpdatable2Object_4DA840();
nox_object_t* nox_xxx_getNextUpdatable2Object_4DA850(nox_object_t* obj);
nox_object_t* nox_server_getFirstObjectUninited_4DA870();
nox_object_t* nox_server_getNextObjectUninited_4DA880(nox_object_t* obj);
void nox_xxx_unitsNewAddToList_4DAC00();
int nox_xxx_gameIsSwitchToSolo_4DB240();
char* nox_xxx_playerCallDisconnect_4DEAB0(int a1, char a2);
void nox_xxx_playerDisconnByPlrID_4DEB00(int a1);
void sub_4DFB50(int a1, int a2);
int nox_xxx_enchantItemTestInventory_4DFBB0(int a1, char a2);
void nox_xxx_effectSpeedEngage_4DFC30(int a1, int a2);
void sub_4DFD10(int a1, int a2);
void nox_xxx_buff_4DFD80(int a1, int a2);
void nox_xxx_checkPoisonProtectEnch_4DFDE0(int a1, int a2);
double nox_xxx_checkFireProtect_4DFE40(uint32_t* a1);
double nox_xxx_checkElectrProtect_4DFF40(uint32_t* a1);
double nox_xxx_getPoisonDmg_4E0040(uint32_t* a1);
void sub_4E0140(int a1, int a2);
void sub_4E0170(int a1, int a2);
float* sub_4E0370(int a1, int a2, int a3, int a4, int a5, float* a6);
float* sub_4E0380(int a1, int a2, int a3, int a4, int a5, float* a6);
int nox_xxx_inversionEffect_4E03D0(int a1, int a2, int a3, int a4, int a5, int* a6);
int nox_xxx_gripEffect_4E0480(int a1, int a2, int a3, int a4, int a5, int* a6);
float* nox_xxx_effectDamageMultiplier_4E04C0(int a1, int a2, int a3, int a4, float* a5);
void nox_xxx_fireEffect_4E0550(void* a1, nox_object_t* a2, nox_object_t* a3, nox_object_t* a4);
void nox_xxx_recoilEffect_4E0640(int a1, int a2, int a3, int a4);
void nox_xxx_lightngEffect_4E06F0(int a1, int a2, int a3, int a4);
int nox_xxx_itemCheckReadinessEffect_4E0960(int a1);
int nox_xxx_effectProjectileSpeed_4E09B0(int a1, int a2, int a3, int a4, int a5);

#endif // NOX_PORT_GAME3_2
