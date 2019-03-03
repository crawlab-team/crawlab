<template>
  <el-row>
    <el-col :span="12">
      <!--last tasks-->
      <el-row>
        <task-table-view title="Latest Tasks"/>
      </el-row>

      <!--last deploys-->
      <el-row>
        <deploy-table-view title="Latest Deploys"/>
      </el-row>
    </el-col>

    <el-col :span="12">
      <!--basic info-->
      <node-info-view/>
    </el-col>
  </el-row>
</template>

<script>
import {
  mapState
} from 'vuex'
import DeployTableView from '../TableView/DeployTableView'
import TaskTableView from '../TableView/TaskTableView'
import NodeInfoView from '../InfoView/NodeInfoView'

export default {
  name: 'NodeOverview',
  components: {
    NodeInfoView,
    DeployTableView,
    TaskTableView
  },
  computed: {
    id () {
      return this.$route.params.id
    },
    ...mapState('node', [
      'nodeForm'
    ])
  },
  methods: {},
  created () {
  }
}
</script>

<style scoped>
  .title {
    margin: 10px 0 3px 0;
  }
</style>
