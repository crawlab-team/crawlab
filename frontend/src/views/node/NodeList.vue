<template>
  <div class="app-container">
    <el-dialog
      :visible.sync="isShowAddNodeInstruction"
      :title="$t('Notification')"
      width="720px"
    >
      <div
        v-html="addNodeInstructionHtml"
        class="content markdown-body"
      >
      </div>
    </el-dialog>

    <el-tabs type="border-card" v-model="activeTab">
      <el-tab-pane :label="$t('Node List')">
        <!--filter-->
        <div class="filter-wrapper">
          <el-button
            size="small"
            type="success"
            icon="el-icon-plus"
            @click="onAddNode"
          >
            {{$t('Add Node')}}
          </el-button>
        </div>
        <!--./filter-->
        <!--table list-->
        <el-table :data="filteredTableData"
                  class="table"
                  :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
                  border
                  @row-click="onRowClick"
                  @expand-change="onRowExpand">
          <el-table-column type="expand">
            <template slot-scope="scope">
              <el-form class="node-detail" :model="scope.row" label-width="120px" inline>
                <el-form-item :label="$t('OS')">
                  <span>{{scope.row.systemInfo ? getOSName(scope.row.systemInfo.os) : ''}}</span>
                </el-form-item>
                <el-form-item :label="$t('ARCH')">
                  <span>{{scope.row.systemInfo ? scope.row.systemInfo.arch : ''}}</span>
                </el-form-item>
                <el-form-item :label="$t('Number of CPU')">
                  <span>{{scope.row.systemInfo ? scope.row.systemInfo.num_cpu : ''}}</span>
                </el-form-item>
              </el-form>
              <el-form class="node-detail executable" :model="scope.row" label-width="120px">
                <el-form-item :label="$t('Executables')">
                  <ul v-if="true" class="executable-list">
                    <li
                      v-for="(ex, index) in getExecutables(scope.row)"
                      :key="index"
                      class="executable-item"
                    >
                      <el-tooltip :content="ex.path">
                        <div>
                          <template v-if="ex.file_name.match(/^python/)">
                            <font-awesome-icon :icon="['fab','python']"/>
                          </template>
                          <template v-else-if="ex.file_name.match(/^java/)">
                            <font-awesome-icon :icon="['fab','java']"/>
                          </template>
                          <template v-else-if="ex.file_name.match(/^bash$|^sh$/)">
                            <font-awesome-icon :icon="['fab','linux']"/>
                          </template>
                          <template v-else-if="ex.file_name.match(/^node/)">
                            <font-awesome-icon :icon="['fab','node-js']"/>
                          </template>
                          <template v-else-if="ex.file_name.match(/^php/)">
                            <font-awesome-icon :icon="['fab','php']"/>
                          </template>
                          <template v-else>
                            <font-awesome-icon :icon="['fas', 'terminal']"/>
                          </template>
                          <span class="executable-label">{{ex.display_name}}</span>
                        </div>
                      </el-tooltip>
                    </li>
                  </ul>
                </el-form-item>
              </el-form>
            </template>
          </el-table-column>
          <template v-for="col in columns">
            <el-table-column v-if="col.name === 'status'"
                             :key="col.name"
                             :label="$t(col.label)"
                             :sortable="col.sortable"
                             :width="col.width">
              <template slot-scope="scope">
                <el-tag type="info" v-if="scope.row.status === 'offline'">{{$t('Offline')}}</el-tag>
                <el-tag type="success" v-else-if="scope.row.status === 'online'">{{$t('Online')}}</el-tag>
                <el-tag type="danger" v-else>{{$t('Unavailable')}}</el-tag>
              </template>
            </el-table-column>
            <el-table-column v-else-if="col.name === 'type'"
                             :key="col.name"
                             :label="$t(col.label)"
                             :sortable="col.sortable"
                             :width="col.width">
              <template slot-scope="scope">
                <el-tag type="primary" v-if="scope.row.is_master">{{$t('Master')}}</el-tag>
                <el-tag type="warning" v-else>{{$t('Worker')}}</el-tag>
              </template>
            </el-table-column>
            <el-table-column v-else
                             :key="col.name"
                             :property="col.name"
                             :label="$t(col.label)"
                             :sortable="col.sortable"
                             :align="col.align || 'left'"
                             :width="col.width">
            </el-table-column>
          </template>
          <el-table-column :label="$t('Action')" align="left" width="160" fixed="right">
            <template slot-scope="scope">
              <el-tooltip :content="$t('View')" placement="top">
                <el-button type="primary" icon="el-icon-search" size="mini" @click="onView(scope.row)"></el-button>
              </el-tooltip>
              <el-tooltip :content="$t('Remove')" placement="top">
                <el-button v-if="scope.row.status !== 'online'" type="danger" icon="el-icon-delete" size="mini"
                           @click="onRemove(scope.row)"></el-button>
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
        <!--./table list-->
      </el-tab-pane>
      <el-tab-pane :label="$t('Network')">
        <node-network :active-tab="activeTab"/>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import showdown from 'showdown'
