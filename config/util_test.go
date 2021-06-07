package config

import "testing"

func TestParseConfigType(t *testing.T) {
	var testData = []struct {
		filePath string
		expected string
	}{
		{
			filePath: "",
			expected: "",
		},
		{
			filePath: ".yaml",
			expected: "yaml",
		},
		{
			filePath: "test.yaml",
			expected: "yaml",
		},
	}

	for i, v := range testData {
		got := parseConfigType(v.filePath)
		if got != v.expected {
			t.Errorf("TestParseConfigType[%d]: got %q, expected %q", i, got, v.expected)
		}
	}
}
