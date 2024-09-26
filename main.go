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
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Order struct {
	RequestID string    `json:"request_id"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	Capacity  int       `json:"capacity"`
}

type Timeslot struct {
	ID       string    `json:"id"`
	Date     time.Time `json:"date"`
	Capacity int       `json:"capacity"`
}

var available = []Timeslot{
	{"timeslot1", day(2024, 8, 13, 8), 3},
	{"timeslot2", day(2024, 8, 13, 10), 2},
	{"timeslot3", day(2024, 8, 13, 12), 2},
	{"timeslot4", day(2024, 8, 13, 14), 0},
}

type Reservation struct {
	RequestID  string `json:"request_id"`
	TimeslotID string `json:"timeslot_id"`
	Capacity   int    `json:"capacity"`
}

var reserved = []Reservation{
	{"request1", "timeslot4", 1},
}

func main() {
	server := http.NewServeMux()
	server.HandleFunc("/slots", handler)

	LogInfo("Server listening on localhost:8080", map[string]string{})

	err := http.ListenAndServe(":8080", server)
	if errors.Is(err, http.ErrServerClosed) {
		LogInfo("Server closed", map[string]string{})
	} else if err != nil {
		LogInfo("Server failed: %s", map[string]string{"err": err.Error()})
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var slot Order
	json.NewDecoder(r.Body).Decode(&slot)

	if slot.From.After(slot.To) {
		w.WriteHeader(200)
		return
	}

	for i, curr := range available {
		if curr.Date.After(slot.From) && curr.Date.Before(slot.To) {
			if curr.Capacity-slot.Capacity < 0 {
				w.WriteHeader(500)
				return
			} else {
				curr.Capacity -= slot.Capacity
				available[i] = curr
				reserved = append(
					reserved,
					Reservation{RequestID: slot.RequestID, TimeslotID: curr.ID, Capacity: slot.Capacity},
				)
			}

			LogInfo("capacity reserved", map[string]string{"id": curr.ID, "cap": fmt.Sprintf("%d", slot.Capacity)})
		}
	}

	LogInfo("operation completed", map[string]string{})
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(slot)
	w.WriteHeader(204)
}

func day(y, m, d, h int) time.Time {
	return time.Date(y, time.Month(m), d, h, 0, 0, 0, time.UTC)
}

var log = slog.Logger{}

func LogInfo(msg string, args map[string]string) {
	var additional string

	for k, v := range args {
		additional += fmt.Sprintf("%s: %s, ", k, v)
	}

	slog.Info(msg, "args", additional)
}
