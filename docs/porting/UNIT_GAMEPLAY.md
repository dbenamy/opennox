# Unit gameplay helpers

## Scope and status

Qualified parent `ae81334f` is pushed. The accepted conversion boundary covers20 C bodies (405 body lines) in seven
files: experience/rewards, AI orders, hurt/damage timers and action metadata,
object use/update callbacks, dialogue and localized item/NPC descriptions.
The two spell creation/start entries from the initial22-body candidate are
deferred to the next batch with their duration/factory fixtures. This keeps the
current unit-helper qualification connected; it is a reversible scope decision.
The tracked selection is authoritative; the caller inventory retains the initial
22-body audit for provenance. Production is
still C; remaining production C is 10,820 lines / 47 files / zero reference C.
The baseline is in progress, not qualified or frozen.

Read-only selection/caller inventories and diagnostic logs are under
`build/port-unit-gameplay`. Existing controls, creature transfer, lifecycle,
spell, trade and death fixtures supply the actual owners. Batch qualification
will include affected default/server/highres contracts, original-C repeat
captures and independent expectations. Reuse parent production evidence only
if production source remains identical; run fresh production/ABI, full-suite
known-failure comparison and relevant headless scenarios after conversion.

## Findings and decisions to review

- The shipped 72-entry AI-name table has an empty entry at index 39, pointing
  to zero-initialized storage at 5D4594+2488508. The other 71 entries point into
  shipped constant data. Existing creature-transfer fixtures instead fill every
  entry with the Go enum label, making index 39 `ai.Action(39)`. New contracts
  relocate the actual shipped strings; existing serializer captures stay intact.
  Literal address/symbol search found the relocation as the only production
  reference to the empty storage. Preserve empty lookup, nil versus empty name,
  restricted fallback 38, general fallback 0, case sensitivity and first match.
- C `sub_534020` is an unreachable duplicate: only its definition and prototype
  refer to the C body. The root Go function with the same name has a Go caller
  in unit_ai_path.go. Searches for the symbolic, hexadecimal and decimal address
  found no C registration or address consumer. Remove the duplicate with the
  batch rather than creating artificial tests for it.
- The hurt-sound helper's char return can contain the low byte of a sound-table
  pointer. Its sole production caller ignores the return. Contracts should
  compare timer, RNG and emitted audio, not incidental pointer bytes.
- XP uses the shipped float32 coefficient 0.01, and C can retain intermediate
  arithmetic in x87 registers. Check returned award, stored XP, protected values,
  level/stat mutations and notifications; inspect compiled arithmetic if needed.
- Preserve required callback ABIs for actual use/update and duration registrations.
  Both AI debug formatter callers are in Go C preambles; migrate them together.

## Initial evidence

`actions-first.log`: 754 original-C cases pass. They cover all shipped names,
all valid IDs, signed extremes/out-of-range IDs, restricted ID selection, unknown,
case and embedded-NUL inputs, and controlled duplicates in both search ranges.
Initial capture SHA-256:
`85d1149befc354744bd00ac5ae551cb0b22c6913d6b6d87c75d4e97a3ea4fc13`.
This is a first capture, not a frozen baseline. Static-memory preflight passes.

The action capture repeats identically in a second process. The initial 1,200 XP
contracts pass (`xp-first.log`, SHA-256
`8252317eec3326e2451dbb02d0209ebca4d196a705b088053b79881496181a1a`).
They cover three XP helpers, adjacent threshold floats, fractional awards,
save-state gates and single-step level advancement. These are initial results;
reward-owner ancestry and cooperative branches remain to cover.

Compiled-parent audit (`xp-c-disassembly.txt`) identifies an additional rounding
boundary: subtraction is spilled to float32 before the coefficient lookup. The
award is then calculated wide, copied to float32 for return/protection, and the
wide value is added to the old XP before the final float32 store. Eight literal
cases in TestUnitGameplayExperienceRounding distinguish these alternatives.
The initial cases did not discriminate the subtraction spill; their expectations
were corrected before freezing. No production change or golden regeneration.

