package daysteps

import (
	"time"
	"errors"
	"strconv"
	"strings" 
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
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
}
