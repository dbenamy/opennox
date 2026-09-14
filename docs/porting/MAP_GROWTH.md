# Map growth and door placement

Active scope: seven routines / **1,082 C section lines**, 4D4790..4D5D20 in
GAME3_2.c. Theme native **14b1caa7** is qualified and pushed. Physical C is
**102,416 / 148 files / zero reference C** after the qualified C prerequisites
below (−3). No growth conversion or locked hashes.

Three external callers need door placement, initial layout and frontier traversal;
recursive dispatch, room fill, hall growth and branch choice can become private Go.
Read-only scope, source, caller audit and table values are in build/port-map-growth.

Fixture bootstrap reuses real room/grid/painting owners and scoped allocation
observation of generator-owned rooms/exclusions. It initializes actual ring and
turn tables, copies the supplied config into the real global config, and provides
distinct door type IDs with the existing door transfer/update layout. Input record
bytes, generated records, grid occupancy, object state, disposal and random state
are captured by the established fixture. The first eight normal/ring-layout
contracts are compiling in bootstrap.log; no behavior changes have been made.

Plan: exercise bounded recursion, all room/hall directions, geometry/rate limits,
obstruction/retry cases, ring/start topology and door selection/orientation. Repeat
complete C captures in default/server/highres and lock/commit/push before converting.
Then run accumulated qualification, production symbol audit, known full suite and
fresh headless gameplay. Count/document/commit/push and continue. Full generation
orchestration and filesystem map save/restore remain a later connected batch.

## Original prerequisite evidence

The original basic fill return-code probe passed, but that did not establish valid
directions. The fixture now captures normalized room/exclusion bytes immediately
before release, then records disposal. The stronger probe finds invalid hallway
kind **1107434383** (fill-directions-original.log). Replace the four separately
indexed direction locals with a real directions[4] array. The correction passes
all 128 blocked-mask/seed contracts plus the original fill and eight bootstrap
contracts, **0.911s** root time (directions-corrected.log). C decreases by three
physical lines to **102,416 / 148 files / zero reference C**.

The original merge probe also fails: actual config mergeRate=100 yields north
connection counts **[0,1]** as the unrelated named global changes from 0 to 100
(merge-original.log). Four reads across fill/hall growth use that named word;
point all four at the actual config blob offset1549844, after recording this
original evidence. All four reads are now corrected. The 192-case regression covers four directions,
0/50/100 configured rates, eight seeds and separately varied named values, checking
connections and random tails. Corrected default map checks are running in
prerequisite-corrected.log.
Qualify both prerequisites across variants and commit/push before native conversion.

Corrected default growth and existing map checks pass in **57.006s**
(prerequisite-corrected.log). The fixture contracts comprise **331 cases**:
eight initial layouts, one explicit direction probe, 128 blocked-direction cases,
two original merge controls and 192 configured-rate cases. All current theme and
older map hashes remain unchanged. Server/highres qualification is running through
qualify-prerequisite.py; no source edits during those runs. Both fixes are still C
prerequisites; no growth routines have been ported yet.

Server/highres prerequisite and existing-map checks pass in **138.224 / 67.847s**
wall (prerequisite-variants.json); all prior mandatory hashes remain unchanged.
The prerequisites are qualified and ready to commit/push. Summary metadata is in
prerequisite-qualification.json. Next expand complete growth C captures from the
unapplied corpus.go.stage draft, add door/frontier cases, repeat and lock all hashes
before native conversion. No compile/test is active.

## Baseline expansion after prerequisite commit

Prerequisites **2fc007db** are committed/pushed. Five C capture groups now pass
**4,068 cases** plus the 331 contracts in **11.330s**: branch choices (3,528),
initial layouts (168), recursion guards (100), room expansion (128), and hall
expansion (144). These are not locked yet. The initial c-first run reached an
uninitialized renderer because the fixture used 1<<22 instead of the existing
GameFlag22 constant (0x200000). Use the named constant, matching population's
fixture; c-fixed.log passes. No production rendering change was needed.

The fixture now captures real waypoint records, list heads and cleanup for door
cases. The four-direction positive door probe is compiling in door-probe.log.
That snapshot extension means earlier captures must be regenerated before final
locking. Add frontier/door and bounded geometry/obstruction coverage, repeat every
C capture across variants, then lock/commit/push before native conversion.

The full eight-group corpus passes **5,765 cases** plus focused contracts in
**25.545s** (c-complete.log), including 1,216 door, 97 repeated-frontier and 384
obstruction cases. Waypoints use their real server owner, are fully captured, and
are disposed between fixture cases.

A translated-room regression then found that two door-scan reads treated integer
grid coordinates as floats. The original positive-offset case6 produced no doors
or waypoints instead of one/two (door-translations-original.log). Restore integer
word reads for GridX/GridY. The correction changes no physical C line count. All
16 translation cases now capture full state and assert the corresponding projected
door-position shift as well as counts. Corrected full growth captures are running
in c-final.log; qualify this correction with the complete baseline before porting.

## Locked corrected C baseline

All **5,781 cases / nine full captures**, plus **335 separate contracts**, pass.
The additional 16 translated-door contracts are included in their capture count.
Default repeat / server / highres reproduce every hash with prior map checks in
**82.945 / 161.009 / 93.334s** wall. The earlier eight hashes are unchanged by the
two integer grid-coordinate reads. All hashes are mandatory in the committed
fixture; the extension bypass is removed. Locked smoke passes in **25.113s**.

Evidence: baseline-captures.json, baseline-variants.json,
baseline-capture-archive.json, early-capture-archive.json and baseline-locked.log
under build/port-map-growth. Complete captures are gzip archived with verified
hashes. Physical C remains **102,416 / 148 files / zero reference**.
Commit/push this corrected baseline before converting the seven routines.
