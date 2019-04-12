package slave

func save(c chan func() bool) {
	method := <-c
	if method() == true {
		// TODO update pos
	}
}
