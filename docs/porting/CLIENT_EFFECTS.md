# Client effects

The connected batch contains **45 routines /2,263 physical C lines**: five whole
files (client__draw__fx, drawrays, lightning, plasma, glowdraw), GAME3_1 plasma
setup 4BA670, drawable update 4CA720, four curve helpers 4BE800–4BEDE0, and
GAME3's spark/glow helpers 4B6880, 4B6970, 4B69F0 and 4B6B80. All production
effects remain C. Production count: **98,506 lines /147 files /zero reference C**.

The provisional ABI audit retains 32 entry points and retires 13. Twenty-one
draw callbacks have Go address references. Keep 4CA720 C-callable while drawable
updates store and dispatch C function pointers. Plasma's synchronous 4BA8B0
callback and curve evaluator can use a private typed Go callback after both
move. The evaluator 4BEDE0 still has a separate GAME2.c caller and retains its
external C callback adapter. The three additional spark helpers have only
in-batch callers and belong with their color-specific draw callbacks.

## Original-C baseline in progress

The asset-independent fixture uses production drawable allocation, type lookup,
list/index ownership, construction and deletion, plus the real NoxRender,
RenderData, viewport and RGB5551 framebuffer. It starts no display and installs
no unrelated server or GUI callbacks. Cleanup releases only owned resources;
it neither closes an unopened sprite bag nor closes shared fonts.

The initial white-orb probe checks nonempty pixels and animation state. Its
locked framebuffer matches repeated C execution and passes standard/server/
highres through the guarded runner (one selected test each; 164.684s,
158.901s and 21.149s wall). Those runs cover only this probe.

Two larger capture groups are now repeated and locked:

| Group | Coverage | Recorded results |
| --- | --- | --- |
| Orb lifetime | 1,056 combinations of radius, period, counter and clip boundary; up to six sequential draws, stopping on deletion | 6,070 snapshots |
| Curve segments | Seven control-point sets, eleven step counts and six finite shifts; signed and large coordinates, empty step counts | 462 cases |

The orb model independently checks exact clip equality, byte-counter wrap,
shrinking/deletion, return values and list/count ownership. Every live drawable
byte must stay unchanged except the expected size/counter fields. Captures
include the complete 512-byte C drawable prefix, complete RenderData and raw
framebuffer hashes. No process pointer occurs in this single-drawable prefix;
the Go-only client handle is excluded.

The curve fixture records exact ordered callback endpoints and userdata, guards
all input storage and requires unchanged control points. Independent checks
cover continuity, constant curves, binary-step Hermite endpoints and a symmetric
midpoint. The first endpoint assertion caught missing fixture initialization:
normal startup loads the Hermite coefficients. The fixture now copies the
production 64-byte table from embedded data and restores the previous contents.
No production algorithm was changed or copied for testing.

The locked hashes live in the root test files. Raw captures, original sections,
ABI audit and qualification metadata are local under build/port-client-effects.
The focused standard check passes all three selected/executed/completed tests
(89.555s wall; early-baseline.json). This is an early
fixture checkpoint, not the complete 45-routine baseline or conversion.

## Remaining work and qualification

Expand coverage to moving orbs, colored sparks, bounce/lifetime/frame wrap,
lightning recursion and packed signed coordinates, plasma state and particles,
curve rasterization, particle construction/failure and list ownership, and ray
packet dispatch/cache state. Capture both RNG indices, framebuffer words,
renderer state, scratch arrays and normalized drawable/list state. Preserve
named C globals separately from mapped words. Full production dispatch and state-environment adapters plus bounded table
initializers are staged locally but are not yet applied. Inspect them before
use; the ray pointer cache needs separate normalization and capture.

Use independent assertions before locking repeated original-C hashes. Complete
the accumulated standard baseline before conversion. The fixture adds shared
client/renderer ownership, so run the full accumulated native matrix under all
three tags, all production builds and a fresh unchanged gameplay replay. All
selected suites use tools/porting/run_tests.py to verify actual execution.
Broaden earlier for unexplained failures or changed shared behavior. Do not
edit source during an active build/test job.

Color compatibility decision (2026-09-14): preserve current original-C effects
output through the current renderer using exact raw framebuffer words. The
historical sprite-color policy remains as documented in FAILURE_DIAGNOSIS.md;
do not regenerate those goldens in this batch. Representative diagnostic PNGs
may support visual review but do not replace the raw-word comparison.
