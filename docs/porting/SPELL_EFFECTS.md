# Instant spell effects, summons and charm

## Scope and baseline

The connected batch covers 42 functions / 1,740 removable C lines in GAME4.c
(500CA0–5017F0) and GAME4_2.c (52BEB0–52E0E0). Production code is still C.
The locked baseline has **1,803 cases / 90 complete captures**. Two independent
C runs match byte-for-byte (8.689s / 9.339s); the hashes are committed in
`src/spell_effects_porttest_test.go`. The enforced run plus all 2,246 existing spell-lifecycle cases passes in 18.423s.

The optional porttest fixture invokes the real C functions through a thin
42-operation dispatcher. It reuses guarded player, inventory, spell, duration,
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
  object; the production selection algorithm remains real C. The force callback
  recorder captures object, raw distance bits and opaque argument in call order.
- 500F40's float-typed output parameter carries pointer bits; the dispatcher
  bitcasts it instead of numerically converting the address. The summon packed
  result includes unaligned position and sequence fields.
- Existing spell-lifecycle fixtures have no Effects spec and retain their locked
  hashes. No production correction has been made while constructing this baseline.

## Conversion and qualification still pending

Retain callback ABIs required by duration dispatch and retained spatial queries;
retire obsolete direct-call bridges before qualification. Preserve intermediate
floating-point precision and explicit store boundaries. Compare all complete C
captures unchanged, then run the accumulated port checks in default/server/highres,
three binary builds, the asset-backed full-suite failure comparison and a fresh
headless repeat-a gameplay run. Record the physical C count and commit/push.

Local evidence and scope inventory: `build/port-spell-effects/`, including
`baseline.json`, `c-proof-a-*.json` and `c-proof-b-*.json`. Earlier diagnostic
captures were losslessly gzip-compressed and verified before removing raw copies.
Current production C: **114,659 lines / 149 files / zero reference C**.
