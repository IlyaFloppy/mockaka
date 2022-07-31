package example_test

import (
	"context"
	"testing"

	"github.com/IlyaFloppy/mockaka"
	"github.com/IlyaFloppy/mockaka/example/stringscache/api/scpb"
	"github.com/IlyaFloppy/mockaka/example/stringscache/api/spb"
	"github.com/IlyaFloppy/mockaka/example/stringscache/stringscacheapp"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//nolint:funlen
func TestCacheService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const stringsServiceAddress = "localhost:1111"

	// mock StringsService.
	m := mockaka.
		New(t).
		AddService(mockaka.Register(spb.RegisterStringsServiceServer), []mockaka.Stub{
			mockaka.NewStub(
				mockaka.UnaryMethod(spb.StringsServiceServer.Reverse),
				&spb.Message{
					Message: "12345",
				},
				&spb.Message{
					Message: "54321",
				},
			).Times(2).Name("12345"),
			mockaka.NewErrStub[*spb.Message, *spb.Message](
				mockaka.UnaryMethod(spb.StringsServiceServer.Reverse),
				&spb.Message{
					Message: "too long",
				},
				status.Errorf(codes.InvalidArgument, "request string is too long"),
			).Times(100).Name("too long"),
		})
	go m.Run(stringsServiceAddress)
	defer m.Stop()

	// run StringsCacheService.
	client, stop := runStringsCacheApp(t, stringsServiceAddress)
	defer stop()

	// make rpc calls to StringsCacheService.
	{
		// first call to StringsService.Reverse.
		// test that Reverse result is cached.
		for i := 0; i < 100; i++ {
			res, err := client.Reverse(ctx, &scpb.Message{
				Message: "12345",
			})
			require.NoError(t, err)
			mockaka.PBEq(t, &scpb.Message{
				Message: "54321",
			}, res)
		}

		// invalidate cache.
		_, err := client.Invalidate(ctx, &scpb.InvalidateRequest{})
		require.NoError(t, err)

		// second call to StringsService.Reverse.
		// test that cache is cleared and Reverse result is cached after next rpc call.
		for i := 0; i < 100; i++ {
			res, err := client.Reverse(ctx, &scpb.Message{
				Message: "12345",
			})
			require.NoError(t, err)
			mockaka.PBEq(t, &scpb.Message{
				Message: "54321",
			}, res)
		}

		// third call to StringsService.Reverse.
		// test that StringsCacheService does not cache errors.
		for i := 0; i < 100; i++ {
			res, err := client.Reverse(ctx, &scpb.Message{
				Message: "too long",
			})
			st, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, codes.InvalidArgument, st.Code())
			require.Nil(t, res)
		}
	}
}

func runStringsCacheApp(t *testing.T, stringsServiceAddress string,
) (client scpb.StringsCacheServiceClient, stop func()) {
	ctx := context.Background()
	const address = ":4444"

	stop = stringscacheapp.Main(ctx, stringscacheapp.Config{
		Address:               address,
		StringsServiceAddress: stringsServiceAddress,
	}) // Main returns when app is ready.

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	require.NoError(t, err)

	return scpb.NewStringsCacheServiceClient(conn), stop
}
