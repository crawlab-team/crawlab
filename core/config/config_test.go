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
	// Create a new Config instance
	c := Config{Name: ""}

	// Initialize the config
	err := c.Init()
	require.NoError(t, err, "Failed to initialize config")

	// Test default values
	assert.Equal(t, "global.edition.community", viper.GetString("edition"), "Unexpected default value for edition")
	assert.Equal(t, "localhost", viper.GetString("mongo.host"), "Unexpected default value for mongo.host")
	assert.Equal(t, 27017, viper.GetInt("mongo.port"), "Unexpected default value for mongo.port")
	assert.Equal(t, "crawlab_test", viper.GetString("mongo.db"), "Unexpected default value for mongo.db")
	assert.Equal(t, "0.0.0.0", viper.GetString("server.host"), "Unexpected default value for server.host")
	assert.Equal(t, 8000, viper.GetInt("server.port"), "Unexpected default value for server.port")
	assert.Equal(t, "localhost:9666", viper.GetString("grpc.address"), "Unexpected default value for grpc.address")
	assert.Equal(t, "Crawlab2021!", viper.GetString("grpc.authKey"), "Unexpected default value for grpc.authKey")
	assert.Equal(t, "http://localhost:8000", viper.GetString("api.endpoint"), "Unexpected default value for api.endpoint")
	assert.Equal(t, "/var/log/crawlab", viper.GetString("log.path"), "Unexpected default value for log.path")

	// Test environment variable override
	os.Setenv("CRAWLAB_MONGO_HOST", "mongodb.example.com")
	defer os.Unsetenv("CRAWLAB_MONGO_HOST")
	assert.Equal(t, "mongodb.example.com", viper.GetString("mongo.host"), "Environment variable should override default value")

	// Test with a config file
	tempDir, err := os.MkdirTemp("", "crawlab-config-test")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	configContent := []byte(`
edition: global.edition.pro
mongo:
  host: mongodb.custom.com
  port: 27018
server:
  port: 8001
`)
	configPath := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, configContent, 0644)
	require.NoError(t, err, "Failed to write config file")

	// Create a new Config instance with the config file
	cWithFile := Config{Name: configPath}
	err = cWithFile.Init()
	require.NoError(t, err, "Failed to initialize config with file")

	// Test values from config file
	assert.Equal(t, "global.edition.pro", viper.GetString("edition"), "Unexpected value for edition from config file")
	assert.Equal(t, "mongodb.custom.com", viper.GetString("mongo.host"), "Unexpected value for mongo.host from config file")
	assert.Equal(t, 27018, viper.GetInt("mongo.port"), "Unexpected value for mongo.port from config file")
	assert.Equal(t, 8001, viper.GetInt("server.port"), "Unexpected value for server.port from config file")

	// Values not in config file should still use defaults
	assert.Equal(t, "Crawlab2021!", viper.GetString("grpc.authKey"), "Unexpected default value for grpc.authKey")
}
