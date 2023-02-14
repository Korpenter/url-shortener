// Package interceptors provides custom interceptors for gprc server.
package interceptors

import (
	"context"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/Mldlr/url-shortener/internal/app/utils/helpers"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Auth is an authentication interceptor.
type Auth struct {
	Config *config.Config
}

// Authenticate authenticates user request by adding user ID and signature cookies to the request
// if they are not present or invalid.
func (a *Auth) AuthInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	usermd, uOk := helpers.CheckMDValue(ctx, "user_id")
	signmd, sOk := helpers.CheckMDValue(ctx, "signature")
	if uOk && sOk && signmd == encoders.HMACString(usermd, a.Config.SecretKey) {
		return handler(ctx, req)
	}
	userID := uuid.New().String()
	signature := encoders.HMACString(userID, a.Config.SecretKey)
	md := metadata.New(map[string]string{
		"user_id":   userID,
		"signature": signature},
	)
	ctxNew := metadata.NewIncomingContext(ctx, md)
	return handler(ctxNew, req)
}
