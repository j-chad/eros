package scheduler

import (
	"backend/internal/testutil"
	"fmt"
	"testing"
	"time"
)

// bits returns a bitfield with the given bit positions set.
func bits(positions ...int) uint64 {
	var b uint64
	for _, p := range positions {
		b |= 1 << p
	}
	return b
}

// bitRange returns a bitfield with all positions from min to max (inclusive) set.
func bitRange(min, max int) uint64 {
	var b uint64
	for i := min; i <= max; i++ {
		b |= 1 << i
	}
	return b
}

// bitStep returns a bitfield with positions from start to max (inclusive), stepping by step.
func bitStep(start, max, step int) uint64 {
	var b uint64
	for i := start; i <= max; i += step {
		b |= 1 << i
	}
	return b
}

// --- Parse: valid expressions ---

func TestParseCronExpression_Wildcard(t *testing.T) {
	expr, err := ParseCronExpression("* * * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bitRange(0, 59))
	testutil.Equal(t, expr.hour, bitRange(0, 23))
	testutil.Equal(t, expr.dayOfMonth, bitRange(1, 31))
	testutil.Equal(t, expr.month, bitRange(1, 12))
	// day-of-week wildcard expands to 0-7, with Sunday normalization setting both bits 0 and 7
	testutil.Equal(t, expr.dayOfWeek, bitRange(0, 7))
}

func TestParseCronExpression_ExactValues(t *testing.T) {
	expr, err := ParseCronExpression("30 12 15 6 3")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(30))
	testutil.Equal(t, expr.hour, bits(12))
	testutil.Equal(t, expr.dayOfMonth, bits(15))
	testutil.Equal(t, expr.month, bits(6))
	testutil.Equal(t, expr.dayOfWeek, bits(3))
}

func TestParseCronExpression_MinBoundary(t *testing.T) {
	expr, err := ParseCronExpression("0 0 1 1 0")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(0))
	testutil.Equal(t, expr.hour, bits(0))
	testutil.Equal(t, expr.dayOfMonth, bits(1))
	testutil.Equal(t, expr.month, bits(1))
	// 0 is Sunday, normalization sets both bit 0 and bit 7
	testutil.Equal(t, expr.dayOfWeek, bits(0, 7))
}

func TestParseCronExpression_MaxBoundary(t *testing.T) {
	expr, err := ParseCronExpression("59 23 31 12 6")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(59))
	testutil.Equal(t, expr.hour, bits(23))
	testutil.Equal(t, expr.dayOfMonth, bits(31))
	testutil.Equal(t, expr.month, bits(12))
	testutil.Equal(t, expr.dayOfWeek, bits(6))
}

