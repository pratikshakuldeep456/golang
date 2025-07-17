package notes

import (
	"bufio"
	"fmt"
	"os"
)

func ReadWrite() {
	path := "/Users/pratikshakuldeep/Documents/targeting-engine/README.md"
	rdData, _ := os.ReadFile(path)
	fmt.Println(string(rdData))
	os.OpenFile(path, os.O_CREATE, 0644)

	f, _ := os.Create("text.txt")
	w := bufio.NewWriter(f)
	w.WriteString(string(rdData))
	fmt.Println(f)

}
