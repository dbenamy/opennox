# Final engine callback exports

Status: qualified. Production baseline `519ce712`; original callback
contracts `eb4c6e21`; independently qualified fixture repair `a853826c`.

This batch replaces the last ten selected legacy C exports with native Go routes:
three player-file section callbacks, map object-data loading, screen particles,
three tooltips, main-menu events and dynamically assigned FlameCleanse updates.
The algorithms already lived in Go. Mutable callback storage and the existing
foreign-callback fallbacks remain supported.

The player-file table uses stable native identities and dispatches its current
value at each old call site. File framing retains the nil pointer foreign ABI;
the fixed-slot metadata reader retains its zero uintptr ABI. Tooltips and object
updates use existing registries. Menu events and the sole production map-loader
consumer use typed functions. Particle traversal still caches the next link before
calling a draw handler.

Review decisions:

- Decode menu event code and arguments before reading the mutable hook; return a
  nil response for zero and preserve signed/raw return bits.
- Preserve table lookup order, framing, nil/zero foreign arguments, fixture IDs
  and captures. Native identities are addresses of nonzero-sized global bytes.
- Read only the first byte of options tooltip data. Keep the conversation tooltip
  as an installed no-op and FlameCleanse as a dynamic update assignment.
- Validate a nil map handler before handling an unsupported section; only a zero
  callback result is failure.

Eight independent callback contracts ran on the original routes. They cover menu
argument/return bits and hook evaluation order; installed options tooltips for all
256 low bytes with nonzero upper bits; the conversation no-op; particle traversal
mutation; foreign section return boundaries; map validation/dispatch order; and
FlameCleanse through actual object update dispatch.

Two issues were found and resolved before acceptance:

1. The first conversion missed `sessionSaveMetadata`, which reads mapped callback
   offset 55956 directly using `CallIntUPtr`. Existing character enumeration tests
   caught it. The correction shares native dispatch while retaining each foreign
   ABI. A table-offset/hex-offset and particle-field audit found no further readers.
   PORT.md now requires offset tracing as well as symbol references. This was a
   primary audit miss, not a helper result; expectations and captures stayed fixed.
2. A server run stalled during GC coordination with the test allocator observer on
   its stack. Review reproduced a separate concrete observer scope violation and
   repaired it in `a853826c`; see [THEME_OBSERVER_SCOPE.md](THEME_OBSERVER_SCOPE.md).
   Exact causation of that stall is unproven. Earlier runs and replays remain
   historical evidence. Final acceptance uses fresh runs after the fixture repair.

A preflight selector union initially treated `^TestFinal` as a literal test name.
Exact-name review caught the omitted eight contracts; the corrected explicit
29-root selector includes callbacks, session regression and observer scope.
Live log buffering also briefly looked like a stall. No signal was sent; completed
events proved normal progress. The driver now flushes root lifecycle events;
selection, execution and acceptance logic are unchanged.

Accepted qualification:

- Original callback code with the repaired fixture: 387 default/high-resolution
  roots and 385 server roots, twice/profile, no skips or failures.
- Converted default corpus: 2,487 passing roots and the one established
  `TestMapPopulationPrerequisiteProbe` skip among 2,488 exact roots.
- Converted focused server/high-resolution: 385/387 roots, no skips or failures.
- All 29 immediate contracts, static checks and safe build pass.
- All three production binaries pass Linux 386/SSE2/cgo and symbol/ABI checks;
  retired exports are absent and test helpers are excluded.
- Fresh headless preview and final production character creation/save/load/resume
  pass. The full asset suite matches the exact known failure multiset and package
  results. All 1,654 original asset hashes remain unchanged.

| Metric | Before | After |
| --- | ---: | ---: |
| Selected project cgo files, client/server | 80 / 81 | 72 / 73 |
| Selected legacy C exports | 10 | 0 |
| Embedded production C bodies | 20 | 20 |
| Tracked headers / physical lines | 157 / 2,741 | 157 / 2,731 |
| Standalone production / test-reference C lines | 0 / 0 | 0 / 0 |

Evidence: [original baseline](final-callback-exports-baseline.json),
[qualification](final-callback-exports-qualification.json) and
[dependency inventory](final-callback-exports-inventory-after.json).
Accepted converted gates share source fingerprints. The inventory's selected
production-file hashes also match; the fixture repair changed only porttest files.
Frozen captures are not regenerated. External SDL2/OpenGL/OpenAL bindings and the
32-bit target remain unchanged; zero selected C exports is not the full internal
C-glue milestone.

Primary owns this batch; Luna is quota-blocked and no substitute model was used.
Read-only next-batch notes identify native scalar types and their caller boundaries.
Local evidence: `build/port-final-callback-exports/`.

Disk recovery preserves original assets and source. Obsolete cache archives and
superseded older binaries were removed only after recorded hash/stat/host-use and
replacement checks. Historical logs, captures and six superseded callback test
binaries are losslessly gzip-archived with SHA-256/mode/path restoration journals.
Restore archives with `gzip -dk`, restore recorded executable modes where needed,
and verify original hashes. Completed cleanup scripts are consumed. Scenario
asset copies are deduplicated only after scenario acceptance; originals and
save/results remain, with per-scenario restore manifests.
