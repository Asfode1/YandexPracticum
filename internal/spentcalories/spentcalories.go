package spentcalories

import (
	"time"
	"errors"
	"strconv"
	"strings"
	"fmt"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining принимает строку формата "3456,Ходьба,3h00m"
// и возвращает количество шагов, вид активности, продолжительность и ошибку.
	func parseTraining(data string) (int, string, time.Duration, error) {
		parts := strings.Split(data, ",")
		if len(parts) != 3 {
			return 0, "", 0, errors.New("неверный формат: должно быть три значения, разделённых запятыми")
		}
	
		// Парсим количество шагов
		steps, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, "", 0, errors.New("ошибка при разборе количества шагов: " + err.Error())
		}
	
		// Вид активности
		activity := parts[1]
	
		// Парсим продолжительность
		duration, err := time.ParseDuration(parts[2])
		if err != nil {
			return 0, "", 0, errors.New("ошибка при разборе продолжительности: " + err.Error())
		}
	
		return steps, activity, duration, nil
	}

// distance вычисляет дистанцию в километрах по количеству шагов и росту
func distance(steps int, height float64) float64 {
	// Рассчитываем длину одного шага
	stepLength := height * stepLengthCoefficient
	
	// Вычисляем общую дистанцию в метрах
	distanceMeters := float64(steps) * stepLength
	
	// Переводим в километры
	return distanceMeters / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
    if duration <= 0 {
        return 0
    }
    
    distanceKm := distance(steps, height)
    durationHours := duration.Hours()
    
    return distanceKm / durationHours
}

// RunningSpentCalories рассчитывает количество калорий, потраченных при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше нуля")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше нуля")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше нуля")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше нуля")
	}

	// Вычисляем среднюю скорость (км/ч)
	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, errors.New("средняя скорость должна быть больше нуля")
	}

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Рассчитываем калории по формуле
	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}


func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше нуля")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше нуля")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше нуля")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше нуля")
	}

	// Вычисляем среднюю скорость (км/ч)
	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, errors.New("средняя скорость должна быть больше нуля")
	}

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Рассчитываем калории по формуле
	calories := (weight * speed * durationMinutes) / minInH

	// Умножаем на корректирующий коэффициент для ходьбы
	calories *= walkingCaloriesCoefficient

	return calories, nil
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

	switch strings.ToLower(activity) {
	case "бег", "running":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "ходьба", "walking":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки: " + activity)
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity,
		duration.Hours(),
		dist,
		speed,
		calories,
	)

	return result, nil
}