# First bounded C inventory — 2026-09-10

The preserved baseline build graphs select 152 repository-owned C translation
units for standard/highres clients and 151 for the server:

| Package below module root | Standard | Highres | Server |
| --- | ---: | ---: | ---: |
| legacy | 148 | 148 | 148 |
| legacy/cnxz | 3 | 3 | 3 |
| legacy/client/audio/ail | 1 | 1 | 0 |

These are build-selected units, not remaining-function counts. Codec sources
included by preambles/headers and external libraries are not additional units in
this count. A file selected for server compilation may still contain client code;
linker retention is a separate measurement. No runtime execution coverage is
inferred from either count.

Reproduce the selection with the target environment from RECOVERY.md:
`go list -deps -json ./cmd/opennox` from src, adding `-tags highres` or `-tags server` before the package argument
for variants (the build driver uses `cmd/opennox` for all three). Feed the JSON and the
corresponding built binary to `tools/porting/c_inventory.py`. For example, from
the repository root:

```bash
python3 tools/porting/c_inventory.py build/baseline/logs/build-graph-default.json \
  --binary build/baseline/bin/opennox \
  --symbol nox_xxx_protectionStringCRC_56FAC0 \
  --symbol nox_xxx_protectionStringCRCLen_56FAE0
```

The script is read-only. Existing `internal/callgraph` invokes source conversion
and removes generated directories; it was not needed for this bounded audit.

## First leaf: protection byte checksum

Both named symbols are defined in the baseline standard, highres and server
binaries (`nm --defined-only`). They originate in `legacy/GAME5_2.c`. This proves
link retention, not execution in the warrior scenario.

The function XORs little-endian 32-bit words, ignoring the final zero to three
bytes. Despite its historical name, it is not a polynomial CRC. The nullable
wrapper returns zero for a null pointer regardless of length. The direct function
only accepts null when there are no complete words. Neither mutates input or
accesses global state. There are no callbacks or downstream subsystem calls.

C callers are the player-name protection path in `GAME1_1.c`, `sub_56FB00`
(validation, called from `GAME3_3.c`), and `sub_56FB60` (item init-data/name
checksum). The nullable wrapper also calls the direct function. The root Go
package already implements the same operation in `protectBytes`; the conversion
can share that implementation while preserving both C ABI entry points.

Tests first compare the existing Go implementation to a test-only historical C
reference, copied from revision 0e9d2e1f without algorithm edits. The reference
uses the `porttest` build tag and is absent from normal builds. Fixed vectors,
all lengths 0–1024 at eight byte alignments, 1,000 larger deterministic generated
buffers, mutation checks, null lengths including UINT_MAX, and fuzz properties
cover the relevant arithmetic and buffer semantics. Invalid non-null pointers
or lengths exceeding their allocation are outside the C contract and are not
used as differential inputs.

Larger alternatives inspected (XP/health, item-name formatting, memfile reads)
have floating-point/global-state dependencies, callbacks, or mismatched existing
Go EOF semantics. Those need different test boundaries and were deferred.
