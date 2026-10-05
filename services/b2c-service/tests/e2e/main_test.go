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

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultEnvFile = "../../config/default.env"
	clinicPhone    = "+7 (495) 000-00-00"
	clinicName     = "Клиника e2e"
	clinicAddress  = "Москва, ул. Тестовая, 1"
	sourceSystem   = "mis-demo"
)

var (
	api *apiClient
	b2b *fakeB2B
	mis *fakeMIS
)

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

	nw, err := network.New(ctx)
	if err != nil {
		return 0, fmt.Errorf("create network: %w", err)
	}
	defer func() { err = errors.Join(err, nw.Remove(ctx)) }()

	b2b = newFakeB2B()
	b2bServer := httptest.NewServer(b2b)
	defer b2bServer.Close()

	mis = newFakeMIS()
	misServer := httptest.NewServer(mis)
	defer misServer.Close()

	b2bPort, err := portOf(b2bServer.URL)
	if err != nil {
		return 0, err
	}
	misPort, err := portOf(misServer.URL)
	if err != nil {
		return 0, err
	}

	vars, err := readEnvFile(defaultEnvFile)
	if err != nil {
		return 0, err
	}
	maps.Copy(vars, map[string]string{
		"B2B_URL":           fmt.Sprintf("http://%s:%d", testcontainers.HostInternal, b2bPort),
		"MIS_URL":           fmt.Sprintf("http://%s:%d", testcontainers.HostInternal, misPort),
		"MIS_SOURCE_SYSTEM": sourceSystem,
		"CLINIC_PHONE":      clinicPhone,
		"CLINIC_NAME":       clinicName,
		"CLINIC_ADDRESS":    clinicAddress,
	})

	ctr, err := testcontainers.Run(ctx, "",
		testcontainers.WithDockerfile(testcontainers.FromDockerfile{Context: "../.."}),
		testcontainers.WithEnv(vars),
		testcontainers.WithExposedPorts("8082/tcp"),
		testcontainers.WithHostPortAccess(b2bPort, misPort),
		network.WithNetwork([]string{"b2c-service"}, nw),
		testcontainers.WithWaitStrategy(wait.ForLog("http server started").WithStartupTimeout(3*time.Minute)),
	)
	if ctr != nil {
		defer func() { err = errors.Join(err, testcontainers.TerminateContainer(ctr)) }()
	}
	if err != nil {
		return 0, fmt.Errorf("start b2c-service: %w", err)
	}

	endpoint, err := ctr.PortEndpoint(ctx, "8082/tcp", "http")
	if err != nil {
		return 0, fmt.Errorf("service endpoint: %w", err)
	}
	api = newAPIClient(endpoint)

	code = m.Run()
	if code != 0 {
		dumpLogs(ctx, ctr)
	}
	return code, nil
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
	fmt.Fprintln(os.Stderr, "----- b2c-service logs -----")
	_, _ = io.Copy(os.Stderr, logs)
}
