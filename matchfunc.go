package mockaka

import "github.com/go-test/deep"

// MatchFunc is a function that matches stubs for request.
// Return value of zero means full match.
type MatchFunc[I protomsg] func(actual I) int

// NewDefaultMatchFunc returns default match function.
func NewDefaultMatchFunc[I protomsg](expected I) MatchFunc[I] {
	return func(actual I) int {
		return len(deep.Equal(expected, actual))
	}
}

// MatchAny is a special matcher that matches any request.
func MatchAny[I any](actual I) int {
	return 0
}
