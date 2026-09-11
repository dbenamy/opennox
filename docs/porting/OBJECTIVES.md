# Objective objects and obelisk recharge

Next connected candidate: 15 functions / 908 C lines. Obelisk, flag, ball and
crown updates; flag pickup and identity helpers; ball owner/collision/home-base
scoring; crown pickup dispatch; ball reset/spawn; flag-mode dispatch and shared pickup buff removal. Keep unrelated map-mode setup
with the later map initialization family. Original bodies/addresses in scope.json.

Reuse world/temporary/effects/inventory fixtures and actual guarded players,
teams, object lists, script/sound/network observers, creation/deletion guards,
spatial queries and clock helper. Add optional objective fields only; preserve
all 22,317 previous focused contracts. Retain actual inventoryCrownPickup,
resource, team, camera, movement and spatial dependencies.

Clock: PlatformTicks can be replaced only inside the test fixture and restored;
C ABI returns unsigned 32-bit ticks even though ball UD stores a uint64 timestamp.
Cover elapsed boundary 20,000, 32-bit wrap and 64-bit stored timestamps. Preserve
that widening/subtraction behavior in native Go. All FPS inputs remain legal;
obelisk requires FPS >= 2 because of division by FPS>>1.

Provide named GameBall/GameBallStart/HomeBase/Flag definitions, guarded ball UD,
owner and player iteration lists, actual team membership, color/name lookup table,
score counters, game-mode flags, obelisk energy/max-mana configuration, and
missing-type/empty-start cases. Save/restore lists, type caches, timestamps,
network state and all touched globals. No algorithm copies in test C.

Test exact/strict deadlines and wrap; owner/dead/stale possession transitions;
nearest/blocked movement, reset spawn selection and seeded RNG; ball pickup,
team possession and home-base score branches; flag friendly return/enemy pickup
and score flows; crown active/inactive owner and pickup dispatch; obelisk class,
range, LOS, energy, HP/mana capacity and recharge cadence. Assert positive
movement, possession, score/messages, spawn counts, mana transfer and energy
conservation. Compare full state/return/event captures from repeated original C,
lock and push before porting. Convert callers/helpers together, qualify once,
update C_LOC/docs, commit/push, summarize and continue.

## Audited original-C scope

- 00417F50: `int sub_417F50(int a1)` (78 lines).
- 004EA490: `void nox_xxx_pickupFlagCtf_4EA490(int a1, int a2)` (124 lines).
- 004EB9B0: `int sub_4EB9B0(int a1, int a2)` (18 lines).
- 004EBA00: `void nox_xxx_collideBall_4EBA00(int a1, int a2)` (58 lines).
- 004EBB50: `int sub_4EBB50(int a1, int a2)` (13 lines).
- 004EBB80: `short nox_xxx_collideHomeBase_4EBB80(int a1, int a2)` (80 lines).
- 004ECBD0: `int sub_4ECBD0(int a1)` (12 lines).
- 004ECC00: `int sub_4ECC00(const char** a1)` (26 lines).
- 0053C580: `signed int nox_xxx_updateObelisk_53C580(int a1)` (177 lines).
- 0053DDF0: `int nox_xxx_updateFlag_53DDF0(int a1)` (31 lines).
- 0053DF40: `void nox_xxx_updateGameBall_53DF40(int a3)` (79 lines).
- 0053E1D0: `void nox_xxx_updateCrown_53E1D0(int a1)` (36 lines).
- 004EA400: `void sub_4EA400(int a1, int a2)` (26 lines).
- 004EA7A0: `int sub_4EA7A0(int a1)` (22 lines).
- 004EA800: `short sub_4EA800(int a1, int a2)` (128 lines).

## Original-C baseline

Before conversion, nineteen complete capture groups, 1,920 cases, repeat
byte-for-byte in independent runs. Hashes are locked in
`src/objectives_porttest_test.go`. The original-C regression run passes all
24,237 accumulated focused cases in 70.385s, including the previous 22,317.
Evidence: `build/port-objectives/c-{first,repeat}-objectives-*.json`,
`final-{first,repeat}.log`, and `c-regression.log`.

The fixture uses actual player and team records, membership lists, inventory,
spell definitions, movement, line-of-sight and object creation. It isolates
object lists, team/type/net-code caches, score limits, status globals, player
state and the tick source, and records transport bytes. Team records have
stable low address bits for retained short return values. Original C still
performs every algorithm under test. All borrowed buffers and player state are
restored; guarded allocation and object cleanup checks remain enabled.

Coverage includes flag identity/deadlines, CTF friendly return/enemy pickup and
scoring, flagball/home-base scoring and owner dispatch, ball pickup/possession/
reset, crown dispatch/following, obelisk idle regeneration/mana/wand recharge,
class/team/death eligibility, blocked and clear sight, and removal of eligible
pickup buffs. Tests include unsigned frame wrap and the C tick API's 32-bit
narrowing before uint64 subtraction. Independent assertions check mana/energy
conservation, ownership, flag position/deadline, CTF score increments and return,
line-of-sight denial, and removal versus preservation of configured buffs.

## Native conversion

Original-C baseline commit `d0094d3a` was pushed before conversion. All 15
functions now live in `legacy/objectives_update.go`, `objectives_scoring.go` and
`objectives_exports.go`; the ball-death handler calls the Go reset owner directly.
The string-table helper's C declaration drops const to match the exported Go ABI;
the pointer representation and read-only behavior are unchanged.

All 1,920 cases / 19 full captures match the original C byte-for-byte in
`native-final-objectives-*.json` (6.585s). The first comparison caught routing of
team-score notification through an existing Go transport entry point; retaining
the original C team-score dependency preserves the exact notification path and
ordering. Native callers use the converted inventory/effects owners directly.

The conversion removes 908 physical C lines: 124,395 production C lines / 149
files / zero reference C. Accumulated port tests, including all 24,237 focused cases, pass in default,
server and high-resolution variants: 145.833s / 122.632s / 124.291s. All three
production binaries verified ELF32/i386, SSE2, CGO enabled. Full-suite failure
multiset exactly matches the existing baseline: 1,553 entries, 15 packages pass,
3 fail, 32 skip. A fresh unchanged repeat-a headless gameplay replay passed in
57.464s. Evidence: `build/port-objectives/qualification.json`, variant logs,
production binaries and `build/baseline/runs/objectives-port`.

Next candidate: the connected player-attack, hit, shooting and reload family;
12 functions / 1,105 C lines. Preserve the same baseline-before-conversion and
one-qualification-per-batch process.
