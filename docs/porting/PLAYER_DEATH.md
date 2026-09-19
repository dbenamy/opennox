# Player death, scoring and corpse creation

## Scope

Seven connected C bodies /677 original body lines: registered PlayerDie handler,
arena/elimination/king-of-the-realm scoring, death notification, respawn corpse
creation and its type-cache initializer. Parent **00131036** is committed and
pushed (26,836 C lines /67 files). Caller audit is currently in
build/port-player-death/selection-candidate.json; its original offsets precede
the C prerequisite below. Refresh positions before using an installer.

The death callback requires a minimal C-to-Go export while its registration uses
a C function pointer. Move the native respawn caller directly to Go; private
scoring/cache helpers need no exported interfaces. Actual callers discard the
historical return values; contracts capture owned state and observable effects.

## C prerequisite — review later

Original C crashes when an unteamed player kills a teamed player in arena mode.
The ordinary scoring branch reads the absent killer team's score before passing
it to the already nil-tolerant Go team service. The isolated contract confirms
SIGSEGV in `nox_xxx_playerUpdateScore_54D980` for victim-team=1/killer-team=0.
Evidence: build/port-player-death/original-teams/tests.jsonl; the test process
terminated and is joined. No production/gameplay process was involved.

Guard only the team-score update when the killer has no team. Keep the player
score increase and all other scoring behavior. This adds two C lines before
freezing; corrected initial arena/corpse-cache contracts pass. A fresh C production qualification
is required because production changed. No expectations have been frozen yet.

Arena suicide/environment paths subtract lessons without incrementing the death
counter. The subtraction service was inspected and does not add that increment.
Preserve this historical rule explicitly; this port does not change scoring policy.

## Contracts in progress

Use the existing real match-roster owner for players, teams and score messages.
The initial arena matrix covers both assigned and unassigned teams. Corpse-cache
contracts use shipped direction suffixes and real definitions, missing definitions,
initial cache values, untouched direction 4, and surrounding storage.

Planned remaining coverage: score/death wrapping, assistants and their distinct
player-info gate, game modes/crown ownership, actual message bytes and order,
registered death dispatch and cleanup, quest lives, recent-assist frame wrap,
corpse creation positions, factory failure and per-object decay RNG consumption.
Reuse the existing controls corpus for integrated respawn behavior.

## Recovery

No native translation is installed. Initial C bridge and contracts are tracked as
working changes. Original bridge/blob/cache draft copies under build/port-player-death
are consumed; actual source takes precedence. Do not edit source while tests run.

Additional fixture review: restore the temporarily non-player actor's class before
calling the typed player snapshot API. Crown team-award notifications precede
player-score notifications, unlike arena's ordering; the expectation was corrected
against C. The real wrapper ability manager must be initialized for death cleanup.
These are fixture corrections, not additional production changes.

## Recovery checkpoint

Seven default-target roots pass, with **608 capture records** (615 test entries).
This is an in-progress checkpoint, **not a frozen or qualified C baseline**.
The only production change is the two-line absent-team guard. Tests, the original
failure and current capture hashes are recoverable from the checkpoint and
[player-death-checkpoint.json](player-death-checkpoint.json). Remaining lifecycle,
quest, assist timing and crown-transfer coverage plus repeat/all-target/production
qualification precede the native conversion. No source/test sessions are active.
