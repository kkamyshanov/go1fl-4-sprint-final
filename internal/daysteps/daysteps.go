package daysteps

import (
	// standart
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	// internal
	sptcal "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// "678,0h50m"
	dataSplit := strings.Split(data, ",")
	if len(dataSplit) != 2 {
		return 0, 0, fmt.Errorf("wrong pars num: %d", len(dataSplit))
	}

	// steps
	steps, err := strconv.Atoi(dataSplit[0])
	if err != nil {
		return 0, 0, err
	} else if steps < 1 {
		return 0, 0, fmt.Errorf("wrong value for steps: %d", steps)
	}

	//duration
	dur, err := time.ParseDuration(dataSplit[1])
	if err != nil {
		return 0, 0, err
	} else if dur < 1 {
		return 0, 0, fmt.Errorf("wrong value for dur: %d", dur)
	}

	// result
	return steps, dur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// steps, duration
	// if steps < 1, we already have err != nil
	steps, dur, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	// distance
	distance := (float64(steps) * stepLength) / mInKm
	cal, err := sptcal.WalkingSpentCalories(steps, weight, height, dur)
	if err != nil {
		log.Println(err)
		return ""
	}

	// result
	return fmt.Sprintf(`Количество шагов: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.
`, steps, distance, cal)
}
