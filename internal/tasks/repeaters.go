package tasks

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type repeater struct {
	modifier string
	value    int
	currDate time.Time
	now      time.Time
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	r, err := parseRepeater(now, date, repeat)
	if err != nil {
		return "", err
	}

	switch r.modifier {
	case "d":
		r.moveDay()
	case "y":
		r.moveYear()
	default:
		return "", errors.New("invalid repeat symbol")
	}

	//now = const 26012024
	//date = current recursive date
	//repeat = days or months to add to date until it will be greater now
	//TODO date + repeate string until it will greater then now

	return r.currDate.String(), nil
}

func parseRepeater(now time.Time, date string, repeat string) (*repeater, error) {
	var repeatVal int
	var err error

	repeatStr := strings.Split(repeat, " ")
	if len(repeatStr) > 1 {
		repeatVal, err = strconv.Atoi(repeatStr[1])
		if err != nil {
			return nil, err
		}
	} else if repeatStr[0] == "y" {
		repeatVal = 1
	} else {
		return nil, errors.New("invalid repeat value")
	}

	if repeatVal > 400 {
		return nil, errors.New("invalid repeat value")
	}

	r := &repeater{modifier: repeatStr[0], value: repeatVal, now: now}
	r.currDate, err = time.Parse("20060102", date)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *repeater) moveDay() {
	for r.currDate.Before(r.now) {
		r.currDate = r.currDate.AddDate(0, 0, r.value)
	}
}

func (r *repeater) moveYear() {
	for r.currDate.Before(r.now) {
		r.currDate = r.currDate.AddDate(r.value, 0, 0)
	}
}

func (r *repeater) moveMonth() {

}
