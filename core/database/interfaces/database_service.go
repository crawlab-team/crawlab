package interfaces

import (
	"github.com/crawlab-team/crawlab/core/database/entity"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DatabaseService interface {
	TestConnection(id primitive.ObjectID) (err error)
	GetMetadata(id primitive.ObjectID) (m *entity.DatabaseMetadata, err error)
	GetMetadataAllDb(id primitive.ObjectID) (m *entity.DatabaseMetadata, err error)
	CreateDatabase(id primitive.ObjectID, databaseName string) (err error)
	DropDatabase(id primitive.ObjectID, databaseName string) (err error)
	GetTableMetadata(id primitive.ObjectID, databaseName, tableName string) (table *entity.DatabaseTable, err error)
	CreateTable(id primitive.ObjectID, databaseName string, table *entity.DatabaseTable) (err error)
	ModifyTable(id primitive.ObjectID, databaseName string, table *entity.DatabaseTable) (err error)
	DropTable(id primitive.ObjectID, databaseName, tableName string) (err error)
	RenameTable(id primitive.ObjectID, databaseName, oldTableName, newTableName string) (err error)
	GetColumnTypes(query string) (types []string)
	ReadRows(id primitive.ObjectID, databaseName, tableName string, filter map[string]interface{}, skip, limit int) ([]map[string]interface{}, int64, error)
	CreateRow(id primitive.ObjectID, databaseName, tableName string, row map[string]interface{}) error
	UpdateRow(id primitive.ObjectID, databaseName, tableName string, filter map[string]interface{}, update map[string]interface{}) error
	DeleteRow(id primitive.ObjectID, databaseName, tableName string, filter map[string]interface{}) error
	Query(id primitive.ObjectID, databaseName, query string) (results *entity.DatabaseQueryResults, err error)
	GetCurrentMetric(id primitive.ObjectID) (m *models.DatabaseMetricV2, err error)
}
