# Go memory and string helpers

Status: qualified Go replacement. Allocation ownership is unchanged.

Replace `alloc.Memset`, `Memcpy`, `Memcmp`, `Strcpy`, `Strcat` and `Strcmp` with Go
while preserving signatures, logging order/spans, destination-pointer returns,
source preservation and valid-input behavior. Keep allocation/free ownership and
all native client bindings unchanged. Shared allocator ownership migration is
separate work in [INTERNAL_C_GLUE.md](INTERNAL_C_GLUE.md).

## Contracts and baseline

New direct API tests exercise guard bytes, destination/source offsets, zero-byte
operations on valid pointers, unsigned and high-byte data, terminating-NUL copies,
empty strings, comparison ordering, and bounded writes up to 65,535 bytes.
The compare captures include input bytes, lengths/offsets and normalized signs;
there are 87,961 memcmp and 107,161 strcmp cases. Independent expectations check
ordering and complete source/guard preservation. Frozen ordering hashes match in
three separate original-libc processes, captured after restoring the original
helpers and before resuming conversion. The safe bridge's 142-record hash remains
unchanged when comparison returns are normalized; its original cases already
returned only signs. `TestShopStockLoading` exercises the actual
safe-profile strcpy consumer.

Default/server/highres original helper contracts and allocator-class tests pass.
A first sandboxed probe exited SIGSYS before qualification; native runs pass.
Its failure log is retained separately as an environment limitation, not a source
regression. The reused pre-conversion production baseline is `2db93f68`; only tests/docs
changed before conversion, so its production evidence was reused.

Since memory clearing/copying is shared infrastructure, qualification after
conversion includes the accumulated root port corpus in all three profiles,
explicit helper/class tests, original safe-profile bridges, safe/static checks,
fresh client/server binaries, exact known-suite comparison and headless save/load.
Raw artifacts are under `build/port-go-memory/`.

## Review and decisions

Luna drafted the six replacements in an ignored directory; the primary owns
original capture and acceptance. The first exact-return baseline is preserved in commit `5cc27785` and the raw
hashes in the baseline evidence. A Go draft exposed libc-dependent magnitudes:
memcmp differs for one-byte versus larger spans, and strcmp differs by optimized
path and mismatch position. Reproducing these incidental values is unnecessary:
a whole-source caller audit found only safe C adapters (plus test bridges and
tooling), with no engine consumer relying on magnitude. The revised contract
preserves negative/zero/positive ordering, with this deliberate, reversible
compatibility decision recorded rather than silently regenerating goldens.
Original helpers were restored before all three revised baseline captures.
Probe sources/results remain under `build/port-go-memory/`.

Luna's audit initially cited a nonexistent fixture filename and conflated sign
assertions with frozen exact-return capture; primary source review corrected both.
In the reverse review, Luna caught missing high-byte copy inputs and unnecessary
nil-pointer zero-length libc calls in the primary's test draft. The primary
removed the undefined-input cases, expanded string payloads, and included case
inputs in exact-return hashes. No measured delegation savings are claimed.

Original memcpy/strcpy/strcat require valid nonoverlapping regions and sufficient
capacity. Do not promise new overlap or invalid-pointer behavior. Memory logger
calls are retained verbatim and reviewed; enabled logging is not a new fixture
claim. Allocation failure and allocator-domain migration are outside this batch.

The original safe-profile bridge/consumer tests pass without skips. Original
fill/copy measurements cover 14 size/value combinations, three samples each,
with zero steady-state allocations. See [baseline evidence](go-memory-c-qualification.json).
Only the documented isolated map-population diagnostic may skip in the broad
corpus; the manifest rejects any other skipped root.

## Performance qualification

The initial byte-loop nonzero fill regressed substantially and was replaced by
progressively doubling an initialized byte span with Go copy. All bounded
contracts pass. Final bounded benchmark samples are retained for native
qualification; these are VM microbenchmarks, not whole-game frame measurements.

Plain Go copy remains slower than libc at 4 KiB and 64 KiB on this 386 runtime.
Whole-source owner review finds the two production Go memcpy callers both copy
60-byte AI path records; the direct memset caller clears 32 bytes. Both sizes are
faster without cgo. The remaining wrappers are safe-profile adapters. Keep the
standard Go implementation and record the large-copy limitation rather than add
a custom architecture-specific memory-copy backend with no identified large
production consumer. This is reversible if a real workload establishes a need.

Median measurements (three 100 ms samples per combination) put the actual
60-byte copy at 11.4 ns versus 110.3 ns original and the 32-byte zero fill at
8.37 ns versus 85.84 ns. Large copy medians are 3.42× original at 4 KiB and
3.86× at 64 KiB; 1 KiB/4 KiB nonzero fills are 1.21×/1.26× original, while
64 KiB nonzero fill improves from 1842 to 1486 ns. All steady-state operations
allocate zero bytes. VM noise limits small differences; no whole-game speedup
is claimed. The rejected byte-loop and final measurements are retained.

The safe build still checks retained C bridges, but its C compiler sanitizer flags
do not automatically instrument the new Go copy/clear operations. Guard-byte and
source-preservation contracts check the exercised write/read boundaries; the
unchanged memory-log calls remain in place. A passing safe build is not a claim
of whole-Go AddressSanitizer coverage.

The accumulated sweeps complete 1,280 default, 1,276 server and 1,280 highres
roots with no failures and only the explicitly allowed map-population diagnostic
skip in each. Their wall times are approximately 932, 913 and 864 seconds,
including discovery/build. These broad shared-helper gates cost about 45 minutes
on this VM; reuse their exact source-qualified evidence for subsequent baseline
owner selections instead of rerunning unchanged tests unnecessarily.

## Completed qualification

Baseline commits are `5cc27785` (raw original observations) and `58f37c6c`
(original-library ordering contract). The corrected frozen tests are byte-for-byte
unchanged after conversion. Direct helper/class tests pass in all three profiles;
safe memory bridges and shop loading pass without skips. The accumulated sweeps
above, safe build/static checks and fresh production binaries/ABI checks pass.
The full suite matches all 304 known failure events and the 17 pass / two fail /
32 skip package outcomes. Headless character creation, explicit save, saved-map
load and resume pass against the retained reference. All phase source fingerprints
match. See [qualification](go-memory-qualification.json).

Six libc memory/string call paths are removed. Standalone production/reference
C remains **0 / 0 lines**; selected production cgo files remain **429** across
**three project packages**, with **79 embedded callback bodies**. Those counts
are unchanged because `alloc.go` still needs libc allocation and no callback
bodies/imports moved in this batch. External native bindings are unchanged.
