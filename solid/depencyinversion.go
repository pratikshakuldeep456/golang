package solid

import "fmt"

// 1. Define the ABSTRACTION that UserService actually needs
type Database interface {
	Save(data string)
}

//  2. MySQLDatabase implements that abstraction — it's still low-level,
//     but now it conforms to a contract UserService can depend on
type MySQLDatabase struct{}

func (m *MySQLDatabase) Save(data string) {
	fmt.Println("Saving to MySQL:", data)
}

// 3. UserService now depends on the INTERFACE, not the concrete struct
type UserService struct {
	db Database
}

//  4. The concrete implementation is INJECTED from outside,
//     NewUserService no longer creates its own dependency
func NewUserService(db Database) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(data string) {
	s.db.Save(data)
}

func main() {
	mysql := &MySQLDatabase{}
	service := NewUserService(mysql) // dependency injected from the caller
	service.Register("alice@example.com")
}
