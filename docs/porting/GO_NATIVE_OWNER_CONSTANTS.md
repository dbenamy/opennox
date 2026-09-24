# Native owners and constants

## Baseline and scope

Baseline source is qualified commit `b505e0ab`. Twelve reviewed source files move
browser, interaction-message, damage and server-options callers to existing Go
owners, preserving temporary string ownership, signed narrowing and return values.
The proposed change removes six production and one test C imports and ten private
Go adapters. All 1,179 actual C exports and 78 embedded C callback bodies stay.
Expected selected production cgo count: 236 to 230; qualification is pending.

The baseline records all 3,052 source fingerprints. Fresh processes using the
qualified binaries pass exactly 195 default, 193 server and 195 highres selected
roots, with no skips. Assertions and frozen expectations are unchanged.
See [baseline](go-native-owner-constants-baseline.json),
[batch commands](go-native-owner-constants-batch.json) and
[test selection](go-native-owner-constants-tests.txt).

## Compatibility review and test plan

The legacy package enables NOX_HIGH_RES for every profile; preserve its compiled
protocol value 0x000F039A, rather than substituting the root package's profile
constant. A 32-bit original-C probe verifies that value, console red=6,
engine godmode=32, int width=4 and short width=2.

Monster blocking policy is extracted without changing its retained C export.
The original AI fixture covers 512 cases and 28 operations, including this policy.
The selected suites include all damage, object-state, AI-state, roster, browser,
interaction, minimap, server-options and resource-definition family roots.
Server numeric edits use the already qualified resourceAtoi parser, with existing
UI clamp/prefix and independent signed-overflow contracts.

Keep the public browser record C() convenience accessor, returning unsafe.Pointer
with the same address/nil behavior. Its C-specific return type is intentionally
removed. Record layout, allocation and the existing 169-byte copy are unchanged.

Luna drafted five bounded files and reviewed seven primary files. Primary expanded
its initially incomplete 25-root test proposal and reviewed exact reconstruction,
whole-source references and preamble dependencies. Local evidence lives under
`build/port-go-native-owner-constants/`; drafts alone are not qualification.

After conversion, run the exact focused roots in all three profiles, safe build
and static checks, three fresh production builds and ABI checks, the exact known
full-suite baseline, and headless character creation/save/load/resume. Verify
source fingerprints, original asset hashes and dependency counts. This localized
caller batch does not change shared layouts or allocation, so another complete
35-minute root corpus is reserved for a meaningful shared-infrastructure boundary.
