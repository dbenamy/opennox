# MP3 PCM synthesis and filterbank state

## Original-C baseline

The actual scalar-SSE2 int16-output C implementation agrees across three native
386 runs and a run with undefined-behavior and float-cast-overflow sanitizers.
There are196,616 PCM conversion cases,390 synthesis pairs,408 synthesis operations
and360 granule sequences containing576 steps. See
[capture provenance](mp3-synthesis-c-capture.json).

Conversion cases include every half-integer boundary in the int16 domain and its
two adjacent float32 values, signed zero, subnormals, largest finite values and
infinities. NaNs are excluded because the original float-to-integer conversion
has no defined result for them. Preserve the source rounding, including its
negative quirk: -1 maps to0, whereas -1.5 maps to-2. Do not replace this with a
standard rounding function under a behavior-preserving port.

Filterbank cases cover mono/stereo, zero and boundary band counts, every spectral
row in both channels at the two active columns, each pair-history coefficient,
finite signed/subnormal patterns, smaller unsaturated amplitudes, retained tails,
and repeated granules with channel changes. The latter retain QMF and scratch
state while refreshing input and PCM sentinels. Mono copies only even QMF words;
the odd words remain unchanged and matter on subsequent stereo calls.

The capture records complete explicitly sized spectral, scratch and QMF arrays,
PCM results/tails and external guards. Primitive pairs additionally check input
immutability. All240 window coefficients independently match C. Source capture
includes the actual header and only imports pattern/I/O helpers from the prior
capture tool; both tool hashes are recorded. No copied C algorithms.

Magic NMP3SYN1 precedes byte opcodes1–4. Integer and float words are little endian;
PCM arrays are little-endian int16. All sequence requests precede all step results
in the frozen fixture. The capture tool defines the complete protocol.

This remains private decoder preparation, unimported by production. Standalone
C remains six lines/one file, reference C zero; the active header and81 preamble
bodies are outside that source-size count. No new complete-decoder PCM claim.
