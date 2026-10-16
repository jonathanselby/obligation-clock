# Obligation Clock — Generated Engineering Design Document

**Evidence commit:** `5ddaf08cc7217595a29fafbadb6eeab26eb8433c`  
**Method:** code-derived review using the `codebase-to-edd` evidence discipline  
**Generated:** 2026-10-17

This document describes the implementation at the evidence commit. The PRD and hand-written EDD are comparison targets, not runtime evidence.

## 1. Executive description

**KNOWN** — `obligation-clock` is a deterministic Go CLI with one material core package, `internal/clock`.

**KNOWN** — it accepts an RFC3339 trigger, duration/unit, explicit caller-supplied calendar and lifecycle label, then returns one deadline or a visible validation/calculation error.

**KNOWN** — it contains no legal applicability logic, regulator-specific periods, holiday feed, persistence or notification.

## 2. Evidence inspected

- `cmd/obligation-clock/main.go` — CLI and exit contract;
- `internal/clock/clock.go` — validation and arithmetic;
- `internal/clock/clock_test.go` — executable behavioural evidence;
- `.github/workflows/test.yml` — engineering gate;
- `examples/bafflewick/` — executable example;
- `go.mod`, `VERSION`, `CHANGELOG.md` — dependency/release boundary.

## 3. Runtime structure

```text
request
  |
  v
CLI -> Calculate
       |
       +-- validate timezone and trigger offset
       +-- validate weekend/holidays/adjustment/lifecycle
       +-- hours | calendar_days | business_days
       +-- reject impossible local DST wall time
       +-- optional next/previous business-day adjustment
       |
       v
     result
```

## 4. Time semantics

**KNOWN** — `hours` uses elapsed duration arithmetic and therefore may change the local wall-clock hour across a DST transition.

**KNOWN** — `calendar_days` uses local date arithmetic and then checks that hour/minute/second/nanosecond still equal the original wall-clock fields. A spring-forward gap is refused rather than silently normalised.

**KNOWN** — `business_days` advances by local calendar dates and counts only dates outside the validated weekend and holiday sets.

## 5. Timezone contract

**KNOWN** — the named IANA timezone is required and loaded from the Go runtime.

**KNOWN** — the explicit offset in the RFC3339 trigger must equal the named timezone's UTC offset at that instant. Contradictory input is refused.

## 6. Calendar contract

**KNOWN** — no weekend input means Saturday/Sunday.

**KNOWN** — custom weekend names must be valid and unique.

**KNOWN** — a seven-day weekend calendar is refused before iteration, preventing a non-terminating business-day or adjustment loop.

**KNOWN** — holiday values must be unique valid `YYYY-MM-DD` dates in the selected timezone.

## 7. Adjustment and lifecycle

**KNOWN** — adjustment supports none/empty, `next_business_day` and `previous_business_day` using the same validated calendar.

**KNOWN** — `Adjusted` is true only when the calculated date moved.

**KNOWN** — lifecycle accepts only `provisional` or `firm` and does not affect arithmetic. `TestLifecycleDoesNotChangeDeadline` proves this.

## 8. Failure model

**KNOWN** — process exits are:

```text
0 success
1 deterministic calculation/calendar-contract refusal
2 usage error
3 malformed JSON or file I/O
```

**KNOWN** — there is no retry loop, fallback calendar or silent input repair.

## 9. State and dependencies

**KNOWN** — the core is stateless and side-effect free.

**KNOWN** — Go 1.26 is declared and no third-party module dependency exists.

**KNOWN** — no database, scheduler, network client or background worker exists.

## 10. Executable evidence

Named tests cover elapsed hours, DST elapsed-hour behaviour, calendar days through DST, DST-gap refusal, holidays, custom weekends, next/previous adjustment, lifecycle invariance, invalid calendar inputs and trigger offset/timezone mismatch.

**KNOWN** — CI runs Go 1.26 formatting with diagnostic diffs, vet, race tests, native build, exact version check, Linux ARM64 cross-build, Bafflewick CLI acceptance and the named calendar-boundary suite.

## 11. Bafflewick path

**KNOWN** — the fixture starts Friday 16 October 2026 at 17:00 in Sydney and adds three business days using the default weekend. Expected output is Wednesday 21 October 2026 at 17:00 +11:00.

**KNOWN** — the fixture makes no claim about any real legal deadline.

## 12. Security/correctness boundary

**KNOWN** — the component calculates from supplied data only; it cannot establish whether an obligation applies or whether a trigger/calender is legally correct.

**KNOWN** — the implementation favours visible refusal over plausible-looking correction when inputs are contradictory or a local wall time does not exist.

## 13. Operational behaviour

**KNOWN** — synchronous short-lived process; no cache, durable state, telemetry framework or network dependency.

**UNKNOWN** — production throughput/latency/memory on the Garage constrained-hardware target. ARM64 CI is compilation proof, not runtime benchmarking.

## 14. Known limits

- no business-hours or partial-day model;
- no holiday/jurisdiction database;
- no automatic extensions or waivers;
- no persistence or notification;
- no leap-second-specific model;
- fall-back ambiguity is represented by the explicit RFC3339 instant supplied by the caller.

## 15. Unknowns

**UNKNOWN** — completeness of production calendars, host timezone-data currency and constrained-hardware runtime performance.

## 16. Recovered architecture

```text
CALLER
  |
  v
cmd/obligation-clock
  |
  v
internal/clock.Calculate
  |
  v
validated deterministic deadline
```

## 17. Evidence conclusion

**KNOWN** — the source implements the declared one-job boundary and the material failure cases are executable in tests/CI. No hidden legal rule or external service is present in the evidence snapshot.
