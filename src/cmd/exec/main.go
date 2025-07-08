package main

import (
  "context"
  "fmt"
  "net"
  "os"

  "go.uber.org/fx"
  "go.uber.org/zap"
  "google.golang.org/grpc"
  "github.com/docker/docker/client"

  "go.mocker.com/src/config"
  "go.mocker.com/src/sandbox"
  "go.mocker.com/src/controller"
  "go.mocker.com/src/handler"
  pb "go.mocker.com/src/proto"
)

func provideDockerClient() (*client.Client, error) {
  return config.NewDockerClient()
}

func newExecGRPCServer(
  lc fx.Lifecycle,
  handler *handler.ExecHandler,
  logger *zap.Logger,
) {
  server := grpc.NewServer()
  pb.RegisterExecServiceServer(server, handler)

  lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("EXEC_PORT")))
  if err != nil {
    logger.Fatal("listen failed", zap.Error(err))
  }

  lc.Append(fx.Hook{
    OnStart: func(context.Context) error {
      go server.Serve(lis)
      logger.Info("ExecService listening", zap.String("port", os.Getenv("EXEC_PORT")))
      return nil
    },
    OnStop: func(ctx context.Context) error {
      server.GracefulStop()
      return nil
    },
  })
}

func main() {
  fx.New(
    fx.Provide(
      zap.NewProduction,
      provideDockerClient,
      func(cli *client.Client) (*sandbox.Pool, error) {
        return sandbox.NewPool(cli, "go-sandbox:latest")
      },
      controller.NewExecController,
      handler.NewExecHandler,
    ),
    fx.Invoke(newExecGRPCServer),
  ).Run()
}

