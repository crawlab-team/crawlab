<template>
  <div class="info-view">
    <el-row>
      <el-form label-width="150px"
               :model="taskForm"
               ref="nodeForm"
               class="node-form"
               label-position="right">
        <el-form-item :label="$t('Task ID')">
          <el-input v-model="taskForm._id" placeholder="Task ID" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Status')">
          <status-tag :status="taskForm.status"/>
          <el-badge
            v-if="errorLogData.length > 0"
            :value="errorLogData.length"
            style="margin-left:10px; cursor:pointer;"
          >
            <el-tag type="danger" @click="onClickLogWithErrors">
              <i class="el-icon-warning"></i>
              {{$t('Log with errors')}}
            </el-tag>
          </el-badge>
          <el-tag
            v-if="taskForm.status === 'finished' && taskForm.result_count === 0"
            type="danger"
            style="margin-left: 10px"
          >
            <i class="el-icon-warning"></i>
            {{$t('Empty results')}}
          </el-tag>
        </el-form-item>
        <el-form-item :label="$t('Log File Path')">
          <el-input v-model="taskForm.log_path" placeholder="Log File Path" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Parameters')">
          <el-input v-model="taskForm.param" placeholder="Parameters" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Create Time')">
          <el-input :value="getTime(taskForm.create_ts)" placeholder="Create Time" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Start Time')">
          <el-input :value="getTime(taskForm.start_ts)" placeholder="Start Time" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Finish Time')">
          <el-input :value="getTime(taskForm.finish_ts)" placeholder="Finish Time" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Wait Duration (sec)')">
          <el-input :value="getWaitDuration(taskForm)" placeholder="Wait Duration" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Runtime Duration (sec)')">
          <el-input :value="getRuntimeDuration(taskForm)" placeholder="Runtime Duration" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Total Duration (sec)')">
          <el-input :value="getTotalDuration(taskForm)" placeholder="Runtime Duration" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Results Count')">
          <el-input v-model="taskForm.result_count" placeholder="Results Count" disabled></el-input>
        </el-form-item>
        <!--<el-form-item :label="$t('Average Results Count per Second')">-->
        <!--<el-input v-model="taskForm.avg_num_results" placeholder="Average Results Count per Second" disabled>-->
        <!--</el-input>-->
        <!--</el-form-item>-->
        <el-form-item :label="$t('Error Message')" v-if="taskForm.status === 'error'">
          <div class="error-message">
            {{ taskForm.error }}
          </div>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container">
      <el-button v-if="isRunning" size="small" type="danger" @click="onStop" icon="el-icon-video-pause">
        {{$t('Stop')}}
      </el-button>
      <!--<el-button type="danger" @click="onRestart">Restart</el-button>-->
    </el-row>
  </div>
</template>

<script>
import {
  mapState,
  mapGetters
} from 'vuex'
import StatusTag from '../Status/StatusTag'
import dayjs from 'dayjs'

export default {
  name: 'NodeInfoView',
  components: { StatusTag },
  computed: {
    ...mapState('task', [
      'taskForm',
      'taskLog'
    ]),
    ...mapGetters('task', [
      'errorLogData'
    ]),
    isRunning () {
      return ['pending', 'running'].includes(this.taskForm.status)
    }
  },
  methods: {
    onRestart () {
    },
    onStop () {
      this.$store.dispatch('task/cancelTask', this.$route.params.id)
        .then(() => {
          this.$message.success(`Task "${this.$route.params.id}" has been sent signal to stop`)
        })
    },
    getTime (str) {
      if (!str || str.match('^0001')) return 'NA'
      return dayjs(str).format('YYYY-MM-DD HH:mm:ss')
    },
    getWaitDuration (row) {
      if (!row.start_ts || row.start_ts.match('^0001')) return 'NA'
      return dayjs(row.start_ts).diff(row.create_ts, 'second')
    },
    getRuntimeDuration (row) {
      if (!row.finish_ts || row.finish_ts.match('^0001')) return 'NA'
      return dayjs(row.finish_ts).diff(row.start_ts, 'second')
    },
    getTotalDuration (row) {
      if (!row.finish_ts || row.finish_ts.match('^0001')) return 'NA'
      return dayjs(row.finish_ts).diff(row.create_ts, 'second')
    },
    onClickLogWithErrors () {
      this.$emit('click-log')
      this.$st.sendEv('任务详情', '概览', '点击日志错误')
    }
  }
}
</script>

<style scoped>
  .node-form {
    padding: 10px;
  }

  .button-container {
    padding: 0 10px;
    width: 100%;
    text-align: right;
  }

  .el-tag {
    height: 36px;
    line-height: 36px;
  }

  .error-message {
    background-color: rgba(245, 108, 108, .1);
    color: #f56c6c;
    border: 1px solid rgba(245, 108, 108, .2);
    border-radius: 4px;
    line-height: 18px;
    padding: 5px 10px;
  }

  .el-form-item {
    text-align: left;
  }
</style>
