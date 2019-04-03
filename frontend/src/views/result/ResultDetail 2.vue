<template>
  <div class="app-container">
    <!--selector-->
    <div class="selector">
      <label class="label">Spider: </label>
      <el-select v-model="spiderForm._id" @change="onSpiderChange">
        <el-option v-for="op in spiderList" :key="op._id" :value="op._id" :label="op.name"></el-option>
      </el-select>
    </div>

    <!--tabs-->
    <el-tabs v-model="activeTabName" @tab-click="onTabClick" type="card">
      <el-tab-pane label="Overview" name="overview">
        <spider-overview/>
      </el-tab-pane>
      <el-tab-pane label="Files" name="files">
        <file-list/>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import FileList from '../../components/FileList/FileList'
import SpiderOverview from '../../components/Overview/SpiderOverview'

export default {
  name: 'NodeDetail',
  components: {
    FileList,
    SpiderOverview
  },
  data () {
    return {
      activeTabName: 'overview'
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderList',
      'spiderForm'
    ]),
    ...mapState('file', [
      'currentPath'
    ]),
    ...mapState('deploy', [
      'deployList'
    ])
  },
  methods: {
    onTabClick () {
    },
    onSpiderChange (id) {
      this.$router.push(`/spiders/${id}`)
    }
  },
  created () {
    // get the list of the spiders
    this.$store.dispatch('spider/getSpiderList')

    // get spider basic info
    this.$store.dispatch('spider/getSpiderData', this.$route.params.id)
      .then(() => {
        // get spider file info
        this.$store.dispatch('file/getFileList', this.spiderForm.src)
      })

    // get spider deploys
    this.$store.dispatch('spider/getDeployList', this.$route.params.id)

    // get spider tasks
    this.$store.dispatch('spider/getTaskList', this.$route.params.id)
  }
}
</script>

<style scoped>

  .selector {
    display: flex;
    align-items: center;
    position: absolute;
    right: 20px;
    /*float: right;*/
    z-index: 999;
    margin-top: -7px;
  }

  .selector .el-select {
    padding-left: 10px;
  }
</style>
