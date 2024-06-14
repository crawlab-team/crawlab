package mongo

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMongoGetDb(t *testing.T) {
	dbName := "test_db"
	viper.Set("mongo.db", dbName)
	err := InitMongo()
	require.Nil(t, err)

	db := GetMongoDb("")
	require.Equal(t, dbName, db.Name())
}
