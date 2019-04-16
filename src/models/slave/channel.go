package slave

func save(c chan func() bool) {
	for {
		method := <-c
		method()
	}
}
