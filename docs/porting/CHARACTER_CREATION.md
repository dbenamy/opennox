# Character creation — qualified C baseline

Qualified parent **6052e5b7** is pushed. Production remains original C for this
batch; the qualified total is **18,044 physical lines /61 files /zero reference C**.
Selected scope:23 UI bodies /1,042 body lines in GAME3 and selclass/selcolor.
Adjacent13 configuration callbacks /217 lines are a separate deferred batch;
source proximity alone did not make them appearance-setting routines.

See [selection](character-creation-selection.json),
[literal references](character-creation-callers.json) and
[named global references](character-creation-global-callers.json).
The global audit excludes vardefs definitions from its outside-use summary, so
also follow the three modifier setters in vardefs.go, called by root modifiers.go.
Their storage can move with the preview while retaining the Go setter interfaces.

## Qualified C baseline

All 173 affected test roots pass on default, server and highres without skips.
All 142 captures /13,807 records match across targets, and all 2,679 source
fingerprints agree. Two independent focused processes produced identical results;
the 16 new captures /889 records are frozen. Static memory checks pass. See
[qualification](character-creation-c-qualification.json) and
[frozen captures](character-creation-captures.json).

- Names: empty/all-whitespace preservation, trimming, first non-space reserved
  character replacement, preservation of later reserved characters, UTF-16 text,
  and maximum entry length. Preserve the original first-character-only rule.
- Palette lookup: every entry in all three32-color banks, absent colors, first-match
  duplicates, packed bank/index words and associated enable/toggle side effects.
  Bank1 absent colors select index9; other absent colors retain index32.
- Palette menu: all32 packed cells, selected indices0/1/31, cancellation32/DEAD/FFFF,
  missing target, hidden state, and return identity. Non-rendering packing tests
  include high bank words without indexing the palette outside its valid range.
- Swatches: every valid bank/index, full software framebuffer captures and an
  independent rectangle extent and RGB555 pixel contract.

The complete focused baseline passes16 roots /889 records. Added coverage includes
all appearance overrides and class/name combinations, default/current initialization,
class selection and shading, empty/single/multiple/255-row quickbar catalogs with
actual RNG consumption, preview materials and clipping, animation destruction,
constructor failures, admission flags, palette event dispatch and outside clicks.

Actual encrypted character files cover all three classes, override masks, slots
0/1/99,100 occupied names, blocked output paths and campaign working-save cleanup.
Checks include the caller's record fields, decrypted section framing/payloads,
working-directory restoration, campaign map names and filename generation. The
legacy byte formatter sees only the first ASCII character of a UTF-16 name; the
fixture confirms H00.plr for Hero. Preserve that behavior during this conversion.

The extended shipped-resource headless scenario and independent repeat pass,
including every class, preview, skin palette/selection, override toggle, outside
click, animation transitions and subsequent gameplay. Palette/selection images
were visually inspected. Constructor failure tests use actual GUI owners with
injected missing-resource/animation-allocation outcomes; normal construction uses
the actual shipped resources in the headless run.

Original-C production is identical to6052e5b7 apart from20 build-tag-only fixture
files. Parent binary hashes were reverified; its production/ABI/full-suite/save-load
evidence can be reused. See [identity](character-creation-production-identity.json).
Native conversion must qualify fresh production and compare the extended scenario.
No production source has changed. Native conversion is next.

The fixtures explicitly own and restore modifierColorsOnce during palette
initialization. No C reference algorithms will be retained after conversion.

Further focused runs passed all initial-appearance and preview/material contracts
(10 roots /771 records). Construction-failure, admission-word and animation-cleanup
contracts also pass. File creation initially exposed fixture mistakes: resetting
the reused quickbar owner destroyed the entry widget; the UTF-16 expected-name
helper omitted its terminator. Both fixture assumptions were corrected; production
C remains unchanged. The expanded run adds campaign working-save cleanup, failed
output creation, palette event routing and inclusive outside-click boundaries.

Frozen file records normalize only identified non-output data: temporary
path, clock, and the two unused tail padding bytes of the1280-byte C save record.
The original constructor clears bytes0..1277;1278..1279 are uninitialized stack
padding and are not serialized by the actual metadata section. Validate serialized
bytes and all defined record fields independently; do not freeze stack padding.
