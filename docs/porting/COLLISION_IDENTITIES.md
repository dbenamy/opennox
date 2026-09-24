# Native collision identities

## Original baseline and intended scope

All 356 selected client roots and 355 server roots pass without skips against
production `5bfa54b5`, with two new original-route contracts. The field-of-view
rendering test has a pre-existing `!server` build constraint. They establish 51
distinct non-null keys across 53 named registrations, the three default aliases,
stability through GC/stack growth while stored in real objects, and activation's
historical low-byte-of-current-key return relationship. Existing frozen captures
remain unchanged. The initial 263-root selection and 91 additional registry
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

Conversion and qualification are pending. One Luna helper owns the bounded legacy
production draft; primary owns registry/result APIs, contracts, fixture routing,
independent review and qualification.

[Original evidence](collision-identities-baseline.json),
[qualification commands](collision-identities-batch.json),
[conversion selection](collision-identities-tests.txt).
Local artifacts: `build/port-collision-identities/`; caller audit:
`build/port-after-catalog-effect/collision-detail.{md,json}`.
