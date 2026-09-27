# Allocator fixture thread scope

The porttest allocator observer intercepted process-wide `calloc`, `free` and
`time` while a fixture was active. That scope also included Go runtime startup
threads: `runtime/cgo/gcc_linux.c` frees its `ThreadStart` before `crosscall1`
initializes the new thread's Go state. Calling a Go observer from that early free
is unsafe. The local runtime source and the test executable's disassembly agree
on this order.

The observer's active flag and fixed clock are now thread-local. Its Go control
helper pins the fixture goroutine to its OS thread while observation is enabled.
Repeated enables/disables retain the existing idempotent cleanup, and the pin is
balanced independently of an enclosing fixture's pin. The fixtures already use
sequential global state; this change does not make them parallel-safe.
Production allocation code is unchanged.

An independent contract starts and joins a real pthread, reading observer state
there without a Go callback. It checks state before, on the fixture thread, on
the foreign thread, and after teardown, including repeated activation/cleanup.
The old fixture fails with `[0 1 1 0]`; the repaired fixture yields `[0 1 0 0]`.

A server test process stalled during GC coordination with `themeTestReleased`
on its stack. This prompted the audit, but the exact cause of that stall remains
unproven. Passing isolated and complete replays did not establish that the
process-wide observer was safe; the independent thread-scope failure did.

Qualification uses a separate source copy of original callback baseline
`eb4c6e21`, differing only in the observer fixture and its new contract. The
expanded focused selection includes all identified observer owner families and
passed twice in separate processes: 387 default, 385 server and 387
high-resolution roots, with no skips or failures.
The full default baseline recorded before this fixture repair is historical
evidence, not a same-source run for the repaired fixture. The callback conversion
requires a fresh full default sweep and its remaining production gates.

Machine-readable evidence: [theme-observer-scope-qualification.json](theme-observer-scope-qualification.json).
Local failure, source copies, binaries and logs:
`build/port-final-callback-exports/{observer-scope-before.*,original-fixed/}`.
