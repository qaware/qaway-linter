package interfaces

var i = "test"

type TestWithoutComments interface { // want `Interface 'TestWithoutComments' is missing required headline comment` `Method 'Method' is missing required comment`
	Method() bool
}

// TODO: should not count
type TestWithTodoComment interface { // want `Interface 'TestWithTodoComment' is missing required headline comment`

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

//
//nolint:qawaylinter
type TestWithEmptyCommentAndDirective interface { // want `Interface 'TestWithEmptyCommentAndDirective' is missing required headline comment` `Method 'Method' is missing required comment`
	//
	//nolint:qawaylinter
	Method() bool
}
