package tasks

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	returnError = "правило повторения указано в неправильном формате"
)

// Структура для модификатора повторений
type repeater struct {
	modifier string    // модификатор повторения d, y, w, m
	days     []int     // значение модификатора дни
	months   []int     // значение модификатора месяцы
	date     time.Time // дата задачи
	now      time.Time // текущая дата
}

// Основной обработчик для поиска следующей даты повторений
func NextDateHandler(now time.Time, date string, repeat string) (string, error) {

	// Парсинг условия повторения
	r, err := parseRepeater(now, date, repeat)
	if err != nil {
		return "", err
	}

	// Распределение повторения
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
		return "", errors.New(returnError)
	}

	return r.date.Format("20060102"), nil
}

// Метод для парсинга правил повторения
func parseRepeater(now time.Time, date string, repeat string) (*repeater, error) {
	var repeatDays []int
	var repeatMonths []int
	var err error

	// Парсинг строки с условием повторения
	repeatStr := strings.Split(repeat, " ")

	// Обработка входных данных, если найдены значения у модификаторов повторения
	if len(repeatStr) > 1 {

		repeatDaysVal := strings.Split(repeatStr[1], ",")

		if repeatDaysVal[0] == "x" {
			return nil, errors.New(returnError)
		}

		// Проверка каждого из значений условия повторения и исключение аномалий
		for _, v := range repeatDaysVal {
			val, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			if repeatStr[0] == "d" && val > 400 {
				return nil, errors.New(returnError)
			}
			if repeatStr[0] == "m" && (val <= -3 || val > 31) {
				return nil, errors.New(returnError)
			}
			if repeatStr[0] == "w" && val > 7 {
				return nil, errors.New(returnError)
			}
			repeatDays = append(repeatDays, val)
		}

		if len(repeatStr) > 2 {
			repeatMonthsVal := strings.Split(repeatStr[2], ",")
			for _, v := range repeatMonthsVal {
				val, err := strconv.Atoi(v)
				if err != nil {
					return nil, err
				}
				repeatMonths = append(repeatMonths, val)
			}
		}

		// Обработка сценария, когда модификатору не требуется значение (смена года)
	} else if repeatStr[0] == "y" {
		repeatDays = append(repeatDays, 1)

		// Возврат ошибки, если указан неверный формат
	} else {
		return nil, errors.New(returnError)
	}

	// Обработка условия, если выбран модификатор повторения месяц
	// без указания конкретных месяцев
	if len(repeatMonths) == 0 && repeatStr[0] == "m" {

		// Если указано 31е число, то добавляем только те месяцы, где 31 день
		if repeatDays[0] == 31 {
			repeatMonths = append(repeatMonths, 1, 3, 5, 7, 8, 10, 12)
		} else {
			// В остальных случаях добавляем все 12 месцев
			for i := 1; i <= 12; i++ {
				repeatMonths = append(repeatMonths, i)
			}
		}
	}

	// Заполняем обработанными данными экземпляр повторения
	r := &repeater{modifier: repeatStr[0], days: repeatDays, months: repeatMonths, now: now}
	r.date, err = time.Parse("20060102", date)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Метод поиска следующей даты проведения в сценарии с днями
func (r *repeater) moveDays() {

	// Если текущая дата меньше даты проведения, добавляем дни к дате проведения
	if r.now.Before(r.date) {
		r.date = r.date.AddDate(0, 0, r.days[0])
		return
	}

	// В противном случае добавляем дни, пока дата проведения не станет больше текущей
	for r.date.Before(r.now) {
		r.date = r.date.AddDate(0, 0, r.days[0])
	}
}

// Метод поиска следующей даты проведения в сценарии с годами
func (r *repeater) moveYears() {

	// Если текущая дата меньше даты проведения, добавляем года к дате проведения
	if r.now.Before(r.date) {
		r.date = r.date.AddDate(r.days[0], 0, 0)
		return
	}

	// В противном случае добавляем года, пока дата проведения не станет больше текущей
	for r.date.Before(r.now) {
		r.date = r.date.AddDate(r.days[0], 0, 0)
	}
}

// Метод поиска следующей даты проведения в сценарии с днями недели
func (r *repeater) moveWeeks() {
	var weekDays string
	// Для каждого из значений дня недели, формируем строку с названиями этих дней
	for _, weekDay := range r.days {
		if weekDay == 7 {
			weekDay = 0
		}
		weekDays += time.Weekday(weekDay).String()
	}

	// Если дата проведения меньше текущей даты, устанавливаем дату проведения на текущую дату
	if r.date.Before(r.now) {
		r.date = r.now
		r.now = r.now.AddDate(0, 0, 1)
	}

	// Добавляем дни пока дата проведения меньше текущей и не будет найден ближайший день проведения
	for r.date.Before(r.now) || !(strings.Contains(weekDays, r.date.Weekday().String())) {
		r.date = r.date.AddDate(0, 0, 1)
	}
}

// Метод поиска следующей даты проведения в сценарии с месяцами
func (r *repeater) moveMonths() {
	var baseDates = make([]time.Time, 0)
	//var newDate time.Time

	// Если дата проведения меньше текущей даты, устанавливаем дату проведения на текущую дату
	if r.date.Before(r.now) {
		r.date = r.now
	}

	// Составление всех возможных базовых пар дня и месяца
	// Выходной список в отсортированном порядке
	baseDates = createDatesSlice(r.now.Year(), r.days, r.months)
	fmt.Println(r.days, r.months)
	i := 1
	for {
		fmt.Println(baseDates)
		for range len(baseDates) {
			if r.date.Before(baseDates[0]) {
				r.date = baseDates[0]
				return
			}
			baseDates = baseDates[1:]
		}
		baseDates = createDatesSlice(r.now.Year()+i, r.days, r.months)
		i++
	}

}

type timeSlice []time.Time

func (s timeSlice) Less(i, j int) bool {
	return s[i].Before(s[j])
}
func (s timeSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s timeSlice) Len() int {
	return len(s)
}

func createDatesSlice(year int, days []int, months []int) []time.Time {
	var dates timeSlice = []time.Time{}
	for _, month := range months {
		for _, day := range days {
			m := month
			if day < 0 {
				day++
				m++
			}
			parseTime := time.Date(year, time.Month(m), day, 0, 0, 0, 0, time.UTC)
			dates = append(dates, parseTime)
		}
	}

	sort.Sort(dates)
	return dates
}
