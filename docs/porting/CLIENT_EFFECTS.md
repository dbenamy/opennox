# Client effects — next connected batch

Candidate scope: 42 routines /2,177 physical C lines. Five whole files:
client__draw__fx, drawrays, lightning, plasma and glowdraw. Include GAME3_1's
plasma setup (4BA670), drawable update (4CA720), and four curve helpers (4BE800
through before 4BF010), plus GAME3's glow helper (4B6B80). Ranges, original
sections and symbol/callback audits are stored here. No effects source or
fixture is applied yet; quest eligibility qualification is complete.

The curve helpers belong in the batch: plasma passes sub_4BA8B0 as a synchronous
callback. A typed Go callback can remove that internal C round trip when both
caller and curve evaluator move. The curve evaluator still has a separate C
caller in GAME2.c, so retain a C wrapper that adapts its external callback.
Keep the drawable's 4CA720 update callback C-callable while the drawable stores
and dispatches C function pointers. Twenty-one draw callbacks have Go address
references and must keep their production C entry points. Recheck all callers
before deciding the final retained/retired set.

Testing approach:
- Use a real NoxRender with explicit RenderData and a bounded RGB5551 pixel
  buffer; capture exact framebuffer words via the existing PixelHash16 oracle.
  Retain renderer state, scratch arrays, input words, drawable state and both RNG
  indices. Save/restore every named C global and touched memmap region.
- Reuse the guarded server/objective owner for frame/RNG state. Add a small
  client fixture owner using production drawable lists/pool/lookup machinery,
  explicit synthetic thing definitions and the real renderer. The root test
  adapter can use the actual root Client factory/deletion methods. Control
  allocation failures as an external dependency; exercise real positive paths.
- First prove one asset-independent pixel-rendering probe through original C,
  then add geometry/clip, color/mode, lifetime/frame-wrap, packet-to-effect,
  allocation/failure and stateful sequences. Include curve callbacks and packed
  signed coordinates. Check float intermediate rounding as in earlier math ports.
- The root package's NewNoxRender wraps client/noxrender.NoxRender and satisfies
  legacy.Render2. A client-side tagged constructor can initialize just the
  required registries/pools/viewport without starting a display or registering
  unrelated server callbacks. Preserve production helpers rather than mocking
  their list/state behavior.

Color compatibility decision (2026-09-14): these higher-level effects preserve
current C output through the current renderer. This does not change the sprite
color conversion backend diagnosed in FAILURE_DIAGNOSIS.md. Establish explicit
current-C raw-framebuffer references; do not regenerate historical sprite PNG
hashes. Export representative diagnostic PNGs only for visual review, separate
from raw framebuffer comparisons.

The fixture adds client/renderer ownership, so broaden qualification beyond the
previous tag-independent eligibility batch. Decide exact variant scope after
the first C probe establishes which paths run under standard/server/highres;
keep all production builds and fresh gameplay, and use the guarded test runner.
