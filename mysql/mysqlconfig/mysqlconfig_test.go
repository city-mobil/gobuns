package mysqlconfig

import (
	"testing"
	"time"
)

func TestEscapeTimezone(t *testing.T) {
	var testData = []struct {
		testName string
		timezone string
		expected string
	}{
		{
			testName: "empty_zone",
			timezone: "",
			expected: "",
		},
		{
			testName: "with_slash",
			timezone: "Europe/Moscow",
			expected: "Europe%2FMoscow",
		},
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			got := escapeTimezone(v.timezone)
			if got != v.expected {
				t.Errorf("got %q, expected %q", got, v.expected)
			}
		})
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	var testData = []struct {
		testName string
		cfg      *DatabaseConfig
		expected string
	}{
		{
			testName: "nil_config",
			cfg:      nil,
			expected: ":@tcp()/?timeout=1s&readTimeout=1s&writeTimeout=1s&interpolateParams=true&charset=utf8mb4&parseTime=false",
		},
		{
			testName: "valid_config",
			cfg: &DatabaseConfig{
				Addr:              "42",
				Username:          "test",
				Password:          "test",
				DatabaseName:      "some_db",
				Charset:           "utf8mb4",
				Collation:         "utf8mb42",
				VitessReplicaType: "replica",
				Driver:            "sql",
				Name:              "testName",
				ReadTimeout:       time.Second,
				WriteTimeout:      time.Second,
				Timeout:           time.Second,
				ParseTime:         false,
				InterpolateParams: false,
				Timezone:          "UTC",
			},
			expected: "test:test@tcp(42)/some_db@replica?timeout=1s&readTimeout=1s&writeTimeout=1s&interpolateParams=false&" +
				"charset=utf8mb4&parseTime=false&collation=utf8mb42&loc=" + escapeTimezone("UTC"),
		},
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			got := v.cfg.DSN()
			if got != v.expected {
				t.Errorf("got %q, expected %q", got, v.expected)
			}
		})
	}
}
