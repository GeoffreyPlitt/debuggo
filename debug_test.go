package debuggo

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDebugEnableDisable(t *testing.T) {
	// Test cases for various DEBUG env var settings
	testCases := []struct {
		envValue      string
		module        string
		expectEnabled bool
		description   string
	}{
		{"", "module", false, "No DEBUG env var should disable all"},
		{"*", "module", true, "Wildcard should enable all modules"},
		{"module", "module", true, "Exact match should enable"},
		{"other", "module", false, "Non-match should disable"},
		{"module:*", "module:submodule", true, "Wildcard should enable submodules"},
		{"*,!module", "module", false, "Negation should override wildcard"},
		{"module,submodule", "module", true, "Multiple modules should work"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Set environment variable
			os.Setenv("DEBUG", tc.envValue)
			ReloadDebugSettings()

			// Check if enabled
			if IsEnabled(tc.module) != tc.expectEnabled {
				t.Errorf("Expected module %s to be enabled=%v with DEBUG=%s",
					tc.module, tc.expectEnabled, tc.envValue)
				return
			}

			// Test output if enabled
			r, w, _ := os.Pipe()
			origStderr := os.Stderr
			os.Stderr = w

			// Call debug
			debug := Debug(tc.module)
			debug("Test message")

			// Need a small delay to ensure buffer gets written before checking
			time.Sleep(10 * time.Millisecond)

			// Close writer to signal EOF to reader
			w.Close()

			// Read the output
			buf := new(bytes.Buffer)
			buf.ReadFrom(r)

			// Restore stderr
			os.Stderr = origStderr

			// Check if output matches expectation
			hasOutput := buf.Len() > 0
			if hasOutput != tc.expectEnabled {
				t.Errorf("Expected output=%v but got output=%v for DEBUG=%s, module=%s",
					tc.expectEnabled, hasOutput, tc.envValue, tc.module)
			}
		})
	}
}

func TestHierarchicalNamespaces(t *testing.T) {
	testCases := []struct {
		envValue      string
		module        string
		expectEnabled bool
		description   string
	}{
		{"app:*", "app", false, "Wildcard applies only to children"},
		{"app:*", "app:server", true, "Wildcard enables direct children"},
		{"app:*", "app:server:http", true, "Wildcard enables all descendants"},
		{"app:server:*", "app:database", false, "Wildcard doesn't affect siblings"},
		{"app:server:*", "app:server:http", true, "Nested wildcard works"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			os.Setenv("DEBUG", tc.envValue)
			ReloadDebugSettings()

			if IsEnabled(tc.module) != tc.expectEnabled {
				t.Errorf("Expected module %s to be enabled=%v with DEBUG=%s",
					tc.module, tc.expectEnabled, tc.envValue)
			}
		})
	}
}

func TestNegation(t *testing.T) {
	testCases := []struct {
		envValue      string
		module        string
		expectEnabled bool
		description   string
	}{
		{"*,!app:server", "app:server", false, "Negation should override wildcard"},
		{"app:*,!app:server", "app:client", true, "Negation should not affect siblings"},
		{"*,!app:*", "app:server", false, "Nested negation should work"},
		{"*,!app:*", "api", true, "Nested negation should not affect others"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			os.Setenv("DEBUG", tc.envValue)
			ReloadDebugSettings()

			if IsEnabled(tc.module) != tc.expectEnabled {
				t.Errorf("Expected module %s to be enabled=%v with DEBUG=%s",
					tc.module, tc.expectEnabled, tc.envValue)
			}
		})
	}
}

func TestPrefixWriter(t *testing.T) {
	// Save original stderr
	origStderr := os.Stderr
	defer func() { os.Stderr = origStderr }()

	r, w, _ := os.Pipe()
	os.Stderr = w

	// Test basic output
	testPrefix := "TEST-PREFIX"
	pw := &PrefixWriter{Prefix: testPrefix}

	msg := "Hello, world\n"
	n, err := pw.Write([]byte(msg))

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if n != len(msg) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(msg), n)
	}

	w.Close()

	// Read output
	buf := &bytes.Buffer{}
	_, _ = buf.ReadFrom(r)

	// Check output
	expected := testPrefix + " " + msg
	if buf.String() != expected {
		t.Errorf("Expected output '%s', got '%s'", expected, buf.String())
	}
}

func TestPrefixWriterIgnore(t *testing.T) {
	// Save original stderr
	origStderr := os.Stderr
	defer func() { os.Stderr = origStderr }()

	// Test ignored phrases
	r, w, _ := os.Pipe()
	os.Stderr = w

	pw := &PrefixWriter{
		Prefix:  "TEST",
		Ignores: []string{"ignore me"},
	}

	// This should be ignored
	pw.Write([]byte("This text contains ignore me phrase"))

	// This should be printed
	pw.Write([]byte("This text should appear"))

	w.Close()

	// Read output
	buf := &bytes.Buffer{}
	_, _ = buf.ReadFrom(r)

	// Check output - should only contain the second message
	if !strings.Contains(buf.String(), "This text should appear") {
		t.Error("Expected text should appear in output")
	}

	if strings.Contains(buf.String(), "ignore me") {
		t.Error("Ignored text should not appear in output")
	}
}

func TestReloadDebugSettings(t *testing.T) {
	// Initial setup
	os.Setenv("DEBUG", "module1")
	ReloadDebugSettings()

	if !IsEnabled("module1") {
		t.Error("module1 should be enabled")
	}

	if IsEnabled("module2") {
		t.Error("module2 should be disabled")
	}

	// Change environment and reload
	os.Setenv("DEBUG", "module2")
	ReloadDebugSettings()

	if IsEnabled("module1") {
		t.Error("module1 should be disabled after reload")
	}

	if !IsEnabled("module2") {
		t.Error("module2 should be enabled after reload")
	}
}
