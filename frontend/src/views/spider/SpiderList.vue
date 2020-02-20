<template>
  <div class="app-container">
    <!--tour-->
    <v-tour
      name="spider-list"
      :steps="tourSteps"
      :callbacks="tourCallbacks"
      :options="$utils.tour.getOptions(true)"
    />
    <v-tour
      name="spider-list-add"
      :steps="tourAddSteps"
      :callbacks="tourAddCallbacks"
      :options="$utils.tour.getOptions(true)"
    />
    <!--./tour-->

    <!--add dialog-->
    <el-dialog :title="$t('Add Spider')"
               width="40%"
               :visible.sync="addDialogVisible"
               :before-close="onAddDialogClose">
      <el-tabs :active-name="activeTabName">
        <!-- customized -->
        <el-tab-pane name="customized" :label="$t('Customized')">
          <el-form :model="spiderForm" ref="addCustomizedForm" inline-message label-width="120px">
            <el-form-item :label="$t('Spider Name')" prop="name" required>
              <el-input id="spider-name" v-model="spiderForm.name" :placeholder="$t('Spider Name')"/>
            </el-form-item>
            <el-form-item :label="$t('Display Name')" prop="display_name" required>
              <el-input id="display-name" v-model="spiderForm.display_name" :placeholder="$t('Display Name')"/>
            </el-form-item>
            <el-form-item :label="$t('Project')" prop="project_id" required>
              <el-select
                v-model="spiderForm.project_id"
                :placeholder="$t('Project')"
                filterable
              >
                <el-option
                  v-for="p in projectList"
                  :key="p._id"
                  :value="p._id"
                  :label="p.name"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('Execute Command')" prop="cmd" required>
              <el-input id="cmd" v-model="spiderForm.cmd" :placeholder="$t('Execute Command')"/>
            </el-form-item>
            <el-form-item :label="$t('Results')" prop="col" required>
              <el-input id="col" v-model="spiderForm.col" :placeholder="$t('Results')"/>
            </el-form-item>
            <el-form-item :label="$t('Upload Zip File')" label-width="120px" name="site">
              <el-upload
                :action="$request.baseUrl + '/spiders'"
                :data="uploadForm"
                :headers="{Authorization:token}"
                :on-success="onUploadSuccess"
                :file-list="fileList"
                :before-upload="beforeUpload"
              >
                <el-button id="upload" size="small" type="primary" icon="el-icon-upload">
                  {{$t('Upload')}}
                </el-button>
              </el-upload>
            </el-form-item>
            <el-form-item :label="$t('Is Git')" prop="is_git">
              <el-switch
                v-model="spiderForm.is_git"
                active-color="#13ce66"
              />
            </el-form-item>
          </el-form>
          <el-alert
            type="warning"
            :closable="false"
            style="margin-bottom: 10px"
          >
            <p>{{$t('You can click "Add" to create an empty spider and upload files later.')}}</p>
            <p>{{$t('OR, you can also click "Upload" and upload a zip file containing your spider project.')}}</p>
            <p style="font-weight: bolder">
              <i class="fa fa-exclamation-triangle"></i> {{$t('NOTE: When uploading a zip file, please zip your' +
              ' spider files from the ROOT DIRECTORY.')}}
            </p>
            <p>
              <template v-if="lang === 'en'">
                You can also upload spiders using <a href="https://docs.crawlab.cn/SDK/CLI.html" target="_blank"
                                                     style="color: #409eff;font-weight: bolder">CLI Tool</a>.
              </template>
              <template v-else-if="lang === 'zh'">
                您也可以利用 <a href="https://docs.crawlab.cn/SDK/CLI.html" target="_blank"
                          style="color: #409eff;font-weight: bolder">CLI 工具</a> 上传爬虫。
              </template>
            </p>
          </el-alert>
          <div class="actions">
            <el-button size="small" type="primary" @click="onAddCustomized">{{$t('Add')}}</el-button>
          </div>
        </el-tab-pane>
        <!-- configurable -->
        <el-tab-pane name="configurable" :label="$t('Configurable')">
          <el-form :model="spiderForm" ref="addConfigurableForm" inline-message label-width="120px">
            <el-form-item :label="$t('Spider Name')" prop="name" required>
              <el-input v-model="spiderForm.name" :placeholder="$t('Spider Name')"/>
            </el-form-item>
            <el-form-item :label="$t('Display Name')" prop="display_name" required>
              <el-input v-model="spiderForm.display_name" :placeholder="$t('Display Name')"/>
            </el-form-item>
            <el-form-item :label="$t('Project')" prop="project_id" required>
              <el-select
                v-model="spiderForm.project_id"
                :placeholder="$t('Project')"
                filterable
              >
                <el-option
                  v-for="p in projectList"
                  :key="p._id"
                  :value="p._id"
                  :label="p.name"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('Template')" prop="template" required>
              <el-select id="template" v-model="spiderForm.template" :value="spiderForm.template"
                         :placeholder="$t('Template')">
                <el-option
                  v-for="template in templateList"
                  :key="template"
                  :label="template"
                  :value="template"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('Results')" prop="col" required>
              <el-input v-model="spiderForm.col" :placeholder="$t('Results')"/>
            </el-form-item>
          </el-form>
          <div class="actions">
            <el-button id="add" size="small" type="primary" @click="onAddConfigurable">{{$t('Add')}}</el-button>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
    <!--./add dialog-->

    <!--running tasks dialog-->
    <el-dialog
      :visible.sync="isRunningTasksDialogVisible"
      :title="`${$t('Latest Tasks')} (${$t('Spider')}: ${activeSpider ? activeSpider.name : ''})`"
      width="920px"
    >
      <el-tabs v-model="activeSpiderTaskStatus">
        <el-tab-pane name="pending" :label="$t('Pending')"/>
        <el-tab-pane name="running" :label="$t('Running')"/>
        <el-tab-pane name="finished" :label="$t('Finished')"/>
        <el-tab-pane name="error" :label="$t('Error')"/>
        <el-tab-pane name="cancelled" :label="$t('Cancelled')"/>
        <el-tab-pane name="abnormal" :label="$t('Abnormal')"/>
      </el-tabs>
      <el-table
        :data="activeNodeList"
        class="table"
        :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
        border
        default-expand-all
      >
        <el-table-column type="expand">
          <template slot-scope="scope">
            <h4 style="margin: 5px 10px">{{$t('Tasks')}}</h4>
            <el-table
              :data="getTasksByNode(scope.row)"
              class="table"
              border
              style="margin: 5px 10px"
              max-height="240px"
              @row-click="onViewTask"
            >
              <el-table-column
                :label="$t('Create Time')"
                prop="create_ts"
                width="140px"
              />
              <el-table-column
                :label="$t('Start Time')"
                prop="start_ts"
                width="140px"
              />
              <el-table-column
                :label="$t('Finish Time')"
                prop="finish_ts"
                width="140px"
              />
              <el-table-column
                :label="$t('Parameters')"
                prop="param"
                width="120px"
              />
              <el-table-column
                :label="$t('Status')"
                width="120px"
              >
                <template slot-scope="scope">
                  <status-tag :status="scope.row.status"/>
                </template>
              </el-table-column>
              <el-table-column
                :label="$t('Results Count')"
                prop="result_count"
                width="80px"
              />
              <el-table-column
                :label="$t('Action')"
                width="auto"
              >
                <template slot-scope="scope">
                  <el-button
                    v-if="['pending', 'running'].includes(scope.row.status)"
                    type="danger"
                    size="mini"
                    icon="el-icon-video-pause"
                    @click="onStop(scope.row, $event)"
                  />
                </template>
              </el-table-column>
            </el-table>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Node')"
          width="150px"
          prop="name"
        />
        <el-table-column
          :label="$t('Status')"
          width="120px"
          prop="status"
        >
          <template slot-scope="scope">
            <el-tag type="info" v-if="scope.row.status === 'offline'">{{$t('Offline')}}</el-tag>
            <el-tag type="success" v-else-if="scope.row.status === 'online'">{{$t('Online')}}</el-tag>
            <el-tag type="danger" v-else>{{$t('Unavailable')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Description')"
          width="auto"
          prop="description"
        />
      </el-table>
      <template slot="footer">
        <el-button type="primary" size="small" @click="isRunningTasksDialogVisible = false">{{$t('Ok')}}</el-button>
      </template>
    </el-dialog>
    <!--./running tasks dialog-->

    <!--crawl confirm dialog-->
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="activeSpiderId"
      @close="crawlConfirmDialogVisible = false"
      @confirm="onCrawlConfirm"
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
              <el-select
                v-model="filter.project_id"
                size="small"
                :placeholder="$t('Project')"
                @change="getList"
              >
                <el-option value="" :label="$t('All Projects')"/>
                <el-option
                  v-for="p in projectList"
                  :key="p._id"
                  :value="p._id"
                  :label="p.name"
                />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-input
                v-model="filter.keyword"
                size="small"
                :placeholder="$t('Spider Name')"
                clearable
                @keyup.enter.native="onSearch"
              >
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
          <el-button
            size="small"
            type="success"
            icon="el-icon-plus"
            class="btn add"
            @click="onAdd"
            style="font-weight: bolder"
          >
            {{$t('Add Spider')}}
          </el-button>

        </div>
      </div>
      <!--./filter-->

      <!--tabs-->
      <el-tabs v-model="filter.type" @tab-click="onClickTab" class="tabs">
        <el-tab-pane :label="$t('All')" name="all" class="all"></el-tab-pane>
        <el-tab-pane :label="$t('Customized')" name="customized" class="customized"></el-tab-pane>
        <el-tab-pane :label="$t('Configurable')" name="configurable" class="configurable"></el-tab-pane>
        <el-tab-pane :label="$t('Long Task')" name="long-task" class="long-task"></el-tab-pane>
      </el-tabs>
      <!--./tabs-->

      <!--legend-->
      <status-legend/>
      <!--./legend-->

      <!--table list-->
      <el-table
        :data="spiderList"
        class="table"
        :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
        border
        @row-click="onRowClick"
        @sort-change="onSortChange"
      >
        <template v-for="col in columns">
          <el-table-column
            v-if="col.name === 'type'"
            :key="col.name"
            :label="$t(col.label)"
            align="left"
            :width="col.width"
            :sortable="col.sortable"
          >
            <template slot-scope="scope">
              {{$t(scope.row.type)}}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'last_5_errors'"
            :key="col.name"
            :label="$t(col.label)"
            :width="col.width"
            :sortable="col.sortable"
            align="center"
          >
            <template slot-scope="scope">
              <div :style="{color:scope.row[col.name]>0?'red':''}">
                {{scope.row[col.name]}}
              </div>
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'cmd'"
            :key="col.name"
            :label="$t(col.label)"
            :width="col.width"
            :sortable="col.sortable"
            align="left"
          >
            <template slot-scope="scope">
              <el-input v-model="scope.row[col.name]"></el-input>
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name.match(/_ts$/)"
            :key="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{getTime(scope.row[col.name])}}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'last_status'"
            :key="col.name"
            :label="$t(col.label)"
            align="left"
            :width="col.width"
            :sortable="col.sortable"
          >
            <template slot-scope="scope">
              <status-tag :status="scope.row.last_status"/>
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="['is_scrapy', 'is_long_task'].includes(col.name)"
            :key="col.name"
            :label="$t(col.label)"
            align="left"
            :width="col.width"
            :sortable="col.sortable"
          >
            <template slot-scope="scope">
              <el-switch
                v-if="scope.row.type === 'customized'"
                v-model="scope.row[col.name]"
                active-color="#13ce66"
                disabled
              />
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'latest_tasks'"
            :key="col.name"
            :label="$t(col.label)"
            :width="col.width"
            :align="col.align"
            class-name="latest-tasks"
          >
            <template slot-scope="scope">
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'pending') > 0"
                type="primary"
                size="small"
              >
                <i class="el-icon-loading"></i>
                {{getTaskCountByStatus(scope.row, 'pending')}}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'running') > 0"
                type="warning"
                size="small"
              >
                <i class="el-icon-loading"></i>
                {{getTaskCountByStatus(scope.row, 'running')}}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'finished') > 0"
                type="success"
                size="small"
              >
                <i class="el-icon-check"></i>
                {{getTaskCountByStatus(scope.row, 'finished')}}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'error') > 0"
                type="danger"
                size="small"
              >
                <i class="el-icon-error"></i>
                {{getTaskCountByStatus(scope.row, 'error')}}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'cancelled') > 0"
                type="info"
                size="small"
              >
                <i class="el-icon-video-pause"></i>
                {{getTaskCountByStatus(scope.row, 'cancelled')}}
              </el-tag>
              <el-tag
                v-if="getTaskCountByStatus(scope.row, 'abnormal') > 0"
                type="danger"
                size="small"
              >
                <i class="el-icon-warning"></i>
                {{getTaskCountByStatus(scope.row, 'abnormal')}}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            v-else
            :key="col.name"
            :property="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            :align="col.align || 'left'"
            :width="col.width"
          >
          </el-table-column>
        </template>
        <el-table-column :label="$t('Action')" align="left" fixed="right" min-width="170px">
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
            <el-tooltip :content="$t('Latest Tasks')" placement="top">
              <el-button
                type="warning"
                icon="fa fa-tasks"
                size="mini"
                @click="onViewRunningTasks(scope.row, $event)"
              />
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
import StatusLegend from '../../components/Status/StatusLegend'

