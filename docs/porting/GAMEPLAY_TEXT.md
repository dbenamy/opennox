# Gameplay text and player iteration candidate

Ten connected routines: the six notification functions from4D9EB0 to4DA4F0,
first/next player iteration4DA7C0/4DA7F0, byte-versus-wide classifier4100F0,
and scripted spatial chat528AC0 (entire server__system__cscrfunc.c).
Roughly300 physical C lines, including one translation unit.
Reuse current guarded player/message fixture with optional text extension,
action range1900+, C dispatch only in baseline. Keep previous hashes unchanged.

Baseline cases: single/all-player narrow and wide text, raw UTF16 code units
including0xff/0x100/surrogates, embedded terminators, zero/1/254/255 length
boundaries within original buffers, flag combinations, real formatter integer/
wide substitutions, empty/sparse/active-without-unit player iteration, null and
non-player guards, all information opcodes including unknown values and input
record mutations (6/10/11-byte branches), private message name lengths0/1/48/49
and player suppression, chat position signed/truncation/word boundaries and
unit-code modes. Capture full direct-message bytes, returns and guarded state.
Do not send any external chat; use only owned fixture queues.

Thin production C varargs adapters remain for the two formatting entry points;
reuse existing C formatting routine and move serialization/fanout into Go.
Retire obsolete Go-to-C helper calls. Classifier remains exported for the client
say routine, which is outside this batch. Use existing native player iteration
via C exports, after comparing the legacy behavior on its valid input domain.
Finish the current reporting qualification/commit before changing source here.

Buffer-domain audit: formatted messages use1032-byte storage (separate formatting
area at520), allowing255 UTF16 units plus terminator. Scripted chat has only520
bytes total: wide text cases must stop at253 units (11-byte header plus508-byte
payload); narrow255 units is within storage and exercises the byte length wrap.
Keep raw UTF16 slices in the new implementation so isolated surrogate code units
and Latin1 bytes are preserved without UTF8 replacement. Original format handling
supports %s as UTF16 and %S as byte text; test these actual conventions.

Baseline review findings: existing nox_vsnwprintf handles %u via nox_itoa(int),
so values at/above0x80000000 print signed decimal. Preserve this existing formatter
behavior, and correct the expected fixture values rather than changing it here.
Suppression masks must be applied/restored to all three fixture-owned player
objects, including inactive players: direct private notifications can target an
inactive record even though broadcast iteration skips it.

## Original-C capture checkpoint

All2,775 cases /ten groups pass their message-byte assertions and repeat exactly;
the hashes are now locked in gameplay_text_porttest_test.go. Captures cover allten
routines, including raw return values, input-record mutations and guarded state.
The final chat-code expectations use actual Immobile0x400000 and
ClientPersist0x20000000 flags; an initial test used the wrong latter mask, corrected
before capturing. No production code has changed. The default accumulated baseline passed84,243 cases /1,020 groups and
contracts (331.687s wall); focused affected server/highres qualification is
running. Long-chat cases below will extend this before the baseline commit.

The audit finds302 C section/file lines in scope. Scripted chat's only outside C
reference is an unused declaration in server__script__builtin.c; remove that with
the Go-only C entry point. Nine original C names remain required (two varargs C
adapters, seven Go exports). The adapters need two new Go serializer exports.
Expected reduction is282 physical C lines with one fewer C translation unit;
measure after conversion. Native drafts are only in build/port-gameplay-text.

Latest completed reporting run asset duplicates have a restoration manifest:
build/baseline/runs/gameplay-reports-port/deduplicated-assets.json. Hash-identical
copies only were removed (556,358,986bytes); original assets and changed run files
are preserved. Restore with the existing deduplicate-run-assets.py --restore.

Additional review before conversion: scripted narrow chat permits up to508 text
units in its520-byte message storage, even though formatted text has a255-unit
limit. Add long narrow cases256/257/507/508 before the baseline commit. The byte
length wraps, and wrapped payloads contain the first count bytes of the original
text, not a newly placed terminator. This is a valid-buffer case and must be
preserved. Extend only long fixture text allocations; existing520-byte records
and their hashes should remain unchanged. Native draft already accounts for it;
production C remains unchanged while current qualification runs.

## Qualified baseline

The final baseline contains **2,791 cases /11 repeated capture hashes**, all
locked and passing together (6.926s). Before the final16 long-chat cases, full
accumulated default ports passed84,243 cases /1,020 groups and contracts
(331.687s); affected reporting/control/text server and highres checks passed
131.178s/56.069s. The16 extra cases repeat on default, with old text/reporting/
control hashes unchanged. Existing records retain their original sizes; only
long text allocates larger guarded storage. All guard families are asserted.
Run all84,259 cases /1,021 groups on every variant after conversion.

There are no production C changes in this baseline. Physical C remains99,791
lines /148 files /zero reference C. Every original function has positive coverage.
The source and ABI audit distinguishes the unused scripted-chat C declaration
from real callers. Native drafts remain local and unapplied until this baseline
is committed and pushed. Local baseline-qualification.json records the check scope.
