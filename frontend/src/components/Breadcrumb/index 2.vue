<template>
  <el-breadcrumb class="app-breadcrumb" separator="/">
    <transition-group name="breadcrumb">
      <el-breadcrumb-item v-for="(item,index) in levelList" :key="item.path">
        <span v-if="item.redirect==='noredirect'||index==levelList.length-1"
              class="no-redirect">{{$t(item.meta.title) }}</span>
        <a v-else @click.prevent="handleLink(item)">{{ $t(item.meta.title) }}</a>
      </el-breadcrumb-item>
    </transition-group>
  </el-breadcrumb>
</template>

<script>
import pathToRegexp from 'path-to-regexp'

export default {
  data () {
    return {
      levelList: null
    }
  },
  watch: {
    $route () {
      this.getBreadcrumb()
    }
  },
  created () {
    this.getBreadcrumb()
  },
  methods: {
    getBreadcrumb () {
      let matched = this.$route.matched.filter(item => item.name)

      const first = matched[0]
      if (first && first.name !== 'Home') {
        matched = [{ path: '/home', meta: { title: 'Home' } }].concat(matched)
      }

      this.levelList = matched.filter(item => item.meta && item.meta.title && item.meta.breadcrumb !== false)
    },
    pathCompile (path) {
      // To solve this problem https://github.com/PanJiaChen/vue-element-admin/issues/561
      const { params } = this.$route
      var toPath = pathToRegexp.compile(path)
      return toPath(params)
    },
    handleLink (item) {
      const { redirect } = item
      if (redirect) {
        this.$router.push(redirect)
        return
      }
      this.$router.push(this.getGoToPath(item))
    },
    getGoToPath (item) {
      if (item.path) {
        var path = item.path
        var startPos = path.indexOf(':')

        if (startPos !== -1) {
          var endPos = path.indexOf('/', startPos)
          var key = path.substring(startPos + 1, endPos)
          path = path.replace(':' + key, this.$route.params[key])
          return path
        }
      }

      return item.redirect || item.path
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
  .app-breadcrumb.el-breadcrumb {
    display: inline-block;
    font-size: 14px;
    line-height: 50px;
    margin-left: 10px;

    .no-redirect {
      color: #97a8be;
      cursor: text;
    }
  }
</style>
