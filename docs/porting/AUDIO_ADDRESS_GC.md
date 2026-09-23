# Audio address words under garbage collection

Status: correction qualified. C remains 114 lines /6 files. No remaining
global/blob storage migration has been applied.

The broad default consumer sweep for the next storage baseline exposed an existing
failure in TestClientAudioStreamsBuffers. An opaque sample address was stored in a
Go pointer-typed field of a foreign audio record. The write barrier queued that
number for GC; if it happened to fall in a freed Go heap span, the runtime rejected
it. The list initializer shown at the crash site was flushing the barrier buffer,
not establishing that the list links produced the bad value.

A new subprocess regression obtains the numeric address of a released Go allocation
and transfers that opaque word through the existing chunk/buffer/voice routines
while another goroutine performs 100 collections. It never dereferences that word.
The test reliably reproduced the old failure before the correction. No frozen audio
expectations are changed, and the broad suite is retained.

Keep the three sample-data fields (buffer, chunk and voice), the chunk initializer's
address argument, and the refill routine's temporary addresses as uint32 ABI words.
Clear/copy/advance those words without Go pointer write barriers. Convert at actual
sample read/copy/device handoff boundaries. Four-byte layout and modulo-32 address
arithmetic are unchanged on the 386 target. Other pointer fields and shared list
algorithms are outside this bounded correction.

The prior 52-owner conversion remains committed as `05182782` with its recorded passing
qualification; the new broad sweep uncovered this additional GC-sensitive path.
Fresh production qualification now covers this fix and the unchanged-C storage
baseline, establishing the starting point for the remaining-owner conversion.

## Qualified result

The red subprocess reproduces the old GC failure; three corrected child processes
pass, followed by the regression in default/server/highres acceptance. All 884 raw
storage patterns and 54,905 numeric patterns retain their frozen captures. The broad
consumer sweep passes 2,291/2,280/2,291 roots without skips; static checks, safe build
and symbol audit, three production binaries/ABI, exact known-suite comparison,
fresh headless gameplay and explicit save/load pass. Safe runtime was not tested.
The known suite comparison preserves existing failures; it is not a green full suite.
All final source fingerprints match and preflight uses the final default production
binary. See [audio-address-native-qualification.json](audio-address-native-qualification.json).
C remains 114 physical lines /6 production files, zero reference C.
