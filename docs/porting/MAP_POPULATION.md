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
bounded real-file integration check for decoding/placement. General game assets
are available, but the supplied image lacks the required AreaMap.lib prefab library. Preserve cache ownership and allocator boundaries.

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
abi-audit.json,named-globals.json}. Initial fixture work is in progress; production C remains unchanged. The shared
painting fixture accepts population globals/dispatch and tracks InitData and
Field189 allocation ownership. Initial groups exercise complete 64-bit returns,
cyclic distance propagation, prefab coordinate candidates, metadata lookup and
progress suppression/timing. Captures are exploratory, not yet locked baselines.

Additional prerequisite candidates from source review: 5224B0 passes a single
32-bit local as an 8-byte point output; 51E1D0 formats a complete spell identifier
into a single-byte local. Establish focused execution evidence before repairs.
The positive progress timing fixture uses the real callback with an already-seen
GUI sequence, which exercises counter/tick updates without drawing. Actual GUI
rendering remains covered by the headless gameplay scenario.

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


## Active prerequisite findings (2026-09-13)

Original production C probes run in separate processes confirm that spell lookup
for `fireball` crashes and the 5224B0 wildcard-book placement path aborts with
stack-smashing detection. The finale returns but places PlayerStart at the
clamped map corner (5851.5, 5852.5), despite the supplied room having an interior
center. The fixture now independently computes the expected center transform.
Evidence: build/port-map-population/probe-{spell-name,point-output,finale-point}.log
and c-probe-prerequisite-finale-point.json. The invalid-book probe first hits
the spell-name defect, so repeat it after that repair to isolate the allocator.

Three production C prerequisites are applied but awaiting qualification: give
spell formatting its full 66-byte output record, reject names of 60 or more
bytes before copying into the 60-byte input record, and use explicit two-float
records in 5224B0 and the finale. Production C is now **106,611 / 149 files /
zero reference C** (+2). No conversion has started and no hashes are locked yet.

The initial 191 scalar/topology/coordinate/timing cases passed alongside all
88 unchanged painting hashes. A further 580 inventory/book/item cases pass;
the exit-name fixture exposed a fixture offset mistake (CollideData is byte 700,
InitData is byte 692), corrected before admitting exit captures. Real pool and
allocation owners are used throughout. Repeat full captures after qualification.


After the stack repairs, all three direct regressions pass and all eight existing
population capture groups remain byte-identical; all 88 painting captures pass.
The invalid-book probe now reaches its allocator and aborts with `free(): invalid
size`. Replace that raw free with the existing engine `objectFreeMem` service.
The regression uses ordinary game flags to isolate pool/UseData disposal and
asserts a zero return plus zero live objects. The four prerequisite probes also
run as normal subtests, so they cannot silently disappear behind a diagnostic
environment variable. Evidence: stack-probe-*.log and prereq-fixed.log.


## Qualified prerequisite checkpoint

All **1,067 cases / 14 capture groups** agree byte-for-byte across default,
server and highres. The new population tests and existing room/painting tests
pass together in **12.619s / 12.302s / 13.366s** package time. Thirteen earlier
population groups also repeat exactly in independent default runs. Every
registered spell name is checked against its exact ID in both upper/lower case;
length boundaries and all four repaired paths have ordinary regression tests.
Evidence: build/port-map-population/prerequisite-qualification.json and
prereq-{repeat,server,highres}.log. Production C: **106,611 / 149 / zero reference**.

Commit this as a prerequisite recovery checkpoint. The full population/prefab
baseline is still incomplete: generated captures remain exploratory rather than
mandatory conversion oracles. Next expand decoded-cache marker scanning, prefab
candidate lists/room allocation/replacement, population ordering and actual
placement. Audit the two six-element candidate arrays in 526550: each is still
represented as a scalar plus a separate five-element local array. A staged
isolated probe is in build/port-map-population/prefab-candidate-probe.go.stage.
Recheck connected-list topology and all positive paths before locking the full
baseline. Production builds/full-suite/gameplay remain batch-boundary checks.


