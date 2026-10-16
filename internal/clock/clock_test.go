package clock

import (
	"strings"
	"testing"
)

func TestHours(t *testing.T) {
	r, err := Calculate(Request{Trigger: "2026-10-17T10:00:00+11:00", Duration: 72, Unit: "hours", Calendar: Calendar{Timezone: "Australia/Sydney"}, Lifecycle: "provisional"})
	if err != nil || !strings.HasPrefix(r.Deadline, "2026-10-20T10:00:00") {
		t.Fatalf("%+v %v", r, err)
	}
}

func TestHoursAreElapsedAcrossDST(t *testing.T) {
	r, err := Calculate(Request{Trigger: "2026-10-03T10:00:00+10:00", Duration: 24, Unit: "hours", Calendar: Calendar{Timezone: "Australia/Sydney"}, Lifecycle: "firm"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Deadline != "2026-10-04T11:00:00+11:00" {
		t.Fatalf("got %s", r.Deadline)
	}
}

func TestCalendarDaysPreserveWallClockAcrossDST(t *testing.T) {
	r, err := Calculate(Request{Trigger: "2026-10-03T10:00:00+10:00", Duration: 2, Unit: "calendar_days", Calendar: Calendar{Timezone: "Australia/Sydney"}, Lifecycle: "firm"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Deadline != "2026-10-05T10:00:00+11:00" {
		t.Fatalf("got %s", r.Deadline)
	}
}

func TestCalendarDayRejectsDSTGap(t *testing.T) {
	_, err := Calculate(Request{Trigger: "2026-10-03T02:30:00+10:00", Duration: 1, Unit: "calendar_days", Calendar: Calendar{Timezone: "Australia/Sydney"}, Lifecycle: "firm"})
	if err == nil || !strings.Contains(err.Error(), "nonexistent local wall time") {
		t.Fatalf("expected DST gap refusal, got %v", err)
	}
}

func TestBusinessDaysHoliday(t *testing.T) {
	r, err := Calculate(Request{Trigger: "2026-12-24T09:00:00+11:00", Duration: 2, Unit: "business_days", Calendar: Calendar{Timezone: "Australia/Sydney", Holidays: []string{"2026-12-25", "2026-12-28"}}, Lifecycle: "firm"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Deadline != "2026-12-30T09:00:00+11:00" {
		t.Fatalf("got %s", r.Deadline)
	}
}

func TestCustomWeekend(t *testing.T) {
	r, err := Calculate(Request{Trigger: "2026-10-15T09:00:00+11:00", Duration: 1, Unit: "business_days", Calendar: Calendar{Timezone: "Australia/Sydney", Weekend: []string{"Friday", "Saturday"}}, Lifecycle: "firm"})
	if err != nil || r.Deadline != "2026-10-18T09:00:00+11:00" {
		t.Fatalf("%+v %v", r, err)
	}
}

func TestAdjustment(t *testing.T) {
	r, err := Calculate(Request{Trigger: "2026-10-16T09:00:00+11:00", Duration: 1, Unit: "calendar_days", Calendar: Calendar{Timezone: "Australia/Sydney", Adjustment: "next_business_day"}, Lifecycle: "firm"})
	if err != nil || r.Deadline != "2026-10-19T09:00:00+11:00" || !r.Adjusted {
		t.Fatalf("%+v %v", r, err)
	}
}

func TestPreviousBusinessDayAdjustment(t *testing.T) {
	r, err := Calculate(Request{Trigger: "2026-10-16T09:00:00+11:00", Duration: 1, Unit: "calendar_days", Calendar: Calendar{Timezone: "Australia/Sydney", Adjustment: "previous_business_day"}, Lifecycle: "firm"})
	if err != nil || r.Deadline != "2026-10-16T09:00:00+11:00" || !r.Adjusted {
		t.Fatalf("%+v %v", r, err)
	}
}

func TestLifecycleDoesNotChangeDeadline(t *testing.T) {
	base := Request{Trigger: "2026-10-17T10:00:00+11:00", Duration: 2, Unit: "calendar_days", Calendar: Calendar{Timezone: "Australia/Sydney"}}
	base.Lifecycle = "provisional"
	p, err := Calculate(base)
	if err != nil {
		t.Fatal(err)
	}
	base.Lifecycle = "firm"
	f, err := Calculate(base)
	if err != nil {
		t.Fatal(err)
	}
	if p.Deadline != f.Deadline || p.Lifecycle == f.Lifecycle {
		t.Fatalf("lifecycle affected deadline: p=%+v f=%+v", p, f)
	}
}

func TestInvalidCalendarAndContract(t *testing.T) {
	base := Request{Trigger: "2026-10-17T10:00:00+11:00", Duration: 1, Unit: "calendar_days", Calendar: Calendar{Timezone: "Australia/Sydney"}, Lifecycle: "firm"}
	cases := []Request{
		func() Request { x := base; x.Duration = -1; return x }(),
		func() Request { x := base; x.Lifecycle = "draft"; return x }(),
		func() Request { x := base; x.Unit = "weeks"; return x }(),
		func() Request { x := base; x.Calendar.Timezone = ""; return x }(),
		func() Request { x := base; x.Calendar.Timezone = "Mars/Olympus"; return x }(),
		func() Request { x := base; x.Calendar.Adjustment = "nearest"; return x }(),
		func() Request { x := base; x.Calendar.Weekend = []string{"Funday"}; return x }(),
		func() Request { x := base; x.Calendar.Weekend = []string{"Saturday", "Saturday"}; return x }(),
		func() Request { x := base; x.Calendar.Weekend = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}; return x }(),
		func() Request { x := base; x.Calendar.Holidays = []string{"2026-12-25", "2026-12-25"}; return x }(),
		func() Request { x := base; x.Calendar.Holidays = []string{"25-12-2026"}; return x }(),
	}
	for i, tc := range cases {
		if _, err := Calculate(tc); err == nil {
			t.Fatalf("case %d unexpectedly accepted", i)
		}
	}
}

func TestTriggerOffsetMustMatchTimezone(t *testing.T) {
	_, err := Calculate(Request{Trigger: "2026-10-17T10:00:00+10:00", Duration: 1, Unit: "hours", Calendar: Calendar{Timezone: "Australia/Sydney"}, Lifecycle: "firm"})
	if err == nil || !strings.Contains(err.Error(), "offset does not match timezone") {
		t.Fatalf("mismatched offset accepted: %v", err)
	}
}
