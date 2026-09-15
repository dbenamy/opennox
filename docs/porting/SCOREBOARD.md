# Scoreboard and rank presentation

Status: **C baseline qualified and frozen**, after qualified/pushed briefing-window
conversion **0b3ed13d**. Production remains **78,657 C lines / 92 files / zero
reference C**. No scoreboard production translation applied yet.

Scope: 31 connected routines, about 1,504 function-block lines, across
client__gui__guirank.c and GAME2_1.c. Includes window construction, player/team
collection and ordering, row/column formatting, rank/status/heading rendering,
objective indicators and mode/visibility helpers. The existing Go mode-cycle and
network-refresh callers remain actual integration. Removing guirank.c must retain
its shared window definition or deliberately migrate its bindings. Function-block
count is not the final net physical C LOC reduction.

The baseline reuses actual GUI/listbox/static-text, renderer/fonts/strings,
32-player and team owners from the presentation fixtures. It owns scoreboard
windows, scratch records, local-player code, original format/color tables and
settings. No sorting, formatting or window algorithm is copied into test-only C.
Expectations are literal SHA-256 values from the original C owners; mismatches
always fail. Eighteen capture groups cover 2,137 results.

Initial fixtures cover repeated construction/cleanup, objective state updates,
row membership, names/classes/status strings, normal and elimination ordering,
rank calculation, color wraparound, populated rendering, mode callers and columns.
Additional fixtures cover real team membership, supported headless-host filtering,
timer/lesson branches, refresh boundaries, objective buffs and populated rows. Exact sort-sentinel inputs need separate
invariant analysis; do not misclassify unsupported state as a port regression.

A fresh headless scoreboard scenario is being developed with the original C
scoreboard in the qualified binary. F9 is the installed default rank key; cycling
in ordinary play requests players, teams, top three and closed. Visual inspection
must confirm the screenshots actually show those states before accepting them as
baseline evidence. Capture names alone are not evidence. Repeat from fresh assets
and saves with golden updates disabled before freezing.

Local drafts, scope/caller audit, logs and captures: build/port-scoreboard.
Original assets and archive remain unchanged. This report is a development
checkpoint, not a qualification claim.


## Development diagnostics

The first three discovery attempts stopped before running roots (43.560s,
44.507s and 44.505s). The test preamble must match the existing cgo declarations,
including the window pointer at 1090100 and the original typedef spellings.
No production or expected behavior changed; all compilers joined before fixes.
The later runs below supersede those discovery failures.

Visual inspection rejected the initial solo capture: the input handler requires
GameOnline for F9, so its four named screenshots did not show a scoreboard.
The existing `-autosrv` path reaches a locally hosted estate map using the
qualified binary; its initial message-of-the-day dialog obscured the first
capture. A second hosted capture explicitly clicks the dialog's OK button.
Neither preliminary capture is accepted as scoreboard baseline evidence.


Development-d compiled but its first root found an unowned listbox palette
pointer table (182.281s total, one of nine roots reached). Reused the existing
production palette initializer with precise save/restore. The real widget strips
a trailing newline, so a newline heading becomes empty stored row text.
Development-e completed all ten roots (42.484s); seven passed and three identified
missing file-qualified strings. Installed guirank.c-qualified keys and removed
case-duplicate logical IDs, matching the real string manager's lookup rules.
Development-f passed all ten roots in **45.156s**, producing **1,927 results /
twelve groups**. All corrections are fixture setup/expectations before freezing;
no production algorithm changed.

Added real team creation, registered client drawable membership, objective buff
lookups and supported headless-host filtering. The headless fixture always owns
an active slot-31 host before testing its exclusion; it does not call the C
sorter with an inconsistent zero-player headless-host count. These additional
fixtures passed and their expectations are now frozen.

## Hosted-game baseline

The corrected local hosted estate route visibly shows Teams and Players, Teams,
Top 3 Players and closed state. Input coordinates refer to the 1280x960 X display,
not the 1024x768 game frame; this explains the first missed dialog click. The
corrected capture is client-scoreboard-hosted-develop3. A separate fresh repeat
completed in **30.099s**, matching all five screenshots with golden updates
disabled. Tracked [scenario](scoreboard-hosted.yaml) and
[pixel manifest](scoreboard-hosted-pixels.json) preserve the route and decoded
reference hashes. Start with `-autosrv -port 18590 -autoexec "load estate"` using
fresh assets/save, seeded GODEBUG and OpenAL null backend as in RECOVERY.md.

This scenario is a local observer with one player row and no populated team rows
or remote participants. Actual-owner fixtures supply populated-player/team
coverage. The original solo and dialog-obscured captures remain rejected as
scoreboard evidence. Original assets/archive are unchanged.


Development-g completed thirteen roots (45.740s): all prior 1,927 results matched,
and the new headless cases passed. Team/objective drawable cases exposed a fixture
accounting mismatch: actual registered sprites were not in the reused lifetime
checker's ownership list. Added explicit ownership and failure-safe cleanup around
the real allocator/registration/deletion calls. Development-h stopped at discovery
(1.259s) on a 386 constant-conversion error in that test accounting; corrected.
Development-i passed **fourteen roots / 47.783s**, yielding **2,025 results /
seventeen groups**, with all previous results identical. Independent contracts now
include actual team/player row colors and strict refresh/wrapped-frame boundaries.
The existing empty callbacks also have direct/nil-response construction contracts.

## C qualification and frozen expectations

Affected qualification passed **286 default roots / 179.337s**, **284 server /
263.888s**, and **286 highres / 185.657s**. All selected roots started and
completed; all 2,025 results in the original seventeen groups repeated exactly.
The affected selection covers scoreboard, briefing, journal, inventory/window,
listbox/slider, renderer/object/meter, gameplay text/report/effect/particle and
objective/team dependencies. This follows the complete three-target accumulated
milestone at 0b3ed13d; it does not claim another full-corpus run.

Exact production fingerprints match 0b3ed13d (1,074 non-test Go/C/header files),
so its builds, interface audit, full-assets known-failure comparison and chapter
gameplay are reused. The new hosted scoreboard repeat supplies additional
integration. Source fingerprints remained unchanged throughout qualification.

Added 112 actual record/draw results for empty, 25–28-character, non-BMP and
mixed names using ordinary and narrow-advance fonts, both score modes and
elimination settings. The independent contracts check unchanged player names and
all adjacent record fields. This passed C before freezing. These additions came
after the broad affected run; the final locked family repeats cover all fifteen
roots / 2,137 results in each target. Do not count that added root in the earlier
286/284/286 totals. A preliminary direct test command used the repository root,
which has no Go module; rerunning from src succeeded without source changes.

The narrow font exposes a legacy layout detail: a player name can contain 27
UTF-16 units while a scoreboard row reserves 26 before its team field. The C
routine copies and clips first, then overwrites the team field. Compatibility
requires preserving the resulting record bytes and actual rendered string, not
silently truncating at the declared name array. Any user-visible correction should
be reviewed separately. Exact sort sentinels and inconsistent headless host/player
counts remain outside supported fixtures; they are not newly validated inputs.

Final locked repeats passed all fifteen roots and all 2,137 frozen results in
each target: **default 50.224s, server 49.019s, highres 49.462s**. Every selected
root started/completed, captures matched and source fingerprints stayed unchanged.
Local evidence lives in
build/port-scoreboard/{c-qualification,locked-c-qualification,frozen-captures}.json.
The original development failures above are retained to distinguish fixture fixes
from production changes; no C algorithm was changed to produce this oracle.
