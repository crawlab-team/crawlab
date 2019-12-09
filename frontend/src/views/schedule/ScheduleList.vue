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
        <el-form-item :label="$t('Node')" prop="node_id" required>
          <el-select v-model="scheduleForm.node_id">
            <!--<el-option :label="$t('All Nodes')" value="000000000000000000000000"></el-option>-->
            <el-option
              v-for="op in nodeList"
              :key="op._id"
              :value="op._id"
              :label="op.name"
              :disabled="op.status === 'offline'"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('Spider')" prop="spider_id" required>
          <el-select v-model="scheduleForm.spider_id" filterable>
            <el-option
              v-for="op in spiderList"
              :key="op._id"
              :value="op._id"
              :label="op.name"
              :disabled="!op.cmd"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <!--:rules="cronRules"-->
        <el-form-item :label="$t('schedules.cron')" prop="cron"  required>
          <!--<template slot="label">-->
            <!--<el-tooltip :content="$t('schedules.cron_format')"-->
                        <!--placement="top">-->
              <!--<span>-->
                <!--{{$t('schedules.cron')}}-->
                <!--<i class="fa fa-exclamation-circle"></i>-->
              <!--</span>-->
            <!--</el-tooltip>-->
          <!--</template>-->
          <el-input style="padding-right:10px"
                    v-model="scheduleForm.cron"
                    :placeholder="$t('schedules.cron')">
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
                class="table"
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
        <el-table-column :label="$t('Action')" align="left" width="180px" fixed="right">
          <template slot-scope="scope">
            <!-- 编辑 -->
            <el-tooltip :content="$t('Edit')" placement="top">
              <el-button type="warning" icon="el-icon-edit" size="mini" @click="onEdit(scope.row)"></el-button>
            </el-tooltip>
            <!-- 删除 -->
            <el-tooltip :content="$t('Remove')" placement="top">
              <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)"></el-button>
            </el-tooltip>
            <el-tooltip content="暂停/运行" placement="top">
              <el-button type="success" icon="fa fa-bug" size="mini" @click="onCrawl(scope.row)"></el-button>
            </el-tooltip>
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
        { name: 'name', label: 'Name', width: '180' },
        { name: 'cron', label: 'schedules.cron', width: '120' },
        { name: 'node_name', label: 'Node', width: '150' },
        { name: 'spider_name', label: 'Spider', width: '150' },
        { name: 'param', label: 'Parameters', width: '150' },
        { name: 'description', label: 'Description', width: 'auto' },
        { name: 'status', label: 'Status', width: 'auto' }
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
      this.$store.commit('schedule/SET_SCHEDULE_FORM', {})
      this.$st.sendEv('定时任务', '添加')
    },
    onAddSubmit () {
      this.$refs.scheduleForm.validate(res => {
        if (res) {
          if (this.isEdit) {
            request.post(`/schedules/${this.scheduleForm._id}`, this.scheduleForm).then(response => {
              if (response.data.error) {
                this.$message.error(response.data.error)
                return
              }
              this.dialogVisible = false
              this.$store.dispatch('schedule/getScheduleList')
            })
          } else {
            request.put('/schedules', this.scheduleForm).then(response => {
              if (response.data.error) {
                this.$message.error(response.data.error)
                return
              }
              this.dialogVisible = false
              this.$store.dispatch('schedule/getScheduleList')
            })
          }
        }
      })
      this.$st.sendEv('定时任务', '提交')
    },
    isShowRun (row) {
    },
    onEdit (row) {
      this.$store.commit('schedule/SET_SCHEDULE_FORM', row)
      this.dialogVisible = true
      this.isEdit = true
      this.$st.sendEv('定时任务', '修改', 'id', row._id)
    },
    onRemove (row) {
      this.$store.dispatch('schedule/removeSchedule', row._id)
        .then(() => {
          setTimeout(() => {
            this.$store.dispatch('schedule/getScheduleList')
            this.$message.success(`Schedule "${row.name}" has been removed`)
          }, 100)
        })
      this.$st.sendEv('定时任务', '删除', 'id', row._id)
    },
    onCrawl (row) {
      // 停止定时任务
      if (!row.status || row.status === 'running') {
        this.$confirm('确定停止定时任务?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.$store.dispatch('schedule/stopSchedule', row._id)
            .then((resp) => {
              if (resp.data.status === 'ok') {
                this.$store.dispatch('schedule/getScheduleList')
                return
              }
              this.$message({
                type: 'error',
                message: resp.data.error
              })
            })
        }).catch(() => {})
      }
      // 运行定时任务
      if (row.status === 'stop') {
        this.$confirm('确定运行定时任务?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.$store.dispatch('schedule/runSchedule', row._id)
            .then((resp) => {
              if (resp.data.status === 'ok') {
                this.$store.dispatch('schedule/getScheduleList')
                return
              }
              this.$message({
                type: 'error',
                message: resp.data.error
              })
            })
        }).catch(() => {})
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
