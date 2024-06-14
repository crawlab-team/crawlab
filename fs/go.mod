module github.com/crawlab-team/crawlab/fs

go 1.22

replace (
	github.com/crawlab-team/crawlab/trace => ../trace
)

require (
	github.com/apex/log v1.9.0
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/crawlab-team/goseaweedfs v0.6.0-beta.20211101.1936.0.20220912021203-dfee5f74dd69
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.6.1
)
