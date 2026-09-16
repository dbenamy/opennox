# Quickbar UI port

## Status and scope

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
- All readers are joined. Native drafts are outside source and unqualified.

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
translating. Baseline preparation is not qualification; no native quickbar source
has been installed.

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
