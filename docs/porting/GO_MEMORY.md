# Go memory and string helpers

Status: original-path baseline qualified; no helper implementation changed.

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
regression. Latest qualified production baseline is `2db93f68`; only test files
are added before conversion, so production baseline evidence can be reused.

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
