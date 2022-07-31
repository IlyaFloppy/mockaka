package mockaka

import (
	"context"
	"net"
	"reflect"
	"sync"
	"testing"

	"google.golang.org/grpc"
)

// Mock is a mock server. Use New to instantiate.
type Mock struct {
	t            *testing.T
	mx           sync.Mutex
	serviceDescs []*grpc.ServiceDesc

	serveErr error

	stop func()
}

// New creates new mock server.
// Use AddService for setup and Run to start listening.
func New(t *testing.T) *Mock {
	t.Helper()

	m := Mock{
		t: t,
	}

	return &m
}

// Reg is a processed representation of RegisterServiceServer function generated by protoc.
type Reg struct {
	desc *grpc.ServiceDesc
	srv  reflect.Type
}

// Register creates Reg from RegisterServiceServer function.
func Register[ServiceServerType any](registerServerFunc registerServerFunc[ServiceServerType]) Reg {
	var ret Reg
	var rc registerCallback = func(desc *grpc.ServiceDesc, srv any) {
		ret.desc = desc
		ret.srv = reflect.TypeOf(registerServerFunc).In(1)
	}

	var nilServer ServiceServerType
	registerServerFunc(rc, nilServer)

	return ret
}

// AddService adds new service with stubs. First argument is mockaka.Register(pb.ServiceServer).
//nolint:funlen,cyclop
func (m *Mock) AddService(reg Reg, stubs []Stub) *Mock {
	m.t.Helper()

	m.mx.Lock()
	defer m.mx.Unlock()

	// Stub slices in this map will be changed by handlers to track stubs usage.
	// Original stubs are not mutated.
	stubsByMethod := make(map[string][]Stub)
	unexpectedCalls := make([]unexpectedCall, 0)
	var unexpectedCallsMx sync.Mutex

	validStubs := make([]bool, len(stubs))

	for i, methodDesc := range reg.desc.Methods {
		method, _ := reg.srv.MethodByName(methodDesc.MethodName)
		if err := matchUnaryHandlerSignature(method); err != nil {
			m.t.Logf("method %q is not an unary handler(%s) and cannot be mocked", method.Name, err.Error())
			continue
		}

		// Method's arguments are (context.Context, *ProtobufMessage).
		// Only the last one is of interest.
		requestType := method.Type.In(method.Type.NumIn() - 1)
		resultType := method.Type.Out(0)

		methodStubs := make([]Stub, 0, len(stubs))
		for j, stub := range stubs {
			if stub.is(method.Name, requestType, resultType) {
				methodStubs = append(methodStubs, stub)
				validStubs[j] = true
			}
		}

		stubsByMethod[method.Name] = methodStubs

		errCall := func(request any) {
			unexpectedCallsMx.Lock()
			unexpectedCalls = append(unexpectedCalls, unexpectedCall{
				service: reg.desc.ServiceName,
				method:  method.Name,
				request: request,
			})
			unexpectedCallsMx.Unlock()
		}

		reg.desc.Methods[i].Handler = m.newUnaryHandler(method.Name, requestType, methodStubs, errCall)
	}

	invalidStubs := make([]Stub, 0, len(stubs))
	for i, valid := range validStubs {
		if !valid {
			invalidStubs = append(invalidStubs, stubs[i])
		}
	}

	if len(invalidStubs) > 0 {
		for _, stub := range invalidStubs {
			m.t.Logf("stub [%s] does not match contract", stub.schema())
		}
		m.t.FailNow()
	}

	m.t.Cleanup(func() {
		m.t.Helper()

		for _, ur := range unexpectedCalls {
			m.t.Fatalf("unexpected call: %s", ur.String())
		}

		for _, stubs := range stubsByMethod {
			for _, stub := range stubs {
				stub.check(m.t)
			}
		}
	})

	m.serviceDescs = append(m.serviceDescs, reg.desc)

	return m
}

// Run mock server using provided stubs. Run blocks until mock server is ready.
func (m *Mock) Run(address string) {
	m.t.Helper()

	m.mx.Lock()
	defer m.mx.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	m.stop = cancel

	lis, err := net.Listen("tcp", address)
	if err != nil {
		m.t.Fatalf("failed to listen on address %q: %s", address, err.Error())
	}

	grpcServer := grpc.NewServer(
		grpc.WriteBufferSize(0),
		grpc.ReadBufferSize(0),
	)

	for _, sd := range m.serviceDescs {
		grpcServer.RegisterService(sd, nil)
	}

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- grpcServer.Serve(lis)
	}()

	m.t.Cleanup(func() {
		if err := <-errCh; err != nil {
			m.t.Fatalf("grpc server failed to serve: %e", err)
		}
	})
}

// Stop mock service.
func (m *Mock) Stop() {
	m.stop()

	// wait for Run to release mutex.
	m.mx.Lock()
	m.mx.Unlock()
}

type registerServerFunc[ServiceServer any] func(s grpc.ServiceRegistrar, srv ServiceServer)
type registerCallback func(desc *grpc.ServiceDesc, srv any)

func (c registerCallback) RegisterService(desc *grpc.ServiceDesc, srv any) {
	c(desc, srv)
}
