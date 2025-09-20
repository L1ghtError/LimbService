package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	testCases := []struct {
		test       string
		configPath string
		expected   error
	}{
		{
			test:       "Load existing .env",
			configPath: "../../.env",
			expected:   nil,
		}, {
			test:       "Load non-exising default conf",
			configPath: "",
			expected:   ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := NewGoDotEnv().Init(tc.configPath)

			if tc.expected != err {
				t.Errorf("NewGoDotEnv().Init(\"%s\") = %v; expected: %v", tc.configPath, err, tc.expected)
			}
		})
	}
}

func TestLoadExample(t *testing.T) {
	testCases := []struct {
		test       string
		configPath string
		configKey  string
		expected   string
	}{
		{
			test:       "Load param `APP_HOST` from example",
			configPath: "../../config_example.env",
			configKey:  AppHostKey,
			expected:   `localhost`,
		}, {
			test:       "Load non-exising param from example",
			configPath: "../../config_example.env",
			configKey:  "NON_EXISTING",
			expected:   ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			conf := NewGoDotEnv()
			err := conf.Init(tc.configPath)
			if err != nil {
				t.Errorf("NewGoDotEnv().Init(\"%s\") returns error: %v", tc.configPath, err)
			}
			val := conf.Get(tc.configKey)
			if val != tc.expected {
				t.Errorf("conf.Get(\"%s\") = %s; expected: %v", tc.configKey, val, tc.expected)
			}
		})
	}
}
