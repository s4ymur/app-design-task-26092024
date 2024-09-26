package main

import (
	"log/slog"
	"reflect"
	"testing"
)

func Test_HandlerSlots_TakeTheWholeSlot(t *testing.T) {
	t.Parallel()

	// arrange
	logger := slog.Default()

	// init reservation slice dependency
	var reserved = ReservationsSlice{}

	timeSlotID := "timeslot1"
	capacty := 3

	var available = []Timeslot{
		{timeSlotID, day(2024, 8, 13, 8), capacty},
	}

	// init handler
	hd := handler{
		logger:         logger,
		resStorage:     &reserved,
		availableSlots: available,
	}

	orderID := "1"

	order := Order{
		RequestID: orderID,
		From:      day(2000, 0, 0, 0),
		To:        day(2070, 0, 0, 0),
		Capacity:  3,
	}

	// act
	err := hd.addReservation(order)

	// assert
	// check err
	if !reflect.DeepEqual(err, nil) {
		t.Fatal("err is not nil")
	}

	// check reservation
	expectedRes := ReservationsSlice{Reservation{RequestID: orderID, TimeslotID: timeSlotID, Capacity: capacty}}
	eq := reflect.DeepEqual(expectedRes, reserved)
	if !eq {
		t.Fatal("reservations are not equal")
	}

	// check time slots
	if available[0].Capacity != 0 {
		t.Fatal("non-zero capacity")
	}
}

func Test_HandlerSlots_TakeTwoSlots_NoRemainingCap(t *testing.T) {
	t.Parallel()

	// arrange
	logger := slog.Default()

	// init reservation slice dependency
	var reserved = ReservationsSlice{}

	timeSlotID1 := "timeslot1"
	timeSlotID2 := "timeslot2"

	firstSlotCapacty := 3
	secondSlotCapacty := 2

	var available = []Timeslot{
		{timeSlotID1, day(2024, 8, 13, 8), firstSlotCapacty},
		{timeSlotID2, day(2024, 9, 13, 8), secondSlotCapacty},
	}

	// init handler
	hd := handler{
		logger:         logger,
		resStorage:     &reserved,
		availableSlots: available,
	}

	orderID := "1"

	order := Order{
		RequestID: orderID,
		From:      day(2000, 0, 0, 0),
		To:        day(2070, 0, 0, 0),
		Capacity:  5,
	}

	// act
	err := hd.addReservation(order)

	// assert
	// check err
	if !reflect.DeepEqual(err, nil) {
		t.Fatal("err is not nil")
	}

	// check reservation
	expectedRes := ReservationsSlice{
		Reservation{
			RequestID: orderID, TimeslotID: timeSlotID1, Capacity: firstSlotCapacty,
		},
		Reservation{
			RequestID: orderID, TimeslotID: timeSlotID2, Capacity: secondSlotCapacty,
		},
	}
	eq := reflect.DeepEqual(expectedRes, reserved)
	if !eq {
		t.Fatal("reservations are not equal")
	}

	// check time slots
	if available[0].Capacity != 0 && available[1].Capacity != 0 {
		t.Fatal("non-zero capacity")
	}

}

func Test_HandlerSlots_TakeTwoSlots_RemainingCap(t *testing.T) {
	t.Parallel()

	// arrange
	logger := slog.Default()

	// init reservation slice dependency
	var reserved = ReservationsSlice{}

	timeSlotID1 := "timeslot1"
	timeSlotID2 := "timeslot2"

	firstSlotCapacty := 3
	secondSlotCapacty := 3

	var available = []Timeslot{
		{timeSlotID1, day(2024, 8, 13, 8), firstSlotCapacty},
		{timeSlotID2, day(2024, 9, 13, 8), secondSlotCapacty},
	}

	// init handler
	hd := handler{
		logger:         logger,
		resStorage:     &reserved,
		availableSlots: available,
	}

	orderID := "1"

	order := Order{
		RequestID: orderID,
		From:      day(2000, 0, 0, 0),
		To:        day(2070, 0, 0, 0),
		Capacity:  5,
	}

	// act
	err := hd.addReservation(order)

	// assert
	// check err
	if !reflect.DeepEqual(err, nil) {
		t.Fatal("err is not nil")
	}

	// check reservation
	expectedRes := ReservationsSlice{
		Reservation{
			RequestID: orderID, TimeslotID: timeSlotID1, Capacity: firstSlotCapacty,
		},
		Reservation{
			RequestID: orderID, TimeslotID: timeSlotID2, Capacity: order.Capacity - firstSlotCapacty,
		},
	}
	eq := reflect.DeepEqual(expectedRes, reserved)
	if !eq {
		t.Fatal("reservations are not equal")
	}

	// check time slots
	if available[0].Capacity != 0 {
		t.Fatal("non-zero capacity")
	}

	if available[1].Capacity != 1 {
		t.Fatal("wrong capacity")
	}

}
