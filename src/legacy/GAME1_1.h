#ifndef NOX_PORT_GAME1_1
#define NOX_PORT_GAME1_1

#include "common__savegame.h"
#include "defs.h"
#include "common__system__team.h"

char* sub_418A40(wchar2_t* a1);
int sub_418BC0(int a1);
uint32_t* nox_xxx_objGetTeamByNetCode_418C80(int a1);
int sub_4190F0(wchar2_t* a1);
int nox_xxx_servObjectHasTeam_419130(int a1);
int nox_xxx_servCompareTeams_419150(int a1, int a2);
char sub_419960(int a1, int a2, short a3);
void sub_419E10(nox_object_t* a1, int a2);
int sub_419E60(nox_object_t* a1);
int sub_419EA0();
int sub_41C200(void* a1, int a2);
int sub_41C280(void* a1);
int nox_xxx_parseFileInfoData_41C3B0(int a1);
int sub_41C780(int a1);
int nox_xxx_netSavePlayer_41CE00();
int sub_41CEE0(void* a1, int a2);
int sub_41E300(int a1);
int* nox_xxx_tileListAddNewSubtile_422160(int a1, int a2, int a3, int a4);
int nox_xxx_tileFreeTile_422200(int a1);
void* nox_xxx_monsterGetSoundSet_424300(nox_object_t* a1);
void* nox_xxx_updateSpellRelated_424830(void* a1, int a2);
int nox_xxx_enchantByName_424880(char* a1);
char* nox_xxx_getEnchantName_4248F0(int a1);
int nox_xxx_getEnchantSpell_424920(int a1);
char* nox_xxx_abilityGetName_425250(int a1);
int nox_xxx_abilityCooldown_4252D0(int a1);
wchar2_t* sub_4252F0(int a1);
nox_video_bag_image_t* nox_xxx_spellGetAbilityIcon_425310(int a1, int a2);
int nox_xxx_bookFirstKnownAbil_425330();
int nox_xxx_bookNextKnownAbil_425350(int a1);
int sub_425450(int a1);
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
char* sub_425CA0(int a1, int a2);
int nox_xxx_mapWriteSectionsMB_426E20(void* a1);
int nox_xxx_mapReadSection_426EA0(void* a1, char* name, uint32_t* a3);
void nox_xxx_comJournalEntryAdd_427500(nox_object_t* a1, char* a2, short a3);
int sub_427980(float4* a1, float4* a2);


#endif // NOX_PORT_GAME1_1
