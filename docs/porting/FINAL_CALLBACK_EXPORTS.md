# Final engine callback exports

Status: original baseline accepted; conversion is an uninstalled draft.
Qualified production before this batch: `519ce712`.

The last ten selected legacy C exports provide three player-file section
callbacks, map object-data loading, screen-particle drawing, three tooltips,
main-menu events and dynamically assigned FlameCleanse updates. Their algorithms
already live in Go; this batch removes the intervening C routes.

The draft keeps the player-file table and particle callback fields mutable. Stable
native identities dispatch to their existing Go owners; unknown foreign callbacks
retain the existing raw fallback. Tooltip and object-update identities use existing
registries. Menu events and map loading have a single production consumer each,
so they use typed Go functions without introducing unused identities.

Review decisions:

- Evaluate menu event code and arguments before reading the mutable hook; return
  a nil response for zero, preserving the original adapter's interface behavior.
- Keep player-file section framing, table lookup order, nil fallback arguments and
  signed return values. Fixture numeric IDs and captures stay unchanged.
- Cache each particle's next link before calling its draw callback; callbacks can
  delete, recycle or insert nodes during traversal.
- Options tooltips read only the first byte of draw data, as before. The conversation
  callback is an installed no-op. FlameCleanse is dynamically registered, without
  inventing a named object-type update.
- Preserve map callback nil validation before unknown-section handling, and treat
  only a zero callback result as failure.

Eight independent new contracts passed a default-profile preflight on the original
routes. They cover menu pointer/argument/return bits and hook ordering; installed
options tooltips for all 256 low bytes with nonzero upper bits; conversation
no-op behavior; particle traversal mutation; foreign section result boundaries;
map validation/dispatch order; and FlameCleanse via actual object-update dispatch.

The focused selection covers 294 root names (292 compiled in server), including
owner families, menu resources, character creation, map sections and orchestration.
Original captures run twice in separate processes in all three profiles; the first
default run is the complete corpus (2,487 compiled roots, including the established
map prerequisite probe skip). Frozen captures are not regenerated. Final acceptance
also requires converted profile runs, static/safe checks, three production builds
and ABI checks, a fresh GUI preview and final save/load scenario, exact known-suite
comparison and original-asset integrity.

Primary owns this batch. Luna remains quota-blocked; no substitute helper model is
used. Local draft, source identities, detailed selector and baseline outputs are
under `build/port-final-callback-exports/`. Eleven obsolete cache archives predating
`519ce712` were removed after hash/stat and host-use checks, reclaiming 497,823,744
allocated bytes. Source, assets and binaries remain; the cleanup script is consumed.

Four superseded production/safe binaries from `62873a54` were also removed after
source, replacement/hash and host-use verification (186,937,344 allocated bytes).
Current `519ce712` binaries and all source, assets and evidence remain.

Original acceptance: default full corpus 2,486 passing roots plus the established
prerequisite skip; focused default 294, server 292 twice and high-resolution 294
twice all pass without skips. Exact root-name sets, source fingerprints and profile
environments are recorded in [the baseline](final-callback-exports-baseline.json).
Production qualification is reused from `519ce712` because only three new test
files changed. All current original runs are joined; source is safe to convert.
