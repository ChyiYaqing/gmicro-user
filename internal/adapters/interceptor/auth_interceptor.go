package interceptor

import (
	"context"
	"log"

	"github.com/chyiyaqing/gmicro-user/internal/adapters/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	AuthInterceptor struct {
		jwtManager      *auth.JWTManager
		accessibleRoles map[string][]string // 映射每个RPC方法定一个可以访问角色列表
	}
)

func NewAuthInterceptor(jwtManager *auth.JWTManager, accessibleRoles map[string][]string) (*AuthInterceptor, error) {
	return &AuthInterceptor{jwtManager: jwtManager, accessibleRoles: accessibleRoles}, nil
}

func (i *AuthInterceptor) UnaryAuthMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Println("--> unary interceptor: ", info.FullMethod)

	err := i.authorize(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC.
func (i *AuthInterceptor) StreamAuthMiddleware(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("--> stream interceptor: ", info.FullMethod)

	err := i.authorize(stream.Context(), info.FullMethod)
	if err != nil {
		return err
	}

	return handler(srv, stream)
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) error {
	accessibleRoles, ok := i.accessibleRoles[method]
	if !ok {
		// everyone can access
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := i.jwtManager.ValidateToken(ctx, accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		if role == claims.Role {
			return nil
		}
	}
	return status.Errorf(codes.PermissionDenied, "no permission to access this RPC")
}
