package main

type ReservationsSlice []Reservation

var _ ReservationsStorage = &ReservationsSlice{}

func (rs *ReservationsSlice) AddReservation(r Reservation) {
	*rs = append(*rs, r)
}
