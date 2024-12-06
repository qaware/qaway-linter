package functions

import (
	"log"
)

func shortMethodNoComments() bool { // want `Method 'shortMethodNoComments' is missing required headline comment` `Method 'shortMethodNoComments' has less than 10% comment density. Actual: 0%` `Method 'shortMethodNoComments' has less than 10% logging density. Actual: 0%`
	return true
}

func longMethodNoComments() bool { // want `is missing required headline comment` `Method 'longMethodNoComments' has less than 10% comment density. Actual: 0%` `Method 'longMethodNoComments' has less than 10% logging density. Actual: 0%`
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
func longMethodWithComments() bool { // want `Method 'longMethodWithComments' has less than 10% logging density. Actual: 0%`
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

/*
Lorem ipsum dolor sit amet, consetetur sadipscing elitr,
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.
Lorem ipsum dolor sit amet, consetetur sadipscing elitr,
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.
At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren,
no sea takimata sanctus est Lorem ipsum dolor sit amet.
*/
func longMethodWithSingleCommentGroup() bool { // want `Method 'longMethodWithSingleCommentGroup' has less than 10% logging density. Actual: 0%`
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
func longMethodWithInlineComments() bool { // want `Method 'longMethodWithInlineComments' has less than 10% logging density. Actual: 0%`
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

// downloads the artifacts
func downloadArtifacts() bool { // want `Method 'downloadArtifacts' has a trivial comment. Similarity to method name: 67%` `Method 'downloadArtifacts' has less than 10% logging density. Actual: 0%`
	// method has a trivial comment
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

// this function has a nice headline comment
// explains everything that is needed
func perfectMethod() bool {
	log.Printf("Hello World")
	s := "abc"
	s += "1"

	// a simple comment
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

	log.Printf("Hello World")
	return false
}
