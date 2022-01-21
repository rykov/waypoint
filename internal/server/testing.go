package server

import (
	"context"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/protocolversion"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
)

// TestServer starts a server and returns a gRPC client to that server.
// We use t.Cleanup to ensure resources are automatically cleaned up.
func TestServer(t testing.T, impl pb.WaypointServer, opts ...TestOption) pb.WaypointClient {
	require := require.New(t)

	c := testConfig{
		ctx: context.Background(),
	}
	for _, opt := range opts {
		opt(&c)
	}

	// Listen on a random port
	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(err)
	t.Cleanup(func() { ln.Close() })

	// We make run a function since we'll call it to restart too
	run := func(ctx context.Context) context.CancelFunc {
		ctx, cancel := context.WithCancel(ctx)
		opts := []Option{
			WithContext(ctx),
			WithGRPC(ln),
			WithImpl(impl),
		}
		if ac, ok := impl.(AuthChecker); ok {
			opts = append(opts, WithAuthentication(ac))
		}

		go Run(opts...)
		t.Cleanup(func() { cancel() })

		return cancel
	}

	// Create the server
	cancel := run(c.ctx)

	// If we have a restart channel, then listen to that for restarts.
	if c.restartCh != nil {
		doneCh := make(chan struct{})
		t.Cleanup(func() { close(doneCh) })

		go func() {
			for {
				select {
				case <-c.restartCh:
					// Cancel the old context
					cancel()

					// This can fail, but it probably won't. Can't think of
					// a cleaner way since gRPC force closes its listener.
					ln, err = net.Listen("tcp", ln.Addr().String())
					if err != nil {
						return
					}

					// Create a new one
					cancel = run(context.Background())

				case <-doneCh:
					return
				}
			}
		}()
	}

	// Get our version info we'll set on the client
	vsnInfo := TestVersionInfoResponse().Info

	// connect is a function since we need to connect multiple times:
	// once to bootstrap, then again with our auth information.
	connect := func(token string) (*grpc.ClientConn, error) {
		opts := []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(protocolversion.UnaryClientInterceptor(vsnInfo)),
			grpc.WithStreamInterceptor(protocolversion.StreamClientInterceptor(vsnInfo)),
			grpc.WithPerRPCCredentials(serverclient.ContextToken(token)),
		}

		return grpc.DialContext(context.Background(), ln.Addr().String(), opts...)
	}

	// Connect, this should retry in the case Run is not going yet
	conn, err := connect("")
	require.NoError(err)
	client := pb.NewWaypointClient(conn)

	// Bootstrap
	tokenResp, err := client.BootstrapToken(context.Background(), &empty.Empty{})
	if status.Code(err) == codes.PermissionDenied {
		// Ignore bootstrap already complete errors
		err = nil
		tokenResp = &pb.NewTokenResponse{Token: ""}
	}
	conn.Close()
	require.NoError(err)

	// Reconnect with a token
	token := c.token
	if token == "" {
		token = tokenResp.Token
	}
	conn, err = connect(token)
	require.NoError(err)
	t.Cleanup(func() { conn.Close() })
	return pb.NewWaypointClient(conn)
}

// TestOption is used with TestServer to configure test behavior.
type TestOption func(*testConfig)

type testConfig struct {
	ctx       context.Context
	restartCh <-chan struct{}
	token     string
}

// TestWithContext specifies a context to use with the test server. When
// this is done then the server will exit.
func TestWithContext(ctx context.Context) TestOption {
	return func(c *testConfig) {
		c.ctx = ctx
	}
}

// TestWithRestart specifies a channel that will be sent to trigger
// a restart. The restart happens asynchronously. If you want to ensure the
// server is shutdown first, use TestWithContext, shut it down, wait for
// errors on the API, then restart.
func TestWithRestart(ch <-chan struct{}) TestOption {
	return func(c *testConfig) {
		c.restartCh = ch
	}
}

// TestWithToken specifies a specific token to use for auth.
func TestWithToken(token string) TestOption {
	return func(c *testConfig) {
		c.token = token
	}
}

// TestVersionInfoResponse generates a valid version info response for testing
func TestVersionInfoResponse() *pb.GetVersionInfoResponse {
	return &pb.GetVersionInfoResponse{
		Info: &pb.VersionInfo{
			Api: &pb.VersionInfo_ProtocolVersion{
				Current: 10,
				Minimum: 1,
			},

			Entrypoint: &pb.VersionInfo_ProtocolVersion{
				Current: 10,
				Minimum: 1,
			},
		},
	}
}
