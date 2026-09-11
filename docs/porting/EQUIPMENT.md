# Weapon and armor equipment port

Current scope: 33 functions / 972 physical C lines, including all of MixPatch.c.
Original-C baseline: `6e13a789` (committed and pushed before conversion).
Native implementation removes **972 C lines**, leaving **128,081 physical C
lines / 149 files / zero reference C**. Qualification is complete.

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
SHA-256 values are locked in equipment_porttest_test.go. The original-C
implementation is available at the baseline commit. These fixtures exercise
the current implementation, including
community shield behavior, actual modifier callbacks, and actual ability lists.
Unlinked armor fixtures retain an inventory holder because the retained equip
packet formatter requires it; list membership and holder relation are distinct.

Ignored scope/caller audit/original blocks: build/port-equipment/{scope.json,
callers.txt,original-c.txt}. Stable captures: c-final-equipment-*.json and
c-repeat-equipment-*.json. Locked equipment plus existing dependency results:
c-locked-dependencies.log. The baseline was committed/pushed before production
edits; the connected family was qualified together.

## Native implementation and review

- equipment.go: strength, inventory queries, modifier dispatch, armor values,
  NPC synchronization, shield selection, sounds, and shared drop-policy table.
- equipment_weapon.go / equipment_armor.go: player and NPC transitions,
  ammo/bow handling, opposing equipment sweeps, and admission checks.
- equipment_exports.go: the 33 existing ABI entry points, used by retained C
  callers and the original contract dispatcher. No copied C algorithms remain.
  Two header const qualifiers are aligned with generated cgo declarations.
- Existing inventory, trade, generator, AI loot and object Go callers route
  directly to the native functions.

All 17 original-C JSON capture groups match native-third byte-for-byte (3.809s).
The baseline test plus inventory/resource/shop dependencies passed in 38.796s.
Broader native qualification results are recorded below.

Preserved details include last-slot modifier return precedence; the community
shield-selection flag and exact equality for fallback shield flags; player
armor dequip's object+48 byte check; signed strength versus unsigned 16-bit
requirements; float32 armor accumulator spills with double helper results;
and the NPC sync helper's raw interior-pointer return. Count with filter zero
includes destroyed items; filtered counts exclude them. Duplicate comparison
4E7DE0 and the existing network serializers remain production dependencies.

The defend callback's writable float uses scoped C-owned storage. Modifier
algorithms remain in their existing production locations for the next connected
batch, while equipment dispatch itself is native.

## Qualification

- Native equipment plus existing inventory/resource/shop/trade: **12,009 cases**,
  all hashes unchanged, 58.052s (concurrent build load).
- Accumulated default/server/highres port corpora: pass in **101.875s /
  91.325s / 96.494s**.
- Three production builds: opennox, opennox-hd, opennox-server, all verified
  ELF32/i386, GOARCH=386, GO386=sse2, CGO_ENABLED=1.
- Whole-module suite: exact known failure multiset, 1,553 entries, zero additions
  or removals; package outcomes 15 pass / 3 fail / 32 skip.
- Fresh equipment-port gameplay: exit 0 in 38.916s, unchanged repeat-a goldens,
  overrides disabled, Xvfb, null audio.

Evidence is under build/port-equipment (native-third captures, native-dependencies,
ports-variant logs, binary-checks, full-suite-comparison), and
build/baseline/runs/equipment-port/result.json. Old regenerable build-cache files
were pruned for disk space; asset archives and gameplay evidence were preserved.

Next: [modifier effects and weapon use](EFFECTS_USE.md), 41 functions / 977 C
lines with one qualification boundary for the connected family.
