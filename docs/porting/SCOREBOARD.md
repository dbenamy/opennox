# Scoreboard and rank presentation

Status: **Go conversion qualified**, from committed/pushed C baseline
**73017352**. The connected 31-routine batch removes **1,537 physical C lines**;
**77,120 C lines / 91 files / zero reference C** remain. All 2,137 frozen results
match in default/server/highres, with production and gameplay qualification below.

Scope: 31 connected routines, about 1,504 function-block lines, across
client__gui__guirank.c and GAME2_1.c. Includes window construction, player/team
collection and ordering, row/column formatting, rank/status/heading rendering,
objective indicators and mode/visibility helpers. The existing Go mode-cycle and
network-refresh callers remain actual integration. Removing guirank.c retains
its shared window definition in vardefs.c. Function-block
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

The headless scoreboard scenario was developed with the original C
scoreboard in the qualified binary. F9 is the installed default rank key; cycling
in ordinary play requests players, teams, top three and closed. Visual inspection
confirmed the corrected screenshots show those states before accepting them as
baseline evidence. Capture names alone were insufficient. A fresh repeat with
golden updates disabled passed before freezing.

Local drafts, scope/caller audit, logs and captures: build/port-scoreboard.
Original assets and archive remain unchanged. The qualification results and original development diagnostics follow.


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

## Go implementation and qualification

C baseline **73017352** is committed and pushed. The connected translation is qualified. It removes 1,537 physical C lines, including
the guirank.c preamble while relocating its one shared window definition to
vardefs.c. Working-tree count: **77,120 C lines / 91 files / zero reference C**.
Five scoreboard C entry points remain for actual C callers. All other scoreboard
entry points, plus the now-private briefing-stage bridge, retire. Go callers invoke
Go directly; typed shared records have compile-time size and field-offset checks.

Native-a stopped during discovery (17.702s): an existing Go collector wrapper
still called the retired C function. Redirected it to the Go collector. No goldens
changed. Native-b passed all fifteen roots in **182.378s** and all 2,137 hashes match.
Review then removed permanent interning of dynamic score/time scratch strings and
added alignment assertions. The completed checks are recorded below.

| Check | Result |
| --- | --- |
| default affected tests | 287 selected/started/completed roots; 276.664s; all captures exact |
| server affected tests | 285 selected/started/completed roots; 265.129s; all captures exact |
| highres affected tests | 287 selected/started/completed roots; 191.406s; all captures exact |
| default post-review family | 15 roots; 182.755s; all captures exact |
| server post-review family | 15 roots; 181.610s; all captures exact |
| highres post-review family | 15 roots; 57.541s; all captures exact |
| opennox production build | 58.493s; ELF32/i386/SSE2/CGO; five retained / 27 retired interfaces; no test helpers |
| opennox-hd production build | 9.425s; ELF32/i386/SSE2/CGO; five retained / 27 retired interfaces; no test helpers |
| opennox-server production build | 60.855s; ELF32/i386/SSE2/CGO; five retained / 27 retired interfaces; no test helpers |
| Full-assets comparison | Exact 1,553 failure entries; 15 passing / 3 failing / 32 skipped packages; expected exit 1 |
| client-scoreboard-port | 30.845s; fresh assets/save; reference comparison passed; updates disabled |
| client-scoreboard-chapter-port | 38.394s; fresh assets/save; reference comparison passed; updates disabled |

After broad qualification, the test adapter's two empty-callback operations were
changed to invoke the production Go callback directly. Constructed parent/group
callbacks already had nil-response contracts. The entire fifteen-root scoreboard
family repeated in all three targets after this cleanup; production fingerprints
are identical to the broader qualification, so builds/full-suite/gameplay evidence
is reused rather than rebuilt for a test-only change. Source fingerprints stayed
unchanged during each run, and all captured JSON was compressed with a verified
lossless round trip. No golden was regenerated for the Go implementation.

The five retained C entries serve existing network/state callers. Private Go
callers invoke Go directly, including the collector wrapper. Removing the last C
briefing-stage caller also removes its C bridge; its Go implementation remains.
The declared 80-byte player and 56-byte team records have size/alignment/offset
checks. UTF-16 clipping, raw record name copying and stale bytes after terminators
remain compatible; dynamic score/time formatting writes existing scratch without
permanent string interning. The shared window variable remains in vardefs.c.

The hosted route checks five visible states; ordinary chapter gameplay checks
all eight existing frames. This does not claim remote-player gameplay, populated
teams in the hosted scene, or audio playback quality. Those row and state branches
are covered by actual-owner fixtures as detailed above. Exact unsupported sorting
sentinels and malformed headless host counts are still separate review items.

Local native evidence: native-qualification.json, post-review-qualification.json,
native-source-fingerprints.json, post-review-source-fingerprints.json, all logs and
compressed captures under build/port-scoreboard. The next candidate is the connected
minimap renderer, fifteen routines / 696 C function-block lines; its C baseline
must be established before translating it.