export default {
  name: 'SpiderList',
  components: {
    StatusLegend,
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
      crawlConfirmDialogVisible: false,
      isRunningTasksDialogVisible: false,
      activeSpiderId: undefined,
      activeSpider: undefined,
      filter: {
        project_id: '',
        keyword: '',
        type: 'all'
      },
      sort: {
        sortKey: '',
        sortDirection: null
      },
      types: [],
      spiderFormRules: {
        name: [{ required: true, message: 'Required Field', trigger: 'change' }]
      },
      fileList: [],
      activeTabName: 'customized',
      tourSteps: [
        {
          target: '#tab-customized',
          content: this.$t('View a list of <strong>Customized Spiders</strong>'),
          params: {
            highlight: false
          }
        },
        {
          target: '#tab-configurable',
          content: this.$t('View a list of <strong>Configurable Spiders</strong>'),
          params: {
            highlight: false
          }
        },
        {
          target: '.table',
          content: this.$t('You can view your created spiders here.<br>Click a table row to view <strong>spider details</strong>.'),
          params: {
            placement: 'top'
          }
        },
        {
          target: '.btn.add',
          content: this.$t('Click to add a new spider.<br><br>You can also add a <strong>Customized Spider</strong> through <a href="https://docs.crawlab.cn/Usage/SDK/CLI.html" target="_blank" style="color: #409EFF">CLI Tool</a>.')
        }
      ],
      tourCallbacks: {
        onStop: () => {
          this.$utils.tour.finishTour('spider-list')
        },
        onPreviousStep: (currentStep) => {
          this.$utils.tour.prevStep('spider-list', currentStep)
        },
        onNextStep: (currentStep) => {
          this.$utils.tour.nextStep('spider-list', currentStep)
        }
      },
      tourAddSteps: [
        {
          target: '#tab-customized',
          content: this.$t('<strong>Customized Spider</strong> is a highly customized spider, which is able to run on any programming language and any web crawler framework.'),
          params: {
            placement: 'bottom',
            highlight: false
          }
        },
        {
          target: '#tab-configurable',
          content: this.$t('<strong>Configurable Spider</strong> is a spider defined by config data, aimed at streamlining spider development and improving dev efficiency.'),
          params: {
            placement: 'bottom',
            highlight: false
          }
        },
        {
          target: '#spider-name',
          content: this.$t('Unique identifier for the spider'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '#display-name',
          content: this.$t('How the spider is displayed on Crawlab'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '#cmd',
          content: this.$t('A shell command to be executed when the spider is triggered to run (only available for <strong>Customized Spider</strong>'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '#col',
          content: this.$t('Where the results are stored in the database'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '#upload',
          content: this.$t('Upload a zip file containing all spider files to create the spider (only available for <strong>Customized Spider</strong>)'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '#template',
          content: this.$t('The spider template to create from (only available for <strong>Configurable Spider</strong>)'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '#add',
          content: this.$t('Click to confirm to add the spider'),
          params: {
            placement: 'right'
          }
        }
      ],
      tourAddCallbacks: {
        onStop: () => {
          this.$utils.tour.finishTour('spider-list-add')
        },
        onPreviousStep: (currentStep) => {
          if (currentStep === 7) {
            this.activeTabName = 'customized'
          }
          this.$utils.tour.prevStep('spider-list-add', currentStep)
        },
        onNextStep: (currentStep) => {
          if (currentStep === 6) {
            this.activeTabName = 'configurable'
          }
          this.$utils.tour.nextStep('spider-list-add', currentStep)
        }
      },
      handle: undefined,
      activeSpiderTaskStatus: 'running'
    }
  },
  computed: {
    ...mapState('spider', [
      'importForm',
      'spiderList',
      'spiderForm',
      'spiderTotal',
      'templateList'
    ]),
    ...mapGetters('user', [
      'token'
    ]),
    ...mapState('lang', [
      'lang'
    ]),
    ...mapState('project', [
      'projectList'
    ]),
    ...mapState('node', [
      'nodeList'
    ]),
    uploadForm () {
      return {
        name: this.spiderForm.name,
        display_name: this.spiderForm.display_name,
        col: this.spiderForm.col,
        cmd: this.spiderForm.cmd
      }
    },
    columns () {
      const columns = []
      columns.push({ name: 'display_name', label: 'Name', width: '160', align: 'left', sortable: true })
      columns.push({ name: 'type', label: 'Spider Type', width: '120', sortable: true })
      columns.push({ name: 'is_long_task', label: 'Is Long Task', width: '80' })
      columns.push({ name: 'is_scrapy', label: 'Is Scrapy', width: '80' })
      columns.push({ name: 'latest_tasks', label: 'Latest Tasks', width: '180' })
      columns.push({ name: 'last_status', label: 'Last Status', width: '120' })
      columns.push({ name: 'last_run_ts', label: 'Last Run', width: '140' })
      columns.push({ name: 'update_ts', label: 'Update Time', width: '140' })
      columns.push({ name: 'create_ts', label: 'Create Time', width: '140' })
      columns.push({ name: 'remark', label: 'Remark', width: '140' })
      return columns
    },
    activeNodeList () {
      return this.nodeList.filter(d => {
        return d.status === 'online'
      })
    }
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
      let projectId = '000000000000000000000000'
      if (this.filter.project_id) {
        projectId = this.filter.project_id
      }
      this.$store.commit('spider/SET_SPIDER_FORM', {
        project_id: projectId,
        template: this.templateList[0]
      })
      this.addDialogVisible = true

      setTimeout(() => {
        if (!this.$utils.tour.isFinishedTour('spider-list-add')) {
          this.$tours['spider-list-add'].start()
          this.$st.sendEv('教程', '开始', 'spider-list-add')
        }
      }, 300)
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
        this.$router.push(`/spiders/${res2.data.data._id}`)
        await this.$store.dispatch('spider/getSpiderList')
        this.$st.sendEv('爬虫列表', '添加爬虫', '可配置爬虫')
      })
    },
    onAddCustomized () {
      this.$refs['addCustomizedForm'].validate(async res => {
        if (!res) return
        let res2
        try {
          res2 = await this.$store.dispatch('spider/addSpider')
        } catch (e) {
          this.$message.error(this.$t('Something wrong happened'))
          return
        }
        this.$router.push(`/spiders/${res2.data.data._id}`)
        await this.$store.dispatch('spider/getSpiderList')
        this.$st.sendEv('爬虫列表', '添加爬虫', '自定义爬虫')
      })
    },
    onRefresh () {
      this.getList()
      this.$st.sendEv('爬虫列表', '刷新')
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
      }).then(async () => {
        await this.$store.dispatch('spider/deleteSpider', row._id)
        this.$message({
          type: 'success',
          message: 'Deleted successfully'
        })
        await this.getList()
        this.$st.sendEv('爬虫列表', '删除爬虫')
      })
    },
    onCrawl (row, ev) {
      ev.stopPropagation()
      this.crawlConfirmDialogVisible = true
      this.activeSpiderId = row._id
      this.$st.sendEv('爬虫列表', '点击运行')
    },
    onCrawlConfirm () {
      setTimeout(() => {
        this.getList()
      }, 1000)
    },
    onView (row, ev) {
      ev.stopPropagation()
      this.$router.push('/spiders/' + row._id)
      this.$st.sendEv('爬虫列表', '查看爬虫')
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
      this.$st.sendEv('爬虫列表', '导入爬虫')
    },
    openImportDialog () {
      this.dialogVisible = true
    },
    isShowRun (row) {
      if (!this.isCustomized(row)) return true
      return !!row.cmd
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
    onUploadSuccess (res) {
      // clear fileList
      this.fileList = []

      // fetch spider list
      setTimeout(() => {
        this.getList()
      }, 500)

      // message
      this.$message.success(this.$t('Uploaded spider files successfully'))

      // navigate to spider detail
      this.$router.push(`/spiders/${res.data._id}`)
    },
    beforeUpload (file) {
      return new Promise((resolve, reject) => {
        this.$refs['addCustomizedForm'].validate(res => {
          if (res) {
            resolve()
          } else {
            reject(new Error('form validation error'))
          }
        })
      })
    },
    getTime (str) {
      if (!str || str.match('^0001')) return 'NA'
      return dayjs(str).format('YYYY-MM-DD HH:mm:ss')
    },
    onRowClick (row, column, event) {
      this.onView(row, event)
    },
    onSortChange ({ column, prop, order }) {
      this.sort.sortKey = order ? prop : ''
      this.sort.sortDirection = order
      this.getList()
    },
    onClickTab (tab) {
      this.filter.type = tab.name
      this.getList()
    },
    async getList () {
      let params = {
        page_num: this.pagination.pageNum,
        page_size: this.pagination.pageSize,
        sort_key: this.sort.sortKey,
        sort_direction: this.sort.sortDirection,
        keyword: this.filter.keyword,
        type: this.filter.type,
        project_id: this.filter.project_id
      }
      await this.$store.dispatch('spider/getSpiderList', params)

      // 更新当前爬虫（任务列表）
      this.updateActiveSpider()
    },
    getTasksByStatus (row, status) {
      if (!row.latest_tasks) return []
      return row.latest_tasks.filter(d => d.status === status)
    },
    getTaskCountByStatus (row, status) {
      return this.getTasksByStatus(row, status).length
    },
    updateActiveSpider () {
      if (this.activeSpider) {
        for (let i = 0; i < this.spiderList.length; i++) {
          const spider = this.spiderList[i]
          if (this.activeSpider._id === spider._id) {
            this.activeSpider = spider
          }
        }
      }
    },
    onViewRunningTasks (row, ev) {
      ev.stopPropagation()
      this.activeSpider = row
      this.isRunningTasksDialogVisible = true
    },
    getTasksByNode (row) {
      if (!this.activeSpider.latest_tasks) {
        return []
      }
      return this.activeSpider.latest_tasks
        .filter(d => d.node_id === row._id && d.status === this.activeSpiderTaskStatus)
        .map(d => {
          d = JSON.parse(JSON.stringify(d))
          d.create_ts = d.create_ts.match('^0001') ? 'NA' : dayjs(d.create_ts).format('YYYY-MM-DD HH:mm:ss')
          d.start_ts = d.start_ts.match('^0001') ? 'NA' : dayjs(d.start_ts).format('YYYY-MM-DD HH:mm:ss')
          d.finish_ts = d.finish_ts.match('^0001') ? 'NA' : dayjs(d.finish_ts).format('YYYY-MM-DD HH:mm:ss')
          return d
        })
    },
    onViewTask (row) {
      this.$router.push(`/tasks/${row._id}`)
      this.$st.sendEv('爬虫列表', '任务列表', '查看任务')
    },
    async onStop (row, ev) {
      ev.stopPropagation()
      const res = await this.$store.dispatch('task/cancelTask', row._id)
      if (!res.data.error) {
        this.$message.success(`Task "${row._id}" has been sent signal to stop`)
        this.getList()
      }
    }
  },
  async created () {
    // fetch project list
    await this.$store.dispatch('project/getProjectList')

    // project id
    if (this.$route.params.project_id) {
      this.filter.project_id = this.$route.params.project_id
    }

    // fetch node list
    await this.$store.dispatch('node/getNodeList')

    // fetch spider list
    await this.getList()

    // fetch template list
    await this.$store.dispatch('spider/getTemplateList')

    // periodically fetch spider list
    this.handle = setInterval(() => {
      this.getList()
    }, 15000)
  },
  mounted () {
    const vm = this
    this.$nextTick(() => {
      vm.$store.commit('spider/SET_SPIDER_FORM', this.spiderForm)
    })

    if (!this.$utils.tour.isFinishedTour('spider-list')) {
      this.$tours['spider-list'].start()
      this.$st.sendEv('教程', '开始', 'spider-list')
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

  .el-table >>> .latest-tasks .el-tag {
    margin: 3px 3px 0 0;
  }
</style>
