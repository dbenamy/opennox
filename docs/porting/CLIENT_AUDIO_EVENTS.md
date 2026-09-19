# Client audio events and playback — native qualified

All 62 selected C bodies are now Go: event allocation and cleanup, sample selection,
cache loading, scheduling, voice ownership and priority, music/dialog controls,
volume/pan updates, AIL sample refill and completion.

**C remaining: 21,083 physical lines /65 files /zero reference C**, a reduction of
**1,327** from the qualified C baseline. Scope contained 1,182 body lines; the
physical reduction also removes obsolete comments and private declarations. The
audevent C file is gone. See [C_LOC.md](C_LOC.md).

## Scope and interfaces

Scope spans GAME1_3.c, GAME2.c and client__audio__audevent.c. Whole-repository caller
and callback review found 35 external roots and all 62 functions reachable.
Adjacent GUI overlay helpers were excluded after inspecting their behavior.

Twenty-one C exports remain: 18 actual C entry points plus three stored voice
callbacks. Forty-one private interfaces and the cache/pool C globals retire.
Metadata, event and handle layouts have size/offset assertions; shared allocations
remain on the C heap. Go callers use native helpers. Test dispatch continues to
call the retained interfaces through C. Historical memory-map metadata remains.
Existing sound observers still run on the actual public playback path; test device
control observes the external queue while native refill and voice callbacks run.

See client-audio-events-selection.json, client-audio-events-callers.json,
client-audio-events-captures.json and the C/native qualification records.

## Contracts and qualification

Fixtures reuse real metadata tables, intrusive lists, event pools, cache/catalog
file handles, voices, timers, music modules and random generators. Independent
contracts cover format mapping; signed/unsigned volume arithmetic and pan clamps;
control switches; metadata defaults and guards; music stack capacity/nesting;
sample shuffle ordering, loop limits and RNG consumption; contiguous/chunked PCM,
scratch-buffer guards and end notifications; event/voice/cache ownership and reuse;
clock narrowing, deadline overflow and exact-deadline waiting; all 200 pool slots,
automatic reclaim, per-sound caps and manager frames; callback transitions, reload
modes, priorities/group eviction; reservation/start failures and retry; public
playback volume/pan/priority effects.

All three original-C targets produced identical **17 captures /820 records** before
freezing. **The first native build matched every frozen capture unchanged.** No
expectations were regenerated to accommodate the translation. Baseline **e30b952e**
and its production qualification **994d7b3c** are committed/pushed.

All native affected default/server/highres runs pass **51 root-package tests plus
five timer tests** each. Their **36 captured artifacts /2,353 records** match C;
existing audio-asset/list tests also validate embedded expectations. All **2,590
source fingerprints** match across the three targets and fresh production. Static memory-map checks
pass. A production-source audit confirms the 41 retired interface names are absent.

Fresh native production passes all three binaries and ABI checks, the exact known
full-suite results (1,553 failure entries; 32 skipped, 15 passed, three failed
packages), gameplay and explicit save/load. Every build/test session is joined.

## Compatibility decisions for review

- Preserve the existing 32-bit clock narrowing before 64-bit deadline arithmetic.
  The exact deadline still waits; only a greater clock value advances playback.
- Empty chunks encountered during refill can end playback before later chunks;
  an initial empty chunk advances once before checking the next result. The
  production cache emits nonempty chunks. Keep this captured convention.
- Serial zero identifies the first event and is also the cleared serial. Its
  handle can still validate after stop until reuse changes identity. Preserve
  existing behavior rather than changing handle semantics during this port.
- Public volume changes update immediately; pan changes set the next timer target.
  Preserve the timer adapter's nil/zero setter return conventions too.
- Preserve repeated selections in the manager, including their RNG consumption.
  Metadata sample counts remain bounded by the existing 32-slot owner contract;
  priorities/buckets retain the production caller ranges.
- The recording test device does not deliver asynchronous completion. The lifecycle
  fixture explicitly delivers the real stream end notification after stop before
  voice reuse. Its initial third-start failure was a fixture ownership issue;
  production behavior was unchanged.

There were no production algorithm corrections in this batch. The C baseline added
ten conditional test-adapter lines (22,400 →22,410 physical lines); they disappear
with the conversion. The first public-playback assertion expected immediate pan
updates and was corrected before freezing. One fixture compile used an incorrect
memory-map accessor signature; it was fixed before qualification.

## Recovery and evidence

Local evidence is under `build/port-client-audio-events`:
`c-final-focused-second`, `c-repeat-{server,highres}`, `c-final-{default,server,highres,production}`,
`native-initial`, `native-final-{default,server,highres,production}`,
`static-final-baseline.log`, `static-native-final.log` and the private-interface audit.
Earlier failed fixture runs are retained as diagnostics, not qualification.

All fixture/native drafts, build-dispatcher.py, freeze.py, install-native.py and
generate-native-adapters.py are consumed. Never replay them: installed source,
frozen expectations and qualification records are authoritative. Physical audible
quality remains a release check, as described in PORT.md.

Four completed C/native scenario copies were verified against original assets and
deduplicated, reclaiming 2,225,495,402 bytes. Their captures, modified save files,
binaries and per-run restoration manifests remain. Audit/apply operations are
consumed; use the recorded helper’s `--restore NAME` only when needed. Original
assets and archive were untouched.
