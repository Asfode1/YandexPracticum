package daysteps

import (
	"time"
	"errors"
	"strconv"
	"strings"
	"fmt"
	"log"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)
// parsePackage разбирает строку формата "3456,30m"
// и возвращает количество шагов, длительность и ошибку
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
	if duration <= 0 {
        return 0, 0, errors.New("продолжительность должна быть больше нуля")
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

	// Вычисляем калории и обрабатываем ошибку
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Ошибка при расчёте калорий:", err)
		return ""
	}

	// Формируем строку результата
	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)

	return result
}