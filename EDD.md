# Obligation Clock — Engineering Design Document

**Version:** 0.1

## 1. Purpose

A small deterministic deadline calculator for elapsed hours, calendar days and business days using an explicit trigger and caller-supplied calendar.

> **Deadlines are data, not code.**

## 2. Design constraints

1. Arithmetic only: no legal applicability engine.
2. Explicit timezone and calendar inputs.
3. Contradictory or impossible time data fails visibly.
4. Elapsed time and calendar time are separate concepts.
5. Provisional/firm is metadata, not arithmetic.
6. Duration is bounded before arithmetic so extreme input cannot overflow hour multiplication or drive unbounded business-day iteration.

## 3. Runtime model

```text
request
 trigger / duration / unit / calendar / lifecycle
                    |
                    v
          validate trigger + timezone
                    |
                    v
          validate weekend/holidays
                    |
                    v
             calculate unit
          /        |         \
      hours   calendar_days  business_days
                    |
                    v
          optional business-day
              adjustment
                    |
                    v
               deadline
```

The CLI is a thin JSON-file wrapper around `internal/clock.Calculate`.

## 4. Trigger and timezone contract

`trigger` is RFC3339 and `calendar.timezone` is a required IANA timezone.

The supplied trigger offset must equal the named timezone's UTC offset at that instant. This prevents contradictory input such as a Sydney timezone paired with an offset that does not apply in Sydney at the trigger instant.

The trigger is then converted into the named location for local-calendar arithmetic.

## 5. Duration and elapsed-hour arithmetic

`Calculate` accepts durations from zero through `MaxDuration` (100,000 units). The bound is checked before unit arithmetic.

`hours` uses `time.Add(duration * time.Hour)`.

This is elapsed-time arithmetic. Around a DST transition, 24 elapsed hours may produce a different local wall-clock hour on the following day. That behaviour is intentional.

## 6. Calendar-day arithmetic

`calendar_days` uses `time.AddDate(0, 0, days)`.

After adding the date, the implementation checks that hour, minute, second and nanosecond still equal the original local wall-clock fields. If Go normalised an impossible local time—such as a spring-forward gap—the calculation returns an error rather than silently accepting the normalised result.

## 7. Business-day arithmetic

`business_days` repeatedly advances by one validated local calendar day and increments the count only when the resulting date is neither a configured weekend weekday nor configured holiday.

The common duration bound also caps this iterative path. This deliberately defines business-day counting in terms of whole local dates, not elapsed 24-hour periods.

## 8. Weekend calendar

If no weekend is supplied, Saturday and Sunday are used.

Custom weekend days must be valid Go weekday names and unique. All seven weekdays cannot be configured as weekend because then no business weekday could ever exist; such a calendar is rejected before calculation rather than allowing an infinite loop.

## 9. Holiday calendar

Holiday values are exact `YYYY-MM-DD` dates interpreted in the selected timezone. Entries must be valid and unique.

V0.1 does not infer, download or apply public-holiday rules. Completeness and legal correctness of the supplied holiday set belong to the caller.

## 10. Adjustment

Supported adjustment values are:

- empty / `none`;
- `next_business_day`;
- `previous_business_day`.

Adjustment uses the same validated calendar and local calendar-day arithmetic as business-day counting. `Result.Adjusted` is true only when at least one date step was required.

## 11. Lifecycle

Only `provisional` and `firm` are accepted.

The lifecycle is copied to the result and does not alter arithmetic. Tests explicitly prove that the same request with only the lifecycle changed produces the same deadline.

## 12. Result contract

`Result` contains:

- RFC3339 `deadline` in the configured timezone;
- lifecycle;
- timezone name;
- whether a business-day adjustment changed the initial result.

## 13. CLI and exits

```text
obligation-clock -input request.json
```

`-version` prints `obligation-clock 0.1.0`.

Exit codes:

```text
0 calculated successfully
1 deterministic calculation/calendar-contract refusal
2 CLI usage error
3 malformed JSON or file I/O
```

The CLI maps calculation errors to exit 1 because the process executed correctly but the supplied deadline contract could not be satisfied.

## 14. State and dependencies

The core is stateless and side-effect free. There is no database, network call, holiday feed, scheduler or persistence layer.

The module targets Go 1.26 and uses only the standard library.

## 15. Test strategy

Named tests cover:

- elapsed hours;
- elapsed hours across DST;
- calendar days across DST;
- DST-gap refusal;
- business days with holidays;
- custom weekend;
- next/previous adjustment;
- lifecycle invariance;
- invalid calendar/contract inputs, including the duration ceiling;
- trigger offset/timezone mismatch.

CI adds format, vet, race tests, native build, exact version, Linux ARM64 build and executable Bafflewick acceptance.

## 16. Known limits

- no legal/regulatory applicability logic;
- no business-hours or partial-day model;
- no holiday source or jurisdiction database;
- no automatic extensions/waivers;
- no persistence or notification;
- no leap-second-specific model;
- ambiguous fall-back wall times are accepted according to the explicit RFC3339 trigger offset when that offset matches the named timezone at the instant;
- ARM64 is build proof, not a constrained-hardware benchmark.

## 17. Acceptance-to-evidence map

| Requirement | Evidence |
|---|---|
| elapsed hours | `TestHours` |
| elapsed hours across DST | `TestHoursAreElapsedAcrossDST` |
| calendar days across DST | `TestCalendarDaysPreserveWallClockAcrossDST` |
| DST-gap refusal | `TestCalendarDayRejectsDSTGap` |
| business days + holiday | `TestBusinessDaysHoliday` |
| custom weekend | `TestCustomWeekend` |
| next/previous adjustment | `TestAdjustment`, `TestPreviousBusinessDayAdjustment` |
| lifecycle invariance | `TestLifecycleDoesNotChangeDeadline` |
| invalid contract/calendar/duration | `TestInvalidCalendarAndContract` |
| offset/timezone mismatch | `TestTriggerOffsetMustMatchTimezone` |
| format/vet/race/build/version/ARM64/Bafflewick | `.github/workflows/test.yml` |
