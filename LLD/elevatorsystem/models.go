package elevatorsystem

type ELEVATORDIRECTION string

const (
	UP   ELEVATORDIRECTION = "up"
	DOWN ELEVATORDIRECTION = "down"
)

type STATE string

const (
	IDLE   STATE = "idle"
	MOVING STATE = "moving"
)

//	type Floor struct {
//		floors []int
//	}
type Elevator struct {
	ID int
	// Floors       *Floor
	Status       STATE
	direction    ELEVATORDIRECTION
	CurrentFloor int
	UpStops      map[int]bool // pending stops above current floor
	DownStops    map[int]bool // pending stops below current floor
}
