package elevatorsystem

import (
	"sync"
)

type ElevatorSystem struct {
	Elevators []*Elevator
	MinFloor  int
	MaxFloor  int
}

var ESInstance *ElevatorSystem
var once sync.Once

func Instance() *ElevatorSystem {
	once.Do(func() {
		elevators := make([]*Elevator, 0, 3)
		for i := 0; i < 3; i++ {
			elevators = append(elevators, &Elevator{
				ID:           i,
				Status:       IDLE,
				direction:    UP,
				CurrentFloor: 0,
				UpStops:      make(map[int]bool),
				DownStops:    make(map[int]bool),
			})
		}
		ESInstance = &ElevatorSystem{
			Elevators: elevators,
			MinFloor:  0,
			MaxFloor:  10,
		}
	})
	return ESInstance
}

func (es *ElevatorSystem) AddRequest(direction ELEVATORDIRECTION, requestAt int) (string, error) {
	data, err := es.IsValidReq(requestAt)
	if !data {
		return "", err
	}
	//
	elevator := es.FindBestElevator(direction, requestAt)
	if elevator == nil {
		return "", ErrNoElevatorAvailable
	}
	elevator.addStop(requestAt)
	return "", nil
}

func (es *ElevatorSystem) FindBestElevator(direction ELEVATORDIRECTION, requestAt int) *Elevator {
	var bestElevator *Elevator
	bestDistance := -1
	for i := 0; i < len(es.Elevators); i++ {
		eligible := false
		distance := 0
		e := es.Elevators[i]
		// if es.Elevators[i].Status == IDLE {
		// 	eligible = true
		// 	distance = abs(requestAt - es.Elevators[i].CurrentFloor)
		// }
		if e.Status == IDLE {
			eligible = true
			distance = abs(requestAt - e.CurrentFloor)
		}

		if eligible && (bestElevator == nil || distance < bestDistance) {
			bestDistance = distance
			bestElevator = e
		}
	}

	return bestElevator

}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func (es *ElevatorSystem) RequestFloor(floor int, eid int) (string, error) {
	valid, err := es.IsValidReq(floor)
	if !valid {
		return "", err
	}
	if eid < 0 || eid >= len(es.Elevators) {
		return "", ErrInvalidElevator
	}
	elevator := es.Elevators[eid]
	elevator.addStop(floor)
	return "", nil
}
func (es *ElevatorSystem) ServeRequat() {

	//stop
	//serve
	//move

}

func (e *Elevator) addStop(floor int) {
	if floor == e.CurrentFloor {
		return
	}
	if floor > e.CurrentFloor {
		e.UpStops[floor] = true
	} else {
		e.DownStops[floor] = true
	}
}
func (es *ElevatorSystem) IsValidReq(floor int) (bool, error) {
	if es.MinFloor <= floor && es.MaxFloor >= floor {
		return true, nil
	}
	return false, InValidRequest
}
