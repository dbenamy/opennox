# Map population and prefabs — next batch

Planned after map-painting qualification: **38 connected routines / 1,758 C
section lines**. Initial virtual-removal audit finds 13 retained ABIs and 25
internal helpers. Recheck Go wrappers/getters before conversion. The batch covers
population selection and ordering, generated inventory/spellbooks/enchantments,
monster and exit placement, room selection, prefab selection/geometry/placement,
and the generation progress callback. Theme parsing and the prefab file decoder
remain services. All eight functions in server__mapgen__generate__populate.c
are included; remove that file if its remaining declarations become unnecessary.

## Baseline and testing plan

Use the newly qualified room and painting owners and real object allocation,
linked inventories, exclusions, type/Xfer checks and startup data. Capture complete
records and pointer-normalized topology, object lifecycle counters, all touched
selection/prefab globals, both random-stream tails, return bits, guards and x87
PC53/nearest state. Seed linked prefab metadata and decoded object lists in the
existing cache so remaining C lookup/enumeration paths execute normally. Add a
bounded asset-backed integration check for decoding/placement; asset availability
is already established locally. Preserve cache ownership and allocator boundaries.

Cover weighted branch thresholds, randomized ordering with full list integrity,
empty/single/multiple lists, explicit and wildcard items/spells, modifier chances
and dependent-slot rules, failed allocation/name/type paths, coordinate rounding,
room bounds/exclusion rejection, all prefab edges/corners/center preferences,
retry limits and candidate exhaustion, exit selection, type-checked exit-name
copies and connected population/placement sequences. Require positive contracts
for actual objects, correct inventory links, room occupancy and nonempty output;
smoke success and hashes alone are insufficient. Keep boundary cases valid for
the documented record capacities, and separately test defined invalid-input gates.

Audit decompiler stack records before locking: the population finale passes &v8
as a two-float PlayerStart coordinate while v9 is a separate local. Prove the
observable defect with an original-C independent check, then replace the pair
with explicit storage as a documented prerequisite if required. Spellbook failure
currently calls free on an object from the engine pool; verify owner semantics
and use the appropriate engine disposal path if the failure test confirms the
mismatch. Neither repair is applied yet. Preserve questionable gameplay behavior
unless a demonstrated correctness/ownership defect requires a reversible fix.

Set PlatformTicks deterministically for progress timing, including wrap/reset
and threshold boundaries. The current input-poll export is a no-op; the progress
GUI callback needs an isolated recorder or initialized headless GUI. Exercise
suppression under the existing generation game flag as well as positive updates.
Other retained callbacks should use existing server services, or thin recording
adapters at service boundaries when full subsystem setup is unrelated to scope.

Repeat and compare complete original/corrected C captures, make all hashes
mandatory, and commit/push the baseline before conversion. Then convert/audit,
run accumulated default/server/highres tests, build all three production targets,
verify retired symbols, compare the known asset-backed full-suite failures and
run unchanged headless gameplay. Count physical C, document, commit/push and
continue. No pending question or new agents.

Local reconstructible inventories: build/port-map-population/{scope.py,audit.py,
candidate-scope.json,candidate-source.txt,candidate-dependencies.json,operations.json,
abi-audit.json,named-globals.json}. Planning only; no fixture or production edits.

## Candidate entry points

| Entry point | C section lines | Initial ABI audit |
| --- | ---: | --- |
| `sub_4D42C0` | 3 | Retire |
| `nox_xxx_mapGenSpellIdByName_51E1D0` | 15 | Retain |
| `nox_xxx_mapgen_521C60` | 26 | Retire |
| `sub_521CB0` | 80 | Retire |
| `nox_xxx_mapgen_521FE0` | 59 | Retire |
| `nox_xxx_mapGenMakeSpellbook_5220E0` | 26 | Retire |
| `nox_xxx_mapGenMakeEnchantedItem_5221A0` | 70 | Retire |
| `sub_522300` | 17 | Retire |
| `nox_xxx_mapgen_522340` | 9 | Retain |
| `nox_xxx_mapgen_522370` | 82 | Retire |
| `nox_xxx_mapgen_5224B0` | 95 | Retire |
| `nox_xxx_mapGenMakeMonsterInRoom_522810` | 33 | Retire |
| `nox_xxx_mapGenFinishPopulate_5228B0_mapgen_populate` | 46 | Retain |
| `nox_xxx_mapGenMakeExit_522A40` | 24 | Retire |
| `nox_xxx_mapgen_522AD0` | 85 | Retire |
| `sub_522C80` | 11 | Retain |
| `sub_522D30` | 83 | Retain |
| `nox_xxx_mapGenTryNextRoom_522F40` | 44 | Retain |
| `nox_xxx_mapGenSetFlags_5235F0` | 20 | Retain |
| `sub_5259D0` | 3 | Retire |
| `sub_5259E0` | 3 | Retire |
| `sub_5259F0` | 62 | Retain |
| `sub_525AF0` | 55 | Retain |
| `sub_525BF0` | 40 | Retire |
| `sub_525C90` | 52 | Retire |
| `nox_xxx_mapGen_InPrefab1_525D20` | 38 | Retain |
| `sub_525DF0` | 143 | Retire |
| `nox_xxx_mapGenPrefabMkRoom_526100` | 48 | Retire |
| `sub_526260` | 17 | Retire |
| `sub_5262F0` | 132 | Retire |
| `sub_526550` | 143 | Retire |
| `nox_xxx_mapGen_InPrefab2_5266F0` | 62 | Retain |
| `nox_xxx_mapGenPlacePrefabs_526830` | 34 | Retain |
| `sub_5268F0` | 17 | Retire |
| `sub_526950` | 52 | Retain |
| `sub_526A90` | 3 | Retire |
| `sub_526AA0` | 7 | Retire |
| `sub_527D50` | 19 | Retire |
