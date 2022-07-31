package mockaka

import "fmt"

type unexpectedCall struct {
	service string
	method  string
	request any
}

func (c unexpectedCall) String() string {
	return fmt.Sprintf("%s.%s(%+v)", c.service, c.method, c.request)
}
