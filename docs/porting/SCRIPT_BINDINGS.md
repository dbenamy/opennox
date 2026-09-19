# Script bindings and object commands — C baseline in progress

Qualified parent: **c425197f**, with **11,409 physical C lines in 49 files**,
zero reference C. The candidate covers 22 bodies (472 original body lines) in
`server__script__builtin.c` and `server__script__script.c`: the remaining builtin
table entries, movement/carry commands, string registry and callback transfer.
Whole-source textual reachability finds all 22 connected; owner and dispatch
review and baseline contracts are still in progress.

## Predicate and signature prerequisite

Two boolean predicates (builtins 184/185) accepted an argument that their actual
`CallIntVoid` caller never supplied. Missing objects caused them to push that
indeterminate argument. The actual legacy table and normal VM dispatch both
reproduce a nonzero result for missing IDs 0, -1, -2 and 2147483647: 16 cases.
Diagnostics record the incorrect boolean result without dumping incidental data.

The correction returns false for a missing object and declares both callbacks
with `(void)` parameters. Valid objects keep their original subclass bits 8 and
7 respectively. The SetRoamFlag callback receives the same signature correction;
its unused initializer read another unpassed argument, but its live path always
replaced that value with the popped low byte.

Independent contracts use the actual VM stack, object resolver/cache and builtin
dispatch. The original C passes 8,192 valid-object combinations, covering both
dispatch paths, the low ten subclass bits and high-bit noise. The corrected
missing-object contract also checks stack balance. Another 6,144 roam cases cover
byte truncation, monster/nonmonster classes, destroyed/disabled flags, preserved
upper bytes and both dispatch paths.

The first roam fixture incorrectly treated a player record as monster data.
Its typed accessor rejected that setup. The corrected fixture owns a separate
monster update record and restores the original object before freeing it.
Those initial failed runs are retained and do not qualify the prerequisite.

All dispatch priorities stay unchanged: registered handlers, existing map
compatibility overrides, legacy table, then root fallback. The compatibility
implementation is outside this batch and remains untouched.

Evidence is under `build/port-script-builtins`: `missing-object-original.log`,
`object-predicates-original.log`, and `prerequisite-final-{default,server,highres}`.
The tracked prerequisite manifest and test selection reproduce the corrected
contracts. All three targets pass all three selected roots with no skips and identical
source fingerprints. This is a recovery prerequisite,
not a completed conversion or a fully qualified production baseline.

Working C is **11,404 physical lines / 49 files / zero reference C**, five fewer
than the qualified parent. Freeze the remaining original-C contracts and run
fresh production qualification before converting the selected functions. Old
ignored selection offsets predate this correction and must be regenerated.

## Remaining baseline work

Cover journal one/all-player dispatch, actual group traversal, experience and
owner/pet commands, talking/trading state, startup inventory changes, shipped
halberd selection, string-registry capacity/lifetime, actual AI action stacks,
mover state and carry/drop selection. Audit host-player preconditions rather than
silently inventing new behavior for unsupported states.

Callback transfer needs actual cryptfile read/write fixtures and VM callback
lookup: signed version checks, empty names, embedded terminators, 1023/1024-byte
read boundaries, editor/runtime modes, exact consumed bytes and untouched fields
on early rejection. Reuse existing object-transfer and compiled-script fixtures.
The broader batch is not yet frozen, converted or qualified.
