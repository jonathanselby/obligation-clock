package clock

import (
	"errors"
	"fmt"
	"time"
)

type Calendar struct {
	Timezone   string   `json:"timezone"`
	Weekend    []string `json:"weekend,omitempty"`
	Holidays   []string `json:"holidays,omitempty"`
	Adjustment string   `json:"adjustment,omitempty"`
}

type Request struct {
	Trigger   string   `json:"trigger"`
	Duration  int      `json:"duration"`
	Unit      string   `json:"unit"`
	Calendar  Calendar `json:"calendar"`
	Lifecycle string   `json:"lifecycle"`
}

type Result struct {
	Deadline  string `json:"deadline"`
	Lifecycle string `json:"lifecycle"`
	Timezone  string `json:"timezone"`
	Adjusted  bool   `json:"adjusted"`
}

func Calculate(r Request) (Result, error) {
	if r.Duration < 0 {
		return Result{}, errors.New("duration cannot be negative")
	}
	if r.Lifecycle != "provisional" && r.Lifecycle != "firm" {
		return Result{}, errors.New("lifecycle must be provisional or firm")
	}
	if r.Calendar.Timezone == "" {
		return Result{}, errors.New("timezone is required")
	}
	loc, err := time.LoadLocation(r.Calendar.Timezone)
	if err != nil {
		return Result{}, fmt.Errorf("timezone: %w", err)
	}
	trigger, err := time.Parse(time.RFC3339, r.Trigger)
	if err != nil {
		return Result{}, fmt.Errorf("trigger: %w", err)
	}
	if err := validateTriggerZone(trigger, loc); err != nil {
		return Result{}, err
	}
	t := trigger.In(loc)
	week, err := weekend(r.Calendar.Weekend)
	if err != nil {
		return Result{}, err
	}
	hol, err := holidays(r.Calendar.Holidays, loc)
	if err != nil {
		return Result{}, err
	}
	if err := validateAdjustment(r.Calendar.Adjustment); err != nil {
		return Result{}, err
	}

	var deadline time.Time
	switch r.Unit {
	case "hours":
		deadline = t.Add(time.Duration(r.Duration) * time.Hour)
	case "calendar_days":
		deadline, err = addCalendarDays(t, r.Duration)
		if err != nil {
			return Result{}, err
		}
	case "business_days":
		deadline = t
		for n := 0; n < r.Duration; {
			deadline, err = addCalendarDays(deadline, 1)
			if err != nil {
				return Result{}, err
			}
			if isBusiness(deadline, week, hol) {
				n++
			}
		}
	default:
		return Result{}, errors.New("unit must be hours, calendar_days or business_days")
	}

	adjusted := false
	switch r.Calendar.Adjustment {
	case "", "none":
	case "next_business_day":
		for !isBusiness(deadline, week, hol) {
			deadline, err = addCalendarDays(deadline, 1)
			if err != nil {
				return Result{}, err
			}
			adjusted = true
		}
	case "previous_business_day":
		for !isBusiness(deadline, week, hol) {
			deadline, err = addCalendarDays(deadline, -1)
			if err != nil {
				return Result{}, err
			}
			adjusted = true
		}
	}

	return Result{
		Deadline:  deadline.Format(time.RFC3339),
		Lifecycle: r.Lifecycle,
		Timezone:  r.Calendar.Timezone,
		Adjusted:  adjusted,
	}, nil
}

func validateTriggerZone(t time.Time, loc *time.Location) error {
	_, suppliedOffset := t.Zone()
	_, zoneOffset := t.In(loc).Zone()
	if suppliedOffset != zoneOffset {
		return fmt.Errorf("trigger offset does not match timezone at that instant")
	}
	return nil
}

func validateAdjustment(a string) error {
	switch a {
	case "", "none", "next_business_day", "previous_business_day":
		return nil
	default:
		return errors.New("invalid adjustment")
	}
}

func addCalendarDays(t time.Time, days int) (time.Time, error) {
	d := t.AddDate(0, 0, days)
	if d.Hour() != t.Hour() || d.Minute() != t.Minute() || d.Second() != t.Second() || d.Nanosecond() != t.Nanosecond() {
		return time.Time{}, errors.New("calendar arithmetic lands on a nonexistent local wall time")
	}
	return d, nil
}

func weekend(xs []string) (map[time.Weekday]bool, error) {
	if len(xs) == 0 {
		xs = []string{"Saturday", "Sunday"}
	}
	m := map[time.Weekday]bool{}
	for _, x := range xs {
		matched := false
		for d := time.Sunday; d <= time.Saturday; d++ {
			if d.String() != x {
				continue
			}
			if m[d] {
				return nil, fmt.Errorf("duplicate weekend day %q", x)
			}
			m[d] = true
			matched = true
			break
		}
		if !matched {
			return nil, fmt.Errorf("invalid weekend day %q", x)
		}
	}
	if len(m) == 7 {
		return nil, errors.New("calendar must contain at least one business weekday")
	}
	return m, nil
}

func holidays(xs []string, loc *time.Location) (map[string]bool, error) {
	m := map[string]bool{}
	for _, x := range xs {
		if _, err := time.ParseInLocation("2006-01-02", x, loc); err != nil {
			return nil, fmt.Errorf("holiday %q: %w", x, err)
		}
		if m[x] {
			return nil, fmt.Errorf("duplicate holiday %q", x)
		}
		m[x] = true
	}
	return m, nil
}

func isBusiness(t time.Time, w map[time.Weekday]bool, h map[string]bool) bool {
	return !w[t.Weekday()] && !h[t.Format("2006-01-02")]
}
