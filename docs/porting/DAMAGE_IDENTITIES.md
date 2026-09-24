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
