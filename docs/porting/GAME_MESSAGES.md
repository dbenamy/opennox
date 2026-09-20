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
