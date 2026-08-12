package ridesharingsystem

type PriceInterface interface {
	CalculateFare(rType RideType, distance int32) int32
}

type CabStrategy struct {
}

func (cs *CabStrategy) CalculateFare(rType RideType, distance int32) int32 {
	return 2 * distance
}

type AutoStrategy struct {
}

func (as *AutoStrategy) CalculateFare(rType RideType, distance int32) int32 {
	return distance
}

type PremiumPricing struct{}

func (p *PremiumPricing) CalculateFare(distance int) int32 {
	return int32(distance) * 25 // rate per km for premium
}

func GetPricingStrategy(rType RideType) PriceInterface {
	switch rType {
	case AUTO:
		return &AutoStrategy{}
	case CAB:
		return &CabStrategy{}
	// case PREMIUM:
	// 	return &PremiumPricing{}
	default:
		return &CabStrategy{}
	}
}
