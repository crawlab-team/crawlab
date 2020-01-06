<template>
  <div class="app-container">
    <!--add popup-->
    <el-dialog
      :title="$t(dialogTitle)"
      :visible.sync="dialogVisible"
      width="60%"
      :before-close="onDialogClose">
      <el-form label-width="180px"
               :model="scheduleForm"
               :inline-message="true"
               ref="scheduleForm"
               label-position="right">
        <el-form-item :label="$t('Schedule Name')" prop="name" required>
          <el-input v-model="scheduleForm.name" :placeholder="$t('Schedule Name')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Run Type')" prop="run_type" required>
          <el-select v-model="scheduleForm.run_type" :placeholder="$t('Run Type')">
            <el-option value="all-nodes" :label="$t('All Nodes')"/>
            <el-option value="selected-nodes" :label="$t('Selected Nodes')"/>
            <el-option value="random" :label="$t('Random')"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="scheduleForm.run_type === 'selected-nodes'" :label="$t('Nodes')" prop="node_ids" required>
          <el-select v-model="scheduleForm.node_ids" :placeholder="$t('Nodes')" multiple filterable>
            <el-option
              v-for="op in nodeList"
              :key="op._id"
              :value="op._id"
              :label="op.name"
              :disabled="op.status === 'offline'"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="!isDisabledSpiderSchedule" :label="$t('Spider')" prop="spider_id" required>
          <el-select
            v-model="scheduleForm.spider_id"
            :placeholder="$t('Spider')"
            filterable
            :disabled="isDisabledSpiderSchedule"
          >
            <el-option
              v-for="op in spiderList"
              :key="op._id"
              :value="op._id"
              :label="`${op.display_name} (${op.name})`"
              :disabled="isDisabledSpider(op)"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-else :label="$t('Spider')" required>
          <el-select
            :value="spiderId"
            :placeholder="$t('Spider')"
            filterable
            :disabled="isDisabledSpiderSchedule"
          >
            <el-option
              v-for="op in spiderList"
              :key="op._id"
              :value="op._id"
              :label="`${op.display_name} (${op.name})`"
              :disabled="isDisabledSpider(op)"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('Cron')" prop="cron" required>
          <el-input
            v-model="scheduleForm.cron"
            :placeholder="`${$t('[minute] [hour] [day] [month] [day of week]')}`"
          >
          </el-input>
          <!--<el-button size="small" style="width:100px" type="primary" @click="onShowCronDialog">{{$t('schedules.add_cron')}}</el-button>-->
        </el-form-item>
        <el-form-item :label="$t('Execute Command')" prop="params">
          <el-input v-model="spider.cmd"
                    :placeholder="$t('Execute Command')"
                    disabled>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('Parameters')" prop="param">
          <el-input v-model="scheduleForm.param"
                    :placeholder="$t('Parameters')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Schedule Description')" prop="description">
          <el-input v-model="scheduleForm.description" type="textarea"
                    :placeholder="$t('Schedule Description')"></el-input>
        </el-form-item>
      </el-form>
      <!--取消、保存-->
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="onCancel">{{$t('Cancel')}}</el-button>
        <el-button size="small" type="primary" @click="onAddSubmit">{{$t('Submit')}}</el-button>
      </span>
    </el-dialog>

    <!--cron generation popup-->
    <!--<el-dialog title="生成 Cron" :visible.sync="showCron">-->
    <!--<vcrontab @hide="showCron=false" @fill="onCrontabFill" :expression="expression"></vcrontab>-->
    <!--</el-dialog>-->

    <el-card style="border-radius: 0">
      <!--filter-->
      <div class="filter">
        <div class="right">
          <el-button size="small" type="primary"
                     icon="el-icon-plus"
                     class="refresh"
                     @click="onAdd">
            {{$t('Add Schedule')}}
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--table list-->
      <el-table :data="filteredTableData"
                class="table" height="500"
                :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
                border>
        <template v-for="col in columns">
          <el-table-column v-if="col.name === 'status'"
                           :key="col.name"
                           :property="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              <el-tooltip v-if="scope.row[col.name] === 'error'" :content="$t(scope.row['message'])" placement="top">
                <el-tag class="status-tag" type="danger">
                  {{scope.row[col.name] ? $t(scope.row[col.name]) : $t('NA')}}
                </el-tag>
              </el-tooltip>
              <el-tag class="status-tag" v-else>
                {{scope.row[col.name] ? $t(scope.row[col.name]) : $t('NA')}}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'run_type'" :key="col.name" :label="$t(col.label)" :width="col.width">
            <template slot-scope="scope">
              <template v-if="scope.row.run_type === 'all-nodes'">{{$t('All Nodes')}}</template>
              <template v-else-if="scope.row.run_type === 'selected-nodes'">{{$t('Selected Nodes')}}</template>
              <template v-else-if="scope.row.run_type === 'random'">{{$t('Random')}}</template>
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'node_names'" :key="col.name" :label="$t(col.label)" :width="col.width">
            <template slot-scope="scope">
              {{scope.row.nodes.map(d => d.name).join(', ')}}
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'enable'" :key="col.name" :label="$t(col.label)" :width="col.width">
            <template slot-scope="scope">
              <el-switch
                v-model="scope.row.enabled"
                active-color="#13ce66"
                inactive-color="#ff4949"
                @change="onEnabledChange(scope.row)"
              />
            </template>
          </el-table-column>
          <el-table-column v-else :key="col.name"
                           :property="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              {{scope.row[col.name]}}
            </template>
          </el-table-column>
        </template>
        <el-table-column :label="$t('Action')" align="left" width="130" fixed="right">
          <template slot-scope="scope">
            <!-- 编辑 -->
            <el-tooltip :content="$t('Edit')" placement="top">
              <el-button type="warning" icon="el-icon-edit" size="mini" @click="onEdit(scope.row)"></el-button>
            </el-tooltip>
            <!-- 删除 -->
            <el-tooltip :content="$t('Remove')" placement="top">
              <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)"></el-button>
            </el-tooltip>
            <!--<el-tooltip :content="$t(getStatusTooltip(scope.row))" placement="top">-->
              <!--<el-button type="success" icon="fa fa-bug" size="mini" @click="onCrawl(scope.row)"></el-button>-->
            <!--</el-tooltip>-->
          </template>
        </el-table-column>
      </el-table>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
