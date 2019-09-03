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
        <el-form-item :label="$t('Node')" prop="node_id">
          <el-select v-model="scheduleForm.node_id">
            <el-option :label="$t('All Nodes')" value="000000000000000000000000"></el-option>
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
        <el-form-item :label="$t('Cron')" prop="cron" :rules="cronRules" required>
          <template slot="label">
            <el-tooltip :content="$t('Cron Format: [second] [minute] [hour] [day of month] [month] [day of week]')"
                        placement="top">
              <span>
                {{$t('Cron')}}
                <i class="fa fa-exclamation-circle"></i>
              </span>
            </el-tooltip>
          </template>
          <el-input style="width:calc(100% - 100px);padding-right:10px"
                    v-model="scheduleForm.cron"
                    :placeholder="$t('Cron')">
          </el-input>
          <el-button size="small" style="width:100px" type="primary" @click="onShowCronDialog">{{$t('生成Cron')}}</el-button>
        </el-form-item>
        <el-form-item :label="$t('Execute Command')" prop="params">
          <el-input v-model="spider.cmd"
                    :placeholder="$t('Execute Command')"
                    disabled>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('Parameters')" prop="params">
          <el-input v-model="scheduleForm.params"
                    :placeholder="$t('Parameters')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Schedule Description')" prop="description">
          <el-input v-model="scheduleForm.description" type="textarea"
                    :placeholder="$t('Schedule Description')"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="onCancel">{{$t('Cancel')}}</el-button>
        <el-button size="small" type="primary" @click="onAddSubmit">{{$t('Submit')}}</el-button>
      </span>
    </el-dialog>

    <!--cron generation popup-->
    <el-dialog title="生成 Cron" :visible.sync="showCron">
      <vcrontab @hide="showCron=false" @fill="onCrontabFill" :expression="expression"></vcrontab>
    </el-dialog>

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
          <el-table-column :key="col.name"
                           :property="col.name"
                           :label="$t(col.label)"
                           :sortable="col.sortable"
                           :align="col.align"
                           :width="col.width">
            <template slot-scope="scope">
              {{$t(scope.row[col.name])}}
            </template>
          </el-table-column>
        </template>
        <el-table-column :label="$t('Action')" align="left" width="250" fixed="right">
          <template slot-scope="scope">
            <el-tooltip :content="$t('Edit')" placement="top">
              <el-button type="warning" icon="el-icon-edit" size="mini" @click="onEdit(scope.row)"></el-button>
            </el-tooltip>
            <el-tooltip :content="$t('Remove')" placement="top">
              <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)"></el-button>
            </el-tooltip>
            <el-tooltip v-if="isShowRun(scope.row)" :content="$t('Run')" placement="top">
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
import vcrontab from 'vcrontab'
import {
  mapState
} from 'vuex'

export default {
  name: 'ScheduleList',
  components: { vcrontab },
  data () {
    const cronValidator = (rule, value, callback) => {
      let patArr = []
      for (let i = 0; i < 6; i++) {
        patArr.push('[/*,0-9-]+')
      }
      const pat = '^' + patArr.join(' ') + '( [/*,0-9-]+)?' + '$'
      if (!value) {
        callback(new Error('cron cannot be empty'))
      } else if (!value.match(pat)) {
        callback(new Error('cron format is invalid'))
      }
      callback()
    }
    return {
      columns: [
        { name: 'name', label: 'Name', width: '180' },
        { name: 'cron', label: 'Cron', width: '120' },
        { name: 'node_name', label: 'Node', width: '150' },
        { name: 'spider_name', label: 'Spider', width: '150' },
        { name: 'description', label: 'Description', width: 'auto' }
      ],
      isEdit: false,
      dialogTitle: '',
      dialogVisible: false,
      cronRules: [
        { validator: cronValidator, trigger: 'blur' }
      ],
      showCron: false,
      expression: ''
    }
  },
  computed: {
    ...mapState('schedule', [
      'scheduleList',
      'scheduleForm'
    ]),
    ...mapState('spider', [
      'spiderList'
    ]),
    ...mapState('node', [
      'nodeList'
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
          let action
          if (this.isEdit) {
            action = 'editSchedule'
          } else {
            action = 'addSchedule'
          }
          this.$store.dispatch('schedule/' + action, this.scheduleForm._id)
            .then(() => {
              this.dialogVisible = false
              setTimeout(() => {
                this.$store.dispatch('schedule/getScheduleList')
              }, 100)
            })
        }
      })
      this.$st.sendEv('定时任务', '提交')
    },
    isShowRun () {
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
    onCrawl () {
    },
    onCrontabFill (value) {
      value = value.replace(/[?]/g, '*')
      this.$set(this.scheduleForm, 'cron', value)

      this.$st.sendEv('定时任务', '提交生成Cron', 'cron', this.scheduleForm.cron)
    },
    onShowCronDialog () {
      this.showCron = true
      if (this.expression.split(' ').length < 7) {
        // this.expression = (this.scheduleForm.cron + ' ').replace(/[?]/g, '*')
        this.expression = this.scheduleForm.cron + ' '
      } else {
        // this.expression = this.scheduleForm.cron.replace(/[?]/g, '*')
        this.expression = this.scheduleForm.cron
      }

      this.$st.sendEv('定时任务', '点击生成Cron', 'cron', this.scheduleForm.cron)
    }
  },
  created () {
    this.$store.dispatch('schedule/getScheduleList')
    this.$store.dispatch('spider/getSpiderList')
    this.$store.dispatch('node/getNodeList')
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
</style>
