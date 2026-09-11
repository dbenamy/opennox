# Monster command, animation and state layer

This connected batch covers monster commands and their stack construction,
animation selection, playerlike animation ranges, movement audio, facing,
capability/aggression/state predicates, mimic transitions and shared type caches
around 533790–534A40. Shield-threat selection, action-name parsing, weapon/strike,
banishment, observation, packet encoding and direction conversion remain shared
engines. Keep C entry points where production C callers remain, and route native
callers directly. No separate qualification per predicate.

## Original-C checkpoint

The shared AI fixture passes 14,336 generated cases and 212 independent contracts,
with identical complete-state hashes on repeated runs:

- Generated: `cce74b363d6c5c4d2210781f9b252f503ee34763880ddf8e14037d7551bab1f5`.
- Contracts: `44bab745a7d8d2e58461eafc5821f61a083f1367898d3dd03e0a5d006a978c03`.

Coverage includes signed empty animation stacks, complete flags and cached type
IDs, nullable health, NaN melee range, aggression/speed boundaries, wrapping frame
arithmetic, running-mutator return values, mimic age/distance and full stacks,
NPC animation pointers and global byte writes, frame-zero and low-byte fallback
selection, command broadcasting and ownership/death gates, source/player/null
variants, second controlled-unit filtering, guard signed direction, escort owner
arguments, player broadcast state, banishment scripts/effects, observation calls,
deferred audio, facing and dot-product precision.

The fixture reuses real C-owned object/state buffers, synthetic types, map/wall
state, packet queues, health, action stacks and deferred audio from previous
batches. Player animation ranges and weapon-animation mappings are deterministic
lookup inputs restored after testing. Observation records its call boundary;
banishment and packet creation execute their production implementations. Result
snapshots include normalized additional objects, update-data mutations, caches,
RNG, animation storage, player command state and packets. Returned sound-helper
scratch values are unused by all callers; tests record its audio effects.

Independent checks found a fixture setup error: alloc.New allocates zeroed
storage but does not copy its argument value. Facing-vector inputs are now
assigned explicitly before the C call. The earlier facing-dot regression already
protects the first-product float32 spill; compiled x87 instructions confirm it.
That old C helper can be retired without keeping a test-only C algorithm.

At this baseline checkpoint production is unchanged: **137,882 physical C lines**,
153 files, zero reference C. Ignored artifacts: `build/port-ai-state`.

## Native conversion

The 29 C bodies are replaced by `legacy/ai_monster_state.go` and the existing
native facing-dot helper. Production C falls by **585 physical lines** to
**137,297**, across 153 files with zero reference C. Live C entry points remain
generated Go bridges; private helpers and their obsolete declarations are
retired. Movement actions now contain no cgo, and combat, navigation, roaming,
commands and lifecycle use the shared native helpers directly.

Both original-C hashes match after conversion. Review and independent contracts
caught the joint Zombie-cache miss rule, Plant's early return before Zombie
lookup, animation flag OR arithmetic, zero/low-byte player animation lookup,
command messages and their ordering, guard sight range and signed direction.
The signed direction must be converted to Go `int` for action argument storage.
Running setters return the post-mutation flags. NaN melee behavior and the
facing-dot float32 spill remain explicit. No test-only C algorithm is retained.

## Qualification

Accumulated port tests pass for default (72.373s), server (39.660s) and highres
(38.845s). All three production binaries build and report ELF32, Intel 80386 and
GO386=sse2. C x87 flags are unchanged. The full asset-backed suite matches the
established failure multiset exactly: 1,553 entries, 15 passing / 3 failing /
32 skipped packages, no added failures. Fresh `ai-monster-state-port` headless
gameplay exits 0 against both preserved screenshots, overrides disabled.

Original-C checkpoint: `6aadf8ef`. Reproduce focused checks with
`go test -tags porttest -run '^TestAIMonsterState' .` from src after loading the
baseline environment. Ignored qualification artifacts are in build/port-ai-state.
Next connected batch: the main monster AI owner, idle vocalization, dangerous-unit
gate, attack-action unwinding, shield-threat scan and dodge reaction (GAME5
5469B0–547C50 plus GAME4_3 533E70/533EB0). Retain casting, damage and map engines.
