type PricingStrategy interface {
	CalculateFee(vType VehicleType, entryTime, exitTime time.Time) int32
}

type TimeBasedPricing struct {
	RatePerHour map[VehicleType]int32
}

func (t *TimeBasedPricing) CalculateFee(vType VehicleType, entryTime, exitTime time.Time) int32 {
	duration := exitTime.Sub(entryTime)
	hours := int32(duration.Hours())
	if duration.Minutes() > float64(hours*60) {
		hours++ // round up partial hour — 43 mins still charges 1 full hour
	}
	if hours == 0 {
		hours = 1 // minimum charge
	}
	return hours * t.RatePerHour[vType]
}

func (pl *ParkingLot) SetPricingStrategy(s PricingStrategy) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.Pricing = s
}