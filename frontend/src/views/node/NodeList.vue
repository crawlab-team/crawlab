<template>
  <div class="app-container">
    <!--filter-->
    <div class="filter">
      <el-input prefix-icon="el-icon-search"
                :placeholder="$t('Search')"
                class="filter-search"
                v-model="filter.keyword"
                @change="onSearch">
      </el-input>
      <div class="right">
        <el-button type="success"
                   icon="el-icon-refresh"
                   class="refresh"
                   @click="onRefresh">
          {{$t('Refresh')}}
        </el-button>
      </div>
    </div>

    <!--table list-->
    <el-table :data="filteredTableData"
              class="table"
              :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
              border>
      <template v-for="col in columns">
        <el-table-column v-if="col.name === 'status'"
                         :key="col.name"
                         :label="$t(col.label)"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
          <template slot-scope="scope">
            <el-tag type="info" v-if="scope.row.status === 'offline'">{{$t('Offline')}}</el-tag>
            <el-tag type="success" v-else-if="scope.row.status === 'online'">{{$t('Online')}}</el-tag>
            <el-tag type="danger" v-else>{{$t('Unavailable')}}</el-tag>
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
      <el-table-column :label="$t('Action')" align="center" width="160">
        <template slot-scope="scope">
          <el-tooltip :content="$t('View')" placement="top">
            <el-button type="primary" icon="el-icon-search" size="mini" @click="onView(scope.row)"></el-button>
          </el-tooltip>
          <!--<el-tooltip :content="$t('Edit')" placement="top">-->
          <!--<el-button type="warning" icon="el-icon-edit" size="mini" @click="onView(scope.row)"></el-button>-->
          <!--</el-tooltip>-->
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
        :current-page.sync="pagination.pageNum"
        :page-sizes="[10, 20, 50, 100]"
        :page-size.sync="pagination.pageSize"
        layout="sizes, prev, pager, next"
        :total="nodeList.length">
      </el-pagination>
    </div>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'NodeList',
  data () {
    return {
      pagination: {
        pageNum: 0,
        pageSize: 10
      },
      isEditMode: false,
      dialogVisible: false,
      filter: {
        keyword: ''
      },
      // tableData,
      columns: [
        { name: 'name', label: 'Name', width: '220' },
        { name: 'ip', label: 'IP', width: '160' },
        { name: 'port', label: 'Port', width: '80' },
        { name: 'status', label: 'Status', width: '120', sortable: true },
        { name: 'description', label: 'Description', width: 'auto' }
      ],
      nodeFormRules: {
        name: [{ required: true, message: 'Required Field', trigger: 'change' }]
      }
    }
  },
  computed: {
    ...mapState('node', [
      'nodeList',
      'nodeForm'
    ]),
    filteredTableData () {
      return this.nodeList.filter(d => {
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
    onSearch () {
    },
    onAdd () {
      this.$store.commit('node/SET_NODE_FORM', [])
      this.isEditMode = false
      this.dialogVisible = true
    },
    onRefresh () {
      this.$store.dispatch('node/getNodeList')
      this.$st.sendEv('节点', '刷新')
    },
    onSubmit () {
      const vm = this
      const formName = 'nodeForm'
      this.$refs[formName].validate((valid) => {
        if (valid) {
          if (this.isEditMode) {
            vm.$store.dispatch('node/editNode')
          } else {
            vm.$store.dispatch('node/addNode')
          }
          vm.dialogVisible = false
        } else {
          return false
        }
      })
    },
    onCancel () {
      this.$store.commit('node/SET_NODE_FORM', {})
      this.dialogVisible = false
    },
    onDialogClose () {
      this.$store.commit('node/SET_NODE_FORM', {})
      this.dialogVisible = false
    },
    onEdit (row) {
      this.isEditMode = true
      this.$store.commit('node/SET_NODE_FORM', row)
      this.dialogVisible = true
    },
    onRemove (row) {
      this.$confirm(this.$t('Are you sure to delete this node?'), this.$t('Notification'), {
        confirmButtonText: this.$t('Confirm'),
        cancelButtonText: this.$t('Cancel'),
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('node/deleteNode', row._id)
          .then(() => {
            this.$message({
              type: 'success',
              message: 'Deleted successfully'
            })
          })
        this.$st.sendEv('节点', '删除', 'id', row._id)
      })
    },
    onView (row) {
      this.$router.push(`/nodes/${row._id}`)

      this.$st.sendEv('节点', '查看', 'id', row._id)
    },
    onPageChange () {
      this.$store.dispatch('node/getNodeList')
    }
  },
  created () {
    this.$store.dispatch('node/getNodeList')
  }
}
</script>

<style scoped lang="scss">
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
  }

  .delete-confirm {
    background-color: red;
  }

  .el-table .el-button {
    padding: 7px;
  }
</style>
