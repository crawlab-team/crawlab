<template>
  <div class="app-container">
    <!--filter-->
    <div class="filter">
      <el-input prefix-icon="el-icon-search"
                placeholder="Search"
                class="filter-search"
                v-model="filter.keyword"
                @change="onSearch">
      </el-input>
      <div class="right">
        <el-button type="success"
                   icon="el-icon-refresh"
                   class="refresh"
                   @click="onRefresh">
          Refresh
        </el-button>
      </div>
    </div>

    <!--table list-->
    <el-table :data="filteredTableData"
              class="table"
              :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
              border>
      <template v-for="col in columns">
        <el-table-column v-if="col.name === 'spider_name'"
                         :key="col.name"
                         :label="col.label"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
          <template slot-scope="scope">
            <a href="javascript:" class="a-tag" @click="onClickSpider(scope.row)">{{scope.row[col.name]}}</a>
          </template>
        </el-table-column>
        <el-table-column v-else-if="col.name === 'node_id'"
                         :key="col.name"
                         :label="col.label"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
          <template slot-scope="scope">
            <a href="javascript:" class="a-tag" @click="onClickNode(scope.row)">{{scope.row[col.name]}}</a>
          </template>
        </el-table-column>
        <el-table-column v-else-if="col.name === 'status'"
                         :key="col.name"
                         :label="col.label"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
          <template slot-scope="scope">
            <el-tag type="success" v-if="scope.row.status === 'SUCCESS'">SUCCESS</el-tag>
            <el-tag type="warning" v-else-if="scope.row.status === 'STARTED'">STARTED</el-tag>
            <el-tag type="danger" v-else-if="scope.row.status === 'FAILURE'">FAILURE</el-tag>
            <el-tag type="info" v-else>{{scope.row[col.name]}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column v-else
                         :key="col.name"
                         :property="col.name"
                         :label="col.label"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
        </el-table-column>
      </template>
      <el-table-column label="Action" align="center" width="auto">
        <template slot-scope="scope">
          <el-tooltip content="View" placement="top">
            <el-button type="primary" icon="el-icon-search" size="mini" @click="onView(scope.row)"></el-button>
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
      isEditMode: false,
      dialogVisible: false,
      filter: {
        keyword: ''
      },
      // tableData,
      columns: [
        { name: 'create_ts', label: 'Create Time', width: '150' },
        { name: 'start_ts', label: 'Start Time', width: '150' },
        { name: 'finish_ts', label: 'Finish Time', width: '150' },
        { name: 'spider_name', label: 'Spider', width: '160' },
        { name: 'node_id', label: 'Node', width: '160' },
        { name: 'status', label: 'Status', width: '160', sortable: true }
      ]
    }
  },
  computed: {
    ...mapState('task', [
      'taskList',
      'taskListTotalCount',
      'taskForm'
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
    },
    onRemove (row) {
      this.$confirm('Are you sure to delete this task?', 'Notification', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('task/deleteTask', row._id)
          .then(() => {
            this.$message({
              type: 'success',
              message: 'Deleted successfully'
            })
          })
      })
    },
    onView (row) {
      this.$router.push(`/tasks/${row._id}`)
    },
    onClickSpider (row) {
      this.$router.push(`/spiders/${row.spider_id}`)
    },
    onClickNode (row) {
      this.$router.push(`/nodes/${row.node_id}`)
    },
    onPageChange () {
      setTimeout(() => {
        this.$store.dispatch('task/getTaskList')
      }, 0)
    }
  },
  created () {
    this.$store.dispatch('task/getTaskList')
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
