#ifndef NOX_PORT_GAME1
#define NOX_PORT_GAME1

#include "defs.h"
#include "common__system__team.h"



void nox_common_setEngineFlag(const nox_engine_flag flags);
void nox_common_resetEngineFlag(const nox_engine_flag flags);
void nox_xxx_servStartCountdown_40A2A0(int a1, char* a2);
int nox_xxx_utilFindSound_40AF50(char* a1);
int* nox_xxx_wallDestroyedByWallid_410520(short a1);
void* nox_xxx_wallGetFirstBreakableCli_410870();
int nox_xxx_wallGetNextBreakableCli_410880(int* a1);
int sub_411490(int a1, int a2);
char* nox_xxx_cliGamedataGet_416590(int a1);
char* sub_4165B0();
void* sub_416640();


#endif // NOX_PORT_GAME1
