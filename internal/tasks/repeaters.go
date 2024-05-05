package tasks

import (
	"errors"
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
		return "", errors.New("правило повторения указано в неправильном формате")
	}

	return r.date.Format("20060102"), nil
}

// TODO change here to pointers, no need to throw shadowed strings
func parseRepeater(now time.Time, date string, repeat string) (*repeater, error) {
	var repeatVal [2][]int
	var negativeVal []int
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
				if repeatStr[0] == "d" && val > 400 {
					return nil, errors.New("правило повторения указано в неправильном формате")
				}
				if repeatStr[0] == "m" && val <= -3 || val > 31 {
					return nil, errors.New("правило повторения указано в неправильном формате")
				}
				if repeatStr[0] == "w" && val > 7 {
					return nil, errors.New("правило повторения указано в неправильном формате")
				}
				if val < 0 {
					negativeVal = append(negativeVal, val)
					continue
				}
				repeatVal[i-1] = append(repeatVal[i-1], val)
			}
		}

		for _, v := range repeatVal {
			sort.Ints(v)
		}
		sort.Ints(negativeVal)

		repeatVal[0] = append(repeatVal[0], negativeVal...)
	} else if repeatStr[0] == "y" {
		repeatVal[0] = append(repeatVal[0], 1)
	} else {
		return nil, errors.New("правило повторения указано в неправильном формате")
	}

	if len(repeatVal[1]) == 0 && repeatStr[0] == "m" {
		if repeatVal[0][0] == 31 {
			repeatVal[1] = append(repeatVal[1], 1, 3, 5, 7, 8, 10, 12)
		}
		for i := 1; i <= 12; i++ {
			repeatVal[1] = append(repeatVal[1], i)
		}
	}

	r := &repeater{modifier: repeatStr[0], value: repeatVal, now: now}
	r.date, err = time.Parse("20060102", date)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *repeater) moveDays() {
	if r.now.Before(r.date) {
		r.date = r.date.AddDate(0, 0, r.value[0][0])
	} else {
		for r.date.Before(r.now) {
			r.date = r.date.AddDate(0, 0, r.value[0][0])
		}
	}
}

func (r *repeater) moveYears() {
	if r.now.Before(r.date) {
		r.date = r.date.AddDate(r.value[0][0], 0, 0)
	} else {
		for r.date.Before(r.now) {
			r.date = r.date.AddDate(r.value[0][0], 0, 0)
		}
	}
}

func (r *repeater) moveWeeks() {
	var weekDays string
	for _, weekDay := range r.value[0] {
		if weekDay == 7 {
			weekDay = 0
		}
		weekDays += time.Weekday(weekDay).String()
	}
	if r.date.Before(r.now) {
		r.date = r.now
		r.now = r.now.AddDate(0, 0, 1)
	}

	for r.date.Before(r.now) || !(strings.Contains(weekDays, r.date.Weekday().String())) {
		r.date = r.date.AddDate(0, 0, 1)
	}
}

func (r *repeater) moveMonths() {
	var dates = make([]time.Time, 0)
	var newDate time.Time

	if r.date.Before(r.now) {
		r.date = r.now
	}

	for _, month := range r.value[1] {
		for _, day := range r.value[0] {
			m := month
			if day < 0 {
				m++
				day++
			}
			parseTime := time.Date(r.now.Year(), time.Month(m), day, 0, 0, 0, 0, time.UTC)
			dates = append(dates, parseTime)
		}
	}

	for len(dates) > 0 {
		if r.date.Before(dates[0]) {
			r.date = dates[0]
			return
		} else {
			// Обработка кейсов с високосными годами
			if dates[0].Year()%4 == 0 {
				newDate = dates[0].AddDate(1, 0, -1)
			} else if dates[0].AddDate(1, 0, 0).Year()%4 == 0 {
				newDate = dates[0].AddDate(1, 0, +1)
			} else {
				newDate = dates[0].AddDate(1, 0, 0)
			}
			dates = dates[1:]
			dates = append(dates, newDate)
		}
	}
}
