package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Set test environment variables
	viper.Set("CRAWLAB_TEST_STRING", "test_string_value")
	viper.Set("CRAWLAB_TEST_INT", 42)
	viper.Set("CRAWLAB_TEST_BOOL", true)
	viper.Set("CRAWLAB_TEST_NESTED_KEY", "nested_value")
}

func TestInitConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "crawlab-config-test")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	// Create a temporary config file
	configContent := []byte(`
log:
  level: info
test:
  string: default_string_value
  int: 0
  bool: false
  nested:
    key: default_nested_value
`)
	configPath := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, configContent, 0644)
	require.NoError(t, err, "Failed to write config file")

	// Set up the test environment
	oldConfigPath := viper.ConfigFileUsed()
	defer viper.SetConfigFile(oldConfigPath)
	viper.SetConfigFile(configPath)

	// Create a new Config instance
	c := Config{Name: configPath}

	// Initialize the config
	err = c.Init()
	require.NoError(t, err, "Failed to initialize config")

	// Test config values
	assert.Equal(t, "default_string_value", viper.GetString("test.string"), "Unexpected value for test.string")
	assert.Equal(t, 0, viper.GetInt("test.int"), "Unexpected value for test.int")
	assert.False(t, viper.GetBool("test.bool"), "Unexpected value for test.bool")
	assert.Equal(t, "default_nested_value", viper.GetString("test.nested.key"), "Unexpected value for test.nested.key")
	assert.Empty(t, viper.GetString("non.existent.key"), "Non-existent key should return empty string")

	// Test environment variable override
	os.Setenv("CRAWLAB_TEST_STRING", "env_string_value")
	defer os.Unsetenv("CRAWLAB_TEST_STRING")
	assert.Equal(t, "env_string_value", viper.GetString("test.string"), "Environment variable should override config value")
}
