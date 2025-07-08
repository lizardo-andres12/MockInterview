package handler

import (
  "context"

  pb "go.mocker.com/src/proto"
  "go.mocker.com/src/controller"
)

type ExecHandler struct {
  ctrl controller.ExecController
  pb.UnimplementedExecServiceServer
}

func NewExecHandler(ctrl controller.ExecController) *ExecHandler {
  return &ExecHandler{ctrl: ctrl}
}

func (h *ExecHandler) Execute(ctx context.Context, req *pb.ExecRequest) (*pb.ExecResponse, error) {
  out, errout, code, err := h.ctrl.Execute(req.Code)
  if err != nil {
    return nil, err
  }
  return &pb.ExecResponse{Stdout: out, Stderr: errout, ExitCode: int32(code)}, nil
}

