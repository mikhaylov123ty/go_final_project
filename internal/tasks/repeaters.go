package tasks

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	returnError = "правило повторения указано в неправильном формате"
)

// структура для модификатора повторений
type repeater struct {
	modifier string    // модификатор повторения d, y, w, m
	value    [2][]int  // значения модификаторов повторения
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

// TODO change here to pointers, no need to throw shadowed strings
func parseRepeater(now time.Time, date string, repeat string) (*repeater, error) {
	var repeatVal [2][]int
	var negativeVal []int
	var err error

	// Парсинг строки с условием повторения
	repeatStr := strings.Split(repeat, " ")

	// Обработка входных данных, если найдены значения у модификаторов повторения
	if len(repeatStr) > 1 {
		for i := 1; i < len(repeatStr); i++ {
			repeatStrVal := strings.Split(repeatStr[i], ",")
			if repeatStrVal[0] == "x" {
				return nil, errors.New(returnError)

			}
			// Проверка каждого из значений условия повторения
			for _, v := range repeatStrVal {
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
				if val < 0 {
					negativeVal = append(negativeVal, val)
					continue
				}
				repeatVal[i-1] = append(repeatVal[i-1], val)
			}
		}

		// Сортировка массива, для вывода упорядоченных дат проведения задачи
		for _, v := range repeatVal {
			sort.Ints(v)
		}

		// Добавление в конце к массиву отсортированных отрицательных значений
		// в любых сценариях они указывают на предпоследний и последний дни месяца
		sort.Ints(negativeVal)
		repeatVal[0] = append(repeatVal[0], negativeVal...)

		// Обработка сценария, когда модификатору не требуется значение (смена года)
	} else if repeatStr[0] == "y" {
		repeatVal[0] = append(repeatVal[0], 1)

		// Возврат ошибки, если указан неверный формат
	} else {
		return nil, errors.New(returnError)
	}

	// Обработка условия, если выбран модификатор повторения месяц
	// без указания конкретных месяцев
	if len(repeatVal[1]) == 0 && repeatStr[0] == "m" {

		// Если указано 31е число, то добавляем только те месяцы, где 31 день
		if repeatVal[0][0] == 31 {
			repeatVal[1] = append(repeatVal[1], 1, 3, 5, 7, 8, 10, 12)
		}

		// В остальных случаях добавляем все 12 месцев
		for i := 1; i <= 12; i++ {
			repeatVal[1] = append(repeatVal[1], i)
		}
	}

	// Заполняем обработанными данными экземпляр повторения
	r := &repeater{modifier: repeatStr[0], value: repeatVal, now: now}
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
		r.date = r.date.AddDate(0, 0, r.value[0][0])
	} else {
		// В противном случае добавляем дни, пока дата проведения не станет больше текущей
		for r.date.Before(r.now) {
			r.date = r.date.AddDate(0, 0, r.value[0][0])
		}
	}
}

// Метод поиска следующей даты проведения в сценарии с годами
func (r *repeater) moveYears() {
	// Если текущая дата меньше даты проведения, добавляем года к дате проведения
	if r.now.Before(r.date) {
		r.date = r.date.AddDate(r.value[0][0], 0, 0)
	} else {
		// В противном случае добавляем года, пока дата проведения не станет больше текущей
		for r.date.Before(r.now) {
			r.date = r.date.AddDate(r.value[0][0], 0, 0)
		}
	}
}

// Метод поиска следующей даты проведения в сценарии с днями недели
func (r *repeater) moveWeeks() {
	var weekDays string
	// Для каждого из значений дня недели, формируем строку с названиями этих дней
	for _, weekDay := range r.value[0] {
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
	var dates = make([]time.Time, 0)
	var newDate time.Time

	// Если дата проведения меньше текущей даты, устанавливаем дату проведения на текущую дату
	if r.date.Before(r.now) {
		r.date = r.now
	}

	// Составление всех возможных пар дня и месяца
	// Выходной список в отсортированном порядке
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

	// Проверка возможных дат проведения
	for len(dates) > 0 {

		// Если текущая дата проведения меньше, то устанавливается ближайшая дата
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
				// Если текущая дата проведения больше, добавляем к ближайшей дате год
				newDate = dates[0].AddDate(1, 0, 0)
			}

			// Исключаем ближайшую дату из слайса и добавляем в конец эту же дату на следующий год
			dates = dates[1:]
			dates = append(dates, newDate)
		}
	}
}
