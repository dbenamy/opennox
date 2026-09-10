# Baseline failure diagnosis — 2026-09-10

This bounded investigation precedes the first C-to-Go conversion. It does not
make the whole suite green. Commands use the 386/CGO environment documented in
[RECOVERY.md](RECOVERY.md), with `NOX_DATA` pointing to the extracted Nox directory.
Run package commands from `src`. Diagnostic dependency changes were made only in
ignored disposable copies; production dependencies and compiler flags are unchanged.

## Blob tooling

Commit `4df9834e` fixes `getStatic` dropping the plus between two dynamic Go
terms. The new regression test checks parsing, idempotence and independently
evaluated arithmetic for 13 expressions at three input assignments. It fails
before the fix and passes afterward; `TestFormatAccesses` also passes.

Reproduce with `go test -count=1 -run 'TestFormatGoOffsetPreservesTerms|TestFormatAccesses' ./internal/blobs`.

`TestReadBlobs` still assumes root `memmap.go`, `memmap.c`, and `GAME_data.c`.
Current storage uses legacy shims, embedded `.dat` data in
`common/memmap/nox/blobdata`, and generated pointer initializers. Merely changing
paths to `legacy/cgo_blobs.c` would lose initialized data by treating it as zero.
Do not use the split/write tooling on the current layout until its reader and
writer preserve embedded bytes and pointer relocations in a tested round trip.
This modernization is not required for the initial read-only C inventory.

## Rendering

The old color library, `github.com/noxworld-dev/opennox-lib` at
`v0.0.0-20221207184240-db559cc95748`, expands maximum five-bit RGB channels to 248.
The current pinned `github.com/opennox/libs` at
`v0.0.0-20241214142037-24bbb986fab8` expands them to 255 in `ColorNRGBA` and
`ColorRGBA` in `color/rgba5551.go`.

All six particle cases produce identical decoded pixels on 386 and amd64.
Applying the historical expansion at PNG export reproduces all six original
PNG hashes. Re-encoding the checked PNGs with Go 1.19.13, 1.23.12 and 1.26.0
also produced identical PNG bytes. Commit `f358c129` replaces only these verified
particle references with SHA-256 over dimensions and little-endian framebuffer
words. Tests check sensitivity to pixel and dimension changes and independence
from stride padding and subimage origin. Particle tests pass on both architectures.
These tests cover framebuffer generation; they do not freeze RGB export behavior.

Reproduce with `go test -count=1 -run 'TestDrawParticle|TestPixelHash16' ./client/noxrender`.

Sprite tests need a separate compatibility decision. Applying old expansion only
at PNG export reproduces 44 sprite cases but does not fix the whole matrix.
The golden-era revision `c62202f5`, with Go 1.19.13 and the same assets, passes
`TestDrawImage/APA00001/default` against the same reference used today. In a
copy of current source, replacing **only** the current dependency's
`color/rgba5551.go` with the historical file makes the entire `TestDrawImage`
matrix pass. The dependency was copied locally and selected using a `go.mod`
replace directive; the module cache was not edited. This isolates the color
conversion change, which affects intermediate rendering inputs as well as export.

To reproduce that experiment, copy the pinned libs module to a writable temporary
directory, substitute that one historical file, point a disposable source copy's
`github.com/opennox/libs` replacement at it, and run
`go test -count=1 -run '^TestDrawImage$' ./client/noxrender` with assets.
Do not use earlier May 2022 tests as evidence for today's goldens: their expected
hashes differ. No sprite goldens or production rendering code were changed.
Before porting this rendering path, decide which color behavior is intended and
establish references for it explicitly; do not blanket regenerate PNG hashes.

## Audio

Three dialogue files were decoded with the same assets using default 386,
amd64, and 386 with `CGO_CFLAGS="-O2 -g -msse2 -mfpmath=sse"`. GCC's default
32-bit configuration uses x87 (`-mfpmath=387`, SSE2 disabled). The decoder already
sets `MINIMP3_NO_SIMD`; that does not prevent compiler floating-point differences.

| File | int16 samples | Samples differing: default 386 vs amd64 | Maximum absolute difference |
| --- | ---: | ---: | ---: |
| C1CAP01E.WAV | 468288 | 206 | 1 |
| C1CAP07E.WAV | 401472 | 205 | 1 |
| C1CAP16E.WAV | 67968 | 52 | 1 |

Sample counts agree in all three configurations. SSE-configured 386 output is
byte-identical to amd64 and matches the original expected hashes for these files.
This demonstrates floating-point evaluation sensitivity for these samples, not
asset corruption. It does not establish equivalence for all dialogue or music,
or assess audible playback quality. No compiler flags, tolerances or PCM goldens
were changed in production.

Existing-test reproduction:

```bash
go test -count=1 -run '^TestAudioDecode$/^dialog$/^c1cap(01|07|16)e.wav$' ./legacy/client/audio/ail
CGO_CFLAGS='-O2 -g -msse2 -mfpmath=sse' go test -count=1 -run '^TestAudioDecode$/^dialog$/^c1cap(01|07|16)e.wav$' ./legacy/client/audio/ail
GOARCH=amd64 PKG_CONFIG_LIBDIR=/usr/lib/x86_64-linux-gnu/pkgconfig:/usr/share/pkgconfig go test -count=1 -run '^TestAudioDecode$/^dialog$/^c1cap(01|07|16)e.wav$' ./legacy/client/audio/ail
```

Before changing decoding code, choose target-specific exact PCM references or a
justified cross-platform comparison policy, and expand differential coverage to
representative codecs, channels, lengths and malformed inputs. Do not change CPU
requirements merely to match a hash.

## Local evidence and next step

Ignored `build/diagnosis` contains before/after formatter logs, architecture
particle comparisons, `render-legacy-colors.log`, `golden-era-sprite.log`,
`sprites-old-color.log`, `audio-summary.json`, PCM outputs, disposable source
copies and `final-suite.jsonl`. These generated artifacts are not recovery
requirements; the findings and reproduction methods above are tracked.

Proceed with a bounded compiled/linked C inventory and select an independent
leaf whose relevant tests pass. Keep the known sprite, PCM and blob-reader
failures visible. Save/load coverage and unexplained baseline Player.plr byte
differences remain separate work. No C implementation has been replaced yet.
