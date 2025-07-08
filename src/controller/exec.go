package controller

import (
	"go.mocker.com/src/sandbox"
)

// ExecController defines execution logic
type ExecController interface {
  Execute(code string) (stdout, stderr string, exitCode int, err error)
}

type execCtrl struct {
  pool *sandbox.Pool
}

// NewExecController constructs a new ExecController
func NewExecController(pool *sandbox.Pool) ExecController {
  return &execCtrl{pool: pool}
}

// Execute runs the code in sandbox
func (c *execCtrl) Execute(code string) (string, string, int, error) {
  out, errout, err := c.pool.Exec(code, "/workspace")
  if err != nil {
    return "", "", 1, err
  }
  return out, errout, 0, nil
}

