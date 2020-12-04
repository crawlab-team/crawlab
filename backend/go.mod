module crawlab

go 1.15

replace (
	github.com/crawlab-team/crawlab-core => /Users/marvzhang/projects/crawlab-team/crawlab-core
	github.com/crawlab-team/crawlab-db => /Users/marvzhang/projects/crawlab-team/crawlab-db
)

require (
	github.com/apex/log v1.9.0
	github.com/crawlab-team/crawlab-core v0.0.0-00010101000000-000000000000
	github.com/crawlab-team/crawlab-db v0.0.2
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0
	github.com/olivere/elastic/v7 v7.0.15
	github.com/spf13/viper v1.7.1
)
