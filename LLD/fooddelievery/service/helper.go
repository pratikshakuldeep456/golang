// // actually resembles a real system (a DB sequence / snowflake ID in prod).
// type idGenerator struct {
// 	counter int64
// }

//	func (g *idGenerator) Next() int {
//		return int(atomic.AddInt64(&g.counter, 1))
package service

import "sync/atomic"

var (
	restaurantIDCounter int32
	menuItemIDCounter   int32
	cartItemIDCounter   int32
	orderIDCounter      int32
	orderItemIDCounter  int32
	paymentIDCounter    int32
	userIDCounter       int32
)

func generateRestaurantID() int { return int(atomic.AddInt32(&restaurantIDCounter, 1)) }
func generateMenuItemID() int   { return int(atomic.AddInt32(&menuItemIDCounter, 1)) }
func generateCartItemID() int   { return int(atomic.AddInt32(&cartItemIDCounter, 1)) }
func generateOrderID() int      { return int(atomic.AddInt32(&orderIDCounter, 1)) }
func generateOrderItemID() int  { return int(atomic.AddInt32(&orderItemIDCounter, 1)) }
func generatePaymentID() int    { return int(atomic.AddInt32(&paymentIDCounter, 1)) }
func generateUserID() int       { return int(atomic.AddInt32(&userIDCounter, 1)) }
