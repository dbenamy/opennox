# Game-message dispatch and player notices

Qualified parent: `0d5a03b2` (pushed), 8,820 C lines in 38 files, zero reference C.
Production algorithms are unchanged. This is original-C baseline development;
no expectations are frozen and no conversion is installed.

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

Build the original-C baseline by behavior family, then qualify the combined scope.
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

## Current work

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
