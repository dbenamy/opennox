# Native collision identities

## Original baseline and intended scope

All 356 selected client roots and 355 server roots pass without skips against
production `5bfa54b5`, with two new original-route contracts. The field-of-view
rendering test has a pre-existing `!server` build constraint. They establish 51
distinct non-null keys across 53 named registrations, the three default aliases,
stability through GC/stack growth while stored in real objects, and activation's
historical low-byte-of-current-key return relationship. Existing frozen captures
remain unchanged. The initial 263-client/262-server selection and 91 additional registry
consumers also passed against matching prebuilt binaries before the new contracts.

Replace those 51 C callback identities with unique addresses in a nonzero-sized
static Go table. Keep object layouts, registration order, data sizes, aliases,
parser bindings and mutable handlers. Named equality checks and fixture identity
maps must use the same new keys; fixture IDs and insertion counts stay unchanged.
The installed Go runtime treats linker-allocated globals as pinned, so these keys
remain valid in existing non-Go object storage. Do not use fabricated pointers or
unrooted heap tokens.

## Return and compatibility decisions

Keep one collision registry with optional typed result handlers. Existing raw
C registrations retain their void routes and raw integer fallback. Native keys
must never reach a C function-pointer call. Migrate the magic-missile expiry path
to result-aware dispatch for every configured key, rather than assuming it always
uses DefaultCollide. Preserve original declared integer, pointer-word and signed
short results. Native void owners explicitly return zero when a result is requested;
the old mismatched C invocation had no defined integer result. Production update
dispatch discards that enclosing return, while custom raw C result contracts must
remain exact. This is a reversible, recorded cleanup decision.

The activation algorithm keeps its low-byte-of-current-key relation. Specific
code-address bits were already linker-dependent; no production caller consuming
that activation result was found. Existing activation captures use a separate
aligned C observer, which remains unchanged. The complete frozen owner corpus and
new relationship-based original contracts qualify the key change.

Direct-helper fixtures will invoke native registry result owners without changing
the actor's stored callback. Their pointer arguments must retain the lifetime
previously provided by crossing C; explicit pinning prevents uintptr conversion
from leaving stack pointers vulnerable during nested Go calls/GC.

[Original evidence](collision-identities-baseline.json),
[qualification commands](collision-identities-batch.json),
[conversion selection](collision-identities-tests.txt).
Local artifacts: `build/port-collision-identities/`; caller audit:
`build/port-after-catalog-effect/collision-detail.{md,json}`.


## Qualified conversion

All 357 default/highres and 356 server focused roots pass without skips. The full
default-client port-test corpus completes all 2,431 roots: 2,430 pass and
only the established `TestMapPopulationPrerequisiteProbe` diagnostic skips.
The full server/highres corpus was not repeated; both profiles received the complete
focused owner/consumer selection. Safe/static checks, three fresh production builds
and retained/retired ABI checks pass. Headless character creation and explicit
save/load/resume pass. The ordinary full asset suite matches the known baseline
exactly (304 failure events; 17 passing, two failing, 32 skipped packages).

All accepted gates have identical source fingerprints and all 27 changed/new files
match the reviewed source. Retained export signatures/bodies, existing assertions,
frozen captures and 1,654 original asset hashes are unchanged. Native result
contracts preserve full uint32 values and argument transport. The existing
`TestTemporaryCollisionAndImpact` independently checks the raw C result fallback.

Selected production cgo files fall 224→221 (242/463 eliminated); C exports fall
978→927 (963/1,890 retired). Embedded production C bodies remain 77. Headers remain
157 files, now 3,652 physical lines. Standalone production/test-reference C remains
zero. External native-library bindings are unchanged.

## Review findings and delegation

The first three-profile run exposed four chest fixture setups that deliberately
reuse PentagramCollide as a Death callback. Their new data identity reached the
raw C death dispatcher. The fixtures now explicitly install a test-only typed
death route to the same Pentagram owner, retaining the same identity and effects.
Only setup changed; assertions and frozen captures did not. Production death
registry behavior is unchanged. First-failure evidence remains under `contracts/`;
accepted evidence is under `contracts-fixed/`. Callback audits must follow keys
into other fields, including fixture setup, rather than only textual C names.

Luna drafted the bounded legacy registration/export changes. Primary independently
checked every owner, sparse fixture ID, argument/result width, alias and parser
binding, and owned the result API, pointer-lifetime contracts and fixture migration.
Review corrected a wrong point-pointer type and an extra wrapper before compilation.
The initial audit missed the chest cross-family use; qualification exposed it.
Luna's suggested missing raw-result test was already covered by the independent
magic-missile contract. Bounded drafts remain useful, with primary acceptance;
no measured subscription savings are claimed.

[Qualification](collision-identities-qualification.json),
[inventory](collision-identities-inventory-after.json).
