package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	stepLengthCoefficient      = 0.45 // коэффициент для расчёта длины шага на основе роста
	mInKm                      = 1000 // метров в километре
	minutesInHour              = 60   // минут в часе
	walkingCaloriesCoefficient = 0.5  // коэффициент для калорий при ходьбе
)

// ActivityType задаёт тип активности
type ActivityType string

const (
	ActivityRunning ActivityType = "running"
	ActivityWalking ActivityType = "walking"
)

// parseTraining разбирает строку формата "3456,Ходьба,3h00m"
// и возвращает количество шагов, тип активности, продолжительность и ошибку
func parseTraining(data string) (steps int, activity ActivityType, duration time.Duration, err error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		err = errors.New("неверный формат: должно быть три значения, разделённых запятыми")
		return
	}

	steps, err = strconv.Atoi(parts[0])
	if err != nil {
		err = fmt.Errorf("ошибка при разборе количества шагов: %w", err)
		return
	}
	if steps <= 0 {
		err = errors.New("количество шагов должно быть больше нуля")
		return
	}

	activity = ActivityType(strings.ToLower(parts[1]))

	duration, err = time.ParseDuration(parts[2])
	if err != nil {
		err = fmt.Errorf("ошибка при разборе продолжительности: %w", err)
		return
	}
	if duration <= 0 {
		err = errors.New("продолжительность должна быть больше нуля")
		return
	}

	return
}

// validateParams проверяет базовые параметры для расчёта калорий
func validateParams(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 {
		return errors.New("количество шагов должно быть больше нуля")
	}
	if weight <= 0 {
		return errors.New("вес должен быть больше нуля")
	}
	if height <= 0 {
		return errors.New("рост должен быть больше нуля")
	}
	if duration <= 0 {
		return errors.New("продолжительность должна быть больше нуля")
	}
	return nil
}

// distance вычисляет дистанцию в километрах по шагам и росту
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return float64(steps)*stepLength / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

// caloriesSpent рассчитывает калории с учётом коэффициента активности
func caloriesSpent(steps int, weight, height float64, duration time.Duration, coefficient float64) (float64, error) {
	if err := validateParams(steps, weight, height, duration); err != nil {
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, errors.New("средняя скорость должна быть больше нуля")
	}

	durationMinutes := duration.Minutes()
	calories := (weight * speed * durationMinutes) / minutesInHour
	calories *= coefficient

	return calories, nil
}

// RunningSpentCalories рассчитывает калории при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	return caloriesSpent(steps, weight, height, duration, 1.0)
}

// WalkingSpentCalories рассчитывает калории при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	return caloriesSpent(steps, weight, height, duration, walkingCaloriesCoefficient)
}

// TrainingInfo формирует строку с информацией о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64
	switch activity {
	case ActivityRunning:
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case ActivityWalking:
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		strings.Title(string(activity)),
		duration.Hours(),
		dist,
		speed,
		calories,
	)

	return result, nil
}