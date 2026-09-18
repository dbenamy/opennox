# Monster control and definitions

## Scope and checkpoint

Parent **0bbaba1d** is committed and pushed. The next connected batch covers
monster action/control state, script unit selection, pending ownership and monster
attribute loading: **40 live C functions / 1,001 body lines**. Three proven orphan
functions add 31 lines. See [monster-control-scope.json](monster-control-scope.json).
Production C is unchanged: **40,218 physical lines / 74 files / zero reference C**.

## Reachability and test design

Comment-aware caller review found that nox_xxx_monsterActionIsCondition_50A010 and
nox_xxx_monsterIsActionScheduled_50A090 are mentioned only in comments on existing
Go methods; those mentions are not live callers. sub_50CAC0 also has no caller.
The complete source/header/Go-preamble/assembly search confirms no registration or callback references; evidence is in build/port-monster-control/orphan-audit.json. Do not
translate them just to create tests, or change the existing Go APIs to reproduce
unused C behavior.

Use actual monster/player records, registered AI actions, map index, scripts,
RNG, ownership and allocation pools. Cover stack boundaries, target refresh,
animation frames, death/decay and observer effects, script selection and control
commands. Pending ownership needs add/resolve/reset/free contracts. Definition
loading needs every shipped property kind, whitespace/comments, ARENA/SOLO filters, EOF,
partial/invalid records and the actual shipped encrypted monster.bin.

Audit shipped numeric data before captures: action metadata at 0x587000:230388
(16-byte records), field schema at 248192 (12-byte records), attack/flag name
pointers at 247464/247536, relocated strings and direction vectors. Exercise
nonzero pointer argument metadata and real parsed fields. Reuse qualified spatial
production for the C baseline only if all production sources are identical;
always run fresh production after conversion.

Local proposals and audits: build/port-monster-control. The C baseline is now
frozen and qualified across all three targets. Prior installers and cleanup apply modes
are consumed; do not replay them. About 20 GiB of disk is available after verified
asset-copy cleanup; original assets and archive remain intact.


Initial fixtures invoke real C commands and registered action-stack operations.
They cover class/flag admission, idle/wait/dead and full stacks, control-byte
boundaries, wrapped wait/flee deadlines, movement/combat commands, copied-frame
state, returned owned pointers and the empty-head legacy word. The first refresh/animation process passed 4,608 and 12,000 cases respectively.
Cache sequence contracts also passed. Death/revival/chapter, pending ownership and
definition fixtures also pass; all captures are now frozen.
No production code has changed.

## Boundaries and review notes

- Animation delay/frame counters wrap as bytes. At delay 255, the byte counter
  never reaches the integer threshold 256. The original function's return may be
  the low byte of an animation address; its sole production caller ignores it.
  Capture the complete four-byte animation state and immutable table instead of
  freezing an address fragment.
- The empty action-head API reads the owned word immediately before the stack;
  it differs from the existing Go nil-head convenience API. Keep that distinction.
- Script cache free-list predecessor/tail words can remain stale after a pop;
  active-list links and free-list forward reachability are the operational contract.
- Actual shipped definition kinds are 0, 1 and 3 through 8. The sound-kind branch
  (2) has no shipped schema entry. Fixtures install original bytes and all 69
  name relocations, plus real registered callback addresses.
- Original C leaks its 248-byte incomplete definition on a rejected callback or
  damage type. The Go conversion will free rejected records while preserving
  return/list behavior. This is a reversible ownership correction for later review;
  original C remains unchanged for capture and production-identity reuse.
- Asset copies and tests use temporary directories; the original encrypted
  monster.bin, archive and asset extraction remain unchanged.

## Qualified C baseline

Fifteen focused roots produce **20,997 records**, repeated byte-for-byte in
separate processes and frozen. The completed corpus includes 144 death-state
cases, 324 revival cases, 18 quest/ability reward cases, four real inventory-drop
cases, 80 chapter reports, all 72 shipped action layouts, cache capacity/links,
pending-owner pool exhaustion/reuse and 88 synthetic definition cases. Three
shipped-file mode captures include all loaded records and successful type lookup.

All affected target sweeps pass with zero skips; exact counts, durations and
capture totals are in [monster-control-c-qualification.json](monster-control-c-qualification.json).
All four gates share the same 2,264-file source manifest. Static mapped-state
checks pass. Production reuse verifies 2,254 original sources unchanged, ten new
porttest-only files and all three parent binary hashes. See
[monster-control-c-production-reuse.json](monster-control-c-production-reuse.json).

Fixture corrections were confined to tests: use real reporting hooks for owner
cleanup; use the actual uninterruptible action ID (61); inventory items are
inactive before the default drop stages them in the world. Frozen expectations
were established only after these independent contracts passed. No production
source or prior golden changed.

Next: replace the forty live functions, remove the three proven orphans, retire
obsolete C interfaces and compare every broad capture. Fresh production,
headless gameplay, save/load and flat-map regeneration remain required.
