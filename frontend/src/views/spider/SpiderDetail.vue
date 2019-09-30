<template>
  <div class="app-container">
    <!--selector-->
    <div class="selector">
      <label class="label">{{$t('Spider')}}: </label>
      <el-select v-model="spiderForm._id" @change="onSpiderChange">
        <el-option v-for="op in spiderList" :key="op._id" :value="op._id" :label="op.name"></el-option>
      </el-select>
    </div>

    <!--tabs-->
    <el-tabs v-model="activeTabName" @tab-click="onTabClick" type="border-card">
      <el-tab-pane :label="$t('Overview')" name="overview">
        <spider-overview/>
      </el-tab-pane>
      <el-tab-pane v-if="isConfigurable" :label="$t('Config')" name="配置">
        <config-list/>
      </el-tab-pane>
      <el-tab-pane v-if="isCustomized" :label="$t('Files')" name="files">
        <file-list/>
      </el-tab-pane>
      <el-tab-pane :label="$t('Environment')" name="environment">
        <environment-list/>
      </el-tab-pane>
      <el-tab-pane :label="$t('Analytics')" name="analytics">
        <spider-stats ref="spider-stats"/>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import FileList from '../../components/File/FileList'
import SpiderOverview from '../../components/Overview/SpiderOverview'
import EnvironmentList from '../../components/Environment/EnvironmentList'
import SpiderStats from '../../components/Stats/SpiderStats'
import ConfigList from '../../components/Config/ConfigList'

export default {
  name: 'SpiderDetail',
  components: {
    ConfigList,
    SpiderStats,
    EnvironmentList,
    FileList,
    SpiderOverview
  },
  data () {
    return {
      activeTabName: 'overview'
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderList',
      'spiderForm'
    ]),
    ...mapState('file', [
      'currentPath'
    ]),
    ...mapState('deploy', [
      'deployList'
    ]),
    isCustomized () {
      return this.spiderForm.type === 'customized'
    },
    isConfigurable () {
      return this.spiderForm.type === 'configurable'
    }
  },
  methods: {
    onTabClick (tab) {
      if (this.activeTabName === 'analytics') {
        setTimeout(() => {
          this.$refs['spider-stats'].update()
        }, 0)
      }
      this.$st.sendEv('爬虫详情', '切换标签', tab.name)
    },
    onSpiderChange (id) {
      this.$router.push(`/spiders/${id}`)
      this.$st.sendEv('爬虫详情', '切换爬虫')
    }
  },
  created () {
    // get the list of the spiders
    // this.$store.dispatch('spider/getSpiderList')

    // get spider basic info
    this.$store.dispatch('spider/getSpiderData', this.$route.params.id)
      .then(() => {
        // get spider file info
        this.$store.dispatch('file/getFileList', this.spiderForm.src)
      })

    // get spider tasks
    this.$store.dispatch('spider/getTaskList', this.$route.params.id)
  }
}
</script>

<style scoped>
  .selector {
    display: flex;
    align-items: center;
    position: absolute;
    right: 48px;
    /*float: right;*/
    z-index: 999;
    margin-top: 5px;
  }

  .selector .el-select {
    height: 30px;
    line-height: 30px;
    padding-left: 10px;
    width: 180px;
    border-radius: 0;
  }

  .selector .el-select >>> .el-input__icon,
  .selector .el-select >>> .el-input__inner {
    border-radius: 0;
    height: 30px;
    line-height: 30px;
  }

  .label {
    text-align: right;
    width: 80px;
    color: #909399;
    font-weight: 100;
  }
</style>
