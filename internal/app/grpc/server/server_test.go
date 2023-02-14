package server

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/Mldlr/url-shortener/internal/app/config"
	pb "github.com/Mldlr/url-shortener/internal/app/grpc/proto"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var server *GRPCServer

func init() {
	repo := storage.NewMockRepo()
	testSubnet := "155.155.5.0/24"
	testPrefix, _ := netip.ParsePrefix(testSubnet)
	cfg := &config.Config{
		GRPCAddress:   ":8888",
		TrustedSubnet: testSubnet,
		SubnetPrefix:  testPrefix,
		SecretKey:     []byte("defaultKeyUrlSHoRtenEr"),
	}
	shortener := service.NewShortenerImpl(repo, cfg)
	server = NewGRPCServer(shortener, cfg)
	go func() {
		server.Run(context.Background())
	}()
}

// TestGRPCServer_Stats checks if server and trusted subnet interceptor are working correctly
func TestGRPCServer_Stats(t *testing.T) {
	tests := []struct {
		name    string
		xRealIP string
		errCode codes.Code
	}{
		{
			name:    "Trusted",
			xRealIP: "155.155.5.6",
			errCode: codes.OK,
		},
		{
			name:    "Untrusted",
			xRealIP: "11.11.5.6",
			errCode: codes.PermissionDenied,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := grpc.Dial(":8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
			assert.NoError(t, err)
			defer conn.Close()
			c := pb.NewShortenerClient(conn)
			ctx := context.Background()
			request := &pb.StatsRequest{}
			assert.NoError(t, err)
			outCtx := metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"X-Real-IP": tt.xRealIP}))
			_, err = c.InternalStats(outCtx, request)
			if statusErr, ok := status.FromError(err); ok {
				assert.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

// TestGRPCServer_DeleteURL checks if server and auth interceptor are working correctly
func TestGRPCServer_DeleteURL(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		signature    string
		request      []string
		errCode      codes.Code
		checkErrCode codes.Code
	}{
		{
			name:         "Delete with wrong signature",
			userID:       "user1",
			signature:    "asdasdasdasdas",
			request:      []string{"aQqomlSbUsE"},
			errCode:      codes.OK,
			checkErrCode: codes.OK,
		},
		{
			name:         "Delete with right signature",
			userID:       "KS097f1lS&F",
			signature:    "e69476e0425466968352c3f5f450516df16461eb4511ce750a7be3e7478cf72a",
			request:      []string{"aQqomlSbUsE"},
			errCode:      codes.OK,
			checkErrCode: codes.Unavailable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := grpc.Dial(":8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
			assert.NoError(t, err)
			defer conn.Close()
			c := pb.NewShortenerClient(conn)
			ctx := context.Background()
			request := &pb.DeleteURLRequest{Urls: tt.request}
			assert.NoError(t, err)
			outCtx := metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"user_id": tt.userID, "signature": tt.signature}))
			_, err = c.DeleteBatch(outCtx, request)
			if statusErr, ok := status.FromError(err); ok {
				assert.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
			checkReq := &pb.ExpandURLRequest{
				ShortURL: tt.request[0],
			}
			var statusCheckErr *status.Status

			time.Sleep(time.Second * 10)
			_, err = c.Expand(outCtx, checkReq)
			statusCheckErr, _ = status.FromError(err)
			assert.Equal(t, tt.checkErrCode.String(), statusCheckErr.Code().String())
		})
	}
}
