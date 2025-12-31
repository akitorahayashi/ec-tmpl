package testsupport

import (
	"context"
	"fmt"
	"time"

	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	apiPort    = "8000/tcp"
	healthPath = "/health"
)

type RunningAPIServer struct {
	container testcontainers.Container
	BaseURL   string
}

type APIServerOptions struct {
	Image          string
	Env            map[string]string
	StartupTimeout time.Duration
}

func StartAPIServer(ctx context.Context, opts APIServerOptions) (*RunningAPIServer, error) {
	if opts.Image == "" {
		return nil, fmt.Errorf("image is required")
	}
	if opts.StartupTimeout == 0 {
		opts.StartupTimeout = 30 * time.Second
	}

	req := testcontainers.ContainerRequest{
		Image:        opts.Image,
		Env:          opts.Env,
		ExposedPorts: []string{apiPort},
		WaitingFor:   wait.ForHTTP(healthPath).WithStartupTimeout(opts.StartupTimeout),
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start container: %w", err)
	}

	host, err := ctr.Host(ctx)
	if err != nil {
		_ = ctr.Terminate(ctx)
		return nil, fmt.Errorf("get container host: %w", err)
	}
	port, err := ctr.MappedPort(ctx, apiPort)
	if err != nil {
		_ = ctr.Terminate(ctx)
		return nil, fmt.Errorf("get container port: %w", err)
	}

	return &RunningAPIServer{
		container: ctr,
		BaseURL:   fmt.Sprintf("http://%s:%s", host, port.Port()),
	}, nil
}

func (s *RunningAPIServer) Terminate(ctx context.Context) error {
	if s == nil || s.container == nil {
		return nil
	}
	return s.container.Terminate(ctx)
}
