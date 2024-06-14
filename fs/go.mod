module github.com/crawlab-team/crawlab/fs

go 1.22

replace github.com/crawlab-team/crawlab/trace => ../trace

require (
	github.com/apex/log v1.9.0
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/crawlab-team/crawlab/trace v0.0.0-20240614094818-e8f694eab76e
	github.com/crawlab-team/goseaweedfs v0.6.0-beta.20211101.1936.0.20220912021203-dfee5f74dd69
	github.com/emirpasic/gods v1.18.1
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/linxGnu/gumble v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/ztrue/tracerr v0.4.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
