# Gameplay text, notifications and player iteration

Ten connected routines now use native Go serialization/iteration: six text and
notification functions from4D9EB0 to4DA4F0, first/next player iteration4DA7C0/
4DA7F0, the byte-versus-wide classifier4100F0, and scripted spatial chat528AC0.
The original C baseline is recoverable at **bab248cf**.

## Implementation

Message construction, encoding, recipient filtering and fanout are Go. Two
production C varargs adapters retain the existing formatter (21 lines total),
then call Go serializers. Seven original entry points have Go-backed C exports;
these adapters require two additional serializer exports. Scripted chat's C ABI
and unused declaration were retired, and its Go caller now invokes Go directly.
Its C translation unit was removed. No C algorithm remains solely for tests.

Raw UTF16 code units are preserved, including isolated surrogates and Latin1
bytes. Original byte-length wrapping is preserved: a narrow508-character chat
message carries the first253 text bytes, rather than a fresh terminator. Existing
%u formatting prints values at/above0x80000000 as signed decimal; this inherited
formatter behavior remains unchanged. Notification input-record mutations, return
values and player/list traversal match the baseline. Player iteration uses the
existing native owner and retains the required C entry points.

## Baseline and coverage

**2,791 cases /11 repeated, locked hashes** cover every routine, full message bytes,
returns, guarded state and input mutations. Cases include byte/UTF16 boundaries,
embedded terminators, formatted substitutions, narrow/wide length limits,
unknown opcodes and full-width IDs, private-message length/suppression rules,
sparse player lists, active players without units, unit-code modes and signed
position truncation. All guard families are asserted. Suppression setup/restoration
includes inactive fixture players; optional larger records are used only for
long text, preserving all old record sizes and hashes.

Before the final16 long-chat cases, accumulated default ports passed84,243 captured
cases /1,020 groups plus contracts (331.687s wall). Affected reporting/control/text
server and highres checks passed131.178s/56.069s. The16 additional valid-buffer
cases repeat on default, with all prior text/reporting/control hashes unchanged.
All2,791 locked C cases pass together in6.926s. Production C was unchanged in the
baseline checkpoint.

## Native qualification

The first focused native text/reporting run passed in27.883s. Final accumulated
runs pass **84,259 captured cases /1,021 groups and applicable contracts**:

| Variant | Wall time | Selected root tests |
| --- | ---: | ---: |
| Standard | 321.386s | 614 |
| Server | 396.167s | 613 |
| High resolution | 335.635s | 614 |

Go discovery confirms these selections; actual root logs show successful nonempty
execution. Server excludes only TestFloorEligibility, explicitly guarded by
!server. All previous capture hashes are unchanged. This valid matrix also
supersedes the earlier reporting matrix whose regex accidentally selected zero
tests; see [the correction](GAMEPLAY_REPORTS.md) and commit713f067b.

Allthree production builds pass ELF386/SSE2, required-export and test-helper audits.
The retired scripted-chat symbol is absent. Full-suite results match all1,553
known failure entries exactly (15 passing /3 failing /32 skipped packages, with
zero added or removed failure entries). Fresh unchanged Xvfb/null-audio gameplay
passes34.664s with golden overrides disabled.

Physical C: **99,509 lines /147 files /zero reference C**, down282 lines and one
file. Local evidence is under build/port-gameplay-text: qualification.json,
variants.json, variant-selection-audit.json, root-execution-audit.json,
binary-verification.json, full-suite-comparison.json and native-first.log. The
scenario result is build/baseline/runs/gameplay-text-port/result.json.

Subsequent matrices must use tools/porting/run_tests.py, which trims the pattern,
rejects empty selections and checks JSON run/completion records for every root
test discovered by Go. Real positive/empty-selection checks and a controlled
partial-execution check validate the guard. The completed matrix above uses its
original logs plus the separate Go selection and nonempty-execution audits.

## Recovery and next batch

Tracked tests contain the capture hashes; recover original C from bab248cf.
Local captures may be SHA-verified gzip; capture-archive.json records the files.
Reporting's completed asset-copy deduplication has a restoration manifest and
preserves original assets/changed run files. Never commit the7z or extracted data.

Next is object-name/ID lookup and its16-entry cache:13 routines /436 C section
lines. The source/ABI audit and staged plan are in build/port-object-lookup.
