# Script bindings and object commands — C baseline qualified

Qualified parent: **c425197f**, with **11,409 physical C lines in 49 files**,
zero reference C. The initial candidate covered 22 bodies (472 original body lines) in
`server__script__builtin.c` and `server__script__script.c`: the remaining builtin
table entries, movement/carry commands, string registry and callback transfer.
The current batch selects **19 bodies /373 corrected C body lines**. Carry/drop,
startup inventory cleanup and halberd replacement will form the following
inventory batch because they share client/server inventory fixture requirements.
The current owner/dispatch audit finds all19 reachable from18 external-reference
roots. The production XP/level helper remains C; this batch moves its VM binding.

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

## C contracts and qualification progress

Thirteen new roots pass together in `initial-contracts-c.log`. Nine captures are
frozen from that original-C run; callback, movement, group and string-registry
captures also repeat identically in independent processes. The new contracts
include432 callback read/write cases,1620 AI movement cases,432 group cases,
32 talking/trading cases,280 journal cases,1027 registry writes plus reset/reuse,
480 ownership cases,128 mover cases, and60 experience cases below the next-level
threshold, alongside the prerequisite cases. The experience binding calls the
unchanged C level helper and exercises real player mutations/notifications.

Initial mover setup used an unregistered type ID; actual state synchronization
rejected it. The corrected fixture uses registered types. Initial experience
setup lacked localization data; the real notification path rejected that setup.
The corrected fixture installs the existing string-manager fixture. Neither
failure required a production change. Preserve the failed logs as diagnostics.

All three `c-affected-*` runs pass183 root tests, no skips, and all nine frozen
captures. Before fresh production gates, translation review adds an independent
journal-name boundary contract: embedded NUL terminates the lookup, invalid
string indexes become empty names, and both can address an existing empty-name
entry. This new root does not modify any frozen capture. Its focused run and final
all-target evidence are still pending. No native implementation is installed.

The original-C body uses extent at object offset40 to associate movers. The
fixture makes extent and network code different and verifies all matching movers,
updatable-list order, and repeated calls without duplicate insertion. The early
read-only prose described this field incorrectly; the source audit and contracts
establish the correct extent lookup.

Two private owners can move with this batch: the script-string counter and cached
mover type. Nineteen selected C interfaces and nine now-private bridge exports
can retire after whole-source checking: the four journal remove/update adapters,
roam-byte callback, object resolver, VM string lookup, callback name and callback
index lookup. Keep actual VM dispatch order and remaining inventory C routes.
The draft at `build/port-script-builtins/native-draft.go` is advisory and has not
been installed or tested. Review against frozen C before use.

Final C target sweeps now pass184 roots each with no skips. All107 captured files
are byte-identical across default/server/highres; the nine new hashes remain
unchanged. The manifest now records all107 for native qualification. The final C
runners enforced the nine new hashes and the selected tests' inherited contracts;
the107-file inventory was compared and frozen after those runs. Source
fingerprints are identical across all three final runs. Static checks pass.
Fresh C production binaries/ABI and exact1553 full-suite failures pass; the two
integration scenarios are still running in production session49299.

## Corrected C baseline qualified

Fresh default/highres/server production binaries and ABI checks pass. The full
suite matches exactly1553 known failure entries and15 pass /3 fail /32 skip
packages. Fresh gameplay and save/load scenarios both pass against the preceding
qualified references. Final184-root target sweeps, production and static source
are identical by fingerprint. All jobs, including production49299, are joined.
See script-bindings-c-qualification.json and script-bindings-batch.json. Current
C remains11,404 physical lines /49 files /zero reference C. This commits the
recoverable C oracle before installing the reviewed Go draft.
