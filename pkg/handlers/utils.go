package handlers

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func ExtractPasskey(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	passwordRaw := md.Get("password")
	if len(passwordRaw) == 0 {
		return ""
	}

	return passwordRaw[0]
}
