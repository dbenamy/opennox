# World grid, wall storage and map serialization

Qualified parent `df2e1155` is pushed: **9,267 C lines /39 files /zero reference C**.
Original-C baseline `d4a6a628` is qualified and pushed. The native conversion is
fully qualified: **8,820 C lines in 38 files**, zero reference C (447 fewer).

The connected selection covers tile-grid allocation/lookup and water eligibility,
door/wall attachment and secret-wall lists, wall bounds and map serialization,
and the remaining private numeric/table helpers used by world motion. Move the
shared grid, tile-definition and secret-list owners with their Go consumers.
Read GRID_LOOKUP.md and FLOAT_INT.md before retiring staged implementations:
tileAtPoint is already qualified; the original C route was retained to avoid
extra callbacks while C callers still existed. Its remaining C water-predicate
caller moves in this batch. Keep the independent grid/IEEE contracts; retire
C-only algorithm and benchmark paths once their real callers have moved.

The current audit also finds float2int16 referenced only by its porttest adapter,
and double2float/wall-bound accessors without production callers. Finish symbol,
address and registration audits before retirement. The int32 and double-to-int C
converters still have client decoder callers and remain. The secret-wall lookup
still has C decoder callers and needs a Go-backed export. Other private callers
should invoke Go directly.

Use the existing real tile/grid, wall, polygon, object and map/cryptfile fixtures.
Add independent list/mutation/lifetime checks and capture numeric scratch state,
including negative zero, fractional rounding and nontrivial initialization
constants. Numeric helpers use mapped scratch words; preserve observable writes
rather than replacing only their return values. Audit actual C arithmetic and
shipped constants before freezing.

Map orchestration currently substitutes serialization through a porttest-only
link wrapper with a real fallback. Preserve that fixture boundary when its private
production callers become Go, and exercise real serialization separately. Audit
format prefixes, wall-coordinate shifting, cryptfile closure and compressed output.
Grid row-free leaves the outer pointer table to its caller; do not infer a complete
free from its name or change allocation ownership without tracing its callers.

The batch selection and caller inventory are provisional until fixture review.
The seven new focused expectations are now frozen after repeated original-C runs;
the broader three-target baseline is qualified.

## Original-C fixture development

At this baseline, production remained identical to `df2e1155`. The first
three focused original-C roots passed in `build/port-world-grid/second.log`:

- Numeric scratch: 8,352 rows covering both pointer aliases, raw IEEE inputs,
  negative early returns, fractional rounding and large values. A separate first
  process produces the same capture. PC53/nearest is explicitly installed and the
  full floating-point environment restored. Signaling NaNs quiet on the C scratch
  store; the qualified production disassembly shows an x87 load/store pair.
- Secret walls: 360 operation sequences cover all four-node insertion orders,
  duplicate and signed-boundary IDs, head/interior/tail/absent/repeated removal,
  clear and nil-next. Independent list and payload contracts use C-owned nodes and
  wall records; returned freed addresses are normalized without dereferencing.
- Trigonometric initialization: all 12,288 words use shipped scale constants,
  with independent range/formula/monotonicity checks, guards and five warm flag
  values. Warm calls must leave the tables untouched.

Attachment and real encrypted map-save fixtures are now included. Map-save tests use the real DebugData writer and actual
cryptfile/compression services; a deliberately failing extra section exercises the
failure path. Other section bodies retain their existing dedicated contracts.
The first save-test setup used an incorrect package import; corrected before it
compiled. This is a fixture issue and does not change production or expectations.

The completed drawable scenario assets were deduplicated after byte verification,
reclaiming 1,112,747,701 bytes. The consumed script and restore manifests are under
`build/port-world-grid`; original assets and the user archive remain untouched.

Review notes: secret-wall lookup promotes its signed 16-bit argument against an
unsigned wall ID, so IDs at or above 32768 do not match. Preserve that observed
behavior in this compatibility port; any gameplay correction is a separate review.
Map-save section failure leaves magic walls in their previous direction; success
restores their new direction. The focused real-save fixture confirms both paths.

The attachment fixture initially used a filtered wall lookup, which intentionally
hides door walls. Assertions now inspect the raw index and use non-broken input
flags. All fixture corrections occurred before baseline freezing.

## Focused baseline reviewed

All seven roots pass in `eighth.log` and in a separate process with three repeated
runs (`repeat.log`). Seven capture files match byte for byte and their hashes are
now literals in the tests and the C qualification manifest. The broader sweeps pass 294 default /293 server /294 highres roots, no skips,
and all 347 captures match exactly. The server-only exclusion is the explicitly
`!server` floor eligibility test. Static checks pass. Production identity is
verified against `df2e1155`; all three qualified binaries were rehashed. See
[world-grid-c-qualification.json](world-grid-c-qualification.json).

