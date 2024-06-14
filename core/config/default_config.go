package config

var DefaultConfigYaml = `
info:
  version: v0.6.3
  edition: global.edition.community
mongo:
  host: localhost
  port: 27017
  db: crawlab_test
  username: ""
  password: ""
  authSource: "admin"
server:
  host: 0.0.0.0
  port: 8000
spider:
  fs: "/spiders"
  workspace: "/workspace"
  repo: "/repo"
task:
  workers: 16
  cancelWaitSeconds: 30
grpc:
  address: localhost:9666
  server:
    address: 0.0.0.0:9666
  authKey: Crawlab2021!
fs:
  filer:
    proxy: http://localhost:8888
    url: http://localhost:8000/filer
    authKey: Crawlab2021!
node:
  master: Y
api:
  endpoint: http://localhost:8000
log:
  path: /var/log/crawlab
`
