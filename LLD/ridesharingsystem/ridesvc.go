package ridesharingsystem

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type RideService struct {
	Riders      map[int]*Rider
	Rides       map[int]*Ride
	Driver      map[int]*Driver
	RideHistory map[int][]*Ride
}

// instance of ridesvc
var RideInstance *RideService
var once sync.Once

func Instance() *RideService {
	once.Do(func() {
		RideInstance = &RideService{
			Riders:      make(map[int]*Rider),
			Rides:       make(map[int]*Ride),
			Driver:      make(map[int]*Driver),
			RideHistory: make(map[int][]*Ride),
		}
	})
	return RideInstance
}

func (rs *RideService) RegisterRider(name, phoneNo string, loc *Location) (*Rider, error) {
	rider := &Rider{
		ID:       int(generateRiderID()),
		Name:     name,
		PhoneNo:  phoneNo,
		Location: loc,
	}
	rs.Riders[rider.ID] = rider
	fmt.Println("rider registered")
	return rider, nil
}
func (rs *RideService) RegisterDriver(name, phoneNo string, loc *Location, rType RideType) (*Driver, error) {
	driver := &Driver{
		ID: int(generateDriverID()),
		//	Name:     name,
		PhoneNo:  phoneNo,
		Location: loc,
		RideType: rType,
		Status:   AVAILABLE,
	}
	fmt.Println("driver registered")
	rs.Driver[driver.ID] = driver
	return driver, nil
}

func (rs *RideService) RequestRide(rId int, src, des *Location, rType RideType) (int32, error) {
	//check rider exists
	if _, exits := rs.Riders[rId]; !exits {
		return 0, ErrRiderNotPresentInSystem
	}
	fmt.Println("ride requested")

	//calcu distance
	// calcu distance
	// notify drivers one at a time

	//create ride
	rider := rs.Riders[rId]
	ride1 := &Ride{
		ID:          int(generateRideID()),
		RequestedBy: rider,
		Src:         src,
		Des:         des,
		Status:      REQUESTED,
		Fare:        0,
		RideType:    rType,
		AcceptedBy:  nil,
		CreatedAt:   time.Now(),
	}
	fmt.Println("ride stored")
	rs.Rides[ride1.ID] = ride1
	return int32(ride1.ID), nil

}

func (rs *RideService) AcceptRide(rid int, dId int, rType RideType) (string, error) {
	if _, exits := rs.Rides[rid]; !exits {
		return "", ErrRiderNotPresentInSystem
	}

	if _, exists := rs.Driver[dId]; !exists {
		return "", ErrDriverNotPresentInSystem
	}

	driver := rs.Driver[dId]
	// check status of driver
	if rs.Driver[dId].Status != AVAILABLE || rType != rs.Driver[dId].RideType {
		return "", ErrRiderBusyOrDifferentRideType
	}
	// if action == ACCEPTED {
	//update ride status
	rs.Rides[rid].AcceptedBy = driver
	rs.Driver[dId].Status = BUSY

	//}
	return "", nil

}
func (rs *RideService) CalculateFare(rID int32) int32 {
	if _, exits := rs.Riders[int(rID)]; !exits {
		return 0
	}
	dis := rs.CalculateDistance(rs.Rides[int(rID)].Src, rs.Rides[int(rID)].Des)
	srateaty := GetPricingStrategy(rs.Rides[int(rID)].RideType)
	fare := srateaty.CalculateFare(rs.Rides[int(rID)].RideType, int32(dis))
	rs.Rides[int(rID)].Fare = fare
	fmt.Println("CALCUALTED FAIR IS", fare)
	return fare

}

func (rs *RideService) CalculateDistance(driver, rider *Location) int {

	return int(math.Sqrt(math.Pow((driver.Lat-rider.Lat), 2) + math.Pow((driver.Long-rider.Long), 2)))
}

func (rs *RideService) FindDriver(rLoc *Location, rType RideType) *Driver {
	var driver *Driver
	var bestDis int = math.MaxInt
	for _, v := range rs.Driver {
		fmt.Println("driver fetcing", v.ID)
		if v.Status == AVAILABLE && rType == v.RideType {
			fmt.Println("driver fetcing", v.ID)
			//calc distance
			best := rs.CalculateDistance(v.Location, rLoc)
			if best <= DistanceLimit {
				if bestDis > best {
					bestDis = best
					driver = v
				} else {
					continue
				}
			} else {
				continue
			}

		}
	}
	fmt.Println("driver we got is", driver.ID)
	return driver
}
