//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

type serviceResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type servicesResponse struct {
	Services []serviceResponse `json:"services"`
}

func TestServiceLifecycle(t *testing.T) {
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve project root: %v", err)
	}

	binaryPath := filepath.Join(t.TempDir(), "platform-service")

	buildCommand := exec.Command(
		"go",
		"build",
		"-o",
		binaryPath,
		"./cmd/api",
	)
	buildCommand.Dir = projectRoot

	if output, err := buildCommand.CombinedOutput(); err != nil {
		t.Fatalf(
			"build API: %v\n%s",
			err,
			string(output),
		)
	}

	port := availablePort(t)
	databasePath := filepath.Join(t.TempDir(), "integration.db")

	ctx, cancel := context.WithCancel(context.Background())

	apiCommand := exec.CommandContext(ctx, binaryPath)
	apiCommand.Dir = projectRoot
	apiCommand.Env = environmentWith(
		"DB_DRIVER=sqlite",
		"DB_DSN="+databasePath,
		"PORT="+strconv.Itoa(port),
	)
	apiCommand.Stdout = os.Stdout
	apiCommand.Stderr = os.Stderr

	if err := apiCommand.Start(); err != nil {
		cancel()
		t.Fatalf("start API: %v", err)
	}

	processFinished := make(chan struct{})
	var processErr error

	go func() {
		processErr = apiCommand.Wait()
		close(processFinished)
	}()

	t.Cleanup(func() {
		cancel()

		select {
		case <-processFinished:
		case <-time.After(2 * time.Second):
			t.Log("API process did not stop within two seconds")
		}
	})

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	client := &http.Client{
		Timeout: time.Second,
	}

	if err := waitForAPI(
		client,
		baseURL+"/health",
		processFinished,
		&processErr,
	); err != nil {
		t.Fatal(err)
	}

	createBody := bytes.NewBufferString(
		`{"name":"integration-api","description":"Created by integration test"}`,
	)

	createRequest, err := http.NewRequest(
		http.MethodPost,
		baseURL+"/services",
		createBody,
	)
	if err != nil {
		t.Fatalf("create POST request: %v", err)
	}

	createRequest.Header.Set("Content-Type", "application/json")

	createResult, err := client.Do(createRequest)
	if err != nil {
		t.Fatalf("POST /services: %v", err)
	}
	defer createResult.Body.Close()

	if createResult.StatusCode != http.StatusCreated {
		t.Fatalf(
			"POST /services status: got %d, want %d",
			createResult.StatusCode,
			http.StatusCreated,
		)
	}

	var created serviceResponse

	if err := json.NewDecoder(createResult.Body).Decode(&created); err != nil {
		t.Fatalf("decode created service: %v", err)
	}

	if created.ID == "" {
		t.Fatal("created service has an empty ID")
	}

	if created.Name != "integration-api" {
		t.Fatalf(
			"created service name: got %q, want %q",
			created.Name,
			"integration-api",
		)
	}

	listResult, err := client.Get(baseURL + "/services")
	if err != nil {
		t.Fatalf("GET /services: %v", err)
	}
	defer listResult.Body.Close()

	if listResult.StatusCode != http.StatusOK {
		t.Fatalf(
			"GET /services status: got %d, want %d",
			listResult.StatusCode,
			http.StatusOK,
		)
	}

	var listed servicesResponse

	if err := json.NewDecoder(listResult.Body).Decode(&listed); err != nil {
		t.Fatalf("decode services list: %v", err)
	}

	if len(listed.Services) != 1 {
		t.Fatalf(
			"service count: got %d, want 1",
			len(listed.Services),
		)
	}

	if listed.Services[0].ID != created.ID {
		t.Fatalf(
			"listed service ID: got %q, want %q",
			listed.Services[0].ID,
			created.ID,
		)
	}
}

func availablePort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find available port: %v", err)
	}
	defer listener.Close()

	address, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener address: %T", listener.Addr())
	}

	return address.Port
}

func waitForAPI(
	client *http.Client,
	healthURL string,
	processFinished <-chan struct{},
	processErr *error,
) error {
	deadline := time.Now().Add(10 * time.Second)

	for time.Now().Before(deadline) {
		select {
		case <-processFinished:
			return fmt.Errorf(
				"API stopped before becoming ready: %v",
				*processErr,
			)
		default:
		}

		response, err := client.Get(healthURL)
		if err == nil {
			response.Body.Close()

			if response.StatusCode == http.StatusOK {
				return nil
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("API did not become ready within 10 seconds")
}

func environmentWith(values ...string) []string {
	replacements := make(map[string]struct{}, len(values))

	for _, value := range values {
		key, _, _ := strings.Cut(value, "=")
		replacements[key] = struct{}{}
	}

	environment := make([]string, 0, len(os.Environ())+len(values))

	for _, value := range os.Environ() {
		key, _, _ := strings.Cut(value, "=")

		if _, replaced := replacements[key]; !replaced {
			environment = append(environment, value)
		}
	}

	return append(environment, values...)
}
