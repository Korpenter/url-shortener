package helpers

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// CheckMDValue checks if theres is a value in metadata in context
func CheckMDValue(ctx context.Context, value string) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get(value)
		if len(arr) > 0 {
			return arr[0], true
		}
	}
	return "", false
}
