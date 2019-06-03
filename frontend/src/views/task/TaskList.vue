<template>
  <div class="app-container">
    <!--filter-->
    <div class="filter">
      <div class="left">
        <el-select class="filter-select"
                   v-model="filter.node_id"
                   :placeholder="$t('Node')"
                   filterable
                   clearable
                   @change="onSelectNode">
          <el-option v-for="op in nodeList" :key="op._id" :value="op._id" :label="op.name"></el-option>
        </el-select>
        <el-select class="filter-select"
                   v-model="filter.spider_id"
                   :placeholder="$t('Spider')"
                   filterable
                   clearable
                   @change="onSelectSpider">
          <el-option v-for="op in spiderList" :key="op._id" :value="op._id" :label="op.name"></el-option>
        </el-select>
        <el-button type="success"
                   icon="el-icon-search"
                   class="refresh"
                   @click="onRefresh">
          {{$t('Search')}}
        </el-button>
      </div>
      <!--<div class="right">-->
      <!--</div>-->
    </div>

    <!--table list-->
    <el-table :data="filteredTableData"
              class="table"
              :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
              border>
      <template v-for="col in columns">
        <el-table-column v-if="col.name === 'spider_name'"
                         :key="col.name"
                         :label="$t(col.label)"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
          <template slot-scope="scope">
            <a href="javascript:" class="a-tag" @click="onClickSpider(scope.row)">{{scope.row[col.name]}}</a>
          </template>
        </el-table-column>
        <el-table-column v-else-if="col.name === 'node_id'"
                         :key="col.name"
                         :label="$t(col.label)"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
          <template slot-scope="scope">
            <a href="javascript:" class="a-tag" @click="onClickNode(scope.row)">{{scope.row[col.name]}}</a>
          </template>
        </el-table-column>
        <el-table-column v-else-if="col.name === 'status'"
                         :key="col.name"
                         :label="$t(col.label)"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
          <template slot-scope="scope">
            <el-tag type="success" v-if="scope.row.status === 'SUCCESS'">{{$t('SUCCESS')}}</el-tag>
            <el-tag type="warning" v-else-if="scope.row.status === 'STARTED'">{{$t('STARTED')}}</el-tag>
            <el-tag type="danger" v-else-if="scope.row.status === 'FAILURE'">{{$t('FAILURE')}}</el-tag>
            <el-tag type="info" v-else>{{$t(scope.row[col.name])}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column v-else
                         :key="col.name"
                         :property="col.name"
                         :label="$t(col.label)"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
        </el-table-column>
      </template>
      <el-table-column :label="$t('Action')" align="left" width="auto" fixed="right">
        <template slot-scope="scope">
          <el-tooltip :content="$t('View')" placement="top">
            <el-button type="primary" icon="el-icon-search" size="mini" @click="onView(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip :content="$t('Remove')" placement="top">
            <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)"></el-button>
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
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'TaskList',
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
        { name: 'create_ts', label: 'Create Time', width: '100' },
        { name: 'start_ts', label: 'Start Time', width: '100' },
        { name: 'finish_ts', label: 'Finish Time', width: '100' },
        { name: 'duration', label: 'Duration (sec)', width: '80' },
        { name: 'spider_name', label: 'Spider', width: '120' },
        { name: 'node_id', label: 'Node', width: '160' },
        { name: 'num_results', label: 'Results Count', width: '80' },
        { name: 'avg_num_results', label: 'Average Results Count per Second', width: '80' },
        { name: 'status', label: 'Status', width: '80' }
      ]
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
      // .filter((d, index) => {
      //   // pagination
      //   const pageNum = this.pageNum
      //   const pageSize = this.pageSize
      //   return (pageSize * (pageNum - 1) <= index) && (index < pageSize * pageNum)
      // })
    }
  },
  methods: {
    onSearch (value) {
      console.log(value)
    },
    onRefresh () {
      this.$store.dispatch('task/getTaskList')
      this.$st.sendEv('任务', '搜索')
    },
    onSelectNode () {
      this.$st.sendEv('任务', '选择节点')
    },
    onSelectSpider () {
      this.$st.sendEv('任务', '选择爬虫')
    },
    onRemove (row) {
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
        this.$st.sendEv('任务', '删除', 'id', row._id)
      })
    },
    onView (row) {
      this.$router.push(`/tasks/${row._id}`)
      this.$st.sendEv('任务', '搜索', 'id', row._id)
    },
    onClickSpider (row) {
      this.$router.push(`/spiders/${row.spider_id}`)
      this.$st.sendEv('任务', '点击爬虫详情', 'id', row.spider_id)
    },
    onClickNode (row) {
      this.$router.push(`/nodes/${row.node_id}`)
      this.$st.sendEv('任务', '点击节点详情', 'id', row.node_id)
    },
    onPageChange () {
      setTimeout(() => {
        this.$store.dispatch('task/getTaskList')
      }, 0)
    }
  },
  created () {
    this.$store.dispatch('task/getTaskList')
    this.$store.dispatch('spider/getSpiderList')
    this.$store.dispatch('node/getNodeList')
  },
  mounted () {
    // request task list every 5 seconds
    this.handle = setInterval(() => {
      this.$store.dispatch('task/getTaskList')
    }, 5000)
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
    margin-top: 20px;
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
