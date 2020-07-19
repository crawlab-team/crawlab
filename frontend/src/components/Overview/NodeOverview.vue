<template>
  <el-row>
    <el-col :span="12">
      <!--last tasks-->
      <el-row class="latest-tasks-wrapper">
        <task-table-view :title="$t('Latest Tasks')" />
      </el-row>
    </el-col>

    <el-col :span="12">
      <el-row class="node-info-view-wrapper">
        <!--basic info-->
        <node-info-view />
      </el-row>
    </el-col>
  </el-row>
</template>

<script>
  import {
    mapState
  } from 'vuex'
  import TaskTableView from '../TableView/TaskTableView'
  import NodeInfoView from '../InfoView/NodeInfoView'

  export default {
    name: 'NodeOverview',
    components: {
      NodeInfoView,
      TaskTableView
    },
    computed: {
      id() {
        return this.$route.params.id
      },
      ...mapState('node', [
        'nodeForm'
      ])
    },
    created() {
    },
    methods: {}
  }
</script>

<style scoped>
  .title {
    margin: 10px 0 3px 0;
  }
</style>
