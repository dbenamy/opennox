# Leaf engine cgo dependencies

First batch of [internal glue removal](INTERNAL_C_GLUE.md): remove the empty
`import "C"` from `internal/binfile/memfile.go` and `server/audio_event.go`, and
replace the Linux socket adapter's C-header `FIONREAD` macro with Go's
`syscall.TIOCINQ`. Preserve the existing ioctl syscall, uint32 result and errno.
The same batch removes 31 unused C imports and declaration-only preambles in
`legacy`, identified by the selected-file inventory and individually reviewed for
C selectors, export directives and build flags. Their Go bodies remain unchanged.
Native client libraries/bindings are unchanged.

## Original behavior and acceptance

The Linux/386 C preprocessor defines FIONREAD as `0x541B`, and the qualified Go
syscall package defines TIOCINQ as `0x541b`. The adapter already invokes the Go
syscall API; no C function invocation is involved in this operation.

Five new socket roots exercise direct `netCanRead` and actual `canReadConn` on
loopback UDP: empty queue, next-datagram byte count with two sends, post-drain,
queued zero-length datagram and invalid-descriptor EBADF. The zero-length cases
use a non-consuming peek to establish arrival before observing a zero count.
Sends in each ordering case share a socket; reads and polling are bounded.
The existing stream handshake test remains in the package selection. All six roots
passed three times against the original C-header path before conversion.

The 32 existing root-package tests in `cgo-leaves-tests.txt` cover MemFile consumers
(sprite/object/image parsers, floor assets, skipped resource sections and audio
assets) and server audio event/owner behavior. Their source and frozen expectations
are unchanged. Run in default/server/highres before and after the change.

Original production evidence is the qualified transfer/sound checkpoint `b034c43e`:
source fingerprints match except for the newly added `_test.go` socket contracts.
Reuse its production baseline; fresh binaries and qualification are required after
conversion. Baseline logs and source identities live in `build/port-cgo-leaves/`.
The selected-build inventory before conversion is committed separately.

## Delegation and review

GPT-6 Luna drafted the socket contracts under an ignored path. Primary review
caught that a zero-count poll could pass before a zero-length packet arrived, and
that the nominal direct queue test used the wrapper for its nonempty checks.
Luna corrected both before integration. Primary checked the original constants,
selected the file/audio consumer roots and owns baseline/production qualification.
No cost or speed saving is asserted.

## Qualified conversion

Original baseline: `a28bdba7`. The unchanged six socket/handshake roots pass three
repetitions with cgo disabled. All 32 consumer roots pass without skips in each
of default/server/highres. Safe build and static checks pass. Three fresh 386/SSE2
production binaries pass ABI checks; the full suite matches all 304 known failure
events and package outcomes (17 pass, two fail, 32 skip). The headless scenario
creates a character, saves, loads and resumes against the original references.

All phase source-fingerprint checks pass. A separate diff check confirms the Go
code in all 34 edited files is unchanged except for the equivalent ioctl constant.
See [qualification](cgo-leaves-qualification.json) and the selected-file
[inventory before](cgo-leaves-inventory-before.json) and
[after](cgo-leaves-inventory-after.json).

Direct project cgo packages fall **6 → 3** in each production profile; selected
cgo files fall **463 → 429**. Remaining packages are `alloc` (one file), `ccall`
(one) and `legacy` (427). `binfile` and `server` still have transitive cgo
requirements, while the Linux `netstr` package itself now tests with cgo off.
Native client bindings are unchanged. Standalone production/reference C remains
**zero**, and embedded production callback bodies remain **79**.

The next draft separates six libc memory/string helpers from the allocation
ownership migration. It remains ignored and unqualified until original-path
captures and review are complete.
