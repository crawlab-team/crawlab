<template>
  <div class="app-container">
    <!--import popup-->
    <el-dialog
      :title="$t('Import Spider')"
      :visible.sync="dialogVisible"
      width="60%"
      :before-close="onDialogClose">
      <el-form label-width="150px"
               :model="importForm"
               ref="importForm"
               label-position="right">
        <el-form-item :label="$t('Source URL')" prop="url" required>
          <el-input v-model="importForm.url" :placeholder="$t('Source URL')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Source Type')" prop="type" required>
          <el-select v-model="importForm.type" placeholder="Source Type">
            <el-option value="github" label="Github"></el-option>
            <el-option value="gitlab" label="Gitlab"></el-option>
            <el-option value="svn" label="SVN" disabled></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="onCancel">{{$t('Cancel')}}</el-button>
        <el-button v-loading="importLoading" type="primary" @click="onImport">{{$t('Import')}}</el-button>
      </span>
    </el-dialog>
    <!--./import popup-->

    <!--add dialog-->
    <el-dialog :title="$t('Add Spider')"
               width="40%"
               :visible.sync="addDialogVisible"
               :before-close="onAddDialogClose">
      <div class="add-spider-wrapper">
        <div @click="onAddConfigurable">
          <el-card shadow="hover" class="add-spider-item success">
            {{$t('Configurable Spider')}}
          </el-card>
        </div>
        <div @click="onAddCustomized">
          <el-card shadow="hover" class="add-spider-item primary">
            {{$t('Customized Spider')}}
          </el-card>
        </div>
      </div>
    </el-dialog>
    <!--./add dialog-->

    <!--configurable spider dialog-->
    <el-dialog :title="$t('Add Configurable Spider')"
               width="40%"
               :visible.sync="addConfigurableDialogVisible"
               :before-close="onAddConfigurableDialogClose">
      <el-form :model="spiderForm" ref="addConfigurableForm" inline-message>
        <el-form-item :label="$t('Spider Name')" label-width="120px" prop="name" required>
          <el-input :placeholder="$t('Spider Name')" v-model="spiderForm.name"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Results Collection')" label-width="120px" name="col">
          <el-input :placeholder="$t('Results Collection')" v-model="spiderForm.col"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Site')" label-width="120px" name="site">
          <el-autocomplete v-model="spiderForm.site"
                           :placeholder="$t('Site')"
                           :fetch-suggestions="fetchSiteSuggestions"
                           @select="onAddConfigurableSiteSelect">
          </el-autocomplete>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="addConfigurableDialogVisible = false">{{$t('Cancel')}}</el-button>
        <el-button v-loading="addConfigurableLoading" type="primary"
                   @click="onAddConfigurableSpider">{{$t('Add')}}</el-button>
      </span>
    </el-dialog>
    <!--./configurable spider dialog-->

    <!--customized spider dialog-->
    <el-dialog :title="$t('Add Customized Spider')"
               width="40%"
               :visible.sync="addCustomizedDialogVisible"
               :before-close="onAddCustomizedDialogClose">
      <el-form :model="spiderForm" ref="addConfigurableForm" inline-message>
        <el-form-item :label="$t('Upload Zip File')" label-width="120px" name="site">
          <el-upload
            :action="$request.baseUrl + '/spiders/manage/upload'"
            :on-success="onUploadSuccess"
            :file-list="fileList">
            <el-button size="small" type="primary">{{$t('Upload')}}</el-button>
          </el-upload>
        </el-form-item>
      </el-form>
    </el-dialog>
    <!--./customized spider dialog-->

    <!--filter-->
    <div class="filter">
      <!--<el-input prefix-icon="el-icon-search"-->
      <!--:placeholder="$t('Search')"-->
      <!--class="filter-search"-->
      <!--v-model="filter.keyword"-->
      <!--@change="onSearch">-->
      <!--</el-input>-->
      <div class="left">
        <el-autocomplete v-model="filterSite"
                         :placeholder="$t('Site')"
                         clearable
                         :fetch-suggestions="fetchSiteSuggestions"
                         @select="onSiteSelect">
        </el-autocomplete>
      </div>
      <div class="right">
        <el-button type="primary" icon="fa fa-cloud" @click="onDeployAll">
          {{$t('Deploy All')}}
        </el-button>
        <el-button type="primary" icon="fa fa-download" @click="openImportDialog">
          {{$t('Import Spiders')}}
        </el-button>
        <el-button type="success"
                   icon="el-icon-plus"
                   class="btn add"
                   @click="onAdd">
          {{$t('Add Spider')}}
        </el-button>
        <el-button type="success"
                   icon="el-icon-refresh"
                   class="btn refresh"
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
        <el-table-column v-if="col.name === 'type'"
                         :key="col.name"
                         :label="$t(col.label)"
                         align="center"
                         :width="col.width">
          <template slot-scope="scope">
            <el-tag type="success" v-if="scope.row.type === 'configurable'">{{$t('Configurable')}}</el-tag>
            <el-tag type="primary" v-else-if="scope.row.type === 'customized'">{{$t('Customized')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column v-else-if="col.name === 'lang'"
                         :key="col.name"
                         :label="$t(col.label)"
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
        <el-table-column v-else-if="col.name === 'last_5_errors'"
                         :key="col.name"
                         :label="$t(col.label)"
                         :width="col.width"
                         align="center">
          <template slot-scope="scope">
            <div :style="{color:scope.row[col.name]>0?'red':''}">
              {{scope.row[col.name]}}
            </div>
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
          <!--<el-tooltip :content="$t('Edit')" placement="top">-->
          <!--<el-button type="warning" icon="el-icon-edit" size="mini" @click="onView(scope.row)"></el-button>-->
          <!--</el-tooltip>-->
          <el-tooltip :content="$t('Remove')" placement="top">
            <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip v-if="scope.row.type === 'customized'" :content="$t('Deploy')" placement="top">
            <el-button type="primary" icon="fa fa-cloud" size="mini" @click="onDeploy(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip v-if="isShowRun(scope.row)" :content="$t('Run')" placement="top">
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
      addConfigurableLoading: false,
      isEditMode: false,
      dialogVisible: false,
      addDialogVisible: false,
      addConfigurableDialogVisible: false,
      addCustomizedDialogVisible: false,
      filter: {
        keyword: ''
      },
      // tableData,
      columns: [
        { name: 'name', label: 'Name', width: '180', align: 'left' },
        { name: 'site_name', label: 'Site', width: '140', align: 'left' },
        { name: 'type', label: 'Spider Type', width: '120' },
        { name: 'lang', label: 'Language', width: '120', sortable: true },
        { name: 'task_ts', label: 'Last Run', width: '160' },
        { name: 'last_7d_tasks', label: 'Last 7-Day Tasks', width: '80' },
        { name: 'last_5_errors', label: 'Last 5-Run Errors', width: '80' }
      ],
      spiderFormRules: {
        name: [{ required: true, message: 'Required Field', trigger: 'change' }]
      },
      fileList: []
    }
  },
  computed: {
    ...mapState('spider', [
      'importForm',
      'spiderList',
      'spiderForm'
    ]),
    filteredTableData () {
      return this.spiderList
        .filter(d => {
          if (this.filterSite) {
            return d.site === this.filterSite
          }
          return true
        })
        .filter((d, index) => {
          return (this.pagination.pageSize * (this.pagination.pageNum - 1)) <= index && (index < this.pagination.pageSize * this.pagination.pageNum)
        })
      // .filter(d => {
      //   if (!this.filter.keyword) return true
      //   for (let i = 0; i < this.columns.length; i++) {
      //     const colName = this.columns[i].name
      //     if (d[colName] && d[colName].toLowerCase().indexOf(this.filter.keyword.toLowerCase()) > -1) {
      //       return true
      //     }
      //   }
      //   return false
      // })
    },
    filterSite: {
      get () {
        return this.$store.state.spider.filterSite
      },
      set (value) {
        this.$store.commit('spider/SET_FILTER_SITE', value)
      }
    }
  },
  methods: {
    onSearch (value) {
      console.log(value)
    },
    onAdd () {
      this.addDialogVisible = true
    },
    onAddConfigurable () {
      this.$store.commit('spider/SET_SPIDER_FORM', {})
      this.addDialogVisible = false
      this.addConfigurableDialogVisible = true
      this.$st.sendEv('爬虫', '添加爬虫-可配置爬虫')
    },
    onAddCustomized () {
      this.addDialogVisible = false
      this.addCustomizedDialogVisible = true
      this.$st.sendEv('爬虫', '添加爬虫-自定义爬虫')
    },
    onRefresh () {
      this.$store.dispatch('spider/getSpiderList')
      this.$st.sendEv('爬虫', '刷新')
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
    onAddCancel () {
      this.addDialogVisible = false
    },
    onDialogClose () {
      this.$store.commit('spider/SET_SPIDER_FORM', {})
      this.dialogVisible = false
    },
    onAddDialogClose () {
      this.addDialogVisible = false
    },
    onAddCustomizedDialogClose () {
      this.addCustomizedDialogVisible = false
    },
    onAddConfigurableDialogClose () {
      this.addConfigurableDialogVisible = false
    },
    onEdit (row) {
      this.isEditMode = true
      this.$store.commit('spider/SET_SPIDER_FORM', row)
      this.dialogVisible = true
    },
    onRemove (row) {
      this.$confirm(this.$t('Are you sure to delete this spider?'), this.$t('Notification'), {
        confirmButtonText: this.$t('Confirm'),
        cancelButtonText: this.$t('Cancel'),
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('spider/deleteSpider', row._id)
          .then(() => {
            this.$message({
              type: 'success',
              message: 'Deleted successfully'
            })
          })
        this.$st.sendEv('爬虫', '删除')
      })
    },
    onDeploy (row) {
      this.$confirm(this.$t('Are you sure to deploy this spider?'), this.$t('Notification'), {
        confirmButtonText: this.$t('Confirm'),
        cancelButtonText: this.$t('Cancel'),
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('spider/deploySpider', row._id)
          .then(() => {
            this.$message({
              type: 'success',
              message: 'Deployed successfully'
            })
          })
        this.$st.sendEv('爬虫', '部署')
      })
    },
    onCrawl (row) {
      this.$confirm(this.$t('Are you sure to run this spider?'), this.$t('Notification'), {
        confirmButtonText: this.$t('Confirm'),
        cancelButtonText: this.$t('Cancel')
      })
        .then(() => {
          this.$store.dispatch('spider/crawlSpider', row._id)
            .then(() => {
              this.$message.success(this.$t(`Spider task has been scheduled`))
            })
          this.$st.sendEv('爬虫', '运行')
        })
    },
    onView (row) {
      this.$router.push('/spiders/' +  row._id)
      this.$st.sendEv('爬虫', '查看')
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
      this.$st.sendEv('爬虫', '导入爬虫')
    },
    openImportDialog () {
      this.dialogVisible = true
    },
    onDeployAll () {
      this.$confirm(this.$t('Are you sure to deploy all spiders to active nodes?'), this.$t('Notification'), {
        confirmButtonText: this.$t('Confirm'),
        cancelButtonText: this.$t('Cancel'),
        type: 'warning'
      })
        .then(() => {
          this.$store.dispatch('spider/deployAll')
            .then(() => {
              this.$message.success(this.$t('Deployed all spiders successfully'))
            })
          this.$st.sendEv('爬虫', '部署所有爬虫')
        })
    },
    isShowRun (row) {
      if (this.isCustomized(row)) {
        // customized spider
        if (!row.deploy_ts) {
          return false
        }
        return !!row.cmd
      } else {
        // configurable spider
        return !!row.fields
      }
    },
    isCustomized (row) {
      return row.type === 'customized'
    },
    fetchSiteSuggestions (keyword, callback) {
      this.$request.get('/sites', {
        keyword: keyword,
        page_num: 1,
        page_size: 100
      }).then(response => {
        const data = response.data.items.map(d => {
          d.value = d.name + ' | ' + d.domain
          return d
        })
        callback(data)
      })
    },
    onSiteSelect (item) {
      this.$store.commit('spider/SET_FILTER_SITE', item._id)
      this.$st.sendEv('爬虫', '搜索网站')
    },
    onAddConfigurableSiteSelect (item) {
      this.spiderForm.site = item._id
    },
    onAddConfigurableSpider () {
      this.$refs['addConfigurableForm'].validate(res => {
        if (res) {
          this.addConfigurableLoading = true
          this.$store.dispatch('spider/addSpider')
            .finally(() => {
              this.addConfigurableLoading = false
              this.addConfigurableDialogVisible = false
            })
        }
      })
    },
    onUploadSuccess () {
    }
  },
  created () {
    // take site from params to filter
    this.$store.commit('spider/SET_FILTER_SITE', this.$route.params.domain)

    // fetch spider list
    this.$store.dispatch('spider/getSpiderList')
  },
  mounted () {
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

  .add-spider-wrapper {
    display: flex;
    justify-content: center;

    .add-spider-item {
      cursor: pointer;
      width: 180px;
      font-size: 18px;
      height: 120px;
      margin: 0 20px;
      flex-basis: 40%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .add-spider-item.primary {
      color: #409eff;
      background: rgba(64, 158, 255, .1);
      border: 1px solid rgba(64, 158, 255, .1);
    }

    .add-spider-item.success {
      color: #67c23a;
      background: rgba(103, 194, 58, .1);
      border: 1px solid rgba(103, 194, 58, .1);
    }
  }

  .el-autocomplete {
    width: 100%;
  }

</style>
