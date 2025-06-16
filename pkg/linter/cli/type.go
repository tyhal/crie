package cli

import (
	"context"
	"github.com/docker/docker/client"
)

type execMode int

const (
	auto execMode = iota
	podman
	docker
	host
)

// Lint defines a predefined command to run against a file
type Lint struct {
	Bin       string
	FrontPar  Par
	EndPar    Par
	Container Container
	execMode  execMode
	cleanedUp chan error
}

// Par represents cli parameters
type Par []string

// Container an image to launch
type Container struct {
	Image        string
	clientDocker *client.Client
	clientPodman *context.Context
	id           string
}
