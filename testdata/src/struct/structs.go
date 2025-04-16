package _struct

// Test is a struct
type Test struct { // want `Field 'MissingComment' is missing required comment`
	MissingComment string
}

type TestWithoutComment struct { // want `Struct 'TestWithoutComment' is missing required headline comment`
	// Comment
	MethodWithComment string
}

//
//nolint:qawaylinter
type TestWithEmptyCommentAndDirective struct { // want `Struct 'TestWithEmptyCommentAndDirective' is missing required headline comment`
}
