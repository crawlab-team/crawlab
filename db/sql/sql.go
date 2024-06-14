package sql

import (
	"errors"
	"fmt"
	"github.com/crawlab-team/go-trace"
	"github.com/jmoiron/sqlx"
)

func GetSqlDatabaseConnectionString(dataSourceType string, host string, port string, username string, password string, database string) (connStr string, err error) {
	if dataSourceType == "mysql" {
		connStr = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)
	} else if dataSourceType == "postgres" {
		connStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, username, database, password, "disable")
	} else {
		err = errors.New(dataSourceType + " is not implemented")
		return connStr, trace.TraceError(err)
	}
	return connStr, nil
}

func GetSqlConn(dataSourceType string, host string, port string, username string, password string, database string) (db *sqlx.DB, err error) {
	// get database connection string
	connStr, err := GetSqlDatabaseConnectionString(dataSourceType, host, port, username, password, database)
	if err != nil {
		return db, trace.TraceError(err)
	}

	// get database instance
	db, err = sqlx.Open(dataSourceType, connStr)
	if err != nil {
		return db, trace.TraceError(err)
	}

	return db, nil
}