func TestParseCronExpression_Step(t *testing.T) {
	expr, err := ParseCronExpression("*/15 * * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(0, 15, 30, 45))
}

func TestParseCronExpression_StepEvery2Hours(t *testing.T) {
	expr, err := ParseCronExpression("0 */2 * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(0))
	testutil.Equal(t, expr.hour, bitStep(0, 23, 2))
}

func TestParseCronExpression_StepEvery6Hours(t *testing.T) {
	expr, err := ParseCronExpression("0 */6 * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.hour, bits(0, 6, 12, 18))
}

func TestParseCronExpression_Range(t *testing.T) {
	expr, err := ParseCronExpression("1-5 * * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bitRange(1, 5))
}

func TestParseCronExpression_RangeFullField(t *testing.T) {
	// 0-59 should be equivalent to *
	expr, err := ParseCronExpression("0-59 * * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bitRange(0, 59))
}

func TestParseCronExpression_RangeWithStep(t *testing.T) {
	expr, err := ParseCronExpression("1-30/5 * * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(1, 6, 11, 16, 21, 26))
}

func TestParseCronExpression_RangeWithStep_DayOfWeek(t *testing.T) {
	// Weekdays: Monday(1) through Friday(5)
	expr, err := ParseCronExpression("0 9 * * 1-5")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.dayOfWeek, bitRange(1, 5))
}

func TestParseCronExpression_List(t *testing.T) {
	expr, err := ParseCronExpression("0,15,30,45 * * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(0, 15, 30, 45))
}

func TestParseCronExpression_ListWithRanges(t *testing.T) {
	// "1-3,7,10-12" = bits 1,2,3,7,10,11,12
	expr, err := ParseCronExpression("1-3,7,10-12 * * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(1, 2, 3, 7, 10, 11, 12))
}

func TestParseCronExpression_ListWithSteps(t *testing.T) {
	// "0-10/5,30-40/5" = bits 0,5,10,30,35,40
	expr, err := ParseCronExpression("0-10/5,30-40/5 * * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(0, 5, 10, 30, 35, 40))
}

func TestParseCronExpression_DailyAt3AM(t *testing.T) {
	expr, err := ParseCronExpression("0 3 * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(0))
	testutil.Equal(t, expr.hour, bits(3))
	testutil.Equal(t, expr.dayOfMonth, bitRange(1, 31))
	testutil.Equal(t, expr.month, bitRange(1, 12))
	testutil.Equal(t, expr.dayOfWeek, bitRange(0, 7))
}

func TestParseCronExpression_WeekdayBusinessHours(t *testing.T) {
	// Every 10 minutes, 9-17, weekdays
	expr, err := ParseCronExpression("*/10 9-17 * * 1-5")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(0, 10, 20, 30, 40, 50))
	testutil.Equal(t, expr.hour, bitRange(9, 17))
	testutil.Equal(t, expr.dayOfWeek, bitRange(1, 5))
}

// --- Parse: error cases ---

func TestParseCronExpression_EmptyString(t *testing.T) {
	_, err := ParseCronExpression("")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_TooFewFields(t *testing.T) {
	_, err := ParseCronExpression("* * *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_TooManyFields(t *testing.T) {
	_, err := ParseCronExpression("* * * * * *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_MinuteOutOfRange(t *testing.T) {
	_, err := ParseCronExpression("60 * * * *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_HourOutOfRange(t *testing.T) {
	_, err := ParseCronExpression("* 24 * * *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_DayOfMonthZero(t *testing.T) {
	_, err := ParseCronExpression("* * 0 * *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_DayOfMonthOutOfRange(t *testing.T) {
	_, err := ParseCronExpression("* * 32 * *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_MonthZero(t *testing.T) {
	_, err := ParseCronExpression("* * * 0 *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_MonthOutOfRange(t *testing.T) {
	_, err := ParseCronExpression("* * * 13 *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_DayOfWeekOutOfRange(t *testing.T) {
	_, err := ParseCronExpression("* * * * 8")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_Sunday7(t *testing.T) {
	// 7 is an alias for 0 (Sunday)
	expr, err := ParseCronExpression("0 0 * * 7")
	testutil.NilErr(t, err)
	// Both bit 0 and bit 7 should be set
	testutil.True(t, expr.dayOfWeek&(1<<0) != 0, "bit 0 (Sunday) should be set")
	testutil.True(t, expr.dayOfWeek&(1<<7) != 0, "bit 7 (Sunday alias) should be set")
}

func TestParseCronExpression_Sunday0And7Equivalent(t *testing.T) {
	// "* * * * 0" and "* * * * 7" should produce the same bitfield
	expr0, err := ParseCronExpression("0 0 * * 0")
	testutil.NilErr(t, err)
	expr7, err := ParseCronExpression("0 0 * * 7")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr0.dayOfWeek, expr7.dayOfWeek)
}

func TestParseCronExpression_Sunday7InRange(t *testing.T) {
	// "5-7" in day-of-week = Friday(5), Saturday(6), Sunday(7→0)
	expr, err := ParseCronExpression("0 0 * * 5-7")
	testutil.NilErr(t, err)
	testutil.True(t, expr.dayOfWeek&(1<<5) != 0, "bit 5 (Friday) should be set")
	testutil.True(t, expr.dayOfWeek&(1<<6) != 0, "bit 6 (Saturday) should be set")
	testutil.True(t, expr.dayOfWeek&(1<<0) != 0, "bit 0 (Sunday) should be set via 7 alias")
}

func TestParseCronExpression_Sunday7InList(t *testing.T) {
	// "0,6,7" — Sunday(0), Saturday(6), Sunday(7) — should have both Sunday bits
	expr, err := ParseCronExpression("0 0 * * 0,6,7")
	testutil.NilErr(t, err)
	testutil.True(t, expr.dayOfWeek&(1<<0) != 0, "bit 0 (Sunday) should be set")
	testutil.True(t, expr.dayOfWeek&(1<<6) != 0, "bit 6 (Saturday) should be set")
	testutil.True(t, expr.dayOfWeek&(1<<7) != 0, "bit 7 (Sunday alias) should be set")
}

func TestParseCronExpression_NonNumeric(t *testing.T) {
	_, err := ParseCronExpression("abc * * * *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_StepZero(t *testing.T) {
	_, err := ParseCronExpression("*/0 * * * *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_NegativeValue(t *testing.T) {
	_, err := ParseCronExpression("-1 * * * *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_ReversedRange(t *testing.T) {
	_, err := ParseCronExpression("5-2 * * * *")
	testutil.NotNilErr(t, err)
}

func TestParseCronExpression_StepOnExactValue(t *testing.T) {
	// "5/2" means "starting at 5, every 2nd value through end of field range"
	// For minutes (0-59): 5, 7, 9, 11, ..., 59
	expr, err := ParseCronExpression("5/2 * * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bitStep(5, 59, 2))
}

func TestParseCronExpression_StepOnExactValue_Hours(t *testing.T) {
	// "6/4" on hours means starting at 6, every 4th: 6, 10, 14, 18, 22
	expr, err := ParseCronExpression("0 6/4 * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.hour, bits(6, 10, 14, 18, 22))
}

func TestParseCronExpression_StepOnExactValue_NoStep(t *testing.T) {
	// Plain "5" without a step should still be just bit 5
	expr, err := ParseCronExpression("5 * * * *")
	testutil.NilErr(t, err)
	testutil.Equal(t, expr.minute, bits(5))
}

// --- Error message quality ---

func TestParseCronExpression_ErrorIdentifiesField(t *testing.T) {
	tests := []struct {
		name  string
		input string
		field string
	}{
		{"minute", "60 0 1 1 0", "minute"},
		{"hour", "0 24 1 1 0", "hour"},
		{"day of month", "0 0 32 1 0", "day of month"},
		{"month", "0 0 1 13 0", "month"},
		{"day of week", "0 0 1 1 8", "day of week"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCronExpression(tt.input)
			testutil.NotNilErr(t, err)
			testutil.True(t, containsSubstring(err.Error(), tt.field),
				fmt.Sprintf("error should mention %q, got: %s", tt.field, err.Error()))
		})
	}
}

// --- Matches ---

func TestCronExpression_Matches_EveryMinute(t *testing.T) {
	expr, err := ParseCronExpression("* * * * *")
	testutil.NilErr(t, err)

	// Should match any time
	testutil.True(t, expr.Matches(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
	testutil.True(t, expr.Matches(time.Date(2025, 6, 15, 12, 30, 0, 0, time.UTC)))
	testutil.True(t, expr.Matches(time.Date(2025, 12, 31, 23, 59, 0, 0, time.UTC)))
}

func TestCronExpression_Matches_ExactTime(t *testing.T) {
	// 30 12 15 6 * = 12:30 on June 15th, any day of week
	expr, err := ParseCronExpression("30 12 15 6 *")
	testutil.NilErr(t, err)

	testutil.True(t, expr.Matches(time.Date(2025, 6, 15, 12, 30, 0, 0, time.UTC)))
	testutil.False(t, expr.Matches(time.Date(2025, 6, 15, 12, 31, 0, 0, time.UTC)), "wrong minute")
	testutil.False(t, expr.Matches(time.Date(2025, 6, 15, 11, 30, 0, 0, time.UTC)), "wrong hour")
	testutil.False(t, expr.Matches(time.Date(2025, 6, 16, 12, 30, 0, 0, time.UTC)), "wrong day")
	testutil.False(t, expr.Matches(time.Date(2025, 7, 15, 12, 30, 0, 0, time.UTC)), "wrong month")
}

func TestCronExpression_Matches_DailyAt3AM(t *testing.T) {
	expr, err := ParseCronExpression("0 3 * * *")
	testutil.NilErr(t, err)

	testutil.True(t, expr.Matches(time.Date(2025, 1, 1, 3, 0, 0, 0, time.UTC)))
	testutil.True(t, expr.Matches(time.Date(2025, 6, 15, 3, 0, 0, 0, time.UTC)))
	testutil.False(t, expr.Matches(time.Date(2025, 1, 1, 3, 1, 0, 0, time.UTC)), "minute must be 0")
	testutil.False(t, expr.Matches(time.Date(2025, 1, 1, 4, 0, 0, 0, time.UTC)), "wrong hour")
}

func TestCronExpression_Matches_Every15Minutes(t *testing.T) {
	expr, err := ParseCronExpression("*/15 * * * *")
	testutil.NilErr(t, err)

	testutil.True(t, expr.Matches(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
	testutil.True(t, expr.Matches(time.Date(2025, 1, 1, 0, 15, 0, 0, time.UTC)))
	testutil.True(t, expr.Matches(time.Date(2025, 1, 1, 0, 30, 0, 0, time.UTC)))
	testutil.True(t, expr.Matches(time.Date(2025, 1, 1, 0, 45, 0, 0, time.UTC)))
	testutil.False(t, expr.Matches(time.Date(2025, 1, 1, 0, 10, 0, 0, time.UTC)))
	testutil.False(t, expr.Matches(time.Date(2025, 1, 1, 0, 59, 0, 0, time.UTC)))
}

func TestCronExpression_Matches_Weekdays(t *testing.T) {
	// 0 9 * * 1-5 = 9:00 AM on weekdays
	expr, err := ParseCronExpression("0 9 * * 1-5")
	testutil.NilErr(t, err)

	// 2025-01-06 is a Monday
	testutil.True(t, expr.Matches(time.Date(2025, 1, 6, 9, 0, 0, 0, time.UTC)), "Monday")
	testutil.True(t, expr.Matches(time.Date(2025, 1, 7, 9, 0, 0, 0, time.UTC)), "Tuesday")
	testutil.True(t, expr.Matches(time.Date(2025, 1, 8, 9, 0, 0, 0, time.UTC)), "Wednesday")
	testutil.True(t, expr.Matches(time.Date(2025, 1, 9, 9, 0, 0, 0, time.UTC)), "Thursday")
	testutil.True(t, expr.Matches(time.Date(2025, 1, 10, 9, 0, 0, 0, time.UTC)), "Friday")
	testutil.False(t, expr.Matches(time.Date(2025, 1, 11, 9, 0, 0, 0, time.UTC)), "Saturday")
	testutil.False(t, expr.Matches(time.Date(2025, 1, 12, 9, 0, 0, 0, time.UTC)), "Sunday")
}

func TestCronExpression_Matches_Sunday0(t *testing.T) {
	// * * * * 0 = every minute on Sunday
	expr, err := ParseCronExpression("* * * * 0")
	testutil.NilErr(t, err)

	// 2025-01-12 is a Sunday
	testutil.True(t, expr.Matches(time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC)))
	// 2025-01-13 is a Monday
	testutil.False(t, expr.Matches(time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC)))
}

func TestCronExpression_Matches_Sunday7(t *testing.T) {
	// * * * * 7 should behave identically to * * * * 0
	expr, err := ParseCronExpression("* * * * 7")
	testutil.NilErr(t, err)

	// 2025-01-12 is a Sunday (Weekday() == 0)
	testutil.True(t, expr.Matches(time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC)))
	testutil.False(t, expr.Matches(time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC)))
}

func TestCronExpression_Matches_IgnoresSeconds(t *testing.T) {
	expr, err := ParseCronExpression("30 12 * * *")
	testutil.NilErr(t, err)

	// Different seconds within the same minute should all match
	testutil.True(t, expr.Matches(time.Date(2025, 1, 1, 12, 30, 0, 0, time.UTC)))
	testutil.True(t, expr.Matches(time.Date(2025, 1, 1, 12, 30, 45, 0, time.UTC)))
	testutil.True(t, expr.Matches(time.Date(2025, 1, 1, 12, 30, 59, 999999999, time.UTC)))
}

func TestCronExpression_Matches_SpecificDayOfMonth(t *testing.T) {
	// 0 0 1 * * = midnight on the 1st of every month
	expr, err := ParseCronExpression("0 0 1 * *")
	testutil.NilErr(t, err)

	testutil.True(t, expr.Matches(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
	testutil.True(t, expr.Matches(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)))
	testutil.False(t, expr.Matches(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)))
}

func TestCronExpression_Matches_SpecificMonth(t *testing.T) {
	// 0 0 25 12 * = midnight on December 25th
	expr, err := ParseCronExpression("0 0 25 12 *")
	testutil.NilErr(t, err)

	testutil.True(t, expr.Matches(time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)))
	testutil.False(t, expr.Matches(time.Date(2025, 12, 24, 0, 0, 0, 0, time.UTC)))
	testutil.False(t, expr.Matches(time.Date(2025, 11, 25, 0, 0, 0, 0, time.UTC)))
}

func TestCronExpression_Matches_RespectsTimezone(t *testing.T) {
	// 0 9 * * * = 9:00 AM
	expr, err := ParseCronExpression("0 9 * * *")
	testutil.NilErr(t, err)

	nzst := time.FixedZone("NZST", 12*60*60)

	// 9:00 AM NZST = 21:00 UTC the previous day
	nzTime := time.Date(2025, 1, 1, 9, 0, 0, 0, nzst)
	utcTime := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)

	// Matches should use the time's own timezone
	testutil.True(t, expr.Matches(nzTime), "should match 9am NZST")
	testutil.True(t, expr.Matches(utcTime), "should match 9am UTC")

	// The same instant in different timezones may not match
	// 9am NZST is 9pm UTC the day before — hour 21 doesn't match
	testutil.False(t, expr.Matches(nzTime.UTC()), "9am NZST converted to UTC is 9pm, should not match")
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && searchSubstring(s, sub)
}

func searchSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
