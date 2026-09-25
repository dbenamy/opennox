# Internal audio-stream callback glue

Retired 20 C exports in `legacy/client_audio_streams_exports.go`. Sixteen
are fixture-only bridges; four provide live internal callback addresses. Preserve
Linux 386 layouts, external SDL/OpenAL bindings, and all existing captures.

## Routes and behavior

`audioStreamContextNew` installs Tick at context offset 216. The root service loop
`sub_486EF0` invokes it after global-enabled, global-mutating and context-flag
checks. Voice initialization installs Data/Loop/End at offsets 276/280/284.
`audioEventSampleRefill` invokes Data/Loop as void callbacks;
`audioEventSampleEnded` invokes End as an int callback and discards the result.
Root completion `sub_43EFD0` invokes End as void after ending the sample, guarded
by its one-time completion flag. All five dispatch sites must recognize the new
identities before any raw C fallback.

The conversion uses four stable byte identities and a common known-callback dispatcher, with
separate void-discard and int32 wrappers. Unknown C callbacks retain their original
call conventions. The separate OnData/OnLoop/OnEnd/OnStop callbacks, driver API,
context initializer, sample storage and external audio libraries remain intact.
Explicit compile-time assertions protect all four callback offsets. Fixtures call
the existing Go owners directly while preserving operation names/numbers, argument
widths and uint32 result bits, including negative status values.

The context initializer is traceable in-tree: blob descriptor slot 94020 points to
`Ptr_sub_43EA20`, which reaches root `sub_43EA20` through the legacy hook. This is
not an unresolved external callback. The root accesses at offsets 216 and 284
were missing from the helper's initial field-name audit; primary found them by
following raw-offset consumers. The corrected audit also fixes Tick's initially
reported offset of 212 (that is the Mutating field).

## Original-path contracts

The selected 86 roots include all direct audio-stream fixture users plus related
audio-event, asset, music, voice and briefing contracts. New root-service tests
exercise ordered contexts, all three outer guards, known tick timing, inner
mutating/period guards and a foreign C observer. New root-completion tests cover
native/foreign callbacks, negative callback results discarded by the root, state
mutation and duplicate-completion suppression. The zero sample handle has a no-op
End implementation in both supported backends; no physical audio device is needed.
Existing refill tests cover data/loop/end routing and byte-exact output across
chunk boundaries. Existing callback tests retain negative return captures.

Original-path baseline passes: all 86 exact roots in each of three profiles,
without skips/failures. Both new root contracts also pass in a second process
per profile. The accepted conversion preserves existing captures.
The original-path baseline's production source matches qualified UI revision
`3d47a346`, allowing its production baseline to be reused. Final qualification
passes focused contracts, safe/static checks, all three
production ABI checks, the exact known asset-suite comparison, and a fresh headless
save/load scenario. The reused default test binary passes 2,459 of 2,460 root
tests with only the established prerequisite-probe skip. All 1,654 original asset
hashes remain unchanged. See [qualification](audio-stream-callbacks-qualification.json).

## Delegation

GPT-6 Luna drafted the 20 fixture cases. Primary handled scope, new contracts,
production callback dispatch, raw root callers, integration and qualification.
The fixture draft was accepted after primary checked every owner/width mapping
and complete name-table coverage. Primary preserved the original dispatcher
default return of zero instead of the draft's new panic. The bounded integration
review found no missed callback consumer or convention mismatch. No extra
forced-GC fixture was needed for process-lifetime nonzero-size static keys. This
was a useful small implementation handoff; no subscription-savings claim is made.

## Result and review decisions

Three production cgo files were eliminated: exports, device registration and voice
registration. Selected project cgo files fall from 170 to 167, selected legacy
exports from 485 to 465. The 157 tracked headers now total 3,222 physical lines;
77 embedded production C bodies remain, and standalone production/test C stay zero.
External native dependency selection is unchanged in all production profiles.

The 16 fixture-only operations invoke the exact pre-existing Go owners. Four live
operations use static identities and signature-aware native dispatch; foreign
callbacks still use the original raw C fallbacks. Pointer arguments and pointer
results retain 32-bit words, numeric arguments retain signed int32 conversion,
and negative results retain their uint32 bit patterns in fixtures. No algorithm,
layout, media backend, assertion or frozen expectation changed. All three focused
profiles compiled and passed on the first attempt.

Baseline commit: `9f6b2b46`; reused production baseline: `3d47a346`. Local commands,
source maps, reviewed hashes, original-path repeats and final results are under
`build/port-audio-stream-callbacks/`. Completed original-path logs can be restored
with the gzip commands in `completed-log-archive.json`.
