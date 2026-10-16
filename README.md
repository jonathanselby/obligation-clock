# Obligation Clock

[![test](https://github.com/jonathanselby/obligation-clock/actions/workflows/test.yml/badge.svg)](https://github.com/jonathanselby/obligation-clock/actions/workflows/test.yml)

**Deadlines are data, not code.**

`obligation-clock` is a small deterministic Go component for calculating elapsed-hour, calendar-day and business-day deadlines from an explicit trigger and user-supplied calendar. It contains no jurisdiction, regulator or legal-applicability logic.

V0.1 is independently usable; CRAB is one consumer, not part of the runtime contract.

## What V0.1 proves

- elapsed-hour deadlines use elapsed time, so local wall time can move across a DST transition;
- calendar-day and business-day arithmetic preserve the trigger's local wall-clock fields when that local time exists;
- a calendar result that would land in a nonexistent DST wall time is rejected rather than silently normalised;
- the RFC3339 trigger's UTC offset must agree with the supplied IANA timezone at that instant;
- business days skip configured weekend weekdays and exact user-supplied holiday dates;
- weekend configuration must leave at least one possible business weekday;
- duplicate/invalid weekend and holiday entries are rejected;
- results can be adjusted to the next or previous business day;
- provisional and firm are labels only: changing lifecycle does not change the deadline.

## Install

Requires Go 1.26+.

```bash
go install github.com/jonathanselby/obligation-clock/cmd/obligation-clock@main
obligation-clock -version
```

Expected:

```text
obligation-clock 0.1.0
```

## CLI

```bash
obligation-clock -input request.json
```

Example request:

```json
{
  "trigger": "2026-10-17T17:00:00+11:00",
  "duration": 2,
  "unit": "business_days",
  "calendar": {
    "timezone": "Australia/Sydney",
    "weekend": ["Saturday", "Sunday"],
    "holidays": []
  },
  "lifecycle": "provisional"
}
```

## Units

`hours` means elapsed hours. `calendar_days` means local calendar arithmetic. `business_days` counts only days not configured as weekend or holiday.

This distinction is deliberate. Around DST, 24 elapsed hours and one calendar day can produce different local clock times.

## Calendar contract

`timezone` is mandatory and must be an IANA timezone recognised by Go. The trigger is RFC3339 and its numeric offset must match that timezone at the trigger instant; contradictory input is rejected.

If `weekend` is omitted, Saturday and Sunday are used. Custom weekend days are supported. All seven days cannot be weekend because business-day calculation would otherwise have no possible terminating weekday.

Holiday values are exact `YYYY-MM-DD` dates in the selected timezone. V0.1 does not download or infer public holidays.

## Exit codes

| Code | Meaning |
|---:|---|
| `0` | deadline calculated |
| `1` | deterministic calculation/calendar contract refusal |
| `2` | invalid CLI usage |
| `3` | malformed JSON/file I/O |

## Boundary

This component calculates a date. It does **not** decide whether an obligation applies, whether the trigger event really occurred, whether a supplied holiday calendar is legally complete, or whether a deadline has been extended/waived by some external rule.

There is no persistence, notification service, holiday feed or legal rules database.

## Engineering evidence

- [`PRD.md`](PRD.md) — product boundary and acceptance criteria;
- [`EDD.md`](EDD.md) — intended arithmetic and failure semantics;
- [`architecture.txt`](architecture.txt) — compact runtime view;
- [`internal/clock`](internal/clock) — implementation and behavioural tests;
- [`examples/bafflewick`](examples/bafflewick) — fictional acceptance case;
- [`generated/EDD.md`](generated/EDD.md) — implementation-derived evidence pinned to a source commit;
- [`generated/COMPARISON.md`](generated/COMPARISON.md) — intended versus recovered design;
- [`RELEASE.md`](RELEASE.md) — release gate.

## Licence

MIT. See [`LICENSE`](LICENSE).
