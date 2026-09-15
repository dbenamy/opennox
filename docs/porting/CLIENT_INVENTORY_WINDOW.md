# Client inventory windows, input and lifecycle

Status: **native Go conversion qualified**, following qualified C baseline
`0842b2d1`. Production now contains **82,632 C lines /
96 files / zero reference C**; 1,605 C lines have been removed.

## Connected scope

34 routines / 1,544 address-block C lines cover construction and image loading,
main drawing/animation, open/close/reset, scroll/mode events, equipment selection,
hover routing, trade selection and item drag/drop/cancellation. Removing the last
five routines from client__gui__guiinv.c also retires its remaining 61
include/extern lines: **1,605 physical C lines** in total.
Paper-doll composition and journal internals remain separate existing owners.

## Baseline and fixtures

**2,120 results / 29 frozen groups**, plus independent contracts, exercise actual
GUI, renderer, player, object-factory, inventory, spell and modifier owners.
All 34 window tests and 18 existing display tests passed together in 45.475s.
Every window capture was byte-identical to the preceding separate-process run.
Cumulative frozen coverage is **455,479 results / 1,185 groups**; repeated runs
are not counted again.

Reuse inventoryDisplayOwner and inventoryTransactionOwner. The display constructor
accepts optional extra type names while preserving its original defaults. The
window fixture adds hidden-row and quest types, bounded image requests, an actual
ImageRef animation, embedded hit rectangles, real globals and trade-grid storage.
Newly authored minimal identify and quantity resources go through the real parser
and actual widgets. Original asset geometry is exercised by headless gameplay.

Coverage includes:

- Constructor children/hidden items, missing resources, real allocation exhaustion,
  event admission, state bytes and all 40,000 valid integer paper-doll points,
  three outside points, and equipment priority including linked shield entries.
- Open/close/identify/repair transitions, visible-dialog and player-animation gates,
  scroll/mode events, toggle and tooltip events, rectangle edges, slider routing,
  map-button behavior and click/drag distance and timing boundaries.
- Opening/closing animation, journal layout, spell/modifier/quest icons and hover,
  paper-doll base-image selection and material colors through the real renderer.
- Stack pickup/restoration order, occupied/empty/outside drops, actual quantity
  dialog acceptance/cancellation, blocked/missing drop requests, trade/repair
  request codes, alternate-weapon events and full-source cancellation fallback.
- Reset deletion order across all 84 cells and idempotence. Client-only world
  selection uses the real spatial index/finder at the 75-unit distance boundary;
  its messages are captured. Other tests also run in dedicated-server builds.

Only declared pointer fields are normalized. Captures retain pixels, text, GUI
and widget records, inventory cells, object lifetimes, allocation/deletion events,
requests, sounds, cursor state, pool counts and shared-state changes.

## C prerequisite

sub_464BD0 assigned only LOBYTE(v14) from a bool rectangle predicate, then tested
the full uninitialized int. The prerequisite assigns the full result. Independent
contracts require false trade hits to preserve the previous selection and true
hits to select the actual grid drawable and apply its stack code. C LOC is
unchanged. Undefined high bytes are not a compatibility oracle.

The geometry audit found no uncovered valid paper-doll point, even with slot0
empty. The equipment index therefore needs no speculative correction.

## C baseline qualification

Affected tests passed: default **313 / 150.307s**, server **311 / 230.942s**,
highres **313 / 154.457s**. Every selected test started and completed.
Production builds passed in 59.895s, 8.664s and 55.593s respectively; all three
verify ELF32/i386/SSE2/CGO, all 34 original interfaces and no test helpers.
Full assets match exactly: **1,553 known failures**, 15 pass / 3 fail / 32 skip
packages. All **1,574 source fingerprints** stayed unchanged during final checks.

[inventory-window.yaml](inventory-window.yaml) extends the four-screen display
scenario with inventory close/reopen, a shirt drag/reorder, scrolling and retained
state after reopening. The original four screens are byte-identical to the prior
qualified reference. All nine screens pass a fresh comparison with reference
replacement disabled, seeded localization, Xvfb 1280×960 and null audio.
The initial game run exited successfully; only its post-run deduplication helper
rejected the new run name. Its explicit list was extended, then verified asset
deduplication and the fresh comparison completed. Original assets stayed unchanged.

Evidence is under build/port-client-inventory-window: c-qualification.json,
baseline-hashes.json, source fingerprints, binary verification, full-suite
comparison and scenario-prefix-comparison.json. Gameplay references/results are
under build/baseline/runs/client-inventory-window-c-{capture,repeat}. Successful
runs have verified asset-restoration manifests. Older development captures have
lossless compression manifests.

## Unexplained identification mismatch

