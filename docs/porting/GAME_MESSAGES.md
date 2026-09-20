# Game-message dispatch and player notices

Qualified parent: `0d5a03b2` (pushed), 8,820 C lines in 38 files, zero reference C.
Original-C client-state/notices baseline **beb38b24** is qualified and pushed:
52,894 captured cases, 32 frozen captures, 265/262/265 affected target roots.
The selected 904-body-line conversion is fully qualified. Current C is
**7,897 lines /37 files**, zero reference C: **923 physical lines removed**.
Final native sweeps pass 266/263/266 roots and every frozen capture, plus ten
independent boundary cases. Fresh production/ABI, exact known full-suite and
headless gameplay/save-load checks pass. See
[client-game-state-native-qualification.json](client-game-state-native-qualification.json),
[the revised boundary](#revised-conversion-boundary-client-state-and-notices), and
[native progress](#native-client-statenotices-conversion-in-progress).

## Scope and audit

Three connected bodies contain 3,841 C body lines: client dispatch (3,064),
server dispatch (449), and the client's player-notice handler (328). The Go outer
dispatchers already handle a separate set of message kinds. The initial comparison
finds 145 named client C cases and 16 server C cases, with no named overlap with
their Go outer dispatchers. Client dispatch calls the notice handler. Complete
numeric-value, fallthrough, shared-owner and callback review before freezing.

Draft selection and complete literal references are under
`build/port-game-messages/{selection,references}-draft.json`; the audit deliberately
retains declarations and return calls for manual classification. These are audit
artifacts, not approved conversion instructions.

## Testing plan

Build original-C baselines by behavior family, then qualify each accepted connected
conversion boundary. The first boundary is client state and notices; the remaining
client/server dispatch follows in subsequent batches.
Reuse actual players, teams, drawable/object pools, inventory, wall, text and
serialization owners. Observe only external effects where an existing fixture
already provides the boundary. Every selected case needs a recorded contract for
its return/consumed length, state mutations and downstream effects; group cases
only when they have the same data layout and observable behavior.

- Notices: exact Unicode text, formatting/signedness, missing players/teams,
  message ring and expiry, sounds, dialog/death-notice routing, sign text and play
  state, unknown kinds, read-only input and valid message lengths.
- Server actions: actual player state and owned objects, eligibility gates,
  returned lengths, coordinate/ID decoding, inventory/waypoint/trade/vote effects.
- Client updates: actual object and wall lifecycle, counters and HUD state,
  inventory/equipment, winner notices, spell/effect/audio routing and timing.
- Outer dispatch: concatenated valid messages and consumed lengths, the existing
  Go-owned cases, and the C fallback boundary. Preserve documented valid-domain
  behavior; review any newly defined rejection behavior separately.

Keep inherited contracts for the called behaviors, and add independent expected
state so broad dispatch tests do not merely replay captured results. Repeat C
captures in separate processes before freezing. The final native batch must pass
three-target comparisons, ABI/export audits, exact known full-suite comparison,
and fresh headless gameplay plus explicit save/load.

## Baseline development history

A thin original-C notice adapter and the first text-notice fixture are installed.
The fixture uses existing real player/text/audio owners and independently checks
full message-ring bytes, console text, expiry, sound requests, consumed length,
unchanged input, player ID width, signed stage formatting, play-state gates and
all unknown notice kinds. Its first focused compile/run is in progress; no hashes
have been accepted. Remaining notice cases and both dispatchers still need tests.

Numeric audit confirms no overlap with Go outer cases. Two C comments are stale:
201 is `MSG_TRADE` (comment says SERVER_QUIT_ACK), and 222 is total mana (comment
says total health). Classify and test by actual numeric kind and called behavior;
do not copy those comments into the Go translation. Evidence:
`build/port-game-messages/outer-dispatch-audit.json`.

The text fixture passes after correcting its sign-string key to the real
`inform.c` namespace. Five notice roots now cover all 22 recognized kinds plus
all unknown byte kinds. Added contracts use real team records and nonempty title
data, including byte-narrowed team IDs; actual award tables and spell definitions;
the real death-feed owner; and the existing dialog-construction boundary. Their
first combined original-C run is in progress (37845). No expectation is frozen.

Alias-domain review: the player record has 255 alias entries; 255 is the reserved
exhaustion sentinel, already excluded by the approved sender fix documented in
NETWORK_ALIASES.md. The C baseline will cover every in-range slot (0–254), with
whole-record independent checks and only deterministic changed bytes captured.
The new receiver should explicitly leave the table and following fields unchanged
for the reserved slot while consuming the existing ten-byte message; add this as
a native-only invalid-domain contract and record it as a reversible decision.
No receiver change is installed yet.

## Partial baseline checkpoint

Seven roots pass 2,210 cases: 458 notice cases, 987 unknown server-kind/subkind
cases, and 765 alias updates (255 slots × three value sets). Three repetitions and
a further independent process pass, with seven byte-identical captures. All jobs
are joined. See [game-messages-initial-fixtures.json](game-messages-initial-fixtures.json).
These recorded hashes are evidence, not yet runtime goldens; the whole dispatch
scope is not qualified and production remains identical to `0d5a03b2`.

Next: reuse the full player-controls owner for actual server actions, then cover
the client dispatch families and nested trade/game variants. The server-storage
draft is consumed and must not overwrite reviewed source. C remains 8,820 lines.

Partial fixture checkpoint `33630cba` is pushed. The next installed contract routes
waypoint messages through the existing complete player-controls owner. It compares
180 cases against the qualified coordinate operation: all three slots, absent or
existing waypoint objects, six full-width state values, and five coordinate pairs
including uint16 bounds. Dispatch return length and unchanged input are checked
before comparing normalized gameplay state; no gameplay field is excluded.
The inherited waypoint roots run alongside it. First job 70551 is in progress.

Owner audit to finish with client conversion: the fade-objects setting is defined
in cdecode but also used by Go drawing, config and the advanced-video checkbox's
pointer getter. Move that owner and preserve all three access paths. Blue/Violet
spark caches have Go initializer/effect consumers and C dispatcher reads; their
remaining GAME2_3 references are declarations. The notice-state word 1200768 is
read/written by the client dispatcher. Move these actual owners, update fixtures,
and audit mapped metadata rather than leaving obsolete C definitions behind.

Waypoint dispatch passes all 180 paired-state cases and inherited waypoint hashes
(job 70551 joined). The optional fixture route checks the original C consumed
length and read-only input, then compares every existing normalized owner field
against the separately qualified direct operation. Eighty-four missing-target
server-action cases are now installed and running with all message roots and
inherited waypoints (53145); positive inventory/spell/trade actions remain.

Build-speed trial: the compiler exceeded its inherited 768 MiB soft heap limit
while compiling the large root test package. This single-job focused run gives
build tools 1536 MiB but launches the actual 386 test through
`-exec='env GOMEMLIMIT=768MiB'`. GOMAXPROCS stays 2. The VM had 5.2 GiB available
before the trial; test-runtime limits and test scope are unchanged. Evaluate
elapsed time and peak memory before adopting this for other runs.

The nine message roots plus inherited waypoints pass (53145 joined, 33.159 seconds
including compilation); all preceding seven captures remain identical. Working
follow-up adds 48 positive/disabled drop and item-use cases against qualified
operations, with independent callback count, kind, actor and item identities.
These prevent an unexercised positive path from matching another no-op. Job 58089
is running; source must remain unchanged until joined.

## Action fixture checkpoint

Ten message roots pass 2,522 cases, plus the inherited waypoint root. Final guard
assertions pass for every owner; repeated and separate-process outputs match all
ten captures. Static checks pass. All jobs are joined. See
[game-messages-action-fixtures.json](game-messages-action-fixtures.json).
This is still a partial original-C baseline, not a conversion milestone.

The repeat with a root-only rebuild took 39.076 seconds. Continue the compiler-only
1536 MiB allowance for single-job focused checks, preserving 768 MiB for the 386
test process. The earlier legacy-plus-root rebuild took 124.990 seconds, so these
numbers are not a controlled speed ratio. Fifteen verified superseded test-cache
archives were removed (1,131,099,228 bytes). Audit/apply scripts are consumed;
qualified binaries, captures, original assets/archive and saves remain.

Client input ownership differs from the server/notice families: the C complex-
object case updates the animation nibble in the message, and several equipment
and creature/NPC cases rewrite unit-code words. Preserve and capture these writes;
do not apply the read-only-input assertion indiscriminately to client dispatch.
The read-only audit list is `build/port-game-messages/client-input-writes.json`.
The real connected gate is mapped word 5D4594+815764; fixture that owner instead
of replacing the connectivity predicate.

The server pickup case deliberately calls the live GameEx `OnLibraryNotice_420`
hook, not the commented-out ordinary inventory placement call. Preserve that
entry and its class/item behavior while GameEx remains C; a similarly named Go
placement API is not a substitute. Include weight-limit and hook-routing coverage
before converting pickup. The remaining C hook is production code, not a test
reference, and can be selected with its GameEx peers in a later batch.

Equipment passes 96 cases with explicit equipment-flag contracts, all prior message
roots and three inherited roots (equipment owners, bow/ammo and waypoints).
Three repetitions pass 42 root runs; eleven captures match across processes.
Evidence is under `equipment-{first,repeat}`; those jobs are joined.

Working pickup probe: 144 cases cover exact/over-limit/maximum carrying values,
player-state gates, full-range ordinary IDs, and the GameEx type/class decision.
The fixture observes the existing placement boundary with exact actor/item/flag
arguments and independent expected call counts. The expected rejection route
uses qualified private-message/audio operations, with a controlled nonempty mapped
class-failure key. Jobs 96744 and 2714 are joined and pass. Twelve message roots cover 2,762
cases; three repetitions with three inherited roots give 45 passing root runs.
All twelve captures match across processes. See
[game-messages-inventory-fixtures.json](game-messages-inventory-fixtures.json).
Production C is unchanged; remaining C stays 8,820 lines in 38 files.

## Client object-update baseline in progress

The first client dispatch fixture reuses the real drawable owner for simple and
complex object updates: connected/disconnected, all 256 status bytes, three frame
boundaries, ordinary/special-monster animation, and uint16 position boundaries.
It independently checks return length, lookup identity, drawable count, animation
state, activity words and exact input mutation. The initial position expectation
missed the qualified update helper's out-of-map fallback to (50,50); inspection
confirmed the 5888 boundary and the test was corrected before accepting captures.
No production source or frozen expectation changed. Creation, lookup variants,
camera and object lifetime branches remain to cover.

Existing-object updates pass 6,144 independent cases; creation adds 96 cases for
ordinary/static identifiers, missing types, coordinate boundaries and camera-only
messages. The first combined repeat passes 60 root runs with fourteen matching
captures across processes (39278 initial fixture failure, 43072 corrected,
68593 creation, 93151 repeat: all joined). The fixture error is explained above.
Height and animation-frame dispatch now pass another 8,288 cases: every height
byte, full-width frame values, present/missing targets, static/dynamic lookup,
connection gates and the special monster frame rule. These compare the complete
raw drawable before/after with only independently specified fields changed;
message input is read-only. Final accumulated repeat passes 63 root runs (21 roots, three repetitions),
with fifteen captures equal across processes and 17,290 message cases total.
All jobs are joined; see game-messages-client-object-fixtures.json.
A second audited cleanup removed fifteen superseded root test-cache archives
(1,130,553,318 bytes); scripts are consumed. Production C remains unchanged.

Client object checkpoint **3e010c73** is pushed. The next lifetime fixture covers
out-of-sight and shadow messages using real object/list owners: connection and
missing-target gates, static/dynamic removal, local-player protection, fading,
animation draw-data exemptions and frame-wrap deadlines. It checks exact consumed
length, unchanged input, object count/identity, active flags, shadow fields and
fade deadlines. All 576 lifetime cases pass. The accumulated repeat passes 66 root runs
(22 roots, three repetitions), 17,866 message cases and sixteen equal captures
across processes. Static mapped-memory checks pass. Jobs 55452/11844 are joined.
See game-messages-client-lifetime-fixtures.json. No production code changed.

The follow-up field fixture includes object enable/disable and byte-sized frame
reports, preserving the conditional draw-callback clearing for class bit 0x40000.
Its independent whole-record contracts pass; no old runtime golden was changed
(the combined baseline is not frozen yet). A new friend-list fixture exercises
add/remove/reset through original C with the qualified list owner, including
capacity, duplicate-first removal and normalization of the identifier high bit.
Its first run includes the inherited friend-list roots. Production is unchanged.

Object controls and friend lists pass the accumulated repeat: 22,130 message
cases /17 roots, plus eight inherited roots; three repetitions give 75 passing
root runs and seventeen equal separate-process captures. Static mapped-memory
checks pass. Jobs 53591/49579/80295 are joined. See
[game-messages-client-control-fixtures.json](game-messages-client-control-fixtures.json).
Only porttest sources changed since qualified production 0d5a03b2. No runtime
goldens are frozen and the combined dispatch scope is not yet qualified.
Remaining C stays 8,820 lines in 38 files. Next families include client health,
wall, equipment and effect reports, plus outstanding positive server actions.

## Client health reports

Health-change queues and meter reports add 6,596 independent cases, including
signed health amounts, full-width health/mana maxima, null player owners, ally
updates, high-bit identifiers, exact player-record writes and actual cooldowns.
Health-change notifications deliberately queue while disconnected; the other
selected reports retain their connection gates. This is established C behavior.
The combined repeat passes 93 root runs (31 roots, three repetitions), 28,726
message cases and nineteen equal captures across processes. Static checks pass;
jobs 34761/25150 are joined. See game-messages-client-health-fixtures.json.
Production is unchanged; remaining C is 8,820 lines /38 files. Next: wall reports.

## Wall reports in progress

Magic-wall add/change/remove uses the real grid and row-index owners. Tests cover
connection gates, even/odd grid positions, door/broken filtering, existing/missing
walls, byte-field extremes and unchanged neighboring wall payload. Secret-wall
open/close uses the actual secret-list head and production lookup, including host
mode suppression, the signed 16-bit lookup boundary, flags and optional close data.
A thin porttest owner borrows one real wall and restores the prior list head; no
production algorithm is substituted.

Review later: the first magic-wall probe used repeated DeleteAtGrid calls for
cleanup. Duplicate-position records exposed an existing row-index inconsistency:
position lookup selects the newest record but deleteByY locates the first record
at that coordinate. This can leave stale row links between fixture cases. Use the
full wall-owner Reset between cases so the dispatcher tests are isolated. The
production deletion behavior is outside this dispatch conversion and has not been
changed; investigate it as a separate wall-owner correction with its own baseline.
The initial fixture failure and corrected run are retained under walls-first and
walls-corrected. A missing closing brace in the new secret fixture was corrected
before its first build; no production code or frozen expectation changed.

Third cache cleanup recovered 1,133,631,892 bytes from fifteen verified superseded
root test archives. Audit/apply scripts are consumed; qualified evidence/assets
remain. Jobs 67029/88187 are joined.

The corrected wall probes pass all 672 cases. Accumulated qualification passes
105 root runs (35 roots, three repetitions), 29,398 message cases and twenty-one
matching separate-process captures. Static mapped-memory checks pass; jobs
25813/36629 are joined. See game-messages-client-wall-fixtures.json. All retained
wall expectations are unchanged. Next: inventory scalar, durability and charge
reports. The inventory-scalars draft is prepared but not installed.

## Client inventory reports

Armor/gold/stat-modifier storage and item durability/charges add 1,924 cases.
Contracts check exact float bit patterns and neighboring storage, full drawable
records, stack lookup and dragged-item precedence, high-bit code normalization,
equipped-state equality, real meter values/visibility, and unchanged input/grid.
Accumulated checks pass 123 root runs (41 roots, three repetitions), 31,322 message
cases and twenty-three matching captures across processes. Static checks pass;
jobs 32355/80982 are joined. See game-messages-client-inventory-fixtures.json.
The inventory-scalars draft is CONSUMED. Next: equipment dispatch using the existing
player/NPC owners and independent insertion/capacity/modifier contracts. The
ignored client-equipment draft is prepared but not yet installed. Production C
is unchanged at 8,820 lines /38 files; the full dispatch baseline is incomplete.

Client inventory report checkpoint **82656d33** is pushed. The player/NPC equipment
fixture is installed and its first probe is running (47271). It retains the
qualified owners and independent mask, capacity, modifier identity and canary
expectations, while calling the original dispatcher. Ordinary equip messages
supply the implicit all-255 modifier bytes; modifiable equip messages retain their
explicit bytes. Static-player identifiers are cleared in place only when connected.
The equipment draft/installer is consumed; never replay over reviewed source.

Player/NPC equip and player unequip pass 1,120 captured cases, plus missing-NPC
assertions. These preserve the distinct aggregate-mask behavior at full capacity,
modifier pointer identity, untouched slot tails, and connectivity gates. Unequip
preserves the input identifier including its high bit; equip clears that bit only
when connected. The accumulated repeat passes 141 root runs (47 roots, three
repetitions), 32,442 message cases and twenty-six matching captures. Static checks
pass; jobs 47271/85335/12730 are joined. See
game-messages-client-equipment-fixtures.json. Production C is unchanged.
Next: NPC appearance and remaining drawable attribute/effect reports.

Equipment checkpoint **6dbeda4e** is pushed. Working follow-up adds NPC appearance
(pool exhaustion, existing record reset, six independent RGB555 conversions and
input-code rewriting), light-color fields and local-player flag replacement.
The latter restores fixture ownership flags after each assertion so arbitrary
wire values cannot alter cleanup. NPC first probe 28720 passed and is joined;
appearance-first 39226 is running. Join before source edits.
A fourth cache audit identifies 43 superseded root test archives /3,262,857,596
bytes before 6dbeda4e; audit 31850 is joined. Apply is prepared but not yet run.

Appearance follow-up passes 13,348 additional cases: NPC appearance 9,216,
light-color fields 4,096, and player flags 36. The accumulated repeat passes
147 root runs (49 roots, three repetitions), 45,790 message cases and twenty-eight
matching captures. Static checks pass; jobs 28720/39226/52358 are joined. See
game-messages-client-appearance-fixtures.json. Fourth audited cache cleanup is
complete (60333 joined): 3,262,857,596 bytes reclaimed; scripts consumed. Next:
light intensity and enchantment reports. Production C remains unchanged.

Appearance checkpoint **2c50982e** is pushed. Light intensity passes 96 routing and
numeric cases alongside the frozen particle-light contracts (40418 joined).
The test uses the actual qualified particle-light entry, not Drawable.SetLightIntensity:
the former preserves the C helper's double-width fixed-point calculation. Independent
clamp/intensity/fixed-point assertions accompany whole-record parity. Working
follow-up covers object/item enchantment gates and restoration to a controlled,
nonzero type intensity. Enchantment first probe 61442 is running; join before edits.
The enchantment draft is consumed. No production C changed.

Light/enchantment follow-up passes 6,624 additional cases, including independently
specified restoration gates and exact local UI word/byte writes. The accumulated
repeat passes 156 root runs (52 roots, three repetitions), 52,414 message cases
and thirty matching separate-process captures. Static checks pass; all jobs
40418/61442/61025 are joined. See game-messages-client-light-fixtures.json. Next:
visual-effect creation dispatch, including cached type ownership and allocation
failure. Production C remains unchanged at 8,820 lines /38 files.

Six simple effect-creation messages pass 288 cases: cached/uncached/overridden type
IDs, allocation failure, signed coordinates, teleport/poof Y offsets, exact cache
writes and complete normalized drawable/list captures. The accumulated repeat
passes 159 root runs (53 roots, three repetitions), 52,702 message cases and
thirty-one matching captures. Static checks pass; jobs 77269/5048 are joined.
See game-messages-client-fx-fixtures.json. Production C remains unchanged.

Boundary review in progress: tested client-state cases plus all notices appear to
form a roughly 900-line connected conversion that need not wait for the remaining
trade, complex effects and server actions. The ignored covered-client-boundary-
review.json is provisional: its literal-label tally does not assign shared effect
blocks correctly. Exclude the entire shared effect block from the first conversion,
and audit goto targets, ownership and remaining callers before accepting this split.
No selection change, baseline freeze or native installation has happened yet.

## Revised conversion boundary: client state and notices

The next conversion will move 42 client-state opcode cases (39 complete switch
blocks, 576 C lines) and the entire 328-line notice dispatcher: 904 body lines.
The wider 3,841-line dispatch scope remains the overall audit, but waiting for
unrelated trade, complex effects and positive server actions is unnecessary.
This is a reversible batch-boundary decision under the user's standing authority;
review it later against qualification cost and actual progress. Shared effect
blocks remain wholly in C for the next batch, including the six tested simple
creation cases. Server action fixtures remain useful inherited coverage.

The exact source-boundary audit is currently under build/port-game-messages:
client-state-boundary-review.json /client-state-cases.txt. All selected shared
labels move together; there are no goto edges into or out of selected blocks.
All 52 non-notice called gameplay helpers already have Go definitions; notice
calls also resolve to Go except formatting/string primitives. The routing adapter
will try the native selected cases and leave the remaining production C dispatcher
as fallback. No C algorithm will be retained solely for testing.

Move the fade-objects owner with its last C consumer, preserving Go drawing,
configuration and the advanced-video checkbox pointer. The selected secret-wall
cases are the last C callers of sub_410550; retire that export and call the private
Go lookup. The notice C entry has only the selected client case and test adapter;
remove its C body/header once converted. Other dispatch globals stay with their
actual remaining C consumers.

The final camera contract passes 192 cases: actual cooldown state, unsigned
coordinates, full-width local-ID comparison, frame wrap, and camera updates even
when sprite creation fails. First probe 32727 and accumulated repeat 93391 are
joined. The repeat passes 162 root runs (54 roots, three repetitions), with 32
identical separate-process captures totaling 52,894 cases. Reviewed expectations
are frozen in game-messages-frozen-captures.json and enforced by the Go capture
helper. No native source is installed yet.

The tracked client-game-state-selection.json preserves exact C groups, source
hashes, and the accepted conversion boundary. The broad affected test selection
also includes combat overlays, client meters/inventory/presentation, player state,
drawables, world grids, interaction, gameplay text and book awards. Default/server/highres pass 265/262/265 selected roots without skips.
All 32 frozen captures and static mapped-memory checks pass on all three targets.
Source fingerprints are unchanged across each run.

Original-C client-state/notices qualification is complete. See
client-game-state-c-qualification.json. All 1,334 tracked production files match
qualified 0d5a03b2, so its production builds/ABI, exact known full-suite result,
and headless gameplay/save-load evidence are reused. Fresh production qualification
is required after native conversion. Fifth audited cache cleanup removed 32
superseded test archives /2,155,476,020 bytes; scripts are consumed. All baseline
and cleanup jobs are joined. Next: native conversion; do not stop at this checkpoint.

## Native client-state/notices conversion in progress

Baseline **beb38b24** is pushed. The selected 42 cases now route through private
Go before the remaining production C fallback. The full notice dispatcher is Go,
using the existing inform.c localization namespace, formatter, audio and dialog
owners. The fade-object setting is Go-owned, including the advanced-video pointer.
Thirty-six obsolete C interfaces retire with their last C callers; tests call Go
directly. An exact whole-repository search finds no remaining C-body/cgo uses of
those interfaces. The native manifest retains the remaining required callbacks
and checks that retired symbols are absent from fresh production binaries.

Implementation review preserves full-wire versus masked identifiers, equip-only
input rewriting, disconnected health-change queuing, signed health deltas/current
health, unsigned creation/camera coordinates, light fixed-point arithmetic,
secret-wall signed lookup, NPC reset/allocation behavior, and actual owner hooks.
The existing particle-light helper is reused instead of the similarly named
Drawable method. Notice kind 15 returns zero for a missing terminator within the
supplied slice; the former C strlen read outside that valid message domain. This
reversible invalid-input choice needs no user decision and remains for later review.

Working C count is **7,897 lines /37 files**, zero reference C (923 lines removed).
This is not yet qualified. The first compile found an unused import left by adapter
retirement; it is removed. Static mapped-memory preflight passes. Default native
frozen/affected qualification is running; remaining target and production gates
are pending. Ignored installation/retirement scripts are consumed.

Native default qualification passes all 265 selected roots, no skips, all 32 frozen
message captures and static checks (73254 joined). The C fallback has no callers
that bypass the updated Go adapter. Production ABI checks explicitly include both
remaining dispatchers and the animation callback used by lifetime handling.
Completed original-C target captures now share verified identical storage: 401
copies /1,070,110,014 bytes reclaimed, with all paths and contents preserved.
The original-state deduplication plan is consumed; use fresh output directories.

All initial native target sweeps pass 265/262/265 roots, no skips, all frozen
captures and static checks. Subsequent source review caught a fixture gap:
unknown spell titles reach the C formatter as NULL and print "(null)". The initial
Go draft discarded the lookup boolean and used an empty title. A standalone
386 probe of the actual noxstring.c formatter confirms "Cast (null)". The native
implementation now distinguishes unknown spells from valid empty titles. Seven
independent title boundary cases and three unterminated-notice cases are added;
no original frozen capture is changed. Final three-target qualification is repeated
on this corrected source before production gates. These extra contracts live in
the native-only test selection; the committed original-C selection stays intact.

Final native sweeps pass 266/263/266 roots, no skips, all 32 frozen captures and
static checks, with matching source fingerprints. Production gates are running.
Final target captures share verified identical storage with C references: 603
copies /1,605,538,655 bytes reclaimed; the preceding superseded native runs also
share 603 copies /1,605,538,655 bytes. Each cleanup has a separate consumed plan.

Next candidate audit (not an accepted/frozen selection): remaining client kinds
72–164 comprise 55 labels /41 whole groups /1,056 C case lines with no incoming
or outgoing goto edges. This connects progress/equipment/winner reports with
visual effects. Five additional live C helpers need consideration alongside the
cases: inventory-name initialization forwarding, white-flash state, player/team
winner score adjustment, and map-generation progress rendering. The latter also
has a Go map-population caller. Include these real dependencies rather than
assuming every remaining dispatcher call already targets Go. Ignored audit files:
remaining-client-groups.json and remaining-progress-effects-helper-audit.json.

## Qualified client-state/notices checkpoint

All final jobs are joined, including production 43757. Default/server/highres
pass 266/263/266 affected roots, without skips, and all 32 frozen captures. Ten
new independent boundary cases pass. Every final target and production gate uses
the same source fingerprint, including the new Go files. Three fresh 386/SSE2/CGO
binaries pass ABI checks: 36 retired interfaces and the former fade global are
absent; actual remaining callbacks and C dispatch fallbacks remain available.
The full asset suite matches exactly all 1,553 known failures and 15 pass /3 fail
/32 skip package results. Character creation/gameplay and explicit save/load pass
against the drawable references, including successful continuation after reload.
See client-game-state-native-qualification.json for commands/artifacts and hashes.

Physical C is now **7,897 lines /37 files**, zero reference C, a reduction of 923.
No original frozen expectation was changed to make the port pass. The next work
is the progress/effects candidate audit above; do not stop after this checkpoint.

Completed scenario asset copies were verified against originals before removing
1,112,747,701 duplicate bytes. Their deduplicated-assets.json manifests and ignored
deduplicate-client-state-assets.py restore command preserve recovery; generated
saves, screenshots, logs and changed files remain. The cleanup is complete and its
plan consumed. Original assets/archive are unchanged.

## Client progress/effects baseline development

Qualified client-state conversion **f9268ef1** is pushed. The next selected scope
has 55 client labels /41 whole groups /1,056 case lines plus six live C helpers
/180 lines, **1,236 body lines** total. The compass-image initializer belongs with
map-progress rendering. Its Go video caller and the map-population progress caller
move together with those helpers. BlueSpark/VioletSpark named globals also have
their last C consumers in this selection; their existing Go initializer and test
owner must move with them. Other mapped progress state remains shared through its
existing addresses. Selection/hashes: client-progress-effects-selection.json.

Player stats pass 768 C cases and inherited inventory name/stats checks. The
fixture verifies every byte of the player record, nonoverlapping field writes,
read-only input, missing players, unsigned values, signed display levels, and the
host-mode distinction between suppressing record writes and refreshing the title.
Lesson, experience and treasure fixtures also pass; no whole-scope baseline is frozen.
Production is unchanged from f9268ef1 and will reuse its qualification for the
next test-only baseline.

The partial progress-report checkpoint now covers 7,820 cases, repeated in three
separate processes with identical four capture hashes and seven completed roots
per run, without skips. Lesson elimination compares signed values while storing
full-width scores. Treasure updates require a strictly newer unsigned frame, but
announcements use the resulting stored values even when host/stale-frame gates
suppress writes; count/max arithmetic wraps. Whole-record, input immutability,
text-ring and audio contracts accompany captures. All jobs are joined. See
client-progress-reports-c-checkpoint.json. C remains 7,897 physical lines; this
checkpoint adds fixtures only. Continue inventory and other selected coverage.

The next fixture checkpoint adds pickup (576), inventory notifications (432),
generator status (1,536) and modifier reports (4,096): **14,460 total new cases**.
Three combined repeats each pass eleven roots without skips and eight identical
capture hashes, with fixed source. Static checks pass. Pickup covers connection
and host modes, masked IDs, stacks, modifier descriptor order/missing IDs and real
allocation exhaustion, including exact failure-report bytes in each queue.
Notification cases cover drops, valid/invalid equipment and missing items.
Generator cases verify full drawable records, countdown rising edges and lighting
flag clearing. Modifier reports cover every byte ID, static/dynamic match and
mismatch, missing objects, exact writes and untouched tail canaries.

Fixture corrections: host mode uses the ordinary client message list, client mode
the reliable queue; RedApple cannot equip; Equipped is uint32. A local repeat
pattern initially had two lines and was rejected before execution; the corrected
one-line union ran all eleven intended roots. These are fixture/tooling fixes,
with no production changes. All jobs joined. See
client-progress-inventory-c-checkpoint.json. C remains **7,897 /37 files**.
Continue remaining selected coverage without pausing at this checkpoint.

## Green-bolt correction during effects baseline development

The new message-152 contract found a typed-pointer arithmetic bug: the C expression
adds 432 to a `nox_drawable*` returned by spriteLoadAdd, advancing 432*512 bytes.
Its intended 13-byte effect record at byte offset 432 remained untouched. Original
probe output and source are retained under progress-effects-third/fourth and Git.
The handler now checks allocation success and advances 432 bytes. Failure consumes
11 bytes without writing an effect. This is a deliberate correction of undefined
writes under the standing reversible-decision policy, recorded in DECISIONS.md.
No old frozen golden was changed. Current C is 7,901 lines (+4); fresh production
qualification is required before freezing this corrected baseline.

Continued owner audit also selected the named map-frame gate
`dword_5d4594_1200804`: its last C reads are the four winner messages, while
network_client.go reads/writes it for map-use/endgame handling. Move that owner
with the winner messages, preserving those Go callers. This raises the selected
named-owner count to three; it does not add C function bodies.

## Qualified green-bolt correction and effects checkpoint

All jobs are joined. Default/server/highres pass **282/279/282 roots**, no skips,
47 exact message captures /75,276 cases, with identical final source fingerprints.
The seven new effect roots add7,922 cases: smoke384, delta-Z6,144, white flash8,
arrow trap288, green bolt432, point sparks288 and ray effects378. Added to the
preceding14,460 progress/inventory cases, the next scope has22,382 captured cases
so far. The previous32 frozen expectations remain unchanged. See
[green-bolt-correction-qualification.json](green-bolt-correction-qualification.json).

Fresh three-target production binaries and ABI checks pass. The full asset suite
matches exactly1,553 known failures (15 pass /3 fail /32 skip packages). Headless
character creation/gameplay and save/load pass against the drawable references,
including continuation after reload. This qualifies the small C correction and
its tests; it does not claim complete baseline coverage or translation of the
remaining progress/effects scope. C is **7,901 lines /37 files /zero reference C**.

Fixture development found and corrected a canary placed in drawable list-link
word107; effect checks now start at byte432 and retain trailing canaries and whole
record comparisons. Named BlueSpark/VioletSpark backing-blob checks use the existing
explicit owner, satisfying static-memory checks. No production expectation was
weakened. Point-spark/ray dispatch comparisons reuse already-qualified particle
helpers and additionally check message fields, allocation schedules and ownership.

Storage: a sixth verified cache cleanup removed32 superseded porttest archives /
2,219,981,278 bytes. Final default/server/highres capture deduplication preserved
all paths and bytes while sharing617,563,932 /616,816,664 /617,563,932 bytes.
All those plans are consumed. Original assets/archive and all qualification
binaries, captures, saves and screenshots remain. Continue the remaining scope;
particle-bursts-draft.go is an uninstalled draft for messages150/163.

## Further particle, stream, summon and shield C fixtures

The next test-only checkpoint adds **1,908 cases**: ricochet/mana-bomb bursts540,
update-stream framing360, summon start/cancel720 and shields288. Three independent
repeats each pass eleven effect roots /9,830 cases, no skips, identical capture
hashes and source fingerprints. The seven previously qualified effect hashes are
unchanged. Static checks pass. New-scope coverage totals **24,290 cases** so far;
remaining production C stays **7,901 lines /37 files**. See
client-effects-second-c-checkpoint.json. This is not full-scope qualification.

Contracts include exact burst RNG use/fields and frame wrapping; float32 radius
narrowing before truncation; stream consumed-length/terminator/early-return behavior
and alias reports even when disconnected; full summon owner IDs and cancellation
of real parent/child allocations; and shield direction, duplicate suppression and
static/dynamic owner lookups. A summon snapshot initially omitted its unlinked
child and failed the existing ownership check; it now supplies that actual owner.
Message164 is update-stream dispatch (a checkpoint label calling it earthquake
was corrected). No production behavior or prior expectation changed.

The completed green-bolt gameplay/save copies now have verified restoration
manifests; 1,112,747,701 bytes of original-identical assets were removed. Use
`deduplicate-green-fix-assets.py --restore NAME` for those historical copies.
Original assets/archive, changed files, saves/screenshots and evidence remain.
All cleanup plans and installed drafts are consumed. Continue with duration-effect
message158; its ignored draft is not yet installed.

## Client progress, winner reports and effects — complete

Baseline a2b8342c freezes all55 selected client labels /41 whole groups and six
helpers:63 captures /122,204 cases, including69,310 cases for this scope. Native
qualification passes378/374/378 default/server/highres roots with no skips and
identical source, all frozen hashes unchanged. Fresh production/ABI checks pass;
the full suite exactly matches1,553 known failure entries and15 pass /3 fail /32
skip packages. Headless gameplay and explicit save/load/continuation both pass.
See [the qualification](client-progress-effects-native-qualification.json) and
[the reproducible manifest](client-progress-effects-native-batch.json).

The conversion moves stats, equipment/progress reports, individual/team winners,
creatures, particles/rays, stream dispatch and map-progress rendering. Six private
C helpers disappear,32 Go-backed C exports are retired, and BlueSpark, VioletSpark
and the map-frame gate move to Go storage. Test adapters invoke Go; there is no C
reference implementation. C remaining: **6,641 physical lines /36 files**, down
**1,260 lines**, zero reference C. Existing client/server fallback behavior stays
in production for the next batches.

Review notes: wide sentry coordinates preserve differing unsigned audio and signed
spark distance arithmetic; winner fixtures cover score underflow, frame gates and
UI/clock ordering. Retiring callback addresses required retaining nil table entries
to keep sorted capture IDs stable. No expectations were regenerated. The earlier
green-bolt correction remains the separately documented intentional behavior fix.

Next, finish the remaining server action dispatcher as one coherent batch, reusing
existing C contracts and extending positive action coverage. The old fixture drafts
and installation/retirement scripts referenced above are consumed.

## Remaining server actions: original-C baseline

The complete remaining server dispatcher has16 top-level cases /476 physical C
lines, plus26 lines in two private pickup helpers. Its baseline passes
**518/514/518 roots** across default/server/highres, no skips, with identical
source fingerprints and all **75 frozen captures**. Twelve new captures cover
**2,716 cases**, including creature commands, spell queues and friendly targets,
collision callbacks, inventory failure, secondary weapons, book lookup priority,
gauntlet respawn, shop admission and trade transactions. Separate contracts cover
vote withdrawal and gauntlet leave routing.

Production remains identical to **f8607de2**:6,641 physical C lines in36 files,
zero reference C. Its production/ABI, full-suite and headless save/load evidence
is reused because only porttest fixtures and documentation changed. See
[the C qualification](server-actions-c-qualification.json) and
[the batch manifest](server-actions-c-batch.json).

Decision for review in the native conversion: reject incomplete actions before
reading their fields. The current caller checks consumed size only afterward.
Complete messages retain the frozen behavior; separate native-only contracts
will check every incomplete prefix without accessing gameplay owners.


## Server actions: native conversion qualified

All16 dispatcher cases and two private pickup helpers now run in Go. The native
path preserves explicit player/update owners, full-width dynamic object codes,
raw inventory-failure codes, spell queue ordering, trade admission and transaction
ordering, alias bytes, vote routing and gauntlet respawn behavior. Thirty-eight
unused C exports are retired; remaining C callers keep their required interfaces.
No C algorithm is retained solely for testing.

Default/server/highres pass **567/563/567 roots**, no skips, with identical source
and all **75 frozen captures /124,920 cases** unchanged. Three fresh production
binaries pass ABI checks. The full suite matches1,553 known failure entries and
15 pass /3 fail /32 skip packages. Headless gameplay and explicit save/load
continuation match their references. See [qualification](server-actions-native-qualification.json),
[manifest](server-actions-native-batch.json) and [retired interfaces](server-actions-retired.json).
C remaining: **6,139 physical lines in35 files**, down **502 lines**, zero reference C.

Decision for review: reject incomplete actions before field/owner access, since
the caller validates consumed size only after dispatch. Independent contracts
cover every incomplete prefix of34 formats; complete messages preserve the C
baseline. Interface retirement also required a fixture-only normalization fix:
nil callback addresses receive no ID, while reserved slots preserve dynamic IDs.
Frozen expectations were not regenerated. A separate replay-buffer copy defect
was identified and queued for its own regression/fix after this conversion.


### Replay adapter correction

After server-action qualification, a focused regression showed that the replay
adapter allocated zero-filled storage without copying the recorded message. It
now clones the input onto the C heap. Empty input bypasses allocation and reaches
the existing empty-message handler; previously it panicked on zero allocation.
Two regressions pass on all three profiles with real player owners, concatenated
alias payloads, frame updates and unchanged input. See [focused results](replay-copy-qualification.json).
C remains6,139 lines /35 files /zero reference C. This is a separate reversible
correctness fix; full-suite/scenario evidence belongs to the preceding conversion.
