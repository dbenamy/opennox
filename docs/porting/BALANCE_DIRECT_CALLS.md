# Direct Go balance lookups

Twenty-seven production call sites now call the existing Go balance getter
functions directly: 19 scalar reads and 8 indexed reads across 14 files. The
lookup implementation, interned string arguments, call ordering, float casts and
surrounding arithmetic are unchanged. Indexed arguments retain signed32
normalization; signed int8 levels stay sign-extended, and unsigned Vampirism
power subtraction still happens before signed conversion. Redundant casts of
literal `4` and already-narrow values are simplified.

This is removal of Go→C→Go round trips, not replacement of balance algorithms.
The qualified target remains 386/SSE2; no new 64-bit compatibility is claimed.
All touched files still use other C interfaces and retain their cgo imports.
The two exported C entry points remain in this checkpoint. A following
whole-repository reachability audit will decide which book/balance exports can
be retired; test-only references must not keep dead production interfaces alive.
No new handwritten C algorithm bodies were introduced for testing.

## Qualification

Baseline `8409eff3` is committed/pushed; [baseline record](balance-direct-calls-c-qualification.json).
The 109-root affected selection passed default/server/highres before conversion,
with no skips. Preexisting production source and all four binaries matched the
preceding qualified book checkpoint, whose production evidence was reused.

New independent contracts run 104 mode/key/index cases, checking four exact
float64 results each through C and direct Go routes. They cover scalar/array
values, nonzero fractions, negatives, signed zero, a finite value beyond float32
range, the next float64 above 1, missing/empty/case-folded names, empty arrays,
negative/end/int32 extreme indices, global fallback and Arena/Solo overrides.
The server fixture overlays actual balance data; it does not replace lookup
behavior. Tests restore flags, the server owner and its prior balance file.

The gameplay selection covers attack, damage, equipment pickup, object collision,
player controls, quest lives, projectile collision, spell effects and spell
lifecycle. Existing source-frozen hashes remain unchanged. Some historical
fixtures do not prove every lookup branch individually (notably quest starting
lives, ForceOfNature staff limits and a zero-valued melee percentage). The new
numeric contract and static review cover the common getter boundary; this is not
an exhaustive claim about every caller state.

After conversion, all 109 roots pass in each of default/server/highres, without
skips or changed frozen hashes. Safe build and static checks, three fresh
production builds and ABI checks, headless character creation and save/load pass.
The full suite matches the existing 304 failure events, with 17 passing, two
failing and 32 skipped packages. All four binaries retain both C exports and omit
the redundant C call bridges; no porttest symbols are present. See the
[conversion record](balance-direct-calls-qualification.json).

## Review and recovery

Luna drafted the 27-call patch, mapped existing consumer coverage, and supplied
the numeric fixture. Primary verified the source/header signatures and local
balance implementation, added precision/tag-precedence/name cases, simplified
redundant casts and ran the gates. Luna's independent review found no material
issue, including string ownership, integer underflow/sign extension, float
conversion sites and fixture cleanup. No measured game-speed claim is made.

Artifacts: `build/port-balance-direct`; initial drafts and coverage review are
under `build/port-book-direct/balance-*`. The batch manifest and qualification
records are committed; local logs remain reproducible working evidence.

For storage, 33 completed historical test logs were losslessly gzip-archived:
1,112,212,170 raw bytes became 67,649,331 compressed bytes. Exact stat/hash and
host open-file checks preceded compression; every decompressed byte stream was
verified before retiring the raw file. The reviewed inventory and restoration
metadata are `old-log-archive-plan.json` and `old-log-archive-record.json` under
this artifact directory. No assets, fixtures, binaries or cache were removed.

Standalone production and test-reference C remain **zero files / zero lines**;
production C preamble bodies remain **79**. Generated bridges and external
libraries are outside that metric.


A later broad scalar inventory was not efficient: several helper Python scans
repeated source searches per exported symbol, ran for minutes and competed with
the default qualification. Their process IDs were10895,10940,11150,11181. Primary
verified ownership through host parent commands; by the final check three had
exited and the remaining11181 was stopped with SIGTERM. A subsequent host check
confirmed no such scans remained. The default109-root qualification passed, but
its wall time is not a performance comparison. PORT.md now requires single-pass
inventories, short subprocess timeouts, retaining/joining tool session IDs and
primary host-namespace process verification. Broad inventory design stays with
the primary; bounded implementation/test drafts remain suitable for Luna.
