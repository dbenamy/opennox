# Final engine callback exports

Status: baseline `eb4c6e21` committed and pushed; conversion installed,
converted preflight/qualification in progress.
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

All six accepted original logs are losslessly gzip-archived, reclaiming
134,705,152 allocated bytes. Restore with `gzip -dk` and verify hashes against
`build/port-final-callback-exports/original-logs-archived.jsonl`.

The first converted server/high-resolution suites caught an incomplete caller
migration in `sessionSaveMetadata`: it reads the metadata callback directly from
mapped offset 55956 and used `CallIntUPtr`, bypassing the new native identities.
`TestSessionEntryCharacterCount/valid` stopped at that route. The remaining default
sweep was deliberately stopped and all jobs joined before editing. Initial logs,
binaries and fingerprints remain under `contracts/`; the passing initial preview
and eight-test preflight are preliminary evidence only.

The correction shares native section dispatch between file framing and metadata
loading while retaining each original foreign fallback (`CallIntPtr` with nil or
`CallIntUPtr` with zero). Table guards, snapshots and return handling are unchanged.
The fixed-offset/hex-offset and particle-field audit found no further consumers.
Existing assertions and captures remain unchanged. A 28-root preflight includes
the original eight contracts plus session-entry regression tests; all source-bound
qualification will be repeated on the corrected code. This was a primary audit
miss, not a helper result. PORT.md now explicitly requires table-offset tracing.

Additional disk recovery removed three superseded `62873a54` test executables
(200,314,880 bytes) after qualified replacement/source/hash/host-use checks, and
losslessly archived the three prior qualified `519ce712` contract logs
(165,019,648 bytes). Restore logs with `gzip -dk` and verify the batch journals.
The first passing preview's 1,654 asset copies were verified and deduplicated
(559,964,160 bytes); originals/saves/results remain. Its cleanup path typo was
corrected before any deletion; completed cleanup scripts are consumed.

The corrected preflight passed all 28 exact roots and the static check. An initial
selector union treated the `^TestFinal` prefix as a literal name and therefore
ran only the 20 session-entry roots; discovery review caught the missing eight.
The corrected explicit selector was rerun and its exact names accepted. These
20-only preliminary results are retained separately as `preflight-fixed/`;
`preflight-final/` is the accepted combined regression run.

The corrected fresh preview also passed. Its 1,654 duplicate asset copies were
verified and deduplicated (559,976,448 bytes). Final full/default and focused
server/high-resolution qualification is now under `contracts-final/`; all earlier
failure and preliminary artifacts remain separate.

Before final qualification, five obsolete root/legacy cache archives predating
`eb4c6e21` were removed after host/compiler and file-use checks (333,918,208 bytes).
Four inactive first-attempt binaries were losslessly archived (121,671,680 bytes);
restore with `gzip -dk`, restore executable mode and verify the journal hashes.
Sixty further historical game-message capture groups were losslessly archived
(355,438,592 bytes), with all 1,033 hard-link paths recorded for restoration.
Original assets and current qualification binaries remain untouched.

A corrected server process later stalled during GC coordination while the
porttest allocator observer called back into Go. SIGQUIT preserved the stack;
it shows `themeTestReleased` / map orchestration release normalization. Isolated
original/converted replays and two complete corrected server replays passed,
as did the complete corrected default corpus and high-resolution suite. These
runs remain historical evidence; they do not qualify the subsequent fixture fix.

Review found a concrete scope defect, though the exact stall cause remains
unproven: a process-global allocation observer can enter Go on a runtime startup
thread before that thread's Go state is initialized. The runtime's `threadentry`
and the actual binary disassembly confirm that free precedes `crosscall1`.
An independent pthread probe reproduces the scope error (`[0 1 1 0]` instead of
`[0 1 0 0]`). The fixture observer/time state is now thread-local, with a balanced
OS-thread pin during the active observation interval. Repeated enable/disable
calls preserve existing idempotent cleanup and outer fixture pins. Production
allocation code and frozen expectations are unchanged.

All 29 callback/session/scope preflight roots pass with the repair. A separate
source-only copy restores original `eb4c6e21` callback code and changes only the
observer fixture plus its new contract. It passed twice per profile
against the expanded 387-root selection (385 server), including every allocator
observer owner family. This repaired-fixture baseline is committed
separately before fresh converted qualification. Nine obsolete rebuildable cache
archives were removed after hash/stat/host-use checks, reclaiming 636,747,776
bytes; journal under `cache-before-observer-proof/`, cleanup consumed.
