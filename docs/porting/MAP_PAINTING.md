# Map painting and door placement

Room conversion is qualified; this is the next baseline-preparation task. The candidate covers 48 connected
routines / 2,255 C section lines: room floor, wall, border and door emission;
coordinate conversion; tile/subtile mutation and recycling; border propagation;
wall direction composition; and shared object placement/orientation helpers.
Prefab loading and randomized population follow separately. Precise physical
removal depends on preserving any shared forward declarations.

The first virtual-removal audit finds 16 ABIs with production references and 32
internal helpers without them, including Go getter/callback references in the
search. The subtile-free-one bridge has only a Go wrapper left: redirect it to
the native owner during conversion, yielding 15 retained ABIs and 33 retired
bridges/helpers. Recheck before qualification. Remaining services include already-native tile
selection, edge mapping, worklists and room exclusions, plus actual wall/object
owners. Local inventories are under `build/port-map-painting`.

## Qualified native conversion

All **48 routines** are now native Go. Fifteen C ABIs remain; 33 internal
helpers/bridges are retired, including the subtile-free-one export whose final
Go caller now calls the native owner. The retired production-reference audit is
empty. No C algorithms remain solely for tests.

All **3,580 cases / 88 full captures** match the locked corrected-C baseline
byte-for-byte on the first complete native run (31.495s). No expected hash changed.
Full byte comparison is in
build/port-map-painting/native-first-comparison.json.

Physical C is **106,609 lines / 149 files / zero reference**: 2,249 lines removed
by conversion, following the separately committed 19-line prerequisite repair.
Retained shared forward declarations account for the difference from section
counts. Implementations are legacy/map_painting{,_tiles,_walls,_borders,_patterns,
_rooms,_objects,_exports}.go. Existing native selection, edge merging, worklists,
room exclusions and server wall/object services are called directly. Shared
subtile allocations continue using the C heap for remaining C owners.

Accumulated **57,578 captured cases / 911 groups** plus four room-precision
contracts pass in default/server/highres: 256.828s / 332.775s / 264.573s wall time. Two additional
native admission contracts (nil name / stale type index) pass in all variants.
Three production binaries are ELF32/i386/SSE2/CGO; all 33 retired symbols are absent.
Asset-backed full-suite failures match exactly (1,553 entries; 15 pass, 3 fail,
32 skip packages). Fresh unchanged repeat-a gameplay passes in 35.310s.
See build/port-map-painting/qualification.json and baseline/runs/map-painting-port.
The next batch is documented in MAP_POPULATION.md.

## Fixture throughput follow-up

Wall-definition snapshots compare the entire 986,560-byte backing array before
reusing its SHA-256. Any changed byte causes a fresh digest; no state is omitted.
A contract checks mutations outside active definitions, restoration and changes
to variation limits during reset. All 88 complete captures remain byte-identical
to the C baseline, and all painting/admission/snapshot checks pass in all variants.
Root package times fall from 31.495s to **8.769s default / 8.551s server / 9.079s
highres** (about 72% less execution time). Compile time is excluded from these
package measurements. Production source and C count are unchanged: 106,609.
Evidence: build/port-map-painting/cached-snapshot-comparison.json and
cached-snapshot-variants.json. The cache is confined to the tagged fixture.

## Baseline prerequisite discovered

Repeated unmodified C processes disagreed in 26 capture groups. Several
coordinate and tile-pattern records were represented by unrelated scalar locals;
optimized code removed fields and callers read unrelated stack memory. Ten
functions now use explicit contiguous records, and the border entry initializes
its pattern. The corrected baseline is locked after two complete, identical process runs.
See DECISIONS.md for evidence and the seed/layout compatibility review item.
This repair reduces physical C by 19 lines to 108,858; the 48 routines still run
in C. The corrected scope is 48 routines / 2,255 C section lines.

## Locked corrected-C baseline

All **3,580 cases / 88 complete capture groups** pass twice in separate processes
(30.530s / 30.924s) with byte-identical serialized state. Hashes are mandatory in
`src/map_painting_porttest_test.go`; full local captures and metadata are under
`build/port-map-painting/c-final-{a,b}` and `c-baseline.json`.
Accumulated default port regression validation passes (root package 252.436s);
the accumulated package checks succeed. See c-accumulated.json and its log.
Including this corpus, the accumulated totals are 57,578 captured cases / 911
groups, plus the four independent room-precision contracts.

The fixture covers every entry point, 13x13 direction composition, coordinate
aliases and clamping, pool growth/recycling and subtile merge/dedup, tile routing
and anchors, rectangle/line/pattern painting, all floor-corner masks, protected
walls and secret-wall removal, room kinds and connected patterns, real object
placement/promotion and repeated movement, door/monster orientation, spellbook
type checks, layout selection and flood-border propagation. Independent positive
contracts supplement the complete-state hashes. Real startup tables include the
monster direction-angle lookup and its preceding word for the legacy invalid
index path. Sparse grid snapshots preserve every changed cell and all row
pointers; guarded pools, links, complete records, globals, RNG tails and x87 state
are included. Only the fresh server handle and known pointer identities normalize.

