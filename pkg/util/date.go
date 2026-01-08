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

func GetShortcutDate(timeType string, reqStartTime, reqEndTime string) (startTime, endTime string) {
	switch timeType {
	case "today":
		startTime, endTime = GetToday()
	case "yesterday": // 昨天
		startTime = GetYesterday()
		endTime = GetYesterday()
	case "thisWeek": // 本周
		startTime = GetThisWeekStartDate()
		endTime = time.Now().Format(TimeLayoutShort)
	case "prevWeek", "lastWeek": // 上周
		startTime = GetLastWeekStartDate()
		endTime = GetLastWeekEndDate()
	case "thisMonth": // 本月
		startTime = GetThisMonthStartDate()
		endTime = time.Now().Format(TimeLayoutShort)
	case "prevMonth", "lastMonth": // 上月
		startTime = GetLastMonthStartDate()
		endTime = GetLastMonthEndDate()
	case "thisYear":
		year, _, _ := time.Now().Date()
		prevYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
		startTime = prevYear.Format(TimeLayoutShort)
		endTime = prevYear.AddDate(1, 0, 0).Format(TimeLayoutShort)
	case "prevYear":
		year, _, _ := time.Now().Date()
		prevYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
		startTime = prevYear.AddDate(-1, 0, 0).Format(TimeLayoutShort)
		endTime = prevYear.Format(TimeLayoutShort)
	case "custom": // 自定义
		if len(reqStartTime) == 0 || len(reqEndTime) == 0 {
			// 默认最近30天
			subTime := time.Hour * 24 * 30
			startTime = time.Now().Add(-subTime).Format(TimeLayoutShort)
			endTime = time.Now().Format(TimeLayoutShort)
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
