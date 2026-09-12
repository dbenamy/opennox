# Instant spell effects, summons and charm

## Scope and baseline

The connected batch covers 42 functions / 1,740 removable C lines in GAME4.c
(500CA0–5017F0) and GAME4_2.c (52BEB0–52E0E0). Baseline `f51f09b6` was
committed and pushed before conversion; all 42 functions are now native Go.
The locked baseline has **1,803 cases / 90 complete captures**. Two independent
C runs match byte-for-byte (8.689s / 9.339s); the hashes are committed in
`src/spell_effects_porttest_test.go`. The enforced run plus all 2,246 existing spell-lifecycle cases passes in 18.423s.

The optional porttest fixture originally invoked the real C functions through a thin
42-operation dispatcher. After conversion it invokes native functions directly,
except the six callbacks whose C duration ABIs remain required. It reuses guarded player, inventory, spell, duration,
object-factory, map and callback owners. The summon capacity owner is explicitly
bound to the real root implementation, replacing the older AI fixture stub. Complete captures include return bits,
object and update records, health/mana/buffs, ownership, packets, audio, creation,
deletion, cache state and surrounding allocation guards. Reference algorithms
are not copied into test-only C.

Coverage includes all entry points, resource and enchantment boundaries, signed summon-cost table words,
projectile directions and levels, summon start/finish/cancel and sequence wrap,
player and monster casters, guide sizes, actual capacity rejection, terrain
admission, unavailable named factories, charm acquisition,
portal slots and replacement, linked doors, radial forces, quake damage,
inversion, healing, fumble and nearest owned glyph selection.

Independent contracts require healing and confusion effects, actual spawned
objects, successful summon completion, charm ownership transfer, locking both
linked doors, positive radial force, force callback arguments, missile ownership
transfer and selection of the nearest eligible glyph.

## Fixture details

- Guide names normally receive their pointer initialization during game startup.
  Headless fixtures explicitly install the real Bat guide name and valid sizes
  1, 2 and 4, then restore both guide tables. Terrain tests reuse the existing
  guarded production-layout tile grid, including the real teleport predicate.
- Summon completion invokes the existing root Go owner, including actual object
  creation, ownership, orders and notifications. The fixture supplies complete
  monster update storage and health, adopts the returned object for capture and
  cleanup, and restores the pending-object list. Normal positive cases explicitly
  select follow order; zero means banish.
- The factory fixture preserves case-exact type names, including Glyph. Fireball
  levels use the actual Fireball, StrongFireball and TitanFireball type names;
  their startup pointer table is initialized and restored explicitly. Every
  direction/level case independently requires a created fireball.
- Spatial inputs carry the engine active flag before registration in the real
  index. Real shape, distance and ray predicates remain in use.
- Glyph detonation is outside this batch. Its callback records the selected
  object; the production selection algorithm is exercised directly. The force callback
  recorder captures object, raw distance bits and opaque argument in call order.
- 500F40's original float-typed output parameter carries pointer bits; the C
  baseline dispatcher bitcasts it instead of numerically converting the address. The summon packed
  result includes unaligned position and sequence fields.
- Existing spell-lifecycle fixtures have no Effects spec and retain their locked
  hashes. No production correction has been made while constructing this baseline.

## Native conversion and qualification

Five `legacy/spell_effects*.go` implementation files own resource/buff effects,
spatial forces and doors, projectile creation, portals/glyph selection and
summon/charm lifecycle. Six duration callback ABIs remain; 36 obsolete exports
are retired. Internal spatial callbacks now use native closures through the same
existing Go map owners. No C algorithms remain solely for testing.

All **90 complete captures match C byte-for-byte**, with no hash changes.
The 1,803 new cases plus 2,246 existing lifecycle cases pass in **18.220s**.
The first native quake comparison caught a center-distance substitution: the
original uses shape-aware distance with a minimum-distance clamp. Reusing the
existing `stateDistance` owner restores exact output. Charm range comparisons
retain double precision and conditional monster health reads retain C ordering.
No production behavior correction was needed for this batch.

Accumulated checks, including **48,608 focused cases / 572 capture groups**, pass
in default / server / highres: **210.312s / 302.784s / 222.884s**. Three production builds are verified
ELF32/i386, SSE2 and CGO enabled; all 36 retired symbols are absent from each.
The asset-backed full-suite failure multiset remains exactly unchanged: 1,553
entries; 15 packages pass, 3 fail, 32 skip. Fresh unchanged repeat-a headless
gameplay passes in **40.050s**, with Xvfb and null audio.

Production C: **112,919 lines / 149 files / zero reference C** (−1,740 lines).
Local evidence: `build/port-spell-effects/qualification.json`, `baseline.json`,
`variants.json`, `binary-verification.json`, `full-suite-comparison.json`, and
complete `c-proof-a`, `c-proof-b`, `native-final` captures (losslessly compressed
and verified). Gameplay: `build/baseline/runs/spell-effects-port`.
