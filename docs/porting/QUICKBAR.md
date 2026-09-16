# Quickbar UI port

## Qualified native result

The complete quickbar owner is native: **87 functions**, **2,494 C lines removed**,
**64,317 C lines / 86 files / zero reference C** remaining. Implementation recovery
commit f7c75139 is followed by final fixture ownership and qualification docs.
All **6,983 records / 49 groups** pass in default/server/highres without changes
to frozen hashes. Final fixture checks take 83.035 / 82.414 / 81.278 driver seconds.
Three accumulated sweeps, three production builds/ABI audits, exact known-suite
comparison and both fresh gameplay replays pass. Source identity and the narrow
fixture-only difference between sweeps are documented at the end of this report.

Review items: the saved-row C prerequisite correction, preserved byte-only flag
writes and stale destructor words, constructor callback/row behavior, unsigned
arithmetic boundaries, and explicit raw-blob versus named-global fixture ownership.
No unresolved question blocks the next batch. Sections below preserve the baseline
and development history; their provisional statuses describe those earlier stages.

## Original C baseline and scope

The original-C baseline is fully qualified. All **2,242 quickbar records in
22 groups** and **4,694 book records in 24 groups** match in default, an independent
repeat, server and highres. All four phase source manifests match the same
**1,801 source files**, unchanged during qualification.

- Default/repeat: 148.864 seconds; server: 161.925 seconds (concurrent jobs).
- Highres: 88.477 seconds; C gameplay build/capture/repeat: 187.324 seconds
  (concurrent jobs).
- Fresh C binary SHA-256:
  87163fc4a80aae957e95f6e11388acb20fc4beb42ea41f84a95a56308cf73cdf.
- Both gameplay runs match all **22 frames**, remove 51 copied maps and
  regenerate the warrior map exactly. See quickbar-replay.json.
- C baseline readers are joined; baseline commit **b487b8dd** is pushed.

The baseline follows qualified spellbook conversion **41b8abfb**.
Production C remains **66,811 lines / 87 files / zero reference C**; the one-line
saved-row prerequisite correction is described below.
The selected connected scope is **87 functions / 2,457 physical lines**:
GAME2.c from 0045D9D0 to its end, GAME2_1.c from 00460D40 up to
00461460, and client__gui__guispell.c. The latter also has file declarations
outside that count. Selection includes slot storage, spell/ability activation,
drag-and-drop, expanded rows, trap rows, initialization/destruction, rendering,
reward presentation and saved slot data.

The spellbook, GUI, input, rendering, metadata, netlist and particle owners provide
the reusable baseline fixtures. Extend them with direct quickbar contracts before
translating. The baseline is complete; the native implementation is now installed
and undergoing qualification.

## Review notes from the audit

- Preserve full 32-bit slot flag words for swaps, but preserve upper three bytes
  when a setter, saved-data operation or row copy writes only the low flag byte.
- Expanding saves the selected row at 1047912; closing reads 1047908. This apparent
  mismatch needs baseline coverage and a caller/state audit before deciding whether
  a correction is warranted. It is not silently corrected by translation.
- Repeated row selection has a timeout and saturation path that writes selected
  row 4 without refreshing the row pointer. Cover the reachable timed sequence
  and explicit boundary state.
- Activation records must exercise the actual client message list, cursor guards,
  timeout owner and frame counter. No duplicated C test algorithms.
- Keep original book gameplay frames and add meaningful quickbar interactions for
  a fresh original-C reference; mere startup frames do not establish slot behavior.

## Native implementation checkpoint

The installed conversion removes **2,494 C lines**, leaving **64,317 / 86 files /
zero reference C** provisionally. Typed records retain the shared 256-byte layout;
actual GUI, renderer, input, save-file, message-list and particle owners are reused.
Go callers now invoke the native helpers directly. Nineteen C interfaces remain
for actual callers and tooltip callbacks; 68 private interfaces are retired.

Preserved details for review: hidden rows reuse nugget 1 and its direction index;
only the main row increments that counter in the original constructor. Both
caster classes check the same availability field when initially showing the
trap/bomber controls. Save-name parsing stops at embedded NUL like the original
C string conversion, while still consuming the complete counted field. These
are preservation choices, not new behavior fixes.

