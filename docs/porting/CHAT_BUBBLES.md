# Client speech bubbles — qualified native

Qualified C baseline: **bf03e4d6**. All thirteen selected C bodies and two private
C globals are replaced by Go. Production C is **25,368 lines /67 files /zero
reference C**, a reduction of **779** from corrected C. The original tail-removal
bug was corrected before the C baseline. Only remove/clear retain C exports for
the live decoder; all other adapters use Go directly.

Native focused contracts match all twelve frozen captures /552 records. Affected
default/highres pass 151 roots /1,707 entries; server 150 /1,706 (its existing
client-only occlusion exclusion). All 65 captures /17,984 records match C and each
other. The three target and production gates use identical 2,514-file source.
Static, all three production binaries/ABI, exact known full-suite failures,
headless gameplay and explicit save/load pass. See
[qualification](chat-bubbles-native-qualification.json) and
build/port-chat-bubbles/native-*.

The sections below preserve the baseline history, fixture corrections and recovery
notes. Every fixture draft, freezer and native installer mentioned is consumed.
Actual source wins; do not replay ignored mutation scripts.

## Scope and strategy

The live bubble lifecycle/layout/drawing graph has13 C bodies /738 lines in
GAME2_3.c. Nearby ordered-message functions are a separate subsystem and are not
included. See [selection/caller audit](chat-bubbles-selection.json).

Use the actual Chat allocation class, mapped head and C tail, real server frame/
tick rate, and existing client renderer/font/drawable owners. Capture list links,
text, lifetime, screen placement and rendered pixels. Preserve frame comparison,
clipping, overlap decisions and draw order. Audit allocation, null/empty state,
replacement, removal, clearing and expiration before freezing C captures.

Initial independent lifecycle contracts are installed. They cover replacement
identity/order, Unicode text, default lifetime thresholds, explicit durations,
frame wrap, head/middle/tail removal, append after removal and clear/destroy/reuse.
The original C explicit-remove path failed the tail invariant when removing the
last node. It left the tail pointing at freed storage, so subsequent append could
be detached from the live list. The two-line correction matches expiry cleanup.
Corrected lifecycle checks pass; production qualification is still required.
Original failure: build/port-chat-bubbles/lifecycle-original-ready/tests.jsonl.

## Geometry contracts and recovery

Overlap and all eight shift contracts use the real font's 11-pixel cap height.
The initial fixture assumed its 13-pixel nominal height; geometry-first retained
that fixture failure, while lifecycle and region checks passed. Corrected geometry
and 144 HUD placement cases are installed. HUD checks cover GUI visibility,
inventory states, summons, offscreen placement, inclusive C-facing rectangle boundaries,
and preservation of the other rectangle coordinates.

No golden is frozen yet. Layout/expiry/pixel coverage, repeated captures,
affected-target qualification and fresh production qualification remain. Actual
source wins over ignored copied drafts; do not replay old fixture drafts.

Initial fixture discovery required the same local create-function declaration used
by the existing client adapter and the actual noxrender.Viewport type. These were
bridge fixes before C execution, not production changes. Older discovery attempts
are retained under lifecycle-original and lifecycle-original-declared.

The root UI helper uses Go's half-open rectangle containment, but the C-facing
export dispatches to geometryRectInt and includes right/bottom edges. The first
placement run retained eight failures from that fixture assumption; the corrected
edge contract passes. No production rectangle behavior changed.

Ten groups pass in attachment-first: lifecycle, overlap/shifts, region boundaries,
HUD placement, viewport/font layout, unsigned expiry, pixel changes, drawable
class/16-bit coordinate attachment, visibility and candidate/arrangement gates.
Additional player-name/team-color, multi-bubble, clipped/hidden/empty drawing tests
are installed. Static mapped-memory checks pass. The complete original-C focused
run now uses the tracked chat-bubbles-c-focused.json manifest, recording captures
and source identity. No capture is frozen yet.

The initial player-color run stopped at the real team-color lookup because the
lightweight renderer had no color definitions. Reuse PortTestMinimapTeamColors to
supply/restore the normal definition input; do not replace the lookup. No production
change was needed. The failed run is c-focused-default; corrected runs use
c-focused-palette-*.

Affected selection includes chat contracts, client object/player rendering, screen
rendering, inventory display/open state, client entry, world geometry, teams and
minimap. These cover the live name-hiding/draw callers and shared renderer, list,
player/team, palette and HUD owners. The current selection is tracked separately
from the focused contracts. Production uses the previous qualified ABI exclusions
plus all thirteen current C bubble symbols, with fresh gameplay and save/load.

