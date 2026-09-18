# Item respawn and object ownership

C baseline follows session-entry native **449ae4c7**. Eight live routines cover
respawn-pool allocation/reset/free/add/remove and periodic update, owner-chain team
comparison, and carried-crown drop/attribution. Seven are the remaining bodies in
GAME3_3.c; the updater lives in server__system__server.c. Original scope is257 C
body lines. The caller audit finds one required C export, sub_4ED050 from GAME5.c;
the other interfaces can retire with direct Go callers. See
[item-respawn-selection.json](item-respawn-selection.json).

## Baseline contracts

Eight roots /461 tests including subtests produce eight captures /3,469 records.
Owners use the existing object/type allocator, report queues, map index, player
records, audio queue and actual fixed respawn pool. Independent assertions cover:

- Record layout and initialization, disabled insertion, duplicate records, return
  identity, both list links, first matching removal, empty/missing removal, reset,
  384-record exhaustion and reuse after removal/reset.
- Class masks, all five modifier words, ammunition byte order, position bits and
  direction narrowing. Raw pointers are normalized only to known fixture identities.
- All small acyclic owner forests, nil roots, identity, unassigned/equal/different
  teams and matches through either owner chain:2,592 independent comparisons.
- Scheduled/missing/destroyed/held objects, monster flags, Crown special treatment,
  definition eligibility, server rule2, mode bypasses, frame wrap and unchanged
  pending records. The updater intentionally disables new insertion until reset.
- Recreation just before/at/after the deadline and high-bit frames; real object
  creation, position/direction, weapon ammunition, old-monster delayed deletion,
  effect packets and deferred audio.
- Relocation at the five-second boundary and around50-unit distance, using actual
  quantized input positions; real map-index movement, health, recharge, ammunition,
  packet output and ordered audio positions.
- Crown filtering and traversal, failed/successful inventory drop, owner removal,
  position and unconditional post-attempt attribution, including zero/high-bit stamps.

## Correction for review

The C removal routine dereferenced head+4 before testing for an empty list. Add a
three-line empty-head guard before freezing. This follows the existing absent-item
no-op behavior; crown/team cleanup can request removal independently of registration.
The original invalid access is established by source inspection, not by deliberately
crashing the test process. Empty and repeated removals now have contracts.
Working C baseline is **31,697 physical lines /69 files /zero reference C** (+3).
Fresh production qualification is required for this correction.

Fixture development corrected the UseDataPtr adapter, a packet-snapshot type and
an item case that incorrectly reused a player object during effect delivery.
Recreation now uses a separate actual item allocation. No port expectations have
been regenerated to hide a translation mismatch. Production is still C.

Original assets/archive remain untouched. Raw runs/drafts live under ignored
build/port-item-respawn; drafts must not be installed before the baseline qualifies.

The seventh focused run and independent repeat pass all eight roots /461 tests;
all eight captures /3,469 records repeat byte-for-byte and are frozen in
[item-respawn-captures.json](item-respawn-captures.json). Static mapped-memory
checks pass. The first accumulated sweeps stopped before discovery because the
pattern file had two lines; it now uses one combined expression, with no source
change. Final default/server each pass344 roots /39,431 tests and214 captures
/55,337 records, including every unchanged parent capture. Remaining gates are
tracked in PORTING_STATE.md.

Verified deduplication of the three completed session-entry native scenarios
reclaimed1,660,044,319 bytes. Each run retains its restoration manifest and changed
files. `build/port-item-respawn/deduplicate-session-entry-native-assets.py` audit/apply
modes are consumed; its restore mode remains available. Original assets/archive
are unchanged.

## Qualified C baseline

All three targets pass344 roots /39,431 tests without skips. All214 captures
/55,337 records agree across targets, with every parent capture unchanged and all
new frozen hashes matching. Four gates share identical2,422-file source. Three
fresh production builds and ABI checks pass; the full asset suite has the exact
known1,553 failure entries and15 pass /3 fail /32 skipped packages. Gameplay,
save/load and map regeneration pass. See [C qualification](item-respawn-c-qualification.json).
All sessions are joined. The Go draft remains separate until this baseline is
committed and pushed.

## Native implementation

The C baseline is committed and pushed as **06ffe2c5**. Eight routines are native
Go, seven C interfaces retire, and only sub_4ED050 remains exported for its live
GAME5.c caller. Four private C globals move to Go. GAME3_3.c is removed entirely;
production C is **31,346 physical lines /68 files /zero reference C** (−351).

The native owner retains the fixed384-record allocator and60-byte layout, with
compile-time size/offset assertions. It keeps duplicate/list ordering, return
identity, mode and eligibility gates, unsigned frame wrap/comparison, relocation
precision and effect ordering. The same actual object/modifier/recharge, audio,
report and deletion owners handle downstream effects. Existing team/session
fixtures now use the native respawn owner; their expectations are unchanged.

The first native focus passed461 tests and all eight frozen captures without a
translation correction. Static checks pass. Accumulated default/server/highres each pass344 roots /39,431 tests without skips;
all214 captures /55,337 records match C and each other. Supplemental spell,
projectile, collision dispatch, inventory and object-state sweeps pass82 roots
/201 test entries per target against existing frozen expectations. These suites
also contain internally enumerated cases beyond the test-entry count.

All seven gates share unchanged2,423-file source. Three fresh production builds,
ABI/symbol checks, exact known full-suite outcomes, gameplay, save/load and map
regeneration pass. See [native qualification](item-respawn-native-qualification.json).
Every tool session is joined; source is editable.
Object-pool exhaustion during item recreation is not injected; fixed respawn-record
pool exhaustion is explicitly covered. No C algorithms are kept solely for tests.

Verified deduplication also reclaimed1,660,044,319 bytes from the three completed
item-respawn C scenarios. Per-run restoration manifests remain. The consumed
`deduplicate-item-respawn-c-assets.py` under build/port-item-respawn must not be
replayed in audit/apply mode; restoration remains available. All original assets
and the archive are unchanged.
