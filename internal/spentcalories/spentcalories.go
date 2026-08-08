package spentcalories

import (
	// standart
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// "3456,Ходьба,3h00m"
	var steps int
	var dur time.Duration
	var err error

	dataSplit := strings.Split(data, ",")
	if len(dataSplit) != 3 {
		return 0, "", 0, fmt.Errorf("wrong pars num: %d", len(dataSplit))
	}

	// steps
	steps, err = strconv.Atoi(dataSplit[0])
	if err != nil {
		return 0, "", 0, err
	} else if steps < 1 {
		return 0, "", 0, fmt.Errorf("wrong value for steps: %d", steps)
	}

	// activity
	// dataSplit[1]

	// duration
	dur, err = time.ParseDuration(dataSplit[2])
	if err != nil {
		return 0, "", 0, err
	} else if dur < 1 {
		return 0, "", 0, fmt.Errorf("wrong value for dur: %d", dur)
	}

	return steps, dataSplit[1], dur, nil
}

func distance(steps int, height float64) float64 {
	var dist float64

	if steps < 1 {
		return 0
	}
	if height <= 0 {
		return 0
	}

	// dist
	stepLength := height * stepLengthCoefficient
	dist = (stepLength * float64(steps)) / mInKm

	// result
	return dist
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	var speed float64

	if steps < 1 {
		return 0
	}
	if height <= 0 {
		return 0
	}
	if duration < 1 {
		return 0
	}

	// distance
	dist := distance(steps, height)

	// speed
	speed = dist / float64(duration.Hours())

	// result
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	//"3456,Ходьба,3h00m"
	// trainType, dur
	steps, trainType, dur, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if weight <= 0 {
		err = fmt.Errorf("wrong value for weight: %f", weight)
		log.Println(err)
		return "", err
	}
	if height <= 0 {
		err = fmt.Errorf("wrong value for height: %f", height)
		log.Println(err)
		return "", err
	}

	// calories
	var calErr error
	var cal float64
	switch trainType {
	case "Бег":
		cal, calErr = RunningSpentCalories(steps, weight, height, dur)
	case "Ходьба":
		cal, calErr = WalkingSpentCalories(steps, weight, height, dur)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
	if calErr != nil {
		log.Println(calErr)
		return "", calErr
	}

	// distance
	dist := distance(steps, height)

	// speed
	speed := meanSpeed(steps, height, dur)

	// result
	return fmt.Sprintf(`Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`, trainType, dur.Hours(), dist, speed, cal), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 1 {
		return 0, fmt.Errorf("wrong value for steps: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("wrong value for weight: %f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("wrong value for height: %f", height)
	}
	if duration < 1 {
		return 0, fmt.Errorf("wrong value for duration: %d", duration)
	}

	speed := meanSpeed(steps, height, duration)

	// calories
	cal := (weight * speed * float64(duration.Minutes())) / minInH

	// result
	return cal, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// calories
	cal, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	cal *= walkingCaloriesCoefficient

	// result
	return cal, nil
}
