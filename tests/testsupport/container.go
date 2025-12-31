package testsupport

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
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
		_ = ctr.Terminate(context.Background())
		return nil, fmt.Errorf("get container host: %w", err)
	}
	port, err := ctr.MappedPort(ctx, apiPort)
	if err != nil {
		_ = ctr.Terminate(context.Background())
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

// TestServerManager manages test server lifecycle with sync.Once initialization
type TestServerManager struct {
	once   sync.Once
	server *RunningAPIServer
	err    error
	url    string
}

// GetBaseURL returns the test server URL, starting the container if needed
func (m *TestServerManager) GetBaseURL(t *testing.T, opts APIServerOptions, logPrefix string) string {
	t.Helper()

	// Check if base URL is provided via environment variable
	baseURL := strings.TrimRight(os.Getenv("EC_TMPL_TEST_BASE_URL"), "/")
	if baseURL != "" {
		return baseURL
	}

	// Check for Docker image environment variable
	image := os.Getenv("EC_TMPL_E2E_IMAGE")
	if image == "" {
		t.Skip("EC_TMPL_E2E_IMAGE is not set")
	}

	// Start server once using sync.Once
	m.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		opts.Image = image
		m.server, m.err = StartAPIServer(ctx, opts)
		if m.err == nil {
			m.url = m.server.BaseURL
			fmt.Printf("\n%s tests running against: %s\n", logPrefix, m.url)
		}
	})

	if m.err != nil {
		t.Fatalf("failed to start %s container: %v", logPrefix, m.err)
	}
	return m.url
}

// Terminate stops the test server
func (m *TestServerManager) Terminate(ctx context.Context) error {
	if m.server != nil {
		return m.server.Terminate(ctx)
	}
	return nil
}
