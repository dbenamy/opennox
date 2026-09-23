# MP3 frame scanning and byte reservoir

This private Go decoder preparation covers frame matching/finding, initialization
and reservoir save/restore. Production continues using the complete C decoder.
Standalone C remains six lines/one file, reference C zero; the active third-party
header and 81 C preamble bodies remain outside this source-size metric.

## Original-C baseline

The self-contained `tools/porting/capture_mp3_stream.py` calls the real production
header. Three native 386 runs and a UBSan run agree on 2,774 records. The initial
2,756-record capture also matched the independent preliminary fixture before
the tool was made self-contained; the subsequent extension is described below. [Capture provenance](mp3-stream-c-capture.json) records source,
compiler flags, hashes and per-operation counts.

- 1,650 frame searches: short/empty packets, ordinary and free-format frames,
  leading junk, exact/truncated lengths, padding, the free-format search limit,
  seeded arbitrary input and valid Layer I/II headers still used by this scanner.
- 190 frame matches: follow-on header checks, incomplete frames and bounded
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

All domains are bounded: packets at most 49,152 bytes, nonnegative frame hints
up to 2,304, stored history/main-data-begin 0–511, current payload up to 2,304,
and scratch bytes up to 2,815. The parser fixtures do not establish behavior for
integer overflow, negative copy lengths or unsupported overlapping C memcpy inputs.

The first baseline (`5a911bce`) had 2,756 records. Caller review then identified the
actual audio reader's 49,152-byte buffer, exceeding the initial 32,768-byte capture
bound. Eighteen appended cases cover that boundary; every earlier fixture byte is
unchanged. The expanded capture repeats across three processes and UBSan. Initial
Go qualification also exposed a redundant-variable cleanup compile error and an
overrestrictive save-position guard. Both are implementation fixes, not changed
C expectations. Capture bounds must not become arbitrary production limits.

## Go implementation qualified

The initial C baseline is `5a911bce`; the caller-size extension is `7bd7e5d6`.
Both were committed/pushed before accepting Go stream helpers. Default/server/
highres/safe and CGO_ENABLED=0 each pass all three MP3 roots:819207 integer,
36864 side-info and2774 stream cases. `go vet` passes. Independent checks cover
no-header offsets and retained hints, exact single-frame acceptance versus sync
matching, partial init including NaN payloads, reservoir concatenation and missing
history behavior, retained tails, byte rounding,511-byte cap and negative remainder.

Primary wrote/reviewed the runner; Luna supplied the implementation draft and
reviewed capture bounds. Primary fixed two late draft edits (a deleted local still
referenced and an overrestrictive save-position guard), and expanded captures after
finding the actual48KiB caller buffer. Failed logs remain in
`build/port-mp3-stream/go-{first-build,save-guard}-failure`; acceptance is in `go/`.
No historical expectation was changed to accommodate an implementation mismatch.

[Qualification](mp3-stream-go-qualification.json) verifies all previous source
hashes and the four production binary hashes. Only private code/tests were added to
the still-unimported package; prior production dependency selections remain valid.
SSE production/ABI/gameplay evidence is explicitly reused. No production build,
full-suite or complete-Go-decoder PCM success is claimed for this unwired chunk.
C stays six lines/one file, reference zero,81 production C preamble bodies.

Next: connected scalefactor byte parsing, quarter-power scaling and scalefactor
float output. Capture float bits under the qualified SSE2 arithmetic and preserve
uint8 wrapping, retained bytes and exact expression rounding before integration.
