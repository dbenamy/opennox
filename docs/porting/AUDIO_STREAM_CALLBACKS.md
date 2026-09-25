# Internal audio-stream callback glue

Scope: retire 20 C exports in `legacy/client_audio_streams_exports.go`. Sixteen
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

Use four stable byte identities and a common known-callback dispatcher, with
separate void-discard and int32 wrappers. Unknown C callbacks retain their original
call conventions. The separate OnData/OnLoop/OnEnd/OnStop callbacks, driver API,
context initializer, sample storage and external audio libraries remain intact.
Add explicit compile-time assertions for all four callback offsets. Fixtures call
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
per profile. Conversion is not installed; existing captures are unchanged.
Production source is identical to qualified UI revision `3d47a346`, allowing its
production baseline to be reused. After conversion run focused contracts,
safe/static checks, production ABI checks, the known asset-suite comparison and a
fresh headless save/load scenario. Reuse the default test binary for a full corpus
check after conversion.

## Delegation

GPT-6 Luna owns an uninstalled draft of the 20 fixture cases. Primary owns scope,
new contracts, production callback dispatch, raw root callers, integration and
qualification. Primary also reviews every fixture argument/result conversion and
checks the complete name table before accepting the draft. Results pending.
