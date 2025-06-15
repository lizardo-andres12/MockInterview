package interceptors

import (
	"context"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// userUUIDKey is the context key for the authenticated user's UUID
// We use an unexported type to avoid collisions.
type userUUIDKey struct{}

// publicMethods maps gRPC method names that should skip auth
var publicMethods = map[string]bool{
	"/proto.AuthService/Register": true,
	"/proto.AuthService/Login":    true,
}

// NewJWTInterceptor returns unary and stream interceptors for JWT auth
func NewJWTInterceptor(secret []byte) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor) {
	authFn := func(ctx context.Context) (context.Context, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		vals := md.Get("Authorization")
		if len(vals) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token required")
		}
		parts := strings.SplitN(vals[0], " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return nil, status.Error(codes.Unauthenticated, "bad authorization header format")
		}
		tokenStr := parts[1]

		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, status.Error(codes.Unauthenticated, "unexpected signing method")
			}
			return secret, nil
		})
		if err != nil || !tok.Valid {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}
		subRaw, ok := claims["sub"]
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "subject claim missing")
		}
		sub, ok := subRaw.(string)
		if !ok || sub == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid subject claim")
		}

		// inject the UUID into context
		return context.WithValue(ctx, userUUIDKey{}, sub), nil
	}

	// Unary interceptor
	unary := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}
		authCtx, err := authFn(ctx)
		if err != nil {
			return nil, err
		}
		return handler(authCtx, req)
	}

	// Stream interceptor
	stream := func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if publicMethods[info.FullMethod] {
			return handler(srv, ss)
		}
		authCtx, err := authFn(ss.Context())
		if err != nil {
			return err
		}

		wrapped := grpcmiddleware.WrapServerStream(ss)
		wrapped.WrappedContext = authCtx
		return handler(srv, wrapped)
	}

	return unary, stream
}

// UserUUIDFromContext retrieves the authenticated user's UUID from context
func UserUUIDFromContext(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(userUUIDKey{}).(string)
	return u, ok
}

