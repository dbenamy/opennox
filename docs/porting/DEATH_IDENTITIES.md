# Native death callback identities

## Baseline and scope

All 140 selected roots pass in default, server and high-resolution profiles,
without skips, against `170b594a` plus a new stored-key lifetime contract.
Production source is identical to the qualified collision checkpoint; its fresh
production/ABI, known-suite and save/load evidence supplies the production baseline.
The new contract copies all 14 registered keys through foreign ObjectType/Object
storage, checks distinctness and declared sizes after GC/stack growth, and invokes
the stored Glyph key. Existing Glyph contracts already cover replacing its mutable
handler between calls, so no duplicate replacement test was added.

Replace 14 C callback addresses with distinct nonzero-sized static Go slots.
Keep the existing typed death cache, raw fallback, registration order, names,
parser bindings and owner closures. CreateObjectDie and SpawnObjectDie retain
132-byte records; the other twelve retain zero callback-data size. Object layout
and allocation ownership stay unchanged.

Move the direct object-creation, object-death and generator fixture calls to their
existing Go owners. Preserve ImpEgg's uint32 bits, SpawnObject's signed int16 result
and GameBall's C-int narrowing. Replace the fourteen fixture identity entries
one-for-one, retaining their IDs and the independent raw C forwarding observer.
The registered death route continues to discard results as before.

Connected armor/weapon death-message cleanup uses native uint16 pointers and
`alloc.InternCString16`. The previous `internWStr` helper calls that exact owner
and adds only a C pointer cast; preserve interning, lifetime, nullness, evaluation
order and formatting behavior. Existing equipment/language captures cover it.

## Qualification plan

After conversion, rerun the exact focused selection in all three profiles,
safe/static checks, three fresh production/ABI checks, the exact known asset-suite
comparison and headless save/load/resume. Verify source identity, all original
asset hashes, external dependencies and updated glue counts. Keep all assertions
and frozen expectations unchanged. The preceding collision batch established the
shared identity mechanism with the full default-client corpus; this batch reuses
an unchanged death dispatch cache and selects all affected owners and consumers.
Broaden if a failure makes the affected scope uncertain.

Conversion is uninstalled. One Luna helper drafted bounded production changes;
primary owns baseline acceptance, fixture migration, UTF-16 cleanup, integration
and final qualification. Review preserves the current Glyph hook and all direct
fixture return conventions. The helper initially reported an unmigrated fixture
from stale context; inspecting the actual ignored overlay resolved that report.
Primary corrected removal-span coordinates to refer to original source rather
than offsets after earlier draft deletions. Neither correction changed behavior.

[Original evidence](death-identities-baseline.json),
[commands](death-identities-batch.json), [selection](death-identities-tests.txt).
Local artifacts: `build/port-death-identities/`.
