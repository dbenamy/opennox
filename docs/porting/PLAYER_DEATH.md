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

Thirteen default-target roots pass, with **1,053 capture records** (1,066 test entries).
This is an in-progress checkpoint, **not a frozen or qualified C baseline**.
The only production change is the two-line absent-team guard. Tests, the original
failure and current capture hashes are recoverable from the checkpoint and
[player-death-checkpoint.json](player-death-checkpoint.json). Repeated captures, all-target affected-corpus checks and fresh production
qualification precede the native conversion. All three focused target repeats passed with identical source and capture bytes
under build/port-player-death/c-{default,server,highres}. Frozen literals are
installed; affected-corpus and fresh production qualification remain.

## Extended contracts

The actual registered PlayerDie callback covers recent-assist expiry (strictly less
than ten seconds), frame wrap, missing/inactive candidates, self/killer exclusion,
online source attribution and weapon-style reset. Separate C-allocated monster and
projectile objects preserve a valid registered player roster. The initial fixture
incorrectly changed a registered player into a monster and interrupted iteration;
this was corrected without changing production. alloc.New supplies zero storage,
so fixture class/type fields are assigned explicitly after allocation.

Quest contracts check life decrement and exhaustion, fractional starting lives,
gold penalty, statistic/reset messages and exact RNG consumption. Ability cleanup
starts with active execution lists, cooldowns and nonzero enchantment arrays.
Cooperative death cancels pending character loading, including frame wrap.

Crown contracts cover held/unheld objects, teams, friendly deaths and transfer
disablement. Preserve the historical asymmetry: the unteamed enemy branch can
assign a dropped crown to the killer; the teamed scoring branch does not. The
second drop argument is a pending-owner pointer, despite the existing adapter's
parameter name `stamp`. Normalize only that identified pointer in captures.

The affected selection includes 317 roots from death/controls, inventory/respawn,
object state, roster/team/gameplay reports, quest/spell lifecycle, statistics,
objectives/rewards, session and orchestration owners. Focused capture repeats
remain 13 roots; broader checks follow frozen literals. Static memory checks pass.

Disk cleanup verified twelve superseded map-section/statistics binaries against
successful reports and recorded hashes, reclaiming 584,677,176 bytes. The latest
map-metadata production binaries remain. Audit and apply are consumed at
build/port-player-death/cleanup-binaries.py; its JSON manifest remains.
