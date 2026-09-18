# Server map/round orchestration

## Scope

The C baseline follows qualified item-respawn native **5822c9c6**. The batch
contains **12 C bodies / 494 body lines** in `server__system__server.c`: ten live
operations and two bodies whose combined loop has no observable effect. The
statistics initializer `sub_426060` stays with its statistics/serialization owners
for a separate connected batch. The tracked selection is
[server-orchestration-selection.json](server-orchestration-selection.json).

All selected outside callers are Go wrappers. The replay loop `sub_4E76C0` calls
an empty `nullsub_25`, computes a read-only object checksum into a discarded local,
traverses objects, and returns zero to a wrapper which discards it. The native conversion removes that
loop/caller and its checksum C export. The existing object-state fixture dispatch
now calls the qualified Go checksum directly with its frozen expectations intact.

## Original-C contracts

**16 roots / 1,819 test entries** pass. All **16 captures / 1,923 records** match
between separate processes (`eighteenth` and `repeat` under
`build/port-server-systems`). Expectations are frozen in the tests and indexed in
[server-orchestration-captures.json](server-orchestration-captures.json).

| Behavior | Records | Independent checks / actual owners |
| --- | ---: | --- |
| Idle timeout | 100 | Strict boundary, uint64 wrap, conditional clock reads, unchanged words |
| Map-load state | 6 | Full-width scalar patterns |
| Difficulty timing | 512 | Participation combinations, frame/FPS products, changed-value scaling, controlled nonzero balance data |
| Player reset | 40 | Exact byte writes, preserved surrounding bytes, glyph inventory count, cooperative bypass |
| Secret walls | 725 | All phases, timer/progress wrapping, two-wall ordering, actual map index, activation queue and audio |
| Restore references | 128 | Main/pending/missile/inventory lookup; missing and destroyed objects; pointer versus net-code fields; distinct early-stop contracts |
| Restore scheduling | 84 | Class precedence, synchronization masks, actual update-list membership and links |
| Restore cleanup | 72 | Protected objects, cooperative/online rules, glyphs, owned pixies, actual deletion queue/order/frame |
| Restore without host | 8 | Inactive or unitless local player leaves state untouched |
| Restore integration | 6 | Saved positions/IDs, AI references, script timer references, pending ownership, creature minimap/reports |
| Carried flags | 24 | Actual inventory removal, pending objects, positions/radius, RNG consumption, flag timestamps and packets |
| Round-start crowns | 72 | Actual team membership, unitless-player rejection, occupied crowns, pickup success/failure, RNG, ownership and buffs |
| Rewards | 105 | Marker/chest replacement, selected/unselected markers, Ankh, bows/quivers, missing quiver, empty rewards, stage boundaries/wrap, nonzero gold, actual factory/inventory/deletion owners |
| Player transition | 32 | Quest-stat transfer, state failure, loadout policy, controls, ability cooldowns, HP/mana, observer/camera behavior |
| Fresh loadout | 8 | Actual gear factory/pickup/ownership; map-state bit selection including high-bit values |
| Open shop transition | 1 | Real session detachment and quest-shop cache ownership |

Difficulty timing uses an empty world object list; downstream health scaling has
separate frozen contracts. Positive FPS is the runtime precondition. The fixture
checks the retained-loadout and actual fresh-warrior-loadout paths; downstream
class-specific item selection and trade policy retain their existing suites.

## Behavior to preserve / review separately

The wall loop carries the preceding activation-scan flag into an unknown/zero
phase on the next wall. The two-wall fixture documents this existing ordering.
Preserve it for this port; changing it is a separate gameplay decision.

Saved monster references have two different failure contracts: the pointer list
writes zero before stopping, while the net-code list stops without overwriting
the missing ID. Both leave the unvisited tail intact.

Fixture development corrected a participation-count assumption, a 386 integer
constant, two field-width assumptions, and the fresh-warrior shirt expectation.
Integration setup now initializes actual ability/type owners. Reward teardown
must detach a chest's inventory before freeing individually tracked items because
`FreeObject` otherwise frees the contents recursively. These were fixture issues;
production C is unchanged.

## Original-C qualification and recovery

All three targets pass **427 roots / 41,317 test entries** without skips. All
**245 captures / 78,253 records** match across targets, including the frozen C
expectations. The gates share identical **2,433-file source**; static checks pass.
The broader pattern adds inventory, object-state/checksum, reward, player-control
and monster-control regressions to the previous item-respawn target set. See
[server-orchestration-c-qualification.json](server-orchestration-c-qualification.json).
This qualification records the original-C baseline before conversion.

Production source is identical to qualified native **5822c9c6**: only ten files
with the `porttest` build constraint were added. Reuse its three production/ABI
builds, exact known full-suite result and headless gameplay/save-load/flat-map
scenarios for this test-only C baseline. The source comparison and referenced
report hash are recorded in
[server-orchestration-production-reuse.json](server-orchestration-production-reuse.json).
The native conversion uses fresh production qualification below.

Baseline **ec5ffeb4** contains **31,346 physical lines / 68 production files /
zero reference C** and was committed and pushed before translation.

Completed item-respawn native scenario asset copies were hash-checked against the
original assets and deduplicated, reclaiming **1,660,044,319 bytes**. Each run has a
`deduplicated-assets.json` restoration manifest. The script under
`build/port-server-systems/deduplicate-item-respawn-native-assets.py` has consumed
its audit/apply passes; only its restore mode is reusable. Original assets,
archive, changed maps, saves, logs, screenshots and production binaries remain.

## Native conversion and qualification

Ten live operations now use Go directly. Twelve selected C bodies retire along
with the previously native checksum's unused C bridge: **13 C function symbols**
and **two private C globals** removed, no new exports. A small address-only shim
retains the actual chest-initializer callback identity. Existing C-owned wall and
object records continue through their real owners. The statistics initializer
stays with its reporting subsystem for a separate batch.

Review preserved byte/word overflow, timer and stage wrapping, strict timeout
comparison and conditional clock reads, distinct missing-reference behavior,
integer audio versus fractional activation coordinates, callback identity and
list mutation order. An initial compile caught a point-value/pointer API mismatch;
it was corrected before qualification. No behavioral mismatch or golden change
was needed. Unknown wall phases retain the existing ordering documented above.

**34 focused roots / 1,837 test entries** pass, including the existing checksum
contracts. All **16 new captures / 1,923 records** match the frozen C baseline.
Each default/server/highres gate passes **427 roots / 41,317 entries** without
skips; all **245 captures / 78,253 records** match C and each other. Final source
identity is checked across target and production gates. Static memory checks,
three fresh production builds/ABI checks, the exact known broader-suite failure
set, and fresh headless gameplay, save/load and flat-map regeneration all pass.
See [server-orchestration-native-qualification.json](server-orchestration-native-qualification.json)
and [server-orchestration-native-batch.json](server-orchestration-native-batch.json).

Final C: **30,819 physical lines / 68 files / zero reference C**, **−527** from
the baseline. This includes removal of stale address comments beside retired
bodies. Native evidence is under `build/port-server-systems/native-*`.

Completed monster-control native scenario asset copies were also hash-checked
and deduplicated, reclaiming **1,660,044,319 bytes** with per-run restoration
manifests. Both cleanup scripts' audit/apply passes are consumed; only restore
mode is reusable. Original assets and archive are unchanged.
