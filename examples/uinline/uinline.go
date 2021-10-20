package uinline

func throw() {
	panic("panic call function")
}

func inline() {
	throw()
}

func Enter() {
	inline()
}
