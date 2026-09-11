# Float-to-integer helpers — 2026-09-11

Scope under investigation: 419A70/419A90/419AB0. All have remaining C callers;
any replacement must keep their declared C ABI (int, short, short). The compiled
386 C bodies save the x87 control word, force truncate mode, perform a 32-bit
FISTP, and restore the control word. The absolute helper applies FABS first.
Short callers consume low16 with normal signed/unsigned C promotion; unspecified
upper EAX bits are not part of the declared short ABI.

## Original-C baseline

4,193,481 ABI conversions pass over 1,397,827 raw float inputs. Coverage includes
every sign/exponent with representative mantissas, every integer from -65,536 to
65,536 and its adjacent float32 values, and one million deterministic raw patterns.
The independent expected result decodes IEEE754 magnitude with integer arithmetic.
It confirms INT32_MIN for masked invalid conversions (NaN, infinity, overflow),
truncation toward zero in range, low16 wrapping, and absolute-before-conversion.

Another 684 conversions pass across all 12 valid PC/RC combinations with exact
control-word preservation, on a locked OS thread with CW restored afterward.
Exception masks retain their normal hosted settings. Sticky exception status
is not tested; no repository reader of x87 exception status was found.
The fixture calls live declarations with ordinary C short-to-int promotion and
checks output guards. It contains no copied conversion algorithm.

Because these helpers have broad C call sites, benchmark the actual C-call
boundary before adopting Go exports. The C baseline is 3.472/3.611/3.791 ns per
primary conversion over three one-second runs. Driver sums are validated; the
loop is in C so it measures the production caller boundary, not repeated Go-to-C
entry overhead. This microbenchmark is not a whole-game performance measurement.

Production C baseline: **140,455 physical lines**, 153 files, zero reference C.
Artifacts: build/port-float-int. Native correctness and boundary-cost comparison
are next; keep the port order open until that evidence is available.

## Measured port-order change

Baseline commit: `eb23d4b2`. The Go-export experiment passes all conversion and
CW checks, but repeated C-caller benchmarks measure 222.9/245.1/285.1 ns/op
(first run 237.9), roughly 60–80× the original C microbenchmark. This does not
establish a whole-game slowdown, but adding that crossing at hundreds of
remaining C call sites is avoidable. Do not adopt those exports yet.

The [tested export experiment](proposals/float-int-go-exports.patch) is saved for
review/reuse. It is not applied to production. Keep the tiny production C
converters for C owners; use the tested native conversion helper when porting
Go owners, starting with grid lookup 411160. Retire the C implementations when
their remaining callers have moved. This is production dependency retention,
not keeping C solely as a test oracle. Production C remains **140,455 lines**.

The next grid audit found `i-1 <= 0` can overflow for INT32_MIN, the observed
invalid-conversion result. Reproduce the caller bug and prepare a bounds repair
before asking the user about the behavior change.
