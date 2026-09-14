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

Current status: repeated C references locked; accumulated standard C baseline
passed in 372.811s: 659 selected/completed tests, 658 passed and one optional
prerequisite probe skipped. Native implementation is staged locally and has not
replaced C yet.
Raw evidence and explicit scope/ABI metadata are in `build/port-client-updates`.
Recover C from the baseline commit accompanying this document; hashes and tests
are tracked, so local capture files are not required for normal verification.
