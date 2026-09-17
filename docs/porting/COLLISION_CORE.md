# Collision event queues, activation and contact dispatch

## Scope and status

Parent **011fcb73** is qualified and pushed. This batch selects **23 reachable C
functions / 993 body lines**, in GAME5.c and GAME4_3.c. All are in the whole-source
live closure. Baseline production is **42,982 C lines / 74 files / zero reference**; native
production is **41,886 / 74 files / zero reference (−1,096)**.
The unchanged-C baseline and native conversion are qualified below.

Scope: Hit class allocation/reset, bucket/pair deduplication and event dispatch;
active collision-object queue membership/removal/pop/traversal; pair eligibility
and retained-target lookup; shape containment and circle/box response; gate/circle
response and angular queue consumption; circle wall scanning and wall-opening
response; elevator and shaft contact dispatch. See collision-core-scope.json.

## Baseline approach

Reuse actual object, wall, Hit and gate owners from the geometry batch. Capture
both event order and bucket links, first-normal retention, callback order and
inverted normals, flags, activation head/tail/links and lazy type caches. Include
real fixed-pool exhaustion and reset/reuse, without replacing allocation logic.
Observe collision/damage callbacks at their real dispatch boundaries.

Supply shipped tables explicitly, including the separate circle/box coincidence
normals at 0x587000:289928, not the circle/circle table used by the previous batch.
Own and restore backing bytes, live globals, object/update records and lists.
Use independent containment/force/list contracts alongside frozen C comparisons.
Exercise team/owner/message gates, angular wrap/clamp/stops, wall opening, height
thresholds and force-suppression paths. Review all actual caller dependencies.

Inspect compiled C before translating float expressions: the prior batch proved
that x87 stores/reloads, rather than local declarations alone, determine observable
rounding. Preserve circle/box's existing interior branch unless a separately
justified correction is made before freezing. Do not substitute textbook collision
rules silently. Activation's retained low-byte returns can contain callback
addresses; use established aligned observation callbacks and identity normalization.

Repeat and freeze the C baseline before conversion. Qualify the selected corpus
on default/server/highres, keeping starts/completions and skips explicit. Reuse
qualified parent production evidence only when actual production source is
unchanged and that identity is recorded. Always run fresh production/gameplay,
save/load and flat-map checks after conversion.

## Qualified C baseline

Sixteen focused roots/captures cover **22,848 records**. Two independent processes
(core-6/core-7) produced identical hashes; expectations are frozen in the tests and
collision-core-batch.json. Static-3 passes. Broader default/server/highres pass
**660/659/660 roots**, **42,918/42,917/42,918 tests including subtests**, zero skips.
All **221 captures / 104,371 records** are identical across targets and match the
manifest. Durations: **308.57/387.42/328.98s**. All three runs share one unchanged
**2,213-file source manifest**; all sessions are joined.

Coverage includes real Hit pool exhaustion, bucket deduplication and retained
normals; activation/FIFO; pair eligibility/retained targets; containment/distance;
callback/damage order; angular drain; gate/circle and circle/box response; indexed
object and wall scanning; wall-opening audio; elevator/shaft heights; special-class
pair dispatch; and the retained radial C caller's strict contact threshold and
PosVec/NewPos distinction. Broader selection adds actual dependent object, player
attack, temporary update, spell effect, generator and spawn-policy fixtures to the
preceding geometry corpus.

The initial height fixture lacked the relocation used by the existing absolute
value helper. It now owns/restores both scratch and relocation, and captures
scratch writes. Production source is unchanged; preserve this legitimate retained
C dependency's side effects during conversion.

The default driver's final hash stage initially had incorrectly renamed inherited
geometry capture filenames. Tests themselves completed successfully. The corrected
manifest matches every saved capture; the independent c-audit verifies full event
completion, unchanged source and all cross-target hashes. The original failed
driver report remains intact. Server/highres used the corrected manifest.

Production evidence is explicitly reused from 011fcb73's native-production:
production-source identity (only test/porttest differences), all three binary
hashes, successful qualification and all 23 selected C symbols are verified.
No fresh production run is claimed for this test-only baseline. Client SHA:
3b0754cf3af539d8447c294af8a64a48429441760495861838d69429233a7737.
Fresh three-binary, gameplay, save/load and flat-map checks follow conversion.

Artifacts: build/port-collision-core/{frozen-focused.json,c-audit.json,
production-reuse-audit.json,c-default,c-server,c-highres}. Scope/body hashes and
frozen test selection/manifest are committed alongside this report.

Before starting this batch, verified duplicate assets in nine completed polygon/
geometry runs were removed, reclaiming **4.638 GiB**. Per-run restoration manifests
retain path, hash and metadata. Changed outputs, original assets and archive are
untouched. The deletion invocation is consumed and must not be rerun.

## Qualified native conversion

Baseline **bc4a791e** is committed/pushed. The installed six-file implementation
replaces all 23 selected functions, retains ten exports for actual C callers and
retires eighteen interfaces (including five geometry exports whose last callers
moved). Eight private C globals, seven mapped type-cache words and the bucket array
move to Go ownership. The Hit pool still uses the real fixed C-backed allocator;
record size and normal offset have compile-time checks. Root golden files remain
byte-identical to the baseline. Focused native-3 passes all 16 roots/captures (22,848 records, 0.378s), unchanged
from C; static native-2 passes. Broader default/server/highres pass **660/659/660 roots**,
**42,918/42,917/42,918 tests including subtests**, zero skips, with **221 identical
captures / 104,371 records** matching C. Durations: **275.41/387.91/328.61s**.
All share an unchanged **2,219-file source manifest**. Fresh production passes in **375.82s**: three binaries/ABI/interface checks, the
exact known 1,553 asset failure entries (15 passing / 3 failing / 32 no-test
packages), gameplay comparison, explicit save/load and flat-map regeneration.
All four gates share the same unchanged 2,219-file source; all sessions are joined. Current C is 41,886 / 74 files / zero reference (−1,096), including
obsolete address-marker and separator cleanup.

The first link check caught removal of a cgo import still needed to emit the
existing update-springs export; restored before retry. Source/interface static
checks pass. The reviewed angular stop uses unsigned modulo, and the retained
absolute-value helper still receives height differences and updates its scratch.
All install scripts/drafts are consumed; never reinstall them.

The final arithmetic review preserves the shaft's wide height subtraction before
its float32 argument store, unsigned angular-stop modulo and original unordered
comparison behavior. Dispatch uses pointer-typed callback arguments for the local
inverse normal, keeping its lifetime visible to cgo. The full interface audit also
checks that every exported Go function still has the cgo import needed to emit it.
No intentional gameplay rule changes or frozen-output changes were made.


Final artifacts: build/port-collision-core/native-{default,server,highres,production},
native-audit.json and native-interface-audit.json. The six production Go files are
collision_core_{queues,dispatch,geometry,gate,heights,exports}.go. Actual Go callers
and state owners use Go directly. Client SHA:
71758f918f3d4ae516cbb3de252bd056f322013095e35da334f087c03753cc33.

Next read-only proposal: world motion, projectile tracing and timed world objects,
31 candidates / 1,069 body lines (30 live plus a one-line potential orphan). Review
and reproduce the sentry-list unsigned membership comparison before freezing that
baseline. No correction to those next-batch functions is included here.
