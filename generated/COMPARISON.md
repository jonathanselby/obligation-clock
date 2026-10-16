# Obligation Clock — Design / Implementation Comparison

**Implementation evidence:** `5ddaf08cc7217595a29fafbadb6eeab26eb8433c`  
**Generated:** 2026-10-17

| Design area | Hand-written intent | Recovered implementation | Assessment |
|---|---|---|---|
| one-job boundary | deterministic deadline arithmetic only | one CLI plus one material `internal/clock` package | MATCH |
| elapsed hours | elapsed time | `time.Add` with DST-crossing test | MATCH |
| calendar days | preserve valid local wall time | `AddDate` plus explicit wall-time check | MATCH |
| DST gaps | fail visibly | impossible spring-forward wall time refused | MATCH |
| business days | caller-supplied calendar | validated weekend/holiday date iteration | MATCH |
| invalid calendar | refuse | invalid/duplicate weekend/holiday values and seven-day weekend refused | MATCH |
| timezone contract | explicit IANA timezone | trigger offset must match named timezone at trigger instant | MATCH |
| adjustment | next/previous business day | same validated calendar reused | MATCH |
| lifecycle | metadata only | test proves no arithmetic effect | MATCH |
| legal boundary | no applicability rules | none present in code/dependencies | MATCH |
| dependencies | Go standard library only | no third-party modules | MATCH |
| failure visibility | stable process contract | refusal 1; usage 2; malformed input/I/O 3 | MATCH |
| CI | format, vet, race, version, ARM64, acceptance | all present; format step prints diffs before failure | MATCH |
| code-derived evidence | pin recovered EDD to source | generated evidence pins `5ddaf08cc7217595a29fafbadb6eeab26eb8433c` | MATCH |

## Corrections made before release

The draft permitted all seven weekdays to be configured as weekend, creating a potential infinite loop. The published source rejects that calendar.

The draft documentation overstated wall-clock preservation. The implementation now explicitly refuses calendar arithmetic that lands in a nonexistent DST wall time.

The GitHub Go 1.26 gate also exposed a formatter change in one long test literal. The evidence commit contains the exact Go 1.26 formatted source and the CI format gate now prints a useful diff when it fails.

## Result

**MATCH after correction.** The implementation and written V0.1 contract agree, and no unexplained design/implementation mismatch remains at the evidence commit.
