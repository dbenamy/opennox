# GPT-6 Luna implementation trial

The requested `gpt-6-luna` helper produced an isolated Go draft for wide/narrow
ASCII-insensitive comparison, base-10 signed 32-bit parsing and UTF-16 copying.
The primary supplied the task boundaries and existing C contracts, independently
reviewed the implementation and callers, and ran the comparison harness.

All four existing scalar roots pass: 393,216 wide comparisons, 131,072 narrow
comparisons, 65,567 decimal cases and 65,536 copies — **655,391 captured cases**.
Complete JSON results match the existing original-C probe3 captures byte for byte.
The independent assertions also cover terminators, saturation, copy guards and
return-pointer identity. Tests ran on the qualified 386/SSE2 target with the
baseline Go environment. The harness only adapts allocation, direct calls and
capture comparison; it does not regenerate expected results.

No behavior corrections were required. The primary requested replacing repeated
uintptr arithmetic with unsafe.Add and formatted the draft. Source review checked
ASCII-only folding, wide result magnitude, narrow sign clamping, whitespace/sign
handling, overflow saturation and copying through NUL. Valid terminated pointers
remain the caller contract. The parser intentionally covers only the live base-10,
no-end-pointer API. Locale and C-long assumptions match the current target.

The first sandboxed test launch failed with `bad system call` before any test ran.
An approved execution outside the sandbox passed all four roots. This is an
execution-environment issue, not a draft behavior mismatch.

Local recovery artifacts are under `build/port-final-formatting/luna-trial/`:
`text_scalar.go`, `review.md`, `prepare_test.py`, `text_scalar_test.go`, and
`test-unsandboxed.log`. Draft SHA-256:
`d7d8ec0614aa01bee9cd09e2cb2b2a7cff219e6d2eca468e0aaad89a7e00bedf`.
These ignored artifacts are disposable; the original tests remain checked in.

Use Luna for further bounded implementation tasks, one helper at a time. Keep
baseline design, review, integration and full qualification with the primary.
This small successful trial does not measure subscription savings or establish
quality for larger stateful subsystems. It is not production integration: the
formatting batch still needs baseline repeats/freeze and final qualification.
Production is unchanged; C remains **1,219 physical lines in 12 files**.
