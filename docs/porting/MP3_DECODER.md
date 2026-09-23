# MP3 decoder baseline and numerical work

The remaining standalone C file has six lines and includes the vendored minimp3
implementation. The current client decoder enables `MINIMP3_ONLY_MP3` and
`MINIMP3_NO_SIMD`, producing signed 16-bit PCM. Server builds exclude this path.
Headers, generated cgo and 81 production preamble bodies remain outside the
standalone C metric. No Go decoder implementation is introduced by this baseline.

## Qualified current-C asset baseline

The new tagged `TestMP3AssetBaseline` discovers RIFF/WAVE containers and selects
MP3 using the parsed decoder format, not the filename extension. It requires the
extracted assets and captures each input SHA256, format, sample rate/channels,
reported sample count/PCM hash, ordered Decode result hash, and full output-buffer
write hash. Each call receives 2,304 int16 slots with eight sentinel words on either
side. The interior is zeroed before each call to make unreported writes observable.
The test preserves the existing wrapper's returned-count semantics.

Three fresh current-C captures are identical. After freezing the expectation,
both default and highres client profiles pass. All 1,246 asset input hashes were
independently rechecked, and all reported PCM hashes match the current-C results
recorded in the preceding production full-suite log. The capture contains
172,369,152 reported samples and 300,498 Decode calls. Frozen SHA256:
`7dc3362576eb00660fd68a836d17af6012c4d1c0595d137958109d62228add0b`.
See [mp3-assets-c-qualification.json](mp3-assets-c-qualification.json).

All shipped MP3 Dialog assets are mono at 22,050 Hz. This baseline does not cover
stereo, other rates, malformed frames, seek behavior or streaming boundaries.
The sentinels detect adjacent overwrites, not arbitrary C memory errors. These
limits must be addressed before accepting a general decoder replacement.

Only a tagged `_test.go` file is new; every previous production source hash is
unchanged. Official default/server/highres/safe package file selections exclude
it, and all four previously qualified production binary hashes were rechecked.
The preceding production/ABI, known-suite and headless/save-load evidence is
explicitly reused on that basis. No historical audio golden changed.

## Scalar-arithmetic diagnostic

Retained 386 decoder disassembly shows x87 arithmetic in the inspected transform
and synthesis routines. An isolated build of the same observation test with
`CGO_CFLAGS=-O2 -g -msse2 -mfpmath=sse` changes all 1,246 PCM/full-buffer hashes,
while preserving every input, metadata, sample-count and Decode-sequence field.
**All 1,246 resulting PCM hashes match the existing historical audio goldens.**
Its capture SHA256 is
`e0688114198dc50b4acf9b5695ea4cd8bf6dbb9f91d53d1e9d550c02ff1e2763`.

This controlled comparison identifies scalar arithmetic selection as the source
of the observed shipped-asset PCM discrepancy on this target. It is not yet a
production flag change, proof of a future Go port's numerical parity, or a speed
benchmark: separate run timings include compilation and other variation.

The next step is a separate, reversible correction scoped to the audio package
on 386: enable SSE2 scalar arithmetic, consistent with the agreed CPU target.
Qualify against the existing audio goldens and unchanged non-audio expectations.
Retain the current-C capture above as historical evidence; any new asset-observation
expectation must be explained as this deliberate arithmetic correction, not a
silent golden update to accommodate a port. The production setting is not applied
in this baseline commit.

## Subsequent decoder port

Establish frame/state/seek and missing format-boundary contracts before replacing
the decoder. Build coherent Go internals before switching production; calling Go
for every bit read from retained C would add overhead at the wrong granularity.
The bit reader mutates parser state, and inactive Layer I/II callers do not define
the active MP3 input domain. Preserve the vendored implementation's license and
attribution when translating its algorithms/tables. Numerical fidelity and final
integration remain primary-agent decisions.

Luna drafted the asset fixture and read-only helper map. Primary removed an
extension prefilter, made unfrozen collection explicit, batched PCM hashing, and
owned repeated captures, independent PCM comparison and the arithmetic diagnostic.
No subscription savings are inferred from these results. The helper's initial
filename-only inference about missing MP3 assets was corrected: Dialog contains
MP3-in-WAV files despite having no `.mp3` filenames.
