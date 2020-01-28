export default {
  // 内容
  addNodeInstruction: `
You cannot add nodes directly on the web interface in Crawlab.

Adding a node is quite simple. The only thing you have to do is to run a Crawlab service on your target machine.

#### Docker Deployment
If you are running Crawlab using Docker, you can start a new \`worker\` container on the target machine, or add a \`worker\` service in the \`docker-compose.yml\`.

\`\`\`bash
docker run -d --restart always --name crawlab_worker \\
  -e CRAWLAB_SERVER_MASTER=N \\
  -e CRAWLAB_MONGO_HOST=xxx.xxx.xxx.xxx \\ # make sure you are connecting to the same MongoDB
  -e CRAWLAB_REDIS_ADDRESS=xxx.xxx.xxx.xxx \\ # make sure you are connecting to the same Redis
  tikazyq/crawlab:latest
\`\`\`

#### Direct Deploy
If you are deploying directly, the only thing you have to do is to run a backend service on the target machine, you can refer to [Direct Deploy](https://docs.crawlab.cn/Installation/Direct.html).

For more information, please refer to the [Official Documentation](https://docs.crawlab.cn).
`
}
