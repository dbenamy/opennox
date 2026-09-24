#ifndef NOX_PORT_GAME1_1
#define NOX_PORT_GAME1_1

#include "common__savegame.h"
#include "defs.h"
#include "common__system__team.h"

int nox_xxx_servObjectHasTeam_419130(int a1);
int nox_xxx_servCompareTeams_419150(int a1, int a2);
int sub_41C280(void* a1);
int nox_xxx_parseFileInfoData_41C3B0(int a1);
int sub_41C780(int a1);
int nox_xxx_netSavePlayer_41CE00();
int sub_41CEE0(void* a1, int a2);
void* nox_xxx_monsterGetSoundSet_424300(nox_object_t* a1);
void* nox_xxx_updateSpellRelated_424830(void* a1, int a2);
void sub_4257F0(int* a1, uint32_t* a2);
nox_list_item_t* nox_common_list_getFirstSafe_425890(nox_list_item_t* list);
nox_list_item_t* nox_common_list_getNextSafe_4258A0(nox_list_item_t* list);
uint32_t* sub_4258C0(uint32_t** a1, int a2);
void nox_common_list_append_4258E0(nox_list_item_t* list, nox_list_item_t* cur);
uint32_t* sub_425900(uint32_t* a1, uint32_t* a2);
void nox_common_list_remove_425920(void* a1);
nox_list_item_t* nox_common_list_getNext_425940(nox_list_item_t* list);
int sub_425960(int a1);
int* sub_425A50();
int* sub_425A60(int* a1);
int* sub_425A70(int a1);
char* sub_425B60(void* lpMem, int a2);


#endif // NOX_PORT_GAME1_1
