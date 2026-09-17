# Map polygon lifecycle, lookup, events and serialization

## Current scope

Parent server-configuration conversion **5ecbdfc2** is qualified and pushed.
The next C baseline covers **33 reachable functions / 1,000 original C body
lines**: GAME1_1 420C40..422140 and GAME1_2 map serialization 428CD0. Whole-source
reference audit plus the selected internal call graph finds all 33 reachable;
private helpers are not orphans merely because they have no outside callers.

Production remains **45,473 C lines / 74 files / zero reference C**. The C baseline is qualified with frozen expectations. Production is unchanged;
the native implementation is still an ignored draft. Read-only scope/audit/plan and
local logs are under build/port-polygons. The combined extraction includes the
serializer; the older core-only extractor does not.

## Baseline design

Own full mapped vertex/polygon/selected-vertex arrays and control/default bytes,
separately from extracted live counters and the remap-list head. Reset and free
only the allocations created by the fixture, then restore the previous pointer-
bearing state. Reuse actual C allocation/removal, cryptfile owners, player/monster
factories, renderer/ambient-color state and real script/report/audio queues.
Seed nonzero default names, RGB values and level bytes so empty blobs do not hide
behavior. Use the existing minimap/object-report fixtures for dependent callers.

Cover sparse allocation/iteration and reuse, remapping, defaults and exact byte
mutations, nearest vertices and ties, selected polygon construction, integer bounds,
point/edge lookup with cache and phase state, player/monster enter/leave scripts,
secret awards/messages/audio, ambient-color updates, and serialization versions
1..4 with exact wire bytes and round trips. Inspect float widths, signedness,
sentinel handling and ownership before translating. Capture addresses only as
owned identities or offsets; no raw pointers in frozen output.

Independent contracts precede freezing. A nil-actor child test checks the explicit
null guards; the monster routine currently loads update data before its guard.
Get-by-index returns raw record addresses except a dedicated sentinel, so fixture
lookup cases must stay in owned record bounds plus that sentinel. Do not invent a
new bounds contract from a guessed getter convention.

Reuse the parent's qualified production result only if polygon baseline work
changes tests/docs alone. Any justified production prerequisite requires fresh
production qualification. Repeat/freeze and commit/push C before conversion.

## Disk cleanup

Before this batch, twelve completed panel/configuration run copies were checked
against original assets. Removing 19,648 byte-identical files reclaimed 6.184 GiB;
per-run deduplicated-assets.json manifests retain hashes and restoration metadata.
Modified maps, saves, captures, binaries and original assets/archive are preserved.
The ignored helper accepts these new run names. Never rerun its deletion mode on
already-deduplicated runs; use its recorded restore command if a full run data copy
is needed again. About 16 GiB was free immediately after cleanup.

Initial focused run passes three roots (0.106s), including nil-player and
nil-monster calls on the qualified compiler. No C null-guard correction was needed.
Get-by-index cases now use owned IDs and the actual sentinel. Expanded contracts
cover nearest-vertex ties/inactive slots/thresholds and float boundaries, runtime
versus editor selection, 0/1/2/3/4/7-vertex construction, bounds, real metadata
allocation and newest-first remapping. Focused-2 passed six roots (0.087s). Final qualification follows below.


## Geometry compatibility finding

The legacy containment test rotates a finite ray among world corners (0 or 5888
on each axis). Its segment helper includes both endpoints, so a ray through a
polygon vertex may count both adjacent edges and report outside for an interior
point. The initial square-at-origin contract exposed this. Preserve the behavior
in this port: changing intersection rules could alter map-region gameplay and
belongs in a separate geometry change. Ordinary interior/exterior contracts use
translated shapes; a separate independent four-phase contract records the shared-
vertex limitation, and the dense origin/boundary matrix retains those cases.
No production C fix or expectation regeneration was used.


## Qualified C baseline

All 16 focused roots pass (0.239s). Eleven polygon captures match separate
process runs before freezing; independent contracts additionally cover borrowed
reset ownership, cache initialization, nil actors, ray degeneracy and ordinary
containment. Static memory checks pass. The full fixtures exercise real C storage,
real script callbacks and event order, private and other-player notification bytes,
secret statistics/audio and repeat suppression, renderer light color, serialization
versions 1..4, vertex-ID remapping and exact modern wire round trips.

Broader default/server/highres: **269/268/269 roots**, **34,988/34,987/34,988
passing tests including subtests**, no skips. **87 captures / 51,433 records** are
identical across targets. Durations: **143.87/225.80/157.58s**. All runs used one
unchanged **2,175-file source manifest**; all sessions joined. Production-source
identity with parent 5ecbdfc2 was checked after excluding test files, so reuse its
qualified builds/ABI, exact known asset-suite failures and gameplay/save-load/flat
map results. No production C correction was made in this baseline.

Selection: map-polygons-tests.txt includes all polygon/map fixtures (excluding the
isolated prerequisite diagnostic), minimap, object reports, quest runtime, world
collisions, AI main/monster state, object rendering, player controls and numeric
conversion. These cover direct lookup/audio/minimap callers, shared polygon state,
map serialization and the real actor/script/render fixture owners. The server has
one fewer root because the corresponding existing test is client-only.
Artifacts: build/port-polygons/c-{default,server,highres}, c-audit.json,
focused-11.log and static-final.log. Full dependent capture hashes are pinned in
map-polygons-batch.json; all 33 functions are listed in map-polygons-scope.json.

## Additional compatibility notes for review

- A remote player's transition out of all polygons retains the previous polygon;
  the host's zero client-cache path clears the server cache without resetting its
  level. These existing branches are explicit contracts, preserved for this port.
- Nearest-vertex squared distance is stored in a signed integer before comparison
  with a float threshold. Preserve truncation and intermediate float widths.
- Record storage uses C allocations and a flag indicating borrowed ownership.
  Reset detaches borrowed arrays without freeing them. The fixture returns detached
  allocations to its real cleanup afterward. Native code must retain compatibility
  with existing C-heap owners, not introduce incompatible allocator bookkeeping.
- Editor callback names use real 256-byte metadata; runtime callbacks use the real
  script VM. Serialized secret counts mutate both current and previous counters,
  and the fixture owns/restores both.

The ignored effects fixture draft is **consumed**; never copy it over the frozen
fixture. Native files under build/port-polygons/native-draft are unfinished and
have not been installed or qualified.
