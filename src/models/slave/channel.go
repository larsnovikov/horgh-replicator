package slave

var AllowHandling = true

func save(c chan func() bool) {
	for {
		if AllowHandling == true {
			method := <-c
			method()
		}
	}
}
