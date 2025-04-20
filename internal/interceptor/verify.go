package interceptor

import (
	"context"
	"fmt"
	"log"

	"github.com/laiker/auth/pkg/access_v1"
	"github.com/laiker/chat-server/internal/config/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	_ "github.com/laiker/auth/pkg/access_v1"
)

func VerifyInterceptor() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
		}

		ctx = metadata.NewOutgoingContext(ctx, md)

		crds, err := credentials.NewClientTLSFromFile("service.pem", "localhost")

		if err != nil {
			log.Printf("failed to load TLS credentials: %v", err)
		}

		config, err := env.NewAuthConfig()

		if err != nil {
			log.Printf("failed to load config: %v", err)
		}

		conn, err := grpc.NewClient(
			fmt.Sprintf("%s:%s", config.Host(), "50052"),
			grpc.WithTransportCredentials(crds),
		)

		if err != nil {
			log.Printf("failed to dial GRPC client: %v", err)
		}

		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {

			}
		}(conn)

		authClient := access_v1.NewAccessV1Client(conn)

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "Authorization token is required")
		}

		hasAccessRequest := &access_v1.CheckRequest{
			EndpointAddress: info.FullMethod,
		}

		_, err = authClient.HasAccess(ctx, hasAccessRequest)

		if err != nil {
			return nil, status.Errorf(codes.Internal, "Error checking access: %v", err)
		}

		return handler(ctx, req)
	}

}
