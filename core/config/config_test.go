package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitConfig(t *testing.T) {
	err := InitConfig()
	require.Nil(t, err)
}
