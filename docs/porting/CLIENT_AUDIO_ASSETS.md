# Client sound-definition assets

Round 1 of the two-round process trial. Six connected C functions (321 physical
block lines) cover client AUD/AVNT readers, a single-record helper, sound-slot and
delay getters, and case-insensitive sample-catalog lookup. This is the existing
small scope used to establish the runner; it is not a demonstration of larger
batch sizing. The live client dispatch in things.go calls these readers. Server
Audio readers are separate and remain unchanged.

## Contracts and baseline

The real 1023 × 200-byte sound table is initialized through Sub_451850 and the
actual Timer implementation. Fixtures fix only the clock provider, count its
calls, and restore globals/table/clock. Catalogs use existing audioStructXxx and
36-byte audioStructYyy owners and their real Free path. Input MemFiles use the
actual allocation and cursor implementation. Only the interned name pointer at
row offset 84 is normalized; its unchanged identity is checked independently.

Eight groups /1,433 frozen records cover:

- Disabled/invalid/edge sound IDs and gated signed delay values.
- Empty, unique and duplicate catalog keys, case folding, read-only ownership,
  and catalog indices around 32767/32768/65534/65535/65536.
- AUD signed fields, timer initialization, positive-only delay updates, mode
  rejection with partial state/cursor, unknown-name and disabled skipping.
- Every AVNT tag byte, reordered/repeated tags, sample-list limits, raw words,
  signed/unsigned fields, and success on unknown sounds despite inner tag failure.
- Sample-index narrowing, missing entries, embedded NUL, extension stripping and
  scratch tails, 127-byte AUD and 255-byte AVNT names, 32-result AVNT cap.
- Bulk signed counts, concatenation and first-error stopping.
- Short nonempty scratch views with sufficient capacity; rejection of empty or
  insufficient-capacity buffers without changing the definition owner.
- Real audio.idx/audio.bag ownership and lookup of every catalog name, then
  selected real samples bound through both parsers. Missing assets must fail
  the qualification driver instead of being counted as successful coverage.

Full definition and scratch hashes, returns, consumed cursor and timer-call counts
are frozen in [the capture manifest](client-audio-assets-captures.json). Independent
assertions distinguish the contracts above. The initial real-catalog test omitted
handle initialization; it failed before capture. Reusing handles.PortTestInit
fixed fixture setup. Production was unchanged.

Valid-input limits: negative signed AUD name lengths, truncated records, and AUD
lists with more than 32 successful samples are outside the original C contract
qualified here. AVNT's bounded list is tested beyond 32 entries. Nullaudio gameplay
exercises startup, not subjective playback quality; existing PCM failures remain
visible in the exact full-suite comparison.

## Qualification

Use [the batch manifest](client-audio-assets-batch.json) with run_batch.py. The C
phase repeats default captures in independent processes and checks server and
highres against the same frozen hashes. The native phase checks affected callers
and owners in all three variants, then all production builds, ABI, exact known
failures and two reference gameplay scenarios. Intermediate commits describe
pending gates; only the completed native report establishes round completion.

Before replacement, production source remains at qualified f0db7a7b; new files
are guarded fixtures. Prior production integration evidence remains applicable.
Retain sound-slot/delay C callbacks for live consumers; retire both reader APIs,
the private record/lookup helpers and the previous AVNT-inner export whose last
C consumer moves in this batch. Go public wrappers retain the backing-capacity
contract established by the preceding readers batch.

C qualification passed: default 24.175s, independent repeat 10.691s, server
105.991s, highres 23.798s. All eight roots completed without skips, all 1,433
records matched each time and source fingerprints stayed unchanged. Native
replacement and its completed-round gates are still pending. Detailed local
artifacts: build/port-client-audio-assets, including failed fixture logs.