The initial compile check passed before integration. The first integrated run
(native-focused-01, 169.342 seconds) matched 45 of 46 groups. Its 18 differing
GUI-frame records contained only missing empty hotkey-label text calls: the live
C string adapter returns a non-null interned pointer even for an empty string.
The Go renderer now preserves those calls and their text-color side effect.
Goldens are unchanged. The corrected focused run (native-focused-02) passes all
46 roots and all 6,936 hashes, without skips, in 185.283 seconds. All 1,811 source
fingerprints remain unchanged. Accumulated default/server milestone runs are
underway; highres and production/gameplay gates still remain. Do not rerun ignored
`native_*.draft` or `integrate.py` over the current source.

## Remaining gates

Native implementation and comparison; production builds/interfaces, known
asset-suite comparison, fresh gameplay replays and the next full accumulated UI
milestone. Commit/push each qualified chunk and continue.

## Baseline development checkpoint

The first candidate had four direct test roots covering 390 slot-search cases,
125 full-word swaps, 270 actual activation-message cases and 48 pending-activation
transition records. Expectations are not frozen yet. The focused build is the only
active source reader; additional test drafts remain under ignored
build/port-quickbar until it joins.

The external-caller audit is preserved in build/port-quickbar/audit.json. The
remaining C callers cluster in startup/character selection, game input,
save handling and client ability updates; many private quickbar helpers can move
to direct Go calls together. Do not retain extra exports just for test dispatch.

## Original-C prerequisite evidence

All seven direct roots pass without skips in default porttest configuration
(build/port-quickbar/c-rows, 28.800 seconds). The initial four-root candidate
also passed (114.012 seconds including the initial C build). The tracked
quickbar-original-captures.json records the seven original hashes and record
counts. These are prerequisite captures, not completed batch qualification.

The complete expanded-row test confirms that starting on any of rows 1–4 and
closing returns to row 0, although expansion saved the actual row at 1047912.
The only use of 1047908 is this read and two ability-loop end-address comparisons;
there is no writer. Correct the closing read to 1047912 before freezing the final
conversion baseline. Preserve these original hashes and the Git revision as evidence.
The correction is reversible and needs no user decision; it restores the row
that the same operation already saved.

The original seven-root baseline is committed as **fdc5048e** (949 records).
The one-line row-read correction is being checked against all seven roots.
Only the expanded-row expectation changes; the six other original hashes remain
enforced. Slot-event coverage is drafted outside source while this reader runs.

## Saved-row correction

Corrected the closing read from 1047908 to 1047912. All seven roots pass
(default porttest, no skips, 115.618 seconds). The six unrelated frozen hashes
remain identical. A field-level comparison of the ten expanded-row records
changes only the closed selected-row byte, its row pointer and the direction
indicators for that restored row. Original captures remain recoverable from
fdc5048e and quickbar-original-captures.json; quickbar-captures.json has the
corrected baseline. C LOC is unchanged at 66,811.

This prerequisite has focused validation. The connected batch's server/highres,
production, gameplay and accumulated checks remain pending. The quickbar port
itself has not yet begun.

## Further fixture findings

The first direct slot-event candidate correctly failed because the metadata
fixture supplied valid definitions but did not enable them; live quickbar input
checks Enabled independently. Enable the actual spell definitions for the
successful-drag cases. With this fixture correction, all 150 slot-event records
pass. The 50 saved-spell-row records pass exact byte and independent round-trip
contracts, including signed nonpositive dimensions and preserving upper flag bytes.

The timed-selection candidate initially assumed saturated presses keep refreshing
the timeout. Original C does not: with half-rate gaps the observed/required rows
are 0,1,2,3,4,4,0,1. The test now explicitly preserves that sequence. On ordinary
reachable sequences the row pointer stays consistent even though the saturated
branch does not rewrite it; no production change is warranted for this branch.

The lifecycle candidate adds actual creation, hiding/showing and destruction
across two widths, all three player classes, missing-player and game-mode cases.
It restores the production startup ability records and image suffix rather than
inventing their contents. The remaining rendering/asset, trap, cross-row input
and saved-ability coverage must still be completed before conversion.

All twelve current roots pass without skips (c-lifecycle2, 33.913 seconds),
covering **1,464 records**. The constructor's missing-player contract required
collecting the fixture's already-created windows before measuring newly created
ones; production C was unchanged. Current captures are recorded in
quickbar-captures.json and asserted by their respective tests. This is a recovery
checkpoint; final repeated/default/server/highres qualification awaits the rest
of the connected baseline coverage.

The quickbar-warrior scenario is undergoing a coordinate pilot using the already
qualified spellbook binary. That pilot is not the final C reference: the latter
will use a fresh binary including the saved-row correction.

## Repeatability correction before final baseline

