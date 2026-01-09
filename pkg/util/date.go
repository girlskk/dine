package util

import (
	"encoding/json"
	"time"

	"github.com/samber/lo"
)

// 当天开始时间：00：00：00
func DayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 当天结束时间：23：59：59
func DayEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

// nowFunc is a package-level function that returns the current time. It is
// set to time.Now by default but can be overridden in tests to make outputs
// deterministic.
var nowFunc = time.Now

type RequestDate time.Time

func (r *RequestDate) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*r = RequestDate(time.Time{})
		return nil
	}
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	if s == "" {
		*r = RequestDate(time.Time{})
		return nil
	}
	t, err := time.ParseInLocation(time.DateOnly, s, time.Local)
	if err != nil {
		return err
	}
	*r = RequestDate(t)
	return nil
}

func (r RequestDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(r).Format(time.DateOnly))
}

func (r RequestDate) ToTime() time.Time {
	return time.Time(r)
}

func (r RequestDate) ToPtrTime() *time.Time {
	if r.ToTime().IsZero() {
		return nil
	}
	return lo.ToPtr(time.Time(r))
}

func (r RequestDate) IsValid() bool {
	return !r.ToTime().IsZero()
}

func (r RequestDate) ToStartOfDay() time.Time {
	return DayStart(r.ToTime())
}

func (r RequestDate) ToEndOfDay() time.Time {
	return DayEnd(r.ToTime())
}

func (r RequestDate) ToPtrStartOfDay() *time.Time {
	if r.ToTime().IsZero() {
		return nil
	}
	return lo.ToPtr(DayStart(r.ToTime()))
}

func (r RequestDate) ToPtrEndOfDay() *time.Time {
	if r.ToTime().IsZero() {
		return nil
	}
	return lo.ToPtr(DayEnd(r.ToTime()))
}

const (
	TimeLayoutShort = "2006-01-02"
)

func ParseDateToPtr(dateStr string) (*time.Time, error) {
	t, err := time.ParseInLocation(TimeLayoutShort, dateStr, time.Local)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func GetShortcutDate(timeType string, reqStartTime, reqEndTime string) (startTime, endTime string) {
	switch timeType {
	case "today":
		startTime, endTime = GetToday()
	case "yesterday": // 昨天
		startTime = GetYesterday()
		endTime = GetYesterday()
	case "thisWeek": // 本周
		startTime = GetThisWeekStartDate()
		endTime = nowFunc().Format(TimeLayoutShort)
	case "prevWeek", "lastWeek": // 上周
		startTime = GetLastWeekStartDate()
		endTime = GetLastWeekEndDate()
	case "thisMonth": // 本月
		startTime = GetThisMonthStartDate()
		endTime = nowFunc().Format(TimeLayoutShort)
	case "prevMonth", "lastMonth": // 上月
		startTime = GetLastMonthStartDate()
		endTime = GetLastMonthEndDate()
	case "thisYear":
		year, _, _ := nowFunc().Date()
		prevYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
		startTime = prevYear.Format(TimeLayoutShort)
		endTime = prevYear.AddDate(1, 0, 0).Format(TimeLayoutShort)
	case "prevYear":
		year, _, _ := nowFunc().Date()
		prevYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
		startTime = prevYear.AddDate(-1, 0, 0).Format(TimeLayoutShort)
		endTime = prevYear.Format(TimeLayoutShort)
	case "custom": // 自定义
		if len(reqStartTime) == 0 || len(reqEndTime) == 0 {
			// 默认最近30天
			subTime := time.Hour * 24 * 30
			startTime = nowFunc().Add(-subTime).Format(TimeLayoutShort)
			endTime = nowFunc().Format(TimeLayoutShort)
		} else {
			startTime = reqStartTime
			endTime = reqEndTime
		}
	default:
		startTime = GetLastWeekStartDate()
		endTime = GetLastWeekEndDate()
	}
	return startTime, endTime
}

// GetToday returns start and end date strings for today (both the same date).
func GetToday() (string, string) {
	t := nowFunc()
	s := t.Format(TimeLayoutShort)
	return s, s
}

// GetYesterday returns the date string for yesterday.
func GetYesterday() string {
	return nowFunc().AddDate(0, 0, -1).Format(TimeLayoutShort)
}

// GetThisWeekStartDate returns the date string for the start of the current week.
// Week starts on Monday.
func GetThisWeekStartDate() string {
	now := nowFunc()
	weekday := int(now.Weekday())
	// Convert so Monday = 0, Sunday = 6
	offset := (weekday + 6) % 7
	start := now.AddDate(0, 0, -offset)
	return start.Format(TimeLayoutShort)
}

// GetLastWeekStartDate returns the date string for the start (Monday) of last week.
func GetLastWeekStartDate() string {
	now := nowFunc()
	weekday := int(now.Weekday())
	offset := (weekday + 6) % 7
	start := now.AddDate(0, 0, -offset-7)
	return start.Format(TimeLayoutShort)
}

// GetLastWeekEndDate returns the date string for the end (Sunday) of last week.
func GetLastWeekEndDate() string {
	now := nowFunc()
	weekday := int(now.Weekday())
	offset := (weekday + 6) % 7
	start := now.AddDate(0, 0, -offset-7)
	end := start.AddDate(0, 0, 6)
	return end.Format(TimeLayoutShort)
}

// GetThisMonthStartDate returns the date string for the first day of the current month.
func GetThisMonthStartDate() string {
	now := nowFunc()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	return start.Format(TimeLayoutShort)
}

// GetLastMonthStartDate returns the date string for the first day of the previous month.
func GetLastMonthStartDate() string {
	now := nowFunc()
	firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	start := firstOfThisMonth.AddDate(0, -1, 0)
	return start.Format(TimeLayoutShort)
}

// GetLastMonthEndDate returns the date string for the last day of the previous month.
func GetLastMonthEndDate() string {
	now := nowFunc()
	firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	end := firstOfThisMonth.AddDate(0, 0, -1)
	return end.Format(TimeLayoutShort)
}