Prerequisites committed/pushed as **683b008f**; continue the full baseline.
The subsequent candidate probe crashes in 526550 when the list has two rooms.
Both scalar-plus-five-element records are now explicit six-element arrays.
One-to-eight-room candidate chains pass order, limit and link checks. C is back
to **106,609 / 149 / zero reference** (−2); this follow-up is not committed yet.
Evidence: candidate-probe.log (original) and candidate-arrays.log (corrected).

The expanded fixture sets each case's game flags before preallocating objects,
fixing a dependence on the preceding case. It tracks occupancy grid rows and
allocated metadata, including freed-state snapshots and zero-sized allocation
identity without reading beyond the allocation. Startup modifier chances and
prefab flag characters come from immutable blob data.

The real prefab-selection service 502D70 reloads the file even when a decoded
cache is present. For bounded selection/geometry tests, a **porttest-only linker
wrapper** supplies an explicitly seeded decoded cache and records load count and
index. Enumeration, metadata processing and room operations still execute their
real production implementations. The wrapper delegates to the real loader when
that fixture option is off; production binaries contain neither the wrapper nor
its linker flag. Asset-backed decoding/placement remains a separate integration
requirement; seeded-cache tests must not be described as exercising file decoding.


Expanded C tests through ordering pass in **7.625s** (ordering.log). Coverage now
includes required/optional prefab selection, nearest-six selection in all four
directions with randomized input order, actual room occupancy/replacement,
weighted inventories, shuffled population lists with full forward/back links,
actual object counts, and real waypoint creation plus all pair connections.
A build-tagged observer at the native room-release boundary records freed rooms
without reading released memory; production uses an empty inlineable observer.
The population fixture owns waypoint cleanup through the existing server owner.

Further tests for actual decoded-cache application, exit placement and modifier
selection are in progress in placement.log. The file decoder alone is replaced
by the explicit seeded-cache service option; the prefab application routine is
being exercised directly. All population captures remain exploratory until the
entire scope and repeatability audit finish. Existing painting/room hashes stay
mandatory and must pass unchanged after fixture extensions.

The placement expansion passes in **9.118s**, covering 27 complete capture groups.
An independent repeat agrees in 26 groups and isolates 356 differences to the
fifth word of enchanted-item attribute data. The C caller initialized only four
words but its setter copies five. A direct regression fails before repair; the
complete buffer is now zero-initialized. Corrected population plus existing
painting/room checks are running in modifier-fixed.log. See the decision record
and modifier-repeat-diagnosis.json; baseline locking remains pending.

## Qualified expanded prerequisite checkpoint

**3,348 cases / 30 complete captures** repeat byte-for-byte and match default,
server and highres. Population-only default passes in **10.138s**; server/highres
with existing painting/room tests pass in **19.320s / 19.779s**. Default painting/room
plus population passed in **18.809s** before the final root-only corpus expansion.
The item fifth word is consistently zero. New coverage includes stable distance
sorting (ties, reverse/permuted lists), theme thresholds, uppercase/case-sensitive
prefab flags, rectangle fit/exclusions, and real hallway waypoint connections.
Evidence: followup-qualification.json and followup-repeat-comparison.json.

Prefab dimensions are truncated to cells before random position selection; an
attempt can pass that coarse check and fail exact rectangle containment. Tests
preserve that behavior, independently assert centered success and oversized/
excluded rejection, and capture all state/RNG on both paths. Do not silently
change this placement policy during conversion.

The supplied 7z contains a filesystem image. Listing the entire image confirms
**no AreaMap.lib**; the extracted tree also lacks it. The loader's filename comes
from the original startup data. Asset-backed prefab decoding is therefore not
currently available; this differs from ordinary asset-backed gameplay validation,
which remains available and required. Seeded-cache application is covered; add a
small synthetic real-file fixture if practical before locking the baseline.
No new downloads or replacement asset assumptions are needed for this checkpoint.

The full population baseline is still **not locked**, and all 38 routines remain
C. Remaining review includes positive prefab connection service calls and
modifier/placement edge coverage. C count is **106,609 / 149 files / zero reference**.
Earlier exploratory captures are gzip-compressed losslessly; their original byte
counts and SHA-256 checksums are in exploratory-capture-archive.json.
