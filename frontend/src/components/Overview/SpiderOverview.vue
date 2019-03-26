<template>
  <el-row>
    <el-col :span="12">
      <!--last tasks-->
      <el-row>
        <task-table-view :title="$t('Latest Tasks')"/>
      </el-row>

      <!--last deploys-->
      <el-row v-if="false">
        <deploy-table-view :title="$t('Latest Deploys')"/>
      </el-row>
    </el-col>

    <el-col :span="12">
      <!--basic info-->
      <spider-info-view/>
    </el-col>
  </el-row>
</template>

<script>
import {
  mapState
} from 'vuex'
import DeployTableView from '../TableView/DeployTableView'
import TaskTableView from '../TableView/TaskTableView'
import SpiderInfoView from '../InfoView/SpiderInfoView'

export default {
  name: 'SpiderOverview',
  components: {
    SpiderInfoView,
    DeployTableView,
    TaskTableView
  },
  data () {
    return {
      // spiderForm: {}
    }
  },
  computed: {
    id () {
      return this.$route.params.id
    },
    ...mapState('spider', [
      'spiderForm'
    ]),
    ...mapState('deploy', [
      'deployList'
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
