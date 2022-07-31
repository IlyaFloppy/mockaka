package mockaka

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Stub contains request or matcher and corresponding response.
type Stub interface {
	Times(n int) Stub
	Name(name string) Stub

	typ(req reflect.Type, res reflect.Type) bool
	resp() (any, error)
	schema() string
	use() Stub
	check(t *testing.T)
	match(method string, request any) int
}

type protomsg interface {
	ProtoMessage()
	ProtoReflect() protoreflect.Message
}

// NewStub creates new stub. It will be used if request matches exactly.
func NewStub[I, O protomsg](method string, request I, response O) Stub {
	return stub[I, O]{
		method:    method,
		request:   request,
		response:  response,
		matchfunc: NewDefaultMatchFunc(request),
		times:     1,
		used:      0,
	}
}

// NewMatchStub creates a new stub using provided MatchFunc.
func NewMatchStub[I, O protomsg](method string, requestMatchFunc MatchFunc[I], response O) Stub {
	return stub[I, O]{
		method:    method,
		response:  response,
		matchfunc: requestMatchFunc,
		times:     1,
		used:      0,
	}
}

// NewErrStub creates new stub that will return an error.
// It will be used if request matches exactly.
// NewErrStub demands explicit type parameters.
func NewErrStub[I, O protomsg](method string, request I, err error) Stub {
	return stub[I, O]{
		method:    method,
		request:   request,
		err:       err,
		matchfunc: NewDefaultMatchFunc(request),
		times:     1,
		used:      0,
	}
}

// NewMatchErrStub creates a new stub that will return an error.
// It uses provided MatchFunc for matching.
// NewMatchErrStub demands explicit type parameters.
func NewMatchErrStub[I, O protomsg](method string, requestMatchFunc MatchFunc[I], err error) Stub {
	return stub[I, O]{
		method:    method,
		err:       err,
		matchfunc: requestMatchFunc,
		times:     1,
		used:      0,
	}
}

type stub[I, O protomsg] struct {
	method   string
	request  I
	response O
	err      error

	matchfunc MatchFunc[I]

	name  string
	times int
	used  int
}

var _ Stub = stub[protomsg, protomsg]{}

func (s stub[I, O]) Times(n int) Stub {
	s.times = n
	return s
}

func (s stub[I, O]) Name(name string) Stub {
	s.name = name
	return s
}

func (s stub[I, O]) typ(req reflect.Type, res reflect.Type) bool {
	return reflect.TypeOf(s.request) == req && reflect.TypeOf(s.response) == res
}

func (s stub[I, O]) resp() (any, error) {
	return s.response, s.err
}

func (s stub[I, O]) schema() string {
	return fmt.Sprintf("%s -> %s", reflect.TypeOf(s.request).Name(), reflect.TypeOf(s.response).Name())
}

func (s stub[I, O]) use() Stub {
	s.used++
	return s
}

func (s stub[I, O]) match(method string, request any) int {
	if method != s.method {
		return math.MaxInt
	}

	if s.used == s.times {
		return math.MaxInt
	}

	if request, ok := request.(I); ok {
		return s.matchfunc(request)
	}

	return math.MaxInt
}

func (s stub[I, O]) check(t *testing.T) {
	t.Helper()

	var name string
	if s.name == "" {
		name = fmt.Sprintf("stub for %q", s.method)
	} else {
		name = fmt.Sprintf("%q (stub for %q)", s.name, s.method)
	}

	switch {
	case s.used > s.times:
		t.Fatalf("%s was used %d extra times", name, s.used-s.times)
	case s.used < s.times:
		t.Fatalf("%s should be used %d more times", name, s.times-s.used)
	}
}
