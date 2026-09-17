# Object reports and recipient updates

## Scope and status

Native conversion qualified: **541 physical C lines removed**, leaving **55,577
lines / 82 production files**, zero reference C. Eleven functions across GAME4_1.c,
GAME4.c and GAME1.c are replaced by three Go files (325 lines), with fifteen C
entry/callback bridges retired. Frozen C baseline: **8b370caa**.

The scope covers type lookup caches, simple/phantom/monster/player object messages,
sprite state, recipient visibility/dirty masks, minimap update timing, polygon
sound levels and observer audio updates. Save/quest transition functions were
excluded after body review. All external production callers are Go wrappers;
selected internal helpers remain live. Moving the selected callers also permits
retirement of three preceding visibility C exports.

## Contracts being established

Reuse the actual registered server/types, C-owned sparse players, object update
storage, message queues, reliable allocator, spatial map and polygon owners.
Freeze packet bytes, recipients, queue kind, returns, dirty/known state, health
history, animation/RNG consumption, minimap scheduling and observer selection.
Test signed coordinate truncation, low words, code limits and all directions.
Keep independent byte/state assertions alongside repeated original-C captures.

Server object updates use Kind2; ordinary reports use Kind1; reliable reports
use the actual reliable-message pool. Queue growth/flush tests require the real
connection owner and valid timestamp prefix, not a mock delivery callback.
Tests must not accidentally invoke a flush on an incompletely initialized owner.

## Recovery

Development evidence is in build/port-object-reports. The primitive dispatcher
and initial encoding tests are a partial baseline, not a qualified conversion.
Do not freeze capture hashes before independent contracts pass and original C
repeats agree. Complete remaining branch/ownership tests before translating.

## Development findings

The first six roots pass 8,246 leaf cases. Health reports use the reliable-message
pool, whereas the corresponding object update uses Kind2 (or the player caller's
chosen Kind1 route). The fixture now observes both actual paths. Immobile sprite
messages use the extent-based high-bit code; the initial fixture expectation used
a dynamic code and was corrected before freeze. No production changes resulted.

Polygon lookups use a shared rotating ray endpoint. The owner now saves/restores
that counter along with polygon and vertex storage. Initial square-center probes
aligned exactly with polygon vertices and did not establish a robust interior
contract; off-diagonal interior probes pass. Preserve and explicitly capture the
vertex-aligned behavior separately rather than silently changing polygon geometry
inside this reporting batch.

## Original-C qualification

Default/server/highres each pass **400 affected roots / 23,241 leaf cases**,
without skips. All **82 groups / 22,549 records** match. These runs share the same
source fingerprints. Default took 207.3s, server 292.8s and highres 224.6s.

Final review added six first-polygon initialization cases as a separate test-only
supplement. They prove the existing C polygon update runs before audio filtering,
and that an observer following another player retains its own polygon state.
The supplement passes on all three targets; the focused repeat passes all
**18 new groups / 13,593 records**. Native qualification will check the combined
**83 groups / 22,555 records** and **401 roots / 23,247 leaf cases** together.
The manifest preserves the original broad phases and supplemental commands.

Production source is identical to qualified conversion `df8b3bb7`; all changes
in this baseline are porttest files and documentation. Reuse that commit's three
production builds, ABI audit, exact known full suite, gameplay/save-load/flat
references and static result rather than rebuilding identical production code.
Native qualification must run all production gates afresh. Its references are
visibility-effects-native, visibility-effects-save-native and
visibility-effects-flat-native; their logs, screenshots and saves survived disk
cleanup. Only identical copied original assets were removed.

## Decisions for review

- Move private `sub_417270` with its sole caller `sub_519710`. The scheduling
  contracts already cover real circular lists from zero through 65 entries,
  strict thresholds and frame wrap. Whole-source references are just its
  definition, declaration and selected caller; no bridge needs to survive.
- Remove the now-unnecessary C-to-Go audio callback bridge along with the ten
  selected entry points and the preceding special/out-of-sight/shadow bridges.
  Existing Go-facing wrappers can remain direct Go calls.
- Preserve 16-bit health-delta narrowing and report suppression for nonnegative
  deltas, plus history updates even when no delta message is sent. Preserve
  self-report history changes when a subsequent ordinary update queue is full.
- The first wall fixture lacked definition bit 1, which report rays specifically
  require with flags 69. Set that actual definition bit; retain a transparent
  variant and verify wall storage is unchanged. This corrected fixture setup,
  not production sight behavior.
- Ordinary queue capacity is checked through real queues, including failure at
  exact byte boundaries. Kind2 flushing itself remains in the existing queue
  implementation; this port changes its callers, and normal Kind2 packet routing
  is exercised directly plus in the headless integration scenarios.

The ignored native drafts/install script are consumed and stale. Do not reinstall
them over the qualified source. Frozen expectations were unchanged throughout.

## Native development

Baseline committed/pushed as **8b370caa**. Three Go files (325 lines) now replace
all eleven selected functions; working C count is **55,577 / 82 files**, zero
reference C, now qualified. Fifteen C symbols are retired, including the
private count helper, audio callback, and three preceding visibility bridges.
No retired symbol references remain, and static checks pass.

The first native focused run passes all 18 groups / 13,593 records and 12,493 leaf
cases unchanged. A stale C prototype in the old server preamble was found by the
source audit and removed before final gates. No production algorithm correction
was needed. The final three target sweeps pass 401 roots / 23,247 leaf cases and all 83 groups /
22,555 records, without skips. All 1,979 source fingerprints agree. Default took
303.4s, server 303.4s and highres 372.0s (concurrent runs, including compilation).
Production/integration qualification passes. All three builds and ABI audits pass;
the full suite exactly retains 1,553 failure entries (15 packages pass / three fail /
32 skip). Gameplay matches 41 frames, actual save/load seven, and flat rendering
fourteen plus exact map regeneration. All gates share 1,979 unchanged source
fingerprints. Production qualification took 381.9s. Client SHA-256:
`b13e16c90cc220af61c48b61ec57daa482f08cd54755e5a74f27d354f736ca2a`. Artifacts: build/port-object-reports/native-production.


Leaf-count correction recorded during reliable-queue qualification: recomputing
unique terminal test names from the original JSON logs gives 23,241 C and 23,247
native affected leaves. The earlier totals were inflated by 18 in the preceding
consumer corpus; the 12,493 new object-report leaves were already correct. All
original test names and capture records remain unchanged.
