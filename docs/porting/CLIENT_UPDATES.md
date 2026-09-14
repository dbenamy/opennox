# Client drawable updates

The next conversion covers 27 routines: twelve whole drawable-update C files,
plus their fireball/charm/cloud/death-ball callers and height/cloud callbacks
in GAME3_1/GAME3_2. Original scope: 1,069 physical C lines. The effects conversion
is pushed at `a2a6e2aa` and supplies the particle primitives for this batch.

Two prerequisite allocation guards add five C lines. Magic-missile tail failure
returns without advancing its previous anchor. Magic-trail tail failure skips
only the tail, preserving the anchor and continuing its four independent sparks.
Previously both dereferenced a failed allocation. This is an intentional,
reversible behavior repair under standing authorization; review the retry policy
if the engine later adopts a different allocation-failure policy.

The focused regression forces tail failure, checks unchanged anchor/no decay
entry, observes surviving independent sparks, restores allocation, verifies the
old anchor is retried and receives the correct deadline, then checks a stationary
anchor does not emit another tail. Production C is temporarily **96,221 lines
/142 files /zero reference C** before conversion.

The shared effects owner now optionally adds nine synthetic drawable types and
normalizes three known draw/update callbacks. Existing type IDs and all existing
effects references remain unchanged. Tests use real drawable pools, ownership
lists, spatial indexes, RNG streams, balance access and the stored cloud callback.
No asset-backed sprite renderer is replaced by this fixture.

Repeated repaired-C captures agree byte-for-byte:

| Group | Results | SHA-256 |
| --- | ---: | --- |
| Construction /20 operations | 1,920 | `eeb05a51bc4138c6511ae3c1bd67d8bd8c2dad802d151f04f187dced6cfe023a` |
| Height/gravity/bounce catch-up | 2,688 | `bc7c837bf2ca7b2865bd1209e0c01105381fecdd07cf8096ef15d3c002e1a03c` |
| Actual stored cloud callback | 120 | `c9fe11b89e517a17837e6815fb975359a31089fec7a9618359fe891d4ebd0704` |
| Charm/heal/drain transfers | 1,536 | `45569e43638732ef34eac69611e6385dccf6f4b3e09ee0e075d679ab66358d05` |

Independent checks cover emission and RNG counts, allocation/ownership, return
values, endpoint targets, missing bindings, pause/emission gates, frame wrap,
height and bounce arithmetic, cloud height wrap, and failed-tail retry. Cases
include coordinate/static/dynamic endpoints, upper code bits, viewport offsets,
clipping, negative/zero counts, density, height, charge age and balance radius.
Total: **6,264 captured results /four groups**, plus the retry contract.

Qualification plan: accumulated standard matrix at baseline and native boundary;
affected effects/update tests under server/highres tags; all three production
builds with ELF32/SSE2/ABI checks; fresh unchanged headless warrior gameplay.
The complete asset-backed suite matched known failures at the immediately prior
effects milestone. Repeat it if shared regressions or another subsystem change
justify doing so, rather than for this state-only callback batch by default.

The native conversion is fully qualified. Baseline `bbac104c` is pushed; recover
original C from that checkpoint. All four locked groups match exactly on the
first native comparison (161.968s), and both tail-allocation and exact squared-
distance 200/205 contracts pass. The latter is an independent source-derived
boundary check added during native review; it is not another C capture group.

Final validation:

| Check | Outcome |
| --- | --- |
| Accumulated standard | 660 selected/completed; 659 pass, one optional prerequisite skip; 448.814s |
| Affected effects/updates server | 27 selected/completed/pass; 161.608s |
| Affected effects/updates highres | 27 selected/completed/pass; 25.727s |
| Production builds | Standard, server and highres pass |
| Binary checks | ELF32/i386, SSE2, CGO; four retired symbols absent, 23 required exports present; test helpers absent |
| Fresh headless gameplay | Unchanged warrior scenario, original assets, override=false, Xvfb/null audio; exit zero in 37.031s |

Twelve C files and the connected GAME3_1/GAME3_2 sections are removed. Four
private C entry points retire; 23 live callbacks/exports remain and are exercised
through their C ABI. Native callers invoke the Go helpers directly. The existing
energy-spark Go wrapper also avoids a C round trip, preserving int16 height
narrowing. No C algorithms remain solely for these tests.

Compatibility details preserved: charm/heal unsigned clipping versus drain's
signed clipping, packed coordinate/static/dynamic bindings, parent-height clearing
by vortex creation, lazy named versus mapped cache storage, exact allocation/RNG
ordering, real cloud callback/list/deadline ownership, and float catch-up/bounce
rounding. The missile density entry gate is unsigned and its loop termination
signed; the native helper preserves that distinction.

Production C: **95,146 physical lines /130 files /zero reference C**. Reduction:
**1,075** from the repaired baseline (1,074 scoped C lines and one unused trailing
separator). Net reduction from the prior effects conversion is 1,070 after the
five prerequisite guard lines. Accumulated coverage is **115,750 captured results
/1,061 groups**, plus independent contracts.

Raw evidence, explicit scope and ABI metadata are in `build/port-client-updates`.
Tracked tests and hashes are sufficient for normal verification after VM recovery.
Next candidate: remaining procedural particle drawing and shared lighting/palette
helpers, with the same real renderer and drawable owner. A staged 18-routine,
671-C-line audit and fixture draft are in `build/port-client-draw-particles`.
