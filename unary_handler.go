package mockaka

import (
	"context"
	"math"
	"reflect"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type unaryHandler = func(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error)

func (m *Mock) newUnaryHandler(
	method string,
	requestType reflect.Type,
	stubs []Stub,
	unexpectedCall func(request any),
) unaryHandler {
	var mx sync.Mutex // protects stubs.

	return func(srv any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
		// pointer.
		request := reflect.New(requestType.Elem()).Interface()

		if err := dec(request); err != nil {
			return nil, err
		}

		var (
			bestDiffLen  = math.MaxInt
			bestResponse any
			bestError    error
			usedStubIdx  = -1
		)

		for i, stub := range stubs {
			diff := stub.match(method, request)
			if diff < bestDiffLen {
				bestDiffLen = diff
				bestResponse, bestError = stub.resp()
				usedStubIdx = i
			}
		}

		mx.Lock()
		defer mx.Unlock()

		switch {
		case bestDiffLen == 0: // exact match.
			stubs[usedStubIdx] = stubs[usedStubIdx].use()
			return bestResponse, bestError //nolint:wrapcheck
		default: // no stub or partial match.
			unexpectedCall(request)
			return nil, status.Error(codes.Unimplemented, "no stub for request") //nolint:wrapcheck
		}
	}
}
