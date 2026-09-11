# Weapon and armor equipment port (in progress)

Current scope: 33 functions / 972 physical C lines, including all of MixPatch.c.
Production is unchanged from inventory commit `6d8a5863`:
129,053 physical C lines / 150 files / zero reference C.

The family includes player/NPC equip/dequip and switching, ammo/bow and shield
interactions, modifier engage/disengage dispatch, armor-value calculation,
inventory admission helpers, strength lookup/checks, drop policy initialization,
and the existing gameex shield-selection additions. Preserve the current C
behavior, including community changes that differ from the original executable.

## Scope

- src/legacy/GAME1.c: 00415C00.
- src/legacy/GAME3_3.c: 004E4B20, 004E7D30, 004E7EC0, 004F2F70, 004F2FB0, 004F2FF0, 004F3030, 004F3180.
- src/legacy/GAME4.c: 004F9FD0.
- src/legacy/GAME4_3.c: 0053A030, 0053A0F0, 0053A140, 0053A2C0, 0053A3D0, 0053A420, 0053A680, 0053A6C0, 0053AAB0, 0053AB90, 0053E2D0, 0053E300, 0053E3A0, 0053E430, 0053E520, 0053E600, 0053E650, 0053E7B0, 0053EAE0, 0053EC40, 0053EC80.
- src/legacy/MixPatch.c: entire file: sub_980523, sub_9805EB.

## Original-C fixture work

The optional Equipment spec extends the inventory/shop fixture with action IDs
400..432. Its C dispatcher calls original production functions; callbacks record
modifier arguments and controlled outcomes, without copying equipment algorithms.
Real C-owned equipment definitions supply strength requirements and armor values.
Definitions and lookup strings are guarded. Shared player/item/object/network,
protection, minimap, cache and callback snapshots remain active.

The fixture scopes gameex_flags, the strength cheat, armor clamp threshold,
modifier definition lists and ability activity. It configures active/secondary
weapons, saved shield, player state, masks and strength as explicit inputs.
The locked original-C corpus contains **1,144 cases / 17 groups**: switching
48, modifier callbacks 64, strength 120, armor values 81, admission 128,
NPC owners 32, bow/ammo 64, shields 90, inventory policies 128, sound/secondary
140, initialized drop tables 20, try-equip owners 54, null guards 16, armor
masks 21, cold drop tables 20, armor search 10, and NPC sync 108.
Independent assertions also check modifier return precedence, signed strength
admission, armor-mask classification, and the NPC sync interior-pointer return.

Full JSON captures c-final and c-repeat match byte-for-byte in all 17 groups;
SHA-256 values are locked in equipment_porttest_test.go. Original production C
is still intact. These fixtures exercise the current implementation, including
community shield behavior, actual modifier callbacks, and actual ability lists.
Unlinked armor fixtures retain an inventory holder because the retained equip
packet formatter requires it; list membership and holder relation are distinct.

Ignored scope/caller audit/original blocks: build/port-equipment/{scope.json,
callers.txt,original-c.txt}. Stable captures: c-final-equipment-*.json and
c-repeat-equipment-*.json. Locked equipment plus existing dependency results:
c-locked-dependencies.log. Commit/push this baseline before production edits.
Convert and qualify the family once at its boundary; update C_LOC/docs and
commit/push before continuing.
