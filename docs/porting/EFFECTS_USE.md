# Modifier effects and weapon use

Connected batch: **41 functions / 977 physical C lines** converted to native Go.
Equipment was the preceding batch. C count falls from 128,081 to **127,104 lines /
149 files / zero reference C**. Original-C baseline `0738dbed` was committed and
pushed before conversion.

## Scope

- src/legacy/GAME3_2.c: 004DFB50, 004DFB80, 004DFBB0, 004DFC30, 004DFCA0, 004DFD10, 004DFD40, 004DFD80, 004DFDB0, 004DFDE0, 004DFE10, 004DFE40, 004DFF40, 004E0040, 004E0140, 004E0170, 004E01D0, 004E02C0, 004E0370, 004E0380, 004E03D0, 004E03F0, 004E0480, 004E04C0, 004E04D0, 004E0640, 004E0670, 004E06F0, 004E0740, 004E07C0, 004E0850, 004E08E0, 004E0960, 004E09B0.
- src/legacy/GAME4_3.c: 0053C520, 0053C940, 0053F290, 0053F480, 0053F4F0, 0053F670, 0053F8E0.

The first family covers item engagement flags, inventory modifier lookup,
speed/protection, regeneration/replenishment, armor/damage modifiers, status
and resource effects, readiness and projectile speed. The second covers item
recharge/rate, fireball/wand projectiles, wand spell acceptance, fire-wand sparks
and use callback dispatch. Keep damage resolution and spell implementations as
production dependencies. The unused grip search has a keep.go reference under
the disabled `none` tag; preserve its behavior within the family for now.

## Testing plan

Extend the existing equipment/inventory/resource/shop fixture with an optional
EffectsUse input and snapshot; keep all existing 12,009 contracts unchanged.
Use action IDs starting at 500, a thin original-C dispatcher, and stable function
pointer identities. Capture actual modifier descriptors, guarded object/update/
use/health memory, buff state, spell acceptance arguments and outcomes, created
projectile/Spark state, packets, audio, RNG state and raw return bits.

Cover all modifier slots and inventory filters; null/class guards; speed and
protection rounding/clamp boundaries; frame-wrap timing; regeneration HP and
holder eligibility; replenishment byte wrap; signed effect arithmetic; status
application and notification order; mana/HP limits; readiness callback identity;
recharge saturation; wand empty-charge/cooldown checks; accept/reject outcomes;
player targeting versus nonplayer coordinates; single/triple fireball directions;
allocation failure; projectile velocity and spark RNG; use callback return values.
Keep original valid-input preconditions explicit (including nonzero divisors).

Repeat and compare full original-C captures, lock hashes, verify existing
contracts, then commit/push the baseline before conversion. Convert the whole
family and qualify once at its boundary: accumulated default/server/highres
corpora, three production builds, unchanged full-suite failure multiset, fresh
headless gameplay. Record actual C_LOC, commit/push, summarize and continue.

Detailed local audit: build/port-equipment/next-effects-use-scope.json and
build/port-effects-use/fixture-plan.md. No user question is pending.

## Original-C baseline

The equipment conversion was committed/pushed as `dec9b1ec`. This baseline
was captured while the entire family remained original C. **4,733 cases /
19 groups** are locked in
`effects_use_porttest_test.go`. Full captures c-tables and c-repeat match
byte-for-byte in every group (11.054s / 10.557s). All prior **12,009** equipment,
inventory, resource, shop and trade cases remain unchanged (37.208s).

- flags: 576 cases.
- arithmetic: 400 cases.
- grip-inversion: 246 cases.
- readiness-recharge: 240 cases.
- recharge-percent: 200 cases.
- protection: 936 cases.
- regeneration: 288 cases.
- replenishment: 160 cases.
- resource-transfers: 324 cases.
- wand-acceptance: 192 cases.
- use-dispatch: 16 cases.
- status-force: 120 cases.
- projectiles: 204 cases.
- fire-wand: 84 cases.
- inventory-lookup: 600 cases.
- speed-rounding: 63 cases.
- wand-classes: 64 cases.
- null-guards: 17 cases.
- fireball-wall: 3 cases.

