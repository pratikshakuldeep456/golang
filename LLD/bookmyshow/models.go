package bookmyshow

import "time"

type User struct {
	ID    int32
	Name  string
	Email string
}

type Movie struct {
	ID   int32
	Name string
	//Location  string
	//ShowTimes map[int]*ShowTime
}
type Venue struct {
	ID int32
	//ShowTimes map[int]*ShowTime
	//Movies   map[int]*Movie
	Location string
}

type ShowTime struct {
	ID       int32
	MovieID  int32
	VenueID  int32
	StartAt  time.Time
	Duration int32
	Seats    map[int32]*Seat
}

//	type Seat struct {
//		ID map[int32]SeatStatus
//	}
type Seat struct {
	ID       int32
	Row      string
	Number   int32
	Status   SeatStatus
	TicketID int32
}

type SeatStatus string

const (
	BOOKED    SeatStatus = "booked"
	AVAILABLE SeatStatus = "available"
	RESERVED  SeatStatus = "reseved"
)

type TicketStatus string

const (
	PENDING   TicketStatus = "pending"
	CONFIRMED TicketStatus = "confirmed"
	CANCELLED TicketStatus = "cancelled"
	FAILED    TicketStatus = "failed"
)

type Ticket struct {
	ID         int32
	UserID     int32
	ShowTime   time.Time
	Status     TicketStatus
	Seats      []*Seat
	ReservedAt time.Time
}