## Testing plan

Lock repeated C captures and push the baseline before conversion. Use a guarded
128x128 legacy tile grid with both half-tiles initially NONE (255), actual wall
allocation/indexing, full room/config/pattern/output records, real object type
lookup/allocation, and complete globals and return bits. Encode the grid sparsely
against its fixed initial state so every mutation is captured without repeating
megabytes of unchanged cells. Preserve complete subtile/pool records and links.

Install actual startup direction/edge tables, save/restore all touched blobs,
selection flags, definitions, pool heads and worklists, and retain guards. Seed a
fresh real platform owner and capture the trailing RNG value. Lock/save/restore
the x87 control word at PC53/nearest. Review generated C assembly when decompiled
float declarations do not explain observable rounding.

Root tests can wrap an isolated core in a real `Server` and temporarily install
it as noxServer and legacy.GetServer. That exercises actual CreateObjectAt,
collider/flags/pending-list updates and ObjectsClearPending, as well as the existing
wall-name selector, which reads noxServer directly. Restore both owners afterward.
Use actual DoorXfer/SpellRewardXfer addresses to cover snapping and type checks.

Independent checks cover coordinate clamping/parity/bounds, painting/no-change,
subtile deduplication, recycling and pool expansion, wall creation/update/removal,
protected walls, seams/corners, all room kinds/directions, floor patterns, border
propagation, default/explicit variations, NONE/missing-name gates, object/door/
monster orientation and pending-object promotion. Include complete connected
painting and placement sequences and a smoke case for every entry point.

Keep allocation bases valid for the owning allocator. Track new 200-byte subtile
blocks, generated exclusion rectangles and secret-wall removal without reading
released storage. The tile grid needs its own pointer-identity range; the small
record fixture's 10,000-byte spacing cannot represent its full storage safely.

After conversion: unchanged full captures, accumulated default/server/highres
checks, all three production builds, exact asset-backed full-suite failure
comparison, fresh unchanged headless gameplay, actual C count, docs, commit/push
and continue. No new agents or outstanding user question.

## Candidate entry points

| Entry point | C section lines | Current ABI audit |
| --- | ---: | --- |
| `nox_xxx_tileListAddNewSubtile_422160` | 26 | Retain |
| `nox_xxx_tileFreeTileOne_4221E0` | 10 | Retire |
| `nox_xxx_tileFreeTile_422200` | 19 | Retain |
| `nox_xxx_wall_42A6C0` | 5 | Retain |
| `nox_xxx_mapGenFixCoords_4D3D90` | 22 | Retain |
| `sub_4D3FF0` | 36 | Retire |
| `sub_51D5E0` | 16 | Retire |
| `sub_51D8F0` | 35 | Retain |
| `sub_51D9C0` | 20 | Retire |
| `sub_51DA70` | 147 | Retire |
| `sub_5244D0` | 12 | Retire |
| `sub_524500` | 23 | Retire |
| `sub_524550` | 22 | Retire |
| `sub_5245A0` | 28 | Retain |
| `sub_524610` | 20 | Retire |
| `nox_xxx_gen_524680` | 157 | Retire |
| `sub_524950` | 12 | Retire |
| `sub_5249C0` | 75 | Retire |
| `sub_524B50` | 93 | Retire |
| `nox_xxx_gen_524E00` | 84 | Retain |
| `sub_524FB0` | 176 | Retire |
| `sub_525330` | 22 | Retire |
| `sub_525370` | 22 | Retire |
| `sub_5253B0` | 76 | Retire |
| `nox_xxx_mapgen_525510` | 27 | Retire |
| `nox_xxx_mapgen_525570` | 58 | Retire |
| `nox_xxx_mapgen_525690` | 23 | Retire |
| `nox_xxx_mapgen_525740` | 31 | Retire |
| `nox_xxx_mapgen_525830` | 23 | Retire |
| `nox_xxx_mapgen_5258E0` | 31 | Retire |
| `sub_526C40` | 9 | Retain |
| `sub_526C80` | 9 | Retire |
| `sub_526D50` | 14 | Retire |
| `sub_526DD0` | 38 | Retire |
| `sub_526E60` | 75 | Retire |
| `sub_527030` | 43 | Retain |
| `sub_527380` | 37 | Retire |
| `sub_527450` | 228 | Retire |
| `nox_xxx_mapGenGetObjID_527940` | 18 | Retain |
| `nox_xxx_mapGenPlaceObj_5279B0` | 18 | Retain |
| `nox_xxx_mapGenMoveObject_527A10` | 33 | Retain |
| `nox_xxx_mapGenOrientObj_527C60` | 62 | Retain |
| `nox_xxx_mapGenFinishSpellbook_527DB0` | 16 | Retain |
| `sub_543680` | 25 | Retire |
| `sub_5437E0` | 110 | Retire |
| `sub_543BC0` | 25 | Retire |
| `nox_xxx_tile_543C50` | 108 | Retire |
| `nox_xxx_tileSubtile_544310` | 36 | Retain |
