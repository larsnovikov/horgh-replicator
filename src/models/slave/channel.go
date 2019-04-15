package slave

func save(c chan func() bool) {
	for {
		method := <-c
		if method() == true {
			// TODO update pos
		}
	}
}
