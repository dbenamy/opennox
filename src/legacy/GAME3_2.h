#ifndef NOX_PORT_GAME3_2
#define NOX_PORT_GAME3_2

#include "defs.h"

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

#endif // NOX_PORT_GAME3_2
