package tasks

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type repeater struct {
	modifier string
	value    [2][]int
	date     time.Time
	now      time.Time
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	r, err := parseRepeater(now, date, repeat)
	if err != nil {
		return "", err
	}

	fmt.Println(r)

	switch r.modifier {
	case "d":
		r.moveDays()
	case "y":
		r.moveYears()
	case "w":
		r.moveWeeks()
	case "m":
		r.moveMonths()
	default:
		return "", errors.New("invalid repeat symbol")
	}

	return r.date.String(), nil
}

func parseRepeater(now time.Time, date string, repeat string) (*repeater, error) {
	var repeatVal [2][]int
	var err error

	repeatStr := strings.Split(repeat, " ")
	if len(repeatStr) > 1 {
		for i := 1; i < len(repeatStr); i++ {
			repeatStrVal := strings.Split(repeatStr[i], ",")
			for _, v := range repeatStrVal {
				val, err := strconv.Atoi(v)
				if err != nil {
					return nil, err
				}
				repeatVal[i-1] = append(repeatVal[i-1], val)
			}
		}
	} else if repeatStr[0] == "y" {
		repeatVal[0] = append(repeatVal[0], 1)
	} else {
		return nil, errors.New("invalid repeat value")
	}

	if repeatVal[0][0] > 400 {
		return nil, errors.New("invalid repeat value")
	}

	if repeatVal[0][0] <= -3 {
		return nil, errors.New("invalid repeat value")
	}

	if len(repeatVal[1]) == 0 {
		for i := 1; i <= 12; i++ {
			repeatVal[1] = append(repeatVal[1], i)
		}
	}

	for _, v := range repeatVal {
		sort.Ints(v)
	}

	r := &repeater{modifier: repeatStr[0], value: repeatVal, now: now}
	r.date, err = time.Parse("20060102", date)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *repeater) moveDays() {
	for r.date.Before(r.now) {
		r.date = r.date.AddDate(0, 0, r.value[0][0])
	}
}

func (r *repeater) moveYears() {
	for r.date.Before(r.now) {
		r.date = r.date.AddDate(r.value[0][0], 0, 0)
	}
}

func (r *repeater) moveWeeks() {
	var weekDays string
	for _, weekDay := range r.value[0] {
		weekDays += time.Weekday(weekDay).String()
	}

	for r.date.Before(r.now) || !(strings.Contains(weekDays, r.date.Weekday().String())) {
		r.date = r.date.AddDate(0, 0, 1)
	}
}

func (r *repeater) moveMonths() {
	dates := make([]time.Time, 0)

	if r.date.Before(r.now) {
		r.date = r.now
	}

	for _, month := range r.value[1] {
		for _, day := range r.value[0] {
			m := month
			if day < 0 {
				m++
			}
			parseTime := time.Date(r.now.Year(), time.Month(m), day, 0, 0, 0, 0, time.UTC)
			dates = append(dates, parseTime)
		}
	}

	for len(dates) > 0 {
		if r.date.Before(dates[0]) {
			r.date = dates[0]
			fmt.Println("ANSWER:", r.date)
			return
		} else {
			newDate := dates[0].AddDate(1, 0, 0)
			dates = dates[1:]
			dates = append(dates, newDate)
		}
	}
}