The 240 slot-rendering records pass, including actual pixels, image/fallback
paths, cooldown/highlight boundaries and signed byte timers. An independent
process exposed unstable addresses in the previously captured lifecycle
destruction records: C leaves child addresses in named 1049524 and mapped
1049528 after their parent is destroyed. All 48 differences were those two
identified pointers in the 24 destruction records.

Remember each actual window identity before calling the destructor, then retain
that identity when snapshotting the stale address. This changes capture
normalization only; do not clear the production words or mask arbitrary values.
The previous lifecycle hash in a6ee4a23 is not a final repeatable golden. Re-freeze
this group after the corrected fixture passes and require an independent repeat.
Raw differences are in build/port-quickbar/lifecycle-repeat-diff.json.

The headless quickbar-pilot completed successfully; inspection confirms both
ability assignments, their swap and Eye of the Wolf activation. Its binary
predates the row prerequisite, so retain it as coordinate evidence only.

## Final C qualification started

All 22 direct quickbar roots pass (2,242 records, c-final-candidate,
42.323 seconds). The corrected lifecycle normalization repeats exactly.
Hotkeys cover valid and blocked animation states and missing players; the
additional mana test binds the actual drawable so disabled/animation/buff and
mana-deficit paths are tested separately. The direction-button fixture selects
the actual self-cast child rather than its key-label sibling.

The reviewed quickbar-batch.json qualifies all 22 quickbar plus 24 spellbook
roots, with explicit hashes for all **6,936 records**. Default/repeat and server
run in separate processes concurrently; highres and the fresh C gameplay binary/
capture/repeat follow. At most two qualification jobs run together on this
4-vCPU/8-GiB VM. Sources remain immutable until all readers join; no parallel
agents are involved. This bounded scheduling choice is reversible and its timings
will be recorded instead of assuming a speedup.

Native translation and its production/interface/known-suite/full accumulated
milestone gates remain pending.

Default/repeat qualification passed all 46 roots and hashes in 148.864 seconds;
server passed in 161.925 seconds. Both source fingerprints are unchanged and the
readers are joined. Highres and the fresh C gameplay build/capture/repeat are now
running concurrently. Initial typed-state and activation drafts live only under
build/port-quickbar; they are not compiled/accepted production source.

Highres passed all 46 roots in 88.258 seconds on unchanged source. Its reader
is joined. The fresh C client built in 61.550 seconds and the forced-map capture
passed in 57.800 seconds. The independent gameplay repeat is still running.
All four phase source manifests match. Three ignored native drafts (state,
activation and row transitions) remain uncompiled and outside source.

The accumulated server milestone passed **1,121 selected/completed tests**
(1,104 in the root package), with only the existing opt-in
`TestMapPopulationPrerequisiteProbe` skipped. All 46 book/quickbar hashes match.
Driver: 815.448 seconds; phase: 815.731 seconds, source unchanged.
Evidence: `build/port-quickbar/native-qualified-server`. Default remains active;
production builds/interface/known-suite/replays have started in the freed slot.

## Final review: additional C boundary contracts

The first accumulated default/server passes (1,125 / 1,121 selected tests) are
superseded by these subsequent source corrections. The production run was
terminated and joined; its partial artifacts are diagnostic only.

- C compares the arrived tray destination to screen height as unsigned; a signed
  Go comparison differed for destinations with the high bit set.
- Previous-row navigation rejects an out-of-range result; sharing the next-row
  wrap condition incorrectly selected row zero for large selected-byte values.
- Trap row navigation computes an address without accessing the row. Preserve
  its byte wrap and address arithmetic rather than adding a Go array-bounds panic.

These are translation corrections, not changes to original game behavior. Three
new test roots (47 records) are being qualified against unchanged C in an isolated
baseline worktree. Original 46-group goldens remain unchanged. Freeze and commit
the supplement before accepting its native results, then repeat qualification on
the corrected source. The final C reduction includes one trailing blank line
removed from the shortened GAME2.c: **2,494**, leaving **64,317**.

The supplemental tray fixture initially used the GUI position setter, which
normalizes overflowing rectangles and therefore changed the coordinate under
test. It now sets the raw coordinate to isolate the original comparison; no
production C was changed. Process adjustment: finish the arithmetic/pointer review
before launching long accumulated gates. Late discovery here requires rerunning
previously passing default/server sweeps rather than accepting stale evidence.

## Supplemental C baseline committed

