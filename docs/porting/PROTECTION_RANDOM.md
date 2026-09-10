# Protection floating RNG and state — 2026-09-10

The remaining four C helpers (56F240, 56FE30, 56FF00, 56FF80) are now native Go.
A private protection.Random owns the five doubles and span/min/max words;
startup and rekey use it directly. No production C caller or external state
reader/writer remained, so all four bridges and both standalone C range globals
are removed. Historical address catalog entries remain metadata; production no
longer accesses those former RNG blob fields. The five constants are the exact
shipped binary64 values, now declared with the pure helper.

The port preserves the literal C behavior: floor(v) had its result discarded,
and the state receives v-v. Seeded finite draws therefore yield positive zero;
the protection range call returns one. This conversion deliberately retains that
behavior and its evolving carry/state. Seed zero maps to UINT_MAX, five xorshift
words initialize the state, and 19 warmups run. Range arithmetic/comparison uses
uint32 wrapping. Explicit float64 conversions enclose each product and sum to
preserve the hosted x87 PC53 evaluation order and prevent fused operations.

Only seeded finite state is reachable through production initialization/rekey;
there is no external state input or writer. Artificial huge/nonfinite state can
expose x87's wider exponent range, and is outside this helper's compatibility
contract. Production does not need arbitrary-precision arithmetic to reproduce
its bounded seeded states, including eventual binary64 underflow.

Original-C baseline: `877401b5`. The pure Go helper existed there but was not
connected to production. 506 seed scenarios, each with an initial snapshot and
90 operations, compared **46,046 snapshots** byte for byte against all four C
entries using shipped constants. Cases include zero/signed-limit/random seeds,
90-call traces, reseeding after state changes, equal/reversed/full-width ranges,
unsigned wrap and eventual underflow. Return values and all range words are
checked alongside the five double bit patterns.

Independent pure tests use math/big precision-53 nearest-even arithmetic, widened
xorshift masking, and explicit binary64 state stores for 506 seeds and 80 draws
each. They run on 386 and amd64. Manager fixtures now snapshot owned Go state;
they consequently use the shipped coefficients rather than uninitialized zero
blob constants. Original C arithmetic remains recoverable in Git, without a
retained C implementation for testing.

Production C: **141,844 physical lines (−70)** in 153 files; test-reference C: **0**.
Local artifacts: `build/port-rng/`.

All accumulated protection tests pass on 386 default/server/highres. All three
production binaries build, and symbol checks confirm the retired four RNG
functions and two globals are absent. The `rng-port` warrior scenario accepts
both preserved screenshots with overrides disabled. Full-suite results exactly
match the object milestone: 15 passing, 3 known failing, 32 skipped/no-test
packages; all 1,553 package/test failure entries are unchanged.
