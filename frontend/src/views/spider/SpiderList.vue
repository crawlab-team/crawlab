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
      <el-tabs :active-name="spiderType">
        <el-tab-pane name="configurable" :label="$t('Configurable')">
          <el-form :model="spiderForm" ref="addConfigurableForm" inline-message label-width="120px">
            <el-form-item :label="$t('Spider Name')" prop="name" required>
              <el-input v-model="spiderForm.name" :placeholder="$t('Spider Name')"/>
            </el-form-item>
            <el-form-item :label="$t('Display Name')" prop="display_name" required>
              <el-input v-model="spiderForm.display_name" :placeholder="$t('Display Name')"/>
            </el-form-item>
            <el-form-item :label="$t('Results')" prop="col" required>
              <el-input v-model="spiderForm.col" :placeholder="$t('Results')"/>
            </el-form-item>
          </el-form>
          <div class="actions">
            <el-button type="primary" @click="onAddConfigurable">{{$t('Add')}}</el-button>
          </div>
        </el-tab-pane>
        <el-tab-pane name="customized" :label="$t('Customized')">
          <el-form :model="spiderForm" ref="addCustomizedForm" inline-message>
            <el-form-item :label="$t('Upload Zip File')" label-width="120px" name="site">
              <el-upload
                :action="$request.baseUrl + '/spiders'"
                :headers="{Authorization:token}"
                :on-change="onUploadChange"
                :on-success="onUploadSuccess"
                :file-list="fileList">
                <el-button size="small" type="primary" icon="el-icon-upload">{{$t('Upload')}}</el-button>
              </el-upload>
            </el-form-item>
          </el-form>
          <el-alert type="error" :title="$t('Please zip your spider files from the root directory')"
                    :closable="false"></el-alert>
        </el-tab-pane>
      </el-tabs>
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

    </el-dialog>
    <!--./customized spider dialog-->

    <!--crawl confirm dialog-->
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="activeSpiderId"
      @close="crawlConfirmDialogVisible = false"
    />
    <!--./crawl confirm dialog-->

    <el-card style="border-radius: 0">
      <!--filter-->
      <div class="filter">
        <div class="left">
          <el-form :inline="true">
            <!--            <el-form-item>-->
            <!--              <el-select clearable @change="onSpiderTypeChange" placeholder="爬虫类型" size="small" v-model="filter.type">-->
            <!--                <el-option v-for="item in types" :value="item.type" :key="item.type"-->
            <!--                           :label="item.type === 'customized'? '自定义':item.type "/>-->
            <!--              </el-select>-->
            <!--            </el-form-item>-->
            <el-form-item>
              <el-input clearable @keyup.enter.native="onSearch" size="small" placeholder="名称" v-model="filter.keyword">
                <i slot="suffix" class="el-input__icon el-icon-search"></i>
              </el-input>
            </el-form-item>
            <el-form-item>
              <el-button size="small" type="success"
                         class="btn refresh"
                         @click="onRefresh">
                {{$t('Search')}}
              </el-button>
            </el-form-item>
          </el-form>
        </div>
        <div class="right">
          <el-button size="small" v-if="false" type="primary" icon="fa fa-download" @click="openImportDialog">
            {{$t('Import Spiders')}}
          </el-button>
          <el-button size="small" type="success"
                     icon="el-icon-plus"
                     class="btn add"
                     @click="onAdd">
            {{$t('Add Spider')}}
          </el-button>

        </div>
      </div>
      <!--./filter-->

      <!--tabs-->
      <el-tabs v-model="filter.type" @tab-click="onClickTab">
        <el-tab-pane :label="$t('All')" name="all"></el-tab-pane>
        <el-tab-pane :label="$t('Configurable')" name="configurable"></el-tab-pane>
        <el-tab-pane :label="$t('Customized')" name="customized"></el-tab-pane>
      </el-tabs>
      <!--./tabs-->

      <!--table list-->
      <el-table :data="spiderList"
                class="table"
                :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
                border
                @row-click="onRowClick"
      >
        <template v-for="col in columns">
          <el-table-column v-if="col.name === 'type'"
                           :key="col.name"
                           :label="$t(col.label)"
                           align="left"
                           :width="col.width">
            <template slot-scope="scope">
              {{$t(scope.row.type)}}
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
          <el-table-column v-else-if="col.name === 'cmd'"
                           :key="col.name"
                           :label="$t(col.label)"
                           :width="col.width"
                           align="left">
            <template slot-scope="scope">
              <el-input v-model="scope.row[col.name]"></el-input>
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
          <el-table-column v-else-if="col.name === 'last_status'"
                           :key="col.name"
                           :label="$t(col.label)"
                           align="left" :width="col.width">
            <template slot-scope="scope">
              <status-tag :status="scope.row.last_status"/>
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
        <el-table-column :label="$t('Action')" align="left" fixed="right">
          <template slot-scope="scope">
            <el-tooltip :content="$t('View')" placement="top">
              <el-button type="primary" icon="el-icon-search" size="mini"
                         @click="onView(scope.row, $event)"></el-button>
            </el-tooltip>
            <el-tooltip :content="$t('Remove')" placement="top">
              <el-button type="danger" icon="el-icon-delete" size="mini"
                         @click="onRemove(scope.row, $event)"></el-button>
            </el-tooltip>
            <el-tooltip v-if="!isShowRun(scope.row)" :content="$t('No command line')" placement="top">
              <el-button disabled type="success" icon="fa fa-bug" size="mini"
                         @click="onCrawl(scope.row, $event)"></el-button>
            </el-tooltip>
            <el-tooltip v-else :content="$t('Run')" placement="top">
              <el-button type="success" icon="fa fa-bug" size="mini" @click="onCrawl(scope.row, $event)"></el-button>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          @current-change="onPageNumChange"
          @size-change="onPageSizeChange"
          :current-page.sync="pagination.pageNum"
          :page-sizes="[10, 20, 50, 100]"
          :page-size.sync="pagination.pageSize"
          layout="sizes, prev, pager, next"
          :total="spiderTotal">
        </el-pagination>
      </div>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
