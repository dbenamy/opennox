# Sustained spells and teleport callbacks

## Scope and locked C baseline

This connected batch follows native instant spell effects (`87fd50dd`). It covers
53 functions / 2,635 removable C lines in GAME4_2.c (52E210–52F820) and GAME4_3.c
(52F8A0–531AF0): mana drain, energy bolt, firewalk, force of nature, greater heal,
channel life, shields, chain lightning, creature tagging, teleport callbacks,
mana bomb, turn undead, plasma and moonglow. The production algorithms now live in native Go.

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
named factories. Two independent expanded C runs match all 129 complete captures byte-for-byte
(11.479s / 11.427s). SHA-256 values are locked
in `src/sustained_spells_porttest_test.go`. The enforced run plus all 1,803 instant
spell-effects and 2,246 lifecycle cases passes in **29.602s**.

Initial baseline `4fad7f26` was pushed before conversion. A final source audit
added the energy-bolt selected-target global (2487880) to save/reset/restore and
both snapshots, with an independent positive selection assertion. This expanded
C baseline was repeated, locked and pushed as `59ef3bc5` before replacing any
algorithms. All seven named globals used by this scope are captured.

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

## Native conversion and qualification

C baselines 4fad7f26 and expanded 59ef3bc5 were pushed before conversion.
All 53 routines / 2,635 C algorithm-section lines are native Go; 11 obsolete
exports are retired and 42 duration ABIs remain. All **2,393 cases / 129 complete
captures** match C byte-for-byte. New plus 1,803 instant and 2,246 lifecycle
cases pass together in 29.641s. No baseline hashes changed.

Production C: **110,283 lines / 149 files / zero reference C** (−2,636 physical
lines, including one trailing blank line). Accumulated **51,001 focused cases /
701 capture groups** pass in default / server / highres:
221.927s / 297.531s / 232.991s.
Three production binaries verified ELF32/i386/SSE2/CGO; 11 retired symbols absent.
Asset-backed full-suite failure multiset unchanged: 1,553 entries, 15 pass,
3 fail, 32 skip. Fresh unchanged repeat-a gameplay passes in 46.968s.

See [SUSTAINED_SPELLS.md](SUSTAINED_SPELLS.md) and
build/port-sustained-spells/qualification.json; gameplay evidence is in
build/baseline/runs/sustained-spells-port. Native implementations preserve the
plasma direction-expression quirk, complete lightning topology and allocation
links, float store boundaries and grouped type-cache initialization. Any gameplay
cleanup is separate from this conversion and noted for later review.

Implementation lives in `legacy/sustained_spells.go`, `sustained_resources.go`,
`sustained_projectiles.go`, `sustained_teleport.go`, `sustained_lightning.go` and
`sustained_exports.go`. Internal fixture calls invoke the native helpers directly;
retained dispatchers only adapt callback signatures, including float pointer bits.
The native damage owner invokes shield absorption directly. No retired source
references remain, and no C algorithms are retained solely as test references.

Review later: plasma's original bitwise OR with 0xC makes its direction predicate
always true. Enemy and interaction checks still apply. The Go port preserves this;
any correction needs a separate gameplay decision. First native qualification
matched every full capture, and final source review additionally preserved paired
and grouped cache initialization even when secondary cache words are already set.
The final repeated native captures match the same locked C oracle unchanged.
