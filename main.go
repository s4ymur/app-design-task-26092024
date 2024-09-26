// Ниже реализован сервис бронирования таймслотов для доставки еды. В предметной области
// существуют два основных понятия:
//
// Timeslot — конкретный промежуток времени и максимальное кол-во
// заказов в этот промежуток. Если заказ, занимает больше кол-во
// времени, чем 1 промежуток (например 3 часа на сбор и доставку) значит, нам необходимо
// забронировать Capacity из 2 промежутов сразу 8:00 - 9:59 и 10:00 - 11:59.
//
// Reservation — это лог всех зарезервированных таймслотов. Для конечного понимания, что за запрос
// зарезервировал время. Тут должны находится только актуальные данные.
// Это очень важно для наших коллег из отдела аналитики!
//
// Задание:
// - провести рефакторинг кода с выделением слоев и абстракций
// - применить best-practices там где это имеет смысл
// - исправить имеющиеся в реализации логические и технические ошибки и неточности
package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
)

const (
	errServerClosed = "Server closed"
	errServerFailed = "Server failed"
)

func main() {
	// init logger
	logger := slog.Default()

	// init reservation slice dependency
	var reserved = ReservationsSlice{
		{"request1", "timeslot4", 1},
	}

	var available = []Timeslot{
		{"timeslot1", day(2024, 8, 13, 8), 3},
		{"timeslot2", day(2024, 8, 13, 10), 2},
		{"timeslot3", day(2024, 8, 13, 12), 2},
		{"timeslot4", day(2024, 8, 13, 14), 0},
	}

	// init handler
	hd := handler{
		logger:         logger,
		resStorage:     &reserved,
		availableSlots: available,
	}

	server := http.NewServeMux()

	server.HandleFunc("/slots", hd.HandlerSlots)

	logger.Info("Server listening on localhost:8080")

	err := http.ListenAndServe(":8080", server)
	if errors.Is(err, http.ErrServerClosed) {
		logger.Error(errServerClosed)
	} else if err != nil {
		logger.Warn("%s err: %s", errServerFailed, err.Error())

		os.Exit(1)
	}
}
