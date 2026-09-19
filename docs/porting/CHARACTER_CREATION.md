# Character creation — qualified Go conversion

C baseline **dd4debb6** is pushed, following qualified parent **6052e5b7**.
The baseline has **18,044 physical C lines /61 files /zero reference C**.
Qualified native source has **16,851 lines /59 files /zero reference C**
(**−1,193** physical C lines).
Selected scope:23 UI bodies /1,042 body lines in GAME3 and selclass/selcolor.
Adjacent13 configuration callbacks /217 lines are a separate deferred batch;
source proximity alone did not make them appearance-setting routines.

See [selection](character-creation-selection.json),
[literal references](character-creation-callers.json) and
[named global references](character-creation-global-callers.json).
The global audit excludes vardefs definitions from its outside-use summary, so
also follow the three modifier setters in vardefs.go, called by root modifiers.go.
Their storage moves with the preview while retaining the Go setter interfaces.

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
That identity applies to the committed C baseline; the native source is separate.

The fixtures explicitly own and restore modifierColorsOnce during palette
initialization. No C reference algorithms will be retained after conversion.

Further focused runs passed all initial-appearance and preview/material contracts
(10 roots /771 records). Construction-failure, admission-word and animation-cleanup
contracts also pass. File creation initially exposed fixture mistakes: resetting
the reused quickbar owner destroyed the entry widget; the UTF-16 expected-name
helper omitted its terminator. Both fixture assumptions were corrected before freezing; production
C was unchanged. The expanded run adds campaign working-save cleanup, failed
output creation, palette event routing and inclusive outside-click boundaries.

Frozen file records normalize only identified non-output data: temporary
path, clock, and the two unused tail padding bytes of the1280-byte C save record.
The original constructor clears bytes0..1277;1278..1279 are uninitialized stack
padding and are not serialized by the actual metadata section. Validate serialized
bytes and all defined record fields independently; do not freeze stack padding.

## Qualified native conversion

C baseline **dd4debb6** is committed and pushed. The 23 selected bodies now have
qualified Go implementations; original C bodies and the two dedicated
translation units are removed. Private state is grouped in a Go owner. Existing
Go modifier setter APIs now address that owner. Test bridges use native helpers;
the original root contracts and frozen expectations are unchanged.

Interface review corrected the initial plan: seven C exports remain, including
five callback addresses required by existing C-layout animation slots. Keep these
thin adapters rather than expanding this batch to change shared animation dispatch.
Sixteen private interfaces and 27 named C owners retire. Revisit animation callback
storage with its own connected batch; all callback logic here is already Go.

Preserve legacy name cleanup, byte-format filenames, previous-root palette lookup,
quickbar RNG consumption and preview layer/material ordering. Native save records
zero the two unused tail-padding bytes; defined/serialized fields stay compatible.

All 174 affected roots pass on default, server and highres without skips. All 142
captures /13,807 records match original C exactly, including the 16 new captures
/889 records. All 2,672 source fingerprints agree across tests and production.
Static memory checks pass. Fresh ELF32/386/SSE2/CGO client/highres/server builds
pass ABI inventories; the full suite matches the known result exactly (1,553
failure entries, 15 pass /3 fail /32 skip packages). Extended character-creation
and explicit save/load headless scenarios pass. See
[native qualification](character-creation-native-qualification.json).

Final evidence: build/port-character-creation/native-reviewed-{default,server,
highres,production}, static-native-reviewed.log. All sessions are joined. No
original root contracts or frozen captures were changed. An additional independent
construction-notification test covers the scenario regression described below.

Disk maintenance: removed 18,704 audited compiler-cache entries older than 12 hours,
with no compiler processes active and identity/access-time checks before deletion.
Reclaimed 4,924,575,701 bytes. Assets, archive, module cache, captures, saves and
production binaries remain intact; audit/apply is consumed.

The first native production scenario found a constructor-only regression despite
matching focused captures: class-event dispatch read argument one as a child
window before checking the event kind. The resource parser emits WindowNewChild
(event22) with a numeric ID, so construction dereferenced address601. Original C
returns before reading that argument. Native dispatch now checks event kind first;
an independent regression covers eight numeric child IDs without UI owners.
No frozen capture changed. All source-dependent gates passed after this fix.
The first three production builds/ABI and exact known full suite passed, but that
source revision is not qualified because its headless scenario failed.

The three completed native scenario copies (including the preserved failed run)
were hash-verified against source assets before removing duplicate files. Restore
manifests/helpers retain file hashes, modes and times. Original assets, saves,
captures and binaries remain. The audit/apply helper is consumed.

Follow-up cleanup: the pre-existing Go-to-C initializer adapter sub_4A5E90_A is
now unreferenced from C. Its Go hook is still used directly. Retire that thin
adapter with the next source batch, alongside the now-empty selclass header, so
cosmetic/adapter cleanup does not force another identical production cycle here.
No C character-creation algorithm remains.
