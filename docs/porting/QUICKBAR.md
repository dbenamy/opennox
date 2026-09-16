# Quickbar UI port

## Status and scope

Baseline preparation follows qualified spellbook conversion **41b8abfb**.
Production C is unchanged: **66,811 lines / 87 files / zero reference C**.
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

Direct frozen C captures and independent contracts; repeat and server/highres C
qualification; committed baseline; native implementation and comparison; production
builds/interfaces, known asset-suite comparison, fresh gameplay replays and the next
full accumulated UI milestone. Commit/push each qualified chunk and continue.

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
