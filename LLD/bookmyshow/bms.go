package bookmyshow

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrMovieNotFound    = errors.New("not found")
	ErrShowTimeNotFound = errors.New(" show not found")
	ErrSeatTimeNotFound = errors.New("seat not found")
	ErrSeatNotAvailable = errors.New("seat not available")
)

const TIMEWINDOW int = 5

const SEATPRICE int32 = 100

type BookingSvc struct {
	//list of movies
	//list of venue
	// tickets
	//payments
	Movies    map[int32]*Movie
	Venue     map[int32]*Venue
	Tickets   map[int32]*Ticket
	ShowTimes map[int32]*ShowTime
	mu        sync.RWMutex
}

var Instance *BookingSvc
var once sync.Once

func BMSInstance() *BookingSvc {
	once.Do(func() {

		Instance = &BookingSvc{Movies: make(map[int32]*Movie),
			Venue:     make(map[int32]*Venue),
			Tickets:   make(map[int32]*Ticket),
			ShowTimes: make(map[int32]*ShowTime)}
	})
	return Instance
}
func (bs *BookingSvc) InsertTestData() {

	Movie1 := &Movie{ID: 1, Name: "spider_man"}
	Movie2 := &Movie{ID: 2, Name: "ABCD"}
	Movie3 := &Movie{ID: 3, Name: "SFIE"}
	Movie4 := &Movie{ID: 4, Name: "KER"}
	bs.Movies[Movie1.ID] = Movie1
	bs.Movies[Movie2.ID] = Movie2
	bs.Movies[Movie3.ID] = Movie3
	bs.Movies[Movie4.ID] = Movie4

	V1 := &Venue{ID: 1, Location: "hsr"}
	V2 := &Venue{ID: 2, Location: "kora"}
	V3 := &Venue{ID: 3, Location: "indranagar"}
	V4 := &Venue{ID: 4, Location: "sjr"}
	bs.Venue[V1.ID] = V1
	bs.Venue[V2.ID] = V2
	bs.Venue[V3.ID] = V3
	bs.Venue[V4.ID] = V4

	loc, _ := time.LoadLocation("Asia/Kolkata")
	ShowTime1 := &ShowTime{
		ID:      IDGenerator(),
		MovieID: 1,
		VenueID: 1,
		StartAt: time.Date(2026, 8, 20, 12, 0, 0, 0, loc),
		Seats:   make(map[int32]*Seat),
	}
	bs.ShowTimes[ShowTime1.ID] = ShowTime1

	seat1 := &Seat{ID: IDGenerator(),
		Row:    "A",
		Number: 1,
		Status: AVAILABLE,
	}
	seat2 := &Seat{ID: IDGenerator(),
		Row:    "A",
		Number: 2,
		Status: AVAILABLE,
	}
	bs.ShowTimes[ShowTime1.ID].Seats[seat2.ID] = seat2
	bs.ShowTimes[ShowTime1.ID].Seats[seat1.ID] = seat1

}
func (bs *BookingSvc) SearchMovie(text string) error {
	//full search
	for _, j := range bs.Movies {
		if j.Name == text {
			fmt.Printf("Movie name: %s, Movie id: %d\n", j.Name, j.ID)
		}
	}

	return nil
}
func (bs *BookingSvc) validateMovieExists(movieID int32) (*Movie, error) {
	movie, exists := bs.Movies[movieID]
	if !exists {
		return nil, ErrMovieNotFound
	}
	return movie, nil
}
func (bs *BookingSvc) ViewShowTimes(movieID int32) {
	_, err := bs.validateMovieExists(movieID)
	if err != nil {
		fmt.Println("error", err)
	}

	for _, val := range bs.ShowTimes {
		if val.MovieID == movieID {
			fmt.Println("show details are", val.ID, val.Seats, val.MovieID)
			for _, j := range val.Seats {
				fmt.Println(j)
			}
		}
	}
}

func (bs *BookingSvc) BookSeat(userID int32, showTimeID int32, seatIDs []int32) (*Ticket, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	//check seat

	showTime, exists := bs.ShowTimes[showTimeID]
	if !exists {
		return nil, ErrShowTimeNotFound
	}
	var selectedSeats []*Seat

	for _, sid := range seatIDs {
		seat, exists := showTime.Seats[sid]
		if !exists {
			return &Ticket{}, ErrSeatTimeNotFound
		}
		if seat.Status == RESERVED {
			//check time
			oldTicket, ok := bs.Tickets[seat.TicketID]
			if ok && time.Since(oldTicket.ReservedAt) > time.Duration(TIMEWINDOW)*time.Minute {
				seat.Status = AVAILABLE
				seat.TicketID = 0
				oldTicket.Status = FAILED
			}
		}
		if seat.Status != AVAILABLE {
			return nil, ErrSeatNotAvailable
		}
		selectedSeats = append(selectedSeats, seat)
		//
	}
	ticket := &Ticket{
		ID:         IDGenerator(),
		UserID:     userID,
		ShowTime:   showTime.StartAt,
		Status:     PENDING,
		Seats:      selectedSeats,
		ReservedAt: time.Now(),
	}

	//update seat..
	for _, seat := range selectedSeats {
		seat.Status = RESERVED
		seat.TicketID = ticket.ID
	}

	bs.Tickets[ticket.ID] = ticket
	return ticket, nil

}

func (bs *BookingSvc) ConfirmBooking(tID int32, pg PaymentGateway) (string, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	ticket, exists := bs.Tickets[tID]
	if !exists {
		return "", errors.New("ticket not presnet in system")
	}
	if ticket.Status != PENDING {
		return "", errors.New("invalid req")
	}
	//check expirry
	//seats := bs.Tickets[tID].Seats
	seats := ticket.Seats
	if time.Since(ticket.ReservedAt) > time.Duration(TIMEWINDOW)*time.Minute {
		// reservation expired...
		bs.Tickets[tID].Status = CANCELLED
		for i := 0; i < len(seats); i++ {
			bs.Tickets[tID].Seats[i].Status = AVAILABLE
		}
		return "", errors.New("booknig window exceeded")
	}
	// check seats
	amount := 0
	for i := 0; i < len(seats); i++ {
		if seats[i].Status != RESERVED {
			return "", errors.New("seat is not reserved")

		} else {
			amount += int(SEATPRICE)
			continue
		}
	}

	// make payment
	pg.Pay(float32(amount), tID)

	//seats := ticket.Seats
	for i := 0; i < len(seats); i++ {
		bs.Tickets[tID].Seats[i].Status = BOOKED
	}
	bs.Tickets[tID].Status = CONFIRMED

	return "BOOKNIG confirmed", nil
}
