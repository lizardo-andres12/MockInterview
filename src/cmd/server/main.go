package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go.mocker.com/src/config"
	"go.mocker.com/src/controller"
	"go.mocker.com/src/handler"
	"go.mocker.com/src/interceptors"
	"go.mocker.com/src/proto"
	"go.mocker.com/src/repository"
)

// provideJWTSecret reads the JWT secret from env
func provideJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

// newGRPCServer creates and run gRPC server at port set by environment variable
func newGRPCServer(
	lc fx.Lifecycle,
	authHandler *handler.AuthHandler,
	jwtSecret []byte,
	logger *zap.Logger,
) *grpc.Server {
	unary, stream := interceptors.NewJWTInterceptor(jwtSecret)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(unary),
		grpc.StreamInterceptor(stream),
	)
	proto.RegisterAuthServiceServer(server, authHandler)
	reflection.Register(server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				logger.Info("gRPC server starting", zap.String("port", os.Getenv("GRPC_PORT")))
				if err := server.Serve(lis); err != nil {
					logger.Fatal("gRPC server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("gRPC server stopping")
			server.GracefulStop()
			return nil
		},
	})

	return server
}

// buildFxApp assembles the Fx application
func buildFxApp() *fx.App {
	return fx.New(
		fx.Provide(
			zap.NewExample, // use production in real env: zap.NewProduction
			config.InitDB,
			repository.NewUserRepo,
			controller.NewAuthController,
			handler.NewAuthHandler,
			provideJWTSecret,
		),
		fx.Invoke(newGRPCServer),
	)
}

func main() {
	if os.Getenv("MYSQL_DSN") == "" || os.Getenv("JWT_SECRET") == "" || os.Getenv("GRPC_PORT") == "" {
		fmt.Println("Environment variables MYSQL_DSN, JWT_SECRET, and GRPC_PORT must be set")
		os.Exit(1)
	}

	app := buildFxApp()
	app.Run()
}