The complete default focused run passes 12 roots /563 test entries, producing
12 captures /552 records. Server/highres repeats are pending. The C-facing tail
fix remains the only production change. Raw tests cover real allocation/list
ownership, viewport/font geometry, default/explicit expiry and wrap, actual
static/dynamic drawable lookup, class flags, text/position pixels, player names,
team colors, multiple bubbles, clipping, hidden and empty lists.

Nine completed map-sections/map-metadata test logs were losslessly compressed,
reclaiming 395,277,454 bytes. Each decompressed stream was verified against the
original before removal. Original paths/hashes and gzip paths are recorded in
build/port-chat-bubbles/compressed-completed-logs.json. Restore an individual log
with gzip -dk if a historical audit needs its old path. Compression script is
consumed. Captures, reports, original assets and archive were unchanged.

## Frozen C checkpoint

All three focused target runs pass 12 roots /563 entries with identical source.
Their 12 captures /552 records are byte-identical. Expected SHA-256 values are now
installed in the unchanged contracts and indexed in chat-bubbles-captures.json.
The freeze script and copied fixture drafts are consumed. Frozen affected-target
and fresh production qualification remain; this is not a converted/qualified batch.

Frozen checkpoint **5bf9ca6b** is committed/pushed. Broader qualification is active.
A further forty verified successful historical target logs were compressed using
the same lossless checks, reclaiming 1,844,486,695 bytes. The separate manifest is
compressed-historical-logs.json under build/port-chat-bubbles; its script is also
consumed. Total log cleanup: 49 files /2,239,764,149 bytes reclaimed.

Affected targets pass on identical 2,509-file source: default/highres 151 roots /
1,707 entries; server 150 /1,706 because TestClientObjectRenderOcclusion is excluded
by its existing client-only build constraint. All 65 captures /17,984 records match.
The first production binary builds, but the ABI manifest incorrectly listed the
still-C bubble symbols as Go-backed exports. Move those thirteen expectations to
retained_c and rerun production; no source or frozen expectation changed.

## Qualified C baseline

Frozen affected tests and fresh production qualification pass on identical
2,509-file source. All 65 captures /17,984 records match; all twelve focused hashes
are unchanged. Three production binaries pass ABI checks; the full suite matches
exactly 1,553 known failure entries and package outcomes. Fresh gameplay and
explicit save/load pass. Evidence: chat-bubbles-c-qualification.json and
build/port-chat-bubbles/c-final-{default,server,highres,production-cabi}.
All build/test sessions are joined. The Go draft remains isolated under build;
production source is still corrected C (26,147 lines /67 files /zero reference).

## Native conversion in progress

Qualified C **bf03e4d6** is pushed. The guarded installer is consumed. Thirteen C
bodies and two private C globals are replaced by five Go files; mapped list-head
storage remains shared with the existing renderer owner. Go adapters call directly;
only remove/clear exports remain for the C decoder. The original allocation class,
692-byte layout and offsets are preserved, with compile-time size/offset checks.
Geometry keeps the original mixed signed/unsigned comparisons and moving-bubble
visibility gate. Native static checks pass; focused validation is pending. Current
C: 25,368 /67 files /zero reference (-779 including obsolete comments/headings).

Completed C scenario copies were verified against original asset hashes and
deduplicated after qualification, reclaiming 1,112,747,701 bytes. Original assets,
changed scenario files and all evidence remain. Restore manifests are in each run;
helper: build/port-chat-bubbles/deduplicate-chat-bubbles-c-assets.py. Its audit/apply
modes are consumed. Native fixtures read original assets, not these scenario copies.

Native focused validation passes on its first build/run: 12 roots /563 entries,
all twelve frozen captures /552 records unchanged. Static checks pass. Affected
three-target native sweeps and fresh production qualification remain.

## Final native review

The first native focused run passed without changing any expected hash. The
conversion retains UTF-16 copy/termination, allocation reuse, list order, unsigned
frame comparison, renderer cap-height padding, inclusive HUD boundaries,
static/dynamic drawable selection, 16-bit coordinate truncation, historical
visibility/overlap decisions, two-pass outlines, text/name colors and draw order.
Compile-time checks enforce the 692-byte storage size and critical field offsets.
No C algorithm is retained solely for testing; remove/clear fixtures also exercise
the two retained live C bridges. Fresh production uses the same source as all
native target gates. The only production correction was the pre-baseline tail fix.

After native qualification, the two native scenario copies were likewise verified
and deduplicated, reclaiming another1,112,747,701 bytes. Their restore manifests
remain in the runs; helper: deduplicate-chat-bubbles-native-assets.py under the
batch build directory. Both C/native dedup audit/apply modes are consumed.
