package sandbox

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

const (
	_baseDockerInstances = 10
)

// Pool manages a set of pre-started containers ready for exec.
type Pool struct {
	cli        *client.Client
	image      string
	containers chan string        // channel of container IDs
}

// NewPool pre-spins _baseDockerInstances many `img` containers
func NewPool(cli *client.Client, img string) (*Pool, error) {
	p := &Pool{
		cli: cli,
		image: img,
		containers: make(chan string, _baseDockerInstances),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	// Pull image
	if _, err := cli.ImagePull(ctx, img, image.PullOptions{}); err != nil {
		return nil, fmt.Errorf("pull image: %w", err)
	}

	// Start containers
	for range _baseDockerInstances {
		resp, err := cli.ContainerCreate(
			ctx,
			&container.Config{
				Image: img,
				Cmd: []string{"sleep", "infinity"},
			},
			nil,
			nil,
			nil,
			"",
		)
		if err != nil {
			return nil, fmt.Errorf("create container: %w", err)
		}
		if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			return nil, fmt.Errorf("start container: %w", err)
		}
		p.containers <- resp.ID
	}
	return p, nil
}

// Exec runs `code` inside any free container, returning stdout/stderr.
func (p *Pool) Exec(code string, workdir string) (string, string, error) {
	ctx := context.Background()

	// pull container from container chan and re-insert when Exec returns
	cid := <-p.containers
	defer func() {
		p.containers <- cid
	}()

	// write code into main.go and run
	execCmd := fmt.Sprintf(
		"sh -c 'mkdir -p %s && echo %s > %s/main.go && go run %s/main.go'",
		workdir,
		shellQuote(code),
		workdir,
		workdir,
	)

	execConfig := container.ExecOptions{
		Cmd:          []string{"sh", "-c", execCmd},
		AttachStdout: true,
		AttachStderr: true,
	}
	execResp, err := p.cli.ContainerExecCreate(ctx, cid, execConfig)
	if err != nil {
		return "", "", err
	}
	attachResp, err := p.cli.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{})
	if err != nil {
		return "", "", err
	}
	defer attachResp.Close()

	out, err := io.ReadAll(attachResp.Reader)
	if err != nil {
		return "", "", err
	}
	// for simplicity, stderr is empty
	return string(out), "", nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}


