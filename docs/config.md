# Config

#### Environment Variable List

Name | Description | Example | Default
---|--- | --- | ---
CRAWLAB_GRPC_ADDRESS| Target gRPC address the nodes are connecting to | 192.168.0.1:9666 | localhost:9666
CRAWLAB_GRPC_SERVER_ADDRESS| Address that the gRPC server is listening to (master node only) | 0.0.0.0:9666 | 0.0.0.0:9666
CRAWLAB_GRPC_AUTHKEY| The token that gRPC clients and server use for authentication | youcanneverguess | Crawlab2021!
CRAWLAB_NODE_MASTER | Whether the current node is a master or worker node (Y: master; N: worker) | Y | Y
CRAWLAB_SERVER_HOST | IP host that the API listens to (master node only) | 0.0.0.0 | 0.0.0.0
CRAWLAB_SERVER_PORT | IP port that the API listens to (master node only) | 8000 | 8000
CRAWLAB_TASK_HANDLER_MAXRUNNERS | Max number of task runners (concurrent spider tasks) that a node can run | 16 | 8
CRAWLAB_FS_FILER_PROXY |Filer API endpoint that Crawlab's filer proxy links to |http://filer-server:8888 | http://localhost:8888
CRAWLAB_FS_FILER_URL |Crawlab's Filer API endpoint |http://crawlab-web-api:8000/filer | http://localhost:8000/filer
CRAWLAB_FS_FILER_AUTHKEY |Crawlab's Filer API auth key token | youcanneverguess | Crawlab2021!

