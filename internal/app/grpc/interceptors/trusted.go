package interceptors

import (
	"context"
	"net/netip"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Trusted is an incterceptor checking if request was made from trusted subnet.
type Trusted struct {
	Config *config.Config
}

// TrustInterceptor verifies request if it came from the trusted network
// if request tries to access internal enpoints
func (t *Trusted) TrustInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod != "/proto.Shortener/InternalStats" {
		return handler(ctx, req)
	}
	if t.Config.TrustedSubnet == "" {
		return nil, status.Error(codes.PermissionDenied, "untrusted user")
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ipmd := md.Get("X-Real-IP")
		if len(ipmd) > 0 {
			xRealIP := ipmd[0]
			netip, err := netip.ParseAddr(xRealIP)
			if err != nil {
				return nil, status.Error(codes.PermissionDenied, err.Error())
			}
			if !t.Config.SubnetPrefix.Contains(netip) {
				return nil, status.Error(codes.PermissionDenied, "untrusted user")
			}
		}
	}
	return handler(ctx, req)
}
