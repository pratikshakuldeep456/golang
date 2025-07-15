package notes

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func Errors() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// your app logic
}
