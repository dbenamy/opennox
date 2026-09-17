# World motion, projectile tracing and timed world objects

## Scope and status

**Corrected-C baseline is frozen and fully qualified. Go conversion is next.**
The sections below retain prerequisite findings; final validation is recorded in
[Qualified corrected-C baseline](#qualified-corrected-c-baseline).

Parent collision-core conversion **950f8f22** is qualified and pushed. This batch
selects **30 live functions / 1,072 C body lines**: sentry registration, beams and
reporting; timed decay queues; activation and velocity updates; radial scans;
fall and projectile contact handling; scorch effects; waypoint movers; shooting
traps and contact triggers. See world-motion-scope.json.

A one-line trace-result getter, sub_537760, had no source caller, callback or
registration: its only references were its definition and header. It is removed
rather than translated solely for tests. Other readers of the underlying trace
state remain intact. The live scope includes the unmarked Hit-reset trampoline,
whose sole caller is Go.

## Corrected-C prerequisite under qualification

The independent real-list contract reproduces a sentry removal defect: after
registering three objects, removing one leaves all three in the active list.
The C function tests an unsigned word with `< 0`, so its unlink branch is never
entered. The compiled original function contains only the membership-bit clear.
Change that condition to test bit 31, matching registration, and preserve the
existing head/previous/next unlinking operations.

This is a narrow, reversible correction under the standing policy. Qualify it
before freezing: six registration permutations, head/middle/tail removal, direct
and destroyed-update paths, repeated removal, backward links and reset. The first
uncorrected run fails as expected (build/port-world-motion/list-1.log). Corrected
fixtures and timed decay contracts are in progress; no goldens are frozen yet.
Because production C changes, this batch requires fresh corrected-C production
qualification rather than reuse of the parent's evidence.

## Planned coverage

Reuse the actual collision/object/index/wall/type/waypoint/RNG/audio/damage owners.
Capture queue membership, flags and callback order, raw float words, actual effect
objects and messages, timer/frame wrap and allocation outcomes. Timed decay checks
independently assert ascending unsigned deadlines, FIFO ties, rescheduling, held
items, expiry and reset. Preserve the existing absolute-deadline wrap semantics.

Cover force/velocity and height flag paths; projectile wall trace outputs and type
exclusions; scorch selection/allocation; sentry beam clipping and contacts; waypoint
transition/wait/disable/missing-target paths; shooting cooldown and actual projectile
creation; trigger size/weight/admission and script events. Review compiled-C
arithmetic stores/reloads and callback-sensitive pointer reloads before translation.
The preceding collision corpus and actual dependent callers will supply broader
three-target coverage. Fresh gameplay/save-load/flat checks follow conversion too.

## Baseline work in progress

physics-4 passes 12 groups after the sentry membership correction. Contracts now
cover list registration/removal, decay ordering and expiry, projectile integration,
activation wrappers, fall/bounce audio thresholds, shaft return, live wall/velocity
paths, beam clipping, strict contact eligibility, actual damage attribution,
viewport reporting and packet rounding/capacity. Captures remain unfrozen while
coverage is extended. A sentry without a player owner attributes damage to itself,
matching the existing FindOwnerChainPlayer contract; the fixture was corrected
after initially expecting nil. Projectile collider coverage explicitly uses the
center shape enum rather than an unspecified shape value.

Disk cleanup removed 14.96 GiB of rebuildable Go cache entries last touched more
than six hours earlier. Recent cache entries, module downloads, assets, logs,
captures and qualified binaries were retained; 24 GiB remained free immediately
afterward. Local audit: build/cleanup/cache-20260917.json.

## Trigger prerequisite correction

The actual-VM one-shot script accepts its first trigger contact and sets its
own disabled flag. On the next contact ScriptCallbackRaw returns nil. The C
trigger unconditionally dereferenced that result, reproducing a null-address
segmentation fault in trigger-disabled-c.log. A nil-result check now treats
that invocation as no admission and leaves existing contact state intact. The
regression also covers a script disabled before its first contact. This is a
reversible correction under the standing policy; fresh corrected-C production
and all target qualification remain required before freezing.

Preserve the C allow-team signed-byte comparison: values 128–255 cannot match
an unsigned object team byte. The deny-team comparison uses signed bytes on
both sides and can reject those values. Contracts cover 127, 128 and 255. This
quirk is captured, not changed.

## Recovery checkpoint validation

The partial corrected-C corpus passes twice: **18 roots, 17 captures / 8,122
records**, no skips. physics-7.log passes in 0.278s. The repeat driver verifies
all 17 hashes against that separate-process run and checks unchanged source.
Mapped-state preflight passes. The recovery manifest records these results;
this is not the frozen baseline or full target/production qualification.

Waypoint mover coverage includes missing/deactivated targets, lazy ID resolution,
pause/resume, waypoint link selection and RNG; scorch coverage uses actual
allocation, pending objects and decay registration with missing types, invalid
sizes, game modes and frame wrap. The trigger regression executes actual script
bytecode and proves that the second one-shot contact no longer crashes.

Shooting traps, radial scan and the remaining spring/tile/monster velocity
integration cases precede the full baseline freeze. Current C count is 41,887
physical lines in 74 files, zero reference C: −3 orphan lines and +4 for the
trigger correction relative to the qualified collision-core parent.

## Completed fixture coverage before freeze

The expanded focused run physics-11 passes all 25 root groups. Real radial
indexing, trap enemy/facing/vision admission, actual projectile allocation and
arrow FX now join the earlier corpus. Projectile allocations are captured and
released between shots so the fixed test object pool remains reusable. The
velocity integration cases exercise real springs (including activation of the
second endpoint), actual tile-grid lookup/Hit sentinel, monster action 67,
freeze/buff force suppression and movement synchronization thresholds.

Two fixture prerequisites were exposed by the positive trap contract: the shipped
576-byte aiming table at 0x587000:202504 must be supplied, and target objects must
have a registered nonzero type. Type zero aliases absent fish IDs in the real
enemy classifier. The owner now reuses PortTestCombatTables and a dedicated
MotionTestTarget type. Shared collision ownership already restores direction
scratch words. No production aiming or enemy-classification behavior changed.

Mover contracts additionally check axis/coincident destinations returning FLT_MIN
(0x00800000), and signed 32-bit speed initialization boundaries. The affected
selection contains the preceding collision-core corpus plus all world-motion,
AI-lifecycle caller and tile-worklist shared-owner contracts (690 root names;
target build tags can omit a root). Full qualification is still pending.

## Qualified corrected-C baseline

All 25 focused roots and 24 captures / 9,633 records repeat unchanged. The affected
sweeps pass 690/689/690 roots (43,124/43,123/43,124 including subtests), zero skips,
and 245 frozen captures each in default/server/highres. Durations are
323.62/402.75/334.67s. Fresh production passes in 380.87s: three 386/SSE2/CGO
binaries, retained/retired symbols and no test helpers, exact known 1,553 failure
entries (15 pass / 3 fail / 32 no-test packages), gameplay/save-load and flat-map
regeneration. All five gates share a 2,236-file source manifest. No source job is
active. world-motion-qualification.json records the evidence.

The initial Go drafts are uninstalled. Review of the newly built C assembly shows
that projectile integration retains both velocity intermediates through position
addition; velocity force sums and radial/mover deltas also remain wide. Declared
float locals alone do not identify rounding boundaries. The draft is being
reviewed against actual stores/reloads before integration; frozen expectations
remain unchanged. The reference audit identifies six needed world-motion C
exports and nine collision-core exports that become unnecessary after conversion.
