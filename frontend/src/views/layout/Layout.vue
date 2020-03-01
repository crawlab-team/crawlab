<template>
  <div :class="classObj" class="app-wrapper">
    <div v-if="device==='mobile'&&sidebar.opened" class="drawer-bg" @click="handleClickOutside"></div>

    <!--sidebar-->
    <sidebar class="sidebar-container"/>
    <!--./sidebar-->

    <!--main container-->
    <div class="main-container">
      <navbar/>
      <tags-view/>
      <app-main/>
    </div>
    <!--./main container-->

    <!--documentation-->
    <div class="documentation">
      <el-tooltip
        :content="$t('Click to view related Documentation')"
      >
        <i class="el-icon-question" @click="onClickDocumentation"></i>
      </el-tooltip>
      <el-drawer
        :title="$t('Related Documentation')"
        :visible.sync="isShowDocumentation"
        :before-close="onCloseDocumentation"
        size="360px"
      >
        <documentation/>
      </el-drawer>
    </div>
    <!--./documentation-->
  </div>
</template>

<script>
import {
  Navbar,
  Sidebar,
  AppMain,
  TagsView
} from './components'
import ResizeMixin from './mixin/ResizeHandler'
import Documentation from '../../components/Documentation/Documentation'

export default {
  name: 'Layout',
  components: {
    Documentation,
    Navbar,
    Sidebar,
    TagsView,
    AppMain
  },
  mixins: [ResizeMixin],
  data () {
    return {
      isShowDocumentation: false
    }
  },
  computed: {
    sidebar () {
      return this.$store.state.app.sidebar
    },
    device () {
      return this.$store.state.app.device
    },
    classObj () {
      return {
        hideSidebar: !this.sidebar.opened,
        openSidebar: this.sidebar.opened,
        withoutAnimation: this.sidebar.withoutAnimation,
        mobile: this.device === 'mobile'
      }
    }
  },
  methods: {
    handleClickOutside () {
      this.$store.dispatch('CloseSideBar', { withoutAnimation: false })
    },
    onClickDocumentation () {
      this.isShowDocumentation = true

      this.$st.sendEv('全局', '点击页面文档')
    },
    onCloseDocumentation () {
      this.isShowDocumentation = false
    }
  },
  async created () {
    await this.$store.dispatch('doc/getDocData')
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
  @import "../../../src/styles/mixin.scss";

  .app-wrapper {
    @include clearfix;
    position: relative;
    height: 100%;
    width: 100%;
    background: white;

    &.mobile.openSidebar {
      position: fixed;
      top: 0;
    }
  }

  .drawer-bg {
    background: #000;
    opacity: 0.3;
    width: 100%;
    top: 0;
    height: 100%;
    position: absolute;
    z-index: 999;
  }

  .documentation {
    z-index: 9999;
    position: fixed;
    right: 25px;
    bottom: 20px;
    font-size: 24px;
    cursor: pointer;
    color: #909399;
  }
</style>

<style scoped>
  .documentation .el-tree {
    margin-left: 20px;
  }

  .documentation >>> span[role="heading"]:focus {
    outline: none;
  }

  .documentation >>> .el-tree-node__content {
    height: 40px;
    line-height: 40px;
  }

  .documentation >>> .custom-tree-node {
    display: block;
    width: 100%;
    height: 40px;
    line-height: 40px;
    font-size: 14px;
  }

  .documentation >>> .custom-tree-node a {
    display: block;
  }

  .documentation >>> .custom-tree-node:hover a {
    text-decoration: underline;
  }
</style>
