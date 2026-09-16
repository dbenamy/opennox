# Complete NXZ map compressor

Status: C baseline qualified in all three targets and an independent repeat.
Go implementation passes focused frozen and independent contracts; production
qualification pending. C baseline commit: **3ee3d454**. Starting size: **71,252 C lines /90 files /zero reference C**.
Scope: `legacy/cnxz/nxz_comp.c` (1,892 lines) and `nxz_common.c` (18 lines).

## Scope and contracts

Complete the codec's Go ownership: sliding dictionary, match selection, adaptive
symbol codes, bit packing, block boundaries and file wrapper. The only production
C entry to compression is the existing Go map wrapper in `legacy/maps.go`, which
calls CompressFile. Compressor-internal symbols have no other callers.

Reuse all decoder contracts and 115 synthetic frozen compressed outputs; add 48
mixed-copy cases around symbol rebuild, 64 KiB window, 521-byte match and
500,000-byte file chunk boundaries. Compare compression with all 50 shipped NXZ
files, using actual filenames. Verify source preservation, destination truncation,
same-path compression/decompression and destination preservation on input errors.
Independent expanded inputs and shipped compressed files supplement C hashes.
Repeat captures in separate processes and qualify default/server/highres.

## Audit notes and review items

The C match search uses five-byte rolling lookahead, including bytes beyond the
current 500,000-byte block. Preserve following file bytes for block-end hashes.
The final file tail can cause C to read a few bytes beyond its allocation during
hash maintenance. Those hashes are unused after the final block; the Go draft uses
explicit zero padding, without changing valid encoded output expectations.

The existing empty-input allocation panic and unchecked file-size arithmetic will
be replaced with explicit empty-file output and addressable-size checks, with
independent contracts. These are bounded, reversible corrections under the user's
standing authorization. Do not regenerate valid-file goldens to absorb differences.

## Evidence

Manifest: `map-encode-batch.json`; affected selection: `map-encode-tests.txt`.
Local initial C checks: `build/port-map-compression/c-focused.log`,
`c-matches.log`. Frozen match capture SHA-256:
`c8c09c83ad7cd5b40da7986ececb7f59fef83ca902d91a1491fcf659bf19728f`.
C gate: `build/port-map-compression/c-qualified/result.json`, all steps passed.
The first native behavioral run passed all frozen expectations (5.071s driver).
Additional native contracts pass (5.520s): empty/sparse-oversized input, every
length 4–521 against 19 distance boundaries, and 12 small block sizes down to one
byte, including repeated table rebuilds. No goldens changed.

Working C size: **69,342 physical lines /88 files /zero reference C**, a reduction
of **1,910 lines**. Initial scope estimate overstated the shared helper by two
lines; the tracked physical-line counter is authoritative. Full gates pending.
