const optimization = {
  splitChunks: {
    chunks: 'initial',
    minSize: 20000,
    minChunks: 1,
    maxAsyncRequests: 3,
    cacheGroups: {
      defaultVendors: {
        test: /[\\/]node_modules[\\/]]/,
        priority: -10,
        reuseExistingChunk: true
      },
      default: {
        minChunks: 2,
        priority: -20,
        reuseExistingChunk: true
      }
    }
  }
}

const config = {
  pages: {
    index: {
      entry: 'src/main.ts',
      template: 'public/index.html',
      filename: 'index.html',
      title: 'Crawlab | Distributed Web Crawler Platform'
    }
  },
  outputDir: './dist',
  configureWebpack: {
    optimization,
    plugins: []
  }
}

module.exports = config
