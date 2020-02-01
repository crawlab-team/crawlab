<template>
  <div class="app-container">
    <!--tour-->
    <v-tour
      name="task-list"
      :steps="tourSteps"
      :callbacks="tourCallbacks"
      :options="$utils.tour.getOptions(true)"
    />
    <!--./tour-->

    <el-card style="border-radius: 0">
      <!--filter-->
      <div class="filter">
        <div class="left">
          <el-form class="filter-form" :model="filter" label-width="100px" label-position="right" inline>
            <el-form-item prop="node_id" :label="$t('Node')">
              <el-select v-model="filter.node_id" size="small" :placeholder="$t('Node')" @change="onFilterChange">
                <el-option value="" :label="$t('All')"/>
                <el-option v-for="node in nodeList" :key="node._id" :value="node._id" :label="node.name"/>
              </el-select>
            </el-form-item>
            <el-form-item prop="spider_id" :label="$t('Spider')">
              <el-select v-model="filter.spider_id" size="small" :placeholder="$t('Spider')" @change="onFilterChange">
                <el-option value="" :label="$t('All')"/>
                <el-option v-for="spider in spiderList" :key="spider._id" :value="spider._id" :label="spider.name"/>
              </el-select>
            </el-form-item>
            <el-form-item prop="status" :label="$t('Status')">
              <el-select v-model="filter.status" size="small" :placeholder="$t('Status')" @change="onFilterChange">
                <el-option value="" :label="$t('All')"></el-option>
                <el-option value="finished" :label="$t('Finished')"></el-option>
                <el-option value="running" :label="$t('Running')"></el-option>
                <el-option value="error" :label="$t('Error')"></el-option>
                <el-option value="cancelled" :label="$t('Cancelled')"></el-option>
              </el-select>
            </el-form-item>
          </el-form>
        </div>
        <div class="right">
          <el-button class="btn-delete" @click="onRemoveMultipleTask" size="small" type="danger">
            删除任务
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--table list-->
      <el-table :data="filteredTableData"
                ref="table"
                class="table"
                :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
                border
                row-key="_id"
                @row-click="onRowClick"
                @selection-change="onSelectionChange">
      >
        <el-table-column type="selection" width="55" reserve-selection/>
        <template v-for="col in columns">
          <el-table-column v-if="col.name === 'spider_name'"
                           :key="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
          >
            <template slot-scope="scope">
              <a href="javascript:" class="a-tag" @click="onClickSpider(scope.row)">{{scope.row[col.name]}}</a>
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name.match(/_ts$/)"
                           :key="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              {{getTime(scope.row[col.name])}}
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'wait_duration'"
                           :key="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              {{getWaitDuration(scope.row)}}
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'runtime_duration'"
                           :key="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              {{getRuntimeDuration(scope.row)}}
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'total_duration'"
                           :key="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              {{getTotalDuration(scope.row)}}
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'node_name'"
                           :key="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              <a href="javascript:" class="a-tag" @click="onClickNode(scope.row)">{{scope.row[col.name]}}</a>
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'status'"
                           :key="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              <status-tag :status="scope.row[col.name]"/>
            </template>
          </el-table-column>
          <el-table-column v-else
                           :key="col.name"
                           :property="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
          </el-table-column>
        </template>
        <el-table-column :label="$t('Action')" align="left" fixed="right" width="120px">
          <template slot-scope="scope">
            <el-tooltip :content="$t('View')" placement="top">
              <el-button type="primary" icon="el-icon-search" size="mini" @click="onView(scope.row)"></el-button>
            </el-tooltip>
            <el-tooltip :content="$t('Remove')" placement="top">
              <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row, $event)"></el-button>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          @current-change="onPageChange"
          @size-change="onPageChange"
          :current-page.sync="pageNum"
          :page-sizes="[10, 20, 50, 100]"
          :page-size.sync="pageSize"
          layout="sizes, prev, pager, next"
          :total="taskListTotalCount">
        </el-pagination>
      </div>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import dayjs from 'dayjs'
import StatusTag from '../../components/Status/StatusTag'