Fixture files: legacy/effects_use_porttest.go, server/effects_use_porttest.go,
and blobdata/effects_use_porttest.go, with optional shared-fixture fields.
Action IDs 500..540 call all original functions through a thin dispatcher.
There are no copied production algorithms in the fixture.

The fixture supplies actual balance arrays and ForceWand/projectile/Spark
lookup definitions. Spark update/collision allocations have checked tail guards.
The spell-acceptance service observer records exact arguments and configurable
accept/reject outcomes. Player damage uses the existing recording callback.
Per-step captures include target and created object data, scalar arguments,
modifier descriptors and callback tables, packets/audio, cache words and state.

New projectiles belong to the enclosing lifecycle fixture: shop cleanup accounts
for exactly those objects, rejects duplicate ownership, and lifecycle cleanup
asserts zero leftovers. The fixture saves/restores the ForceWand/buff lookup
caches and the original six-row inventory modifier table at 587000+200160,
relocating exactly the six callback addresses used by production initialization.
An independent lookup assertion caught missing initialization during fixture
construction; the locked baseline includes working positive/negative lookups.

Independent assertions also check exact projectile origins for clear, blocked
and transparent walls. Regeneration includes successful non-armor healing at
frames 90/180. Readiness does not accept nil: its C code dereferences item data
before its apparent null check. Null tests cover supported inputs. Regeneration
and replenishment use nonzero divisors. Missing Spark lookup is exercised;
invalid type IDs are not used to simulate allocation failure because the retained
Go indexed allocator requires a valid definition and does not return a recoverable
nil allocation. The C wand-shot null-result branch is retained behavior to port,
but has no safely reachable fixture input through that existing allocator.

Local evidence: build/port-effects-use/{scope.json,original-c.txt,callers.txt,
c-tables.log,c-repeat.log,c-locked.log,c-dependencies.log} and corresponding
c-tables/c-repeat JSON captures. Obsolete intermediate captures were removed
for disk space; logs, final evidence and user assets remain intact.

## Native implementation

`legacy/effects_modifiers.go` owns inventory modifier lookup, engagement flags,
speed/protection, regeneration, replenishment, grip, status/resource effects and
readiness. `legacy/effects_weapon_use.go` owns recharge, projectile creation,
fireball/wand use and callback dispatch. `legacy/effects_exports.go` retains all
41 ABI entries for C callers, registered modifier addresses and the test dispatcher.
No original C algorithm remains solely for tests.

Inventory/AI use callers, shop recharge and poison resistance call native Go.
Three wand registrations use native Go functions with the original callback
addresses retained. Generic use dispatch still preserves raw integer callback
results; it cannot substitute the registry's bool result for arbitrary callbacks.

Native full captures match all 19 original-C groups byte-for-byte (10.836s).
The first implementation exposed one projectile-velocity difference: the original
compiler spills the X product across its second direction lookup but keeps the
Y product until owner velocity is added. Explicit Go float32/float64 operations
preserve this distinction. Speed/protection caps, signed 64-bit effect truncation,
byte charge wrap, cooldown wrap, spell-acceptance ordering, callback return bits
and RNG sequences retain the locked behavior. Tests and baseline hashes were
not changed to accept native differences.

Combined **16,742** effects/equipment/inventory/resource/shop/trade contracts pass
with all older hashes unchanged (48.932s). Three production builds are verified
ELF32/i386, GOARCH=386, GO386=sse2 and CGO enabled. Full-suite failure multiset is
unchanged: **1,553 entries; 15 pass / 3 fail / 32 skip packages**. Fresh
`effects-use-port` gameplay passes unchanged repeat-a goldens in **57.062s**, with
overrides disabled, Xvfb and null audio.

Accumulated default/server/highres corpora pass in **121.800s / 105.062s /
104.722s**. The entire connected family is qualified. Local evidence: native-final captures/log,
native-dependencies.log, ports-*.log, binary-checks.json, full-suite-comparison.json,
and build/baseline/runs/effects-use-port/result.json. Obsolete intermediate
captures and old regenerable build-cache entries were removed for disk space;
stable baseline/native captures, logs, gameplay evidence and user assets remain.
