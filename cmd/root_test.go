package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gkwa/littlewill/internal/logger"
	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCustomLogger(t *testing.T) {
	var buf bytes.Buffer

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLogger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapConfig.EncoderConfig),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	))

	customLogger := zapr.NewLogger(zapLogger)

	cliLogger = customLogger

	// Create a test command that uses the logger
	testCmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			logger := LoggerFrom(cmd.Context())
			logger.Info("Test log message")
		},
	}
	rootCmd.AddCommand(testCmd)

	// Execute the test command
	cmd := rootCmd
	cmd.SetArgs([]string{"test"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	logOutput := buf.String()
	if logOutput == "" {
		t.Error("Expected log output, but got none")
	}

	if !strings.Contains(logOutput, "Test log message") {
		t.Errorf("Expected log output to contain 'Test log message', but got: %s", logOutput)
	}

	t.Logf("Log output: %s", logOutput)
}

func TestJSONLogger(t *testing.T) {
	oldVerbose, oldLogFormat := verbose, logFormat
	verbose, logFormat = true, "json"
	defer func() {
		verbose, logFormat = oldVerbose, oldLogFormat
	}()

	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	customLogger := logger.NewConsoleLogger(verbose, logFormat == "json")
	cliLogger = customLogger

	// Create a test command that uses the logger
	testCmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			logger := LoggerFrom(cmd.Context())
			logger.Info("Test log message")
		},
	}
	rootCmd.AddCommand(testCmd)

	// Execute the test command
	cmd := rootCmd
	cmd.SetArgs([]string{"test"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to copy log output: %v", err)
	}
	logOutput := buf.String()

	if logOutput == "" {
		t.Error("Expected log output, but got none")
	}

	lines := strings.Split(strings.TrimSpace(logOutput), "\n")
	for _, line := range lines {
		var jsonMap map[string]interface{}
		err := json.Unmarshal([]byte(line), &jsonMap)
		if err != nil {
			t.Errorf("Expected valid JSON, but got error: %v", err)
		}

		if message, ok := jsonMap["message"]; ok {
			if message != "Test log message" {
				t.Errorf("Expected log message 'Test log message', but got: %v", message)
			}
		} else {
			t.Error("Expected 'message' field in JSON log output, but it was not found")
		}
	}

	t.Logf("Log output: %s", logOutput)
}
