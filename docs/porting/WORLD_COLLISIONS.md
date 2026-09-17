# World collisions and interactions

## C baseline

Baseline after qualified reliable-queue conversion `fec78429`: **21 live functions /
1,016 physical C lines** in GAME3_3.c, 004E86E0 through 004EBF40, ending before
004EC520. The count includes an obsolete commented projectile implementation.
Every selected entry has live callers or callback registration. Production remains
**54,962 C lines / 82 files**, zero reference C. No production correction was needed.

Default, server, highres and an independent default repeat each pass **28 roots /
10,515 unique leaf cases**, without skips. All **28 groups / 10,541 frozen records**
match. All four gates have the same **2,013 source fingerprints**. Wall times were
35.77s, 118.74s, 44.85s and 6.54s respectively.

The source audit proves all 1,991 preceding source files are unchanged; the 22
additions are guarded by `porttest`. Reuse the production qualification from
`build/port-reliable-reports/native-production`: all three builds and ABI audits,
exact known full-suite failures (1,553 failure entries), gameplay (41 frames),
actual save/load (7 frames) and flat rendering (14 frames), with exact map
regeneration. Reused client SHA-256:
`a501fc9e201c9df1880ac25ef09b4552066f3637245dbe8b6bffbd19e8081b34`.
This is documented production-identity reuse, not a fresh production run for this
baseline. Fresh production qualification is required after conversion.

## Contracts and fidelity notes

Fixtures use real server/player owners, C-owned object records, inventory and
factory allocation, map indices, script VM callbacks, deferred audio, message
queues, spell delivery, observer transitions and queued map/save requests.
The damage callback observer records arguments only; it contains no damage
algorithm and is excluded from production. The tradable-ankh fixture uses the real
registered pickup. C callback identities remain necessary where other C compares
addresses; migrating the callback registry is outside this batch.

- Physics covers floating-point velocity stores, equal-mass exchange and independent
  momentum checks. Both-zero mass is outside the defined division domain.
- Doors cover four actual key types, adjacent-door updates, magic-lock expiry,
  missing-key wire encoding and exact clock throttles. Valid door angles are
  0/8/16/24; the original leaves adjacent coordinates uninitialized otherwise.
- Pickup covers unsigned frame age, signed FPS shifts, admission flags and real
  inventory insertion. Chests cover admission, silver-key consumption, callbacks,
  inventory unlinking and actual pending-world placement.
- Exits cover readiness, host exclusion, stage thresholds, countdowns, ordinary
  map/waypoint requests, quest observer/warp transitions, partial-party waiting,
  glyph cleanup and saturating trap counts. Cooperative exits queue a WORKING save
  with the portal and frame; actual later saving/loading is covered by production
  integration. Map selection uses the real quest catalog.
- Projectiles cover target selection, wall reflection, shield/sword and optional
  behavior flags, actual inversion modifiers, ownership, RNG, deletion and real
  cure-poison delivery. Scripts exercise direct and registered callbacks, disabled
  blocks, arm/deadline state, geometry and ability-based trap admission.
- Spell awards, soul gates and extra lives cover state and messages, identity
  history, expiry boundaries, a full player history and the 64-entry ring wrap.

The C platform clock adapter truncates to 32 bits before unsigned 64-bit throttle
subtraction/storage; preserve that boundary. Countdown writes use the real root
server clock: tests bracket it and normalize only the deadline timestamp. Fixtures
restore shared scratch state. Development fixes were fixture-only: the trap class
is Door, UseData is a pointer wrapper, player classes must be restored between
cases, non-player teleport targets need a separate allocation from the player
roster, and the save-portal flag needs an explicit fixture owner.

## Native implementation review

The Go implementation replaces the entire 1,016-line block with six production
files. Twenty C entry points remain for live callers/callbacks; the charge caller
uses the mass-exchange Go helper directly and its C interface is retired.
Qualified production C is **53,946 lines in 82 files**, zero reference C.
The six production Go files contain 773 physical lines.

The first focused comparison matched 27 groups and exposed a floating-point
rounding difference in mass exchange. The qualified original client instructions
confirm that the mass sum and doubled-second-mass coefficient stay in x87 PC53
registers, whereas the second object's difference coefficient is stored to a
32-bit temporary before reuse. Native Go explicitly reproduces those boundaries.
The original 175 frozen physics records and independent momentum/equal-mass
contracts remain unchanged. Instructions and failed-run diagnostics are retained
under build/port-world-collisions; no C algorithm is retained for testing.

The review also preserves the platform clock's uint32 truncation, signed countdown
reduction, short stage-message arguments, raw UTF-16 identity comparison, string
copies only through their terminator, and reads after gameplay callbacks. Invalid
door angles formerly used uninitialized adjacent coordinates; Go initializes them
to zero, while the valid four-angle domain remains the compatibility contract.
Projectile argument fields that C left uninitialized are zeroed; its target object
continues to be supplied to the actual spell delivery path.

## Frozen groups

| Group | Records |
| --- | ---: |
| ankh-history | 96 |
| audio-frames | 504 |
| chest-contents | 80 |
| chest-key | 16 |
| clock-throttle | 168 |
| coop-exit-save | 9 |
| countdown | 100 |
| door-keys | 64 |
| door-magic | 180 |
| door-missing-key | 8 |
| exit-admission | 30 |
| exit-glyphs | 24 |
| mass | 175 |
| pentagram | 6 |
| pickup-admission | 1728 |
| pickup-integration | 1 |
| quest-exit | 8 |
| quest-readiness | 6000 |
| scripts | 144 |
| sign | 96 |
| soul-gate | 480 |
| spell-award | 120 |
| spell-projectile | 96 |
| spell-wall | 6 |
| teleport | 24 |
| trap-ability | 10 |
| trap-geometry | 240 |
| undead-damage | 128 |

## Native qualification

Default/server/highres each pass **444 roots / 40,751 unique leaf cases**, without
skips, and match **126 frozen groups / 42,128 records**. Wall times are 219.61s,
321.76s and 261.40s. All three share the same **2,019 source fingerprints**; see
`build/port-world-collisions/native-coverage-audit.json`. The focused retry passes
all 28 groups in 119.94s. No expected hash was changed for the Go implementation.

All three fresh production builds pass their ELF32/i386/SSE2/CGO and symbol audits.
The full asset suite has exactly the known 1,553 failure entries, with 15 passing,
3 failing and 32 no-test packages. Fresh headless gameplay (41 frames), actual save/load (7 frames) and flat rendering
(14 frames) all pass against the prior reliable-queue native references. The flat
scenario removes 51 uncompressed maps and verifies exact regeneration of the
loaded map. Full production qualification takes 371.69s, with the same 2,019
source fingerprints as the affected sweeps. Client SHA-256:
`bd58424e6816bd0b5c1690ebcb2ca21464bd5f8f3b0e7c24427f4a27c039637c`.
Evidence is under `build/port-world-collisions/native-production`.

## Reproduction and next gates

Source `build/baseline/env.sh`, then use `tools/porting/run_batch.py` with
`docs/porting/world-collisions-batch.json`, one of `c-default`, `c-server`,
`c-highres` or `c-repeat`, and a fresh output directory. Frozen hashes live in
both the tests and manifest. Baseline evidence is under
`build/port-world-collisions/c-*`; the production reuse audit is
`c-production-reuse-audit.json`. Raw diagnostics stay ignored.

The native phases are `native-default`, `native-server`, `native-highres` and
`production`, each with a fresh output directory. The conversion is complete;
next candidate is quest runtime/statistics and difficulty scaling. The one-shot
freeze script is consumed; installed exit fixture drafts are stale.
