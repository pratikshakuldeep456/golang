package solid

type User struct {
	Name string
	Age  int
}
type AgeStruct struct {
	Age int
}

// func (u *User) GetAge() int {
// 	return u.Age
// }

type category string

const (
	SMALL  category = "small"
	MEDIUM category = "medium"
	LARGE  category = "large"
)

func (u *User) CheckCategory() (string, error) {

	if u.Age < 18 {
		return string(SMALL), nil
	}
	return string(MEDIUM), nil
}
