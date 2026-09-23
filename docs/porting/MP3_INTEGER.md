# MP3 bit reader and header arithmetic

This is preparation for a complete Go decoder. The production minimp3 decoder
remains active; these private helpers will not be connected through per-read
C-to-Go calls. Standalone C remains six lines in one file, reference C zero.
The included third-party header and C preambles remain outside that metric.

## Frozen original-C baseline

`tools/porting/capture_mp3_integer.py` includes the actual production header and
calls its helpers. Its C shim only serializes requests/results and initializes
synthetic input bytes; it contains no copied decoder algorithms. To recapture,
run `python3 tools/porting/capture_mp3_integer.py build/<new-directory>` on a
revision retaining that header, with a working native 386 GCC runtime.

Three separate processes produced identical output. A fourth, compiled with
`-fsanitize=undefined -fno-sanitize-recover=all`, also completed without diagnostics
and produced identical bytes. The initial sandbox run failed with SIGSYS before
capture; the successful native runs are in `build/port-mp3-integer/capture-host`.
Provenance and all three output hashes are recorded in
[mp3-integer-c-capture.json](mp3-integer-c-capture.json).

The compressed fixture contains 819,207 records:

- 196,864 validity checks: every second/third header-byte combination for three
  sync-prefix values, plus every first-byte value against one valid header.
- 524,288 comparisons: four representative base headers, every second/third-byte
  combination, in both argument orders. Includes free format, MPEG-2.5 and an
  invalid first sync byte; validation is deliberately asymmetric.
- 30,240 metric cases: all structurally valid version/layer bits, bitrate and rate
  indices, padding/private bits, two final-byte values and six signed free-format
  lengths including int32 extremes.
- 67,815 reader cases: five byte patterns, eleven lengths through 2,815 bytes,
  widths 0–32, initial offsets and positions near/beyond the logical end.
  Physical padding keeps the original C's zero-width read and pre-check pointer
  construction within allocated storage. Negative/overflowing state and widths
  above 32 are outside the contract.

The fixture protocol is `NMP3INT1`, then byte opcodes and little-endian uint32
fields. Operations are header/valid, two headers/compare, header/free-format/five
metrics, and pattern/length/start/width/value/final-position/limit/initial-position.
The capture tool is the precise wire-format reference. Fixture tests will consume
these frozen bytes without a C compiler or copied C reference implementation.

All 90 bitrate lookup entries were independently compared with the original C.
These cases do not qualify full MP3 decoding, floating-point transforms, seeking,
stereo output or malformed-stream behavior. Those remain subsequent work.
