package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	mInKm                      = 1000
	minInH                     = 60
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
	runCaloriesCoefficient     = 1.0
	cyclingCaloriesCoefficient = 8.0
)

// parseTraining парсит строку вида "3456,Ходьба,3h00m"
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных, ожидалось 3 поля, получено %d", len(parts))
	}

	stepsStr := parts[0]
	activity := parts[1]
	durationStr := parts[2]

	if strings.ContainsAny(stepsStr, " \t\n\r") ||
		strings.ContainsAny(activity, " \t\n\r") ||
		strings.ContainsAny(durationStr, " \t\n\r") {
		return 0, "", 0, fmt.Errorf("пробелы в данных недопустимы")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, activity, duration, nil // ← activity — строка, НЕ рост!
}

// distance вычисляет дистанцию в км на основе роста
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / duration.Hours() // ← КЛЮЧЕВО: .Hours(), не .Seconds()
}

// WalkingSpentCalories — калории при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("некорректные параметры: вес и рост должны быть положительные")
	}
	if steps <= 0 {
		return 0, fmt.Errorf("некорректные параметры: количество шагов должно быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры: продолжительновсть должна быть больше 0")
	}

	meanSpeedValue := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := walkingCaloriesCoefficient * (weight * meanSpeedValue * durationInMinutes) / minInH
	return calories, nil
}

// RunningSpentCalories — калории при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("некорректные параметры: вес и рост должны быть положительные")
	}
	if steps <= 0 {
		return 0, fmt.Errorf("некорректные параметры: количество шагов должно быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры: продолжительность должна быть больше 0")
	}

	meanSpeedValue := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := runCaloriesCoefficient * (weight * meanSpeedValue * durationInMinutes) / minInH
	return calories, nil
}

// TrainingInfo — информация о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64
	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Велосипед":
		calories = cyclingCaloriesCoefficient * weight * duration.Hours()
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity, duration.Hours(), dist, speed, calories,
	)
	return result, nil
}

// 2 part of final task sprint 4
