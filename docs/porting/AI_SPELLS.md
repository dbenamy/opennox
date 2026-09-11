# Monster spell decisions and cast actions

Connected scope: GAME4_3 5408A0–541490, fifteen functions covering candidate
selection, inversion scans, summon retries, enchant rejection, healing, cast
animation, morphing and recoil. Existing MonsterCast, spell definitions, duration
lists, map queries and RNG remain the shared production implementations.

## Original-C checkpoint

3,840 generated cases and 172 independent contracts pass against the original C.
The guarded shared AI fixture captures actions, object/update writes, both RNG
indices, callbacks, spatial scan selection, duration/health memory and all 136
permission slots. Existing main-AI and monster-state hashes pass unchanged.

- Corpus: `984b1523ba137d057e07c021202a8dc167c7bd067e27c6257a9f91c17993dc0f`.
- Contracts: `3a1f3cbadc33dd52919d8c064995513c5e4526e0710a8cbb1b70483cef5b51ee`.

Contracts distinguish unsigned cooldowns across frame wrap, permission endpoints,
whole-list rejection for active enchants, summon retries and limits, last eligible
heal target, integer half-health for odd maxima, cast animation/mute ordering and
morph class at the cast callback. Five recoil golden outputs cover exact float
bits and RNG consumption. Accuracy -16777216 distinguishes computing the upper
random bound from the unrounded double inverse accuracy: deriving it from the
rounded float would collapse the range and skip a draw.

Compiled x87 arithmetic keeps inverse accuracy in double for the upper bound,
then spills the lower bound, both intermediate target coordinates and jitter
scale to float32. Both random entry points use Logic. C flags remain unchanged;
Go targets 386/SSE2. Private char/pointer scratch returns have no production users.

Reproduce with build/baseline/env.sh loaded, from src:
`go test -tags porttest -run '^TestAISpell' .`.
OPENNOX_SPELL_CAPTURE writes optional ignored snapshots; OPENNOX_SPELL_CASE narrows
generated cases. Artifacts and assembly: build/port-ai-spells. Production C at this
checkpoint is unchanged: **136,741 physical lines**, 153 files, zero reference C.

## Native conversion

Original-C checkpoint: `326f2a08`. All fifteen C bodies and their prototypes are
removed; `legacy/ai_spells.go` owns selection, scans, healing and cast execution.
Main-AI, combat, navigation and action wrappers call Go directly. No C exports or
reference C were needed. The shared selector preserves candidate order and the
post-callback cooldown reads. Cast arguments remain C-owned across retained
engine callbacks. Healing retains its shared selected-target global and original
integer half-health comparison and last-qualifying-target behavior.

The native focused run matches both locked hashes and all 172 contracts,
including every recoil bit pattern. Production C is **136,242 physical lines**,
a reduction of **499**, across 153 files; zero test-reference C remains.

## Qualification

Accumulated port tests pass in default (70.907s), server (41.240s) and highres
(42.561s) configurations. All three production binaries build and report ELF32,
Intel 80386 and GO386=sse2. The full-suite comparison has exactly the same 1,553
failure entries as baseline: 15 passing, 3 failing and 32 skipped packages, with
no added or removed failures. Fresh headless `ai-spell-port` gameplay exits 0
against both preserved screenshots with overrides disabled and null audio.

Next larger batch: MonsterDef callback loading, strikes, death effects and loot,
GAME5 549040–54A950 (948 physical C lines, stopping before geometry 54A990).
Use separate focused contract groups within one final qualification cycle.
