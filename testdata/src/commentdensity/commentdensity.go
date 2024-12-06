package commentdensity

func shortMethodNoComments() bool {
	return true
}

func longMethodNoComments() bool {
	s := "abc"
	s += "1"

	if s == "" {
		return true
	}

	s += "2"
	if s == "" {
		return true
	}

	s += "3"
	if s == "" {
		return true
	}

	s += "4"
	if s == "" {
		return true
	}

	return false
}

// Lorem ipsum dolor sit amet, consetetur sadipscing elitr,
// sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
// sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.
// Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.
// Lorem ipsum dolor sit amet, consetetur sadipscing elitr,
// sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
// sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.
// Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.
func longMethodWithComments() bool {
	s := "abc"
	s += "1"

	if s == "" {
		return true
	}

	s += "2"
	if s == "" {
		return true
	}

	s += "3"
	if s == "" {
		return true
	}

	s += "4"
	if s == "" {
		return true
	}

	return false
}

// mininmal header
func longMethodWithInlineComments() bool {
	s := "abc"
	s += "1"

	// Lorem ipsum dolor sit amet, consetetur sadipscing elitr,
	// sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
	// sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.
	// Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.
	// Lorem ipsum dolor sit amet, consetetur sadipscing elitr,
	// sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
	// sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.
	// Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.
	if s == "" {
		return true
	}

	s += "2"
	if s == "" {
		return true
	}

	s += "3"
	if s == "" {
		return true
	}

	s += "4"
	if s == "" {
		return true
	}

	return false
}
