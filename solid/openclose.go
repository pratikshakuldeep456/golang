package solid

type Shape interface {
	Area() float64
}

type Rectange struct {
	LENGTH  float64
	BREADTH float64
}

func (r *Rectange) Area() float64 {
	return r.LENGTH * r.BREADTH
}

type Circle struct {
	RADIUS float64
}

func (c *Circle) Area() float64 {
	return 3.14 * c.RADIUS * c.RADIUS
}

func GetArea(s Shape) float64 {
	return s.Area()
}

//////////////////

type Account struct {
}
