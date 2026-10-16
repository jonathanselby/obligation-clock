# Obligation Clock — Skill

## Use this when

Use this component when an application needs deterministic elapsed-hour, calendar-day or business-day deadline arithmetic from explicit trigger and calendar data.

## Core rule

> **Deadlines are data, not code.**

## Change discipline

1. Never embed legal applicability or jurisdiction rules in the calculator.
2. Keep elapsed time distinct from local calendar arithmetic.
3. Reject trigger offset/timezone contradictions.
4. Reject impossible local wall-time results rather than silently normalising them.
5. Validate weekend and holiday calendars before iterating.
6. Never allow a calendar with no possible business weekday.
7. Keep provisional/firm as metadata only.
8. Run format, vet, race, native and ARM64 gates after changes.
9. Keep the Bafflewick example executable and consistent with its prose.
10. Regenerate code-derived evidence before release.
