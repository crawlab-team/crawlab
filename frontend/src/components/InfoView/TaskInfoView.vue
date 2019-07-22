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
        </el-form-item>
        <el-form-item :label="$t('Log File Path')">
          <el-input v-model="taskForm.log_stdout_path" placeholder="Log File Path" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Create Timestamp')">
          <el-input v-model="taskForm.create_ts" placeholder="Create Timestamp" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Finish Timestamp')">
          <el-input v-model="taskForm.finish_ts" placeholder="Finish Timestamp" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Duration (sec)')">
          <el-input v-model="taskForm.duration" placeholder="Duration" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Results Count')">
          <el-input v-model="taskForm.num_results" placeholder="Results Count" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Average Results Count per Second')">
          <el-input v-model="taskForm.avg_num_results" placeholder="Average Results Count per Second" disabled>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('Error Message')" v-if="taskForm.status === 'error'">
          <div class="error-message">
            {{ taskForm.error }}
          </div>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container">
      <el-button v-if="isRunning" type="danger" @click="onStop">{{$t('Stop')}}</el-button>
      <!--<el-button type="danger" @click="onRestart">Restart</el-button>-->
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import StatusTag from '../Status/StatusTag'

export default {
  name: 'NodeInfoView',
  components: { StatusTag },
  computed: {
    ...mapState('task', [
      'taskForm'
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
</style>
