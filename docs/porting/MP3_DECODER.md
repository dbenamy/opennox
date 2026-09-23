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

## Qualified SSE2 correction

The audio package now adds `#cgo 386 CFLAGS: -msse2 -mfpmath=sse`. This is scoped
to that package on the supported 386/SSE2 target; it does not set global C flags.
The retained client transform disassembly contains XMM instructions. All 1,246
historical dialog PCM goldens pass unchanged in default and highres, alongside
the new guarded observation capture. The latter intentionally moves from the
recorded x87 hash to the SSE hash above because this step corrects arithmetic
behavior; no historical TestAudioDecode expectation changes.

Audio stream lifetime checks pass in both client profiles. The first invocation
mistakenly selected the root package, found zero tests and failed discovery; it
is preserved and is not qualification. The reviewed manifest selects `./legacy`.
Safe build/static, three fresh production builds/ABI, headless character creation
and explicit save/reload/resumption pass. The full suite removes exactly 1,249
audio failure events, while all 304 non-audio failure events remain unchanged.
Package outcomes are 16 pass, two fail and 32 skip. Actual failures are compared
without filtering against a stricter expectation derived from the prior log:
the audio package must pass and every non-audio failure must remain identical.
See [mp3-sse-qualification.json](mp3-sse-qualification.json).

Only the audio compiler directive and the explicit new observation expectation
change from the prior source fingerprint. Standalone C stays **six lines/one file**,
zero standalone reference C; production preamble bodies remain 81. This is a
correctness prerequisite for the decoder port, not a decoder algorithm conversion.
No controlled performance result is claimed.

Ten superseded safe/callback executables were removed after hash/qualification
checks and checking 163 host processes, recovering 489,093,068 bytes. Current
address-adapter binaries and both callback benchmark binaries were hash-checked
and retained. Source, manifests, logs, captures and original assets remain.
Historical safe/callback/address-getter C-baseline finalizers require rebuilding
their old executables before replay; the consumed record is
build/port-artifact-cleanup/superseded-safe-callback-removed.json. Successful new
scenario asset copies have separate verified deduplication/restoration records.