The first broader default run completed 313 tests in 158.645s with one mismatch
in the existing display-identification test. Got
`7f44a2d7d1a3e37bce6662cef939c112601768bff9ddb90b093b812bd6af6c56`, expected
`1e16c910126c894d88e21b5f0a39c65154f7dd142fb674bdcb48f93627c787b7`.
No complete raw capture was requested for that first failure. An unchanged-source
replay of the same 49 preceding tests plus identification passed in 23.929s;
identification alone passed in 8.737s. The repeated broader run passed as above.

No cause is established; do not call this fixed or pre-existing. The display
capture helper now automatically preserves a private, uniquely named full
mismatch under build/port-failures when explicit capture was not requested.
Expectations and production code stayed unchanged. The production builds were
reused across that one-file diagnostic-only fixture change, recorded in
capture-diagnostic-delta.json. Inspect complete captures if it recurs. This is
separate from the earlier hallway mismatch.

## Cancellation ownership follow-up

sub_467CD0 clears a non-equipped drag pointer without deleting its temporary
drawable. Placement/pickup copies into a different inventory drawable. Each of
four simple cancellation captures retains five live case drawables: three hidden
items, the inventory item and one extra drag copy. Pool deltas and the full event
ledger record this; no normalization hides it.

Preserve this deterministic behavior during the conversion, then correct it in a
separate Go cleanup chunk with repeated-cancellation and equipped-item lifetime
contracts. Local cancel-ownership-audit.json records the source/capture evidence.

## Development corrections and resumption

Development runs c through o expanded from four to 34 tests. Fixture corrections
included a missing C declaration, a pointer/integer declaration mismatch, a missing
helper argument, duplicate initialization of an already prepared weapon registry,
and an image-handle conversion. World-selection fixtures remove their temporary
index membership after executing the real selection and before the shared
unlinked-object snapshot. The first quantity test incorrectly expected the dragged
code first; actual restoration appends it to the stack. Removal/restoration order
now has independent checks. These fixture corrections did not change production
logic. All source-reading jobs, including failures, were joined before edits.

Ignored native-*.go.stage drafts are stale after integration and compile
corrections. Actual source and committed expectations take precedence. The
integration script is already applied; never rerun it.

## Native conversion

Baseline **0842b2d1 is committed and pushed**. Seven Go owners now replace all
34 routines. The obsolete inventory translation unit is deleted: **1,605 C lines
removed**, leaving **82,632 / 96 files / zero reference C** in the working tree.
Eleven interfaces remain for C callers or tooltip/quantity callbacks; 23 internal
functions retire. The source audit finds no remaining C reference to a retired
function. Go callers use Go directly. Binary qualification confirms those interfaces in all three targets.

Native-a/b failed discovery on missing headers; native-c exposed C/Go pointer,
UTF-16 and integer conversions. Native-d compiled and completed all 52 selected
window/display tests, with 51 passing (180.993s). Its sole mismatch had exactly
two fields: fallback cancellation wrote the equipped flag to a drawable instead
of the cell. The pointer expression was corrected and an independent equipped-cell
assertion added. No frozen expectation changed. Native-e passed all 52
tests in 180.632s; all 29 raw capture hashes match the frozen C baseline. The
accumulated default suite passed all 874 selected root tests in 473.503s.
All failed jobs were joined before editing; ignored drafts are now stale.

Native qualification passed: accumulated default **874 / 473.503s**,
affected server **311 / 229.165s**, highres
**313 / 152.932s**. Every selected root test started
and finished. All three production binaries verify ELF32/i386/SSE2/CGO, 11 retained
C interfaces, 23 retired interfaces and no test helpers. Full assets retain exactly
**1,553 known failures**, 15 pass / 3 fail / 32 skip packages. The nine-screen
headless scenario matches the committed C baseline's reference, with replacement
disabled, in **51.206s**. All **1,580 source fingerprints** stayed
unchanged throughout qualification. No hallway or identification mismatch recurred.

Evidence: native-qualification.json, native-focused-capture-verification.json,
native-{default,server,highres}-result.json, native-binary-verification.json,
native-source-interface-audit.json and full-suite-comparison.json under
build/port-client-inventory-window. The fresh gameplay run is
build/baseline/runs/client-inventory-window-port. Completed development captures
have verified lossless compression manifests. Original assets remain unchanged.

The deterministic cancellation leak remains deliberately unchanged in this
conversion. Its next, separate cleanup will first demonstrate the leak with
independent lifetime contracts, then verify precise ownership changes without
changing unrelated gameplay expectations. C baseline and original hashes remain
recoverable from Git.

The cancellation ownership follow-up is now separately qualified; see
[CLIENT_INVENTORY_CANCEL.md](CLIENT_INVENTORY_CANCEL.md). The conversion's
unchanged C expectations and leak behavior remain recoverable at 24f3f67e; only
the six explicitly reviewed ownership results change in the follow-up.
