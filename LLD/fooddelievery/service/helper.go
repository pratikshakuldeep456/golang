// // actually resembles a real system (a DB sequence / snowflake ID in prod).
// type idGenerator struct {
// 	counter int64
// }

// func (g *idGenerator) Next() int {
// 	return int(atomic.AddInt64(&g.counter, 1))
// }

package service
