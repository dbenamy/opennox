# Wall-list deletion correction

Generator integration exposed a pre-existing bug in the Go wall manager:
`serverWalls.find` started at the global wall-list head but advanced through
`NextByY24`, the row-list link. Deleting an older wall could drop other live walls
from the global list and eventually corrupt wall reuse. A ring-map integration
case timed out inside `GetWallAtGrid` after the position chain became cyclic.

The correction follows `Next20`, the global-list link. This changes no C source
and is an independent correctness fix, not part of translating orchestration.

## Direct and integration evidence

The original focused test inserted four walls and deleted each position in turn.
Deleting the oldest/second wall left only one/two global-list walls instead of
three. With the correction, all cases pass. The expanded regression has 12 cases:
four deletion positions across different rows, one shared row, and one shared
position bucket. It validates all indices, live/free membership, repeated deletion
and reuse; the corrected run passes in 0.032s.

The original C generator now completes 128 normal and 16 ring maps, plus seven
invalid-theme paths. These use real keyed synthetic theme files and the real
room, painting, object and waypoint owners. The single positive probe also requires
a PlayerStart and observes nine released room records. Three complete C integration
capture groups are mandatory; coverage of outer start/retry/save orchestration is
still pending. No orchestration routines have been converted.

## Audited historical capture updates

The painting conversion originally matched all 88 C capture groups. Fixing the
shared wall service changes 11 of those groups, covering 350 changed cases:
520 `Next20` words and 31 global wall-list heads. **Every other field is identical**,
including object state, geometry, RNG tails, floor state and the other wall links.
The updated captures pass live/free partition and both wall-index consistency
checks across all 673 steps in the affected groups.

Affected groups: connected-rooms, object-placement, rooms-19, rooms-24, rooms-25,
smoke-25 through smoke-29, and wall-floor-masks. Only these 11 expected hashes are
updated. Hash checks remain mandatory and fail immediately; the temporary
nonfatal mismatch reporting used to collect all differences was removed.
Historical C captures and their original hashes remain preserved locally and in
the earlier fixture revision. This is an intentional wall-service correction,
not a newly discovered C-to-Go difference in the painting algorithms.

## Qualification and recovery

Full qualification passes (`build/fix-wall-list/qualification.json`):

- 73,635 accumulated captured cases / 981 groups plus contracts pass in default,
  server and highres variants: 299.109s / 367.549s / 307.546s.
- Normal/HD/server production builds pass in 71.431s / 9.508s / 68.330s;
  ELF32/i386, SSE2/CGO and prior retired-symbol checks pass, with fixtures absent.
- Asset-backed full suite matches all 1,553 known failure entries exactly:
  15 packages pass, three fail and 32 skip, no added/removed failures.
- Fresh unchanged headless gameplay passes in 37.170s with null audio and
  `NOX_E2E_OVERRIDE=false`; run is `build/baseline/runs/wall-list-fix`.

Remaining C: **101,335 lines / 148 files / zero reference C**, no line delta.
The lookup-cache optimization was already qualified and pushed separately as
`55a4b2ef`; the completed growth conversion is `bd256d10`.

Evidence under `build/port-map-orchestration`: wall-list-original.log,
wall-list-corrected.log, wall-list-expanded.log, c-cached.log (ring timeout),
c-wall-fixed.log (ring success and first old-hash differences),
wall-capture-audit.json, wall-capture-deltas.json and
wall-list-capture-integrity.json. All original painting captures remain under
build/port-map-painting. Current expected hashes are in the committed fixture
source once this correction is committed.

A separate ownership review remains: the existing theme cleanup frees decoration
heads but leaves nested wall/floor records allocated. The integration fixture
records those survivors and releases them at fixture teardown. This correction
does not claim to fix that separate cleanup behavior.

## Stable integration tracing

The initial integration captures exposed a test-only identity problem: C-allocated
release snapshots could reuse storage from freed engine records. The normalizer
then identified dangling references as diagnostic records. Release snapshots now
use retained Go-owned buffers, never passed to C, and grid-row IDs are recorded
before release. The audit changes 66 normal-map trace words and 1,088 ring-map
trace/global references; every other field is identical. Changed globals are the
already-disposed grid index and scratch buffer. No production behavior changed.

Eight allocator layouts each compare complete captures for a normal map and two
ring maps. Those checks and all 151 integration cases repeat three times with
mandatory hashes in 11.036s. The final hashes are in
`src/map_orchestration_corpus_porttest_test.go`; the earlier preliminary captures
remain local for the audit (`trace-identity-audit.json`). Full combined variant
qualification must also reproduce these hashes before the C baseline is committed.


The first combined run after trace isolation exposed two text-table words that
numerically overlapped Go trace-buffer addresses. Diagnostic buffers now carry a
`captureOnly` marker, excluding them from engine-pointer lookup while retaining
all their captured contents and guard checks. A direct identity regression proves
both exclusion of diagnostic addresses and preservation of engine byte offsets.
The complete combined/isolated difference was exactly those two text words; no
map state differed. `combined-isolated-delta.json` preserves the evidence.

The focused rerun with diagnostic exclusion passes three repetitions: integration
11.658s, legacy identity contract 0.015s. Locked hashes remain unchanged. Full
qualification was restarted after that pass.

All three final combined variants reproduce the locked C hashes.
