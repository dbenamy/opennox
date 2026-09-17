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

## Reproduction and next gates

Source `build/baseline/env.sh`, then use `tools/porting/run_batch.py` with
`docs/porting/world-collisions-batch.json`, one of `c-default`, `c-server`,
`c-highres` or `c-repeat`, and a fresh output directory. Frozen hashes live in
both the tests and manifest. Baseline evidence is under
`build/port-world-collisions/c-*`; the production reuse audit is
`c-production-reuse-audit.json`. Raw diagnostics stay ignored.

Next: translate the connected block, retire only the unused mass-exchange C
interface, preserve the other callback/caller interfaces, then run the broadened
consumer suites and production integration. Expected affected scope is 444 roots /
40,751 leaf cases and 126 groups / 42,128 records; verify counts from the actual
qualification logs. Update C LOC and commit/push the qualified conversion.
The one-shot freeze script is consumed. Installed exit fixture drafts are stale.
