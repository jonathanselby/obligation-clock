# Obligation Clock — Design / Implementation Comparison

**Implementation evidence:** `dcacea90dd09106cb6f432dee488f18cce3d9b22`  
**Generated:** 2026-10-17

| Design area | Hand-written intent | Recovered implementation | Assessment |
|---|---|---|---|
| one-job boundary | deterministic deadline arithmetic only | one CLI plus one material `internal/clock` package; no legal applicability or notification logic | MATCH |
| elapsed hours | elapsed time, not wall-clock preservation | `time.Add(duration * time.Hour)`; DST test proves elapsed semantics | MATCH |
| calendar days | preserve local wall-clock time | `AddDate` plus explicit wall-time equality check | MATCH |
| impossible DST time | fail visibly | spring-forward gap is detected and refused | MATCH |
| business days | caller-supplied weekend/holiday calendar | local-date iteration against validated weekend/holiday sets | MATCH |
| invalid calendar | refuse, do not loop or infer | invalid weekdays, duplicates, bad holidays and seven-day weekend are rejected | MATCH |
| timezone contract | explicit IANA timezone | trigger offset must match named timezone at trigger instant | MATCH |
| adjustment | next/previous business day | same validated calendar semantics used by adjustment | MATCH |
| lifecycle | provisional/firm metadata only | test proves lifecycle cannot change deadline | MATCH |
| legal boundary | no legal applicability rules | none present in source or dependencies | MATCH |
| dependencies | Go standard library only | no third-party modules | MATCH |
| CLI failure visibility | stable exit contract | calculation refusal 1; usage 2; malformed input/I/O 3 | MATCH |
| CI | format, vet, race, version, ARM64, acceptance | all encoded in `.github/workflows/test.yml` | MATCH |
| code-derived evidence | recovered design pinned to source | `generated/EDD.md` pins `dcacea90dd09106cb6f432dee488f18cce3d9b22` | MATCH |

## Material corrections made before publication

The draft implementation allowed a calendar containing all seven weekdays as weekend. Business-day calculation or business-day adjustment could then iterate forever. The published V0.1 validates and refuses that calendar before calculation.

The draft documentation also implied unconditional wall-clock preservation for calendar-day arithmetic. The published contract is more precise: wall time is preserved **when that local wall time exists**; calculations landing in a DST spring-forward gap are refused rather than silently normalised.

The first GitHub publication also exposed a Go 1.26 formatting difference in the CLI import block. The evidence commit is the corrected, `gofmt`-conformant source snapshot.

## Result

**MATCH after correction.** The final source implements the intended one-job component with stronger failure semantics than the initial draft. No unexplained design/implementation mismatch remains at the evidence commit.

The comparison does not claim legal correctness of caller-supplied obligations or calendars, and it does not claim constrained-hardware runtime performance; those remain outside V0.1 evidence.
