package slave

import (
	"runtime"
)

var AllowHandling = true

func save(c chan func() bool) {
	for {
		if AllowHandling == true {
			method := <-c
			method()
		} else {
			runtime.Goexit()
		}
	}
}
