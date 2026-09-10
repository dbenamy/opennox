# Rule-writer baseline and online bug — 2026-09-10

The rule writer 57AAA0 now runs Go after offline baseline `36d8fa66`. Its
online memory-layout bug is corrected with the user’s explicit approval. The
original-C evidence below explains the intended behavior change; conversion
and final counts are recorded at the end.

## Concrete finding

The writer declares `v19[24]`, `v20[36]`, `v21[24]`, `v22[36]` as separate local
arrays. Decompiled offset comments imply two contiguous 60-byte Settings2 values.
The current GCC build does not preserve that implied layout. With B representing
ESP after the function prologue, disassembly of the pre-loader-port 386 binary
`build/port-waypoint/bin/opennox` shows:

| Operation | Actual buffer address |
| --- | --- |
| Internet settings load (selection4), v21 | B + 0x3c |
| User/map settings load (selection3), v19 | B + 0x24 |
| First online spell-mask check, v22 | B + 0x78 |
| Second online spell-mask check, v20 | B + 0x54 |

Settings spell data starts at offset24. The actual loader destinations therefore
have spell data at B+0x54 and B+0x3c, respectively. They overlap, and the first
online check reads uninitialized storage at B+0x78. Source comments establish
intention, not actual compiled placement. Any conversion must be reviewed against
these compiled addresses and observable behavior.

The intended filtering condition visible in the C expression is:
`!online || internetHas(spell) || !userOrMapHas(spell)`.
A disabled valid rule spell should be emitted when this condition is true.
Two explicit Settings2 temporaries would implement it without memory overlap.

The isolated current-C writer probe used DEATHMATCH, disabled BLINK and FIREBALL,
and independently toggled both spells in internet.rul and user.rul. Offline
writing consistently emitted both disabled spells. In online mode, with neither
input file disabling either spell, the writer emitted only the header, omitting
both restrictions. Repeated identical combinations also varied: the combination
with both input files disabling both spells emitted only the header in trials0/1,
but emitted both spell directives in trial2. This is observable output instability,
not merely a theoretical C-language concern.

## Historical reproduction

An asset-free diagnostic is preserved at
[writer_online_probe_test.go.txt](probes/writer_online_probe_test.go.txt). Copy it
to src/rules_writer_probe_test.go, load the documented 386 baseline environment,
then run from src at original-C checkpoint `36d8fa66`:

```
go test -tags porttest -count=1 -v -run '^TestRulesWriterProbe$' .
```

Remove that temporary test afterward. It is a diagnostic that prints the output
matrix, not an assertion that the broken online output should be preserved.
Exact garbage-dependent output can change with stack history or build layout.
The fixture isolates files, server registries, settings, flags, list ownership,
handle arena and current rule context. Local observed output is in
build/port-rules/writer-probe.log; no raw environment logs are committed.

The user explicitly approved correcting this during the writer port. The fix
uses two complete Settings2 values and independent intended-filter tests. The
garbage-dependent online output is not a compatibility target.

## Offline baseline

TestRulesWriterOfflineABI checks 80 combinations of flags, restriction patterns
and optional rejected lists, plus a failed-create case. It compares exact output
bytes and order, all settings bytes and surrounding guards, unchanged context,
input files, list links, header table and registered file-handle count. Spell
output is ordered by ID; armor and weapon output is ordered by individual bits.
The existing C bridges return non-nil empty strings for unknown equipment bits,
so the original writer emits empty-name directives for disabled unassigned bits;
that defined behavior is represented in the baseline. Rejected UTF-16 text is
narrowed byte-wise, including a low-byte NUL that can suppress the following LF.

All accumulated protection/network/waypoint/rule tests, including the 81 offline
writer cases, pass on 386 default, server and highres. At baseline checkpoint `36d8fa66`, production files were
unchanged since the fully built/gameplay-validated loader commit; that checkpoint
added fixtures, baseline tests and recovery documentation only.

## Go conversion and approved fix

The user approved the online fix explicitly. `legacy/rules_writer.go` replaces
57AAA0 with Go, retaining its live C bridge; the root Go caller now enters the
native writer directly. The signed low-byte Field52 gate and historical zero
return on success or failed creation remain unchanged. The writer creates the
output before loading comparison rules, preserving truncation behavior when
saving directly over user.rul. Rejected nodes remain owned by the caller's C list.

Two complete independent Settings2 values replace the split arrays. Online
filtering is tested using the intended condition rather than the old invalid
stack reads. 576 cases independently vary internet restrictions, local
restrictions, current restrictions, user/map/empty-user selection and repetition;
they exercise both the live C ABI and direct Go caller. Eight additional cases
cover missing inputs, saving over user.rul, blocked saves and failed creation
through both callers. The 81 original-C offline cases remain unchanged.

Production C after conversion: **141,351 physical lines (−104)** in 153 files,
zero test-reference C. Local artifacts: `build/port-writer/`.


Final validation: all accumulated ABI tests pass on386 default/server/highres;
all three binaries build; writer-port passes both preserved gameplay screenshots
with overrides disabled. Full-suite results exactly match the loader milestone:
15 passing,3 known failing,32 skipped/no-test packages and the same1,553 failure
entries. Symbol inspection confirms the live57AAA0 bridge and native ruleWrite.
No C writer implementation is retained solely for tests.
