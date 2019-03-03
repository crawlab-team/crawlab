<template>
  <div class="app-container">
    <!--selector-->
    <div class="selector">
      <label class="label">Node: </label>
      <el-select v-model="nodeForm._id" @change="onNodeChange">
        <el-option v-for="op in nodeList" :key="op._id" :value="op._id" :label="op.name"></el-option>
      </el-select>
    </div>

    <!--tabs-->
    <el-tabs v-model="activeTabName" @tab-click="onTabClick" type="card">
      <el-tab-pane label="Overview" name="overview">
        <node-overview></node-overview>
      </el-tab-pane>
      <el-tab-pane label="Deployed Spiders" name="spiders">
        Deployed Spiders
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import NodeOverview from '../../components/Overview/NodeOverview'

export default {
  name: 'NodeDetail',
  components: {
    NodeOverview
  },
  data () {
    return {
      activeTabName: 'overview'
    }
  },
  computed: {
    ...mapState('node', [
      'nodeList',
      'nodeForm'
    ])
  },
  methods: {
    onTabClick () {
    },
    onNodeChange (id) {
      this.$router.push(`/nodes/${id}`)
    }
  },
  created () {
    // get list of nodes
    this.$store.dispatch('node/getNodeList')

    // get node basic info
    this.$store.dispatch('node/getNodeData', this.$route.params.id)

    // get node deploy list
    this.$store.dispatch('node/getDeployList', this.$route.params.id)

    // get node task list
    this.$store.dispatch('node/getTaskList', this.$route.params.id)
  }
}
</script>

<style scoped>
  .selector {
    display: flex;
    align-items: center;
    position: absolute;
    right: 20px;
    margin-top: -7px;
    /*float: right;*/
    z-index: 999;
  }

  .selector .el-select {
    padding-left: 10px;
  }
</style>
<style lang="scss">
  .selector {
    .el-select {
      .el-input {
        .el-input_inner {
          height: 26px;
        }
      }
    }
  }
</style>
