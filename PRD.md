# Obligation Clock — Product Requirements Document

**Version:** 0.1

## 1. Problem

Legal, contractual and operational deadlines use different units, timezones and working calendars. Hard-coding jurisdictional rules into application code makes the arithmetic difficult to reuse and, more dangerously, encourages silent assumptions around DST, holiday data and trigger timezones.

## 2. Product objective

Provide a deliberately small Go component that calculates a deadline from an explicit trigger, duration, unit and user-supplied calendar while refusing contradictory or impossible calendar inputs.

> **Deadlines are data, not code.**

## 3. In scope

- elapsed-hour deadlines;
- local calendar-day deadlines;
- business-day deadlines;
- IANA timezone selection;
- explicit trigger offset/timezone consistency;
- configurable weekend weekdays and holidays;
- next/previous business-day adjustment;
- provisional/firm lifecycle labels;
- bounded duration input;
- deterministic validation and structured result;
- Go 1.26 and Linux ARM64 build proof.

## 4. Out of scope

Legal applicability, jurisdiction/regulator rules, legal advice, trigger-event truth, holiday-calendar acquisition/completeness, business hours/partial days, notification delivery, clock persistence, extensions/waivers and leap-second modelling.

## 5. Functional requirements

### FR1 — Trigger/timezone consistency
The trigger is RFC3339. The selected IANA timezone is mandatory. The trigger's numeric UTC offset must equal the selected timezone's offset at that instant; contradictory input is rejected.

### FR2 — Elapsed hours
`hours` adds elapsed hours. Across DST this may change local wall-clock time.

### FR3 — Calendar days
`calendar_days` uses local calendar-day arithmetic and preserves hour/minute/second/nanosecond when the destination wall time exists. A destination in a DST gap is rejected rather than normalised silently.

### FR4 — Business days
`business_days` advances one local calendar day at a time and counts only dates that are neither configured weekend weekdays nor configured holidays.

### FR5 — Valid business calendar
Weekend names are valid Go weekday names, are unique, and cannot cover all seven weekdays. Holiday dates are unique valid `YYYY-MM-DD` values.

### FR6 — Adjustment
`none`, `next_business_day` and `previous_business_day` are supported. Adjustment uses the same calendar-day arithmetic and business calendar.

### FR7 — Lifecycle is metadata
`provisional` and `firm` are the only lifecycle values. Changing the lifecycle label cannot change the calculated deadline.

### FR8 — Visible failure
Negative duration, duration above 100,000 units, invalid unit, lifecycle, timezone, trigger, offset mismatch, calendar entry or adjustment is rejected rather than silently coerced.

## 6. Acceptance criteria

Automated evidence must show:

1. elapsed-hour calculation;
2. 24 elapsed hours across Sydney DST changes wall-clock time as expected;
3. calendar days preserve wall-clock time across a normal DST transition;
4. calendar arithmetic landing in a DST spring-forward gap is rejected;
5. business days skip weekend and configured holidays;
6. custom weekend weekdays work;
7. next and previous business-day adjustment work;
8. provisional and firm yield the same deadline;
9. invalid unit/lifecycle/timezone/adjustment/weekend/holiday/duration, including over-limit duration, are rejected;
10. duplicate weekend/holiday entries and seven-day weekends are rejected;
11. trigger offset/timezone mismatch is rejected;
12. format, vet, race, native build, exact version and Linux ARM64 build pass.

## 7. Definition of done

All acceptance criteria have executable evidence; calendar arithmetic and DST-gap semantics are explicit; no jurisdictional data is embedded in code; Go 1.26 CI is green; and generated design evidence is recovered from the exact release-candidate source commit.
