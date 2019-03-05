<template>
  <div class="app-container">
    <!--add popup-->
    <el-dialog
      title="Import Spider"
      :visible.sync="dialogVisible"
      width="60%"
      :before-close="onDialogClose">
      <el-form label-width="150px"
               :model="importForm"
               ref="importForm"
               label-position="right">
        <el-form-item label="Source URL" prop="url" required>
          <el-input v-model="importForm.url" placeholder="Source URL"></el-input>
        </el-form-item>
        <el-form-item label="Source Type" prop="type" required>
          <el-select v-model="importForm.type" placeholder="Source Type">
            <el-option value="github" label="Github"></el-option>
            <el-option value="gitlab" label="Gitlab"></el-option>
            <el-option value="svn" label="SVN" disabled></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="onCancel">Cancel</el-button>
        <el-button v-loading="importLoading" type="primary" @click="onImport">Import</el-button>
      </span>
    </el-dialog>

    <!--filter-->
    <div class="filter">
      <el-input prefix-icon="el-icon-search"
                placeholder="Search"
                class="filter-search"
                v-model="filter.keyword"
                @change="onSearch">
      </el-input>
      <div class="right">
        <el-button type="primary" icon="fa fa-cloud" @click="onDeployAll">
          Deploy All
        </el-button>
        <el-button type="primary" icon="fa fa-download" @click="openImportDialog">
          Import Spiders
        </el-button>
        <el-button type="success"
                   icon="el-icon-refresh"
                   class="btn refresh"
                   @click="onRefresh">
          Refresh
        </el-button>
        <el-button type="primary"
                   v-if="false"
                   icon="el-icon-plus"
                   class="add"
                   @click="onAdd">
          Add Spider
        </el-button>
      </div>
    </div>

    <!--table list-->
    <el-table :data="filteredTableData"
              class="table"
              :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
              border>
      <template v-for="col in columns">
        <el-table-column v-if="col.name === 'type'"
                         :key="col.name"
                         :label="col.label"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.type === 'scrapy'">Scrapy</el-tag>
            <el-tag type="warning" v-else-if="scope.row.type === 'pyspider'">PySpider</el-tag>
            <el-tag type="info" v-else-if="scope.row.type === 'webmagic'">WebMagic</el-tag>
            <el-tag type="success" v-else-if="scope.row.type">{{scope.row.type}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column v-else-if="col.name === 'lang'"
                         :key="col.name"
                         :label="col.label"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
          <template slot-scope="scope">
            <el-tag type="warning" v-if="scope.row.lang === 'python'">Python</el-tag>
            <el-tag type="primary" v-else-if="scope.row.lang === 'javascript'">JavaScript</el-tag>
            <el-tag type="info" v-else-if="scope.row.lang === 'java'">Java</el-tag>
            <el-tag type="danger" v-else-if="scope.row.lang === 'go'">Go</el-tag>
            <el-tag type="success" v-else-if="scope.row.lang">{{scope.row.lang}}</el-tag>
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
      <el-table-column label="Action" align="center" width="250">
        <template slot-scope="scope">
          <el-tooltip content="View" placement="top">
            <el-button type="primary" icon="el-icon-search" size="mini" @click="onView(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip content="Edit" placement="top">
            <el-button type="warning" icon="el-icon-edit" size="mini" @click="onView(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip content="Remove" placement="top">
            <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip content="Deploy" placement="top">
            <el-button type="primary" icon="fa fa-cloud" size="mini" @click="onDeploy(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip content="Run" placement="top">
            <el-button type="success" icon="fa fa-bug" size="mini" @click="onCrawl(scope.row)"></el-button>
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
        :total="spiderList.length">
      </el-pagination>
    </div>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'SpiderList',
  data () {
    return {
      pagination: {
        pageNum: 0,
        pageSize: 10
      },
      importLoading: false,
      isEditMode: false,
      dialogVisible: false,
      filter: {
        keyword: ''
      },
      // tableData,
      columns: [
        { name: 'name', label: 'Name', width: 'auto' },
        { name: 'type', label: 'Spider Type', width: '160', sortable: true },
        { name: 'lang', label: 'Language', width: '160', sortable: true },
        { name: 'update_ts', label: 'Last Update', width: '120' }
      ],
      spiderFormRules: {
        name: [{ required: true, message: 'Required Field', trigger: 'change' }]
      }
    }
  },
  computed: {
    ...mapState('spider', [
      'importForm',
      'spiderList',
      'spiderForm'
    ]),
    filteredTableData () {
      return this.spiderList.filter(d => {
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
      console.log(value)
    },
    onAdd () {
      this.$store.commit('spider/SET_SPIDER_FORM', {})
      this.isEditMode = false
      this.dialogVisible = true
    },
    onRefresh () {
      this.$store.dispatch('spider/getSpiderList')
    },
    onSubmit () {
      const vm = this
      const formName = 'spiderForm'
      this.$refs[formName].validate((valid) => {
        if (valid) {
          if (this.isEditMode) {
            vm.$store.dispatch('spider/editSpider')
          } else {
            vm.$store.dispatch('spider/addSpider')
          }
          vm.dialogVisible = false
        } else {
          return false
        }
      })
    },
    onCancel () {
      this.$store.commit('spider/SET_SPIDER_FORM', {})
      this.dialogVisible = false
    },
    onDialogClose () {
      this.$store.commit('spider/SET_SPIDER_FORM', {})
      this.dialogVisible = false
    },
    onEdit (row) {
      this.isEditMode = true
      this.$store.commit('spider/SET_SPIDER_FORM', row)
      this.dialogVisible = true
    },
    onRemove (row) {
      this.$confirm('Are you sure to delete this spider?', 'Notification', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('spider/deleteSpider', row._id.$oid)
          .then(() => {
            this.$message({
              type: 'success',
              message: 'Deleted successfully'
            })
          })
      })
    },
    onDeploy (row) {
      this.$store.dispatch('spider/getSpiderData', row._id.$oid)
      this.$store.commit('dialogView/SET_DIALOG_VISIBLE', true)
      this.$store.commit('dialogView/SET_DIALOG_TYPE', 'spiderDeploy')
    },
    onCrawl (row) {
      this.$confirm('Are you sure to run this spider', 'Notice', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel'
      })
        .then(() => {
          this.$store.dispatch('spider/crawlSpider', { id: row._id.$oid })
            .then(() => {
              this.$message.success(`Running spider "${row._id.$oid}" has been scheduled`)
            })
        })
    },
    onView (row) {
      this.$router.push(`/spiders/${row._id.$oid}`)
    },
    onPageChange () {
      this.$store.dispatch('spider/getSpiderList')
    },
    onImport () {
      this.$refs.importForm.validate(valid => {
        if (valid) {
          this.importLoading = true
          // TODO: switch between github / gitlab / svn
          this.$store.dispatch('spider/importGithub')
            .then(response => {
              this.$message.success('Import repo successfully')
              this.$store.dispatch('spider/getSpiderList')
            })
            .catch(response => {
              this.$message.error(response.data.error)
            })
            .finally(() => {
              this.dialogVisible = false
              this.importLoading = false
            })
        }
      })
    },
    openImportDialog () {
      this.dialogVisible = true
    },
    onDeployAll () {
      this.$confirm('Are you sure to deploy all spiders to active nodes?', 'Notice', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      })
        .then(() => {
          this.$store.dispatch('spider/deployAll')
            .then(() => {
              this.$message.success('Deployed all spiders successfully')
            })
        })
    }
  },
  created () {
    this.$store.dispatch('spider/getSpiderList')
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

    .right {
      .btn {
        margin-left: 10px;
      }
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

</style>
