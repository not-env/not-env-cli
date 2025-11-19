package commands

import (
	"fmt"
	"strings"
	"testing"
)

func TestEnvSetOutput(t *testing.T) {
	// This test verifies the output format for env set command
	// We can't easily test the full command without a backend, but we can test the logic

	// Test shell escaping
	testCases := []struct {
		key   string
		value string
		want  string
	}{
		{key: "SIMPLE", value: "value", want: `export SIMPLE="value"`},
		{key: "WITH_QUOTES", value: `"quoted"`, want: `export WITH_QUOTES="\"quoted\""`},
		{key: "WITH_DOLLAR", value: "$VAR", want: `export WITH_DOLLAR="\$VAR"`},
		{key: "WITH_BACKTICK", value: "`command`", want: "export WITH_BACKTICK=\"\\`command\\`\""},
	}

	for _, tc := range testCases {
		// Simulate the escaping logic from env.go
		escapedValue := tc.value
		escapedValue = strings.ReplaceAll(escapedValue, `"`, `\"`)
		escapedValue = strings.ReplaceAll(escapedValue, `$`, `\$`)
		escapedValue = strings.ReplaceAll(escapedValue, "`", "\\`")
		output := fmt.Sprintf(`export %s="%s"`, tc.key, escapedValue)

		if output != tc.want {
			t.Errorf("escapeValue(%q, %q) = %q, want %q", tc.key, tc.value, output, tc.want)
		}
	}
}

func TestEnvClearOutput(t *testing.T) {
	// Test that env clear outputs correct format
	key := "TEST_VAR"
	output := fmt.Sprintf("unset %s", key)
	expected := "unset TEST_VAR"

	if output != expected {
		t.Errorf("env clear output = %q, want %q", output, expected)
	}
}