// import vcrontab from 'vcrontab'
import request from '../../api/request'
import {
  mapState
} from 'vuex'

export default {
  name: 'ScheduleList',
  data () {
    return {
      columns: [
        { name: 'name', label: 'Name', width: '150px' },
        { name: 'cron', label: 'Cron', width: '120px' },
        { name: 'run_type', label: 'Run Type', width: '120px' },
        { name: 'node_names', label: 'Node', width: '150px' },
        { name: 'spider_name', label: 'Spider', width: '150px' },
        { name: 'param', label: 'Parameters', width: '150px' },
        { name: 'description', label: 'Description', width: '200px' },
        { name: 'enable', label: 'Enable/Disable', width: '120px' }
        // { name: 'status', label: 'Status', width: '100px' }
      ],
      isEdit: false,
      dialogTitle: '',
      dialogVisible: false,
      showCron: false,
      expression: '',
      spiderList: [],
      nodeList: []
    }
  },
  computed: {
    ...mapState('schedule', [
      'scheduleList',
      'scheduleForm'
    ]),
    filteredTableData () {
      return this.scheduleList
    },
    spider () {
      for (let i = 0; i < this.spiderList.length; i++) {
        if (this.spiderList[i]._id === this.scheduleForm.spider_id) {
          return this.spiderList[i]
        }
      }
      return {}
    },
    isDisabledSpiderSchedule () {
      return false
    }
  },
  methods: {
    onDialogClose () {
      this.dialogVisible = false
    },
    onCancel () {
      this.dialogVisible = false
    },
    onAdd () {
      this.isEdit = false
      this.dialogVisible = true
      this.$store.commit('schedule/SET_SCHEDULE_FORM', { node_ids: [] })
      this.$st.sendEv('定时任务', '添加定时任务')
    },
    onAddSubmit () {
      this.$refs.scheduleForm.validate(res => {
        if (res) {
          const form = JSON.parse(JSON.stringify(this.scheduleForm))
          form.cron = '0 ' + this.scheduleForm.cron
          if (this.isEdit) {
            request.post(`/schedules/${this.scheduleForm._id}`, form).then(response => {
              if (response.data.error) {
                this.$message.error(response.data.error)
                return
              }
              this.dialogVisible = false
              this.$store.dispatch('schedule/getScheduleList')
              this.$message.success(this.$t('The schedule has been saved'))
            })
          } else {
            request.put('/schedules', form).then(response => {
              if (response.data.error) {
                this.$message.error(response.data.error)
                return
              }
              this.dialogVisible = false
              this.$store.dispatch('schedule/getScheduleList')
              this.$message.success(this.$t('The schedule has been added'))
            })
          }
        }
      })
      this.$st.sendEv('定时任务', '提交定时任务')
    },
    isShowRun (row) {
    },
    onEdit (row) {
      this.$store.commit('schedule/SET_SCHEDULE_FORM', row)
      this.dialogVisible = true
      this.isEdit = true
      this.$st.sendEv('定时任务', '修改定时任务')
    },
    onRemove (row) {
      this.$confirm(this.$t('Are you sure to delete the schedule task?'), this.$t('Notification'), {
        confirmButtonText: this.$t('Confirm'),
        cancelButtonText: this.$t('Cancel'),
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('schedule/removeSchedule', row._id)
          .then(() => {
            setTimeout(() => {
              this.$store.dispatch('schedule/getScheduleList')
              this.$message.success(this.$t('The schedule has been removed'))
            }, 100)
          })
      }).catch(() => {
      })
      this.$st.sendEv('定时任务', '删除定时任务')
    },
    isDisabledSpider (spider) {
      if (spider.type === 'customized') {
        return !spider.cmd
      } else {
        return false
      }
    },
    getStatusTooltip (row) {
      if (row.status === 'stop') {
        return 'Start'
      } else if (row.status === 'running') {
        return 'Stop'
      } else if (row.status === 'error') {
        return 'Start'
      }
    },
    async onEnabledChange (row) {
      let res
      if (row.enabled) {
        res = await this.$store.dispatch('schedule/enableSchedule', row._id)
      } else {
        res = await this.$store.dispatch('schedule/disableSchedule', row._id)
      }
      if (!res || res.data.error) {
        this.$message.error(this.$t(`${row.enabled ? 'Enabling' : 'Disabling'} the schedule unsuccessful`))
      } else {
        this.$message.success(this.$t(`${row.enabled ? 'Enabling' : 'Disabling'} the schedule successful`))
      }
    }
  },
  created () {
    this.$store.dispatch('schedule/getScheduleList')

    // 节点列表
    request.get('/nodes', {}).then(response => {
      this.nodeList = response.data.data.map(d => {
        d.systemInfo = {
          os: '',
          arch: '',
          num_cpu: '',
          executables: []
        }
        return d
      })
    })

    // 爬虫列表
    request.get('/spiders', {})
      .then(response => {
        this.spiderList = response.data.data.list || []
      })
  }
}
</script>

<style scoped>
  .filter .right {
    text-align: right;
  }

  .table {
    min-height: 360px;
    margin-top: 10px;
  }

  .status-tag {
    cursor: pointer;
  }
</style>
