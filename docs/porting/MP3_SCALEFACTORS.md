# MP3 scalefactor bytes and quarter-power scaling

This connected decoder preparation covers `L3_read_scalefactors`, `L3_ldexp_q2`
and `L3_decode_scalefactors`. The complete C decoder remains active until the Go
path is assembled and qualified. Standalone C remains six lines/one file, reference
zero; the third-party implementation header and 81 preamble bodies remain separate.

## Original-C baseline

`tools/porting/capture_mp3_scalefactors.py` calls the real production header under
`-m32 -msse2 -mfpmath=sse`. Three native processes and an additional UBSan process
produce identical results for 36,842 records. The preceding SSE correction supplies
the arithmetic baseline; x87 outputs are intentionally not used. Provenance is in
[mp3-scalefactors-c-capture.json](mp3-scalefactors-c-capture.json).

- 14,112 byte-reader cases: actual partition shapes plus early termination,
widths 0–5, reuse masks and negative intensity mode, logical truncation, several
bit offsets and input patterns. Capture both complete 40-byte output buffers,
reader state and surrounding canaries, exposing retained tails and coded sentinels.
- 10,250 quarter-power cases: exponents 0–1024 over finite float bit patterns including
signed zero, smallest/largest subnormals, smallest normal, 1, 2, 2048, largest finite
and negative 1. Expected output is the exact IEEE float32 word, not a tolerance.
- 12,480 complete scalefactor cases: every compressor value for MPEG-1/2/2.5,
long/mixed/short layouts, normal/intensity/MS/mono contexts, both scales, subblock
gains, pre-emphasis, previous intensity bytes, truncation and CRC-position inputs.
All 40 output float words and intensity bytes are observed, alongside reader state,
canaries and unchanged granule information. Contexts are directed, not Cartesian.

The magic is `NMP3SCF1`; each byte opcode precedes its inputs and outputs. Words
are little-endian uint32, with signed int/float values serialized as their bit
patterns. The capture tool defines exact layout and synthetic initialization.
No shipped asset bytes, C algorithm copies or raw process pointers are retained.

All 134 integer lookup-table entries were compared independently against C.
The active gain path uses MAX_SCFI 44 and nonnegative quarter exponents; integer
shifts stay within the original defined domain. Decoder-output acceptance still
requires complete PCM/stream tests later; these helper results do not establish
end-to-end Go audio playback or a performance result.
