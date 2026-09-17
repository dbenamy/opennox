# World motion, projectile tracing and timed world objects

## Result and scope

**The Go conversion is fully qualified.** Corrected-C baseline e7174c35 is pushed;
it follows collision-core parent 950f8f22 and recovery checkpoint 71888d5e.
Thirty live functions / 1,072 C body lines move to Go: sentry lists, beams and
reporting; timed decay; activation and velocity updates; radial scans; fall and
projectile contacts; scorch effects; waypoint movers; shooting traps and triggers.
See world-motion-scope.json for the corrected-C definitions.

Production C is **40,777 physical lines / 74 files / zero reference C**, down
**1,110 lines** from the corrected-C baseline. Six exports remain for actual C
callers and registered callbacks. Twenty-four obsolete batch interfaces, nine
collision-core interfaces and two C globals are retired. Decay/sentry heads and
velocity type IDs are Go-owned. The last collision-core export, sub_547DB0, still
serves the remaining C spatial candidate helper. No converted C algorithm is kept
solely to supply test expectations.

## Corrected-C prerequisites

Sentry removal tested an unsigned flag word with `<0`, so the compiler removed
its unlink branch. Real three-member lists reproduced retained head/middle/tail
membership. Testing bit 31 restores the intended unlinking; contracts cover all
registration permutations, repeated removal and destroyed-object updates.

An actual one-shot trigger script accepts its first contact and disables itself.
The next call returns nil, which C unconditionally dereferenced. The original
crash is recorded in trigger-disabled-c.log. Nil now means no admission, preserving
prior contact state. Regression cases include a script disabled before its first
contact. Both corrections were qualified in C before freezing, under the standing
policy for confident reversible fixes. See DECISIONS.md for later review.

The allow-team signed-byte quirk is preserved: values 128–255 cannot match an
unsigned object team byte. Deny-team compares signed bytes on both sides. Tests
cover 127, 128 and 255. Proven orphan sub_537760 was removed; other trace readers
remain intact. Prerequisites changed the parent's C count by +1 overall.

## Contracts and ownership

The focused corpus contains **25 root groups, 24 captures / 9,633 records**,
including a separate actual-VM one-shot regression. It uses actual object pools,
spatial indices, walls, waypoints, RNG, audio, message queues and script execution.
Independent checks supplement exact captures for list order/link integrity,
unsigned deadline ordering/FIFO ties/wrap, callback order, force and movement,
strict geometry boundaries, admission and allocation effects.

Coverage includes sentry clipping/damage/viewport reports and packet capacity;
projectile integration, sampling, wall normals and bilateral contacts; falling,
bouncing and shaft return; real springs and endpoint activation; actual tile
lookup and floor Hit sentinel; monster action 67, buffs/freezing and sync
thresholds; waypoint transitions, RNG, FLT_MIN axis/coincident velocities and
signed speed boundaries; scorch allocation/decay; trap target admission,
projectiles, timers and FX; and trigger mass/class/team/script gates.

Positive trap admission exposed two fixture prerequisites: supply the shipped
576-byte aiming table at 0x587000:202504, and use a registered nonzero target type.
Type zero aliases absent fish IDs in the real enemy classifier. Fixtures reuse
PortTestCombatTables and shared direction-scratch ownership. Projectile records
are captured and released between shots to reuse the fixed test object pool.
These were fixture changes, not production behavior changes.

## Qualification

| Gate | Default | Server | High-resolution |
| --- | ---: | ---: | ---: |
| Root groups | 690 | 689 | 690 |
| Including subtests | 43,124 | 43,123 | 43,124 |
| Frozen captures | 245 | 245 | 245 |
| Captured records | 114,004 | 114,004 | 114,004 |
| Native duration | 278.76s | 386.65s | 330.78s |

All captures match corrected C, with zero skips. The selection includes the
qualified collision corpus, affected callers/owners, AI lifecycle and shared
tile-worklist contracts; see world-motion-selection.json. Native-3 focused
contracts pass in 0.373s. Static mapped-memory checks pass.

Fresh native production passes in **372.74s**: three 386/SSE2/CGO binaries, ABI and
retained/retired symbols, no test helpers, exact known **1,553 failure entries**
(15 passing / 3 failing / 32 no-test packages), gameplay, save/load and flat-map
regeneration. All four native gates share an unchanged **2,244-file manifest**.
Evidence: world-motion-native-qualification.json and ignored native-* artifacts.
The earlier corrected-C gates share a 2,236-file manifest; their evidence remains
in world-motion-qualification.json and c-* artifacts.

## Implementation findings and remaining dependencies

Compiled-C review places float32 rounding at actual stores/reloads. Projectile
velocity intermediates, force sums and radial/mover/fall deltas remain wide;
mover/fall denominators, velocity's Y sync threshold, and trace X before step
division have explicit spills. Declared C float locals alone were insufficient.
Frozen raw-float records remain unchanged.

Validation caught and corrected two translation errors: delayed deletion belongs
to the outer server wrapper, and alloc.New selects a type without initializing it
from its argument. Explicit point initialization restored trace/dispatch contracts.
The temporary C-backed trace record remains necessary for the retained spatial
callback; the next connected spatial-targeting batch can remove that allocation
and the final collision-core export.

The cache cleanup reclaimed 14.96 GiB while retaining original assets/archive,
module downloads, logs, captures and qualified binaries. About 18 GiB remains
free after this qualification. The local audit is build/cleanup/cache-20260917.json.
