# Native creation and initialization identities

## Original baseline and scope

All 163 selected roots pass in default, server and high-resolution profiles,
without skips, against `e7ca9d31` plus two new contracts. Production source is
identical to the qualified death checkpoint; its production/ABI, known-suite and
save/load checks supply the production baseline.

The lifetime contract stores keys in foreign ObjectType/Object storage, forces
stack growth and GC, and dispatches stored monster initializer keys afterward.
Creation keys remain on ObjectType; the factory invokes them rather than copying
them onto Object. A second contract covers the public gold Get/Set/Change methods:
initializer-based classification, signed values, 32-bit wraparound and nil objects.
Existing contracts cover real factories, nil/raw callbacks, late replacement,
pending initialization, arguments, and all creation and reward initializer owners.

Replace nineteen C addresses across twenty named registrations: eight create and
twelve init names. MonsterInit and ShopkeeperInit share one identity, retaining
separate declared data sizes. Nil built-ins remain unchanged. Reuse the existing
server typed caches and preserve raw fallback; retain mutable monster hooks inside
closures so later replacement remains visible.

Factor five small creation bodies and Boulder initialization into typed Go owners.
Preserve write order, integer widths and pointer-shaped results in direct fixtures;
registered callback paths discard results as before. PlayerInit's direct fixture
keeps signed-char extension into uint32. ArmorCreate's mixed address/integer return
remains an integer. Move the Gold getter, Chest classifier and quest PlayerInit
caller together with their callbacks. Keep every fixture ID, order and frozen
expectation unchanged, including sparse controls operation 20.

## Qualification plan

Rerun the exact selection in all three profiles, safe/static checks, three fresh
production/ABI checks, the exact known asset-suite comparison and headless
character creation/save/load/resume. Check all retired symbols, retained exports,
source fingerprints, external bindings and 1,654 original asset hashes. Broaden
if a failure leaves affected scope uncertain. The earlier collision milestone
already ran the full default port corpus; this batch changes no shared cache or
fallback implementation.

One Luna helper drafts bounded production changes; primary owns new contracts,
fixture migrations, independent mapping/width review and acceptance. Original
source hashes and removal coordinates are checked before installation. Follow
getter-returned identities into consumers, not just original C symbol references.
The initial audit's claim that Create is copied onto Object was corrected by
reading the factory; it did not describe an implementation change.

[Original evidence](create-init-identities-baseline.json),
[commands](create-init-identities-batch.json), [selection](create-init-identities-tests.txt).
Local artifacts: `build/port-create-init-identities/`.

## Qualified conversion

All 163 focused roots pass in each default/server/high-resolution profile without
skips. Safe/static checks, three production/ABI builds, headless character creation
and explicit save/load/resume pass. The ordinary full asset suite matches its
known results exactly: 304 failure events, 17 passing, two failing and 32 skipped
packages. No assertion or golden was changed. The full port corpus was not
repeated; the preceding collision milestone remains the shared identity check.

All accepted phases have identical source fingerprints. Seventeen changed, new
or deleted source files match review. All retained export signatures/bodies and
1,654 original asset hashes are unchanged. The nineteen exports and prototypes
are gone. Selected production cgo files fall 218→216 (247/463 eliminated); legacy
exports fall 913→894 (996/1,890 retired). Three test cgo imports also disappear.
Embedded production C bodies remain 77; headers remain 157 files, now 3,619
physical lines. Standalone production and test-reference C remain zero. External
SDL2/OpenGL/OpenAL bindings remain unchanged.

Luna supplied the bounded production draft; primary verified the twenty original
registration mappings, nineteen unique keys, native owner extraction semantics,
header removals and fixture operation widths/IDs before installation. Primary
owned the five connected fixture migrations and original contracts. This is
useful delegation evidence, not a measured subscription-cost saving.

[Qualification](create-init-identities-qualification.json),
[inventory](create-init-identities-inventory-after.json).
