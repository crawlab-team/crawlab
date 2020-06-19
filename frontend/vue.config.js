'use strict'
const path = require('path')

function resolve(dir) {
  return path.join(__dirname, dir)
}

const isDev = process.env.NODE_ENV === 'development'
module.exports = {
  // TODO: need to configure output static files with hash
  lintOnSave: isDev,

  productionSourceMap: false,

  configureWebpack: {
    // provide the app's title in webpack's name field, so that
    // it can be accessed in index.html to inject the correct title.
    // name: name,
    devtool: 'source-map',
    resolve: {
      alias: {
        '@': resolve('src')
      }
    }
  },

  chainWebpack: config => {

  },

  css: {
    sourceMap: true
  }
}
