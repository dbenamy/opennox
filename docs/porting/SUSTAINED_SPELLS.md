# Sustained spells and teleport callbacks

## Scope and locked C baseline

This connected batch follows native instant spell effects (`87fd50dd`). It covers
53 functions / 2,635 removable C lines in GAME4_2.c (52E210–52F820) and GAME4_3.c
(52F8A0–531AF0): mana drain, energy bolt, firewalk, force of nature, greater heal,
channel life, shields, chain lightning, creature tagging, teleport callbacks,
mana bomb, turn undead, plasma and moonglow. Production algorithms remain C.

The optional `Effects.Sustained` fixture reuses guarded duration, object, player,
resource, map and factory owners. Its thin C dispatcher explicitly bitcasts the
float-typed duration pointers. Complete captures include return values, record
and output bytes, callback ordering, resource/object state, packets, creation,
initialization guards, cache words, named lightning globals and live child
segments from the real Go duration allocator. Damage callbacks for otherwise
unconfigured test objects use the existing callback recorder; target-selection
and damage calculations remain the production algorithms under test.

The locked C baseline has **2,393 cases / 129 complete capture groups**.
Coverage includes every entry point, duration/dead/nil boundaries,
positive spawn and enchantment checks, mana transfer and shield arithmetic,
spatial distance/ray boundaries, glyph slots and teleport movement, wand state,
firewalk placement, multistep lightning, healing/channel life and unavailable
named factories. Two independent C runs match all 129 complete captures byte-for-byte
(11.977s / 11.281s). SHA-256 values are locked in
`src/sustained_spells_porttest_test.go`. The enforced run plus all 1,803 instant
spell-effects and 2,246 lifecycle cases passes in **28.661s**. No native
conversion has been applied. Commit and push this baseline before conversion.

## Fixture and source decisions

- Move creature-tag's caster-data read below the existing nil guard before
  baseline capture. The current optimized C build already passes the nil-caster
  test; this makes source behavior defined, rather than claiming a reproduced
  runtime crash. See [DECISIONS.md](DECISIONS.md).
- Normal transition targets carry sufficient health to avoid invoking unrelated
  death/lifecycle owners. Dedicated health and shield cases use the existing
  guarded resource fixture and its death recorder.
- Save/reset/restore named lightning globals as well as blob-mapped cache words.
  Capture actual linked child segments and normalize their real allocation IDs.
  Repeated full captures exposed stale address labels from earlier cases; always
  derive live segment labels from DurSpell.ID, retaining links and all 120 bytes.
- Supply per-level healing input and the actual lightning-count startup table
  (1, 2, 3, 4, 5 from blob_587000.dat). The inherited fixture leaves the latter
  uninitialized; a stronger contract exposed every level creating only one
  segment. Require the expected count with one through five spatial targets.
  All added targets have valid explicit collision shapes. Preserve direction
  tables. Created effect initialization storage includes
  destination pointers/positions and explicit trailing guard checks.
- Clone the AI target's player update and player metadata into separate guarded
  regions for this optional fixture. The older AI fixture aliases that storage to
  player zero, which made an apparent mana drain write the caster's own data.
  An independent two-tick transfer contract exposed the alias. Restore the
  original pointers on cleanup; unrelated old fixtures remain unchanged.

## Conversion gate

The ABI audit finds 42 duration callback getters that require C entry points.
Eleven helpers can retire their ABIs, including the shield-damage bridge whose
only remaining production caller is native Go. Preserve pointer getter callers
when deciding which exports are required. Keep the three embedded declarations
for lightning segment creation, waypoint selection and buff cleanup.

After conversion, compare every complete capture unchanged, run accumulated port
checks in default/server/highres, build/verify three production binaries, compare
the asset-backed full-suite failure multiset, and run fresh unchanged headless
repeat-a gameplay. Then record the physical C count, commit and push.

Current production C: **112,919 lines / 149 files / zero reference C**. Local
scope, ABI audit, diagnostic captures and logs: `build/port-sustained-spells/`.
