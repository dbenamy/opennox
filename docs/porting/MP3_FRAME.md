# Complete MP3 frame decoding

## Original-C baseline

The original scalar-SSE2 decoder agrees across three native386 processes,
UBSan with float-cast-overflow checking, and a history-poison check on2,682 calls
in544 sequences. The generated streams cover all MPEG versions/rates, stereo,
intensity/MS and mono, CRC/padding, block layouts, nonzero coded/escape spectra,
missing/available reservoir history, invalid side info, truncated/prefixed scans,
free format, unsupported Layer I/II, reset, metadata-only and the real48KiB input
extent. Inputs are synthetic, not shipped asset bytes.

Each call records return count, all frame-info fields, persistent decoder state,
full2,368-sample output storage including tails, and immutability/external guards.
The NMP3FRM1 fixture serializes all requests before each sequence's outputs;
[provenance](mp3-frame-c-capture.json) records compiler/source/tool hashes.
No C algorithm copies or changed production algorithms are used by the capture.

## Original storage behavior: two findings

Thirty QMF words are unspecified retained storage: indices896+4*i+2 and+3 for
0<=i<15. Scalar synthesis overwrites them before reading: pair operations access
only row offsets60..63, and each descending filter loop writes history row14
lanes2/3 before its matching tap reads. Normalizing only these30 words in fixture
serialization removes process-dependent stack bytes. All other state and every
PCM byte remain compared. Poisoning these slots before every C call produces
identical observable results. Primary and Luna independently checked the lifetime.

The shared synthetic side-info builder originally set MPEG-1 private bits. The C
parser includes those bits in its SCFSI shift register, which can feed nonzero
reuse flags to the first granule before the corresponding intensity positions
are initialized. Six intensity/escape sequences (30 calls) then depend on stack
contents; native vs sanitizer/zero-auto-init diagnostics differ in PCM as well as
state. [Diagnostic hashes](mp3-frame-private-bits-diagnostic.json) preserve this
finding. Private bits are normally ignored by a decoder, so this is an original
implementation bug rather than an intended format requirement.

The frozen generated streams clear private bits while retaining all actual SCFSI
reuse bits. Their complete results exactly match the earlier zero-auto-init
diagnostic too. Go scratch is deterministically zero-initialized; recreating
uninitialized C stack contents is not a port requirement. Review this behavioral
clarification at final integration; do not silently claim bit-identical results
for C executions that depend on uninitialized storage.

The Go workspace keeps1,152 spectral floats adjacent to40 scalefactors, preserving
the previously captured8kHz mixed-block extent. Scratch lookahead preserves the
actual2,815-byte capacity. Production remains on C until whole-decoder asset and
integration qualification. Standalone C remains six lines/one file, reference
zero, plus the active header and81 preamble bodies.
