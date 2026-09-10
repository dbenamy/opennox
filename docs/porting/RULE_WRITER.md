# Rule-writer baseline and online bug — 2026-09-10

The loader port is completed at `0a2766f7`. The next candidate, C rule writer
57AAA0, has an existing online-mode memory-layout bug. No writer production code
has been changed. Production C remains **141,455 physical lines**, 153 files,
zero test-reference C.

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

## Reproduction and next decision

An asset-free diagnostic is preserved at
[writer_online_probe_test.go.txt](probes/writer_online_probe_test.go.txt). Copy it
to src/rules_writer_probe_test.go, load the documented 386 baseline environment,
then run from src:

```
go test -tags porttest -count=1 -v -run '^TestRulesWriterProbe$' .
```

Remove that temporary test afterward. It is a diagnostic that prints the output
matrix, not an assertion that the broken online output should be preserved.
Exact garbage-dependent output can change with stack history or build layout.
The fixture isolates files, server registries, settings, flags, list ownership,
handle arena and current rule context. Local observed output is in
build/port-rules/writer-probe.log; no raw environment logs are committed.

The user has been asked whether to correct this while porting the writer
(recommended), using two explicit Settings2 buffers and independent intended-
filter tests, or postpone the writer and port another section. The online
behavior change is pending that answer. Offline baseline work is independent.

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
writer cases, pass on 386 default, server and highres. Production files are
unchanged since the fully built/gameplay-validated loader commit; this checkpoint
adds fixtures, baseline tests and recovery documentation only.
