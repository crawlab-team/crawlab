module crawlab

go 1.15

replace (
	github.com/crawlab-team/crawlab-core => /Users/marvzhang/projects/crawlab-team/crawlab-core
	github.com/crawlab-team/crawlab-db => /Users/marvzhang/projects/crawlab-team/crawlab-db
	github.com/crawlab-team/crawlab-fs => /Users/marvzhang/projects/crawlab-team/crawlab-fs
	github.com/crawlab-team/crawlab-grpc => /Users/marvzhang/projects/crawlab-team/crawlab-grpc/dist/go
	github.com/crawlab-team/crawlab-log => /Users/marvzhang/projects/crawlab-team/crawlab-log
	github.com/crawlab-team/crawlab-vcs => /Users/marvzhang/projects/crawlab-team/crawlab-vcs
	github.com/crawlab-team/goseaweedfs => /Users/marvzhang/projects/crawlab-team/goseaweedfs
	github.com/crawlab-team/go-trace => /Users/marvzhang/projects/crawlab-team/go-trace
)

require (
	github.com/apex/log v1.9.0
	github.com/crawlab-team/crawlab-core v0.0.0
	github.com/crawlab-team/crawlab-db v0.0.2
	github.com/crawlab-team/go-trace v0.1.0
	github.com/gin-gonic/gin v1.6.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	go.uber.org/dig v1.10.0
)
