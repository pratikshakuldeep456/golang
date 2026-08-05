package notes

import (
	"fmt"
	"net/url"
)

func Parse() {
	s := "https://example.com/search?q=golang&sort=desc"
	u, _ := url.Parse(s)
	fmt.Println(u.Path, u.Host, u.Scheme, u.Query())

}
