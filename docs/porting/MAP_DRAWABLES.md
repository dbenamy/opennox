# Client map drawable readers

Starting point: qualified and pushed options conversion **a1e88a75**, **61,682 C
lines / 82 files / zero reference C**. This batch covers nine connected map
readers in GAME3.c, 004ABDA0–004AD570 exclusive, initially **617 physical C lines**.
The unrelated disconnect and server-options windows remain outside this scope.

The readers construct client objects from saved map records: common fields,
colored lights, doors/box objects, team bases, pressure plates and minimap markers.
The actual map-section callback is referenced by legacy/maps.go; other functions
currently have only internal callers. Verify the final interface audit before
retiring helpers.

## Baseline and independent contracts

Use real file streams and the existing drawable allocation/type/index owner.
Exercise version boundaries, consumed-byte totals and file position, allocation
failure, object flags/team defaults, shape/light state, active lists and modifier
resolution. Add real asset loading and gameplay coverage. Do not freeze pointer
addresses, uninitialized stack bytes or unrelated allocation padding.

Audit the complete record before freezing: old and modern common headers differ,
outer and inner versions are signed 16-bit values, and some old count fields use
raw low bits rather than a floating-point numeric conversion. Older light records
have default arrays/radius correction; later records add angle data. Pressure
plate defaults and later overrides, four team-base modifier names, door wall
registration and minimap membership need separate contracts.

## Old-format defaults prerequisite

The old common reader sub_4ABDA0 stores its team byte and extra-flags word even
when the selected version does not read those fields. The team local is read only
for outer version >=20; extra flags only for outer >=40 and inner >=3. Earlier
records therefore copy uninitialized locals into drawable state.

The first independent regression uses complete valid old record layouts and the
real file and drawable owners. It expects absent fields to default to zero,
checks exact consumed bytes/file position, and prevents unrelated team registration
while inspecting the resulting team byte. Correct the C defaults before capturing
old-format results. This is a scoped, reversible prerequisite, not a reason to
freeze unstable bytes or ask the user to choose their values.

All **42** independent missing-field cases fail against unchanged C. The two
locals are now initialized to zero, without changing C LOC. All 42 pass afterward.
No map-reader Go translation is installed yet. The correction is qualified in all targets and real gameplay, as recorded below.

## Frozen C capture set

The current set contains **2,484 records / seven groups**: 732 common-record cases,
240 typed-object cases, 36 map sections, 168 count boundaries, 96 team-registration
cases, 1,176 light boundaries and 36 door-coordinate boundaries. Additional
contracts cover the 42 absent-field cases and write-mode behavior. Hashes are in
map-drawables-captures.json; commands are in map-drawables-batch.json.

Typed records exercise actual drawable factories, light/shape operations,
modifier-name lookup, minimap lists and door wall registration. Default-zero
fixtures use complete records; no truncated inputs or uninitialized output bytes
are used as an oracle. The actual drawable factory retains its active bit, and
coordinate conversion truncates using the already-qualified converter. The
modern record orders its team byte differently from the older layout. These
fixture assumptions were corrected against the source before freezing.

Count boundaries include uint32 accumulated-byte wrap and counts around 16,384.
The old reader narrows four times the count to uint16 before skipping; the modern
reader retains the full product. This difference is intentional compatibility.

Review later: if drawable allocation fails after consuming part of a section
record, the legacy typed reader returns zero. The section reader then skips the
full declared length and can misread the next record boundary. The controlled
failure cases stop with return zero and a load-error flag, at reproducible file
positions. Preserve that existing behavior during this translation; any later
change should define failure recovery for the whole map-loading owner. This is
separate from the absent-field prerequisite.

## Final fixture review

Team registration uses actual membership lists, including creation of a missing
team and cleanup. It covers host/client mode, the legacy class/mode exception and
the modern FlagMarker exception/cache. The original drawing fixture omitted team
creation's string manager and printer; its new focused owner supplies those real
dependencies. This was a fixture correction, not an engine change.

Independent light-boundary assertions exposed missing production constants in the
initial fixture. A focused owner now copies and restores the actual embedded
intensity/angle and door-direction tables. No zero-filled stand-in is used. Light
cases cover both legacy intensity representations around 63, signed fixed-point
values, NaN/infinities and the entire range of angle conversion boundaries.
The typed capture was refreshed in C after adding real tables and explicit wall
coordinates; the earlier development capture is not the final oracle.

The door reader stores X into a signed local before dividing by 23, while its
inline Y expression divides unsigned. Boundary cases near the origin distinguish
both rules, including omitted wall registration when adjusted Y is negative.
Preserve this asymmetric arithmetic explicitly in Go.

All eight root fixtures pass together in c-doors-develop.log. The final seven
hashes are frozen. Static final preflight passes. Final C qualification is complete. No Go map-reader
implementation is installed. Ignored native drafts are incomplete review artifacts
until the baseline is committed and pushed; do not copy them into source blindly.

Evidence and historical drafts: build/port-map-drawables.

## Qualified C checkpoint

All eight roots complete without skips in default, repeat, server and highres;
all **2,484 records / seven hashes** match. Driver seconds: **31.404 / 10.443 /
121.987 / 36.897**. The affected selection completes **44 roots**, no skips, in
21.557s. It includes object rendering, door/geometry/equipment helpers, minimap,
world walls and light properties. Static final preflight passes.

Real C map loading matches the previous qualified options build's **41 frames**,
then repeats those 41 frames exactly. The flat scenario also matches, with exact
warrior map regeneration. Process seconds: **59.842 / 57.483 / 48.695**; full C
gameplay phase **231.087s**. Metadata is tracked in map-drawables-replay.json.
The runner's inherited flat display-description label was corrected in the
manifest; it is diagnostic metadata and changes no scenario behavior.

The C initializers change no physical line count: **61,682 / 82 / zero reference
C**. Final c-index-proof.json records source identity against all six phases.
All source readers are joined. Commit/push this checkpoint before installing Go.

Disk maintenance deduplicated nine older completed spellbook/catalog/list replay
asset copies, reclaiming about 4.9 GB with verified restoration manifests. Original
assets and active runs were untouched; screenshots, logs and changed outputs stay
local. No qualification failed from disk exhaustion.
