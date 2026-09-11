# Movement-path execution and waypoint graph — 2026-09-11

Connected scope: 50D2A0, 50D2E0, 50D3B0, 50D5A0, 547F10,
547F20 and BuildWaypointPath547F70. Their remaining owners are native
navigation and roaming. Port the whole group and remove the unused
waypoint-eligibility C export. Retain shared direction conversion and the
already-native detailed pathfinder/map-ray implementations.

Original C passes 17,585 execution cases and 2,132 graph cases. The shared
C-backed guarded AI fixture now covers explicit path points/indices, waypoint
lists, cached targets, unsigned timers/wrap, selected-index scratch, actual solid
and transparent wall rays, running multipliers and nonfinite speeds. Independent
assertions check wall behavior, graph-owner counts, nearest-point rounding and
the unrounded opening distance cutoff. Successful graph construction also runs
through the real movement owner, including paths beyond its 16-slot buffer.

Execution state SHA256:
f2ae3862c6ea077add1164716de01de0ec26b2e5ccf0757517d3e3daf15c8195.
Graph state SHA256:
9aaae7e2192fbc77655d061d7399e5f24dc2e45bccbd348db98ec2d54a09f3d2.

The graph fixture saves/restores the standalone epoch, 256-word scratch array
and one-shot flag, guards each waypoint/output allocation, and normalizes all
pointer-bearing state. An independent layer traversal checks route ordering,
cycles, duplicate edges, epoch collisions/wrap and capacity behavior. Test chains
include 255/256/257/258 nodes. Waypoint traversal state lives at offsets 504
(epoch), 508 (parent), 512 (next frontier); unrelated map links at 496/500 must
remain unchanged. C writes output index capacity before checking capacity;
direct tests allocate capacity+1 words and verify that extra write. The normal
16-slot owner aliases it with Field91 and then resets Field91. Preserve it.

Disassembly confirms a float32 nearest-distance threshold, full-double selected
length calculation, float32 denominator and force-X delta, and float64 force-Y
delta. The runtime epsilon at blob 581450+10288 is double bits3f847ae140000000
(float32(0.01) promoted to double). Isolated tests otherwise leave it zero, so
this fixture explicitly installs, verifies and restores the runtime value.
The waypoint follower overwrites Path[count-1] after the detailed-path call;
fixture success results obey its nonempty-success contract. Actual movement
reads Path[1] for direction when index is zero, including one-point active paths;
the allocated slot is populated in the fixture. No algorithm repair is bundled.

Two 200,000-call checks preserve complete final state for actual movement and
its move-path owner. Baseline times were 2,686 and 3,020 ns/call, respectively;
these fixed-position caller measurements do not predict whole-game speed.
Both final-state hashes are
64e97f956b16ea81f659b3bf00ec23fdc11038c26102a63341bfbbefe8f4af4d.
Existing navigation, roaming, guard/escort and waypoint checks pass alongside
the new corpora. Artifacts: build/port-ai-path-execution (ignored); C snapshots
are c-path-cases.json and c-graph-cases.json, and logs are c-qualified.log and
c-repeated.log. Recover the original implementation through the baseline commit.

At the original baseline, production C was 139,165 lines, 153 files, zero
reference C. Prior navigation batch is
fully qualified, committed and pushed as e32982f7. No user question is pending.

Native conversion matches original-C baseline bc1e662a: both corpora, independent
assertions, repeated-update hashes and all existing AI/waypoint checks pass.
Seven C bodies and four unused C exports are retired (waypoint eligibility,
nearest-waypoint lookup, detailed path setup and retreat generation). Native
roaming/navigation call private Go path helpers directly. The graph still shares
the original epoch/scratch/one-shot storage; no test-only C body remains.
Production C is 138,824 physical lines (minus 341), 153 files, zero reference C.
First native fixed-position costs were 1,803/1,847 ns per actual-move/move-path
call, versus C 2,686/3,020; these are VM microbenchmarks, not gameplay speed claims.

Review checked the detailed-path success contract in unit_ai_path.go and
server/object_ai_path.go: status zero appends at least the target point. The
follower reloads the returned count before overwriting its last point. Cached
target writes retain C ordering even for an aliased point. Accumulated tests
and production builds pass in default/server/highres; all binaries identify as
ELF32/80386 with GO386=sse2. The full suite matches exactly 1,553 known failure
entries (15 passing/3 failing/32 skipped-no-test packages), none added or removed.
Fresh ai-path-execution-port gameplay exits 0 against both preserved screenshots,
overrides off. Qualification is complete.

Next: group the six combat actions, their decision/stack helpers, scan callbacks
and private buff cleanup. Record script/audio/strike/projectile effects before
porting; retain broader damage, player-attack and spell engines.
