<template>
  <div class="app-container">
    <!--tabs-->
    <el-tabs v-model="activeTabName" @tab-click="onTabClick" type="card">
      <el-tab-pane :label="$t('Overview')" name="overview">
        <task-overview/>
      </el-tab-pane>
      <el-tab-pane :label="$t('Log')" name="log">
        <el-card>
          <log-view :data="taskLog"/>
        </el-card>
      </el-tab-pane>
      <el-tab-pane :label="$t('Results')" name="results">
        <div class="button-group">
          <el-button type="primary" icon="el-icon-download" @click="downloadCSV">
            {{$t('Download CSV')}}
          </el-button>
        </div>
        <general-table-view :data="taskResultsData"
                            :columns="taskResultsColumns"
                            :page-num="resultsPageNum"
                            :page-size="resultsPageSize"
                            :total="taskResultsTotalCount"
                            @page-change="onResultsPageChange"/>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import {
  mapState,
  mapGetters
} from 'vuex'
import TaskOverview from '../../components/Overview/TaskOverview'
import GeneralTableView from '../../components/TableView/GeneralTableView'
import LogView from '../../components/ScrollView/LogView'

export default {
  name: 'TaskDetail',
  components: {
    LogView,
    GeneralTableView,
    TaskOverview
  },
  data () {
    return {
      activeTabName: 'overview',
      handle: undefined
    }
  },
  computed: {
    ...mapState('task', [
      'taskLog',
      'taskResultsData',
      'taskResultsTotalCount'
    ]),
    ...mapGetters('task', [
      'taskResultsColumns'
    ]),
    ...mapState('file', [
      'currentPath'
    ]),
    ...mapState('deploy', [
      'deployList'
    ]),
    resultsPageNum: {
      get () {
        return this.$store.state.task.resultsPageNum
      },
      set (value) {
        this.$store.commit('task/SET_RESULTS_PAGE_NUM', value)
      }
    },
    resultsPageSize: {
      get () {
        return this.$store.state.task.resultsPageSize
      },
      set (value) {
        this.$store.commit('task/SET_RESULTS_PAGE_SIZE', value)
      }
    }
  },
  methods: {
    onTabClick (tab) {
      this.$st.sendEv('任务详情', '切换标签', tab.name)
    },
    onSpiderChange (id) {
      this.$router.push(`/spiders/${id}`)
    },
    onResultsPageChange (payload) {
      const { pageNum, pageSize } = payload
      this.resultsPageNum = pageNum
      this.resultsPageSize = pageSize
      this.$store.dispatch('task/getTaskResults', this.$route.params.id)
    },
    downloadCSV () {
      this.$store.dispatch('task/getTaskResultExcel', this.$route.params.id)
      this.$st.sendEv('任务详情-结果', '下载CSV')
    }
  },
  created () {
    this.$store.dispatch('task/getTaskData', this.$route.params.id)
    this.$store.dispatch('task/getTaskLog', this.$route.params.id)
    this.$store.dispatch('task/getTaskResults', this.$route.params.id)

    if (['running'].includes(this.taskForm.status)) {
      this.handle = setInterval(() => {
        this.$store.dispatch('task/getTaskLog', this.$route.params.id)
      }, 5000)
    }
  },
  destroyed () {
    clearInterval(this.handle)
  }
}
</script>

<style scoped>

  .selector {
    display: flex;
    align-items: center;
    position: absolute;
    right: 20px;
    /*float: right;*/
    z-index: 999;
    margin-top: -7px;
  }

  .selector .el-select {
    padding-left: 10px;
  }

  .log-view {
    margin: 20px;
    height: 640px;
  }

  .log-view pre {
    height: 100%;
    overflow-x: auto;
    overflow-y: auto;
  }

  .button-group {
    margin-bottom: 10px;
    text-align: right;
  }
</style>