export default {
  name: 'TaskList',
  components: { StatusTag },
  data () {
    return {
      // setInterval handle
      handle: undefined,

      // determine if is edit mode
      isEditMode: false,

      // dialog visibility
      dialogVisible: false,

      // table columns
      columns: [
        { name: 'node_name', label: 'Node', width: '120' },
        { name: 'spider_name', label: 'Spider', width: '120' },
        { name: 'status', label: 'Status', width: '120' },
        { name: 'param', label: 'Parameters', width: '120' },
        // { name: 'create_ts', label: 'Create Time', width: '100' },
        { name: 'start_ts', label: 'Start Time', width: '100' },
        { name: 'finish_ts', label: 'Finish Time', width: '100' },
        { name: 'wait_duration', label: 'Wait Duration (sec)', align: 'right' },
        { name: 'runtime_duration', label: 'Runtime Duration (sec)', align: 'right' },
        { name: 'total_duration', label: 'Total Duration (sec)', width: '80', align: 'right' },
        { name: 'result_count', label: 'Results Count', width: '80' }
        // { name: 'avg_num_results', label: 'Average Results Count per Second', width: '80' }
      ],

      multipleSelection: [],

      // tutorial
      tourSteps: [
        {
          target: '.filter-form',
          content: this.$t('You can filter tasks from this area.')
        },
        {
          target: '.table',
          content: this.$t('This is a list of spider tasks executed sorted in a descending order.')
        },
        {
          target: '.table .el-table__body-wrapper tr:nth-child(1)',
          content: this.$t('Click the row to or the view button to view the task detail.')
        },
        {
          target: '.table tr td:nth-child(1)',
          content: this.$t('Tick and select the tasks you would like to delete in batches.'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '.btn-delete',
          content: this.$t('Click this button to delete selected tasks.'),
          params: {
            placement: 'left'
          }
        }
      ],
      tourCallbacks: {
        onStop: () => {
          this.$utils.tour.finishTour('task-list')
        },
        onPreviousStep: (currentStep) => {
        },
        onNextStep: (currentStep) => {
        }
      }
    }
  },
  computed: {
    ...mapState('task', [
      'filter',
      'taskList',
      'taskListTotalCount',
      'taskForm'
    ]),
    ...mapState('spider', [
      'spiderList'
    ]),
    ...mapState('node', [
      'nodeList'
    ]),
    pageNum: {
      get () {
        return this.$store.state.task.pageNum
      },
      set (value) {
        this.$store.commit('task/SET_PAGE_NUM', value)
      }
    },
    pageSize: {
      get () {
        return this.$store.state.task.pageSize
      },
      set (value) {
        this.$store.commit('task/SET_PAGE_SIZE', value)
      }
    },
    filteredTableData () {
      return this.taskList
        .map(d => d)
        .sort((a, b) => a.create_ts < b.create_ts ? 1 : -1)
        .filter(d => {
          // keyword
          if (!this.filter.keyword) return true
          for (let i = 0; i < this.columns.length; i++) {
            const colName = this.columns[i].name
            if (d[colName] && d[colName].toLowerCase().indexOf(this.filter.keyword.toLowerCase()) > -1) {
              return true
            }
          }
          return false
        })
    }
  },
  methods: {
    onSearch (value) {
    },
    onRefresh () {
      this.$store.dispatch('task/getTaskList')
      this.$st.sendEv('任务列表', '搜索')
    },
    onRemoveMultipleTask () {
      if (this.multipleSelection.length === 0) {
        this.$message({
          type: 'error',
          message: '请选择要删除的任务'
        })
        return
      }
      this.$confirm('确定删除任务', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        let ids = this.multipleSelection.map(item => item._id)
        this.$store.dispatch('task/deleteTaskMultiple', ids).then((resp) => {
          if (resp.data.status === 'ok') {
            this.$message({
              type: 'success',
              message: '删除任务成功'
            })
            this.$store.dispatch('task/getTaskList')
            this.$refs['table'].clearSelection()
            return
          }
          this.$message({
            type: 'error',
            message: resp.data.error
          })
        })
      }).catch(() => {})
    },
    onRemove (row, ev) {
      ev.stopPropagation()
      this.$confirm(this.$t('Are you sure to delete this task?'), this.$t('Notification'), {
        confirmButtonText: this.$t('Confirm'),
        cancelButtonText: this.$t('Cancel'),
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('task/deleteTask', row._id)
          .then(() => {
            this.$message({
              type: 'success',
              message: 'Deleted successfully'
            })
          })
        this.$st.sendEv('任务列表', '删除任务')
      })
    },
    onView (row) {
      this.$router.push(`/tasks/${row._id}`)
      this.$st.sendEv('任务列表', '查看任务')
    },
    onClickSpider (row) {
      this.$router.push(`/spiders/${row.spider_id}`)
      this.$st.sendEv('任务列表', '点击爬虫详情')
    },
    onClickNode (row) {
      this.$router.push(`/nodes/${row.node_id}`)
      this.$st.sendEv('任务列表', '点击节点详情')
    },
    onPageChange () {
      setTimeout(() => {
        this.$store.dispatch('task/getTaskList')
      }, 0)
    },
    getTime (str) {
      if (str.match('^0001')) return 'NA'
      return dayjs(str).format('YYYY-MM-DD HH:mm:ss')
    },
    getWaitDuration (row) {
      if (row.start_ts.match('^0001')) return 'NA'
      return dayjs(row.start_ts).diff(row.create_ts, 'second')
    },
    getRuntimeDuration (row) {
      if (row.finish_ts.match('^0001')) return 'NA'
      return dayjs(row.finish_ts).diff(row.start_ts, 'second')
    },
    getTotalDuration (row) {
      if (row.finish_ts.match('^0001')) return 'NA'
      return dayjs(row.finish_ts).diff(row.create_ts, 'second')
    },
    onRowClick (row, event, column) {
      if (column.label !== this.$t('Action')) {
        this.onView(row)
      }
    },
    onSelectionChange (val) {
      this.multipleSelection = val
    },
    onFilterChange () {
      this.$store.dispatch('task/getTaskList')
      this.$st.sendEv('任务列表', '筛选任务')
    }
  },
  created () {
    this.$store.dispatch('task/getTaskList')
    this.$store.dispatch('spider/getSpiderList')
    this.$store.dispatch('node/getNodeList')
  },
  mounted () {
    this.handle = setInterval(() => {
      this.$store.dispatch('task/getTaskList')
    }, 5000)

    if (!this.$utils.tour.isFinishedTour('task-list')) {
      this.$tours['task-list'].start()
    }
  },
  destroyed () {
    clearInterval(this.handle)
  }
}
</script>

<style scoped lang="scss">
  .el-dialog {
    .el-select {
      width: 100%;
    }
  }

  .filter {
    display: flex;
    justify-content: space-between;

    .left {
      .filter-select {
        width: 180px;
        margin-right: 10px;
      }
    }

    .filter-search {
      width: 240px;
    }

    .add {
    }
  }

  .table {
    margin-top: 8px;
    border-radius: 5px;

    .el-button {
      padding: 7px;
    }
  }

  .delete-confirm {
    background-color: red;
  }

  .el-table .a-tag {
    text-decoration: underline;
  }

  .pagination {
    margin-top: 10px;
    text-align: right;
  }
</style>

<style scoped>
  .el-table >>> tr {
    cursor: pointer;
  }
</style>
