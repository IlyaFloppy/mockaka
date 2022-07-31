package stringscacheapp

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/IlyaFloppy/mockaka/example/stringscache/api/scpb"
	"github.com/IlyaFloppy/mockaka/example/stringscache/api/spb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Service is a StringsCache app.
type Service struct {
	address        string
	stringsAddress string
	stringsClient  spb.StringsServiceClient

	cache map[string]string
	mx    sync.Mutex

	scpb.UnimplementedStringsCacheServiceServer

	readyCh chan struct{}
}

// NewService creates new StringsCache app.
func NewService(address string, stringsAddress string) *Service {
	return &Service{
		address:        address,
		stringsAddress: stringsAddress,
		cache:          make(map[string]string),
		readyCh:        make(chan struct{}),
	}
}

// Run runs server. It does not return until ctx is done or error is encountered.
func (s *Service) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen on address %q: %s", s.address, err.Error())
	}

	conn, err := grpc.DialContext(ctx, s.stringsAddress, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to dial strings service")
	}

	stringsClient := spb.NewStringsServiceClient(conn)
	s.stringsClient = stringsClient

	grpcServer := grpc.NewServer(
		grpc.WriteBufferSize(0),
		grpc.ReadBufferSize(0),
	)

	scpb.RegisterStringsCacheServiceServer(grpcServer, s)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	close(s.readyCh)

	return errors.Wrap(grpcServer.Serve(lis), "serve failed")
}

// Ready returns a channel that is closed when service is ready.
func (s *Service) Ready() <-chan struct{} {
	return s.readyCh
}
