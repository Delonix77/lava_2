package testproto

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InMemoryClientConn creates a grpc client connection to a grpc server which is listening on a random port.
func InMemoryClientConn(t *testing.T, grpcSrv *grpc.Server) *grpc.ClientConn {
	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go func() {
		defer lis.Close()
		if err := grpcSrv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// returns the listening port
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return conn
}
