package config

import (
  	"github.com/docker/docker/client"
)

// NewDockerClient constructs a Docker SDK client
func NewDockerClient() (*client.Client, error) {
  return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

