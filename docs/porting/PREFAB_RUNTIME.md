# Prefab and map runtime

## Scope and status

Forty connected functions cover prefab metadata/files, cache construction and
placement, group dispatch/serialization, waypoint serialization and generation.
The C baseline **86654f0f** is committed and pushed. The native conversion is now qualified. Historical prerequisite
and baseline evidence follows. Current action state is in
[PORTING_STATE.md](../../PORTING_STATE.md).

## Reproduced prerequisite defects

- A real call through `nox_xxx_mapReadSection_426EA0` panics at its `TODO`.
  The C prefab loader calls this for every nonempty section. Connect the bridge
  to the existing Go map-section registry, using the active file and bounds.
  Unknown names return 0/error 0 for the caller's object fallback; recognized
  successful sections return 1/error 0, and reader errors return 0/error 1.
  The error output is assigned on every call. Evidence:
  `section-reproduce.{json,log}` under the ignored batch directory.
- A C-created object cache node passed to the real Go cache destructor panics
  with `incorrect free`: C `calloc` bypasses the tracked Go allocator used by
  that destructor. Four cache-node constructors have this same mismatch.
  Use the same C allocator in the four Go node-release paths. C constructors
  and removal paths remain unchanged; group references retain their separate
  tracked allocator. Payload ownership
  is separate and still needs the larger batch's lifecycle audit. Evidence:
  `cleanup-reproduce.{json,log}`.

The changes are reversible prerequisites under the standing authorization for
confident local corrections. They do not retain C algorithms solely for tests.

## Independent contracts

Four root tests exercise the actual C entrypoints and real Go owners:

- Unknown/empty/case-sensitive section names and stale error-output canaries.
- Successful DebugData decoding into the real owner, empty data, rejected
  version, missing version, bounds preservation and existing short-read behavior.
- Real nonempty prefab files with successful/rejected sections, decoded-cache
  status, proof the seeded loader was bypassed, and file closure.
- All four C cache-node constructors followed by the real Go destructor and a
  second empty cleanup. Original C payload allocations are explicitly released
  by the fixture; this test establishes wrapper ownership, not payload disposal.

The existing DebugData reader accepts some short reads: a one-byte version can
supply version 1, and a partially read string is zero-padded then NUL-trimmed.
The first expectations assumed stricter I/O handling and failed; the bridge does
not implement a new reader, so the contract now explicitly checks the existing
reader's behavior. Review stricter section validation separately. No frozen
expectations were changed to accommodate this finding.

Initial fixture compile corrections were a required `unsafe.Pointer` conversion
for a C string and a 386 constant-width cast. These did not change production.

## Qualification

All three targets pass **111 roots / 1,756 tests including subtests**, with no
skips: default **22.787s**, server **22.719s**, highres **22.049s** package time.
All **59 captures / 21,698 records** match byte-for-byte across targets, and the
existing frozen expectations pass unchanged. The prerequisite manifest selects
connected tests covering
prefab/population, map rooms/painting/polygons, object transfer and waypoint
contracts. Discovery/execution must agree and no selected test may skip. Fresh
production qualification is required because these prerequisites change
production source. Preserve the exact known asset-suite failure set and run
headless gameplay, save/load and flat-map regeneration.

An initial tracked-node allocation approach passed the focused contracts but
failed the pre-existing marker-removal fixture, which supplies C-allocated cache
nodes. The narrower final repair matches the destructor to the existing C
allocation/removal ownership. The interrupted server/production gates and failed
default run are superseded; none qualify the final source.

Fresh production qualification passes in **377.678s**: all three binaries and
ABI/interfaces, exact known asset-suite failures (1,553 entries / three packages),
headless gameplay, save/load and flat-map regeneration. All four gates used the
same **2,289 source files**, unchanged throughout; all sessions are joined.
Production C remains **38,498 physical lines / 74 files / zero reference C**.
See [qualification evidence](prefab-runtime-prerequisite-qualification.json).

