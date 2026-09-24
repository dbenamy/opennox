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
empty strings, exact comparison results, and bounded writes up to 65,535 bytes.
The compare captures include input bytes, lengths/offsets and raw integer returns;
there are 87,961 memcmp and 107,161 strcmp cases. Independent expectations check
ordering and complete source/guard preservation. Exact frozen hashes match in
three separate original-libc processes; expectations are frozen before conversion.
The existing `TestSafeMemoryBridges` also freezes exact return magnitudes in its
142 records, although its independent comparison assertions check only sign.
Keep that historical hash unchanged. `TestShopStockLoading` exercises the actual
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
original capture and acceptance. Comparison magnitudes in the draft are
provisional until verified against original results. Do not substitute sign-only
returns merely because libc's portable contract permits them.

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
