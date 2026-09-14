# Map population and prefabs — qualified C baseline

Current connected batch after map-painting qualification: **38 connected routines / 1,758 C
section lines**. The refreshed virtual-removal audit finds 13 retained ABIs and 25
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

The stack-record/ownership audit produced two pushed prerequisite checkpoints:
683b008f repairs spell-name storage, point records and invalid-book disposal;
41985ada repairs candidate arrays and initializes the complete item-attribute
buffer. Focused original-C failures and corrected regressions support those
changes. Preserve other questionable gameplay behavior during conversion,
including grid truncation before exact containment and case-sensitive flags.

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
abi-audit.json,named-globals.json}. Original inventories are archived locally.
The shared
painting fixture accepts population globals/dispatch and tracks InitData and
Field189 allocation ownership. Initial groups exercise complete 64-bit returns,
cyclic distance propagation, prefab coordinate candidates, metadata lookup and
progress suppression/timing. Captures are exploratory, not yet locked baselines.

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

Checkpoint **41985ada** is pushed. The next expansion adds a synthetic, empty
prefab-file fixture through the real unwrapped loader (including normal file
handle initialization, header/bounds and section terminator). This does not
substitute for unavailable game prefab assets. Together with 576 exact modifier
boundary cases and the first successful hallway connection it passes in 1.526s
(connections.log). Four-direction connection gaps/widths/missing candidates and
complete service-global capture are now being qualified in complete.log.
Refreshed scope remains 38 routines / 1,758 C section lines / 13 retained ABIs;
original inventories are preserved under original-inventory.

Final baseline review adds density clamping, prefab placement retries followed
by exclusion cleanup, and removal of marker objects plus their decoded-cache
list nodes. The marker fixture records disposals at the service boundary before
snapshotting; disposed objects are excluded from live-object traversal. This
fixes a fixture-only attempt to inspect a freed object, without production edits.
Current qualification is in final.log; lock-baseline.py is prepared but has not
run. Do not apply conversion edits before repeated captures are locked/pushed.

Final default qualification passes in **22.538s** with existing painting/room
checks. All **4,221 cases / 36 complete captures** repeat byte-for-byte in an
independent population run (**12.997s**); see final-repeat-comparison.json.
Server/highres verification is running before hashes are made mandatory. The
remaining differences from checkpoint 41985ada are test fixtures/corpora and
documentation; production C is unchanged.

## Locked corrected-C baseline

All **4,221 cases / 36 complete capture groups** repeat byte-for-byte and match
default/server/highres. Population with existing painting/room tests passes in
**22.538s / 21.741s / 22.830s**; the independent default population repeat takes
**12.997s**. Every hash is mandatory in
src/map_population_baseline_porttest_test.go. Evidence:
build/port-map-population/baseline-qualification.json, final-repeat-comparison.json
and final{,-repeat,-server,-highres}.log. The locked-hash smoke is in
locked-baseline.log. Commit/push this checkpoint before editing production C.

The baseline includes explicit density clamps, bounded placement retries and
exclusion cleanup, plus all 16 marker combinations with surviving-cache topology
and exact object/node disposal counts. It uses real waypoint and room-connection
services and a synthetic file decoded through the real loader. Unavailable
AreaMap.lib game content remains an integration limitation, not an assumed pass.
The ordinary asset-backed headless scenario remains required after conversion.

To reproduce the focused checks after setting the documented 386/CGO environment:

```sh
cd src
GOMAXPROCS=2 go test -p 2 -tags porttest -count=1 -run '^TestMap(Population|Painting|Room)' .
GOMAXPROCS=2 go test -p 2 -tags porttest,server -count=1 -run '^TestMap(Population|Painting|Room)' .
GOMAXPROCS=2 go test -p 2 -tags porttest,highres -count=1 -run '^TestMap(Population|Painting|Room)' .
```

Set OPENNOX_MAP_POPULATION_CAPTURE to an absolute filename prefix to also retain
complete JSON snapshots for diagnosis; mandatory hash checks remain enabled.
The baseline hashes live in Git, so they survive loss of the ignored build tree.
Production C is **106,609 / 149 files / zero reference C**; all 38 routines still
use C at this checkpoint. The refreshed ABI audit retains 13 and retires 25.

Baseline **38c905a9** is committed/pushed; mandatory hashes pass in **13.906s**.
The 38-routine native implementation is now written and switched into the tagged
dispatcher. The old C algorithms and populate.c are removed, with 13 retained
exports and 25 retired internal helpers. Current unqualified C count:
**104,839 / 148 files / zero reference**, −1,770 physical lines. First native
build/comparison is in native-first.log. No baseline hashes have changed and
qualification remains pending. Do not rerun apply-native.py.

## Native comparison

All **4,221 cases / 36 complete captures** match corrected C byte-for-byte; no
expected hashes changed. The first native run matched 35/36. The only difference
was progress timestamp rollover: C's mixed signed/unsigned expression compares
as unsigned. Correcting that comparison produces the exact expected state.
Native population plus painting/room checks pass in **23.322s**. Evidence:
native-corrected.log, native-comparison.json and native-first-progress-differences.json.

Source audit finds zero references to the 25 retired helper names in production
C/Go/headers. Full accumulated tests, three production builds/symbol checks,
asset-backed full-suite comparison and unchanged headless gameplay are running
under qualify-all.py. Qualification and the native commit/push remain pending.
Qualified C baseline duplicate captures were compressed losslessly; checksums
are in baseline-capture-archive.json, separate from exploratory archives.

## Completed native qualification — 2026-09-14

All 38 routines are native; populate.c is removed. **104,839 physical production
C lines / 148 files / zero reference C**, a reduction of **1,770 lines**.
All 4,221 cases / 36 complete captures match baseline **38c905a9** unchanged.
Accumulated 61,799 cases / 947 groups pass in default/server/highres:
250.004s / 325.610s / 261.790s, with the additional room and painting contracts.
All three production builds pass ELF32/i386/SSE2/CGO and symbol verification:
25 retired helpers absent, 13 retained entry points present, test loader adapter
absent. The asset-backed full suite exactly matches the known 1,553 failure
entries (15 pass, 3 fail, 32 skip packages). Fresh unchanged repeat-a gameplay
passes with Xvfb/null audio in 36.118s. See qualification.json under
build/port-map-population for machine-readable results. The AreaMap.lib limitation
above remains. This completed result supersedes the intermediate pending status.

Next: the six connected hallway-routing routines in GAME5.c. Preserve the locked
population hashes while extending bent-route, obstruction and fallback coverage.
Do not rerun population scope.py after C removal; recover inventories from Git.
