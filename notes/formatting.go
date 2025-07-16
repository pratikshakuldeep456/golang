package notes

import "fmt"

func FormatString() {
	name := "Go"
	str := fmt.Sprintf("Hello %s", name)
	fmt.Println(str)
	count := 3
	flt := 3.5
	fmt.Printf("name is %s with len of %d and %f with a type is %T\n", name, count, flt, name)
	u := struct {
		Name string
		Age  int
	}{
		"Alice",
		25,
	}
	fmt.Printf("struct vall is %v and all details are %+v", u, u)

}
