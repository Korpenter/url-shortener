package interceptors

import (
	"context"
	"net/netip"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/utils/helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	usermd, ok := helpers.CheckMDValue(ctx, "X-Real-IP")
	if ok {
		netip, err := netip.ParseAddr(usermd)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if !t.Config.SubnetPrefix.Contains(netip) {
			return nil, status.Error(codes.PermissionDenied, "untrusted user")
		}
	}
	return handler(ctx, req)
}
