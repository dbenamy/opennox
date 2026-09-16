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
The correction was qualified in all targets and real gameplay before translation, as recorded below.

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
hashes are frozen. Static final preflight passes. Final C qualification is complete. The baseline was committed and pushed as
5777c6d3 before the reviewed native drafts were installed.

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
All C source readers were joined before baseline commit/push and Go integration.

Disk maintenance deduplicated nine older completed spellbook/catalog/list replay
asset copies, reclaiming about 4.9 GB with verified restoration manifests. Original
assets and active runs were untouched; screenshots, logs and changed outputs stay
local. No qualification failed from disk exhaustion.

## Native conversion

The nine readers are implemented in map_drawables_state.go and
map_drawables_types.go; map_drawables_exports.go retains only the actual raw
map-section callback. Eight private C interfaces and **617 C lines** are removed.
Working count: **61,065 lines / 82 files / zero reference C**.

The first native focused run passes all eight roots and seven frozen hashes in
112.039s. No compiler or behavior corrections were needed. Source review checked
read versus seek accounting, allocation order, version/count signedness, team
registration, class dispatch, door coordinate arithmetic and light calculations
against C. Static final checks pass. Completed target and production qualification is recorded below.

Completed C replay asset copies were also deduplicated using verified restoration
manifests, reclaiming about 1.66 GB. Reference frames, logs and modified outputs
remain available; original assets are untouched.

Affected native qualification passes **44 default / 42 server / 44 highres** roots,
with no skips and all seven capture hashes matching. Driver seconds:
**146.035 / 143.357 / 201.454**. The server build excludes
TestClientObjectRenderOcclusion and TestWorldWallsFieldOfView at compilation;
these are not skipped tests. All target source readers are joined, and each phase
reports unchanged source. Production qualification also passes, as recorded below.

Production qualification passes all three 386/SSE2/CGO builds, ABI checks, and
absence of test helpers. One map-section callback is retained; all eight private
C interfaces are absent from source and production binaries. The full asset suite
matches exactly **1,553 known failure entries**, with **15 pass / 3 fail / 32 skip**
package outcomes. All **41 gameplay frames** match the qualified C baseline;
flat gameplay matches all **14 frames** and exact warrior map regeneration.
Production driver took **277.154s**, then flat gameplay **49.863s**. Gameplay
process times were **63.020s / 49.320s**; exact timings and binary/frame
hashes are recorded in map-drawables-replay.json.

Final native-index-proof.json checks all **1,863 staged source files** against the
three affected phases and production. All four phases report unchanged source;
all source readers are joined. Measured final C: **61,065 lines / 82 files / zero
reference C**, a reduction of **617 lines**. No frozen expectations changed.

The next candidate is colored-light animation and viewport updates, which consume
these loaded fields. Its six functions form a small coherent follow-up with
reusable real drawable/light owners; do not expand it into unrelated resource UI
merely to meet the target batch size.