Static mapped-memory preflight passes. Raw logs remain under
`build/port-prefab-runtime`. Do not rerun consumed source installation commands;
read actual source and the latest checkpoint first.

## Larger-batch audit notes

The shipped script descriptor table contains 36 named records (IDs 0–35) and a
zero-name sentinel at 36; arguments include all eight kinds. Use the shipped
nonzero table in script-file contracts, not zero-filled fixture storage. Preserve
C's tile-orientation low-byte interpretation, the wall-group dispatch fallthrough,
version signedness and file positioning until independently qualified. Audit
cache payload lifetime and the unknown-object loader failure's file closure.


## Larger C baseline in progress

Prerequisite commit **4f103725 is pushed**. New test-only dispatch/global/table
helpers and scalar/file/group/waypoint contracts are being added; no captures are
frozen. Metadata, paths, library enumeration and all shipped script descriptors
currently pass 2,880 capture records. Group read/write contracts and actual
waypoint link resolution are also passing.

A subsequent real lifecycle regression reproduced `incorrect free` when a
prefab-created waypoint is transferred into the server and destroyed by normal
server teardown (`c-lifetime-before.log`). Its allocation helper used raw C
allocation while the receiving owner uses the tracked Go allocator. A repair
changes that existing Go helper to the receiving owner's allocator and releases
unplaced waypoint payloads during cache cleanup. The allocation-byte ABI fixture
uses the same matching release. The new disposal contract observes actual
allocation/free balance for both placement states. This repair is **not yet
qualified**; a fresh production gate is required before freezing the larger C
baseline. The previously qualified wrapper repair remains unchanged.

Baseline fixture corrections so far: anchor generated adapter prototypes to
actual declarations, check both path setters' defined 1/0 returns and nil clearing,
and use `GameHost | GameFlag22` for the actual group's allocation admission.
These corrected independent fixture assumptions, not production algorithms.


The expanded exploratory suite passes **23 root tests / 8,823 capture records**
before the loader expansion. Coverage now includes complete cache allocation
records, file open/rewind/seek/close, actual object pending-list placement, group
traversal through real owners, waypoint generation, all 32 serialized waypoint
link slots, signed section versions, introductions, and relative-coordinate
mutation. Floor/wall placement reuses guarded painting owners. Its overlay input
was corrected from an undefined fixture border index to a real defined border;
a separate assertion now requires the overlay allocation to be observable.

Complete prefab-file contracts reproduce open files on bad magic and unknown
object types in every tested header version (`c-loader-before.log`). Two C failure
exits now invoke the existing close helper; the focused rebuilt run is pending.
This is an intentional pre-baseline correction for review, alongside the waypoint
allocator correction. It adds two temporary C lines before conversion. Valid
DebugData, GroupData and WayPoints sections, bounds offsets and optional header
extensions passed the same independent contracts.

Payload audit: the old destructor also drops tile/wall payload pointers without
freeing them. Placed secret-wall data still refers back to the cached wall record,
so simply freeing all wall payloads would introduce a dangling pointer. Preserve
this evidence while qualifying the C baseline; resolve the transfer/back-reference
and allocation balance together in the native ownership work, with an independent
regression. This known issue does not justify changing frozen functional captures.


## Qualified larger C baseline

The final original-C suite repeats **18 new captures / 8,880 records** in separate
processes before locking expected hashes. All three target sweeps pass **133 root
tests / 1,780 including subtests**, with **77 byte-identical captures / 30,578
records** each. Package times: default **27.100s**, server **26.746s**, highres
**27.184s**. The previous ownership and file-close repairs now pass their focused
contracts and all broader gates. Initialization, direct instantiation and the
inverted selected-loader return convention are covered too.

Fresh production passes all three builds/ABI/interfaces, the exact 1,553 known
asset-suite failures (15 passing, three failing, 32 skipped packages), headless
gameplay, save/load and flat-map regeneration. The first production invocation
stopped before gameplay because a copied scenario name already existed. Corrected
names resumed with source/binary/gate identity checks, reusing the fresh successful
build and suite evidence. All final gates share the same **2,314 source files**;
static mapped-memory checks pass. No sources changed during any gate.