**a6c1640e** is pushed. The 47 new records (14 row, 21 tray, 12 trap-row) pass
unchanged C in default/repeat/server/highres: phase times 29.810 / 106.170 /
23.662 seconds. All readers joined. The staged and committed 1,802-file source
tree matches every supplemental C manifest exactly. No native implementation
was included in that commit; source proof is c-boundary-index-proof.json.

The corrected native production code still matches every original record
(native-focused-03). Native-focused-04 and native-final-server now include the
committed supplement: **49 groups / 6,983 records**, with all original hashes
unchanged. Full final qualification remains pending. The source-review ordering
adjustment is recorded in PORT.md.

## Native recovery checkpoint

The full **49 groups / 6,983 records** pass in native-focused-04 without skips
or golden changes (92.620 seconds), with 1,812 unchanged source fingerprints.
Save the native implementation and docs as a recovery commit now, while the
fresh default/server accumulated sweeps run. This protects the completed
implementation work from VM loss; it is explicitly **not final qualification**.
Highres, production/ABI, known-suite and replay gates remain, followed by the
completed-chunk documentation commit. Do not begin the next port until those pass.

The native recovery commit is **f7c75139**, pushed. The final sweeps use its
unchanged 1,812-file source. During default/server execution, both test processes
used about 2.2 CPU cores combined and the VM reported roughly 4.9 GiB available
memory. Start the independent highres sweep with GOMAXPROCS=1 and the existing
768MiB heap bound in the remaining capacity. This replaces the earlier two-job
cap for these isolated target sweeps; do not add a fourth job. No agents or extra
checks are involved. Record actual timings and outcomes rather than claiming a
speedup from this resource observation alone. This choice is reversible.

Final server sweep passed all **1,124 selected tests** / 1,107 root-package tests
with only the known opt-in diagnostic skip and all 49 hashes unchanged.
Driver time: 813.561 seconds. Its reader is joined; final production qualification
starts in the freed slot while default/highres continue. Source is unchanged.

Final default sweep passed all **1,128 selected tests** / 1,111 root-package
tests, with the same single existing skip and all 49 hashes unchanged. Driver:
760.999 seconds; phase: 761.323 seconds.
The reader is joined. Highres and production/replay qualification remain active.


## Final fixture ownership correction

The final highres accumulated sweep passes all 1,128 selected tests / 1,111
root tests with the same diagnostic skip and all 49 hashes (712.226 driver
seconds). Default/server/highres all qualify the production source in f7c75139.

Production-final built and audited all three binaries, but the full asset suite
reported 1,555 failures instead of 1,553: the two additional entries are
TestCodeStatic and its package. The quickbar fixture's broad backing-region
snapshot began with a typed memmap access at an extracted named-global address.
That fixture existed in the C baseline; no gameplay mismatch was involved.

The fixture now explicitly slices Blob.Data for its raw backing snapshot and
registers addresses from that same slice. Actual named globals remain separately
owned by PortTestQuickbarWords. This preserves the captured bytes, including
inert backing bytes at extracted offsets; it does not reinterpret separately
allocated globals as contiguous storage. TestCodeStatic passes without changes
to the checker. PORT.md adds this cheap check before long qualification.

Only quickbar_owner_porttest_test.go changed after the accumulated sweeps.
Recheck all 49 focused captures in default/server/highres in fixture-final, and
rerun all production gates in production-final-02. Do not repeat the accumulated
corpus for this fixture-only equivalent access change; retain exact source-diff
proof connecting its results to the final tree. Production code is unchanged.
The failed production-final remains diagnostic evidence and is not acceptance.


Final production gate `production-final-02` passes in 176.601 phase seconds.
All three rebuilt binaries are byte-identical to the earlier builds, confirming
the fixture change does not affect production. ABI audits retain 19 interfaces,
retire 68, and find no test helpers. The asset suite has exactly 1,553 known
failure entries and 15/3/32 passing/failing/skipped packages. Fresh quickbar and
flat replays both pass (63.152 / 47.715 driver seconds), each removing 51 copied
maps and regenerating the warrior map exactly. The quickbar replay compares all
22 frames to quickbar-c; the flat replay uses lists-flat-native.

`fixture-equivalence-proof.json` compares all 1,812 fingerprints from each
accumulated sweep with the final source: only quickbar_owner_porttest_test.go
changed. It includes that exact ownership-only diff. Final focused and production
manifests match each other exactly. To reproduce the affected fixture checks,
use the native-focused phase's command and 49 hashes with tags porttest,
porttest,server and porttest,highres, each with separate output log/result paths.
The ignored fixture-batch.json records those three steps explicitly.
