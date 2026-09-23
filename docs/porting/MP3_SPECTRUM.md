# MP3 stereo, reordering and antialiasing

This connected private decoder preparation covers the active scalar stereo helpers,
short-window reordering and antialiasing. Production still uses the complete C
decoder. Standalone C remains six lines/one file, reference zero; the implementation
header and 81 production preamble bodies remain outside that source-size metric.

## Original-C baseline

`tools/porting/capture_mp3_spectrum.py` includes the actual header and uses the
committed side-info input generator to obtain original-C band tables. Three native
386 scalar-SSE2 processes and UBSan agree on 4,364 records. The fixture records exact
float bits, retained tails, intensity-position bytes and guard/immutability checks.
[Provenance](mp3-spectrum-c-capture.json) includes both tool hashes and compiler flags.

The cases cover 60 mid/side operations, 126 intensity-band operations, 378 top-band
searches, 3,600 complete intensity passes, 144 reorder operations and 56 antialias
operations. Contexts include all three MPEG versions/rates and long/mixed/short
layouts, raw MS-extension bit clear/set, both MPEG-2 intensity shifts, valid and
sentinel intensity positions, zero/dense/sparse right channels, both channel offsets,
empty/boundary band counts and finite signed/subnormal bit patterns. Luna separately
reviewed domain bounds, serialization and meaningful caller contexts before freeze.
All 30 pan/antialias float constants independently match the C float32 values.

The magic is `NMP3SPC1`; byte opcodes1–6 precede requests/results, with little-endian
words for integers and float bit patterns. The capture tool defines the complete
protocol. No C algorithm copies, raw pointers or shipped asset bytes are retained.

## Mixed-block extent: preserve and review at integration

Primary and Luna independently confirmed a reachable source quirk for MPEG-2.5
sample-rate index2 (8kHz), mixed blocks. The parser selects mixed row1 and six long
bands, whose widths sum48. The caller nevertheless advances by72 spectral words;
the selected reorder suffix contains528 words. The operation therefore ends at
channel-relative offset600, crossing24 words past the nominal576-word channel.
Channel0 touches channel1; channel1 touches the following scalefactor work area.

A native386 compile-time assertion verifies that the original C scratch layout
places40 scalefactor floats immediately after the1152 channel floats. The capture
uses one explicitly sized1192-float array plus canaries, making the helper operation
defined while observing all those writes. This is not a claim that the original
caller's per-channel pointer arithmetic is safe in the C abstract machine.

For this helper port, preserve the observed copy extent in bounded Go slices; do
not silently clamp it. At complete-decoder integration, represent the adjacent
spectral/scalefactor workspace explicitly so these writes can be reproduced without
unsafe access, or qualify a separate behavior correction. This is a documented,
reversible compatibility decision for review, not a new decoder-format guarantee.