Production C is **38,500 physical lines / 74 files / zero reference C**. The scope
is now 1,371 live body lines after two failure-close statements. See
[selection](prefab-runtime-selection.json), [captures](prefab-runtime-captures.json)
and [qualification](prefab-runtime-c-qualification.json). The 40 algorithm bodies
remain C at this baseline. Proceed with the native conversion and ownership audit;
keep all frozen expectations unchanged. Baseline freezer/finalizer are consumed.

Disk maintenance reclaimed **3.092 GiB** from six completed prerequisite/quest
scenario copies after SHA-256 comparison with original assets. Changed maps/saves,
reports, binaries and original assets/archive remain. Per-run restoration manifests
are saved; build/port-prefab-runtime/deduplicate-completed-assets.py --apply is
consumed and must not be repeated. About **8.1 GiB** was free after cleanup.

## Native implementation and ownership corrections

Forty algorithms now run in Go. Ten C exports serve remaining callers; thirty
interfaces and two proven orphan bodies (`sub_51D100`, `nox_strnicmp`) are retired.
Fourteen globals move into Go; selected-prefab and instance counters remain shared
with live C generation/script callers. Go callers use direct adapters. The obsolete
population fixture linker wrapper is replaced by a build-tagged Go loader hook.
There are no retained C reference algorithms. Working C is **36,917 lines / 73
files**, a reduction of **1,583** from the qualified baseline.

The first linked native run (`native-third`) passed all 26 existing roots and
matched all **18 captures / 8,880 records byte for byte**. New independent contracts
then reproduced the known tile/wall payload leaks and two connected wall-transfer
defects: secret-wall data referenced the cache record/relative coordinates, and
breakable registration received wall data rather than the actual world wall.
Special-wall identity also failed to transfer. These were observed with real world
owners and allocation/free observers, for secret, breakable and combined flags,
successful placement and partial failure, plus repeated cleanup.

The correction transfers secret-data ownership to the world, updates its back-reference
and grid coordinates, transfers special-wall identity, registers the actual world
wall, and clears replaced secret state. Clearing the cache's data pointer records
ownership even if a later wall cannot be placed. The destructor now frees tile and
wall payloads, plus secret data that was never transferred. A failed wall-payload
allocation unlinks/releases its wrapper and returns failure. These are intentional,
reversible corrections; frozen functional captures are not changed to hide them.

Native bounds checks reject metadata names exceeding 63 bytes, group/waypoint names
exceeding their 75-byte buffers, missing version-2 group name tokens, script names
exceeding 1,024 bytes, out-of-table script opcodes, invalid record lengths, and more
than 32 serialized waypoint links. Defined signed-version behavior, opcode 36's
sentinel row, the 31-link interactive limit, float rounding and short DebugData
reader behavior remain covered. Newly detected I/O failures return failure instead
of depending on uninitialized C locals. These undefined-input choices are recorded
for review and are not claimed as exact compatibility for malformed files.

## Qualified native conversion

Default, server and highres each pass **136 root tests / 1,791 including subtests**,
without skips. All **77 captures / 30,578 records** match the original C baseline
byte for byte. Package times are **23.492s / 23.825s / 24.097s**. The independent
secret-to-plain replacement contract also passes with no remaining allocation.

Fresh production passes all three ELF32/SSE2/CGO builds, required/retired interface
checks, the exact known full asset-suite failures, headless gameplay, save/load and
flat-map regeneration. All four gates share unchanged **2,323-file source**; static
mapped-memory checks pass. See [native qualification](prefab-runtime-native-qualification.json).
No frozen expectations changed. Production C is **36,917 lines / 73 files / zero
reference C**, **−1,583** from baseline. Baseline qualification remains recorded
separately; consumed installers/finalizers must not be replayed.
