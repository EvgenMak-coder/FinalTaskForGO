package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

var (
	ErrEmptyRepeat      = errors.New("empty repeat rule")
	ErrInvalidFormat    = errors.New("invalid repeat format")
	ErrInvalidDate      = errors.New("invalid date")
	ErrInvalidDay       = errors.New("invalid day")
	ErrInvalidMonth     = errors.New("invalid month")
	ErrInvalidWeekday   = errors.New("invalid weekday")
	ErrIntervalTooLarge = errors.New("interval exceeds maximum value")
)

func NextDate(now time.Time, dStart string, repeat string) (string, error) {
	if repeat == "" {
		return "", ErrEmptyRepeat
	}

	date, err := time.Parse(dateFormat, dStart)
	if err != nil {
		return "", ErrInvalidDate
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", ErrInvalidFormat
	}

	switch parts[0] {
	case "d":
		return nextDateDaily(now, date, parts)
	case "y":
		return nextDateYearly(now, date)
	case "w":
		return nextDateWeekly(now, date, parts)
	case "m":
		return nextDateMonthly(now, date, parts)
	default:
		return "", ErrInvalidFormat
	}
}

func nextDateDaily(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", ErrInvalidFormat
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", ErrInvalidFormat
	}

	if interval <= 0 || interval > 400 {
		return "", ErrIntervalTooLarge
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
}

func nextDateYearly(now time.Time, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
}

func nextDateWeekly(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", ErrInvalidFormat
	}

	weekdays := strings.Split(parts[1], ",")
	days := make(map[int]bool, len(weekdays))

	for _, wd := range weekdays {
		d, err := strconv.Atoi(wd)
		if err != nil || d < 1 || d > 7 {
			return "", ErrInvalidWeekday
		}
		days[d] = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) {
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if days[weekday] {
				break
			}
		}
	}

	return date.Format(dateFormat), nil
}

func nextDateMonthly(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", ErrInvalidFormat
	}

	days := parseDays(parts[1])
	if len(days) == 0 {
		return "", ErrInvalidDay
	}

	var months []int
	if len(parts) == 3 {
		months = parseMonths(parts[2])
		if len(months) == 0 {
			return "", ErrInvalidMonth
		}
	}

	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) {
			day := date.Day()
			month := int(date.Month())

			lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
			dayCheck := day
			if contains(days, -1) && day == lastDay {
				dayCheck = -1
			} else if contains(days, -2) && day == lastDay-1 {
				dayCheck = -2
			}

			if contains(days, dayCheck) && (len(months) == 0 || contains(months, month)) {
				break
			}
		}
	}

	return date.Format(dateFormat), nil
}

func parseDays(s string) []int {
	parts := strings.Split(s, ",")
	var days []int

	for _, p := range parts {
		d, err := strconv.Atoi(p)
		if err == nil {
			if d == -1 || d == -2 || (d >= 1 && d <= 31) {
				days = append(days, d)
			}
		}
	}

	return days
}

func parseMonths(s string) []int {
	parts := strings.Split(s, ",")
	var months []int

	for _, p := range parts {
		m, err := strconv.Atoi(p)
		if err == nil && m >= 1 && m <= 12 {
			months = append(months, m)
		}
	}

	return months
}

func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func afterNow(a, b time.Time) bool {
	return a.Year() > b.Year() ||
		(a.Year() == b.Year() && a.YearDay() >= b.YearDay())
}
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "Invalid 'now' date format", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}
