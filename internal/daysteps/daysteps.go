package daysteps

import (
	"time"
	"errors"
	"strconv"
	"strings"
	"fmt"
	"log"
	"github.com/Asfode1/YandexPracticum/internal/spentcalories"
	
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат входных данных: должна быть одна запятая между шагами и продолжительностью")
	}

	// Парсим количество шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.New("ошибка при разборе количества шагов: " + err.Error())
	}
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше нуля")
	}
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка при разборе продолжительности: %v", err)
	}

	return steps, duration, nil
}
// DayActionInfo парсит данные, вычисляет дистанцию и калории, возвращает форматированную строку
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка парсинга данных:", err)
		return ""
	}

	if steps <= 0 {
		return ""
	}
	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength
	// Переводим в километры
	distanceKm := distanceMeters / mInKm

	// Вычисляем калории
	calories := spentcalories.WalkingSpentCalories(weight, height, duration)

	// Формируем строку результата
	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distanceKm, calories,
	)

	return result
}
