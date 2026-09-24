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

One Luna helper drafted bounded production changes;
primary owns baseline acceptance, fixture migration, UTF-16 cleanup, integration
and final qualification. Review preserves the current Glyph hook and all direct
fixture return conventions. The helper initially reported an unmigrated fixture
from stale context; inspecting the actual ignored overlay resolved that report.
Primary corrected removal-span coordinates to refer to original source rather
than offsets after earlier draft deletions. Neither correction changed behavior.

[Original evidence](death-identities-baseline.json),
[commands](death-identities-batch.json), [selection](death-identities-tests.txt).
Local artifacts: `build/port-death-identities/`.


## Qualified conversion

All 140 focused roots pass in default, server and high-resolution profiles without
skips. Safe/static checks, three fresh production builds and retained/retired ABI
checks pass. Headless character creation and explicit save/load/resume pass.
The ordinary full asset suite matches the known baseline exactly: 304 failure
events, with 17 passing, two failing and 32 skipped packages. The previous collision
milestone's full default corpus remains the broad shared-identity check; this
batch did not repeat the full port corpus.

All accepted phases use identical source fingerprints. All nineteen changed,
new or deleted source files match the reviewed draft. Retained export signatures
and bodies, existing assertions, frozen captures and 1,654 original asset hashes
are unchanged. The first profile run exposed five player-death tests calling a registry key
through raw C after it became data. Their getter now returns the registered typed
DeathFunc, and callers invoke it without changing the victim's Death slot. Only
fixture routing and imports changed; frozen assertions are unchanged. The failed
run remains under `contracts/`; accepted evidence is under `contracts-fixed/`.

Selected production cgo files fall 221→218 (245/463 eliminated); legacy exports
fall 927→913 (977/1,890 retired). The death registry and equipment death owner lose
cgo; the PlayerDie export-only file is deleted. One test cgo import is removed.
Embedded production C bodies remain 77, and headers remain 157 files with 3,638
physical lines. Standalone production and test-reference C remain zero. External
native-library bindings are unchanged.

Luna's production draft preserved all fourteen registration closures and parser
bindings exactly. Primary independently reconstructed mappings/removals, corrected
source-coordinate metadata, migrated fixture calls with explicit widths, and
removed the connected UTF-16 type adapter. The helper's stale-overlay warning
was resolved without code changes. Final compilation and qualification validate
the integrated result; no subscription savings are inferred from these outcomes.

[Qualification](death-identities-qualification.json),
[inventory](death-identities-inventory-after.json).
