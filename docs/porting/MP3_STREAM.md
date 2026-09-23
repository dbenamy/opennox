# MP3 frame scanning and byte reservoir

This private Go decoder preparation covers frame matching/finding, initialization
and reservoir save/restore. Production continues using the complete C decoder.
Standalone C remains six lines/one file, reference C zero; the active third-party
header and 81 C preamble bodies remain outside this source-size metric.

## Original-C baseline

The self-contained `tools/porting/capture_mp3_stream.py` calls the real production
header. Three native 386 runs and a UBSan run agree on 2,756 records. The independent
preliminary capture also produced identical fixture bytes before the tool was
made self-contained. [Capture provenance](mp3-stream-c-capture.json) records source,
compiler flags, hashes and per-operation counts.

- 1,637 frame searches: short/empty packets, ordinary and free-format frames,
  leading junk, exact/truncated lengths, padding, the free-format search limit,
  seeded arbitrary input and valid Layer I/II headers still used by this scanner.
- 185 frame matches: follow-on header checks, incomplete frames and bounded
  zero-progress free-format matching. The first header is structurally valid.
- Five initialization patterns: only header byte zero changes; all remaining
  decoder bytes retain their initialized values, including overlap/filter state.
- 425 reservoir restores: absent/partial/full history, byte and bit offsets,
  zero-length and maximum current payloads. Capture includes the full scratch
  byte array, reader state/alias, and unchanged decoder/outer-reader checks.
- 504 reservoir saves: byte rounding, the 511-byte cap, retained output tails,
  empty and negative remainders, and unchanged scratch-state checks. Negative
  saved remainders are tested only in this defined operation and never passed
  to the C restore function's copy-length arguments.

The fixture magic is `NMP3STR1`. Byte opcodes 1–5 precede the request fields and
expected output; uint32 fields are little-endian and signed results use their
int32 bit patterns. The capture tool defines the full protocol and input patterns.
No raw pointers, struct padding, copied C algorithm or shipped asset bytes are
retained in the fixture. Init uses a boolean comparison of initialized C bytes;
Go independent checks must verify the corresponding complete typed state.

All domains are bounded: packets at most 32,768 bytes, nonnegative frame hints
up to 2,304, stored history/main-data-begin 0–511, current payload up to 2,304,
and scratch bytes up to 2,815. The parser fixtures do not establish behavior for
integer overflow, negative copy lengths or unsupported overlapping C memcpy inputs.
