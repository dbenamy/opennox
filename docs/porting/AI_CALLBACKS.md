# Monster callback registration, strikes, death effects and loot

This connected batch covers GAME5 549040–54A950: 36 C functions and 948 physical
lines, stopping before geometry 54A990. Three loaders serve the live C MonsterDef
parser. Eleven strike and fourteen death/loot callback addresses are installed in
the production blob tables; private target, poison, area and creation helpers
move with their owners. The already-Go BomberDead callback remains shared.

## Original-C baseline (`6923e2ac`)

The shared AI fixture runs the real callback table addresses and retained force,
poison, enchant, area damage, visibility, allocation, modifier and packet engines.
It installs only the shipped table/name region and restores it exactly. Synthetic
C-owned types, modifier descriptors, cloud update data, item initializer/use data
and guarded actor/target/update/definition memory keep tests asset-independent.
The object allocator has 64 slots for debris bursts; absent type names simulate
allocation failure. Created auxiliary buffers are explicitly freed.

There are 4,224 generated cases, 33 callback smoke cases, 172 loot contracts,
100 strike contracts, 36 poison contracts, 48 debris contracts, two precision
contracts, 20 cloud lifetime contracts, and 88 callback-loader checks. Existing
spell, main-AI and monster-state hashes pass unchanged. Contracts cover:

- Exact case-sensitive callback names, case-insensitive NULL, empty/unknown
  tables, return values, callback slot identity, field-only mutation and guards.
- Targets, angular/range rejection, opaque/transparent walls, strict nearest
  ties, scan order, damage parameters and force rereads after damage callbacks.
- Poison chance boundaries, immunity, RNG indices, player status timestamps,
  and wasp's poison-before-force ordering versus other poisonous strikes.
- Loot probability endpoints, game-mode gates, exact modifiers and ammo five;
  debris allocation failures, creation order and rotating indices.
- Signed FPS conversion and float32 spill before cloud lifetime truncation.

Hashes are locked in ai_callbacks_porttest_test.go:

- smoke: `9c728bf2e7aca4122a4fe999d773a25a95f11a5a8d56dd52ef622865d8e819df`.
- corpus: `dfb3eabe2a2740dce4e5b72ce5f58b646450aee06716bd1702a8703920263376`.
- loot: `656c517c1c229997a5451ed1e965f0552c538935c8c90e495b37ca141fa9af70`.
- strike: `643d99fbd75072d3379992f93d7252c457440574afc170c1f532429d2cc2dffb`.
- poison: `0935568717e4b8b976cb6ea4ffc2b1acf30e1fb5619f5da0b04d75b6678c2e2e`.
- debris: `56b8186e6e254e6ce2cd46ce1fb6af07cf271d9aa40d63c85562b4ba9637d196`.
- precision: `313607d3b331bd5ba6e4e8c551c65edc8bdf2aab37871762e14a08445d6f933d`.
- cloud: `a8ab0325f416409b1ffcc70e7cd9537fdf83a1e00f7a6c6b233c8dbb5785f987`.

Loader hash: `4c630d08dca8d187a0634cfdb80e6f25b3fae8c022aec015fe61ebe017380c67`.

## Compiled arithmetic and undefined padding

sub5494C0 numerically converts the float interpretations of object class/flags
to bytes; it does not reinterpret them as integer bitfields. Compiled x87 uses
a truncating signed-word conversion and tests its low byte. Independent cases
make replacing that with ordinary class/flag checks fail. Target selection keeps
X delta in double, spills Y delta to float32, spills distance and the radius-adjusted
gap to float32, then uses strict comparisons. Two large-coordinate direct-callback
contracts distinguish premature X rounding and missing Y rounding without map
indexing. Ogre keeps both deltas full precision through its magnitude, spills
Y for the dot product, and compares the unspilled radius-adjusted gap. Assembly
artifacts are ignored under build/port-ai-callbacks.

sub54A390 initializes four modifier descriptors in a 20-byte local and the shared
setter copies all 20 bytes. The fifth word (ModifierInitData.Field16) is undefined
stack content. No semantic reader was found; generic respawn code copies it.
The baseline checks the four defined slots and ammo. Native code initializes
the reserved word to zero; it cannot preserve undefined stack contents.

## Native conversion

All 36 function bodies are removed from GAME5.c. The 28 public parser/table
entry points use generated C ABI bridges into ai_callbacks.go; eight private
helpers are Go-only and their obsolete declarations are removed. Strike ABI
float arguments transport raw object-pointer bits, as in the original C.
The original-C hashes all pass unchanged (2.405s). No reference C is retained.

Production C is **135,294 physical lines**, down **948**, across 153 files;
zero test-reference C lines. Run from src with build/baseline/env.sh loaded:
`go test -tags porttest -run '^TestAICallback' .`.
OPENNOX_CALLBACK_CAPTURE optionally saves ignored snapshots; CASE and CORPUS_CASE
suffixes narrow the respective case lists.

Accumulated port tests pass on default (55.763s), server (43.408s) and highres
(44.000s). All three production builds pass and are ELF32/i386 with GO386=sse2.
Full-suite metadata matches exactly the 1,553 known failure entries: 15 passed,
3 failed and 32 skipped packages, with no added or removed failures. Artifacts
are under build/port-ai-callbacks. Fresh `ai-callback-port` headless gameplay exits 0 against preserved screenshots,
overrides off and null audio.
