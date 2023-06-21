package request

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const (
	ContextKeyRequestID contextKey = "requestID"
)

func Populate() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			requestID := extractRequestID(ctx)

			// Set the request ID in the context
			ctx = context.WithValue(ctx, ContextKeyRequestID, requestID)

			// Recover from any panics and log the request ID
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic occurred for request ID: %s\n", requestID)
					panic(r) // re-throw the panic
				}
			}()

			// Pass the modified context to the next middleware or handler
			return handler(ctx, req)
		}
	}
}

func extractRequestID(ctx context.Context) string {
	// Extract the request ID from the gRPC metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if requestIDs := md.Get("requestID"); len(requestIDs) > 0 {
			return requestIDs[0]
		}
	}

	// Extract the request ID from the HTTP headers

	// Extract the request ID from the HTTP headers
	if tr, ok := transport.FromServerContext(ctx); ok && tr.Kind() == transport.KindHTTP {
		if r, ok := http.RequestFromServerContext(ctx); ok {
			if requestIDs := r.Header.Get("requestID"); requestIDs != "" {
				return requestIDs
			}
		}
	}

	// Generate a new request ID
	return uuid.New().String()
}
