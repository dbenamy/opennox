#ifndef NOX_PORT_GAME1
#define NOX_PORT_GAME1

#include "defs.h"
#include "common__system__team.h"



void nox_common_setEngineFlag(const nox_engine_flag flags);
void nox_common_resetEngineFlag(const nox_engine_flag flags);
bool nox_common_getEngineFlag(const nox_engine_flag flags);
char* nox_server_currentMapGetFilename_409B30();
char* nox_xxx_gameSetMapPath_409D70(char* a1);
int sub_409F40(int a1);
int sub_40A1A0();
void nox_xxx_servStartCountdown_40A2A0(int a1, char* a2);
int sub_40A300();
int nox_xxx_utilFindSound_40AF50(char* a1);
int* nox_xxx_wallDestroyedByWallid_410520(short a1);
void* nox_xxx_wallGetFirstBreakableCli_410870();
int nox_xxx_wallGetNextBreakableCli_410880(int* a1);
int sub_411490(int a1, int a2);
int nox_xxx_modifGetIdByName_413290(char* a1);
void* nox_xxx_modifGetDescById_413330(int a1);
void sub_413A00(int a1);
int nox_xxx_weaponInventoryEquipFlags_415820(nox_object_t* item);
int sub_415840(int a1);
int nox_xxx_ammoCheck_415880(int typ_ind);
double nox_xxx_itemApplyDefendEffect_415C00(int a1);
int nox_xxx_unitArmorInventoryEquipFlags_415C70(nox_object_t* item);
int sub_416150(int a1, int a2);
int sub_416170(int a1);
char* nox_xxx_cliGamedataGet_416590(int a1);
char* sub_4165B0();
char* sub_4165D0(int a1);
void* sub_416640();
void nox_ticks_reset_416D40();
void nox_xxx_netMarkMinimapObject_417190(int a1, nox_object_t* a2, unsigned int a3);
void nox_xxx_netUnmarkMinimapObj_417300(int a1, nox_object_t* a2, unsigned int a3);
bool nox_xxx_CheckGameplayFlags_417DA0(int a1);
int sub_417F50(int a1);


#endif // NOX_PORT_GAME1
