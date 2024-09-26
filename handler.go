package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type ReservationsStorage interface {
	AddReservation(res Reservation)
}

type handler struct {
	logger         *slog.Logger
	resStorage     ReservationsStorage
	availableSlots []Timeslot
}

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

type Reservation struct {
	RequestID  string `json:"request_id"`
	TimeslotID string `json:"timeslot_id"`
	Capacity   int    `json:"capacity"`
}

func (hr *handler) HandlerSlots(w http.ResponseWriter, r *http.Request) {
	// decode order
	var rqOrder Order

	err := json.NewDecoder(r.Body).Decode(&rqOrder)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		hr.logger.Error(err.Error())

		return
	}

	// check request params
	if rqOrder.From.After(rqOrder.To) {
		w.WriteHeader(http.StatusBadRequest)

		_, err := w.Write([]byte("Error: 'From' date cannot be after 'To' date"))
		if err != nil {
			hr.logger.Error(err.Error())
		}

		return
	}

	// range over available
	err = hr.addReservation(rqOrder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	// return slot in response
	w.Header().Add("content-type", "application/json")

	err = json.NewEncoder(w).Encode(rqOrder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	hr.logger.Info("operation completed")
}

func (hr *handler) addReservation(rqOrder Order) error {
	capacityToShare := rqOrder.Capacity

	for i, curr := range hr.availableSlots {
		if !isBetweenDates(curr.Date, rqOrder.From, rqOrder.To) {
			continue
		}

		// order (or its remains) fits in the slot capacity
		if curr.Capacity >= capacityToShare {
			curr.Capacity -= capacityToShare
			hr.availableSlots[i] = curr

			hr.resStorage.AddReservation(Reservation{RequestID: rqOrder.RequestID, TimeslotID: curr.ID, Capacity: capacityToShare})

			hr.logger.Info("capacity reserved id: %s, cap: %d", curr.ID, capacityToShare)

			return nil
		}

		// order DOES NOT fit in the slot capacity
		hr.resStorage.AddReservation(Reservation{RequestID: rqOrder.RequestID, TimeslotID: curr.ID, Capacity: curr.Capacity})

		hr.logger.Info("capacity reserved id: %s, cap: %d", curr.ID, curr.Capacity)

		capacityToShare -= curr.Capacity
		curr.Capacity = 0 // no more capacity available
		hr.availableSlots[i] = curr
	}

	if capacityToShare != 0 {
		return errors.New("not enough time slots for order")
	}

	return nil
}

func isBetweenDates(t, start, end time.Time) bool {
	// Check if t is after start and before end
	return (t.After(start) || t.Equal(start)) && (t.Before(end) || t.Equal(end))
}

func day(y, m, d, h int) time.Time {
	return time.Date(y, time.Month(m), d, h, 0, 0, 0, time.UTC)
}
