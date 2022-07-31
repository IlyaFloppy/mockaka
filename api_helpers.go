package mockaka

import (
	"context"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// UnaryMethod gets rpc call's name. Use like this:
// mockaka.UnaryMethod(pb.ServiceServer.Method).
func UnaryMethod[S any, I, O protomsg](m func(S, context.Context, I) (O, error)) string {
	full := runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
	parts := strings.Split(full, ".")

	return parts[len(parts)-1]
}

// PBEq is equivalent to require.Equal for protobuf structs.
// It only compares exported fields unlike require.Equal.
func PBEq[T protomsg](t *testing.T, expected, actual T, msgAndArgs ...any) {
	expected = PB(expected)
	actual = PB(actual)

	require.Equal(t, expected, actual, msgAndArgs...)
}

// PB prepares protobuf struct for comparison. It returns a copy of it's argument
// with all unexported fields set to zero values.
func PB[T protomsg](t T) T {
	var p T
	ret := reflect.New(reflect.TypeOf(p).Elem())

	typ := reflect.TypeOf(t).Elem()
	val := reflect.ValueOf(t)

	for i := 0; i < typ.NumField(); i++ {
		if f := ret.Elem().Field(i); f.CanSet() {
			f.Set(val.Elem().Field(i))
		}
	}

	return ret.Interface().(T)
}
