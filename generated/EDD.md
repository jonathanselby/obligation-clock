# Obligation Clock — Generated Engineering Design Document

**Evidence commit:** `6086d377153afffffbd28553126e91273f65ba44`  
**Method:** code-derived review using the `codebase-to-edd` evidence discipline  
**Generated:** 2026-10-17

This document describes the implementation at the evidence commit. `PRD.md` and the hand-written `EDD.md` are comparison targets, not evidence for runtime behaviour.

## 1. Executive system description

**KNOWN** — `obligation-clock` is a deterministic Go CLI with one material core package: `internal/clock`.

**KNOWN** — the core accepts an RFC3339 trigger, explicit duration/unit, caller-supplied local calendar and lifecycle label, then returns one calculated deadline or a visible validation/calculation error.

**KNOWN** — it does not contain legal applicability rules, regulator-specific periods, a holiday feed, persistence or notification.

## 2. Evidence inspected

- `cmd/obligation-clock/main.go` — process boundary and exit mapping;
- `internal/clock/clock.go` — validation and arithmetic;
- `internal/clock/clock_test.go` — executable behavioural evidence;
- `.github/workflows/test.yml` — engineering/release gate;
- `examples/bafflewick/` — representative acceptance fixture;
- `go.mod` — dependency and language baseline;
- `VERSION` / `CHANGELOG.md` — release metadata.

## 3. Runtime structure

```text
JSON request
    |
    v
cmd/obligation-clock
    |
    v
internal/clock.Calculate
    |
    +-- validate IANA timezone
    +-- validate trigger offset against timezone at instant
    +-- validate weekend and holiday calendar
    +-- apply unit-specific arithmetic
    +-- reject impossible DST wall-time result
    +-- apply optional business-day adjustment
    |
    v
JSON Result
```

**KNOWN** — material calculation logic is outside the CLI.

## 4. Input contract

**KNOWN** — the request carries trigger, duration, unit, calendar, lifecycle and optional adjustment.

**KNOWN** — the named timezone is loaded from the Go runtime timezone database. The trigger's explicit UTC offset must equal the named timezone's offset at the trigger instant. Contradictory trigger/timezone input is refused.

**KNOWN** — negative duration is refused.

## 5. Unit semantics

### Hours

**KNOWN** — `hours` uses elapsed duration arithmetic (`time.Add(duration * time.Hour)`). It intentionally does not preserve local wall-clock hour across DST changes.

### Calendar days

**KNOWN** — `calendar_days` uses local calendar-date arithmetic (`AddDate`). After arithmetic, the implementation checks that the original local wall-clock fields still exist. If a spring-forward gap caused Go to normalise to a different wall time, the component returns an error.

### Business days

**KNOWN** — `business_days` advances by validated local calendar days and counts dates not present in the weekend set and holiday set.

**KNOWN** — no elapsed-24-hour assumption is used for business-day counting.

## 6. Calendar validation

**KNOWN** — when no weekend is provided, Saturday/Sunday is used.

**KNOWN** — custom weekend names must be recognised Go weekdays and unique.

**KNOWN** — all seven days cannot be configured as weekend. This guard prevents an otherwise non-terminating business-day/adjustment loop.

**KNOWN** — holidays must be unique, valid `YYYY-MM-DD` dates in the selected timezone.

## 7. Adjustment semantics

**KNOWN** — supported values are none/empty, `next_business_day` and `previous_business_day`.

**KNOWN** — adjustment uses the same validated calendar and local-day stepping as business-day counting.

**KNOWN** — `Result.Adjusted` is true only when the initially calculated date was moved.

## 8. Lifecycle semantics

**KNOWN** — only `provisional` and `firm` are accepted.

**KNOWN** — lifecycle is result metadata only. `TestLifecycleDoesNotChangeDeadline` proves that changing only lifecycle does not alter deadline arithmetic.

## 9. CLI and failure model

**KNOWN** — `-version` prints `obligation-clock 0.1.0`.

**KNOWN** — the CLI uses:

```text
0 success
1 deterministic calculation/calendar-contract refusal
2 usage error
3 malformed JSON or file I/O
```

**KNOWN** — there is no retry loop, hidden fallback calendar or automatic input repair.

## 10. State and dependencies

**KNOWN** — the core is stateless and side-effect free.

**KNOWN** — `go.mod` targets Go 1.26 and declares no third-party modules.

**KNOWN** — there is no database, scheduler, network client, holiday service or notification adapter.

## 11. Security and correctness boundary

**KNOWN** — the component calculates a deadline from supplied data; it does not decide whether an obligation applies or whether the supplied trigger is legally correct.

**KNOWN** — rejecting contradictory offset/timezone input prevents the engine from silently interpreting an impossible request as a plausible date.

**KNOWN** — rejecting a seven-day weekend prevents unbounded iteration from invalid caller configuration.

**INFERRED** — callers should validate legal/jurisdictional source data before invoking this component because no such validation exists here.

## 12. Executable evidence

The package tests include:

- `TestHours`;
- `TestHoursAreElapsedAcrossDST`;
- `TestCalendarDaysPreserveWallClockAcrossDST`;
- `TestCalendarDayRejectsDSTGap`;
- `TestBusinessDaysHoliday`;
- `TestCustomWeekend`;
- `TestAdjustment`;
- `TestPreviousBusinessDayAdjustment`;
- `TestLifecycleDoesNotChangeDeadline`;
- `TestInvalidCalendarAndContract`;
- `TestTriggerOffsetMustMatchTimezone`.

CI additionally executes formatting, vet, race tests, native build, exact version, Linux ARM64 build and the Bafflewick CLI fixture.

## 13. Bafflewick acceptance path

**KNOWN** — the executable fixture starts at Friday 16 October 2026 17:00 in `Australia/Sydney`, adds three business days with the default Saturday/Sunday weekend and expects Wednesday 21 October 2026 17:00 +11:00.

**KNOWN** — the example makes no claim that a real legal obligation uses that period.

## 14. Operational behaviour

**KNOWN** — the process is short-lived and synchronous.

**KNOWN** — there is no internal telemetry, cache, durable state, service discovery or background worker.

**UNKNOWN** — production throughput, latency and memory behaviour on the Garage constrained-hardware target have not been benchmarked in V0.1. Linux ARM64 CI proves compilation only.

## 15. Known limits

**KNOWN** — no legal/regulatory applicability logic.

**KNOWN** — no business-hours or partial-day calendar.

**KNOWN** — no holiday source or jurisdiction database.

**KNOWN** — no automatic extensions, waivers, persistence or notification.

**KNOWN** — no leap-second-specific model.

**KNOWN** — ambiguous fall-back wall times are resolved by the explicit RFC3339 instant when its offset matches the configured timezone.

## 16. Unknowns

**UNKNOWN** — completeness/correctness of calendars supplied by a production caller.

**UNKNOWN** — deployment host timezone-data currency and operational hardening.

**UNKNOWN** — constrained-hardware runtime benchmark.

None of these unknowns is hidden by the implementation.

## 17. Recovered architecture

```text
CALLER
  |
  v
+----------------------------+
| cmd/obligation-clock       |
| JSON/file + exit contract  |
+-------------+--------------+
              |
              v
+----------------------------+
| internal/clock.Calculate   |
| validate -> calculate      |
| -> adjust -> result        |
+-------------+--------------+
              |
              v
        deterministic result
```