import {
  mapState
} from 'vuex'
import 'github-markdown-css/github-markdown.css'
import NodeNetwork from '../../components/Node/NodeNetwork'

export default {
  name: 'NodeList',
  components: { NodeNetwork },
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
        { name: 'type', label: 'nodeList.type', width: '120' },
        // { name: 'port', label: 'Port', width: '80' },
        { name: 'status', label: 'Status', width: '120' },
        { name: 'description', label: 'Description', width: 'auto' }
      ],
      nodeFormRules: {
        name: [{ required: true, message: 'Required Field', trigger: 'change' }]
      },
      activeTab: undefined,
      isButtonClicked: false,
      isShowAddNodeInstruction: false,
      converter: new showdown.Converter(),
      addNodeInstructionMarkdown: 'addNodeInstruction'
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
    },
    addNodeInstructionHtml () {
      if (!this.converter) return
      return this.converter.makeHtml(this.$t(this.addNodeInstructionMarkdown))
    }
  },
  methods: {
    onSearch () {
    },
    onAddNode () {
      this.isShowAddNodeInstruction = true
    },
    // onAdd () {
    //   this.$store.commit('node/SET_NODE_FORM', [])
    //   this.isEditMode = false
    //   this.dialogVisible = true
    // },
    onRefresh () {
      this.$store.dispatch('node/getNodeList')
      this.$st.sendEv('节点列表', '刷新')
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
      this.isButtonClicked = true
      setTimeout(() => {
        this.isButtonClicked = false
      }, 100)

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
        this.$st.sendEv('节点列表', '删除节点')
      })
    },
    onView (row) {
      this.isButtonClicked = true
      setTimeout(() => {
        this.isButtonClicked = false
      }, 100)

      this.$router.push(`/nodes/${row._id}`)

      this.$st.sendEv('节点列表', '查看节点')
    },
    onPageChange () {
      this.$store.dispatch('node/getNodeList')
    },
    onRowExpand (row) {
      this.$store.dispatch('node/getNodeSystemInfo', row._id)
    },
    onRowClick (row) {
      if (this.isButtonClicked) return
      this.onView(row)
    },
    getExecutables (row) {
      if (!row.systemInfo || !row.systemInfo.executables) return []
      return row.systemInfo.executables
    },
    getOSName (os) {
      if (os === 'linux') {
        return 'Linux'
      } else if (os === 'windows') {
        return 'Windows'
      } else if (os === 'darwin') {
        return 'Mac OS (darwin)'
      } else {
        return os
      }
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

  .filter-wrapper {
    text-align: right;
  }

  .content {
    word-break: break-word;
  }
</style>
<style>
  .node-detail .el-form-item {
    height: 25px;
    margin-bottom: 0;
  }

  .node-detail .el-form-item .el-form-item__label {
    line-height: 25px;
    height: 25px;
    font-size: 12px;
  }

  .node-detail .el-form-item .el-form-item__content {
    line-height: 25px;
    height: 25px;
    font-size: 12px;
  }

  .node-detail.executable .el-form-item,
  .node-detail.executable .el-form-item .el-form-item__label,
  .node-detail.executable .el-form-item .el-form-item__content {
    line-height: 25px;
    height: auto;
  }

  .node-detail .executable-list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-wrap: wrap;
  }

  .node-detail .executable-list .executable-item {
    margin-right: 10px;
    cursor: pointer;
  }

  .node-detail .executable-list .executable-item:hover {
    text-decoration: underline;
  }

  .node-detail .executable-list .executable-label {
    margin-left: 5px;
  }

  .table.el-table tr {
    cursor: pointer;
  }
</style>
