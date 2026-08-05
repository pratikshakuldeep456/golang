package notes

import (
	"encoding/json"
	"fmt"
	"time"
)

// Use json.Encoder and json.Decoder for streaming JSON (I/O).
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Json() {

	tm := time.Date(2001, 12, 1, 11, 1, 2, 0, time.UTC)
	fmt.Println(tm)
	nw := time.Now()
	fmt.Println(nw.Sub(tm))
	fmt.Println(tm.Add(1 * time.Hour))
	fmt.Println()
	u1 := User{Name: "pratiksha",
		Age: 24}
	data, err := json.Marshal(u1)
	if err != nil {
		panic("cant marsal data")
	}

	jsonData := string(data)
	fmt.Println(jsonData)
	data1 := `{"name":"pratiksha","age":24}`
	var u2 User
	json.Unmarshal([]byte(data1), &u2)
	fmt.Printf("expected user is %+v\n", u2)

}
