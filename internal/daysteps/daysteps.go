package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65 // Средняя длина шага в метрах
	mInKm      = 1000 // Метров в одном километре
)

// parsePackage парсит строку вида "678,0h50m"
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных, ожидалось 2 поля, получено %d", len(parts))
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	// Запрет пробелов: строка должна совпадать со своим TrimSpace
	if stepsStr != strings.TrimSpace(stepsStr) || durationStr != strings.TrimSpace(durationStr) {
		return 0, 0, fmt.Errorf("пробелы в данных недопустимы")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверное значение количества шагов: %w", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат продолжительности: %w", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, duration, nil
}

// DayActionInfo возвращает информацию о дневной активности
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("ошибка обработки данных: %v", err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKilometers := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps, distanceKilometers, calories,
	)
	return result
}

// part of the final project for Sprint
