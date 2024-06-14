package constants

const (
	DataSourceTypeMongo         = "mongo"
	DataSourceTypeMysql         = "mysql"
	DataSourceTypePostgresql    = "postgresql"
	DataSourceTypeMssql         = "mssql"
	DataSourceTypeSqlite        = "sqlite"
	DataSourceTypeCockroachdb   = "cockroachdb"
	DataSourceTypeElasticSearch = "elasticsearch"
	DataSourceTypeKafka         = "kafka"
)

const (
	DefaultHost = "localhost"
)

const (
	DefaultMongoPort         = "27017"
	DefaultMysqlPort         = "3306"
	DefaultPostgresqlPort    = "5432"
	DefaultMssqlPort         = "1433"
	DefaultCockroachdbPort   = "26257"
	DefaultElasticsearchPort = "9200"
	DefaultKafkaPort         = "9092"
)

const (
	DataSourceStatusOnline  = "on"
	DataSourceStatusOffline = "off"
)
