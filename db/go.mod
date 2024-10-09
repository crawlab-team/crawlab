module github.com/crawlab-team/crawlab/db

go 1.22

replace github.com/crawlab-team/crawlab/trace => ../trace

require (
	github.com/apex/log v1.9.0
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/crawlab-team/crawlab/trace v0.0.0-20240614095218-7b4ee8399ab0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
	go.mongodb.org/mongo-driver v1.15.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/ztrue/tracerr v0.4.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