import {
  mapState,
  mapGetters
} from 'vuex'
import dayjs from 'dayjs'
import CrawlConfirmDialog from '../../components/Common/CrawlConfirmDialog'
import StatusTag from '../../components/Status/StatusTag'
import request from '../../api/request'

export default {
  name: 'SpiderList',
  components: {
    CrawlConfirmDialog,
    StatusTag
  },
  data () {
    return {
      pagination: {
        pageNum: 1,
        pageSize: 10
      },
      importLoading: false,
      addConfigurableLoading: false,
      isEditMode: false,
      dialogVisible: false,
      addDialogVisible: false,
      addConfigurableDialogVisible: false,
      addCustomizedDialogVisible: false,
      crawlConfirmDialogVisible: false,
      activeSpiderId: undefined,
      filter: {
        keyword: '',
        type: 'all'
      },
      types: [],
      columns: [
        { name: 'display_name', label: 'Name', width: '160', align: 'left' },
        { name: 'type', label: 'Spider Type', width: '120' },
        { name: 'last_status', label: 'Last Status', width: '120' },
        { name: 'last_run_ts', label: 'Last Run', width: '140' },
        // { name: 'update_ts', label: 'Update Time', width: '140' },
        { name: 'remark', label: 'Remark', width: '140' }
      ],
      spiderFormRules: {
        name: [{ required: true, message: 'Required Field', trigger: 'change' }]
      },
      fileList: [],
      spiderType: 'configurable'
    }
  },
  computed: {
    ...mapState('spider', [
      'importForm',
      'spiderList',
      'spiderForm',
      'spiderTotal'
    ]),
    ...mapGetters('user', [
      'token'
    ])
  },
  methods: {
    onSpiderTypeChange (val) {
      this.filter.type = val
      this.getList()
    },
    onPageSizeChange (val) {
      this.pagination.pageSize = val
      this.getList()
    },
    onPageNumChange (val) {
      this.pagination.pageNum = val
      this.getList()
    },
    onSearch () {
      this.getList()
    },
    onAdd () {
      this.$store.commit('spider/SET_SPIDER_FORM', {})
      this.addDialogVisible = true
    },
    onAddConfigurable () {
      this.$refs['addConfigurableForm'].validate(async res => {
        if (!res) return

        let res2
        try {
          res2 = await this.$store.dispatch('spider/addConfigSpider')
        } catch (e) {
          this.$message.error(this.$t('Something wrong happened'))
          return
        }
        await this.$store.dispatch('spider/getSpiderList')
        this.$router.push(`/spiders/${res2.data.data._id}`)
        this.$st.sendEv('爬虫', '添加爬虫-可配置爬虫')
      })
    },
    onAddCustomized () {
      this.addDialogVisible = false
      this.addCustomizedDialogVisible = true
      this.$st.sendEv('爬虫', '添加爬虫-自定义爬虫')
    },
    onRefresh () {
      this.getList()
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
    onRemove (row, ev) {
      ev.stopPropagation()
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
    onCrawl (row, ev) {
      ev.stopPropagation()
      this.crawlConfirmDialogVisible = true
      this.activeSpiderId = row._id
      this.$st.sendEv('爬虫', '点击运行')
    },
    onView (row, ev) {
      ev.stopPropagation()
      this.$router.push('/spiders/' + row._id)
      this.$st.sendEv('爬虫', '查看')
    },
    onImport () {
      this.$refs.importForm.validate(valid => {
        if (valid) {
          this.importLoading = true
          // TODO: switch between github / gitlab / svn
          this.$store.dispatch('spider/importGithub')
            .then(response => {
              this.$message.success('Import repo successfully')
              this.getList()
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
    isShowRun (row) {
      if (row.cmd) {
        return true
      } else {
        return false
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
    onUploadChange () {
    },
    onUploadSuccess () {
      // clear fileList
      this.fileList = []

      // fetch spider list
      setTimeout(() => {
        this.getList()
      }, 500)

      // close popup
      this.addCustomizedDialogVisible = false
    },
    getTime (str) {
      if (!str || str.match('^0001')) return 'NA'
      return dayjs(str).format('YYYY-MM-DD HH:mm:ss')
    },
    onRowClick (row, column, event) {
      this.onView(row, event)
    },
    onClickTab (tab) {
      this.filter.type = tab.name
      this.getList()
    },
    getList () {
      let params = {
        pageNum: this.pagination.pageNum,
        pageSize: this.pagination.pageSize,
        keyword: this.filter.keyword,
        type: this.filter.type
      }
      this.$store.dispatch('spider/getSpiderList', params)
    },
    getTypes () {
      request.get(`/spider/types`).then(resp => {
        this.types = resp.data.data
      })
    }
  },
  created () {
    this.getTypes()
    // fetch spider list
    this.getList()
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
    margin-top: 8px;
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

    .add-spider-item.info {
      color: #909399;
      background: #f4f4f5;
      border: 1px solid #e9e9eb;
    }

  }

  .el-autocomplete {
    width: 100%;
  }

</style>

<style scoped>
  .el-table >>> tr {
    cursor: pointer;
  }

  .actions {
    text-align: right;
  }
</style>
