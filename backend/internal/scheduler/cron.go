package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CronExpression struct {
	minute, hour, dayOfMonth, month, dayOfWeek uint64
}

func (ce *CronExpression) Matches(t time.Time) bool {
	minuteMatch := (ce.minute & (1 << t.Minute())) != 0
	hourMatch := (ce.hour & (1 << t.Hour())) != 0
	dayOfMonthMatch := (ce.dayOfMonth & (1 << t.Day())) != 0
	monthMatch := (ce.month & (1 << (t.Month()))) != 0
	dayOfWeekMatch := (ce.dayOfWeek & (1 << t.Weekday())) != 0

	return minuteMatch && hourMatch && dayOfMonthMatch && monthMatch && dayOfWeekMatch
}

func ParseCronExpression(s string) (CronExpression, error) {
	fields := strings.Fields(s)
	if len(fields) != 5 {
		return CronExpression{}, fmt.Errorf("expected 5 fields, got %d", len(fields))
	}

	minuteBitfield, err := parseCronField(fields[0], 0, 59)
	if err != nil {
		return CronExpression{}, fmt.Errorf("invalid minute field: %w", err)
	}

	hourBitfield, err := parseCronField(fields[1], 0, 23)
	if err != nil {
		return CronExpression{}, fmt.Errorf("invalid hour field: %w", err)
	}

	dayOfMonthBitfield, err := parseCronField(fields[2], 1, 31)
	if err != nil {
		return CronExpression{}, fmt.Errorf("invalid day of month field: %w", err)
	}

	monthBitfield, err := parseCronField(fields[3], 1, 12)
	if err != nil {
		return CronExpression{}, fmt.Errorf("invalid month field: %w", err)
	}

	dayOfWeekBitfield, err := parseCronField(fields[4], 0, 7)
	if err != nil {
		return CronExpression{}, fmt.Errorf("invalid day of week field: %w", err)
	}

	// day of week is 0-6 where 0 and 7 both represent Sunday
	if (dayOfWeekBitfield & 0x81) != 0 {
		dayOfWeekBitfield |= 0x81
	}

	return CronExpression{
		minute:     minuteBitfield,
		hour:       hourBitfield,
		dayOfMonth: dayOfMonthBitfield,
		month:      monthBitfield,
		dayOfWeek:  dayOfWeekBitfield,
	}, nil
}

func parseCronField(field string, min, max uint64) (uint64, error) {
	bitfield := uint64(0)

	// lists
	parts := strings.Split(field, ",")
	for _, part := range parts {
		// steps
		base, step, hasStep := strings.Cut(part, "/")
		if !hasStep {
			step = "1"
		}

		// wildcard
		if base == "*" {
			base = fmt.Sprintf("%d-%d", min, max)
		}

		// ranges
		minBase, maxBase, hasRange := strings.Cut(base, "-")
		if !hasRange {
			if hasStep {
				// handle ambiguous case like "3/15" where it really means "3-59/15"
				maxBase = strconv.FormatUint(max, 10)
			} else {
				maxBase = minBase
			}
		}

		// convert to integer
		minBaseInt, err := strconv.ParseUint(minBase, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %w", err)
		}
		maxBaseInt, err := strconv.ParseUint(maxBase, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %w", err)
		}
		stepInt, err := strconv.ParseUint(step, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %w", err)
		}

		// validity checks
		if stepInt == 0 {
			return 0, fmt.Errorf("step cannot be zero")
		}

		if minBaseInt > maxBaseInt {
			return 0, fmt.Errorf("min base must be less than max base number")
		}

		// generate part bitfield
		partBitfield := uint64(0)
		for i := minBaseInt; i <= maxBaseInt; i += stepInt {
			if i < min || i > max {
				return 0, fmt.Errorf("value %d out of range [%d-%d]", i, min, max)
			}
			partBitfield |= 1 << i
		}

		bitfield |= partBitfield
	}

	return bitfield, nil
}