Expanded timer/order/dialogue contracts are installed. The first expanded build
failed on a missing test-only declaration for the private dialogue completion
function; its actual signature is now declared in the fixture preamble. The
second run is pending. Static-memory check passes with the expanded fixtures.
Dialogue fixtures invoke the actual VM's registered Go callbacks and observe
caller/trigger, response/partner state and unfreeze ordering.

Fourth build succeeds. Split test groups establish: action names754, AI debug400,
orders36, dialogue1,120, XP1,200 and literal rounding8 pass. See fourth-text/tests.log
and fourth-xp-rest.log. Dialogue initially changed the synthetic player's class
for a rejection case without restoring it before the outer typed snapshot; fixed.
The timer fixture similarly accessed a typed monster view after changing class;
fixed. Reward fixture omitted the coefficient; now both XP and reward fixtures
load the shipped coefficient. The temporary update fixture initially placed its
target in init data instead of collision data at object+700; corrected. These
were fixture defects, not production corrections. All failed logs remain.

Final focused baseline: **14 roots /8,464 cases**, all pass twice in separate
processes with identical captures. Thirteen captures were frozen after their
repeats; the final20 read-message boundary cases also pass/repeat and are now
frozen. Tracked unit-gameplay-batch.json requires all14 hashes. Evidence:
`focused-c-{a,b}`, `focused-c-hashes.json`, `seventh-read{,-repeat}`.

Default/server/highres qualification passes **238 roots and242 captures each**,
with no skips and identical source fingerprints/capture inventories. The drivers
enforced14 new hashes plus inherited literal expectations; the full242-capture
inventory was compared and frozen after the runs. The tracked manifest now
requires that full inventory for future reruns. Static-current passes.

Production identity audit confirms all23 changed source files are porttest-only,
all three parent binary hashes match, and all22 audited symbols (20 selected
bodies, the debug bridge, the retained spell-control owner) exist in each binary.
Parent `ae81334f` supplies qualified production/ABI, exact1,553 known full-suite
failures (15 pass/3 fail/32 skip packages), and gameplay/save-load comparisons.
No new production/scenario runs are claimed for this test-only baseline.
See unit-gameplay-c-qualification.json and `c-affected-{default,server,highres}`.
All test/build/scenario jobs are joined.

Remaining fixture corrections, now resolved before freezing:
- The string manager folds localization keys case-insensitively; NPC key
  construction itself preserves case. Contracts now test both properties.
- Timer/update specs must give the outer clock owner the same frame as the
  called function's fixture state.
- ReadUse borrows a guarded UseData allocation, restoring the original pointer
  directly. It must not transfer/free that interior pointer through SetPtr.
- Inventory setup installs the final wall owner after combat setup; blocking
  read cases now configure Inventory.WallMode and verify actual visibility.

The cooperative XP presentation cases run the real pause/presentation routine,
with owned globals, deterministic ticks and successful/missing actual factory
objects. The inventory fixture still observes the existing player-state boundary;
these cases do not claim to retest the whole player-state implementation.

Conversion planning found GAME4.c also owns nox_cheat_charmall. Keep its minimal
C definition for the following spell batch, which owns its Go consumers; remove
its obsolete order body and surrounding dead declarations here. Do not delete
that remaining live owner just because the last C algorithm in its file moved.
Ignored native-*-draft.go files are NOT installed or qualified.

Completed inventory scenario copies were verified against the original assets
and deduplicated:1,112,747,701 bytes reclaimed. The script's apply mode is
CONSUMED; restoration manifests remain in both completed run directories.

Completed C capture deduplication reclaimed1,969,117,384 bytes across490 identical
files. All paths and hashes remain; contents are immutable. Audit/apply manifests
are under build/port-unit-gameplay; apply mode is CONSUMED.
