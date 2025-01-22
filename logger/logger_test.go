package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/chtc/chtc-go-logger/config"
)

// TestContextAwareLogger validates that the logger correctly extracts context attributes,
// merges them with additional attributes, and logs them with the appropriate log level.
func TestContextAwareLogger(t *testing.T) {
	tempFile, err := os.CreateTemp("", "logfile-*.log")
	if err != nil {
		t.Fatalf("failed to create temporary log file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the temp file

	// Create a configuration with the overridden file name
	cfg := &config.Config{
		FileOutput: config.FileOutputConfig{
			Enabled:  true,
			FilePath: tempFile.Name(), // Override file path
		},
	}

	// Initialize the logger
	contextLogger, err := NewContextAwareLogger(cfg)
	if err != nil {
		t.Fatalf("failed to initialize context-aware logger: %v", err)
	}

	ctx := context.WithValue(context.Background(), LogAttrsKey, map[string]string{
		"user_id":    "12345",
		"request_id": "abcde",
	})

	contextLogger.Info(ctx, "Test message", slog.String("extra_key", "extra_value"))

	tempFile.Close()

	// Read the log file and validate its contents
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to read temporary log file: %v", err)
	}

	logContents := string(content)
	expectedValues := []string{
		`"user_id":"12345"`,
		`"request_id":"abcde"`,
		`"extra_key":"extra_value"`,
		`"msg":"Test message"`,
		`"level":"INFO"`,
	}
	for _, value := range expectedValues {
		if !contains(logContents, value) {
			t.Errorf("log does not contain expected value: %s", value)
		}
	}
}

// Helper function to check if a string is contained
func contains(content, substring string) bool {
	return len(content) >= len(substring) && strings.Contains(content, substring)
}
