# Native damage and damage-sound identities

## Original baseline and scope

All 76 selected roots pass in default, server and high-resolution profiles without
skips against `ed82a5a8` plus two new contracts. Production source is identical
to that qualified checkpoint; its production/ABI, known-suite and save/load
evidence supplies the production baseline.

Retire all 26 damage_exports.go wrappers plus two mutable sound export bridges.
Thirteen distinct stable identities replace eleven damage and two sound addresses;
Stone and Default retain distinct keys even though they share damageDefault.
Reuse both boolean/value damage caches and the sound map; raw fallback is unchanged.
Defaults in object_type and the nil-slot fallback in damageDefault move with them.
The nil-slot fallback must call the fixed default sound hook, not the mutable
server.DefaultDamageSound identity. Existing sound-owner contracts distinguish them.

The 76-root baseline covers all Damage*, XferSoundRegistry*, Projectile* and
ResourceDefinitions* tests. New contracts verify foreign ObjectType/Object key
storage across GC/stack growth, stored Ball and sound dispatch, all thirteen keys,
and damage-name first-NUL behavior. Existing roots cover full signed results,
raw argument/result forwarding, adapter evaluation order, real audio owners,
late-bound sound overrides, rejection with nil objects and direct owner effects.

Primary migrates the 26-operation fixture one-for-one. Preserve uint32 return
bits, uint16 item-destroyed arguments, float32 payload interpretation and float64
conductivity bits. First-NUL semantics use the exact alloc.GoString owner and
existing string interning. Preserve the original evaluation order of fixture
argument lookups. BlackPowder's rejected-kind early return also exists at the
start of the native owner. Keep independent C pre/defend/scalar observers.

Only eleven of the original 26 fixture identities remain meaningful. Reserve
fifteen removed slots and skip their nil lookup results rather than inserting nil
into the ID map. The eleven keys retain original sparse operation IDs. The map's
cardinality contributes to subsequent dynamically allocated capture IDs.

The source audit found no other in-tree consumers for the fifteen adjacent
wrappers; translating the shared fixture once permits deletion of the full export
file. Internal export retirement is within the agreed engine-glue milestone.
No callback ownership/layout/API modernization is included. The helper's claim
that sound.go must retain cgo was corrected against its actual contents: after
its two exports disappear it has only native Go hooks/functions.

After conversion, repeat all focused profiles, safe/static checks, three production
and ABI builds, exact known-suite results, headless character creation/save/load,
source fingerprints, inventory and original-asset hashes. Existing assertions and
frozen captures remain unchanged. Broaden if affected scope becomes uncertain.


One Luna helper drafts the bounded production changes; primary owns contracts,
fixture migration, independent review and acceptance. Local artifacts are under
`build/port-damage-identities/`.

[Original evidence](damage-identities-baseline.json),
[commands](damage-identities-batch.json), [selection](damage-identities-tests.txt).


## Qualified conversion

All 76 focused roots pass in each default/server/high-resolution profile without
skips. Safe/static checks, three fresh production builds and ABI checks, headless
character creation and explicit save/load/resume pass. The ordinary full asset
suite matches the known baseline exactly: 304 failure events, 17 passing, two
failing and 32 skipped packages. Existing assertions and goldens are unchanged.
The full port corpus was not repeated; the earlier collision milestone remains
the broad shared-identity check, and no shared registry/cache implementation changes.

All accepted phases have identical source fingerprints. Eleven changed, new or
deleted source files match review; retained export signatures/bodies and all
1,654 original asset hashes are unchanged. All 28 retired exports are absent.
Production cgo files fall 216→213 (250/463 eliminated); selected legacy exports
fall 894→866 (1,024/1,890 retired). No test cgo import is removed: the fixtures
still use independent C observers. Embedded production C bodies remain 77.
Headers remain 157 files with 3,592 physical lines; 27 prototypes are removed
because PlayerDamageSound had only inline declarations. Standalone production
and test-reference C remain zero. External native-library bindings are unchanged.

Luna supplied the bounded production draft and independently reviewed the two
primary fixture migrations. Primary verified every registration mapping, sparse
fixture key and prototype removal, retained-export identity and argument/result
widths. The sound.go cgo correction was made before installation. Qualification
found no behavior differences and required no assertion or capture changes.

[Qualification](damage-identities-qualification.json),
[inventory](damage-identities-inventory-after.json).
