package interfaces

var i = "test"

type TestWithoutComments interface { // want `Interface 'TestWithoutComments' is missing required headline comment` `Method 'Method' is missing required comment`
	Method() bool
}

// This has a sample comment
type TestWithHeadlineComments interface { // want `Method 'Method' is missing required comment`
	Method() bool
}

// This is a comment
type TestWithComments interface {
	// Comment
	Method() bool
}