Additional coverage: seven allocator outcomes (outer failure, early/middle/final
row failures, success and no-trigger limit), 42 wall attachments, 12 encrypted
map-save cases and 216 drawable-save cases. The latter uses the real client list,
server type policies and object allocator; only the downstream xfer call is
observed. It checks copied fields, float narrowing, ignored xfer failure, unchanged
drawables, live-object balance and persistent pool netcodes. Map saving uses the
real DebugData section and actual encryption/compression services.

The first allocation observer reused global Go callbacks. Two combined probes
stalled in GC after that test; the precise runtime cause is not established.
The final observer records counters in C, is restricted to the locked test thread,
and injects failure only for the grid's exact allocation shapes. Existing theme
observation is unchanged. Four focused passes completed after that change.

Two independent fixture assumptions were corrected before freezing: an interior
point on the polygon's corner ray hits its already documented legacy containment
limitation, and the object allocator returns freed objects to the back of its pool.
The attachment fixture now uses a non-degenerate interior point; netcode checking
surveys every real pool slot instead of assuming immediate reuse.

Removed 15 verified superseded test-cache archives, reclaiming 1,111,445,656 bytes.
The audit/apply scripts under `build/port-world-grid` are consumed; production
archives, current test archives, qualified binaries and evidence were retained.

## Recovery

All original-C qualification jobs are joined.
The seven new hashes were enforced during every sweep; existing tests also checked
their committed expectations. The complete equal 347-file inventory was then
frozen for native comparison. Completed server/highres captures share storage with
default after hash verification; treat them as immutable and use fresh output
paths. The two capture deduplication plans are consumed (1,185,269,604 bytes saved).
The ignored numeric implementation draft and owner installer are consumed and
stale; do not replay them or fixture drafts over the reviewed source.

## Native conversion and review notes

C baseline `d4a6a628` is committed and pushed. The working conversion moves the
selected live behaviors and grid/tile/secret-list owners to Go. Four unused or
test-only C helpers retire, along with the now-unused Go-backed `sub_411350` C
bridge. The remaining C decoder secret-wall lookup keeps a Go-backed export.
`GAME4_2.c` is empty and removed. Working C is **8,820 lines /38 files /zero
reference C** (447 fewer), now fully qualified.

The first compile found an adapter pointer conversion and an unused import; both
were fixed before comparison. No expectation has changed. The current focused
native run includes the seven new roots and the existing grid/IEEE/subtile checks.
Ignored installers and the numeric draft are now consumed and must not be replayed.

Reversible decision for review: map saving rejects paths shorter than four bytes
or at least 1024 bytes. The old implementation subtracted four from an unsigned
string length and copied into fixed 1024-byte C buffers; those inputs were outside
its valid domain. Two native-only contracts check clean rejection without opening
files. Valid-path behavior stays subject to the frozen encrypted output captures.
The Go save path reacquires the global cryptfile after callbacks, matching the
original dispatch behavior. The orchestration fixture now substitutes the same
serialization boundary through a Go service instead of retaining a C link wrapper.

The first broad native default run completed 265 roots before exposing a Go write
barrier bug in the numeric scratch address store. The slot's old value is an
arbitrary ABI word (the test uses `0x13572468`), not necessarily a managed pointer.
Writing it through `*unsafe.Pointer` made GC inspect the old bits as a pointer.
The store now uses `uint32(uintptr(address))`, matching the original C uint32_t
write. This is a production conversion correction, not a changed expectation.
The failed output remains under `native-default`; successful reruns will use fresh
paths. The corrected focused comparison passes, followed by successful final
default and server sweeps.

A further 15 pre-conversion test-cache archives were verified and removed,
reclaiming 1,122,457,854 bytes. The corresponding cleanup plan is consumed.

The scratch correction passes all 14 focused roots. List insertion and predecessor
relinking also now write raw uint32 words. An additional native-only independent
contract seeds an uninserted node's next word with integer data; insertion must
overwrite it without inspecting the prior value, as the original C assignment did.
The resulting list/payload captures are unchanged. Final broad sweeps use fresh
`native-default-final` and `native-server` directories.

## Final qualification

Default/server/highres pass 294/293/294 roots without skips; all 347 captures
match the original-C baseline. Final static checks pass. Three fresh production
binaries pass build and ABI/export audits. The full suite matches exactly the
1,553 known failures. Headless gameplay and explicit save/load pass against the
previous drawable references. All source fingerprints match the final checkout;
no frozen expectation changed. All jobs are joined. See
[world-grid-native-qualification.json](world-grid-native-qualification.json).

Successful native capture sets share storage with the C baseline after byte
verification, saving another 1,777,904,406 bytes. Shared completed evidence is
immutable. Applied cleanup plans are consumed; use fresh output directories.

World-grid completed scenario cleanup reclaimed 1,112,747,701 bytes after verifying
original asset hashes. Per-run restore manifests preserve recovery; original
assets/archive and changed saves remain. The cleanup script is consumed.
