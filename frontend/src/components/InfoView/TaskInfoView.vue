<template>
  <div class="info-view">
    <el-row>
      <el-form label-width="150px"
               :model="taskForm"
               ref="nodeForm"
               class="node-form"
               label-position="right">
        <el-form-item label="Task ID">
          <el-input v-model="taskForm._id" placeholder="Task ID" disabled></el-input>
        </el-form-item>
        <el-form-item label="Status">
          <el-tag type="success" v-if="taskForm.status === 'SUCCESS'">SUCCESS</el-tag>
          <el-tag type="warning" v-else-if="taskForm.status === 'STARTED'">STARTED</el-tag>
          <el-tag type="danger" v-else-if="taskForm.status === 'FAILURE'">FAILURE</el-tag>
          <el-tag type="info" v-else>{{taskForm.status}}</el-tag>
        </el-form-item>
        <!--<el-form-item label="Spider Version">-->
        <!--<el-input v-model="taskForm.spider_version" placeholder="Spider Version" disabled></el-input>-->
        <!--</el-form-item>-->
        <el-form-item label="Log File Path">
          <el-input v-model="taskForm.log_file_path" placeholder="Log File Path" disabled></el-input>
        </el-form-item>
        <el-form-item label="Create Timestamp">
          <el-input v-model="taskForm.create_ts" placeholder="Create Timestamp" disabled></el-input>
        </el-form-item>
        <el-form-item label="Finish Timestamp">
          <el-input v-model="taskForm.finish_ts" placeholder="Finish Timestamp" disabled></el-input>
        </el-form-item>
        <el-form-item label="Duration (sec)">
          <el-input v-model="taskForm.duration" placeholder="Duration" disabled></el-input>
        </el-form-item>
        <el-form-item label="Error Message" v-if="taskForm.status === 'FAILURE'">
          <div class="error-message">
            {{taskForm.result}}
          </div>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container">
      <el-button v-if="isRunning" type="danger" @click="onStop">Stop</el-button>
      <!--<el-button type="danger" @click="onRestart">Restart</el-button>-->
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'NodeInfoView',
  computed: {
    ...mapState('task', [
      'taskForm'
    ]),
    isRunning () {
      return !['SUCCESS', 'FAILURE'].includes(this.taskForm.status)
    }
  },
  methods: {
    onRestart () {
    },
    onStop () {
      this.$store.dispatch('task/stopTask', this.$route.params.id)
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
