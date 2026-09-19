# Client audio streams, cache and driver queues — native qualified

All 70 selected C bodies are now Go: sample metadata and bag/WAV reading, fixed
block pools, cached sample chunks and references, driver/context/voice ownership,
voice selection, buffer traversal, callbacks and timer dispatch.

**C remaining: 22,400 physical lines / 66 files /zero reference C**, a reduction of
**1,261** from the qualified C baseline. Original selection:1,109 body lines;
prerequisite corrections increased that to1,122. See [C_LOC.md](C_LOC.md).

## Scope and interfaces

Scope was GAME2_2 addresses 486640–487D60 and GAME3_1 addresses 4BD280–4BDC00.
The whole-repository caller/callback audit found 25 external roots and all 70 bodies
reachable. No orphan pruning was appropriate. The initial provisional selector
included three indented call sites; the corrected selector requires a declaration
at column zero and unique names. Only the corrected tracked audit is authoritative.

Layout-asserted Go types retain the catalog, chunks/cache entries, drivers,
contexts and voices needed by remaining C owners. Allocation ownership remains
on the C heap. Go wrappers call Go helpers directly. Twenty-two C interfaces remain:
18 external C calls and four stored callbacks. Forty-eight private interfaces and
the driver C global `dword_587000_155144` retire. Its Go getter resolves the same
mapped storage; the fixture replaces that actual owner. The test bridge continues
to exercise retained interfaces through C.

See client-audio-streams-selection.json, client-audio-streams-callers.json and
client-audio-streams-native-qualification.json for exact scope and interfaces.

## Independent contracts and frozen comparisons

The fixtures use actual allocations, intrusive lists, catalog file handles and
timers, plus controlled external driver callbacks. The 16 frozen captures cover 927
cases: pool layout/exhaustion/LIFO reuse and payload preservation; reference-count
boundaries; sample formats, unsigned volume multiplication and signed byte rates;
partial/short reads and signed length limits; WAV override/fallback/close ownership;
cache hits, eviction, pinned entries and failed fills; buffer traversal; multiple
devices, context reuse and format replacement; voice capacity and constructor failure
cleanup; bulk creation, priority/level selection, callback results and flag changes;
timer composition, conditional driver updates and clock boundaries.

All three original-C target captures matched before freezing. **The first native
build passed all 16 captures /927 records unchanged.** No goldens were regenerated
to accommodate the translation. The C baseline is **2f4dbfea**, and its full
qualification checkpoint is **419a85dc**; both are pushed.

All affected default/server/highres targets pass 34 root-package tests plus 5 timer
package tests each. Their 19 capture files /1,533 records match C exactly. Existing
audio-asset/list tests also validate embedded frozen hashes; their optional artifact
output was not enabled. All 2,569 final native source fingerprints match across
targets and production. The static memory-map test passes.

Three fresh native binaries pass ABI checks, including retained callbacks and
retired symbols/global. The full suite matches the exact known 1,553 failure entries
and package outcomes 32 skipped /15 passed /3 failed. Fresh audio-enabled gameplay
and explicit save/load pass. Physical audible quality remains a release check.

## Corrections and compatibility decisions for review

- The original WAV reader provided only 36 bytes for an override path and aborted
  on an ordinary longer fixture directory. Before freezing, C was corrected to use
  a separate bounded 280-byte path; unsupported lengths fall back to packed audio.
  Missing, short, truncated or zero-channel format chunks now close the override
  and preserve the bag fallback, instead of using incomplete format metadata.
  These reversible prerequisite corrections added 13 physical C lines and received
  fresh C production qualification before translation.
- The C clock adapter returns 32 bits before promotion into 64-bit context timing.
  The port preserves that narrowing, including high-word and wrap behavior. The
  fixture controls the legacy clock separately from the Go timer clock.
- A failed cache fill can retain the active catalog stream until the next open or
  close. Preserve that existing ownership behavior; cleanup/reuse closes it.
- A zero bulk-voice request fills available capacity because of the original
  do/while behavior. Preserve it; production initialization uses positive counts.
- Low-level pool counts remain positive, as required by actual callers. Signed
  read-boundary fixtures provide enough storage for the effective read request;
  they do not force artificial multi-gigabyte allocations.

## Evidence and recovery

Artifacts are under `build/port-client-audio-streams`:

- `dispatch-final`, `c-repeat-server`, `c-repeat-highres`: matching pre-freeze C.
- `c-final-{default,server,highres}`, `c-qualified-production`: qualified frozen C.
- `native-focused`, `native-final-{default,server,highres,production}`: native gates.
- `static-c-qualified.log`, `static-native.log`: static memory-map checks.
- Tracked capture index and C/native qualification JSON files summarize evidence.

The first test-only callback bridge needed a C accessor for static callback
addresses. The first C production command referenced the previous batch manifest;
its scenario name was already in use. The corrected manifest was rerun successfully;
only `c-qualified-production` qualifies C. These failures are retained locally.

All fixture drafts, native drafts, `freeze.py`, `generate-native-adapters.py` and
`install-native.py` are **consumed**. Never replay them; installed source wins.
The original archive/assets and qualification captures remain untouched.

Before this batch, 164 regenerable compiler-cache artifacts older than six hours
were removed with no compilers active, reclaiming 8,618,371,198 bytes. The manifest
is `removed-stale-go-cache.json`; that cleanup is consumed. Completed C/native scenario copies reclaimed another 2,225,495,402 bytes after
hash/inactivity checks. `deduplicate-client-audio-streams-completed-assets.py` audit
and apply modes are consumed; `--restore NAME` and per-run manifests preserve
restoration. All cleanup sessions are joined.
