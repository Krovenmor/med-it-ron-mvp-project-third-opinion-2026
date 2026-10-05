//go:build e2e

package e2e

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/tests/e2e/queries"
)

const (
	defaultEnvFile = "../../config/default.env"

	pgUser     = "meditron"
	pgPassword = "meditron"
	pgDatabase = "meditron"
	pgAlias    = "postgres"
)

var (
	api   *apiClient
	db    *pgxpool.Pool
	query queries.Queries
	ai    *fakeAI
	mis   *fakeMIS
	env   *stack
)

type stack struct {
	network *testcontainers.DockerNetwork
	aiURL   string
	misURL  string
	ports   []int
}

type service struct {
	container *testcontainers.DockerContainer
	api       *apiClient
}

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Fprintln(os.Stderr, "e2e setup:", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	ctx := context.Background()

	query, err = queries.Load()
	if err != nil {
		return 0, err
	}

	nw, err := network.New(ctx)
	if err != nil {
		return 0, fmt.Errorf("create network: %w", err)
	}
	defer func() { err = errors.Join(err, nw.Remove(ctx)) }()

	pg, err := tcpostgres.Run(ctx, "postgres:17-alpine",
		tcpostgres.WithDatabase(pgDatabase),
		tcpostgres.WithUsername(pgUser),
		tcpostgres.WithPassword(pgPassword),
		tcpostgres.BasicWaitStrategies(),
		network.WithNetwork([]string{pgAlias}, nw),
	)
	defer func() { err = errors.Join(err, testcontainers.TerminateContainer(pg)) }()
	if err != nil {
		return 0, fmt.Errorf("start postgres: %w", err)
	}

	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return 0, fmt.Errorf("postgres dsn: %w", err)
	}
	db, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return 0, fmt.Errorf("connect postgres: %w", err)
	}
	defer db.Close()

	ai = newFakeAI()
	aiServer := httptest.NewServer(ai)
	defer aiServer.Close()

	mis = newFakeMIS()
	misServer := httptest.NewServer(mis)
	defer misServer.Close()

	env, err = newStack(nw, aiServer.URL, misServer.URL)
	if err != nil {
		return 0, err
	}

	svc, err := env.startService(ctx, map[string]string{
		"POSTGRES_DSN":    env.dsn(pgDatabase),
		"AI_SERVICE_MOCK": "false",
		"AI_SERVICE_URL":  env.aiURL,
		"DEMO_MODE":       "true",
	})
	if svc != nil {
		defer func() { err = errors.Join(err, testcontainers.TerminateContainer(svc.container)) }()
	}
	if err != nil {
		return 0, fmt.Errorf("start b2b-service: %w", err)
	}
	api = svc.api

	code = m.Run()
	if code != 0 {
		dumpLogs(ctx, svc.container)
	}
	return code, nil
}

func newStack(nw *testcontainers.DockerNetwork, aiServerURL, misServerURL string) (*stack, error) {
	aiPort, err := portOf(aiServerURL)
	if err != nil {
		return nil, err
	}
	misPort, err := portOf(misServerURL)
	if err != nil {
		return nil, err
	}
	return &stack{
		network: nw,
		aiURL:   fmt.Sprintf("http://%s:%d", testcontainers.HostInternal, aiPort),
		misURL:  fmt.Sprintf("http://%s:%d", testcontainers.HostInternal, misPort),
		ports:   []int{aiPort, misPort},
	}, nil
}

func (s *stack) dsn(database string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", pgUser, pgPassword, pgAlias, database)
}

func (s *stack) startService(ctx context.Context, overrides map[string]string) (*service, error) {
	vars, err := readEnvFile(defaultEnvFile)
	if err != nil {
		return nil, err
	}
	maps.Copy(vars, map[string]string{
		"MIS_URL":              s.misURL,
		"AI_SERVICE_TIMEOUT":   "5s",
		"WORKER_POLL_INTERVAL": "100ms",
		"WORKER_BACKOFF_BASE":  "100ms",
		"WORKER_BACKOFF_MAX":   "1s",
		"WORKER_MAX_ATTEMPTS":  strconv.Itoa(maxAttempts),
		"WORKER_JOB_TIMEOUT":   "10s",
		"WORKER_LEASE":         "30s",
	})
	maps.Copy(vars, overrides)

	ctr, err := testcontainers.Run(ctx, "",
		testcontainers.WithDockerfile(testcontainers.FromDockerfile{Context: "../.."}),
		testcontainers.WithEnv(vars),
		testcontainers.WithExposedPorts("8080/tcp"),
		testcontainers.WithHostPortAccess(s.ports...),
		network.WithNetwork([]string{"b2b-service"}, s.network),
		testcontainers.WithWaitStrategy(wait.ForLog("http server started").WithStartupTimeout(3*time.Minute)),
	)
	if err != nil {
		if ctr != nil {
			dumpLogs(ctx, ctr)
		}
		return &service{container: ctr}, err
	}

	endpoint, err := ctr.PortEndpoint(ctx, "8080/tcp", "http")
	if err != nil {
		return &service{container: ctr}, fmt.Errorf("service endpoint: %w", err)
	}
	return &service{container: ctr, api: newAPIClient(endpoint)}, nil
}

func readEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}
	vars := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("%s: malformed line %q", path, line)
		}
		vars[key] = value
	}
	return vars, nil
}

func portOf(rawURL string) (int, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return 0, fmt.Errorf("parse %q: %w", rawURL, err)
	}
	return strconv.Atoi(u.Port())
}

func dumpLogs(ctx context.Context, ctr *testcontainers.DockerContainer) {
	logs, err := ctr.Logs(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read service logs:", err)
		return
	}
	defer logs.Close()
	fmt.Fprintln(os.Stderr, "----- b2b-service logs -----")
	_, _ = io.Copy(os.Stderr, logs)
}
