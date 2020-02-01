<template>
  <div class="app-container">
    <!--tour-->
    <v-tour
      name="node-detail"
      :steps="tourSteps"
      :callbacks="tourCallbacks"
      :options="$utils.tour.getOptions(true)"
    />
    <!--./tour-->

    <!--selector-->
    <div class="selector">
      <label class="label">{{$t('Node')}}: </label>
      <el-select size="small" v-model="nodeForm._id" @change="onNodeChange">
        <el-option v-for="op in nodeList" :key="op._id" :value="op._id" :label="op.name"></el-option>
      </el-select>
    </div>

    <!--tabs-->
    <el-tabs v-model="activeTabName" @tab-click="onTabClick" type="card">
      <el-tab-pane :label="$t('Overview')" name="overview">
        <node-overview></node-overview>
      </el-tab-pane>
      <el-tab-pane :label="$t('Installation')" name="installation">
        <node-installation></node-installation>
      </el-tab-pane>
      <el-tab-pane :label="$t('Deployed Spiders')" name="spiders" v-if="false">
        {{$t('Deployed Spiders')}}
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import NodeOverview from '../../components/Overview/NodeOverview'
import NodeInstallation from '../../components/Node/NodeInstallation'

export default {
  name: 'NodeDetail',
  components: {
    NodeOverview,
    NodeInstallation
  },
  data () {
    return {
      activeTabName: 'overview',
      tourSteps: [
        // overview
        {
          target: '.selector',
          content: this.$t('Switch between different nodes.')
        },
        {
          target: '.latest-tasks-wrapper',
          content: this.$t('You can view the latest executed spider tasks.'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '.node-info-view-wrapper',
          content: this.$t('This is the detailed node info.'),
          params: {
            placement: 'left'
          }
        },
        // installation
        {
          target: '#tab-installation',
          content: this.$t('Here you can install<br> dependencies and modules<br> that are required<br> in your spiders.')
        },
        {
          target: '.search-box',
          content: this.$t('You can search dependencies in the search box and install them by clicking the "Install" button below.')
        }
      ],
      tourCallbacks: {
        onStop: () => {
          this.$utils.tour.finishTour('node-detail')
        },
        onPreviousStep: (currentStep) => {
          if (currentStep === 3) {
            this.activeTabName = 'overview'
          }
        },
        onNextStep: (currentStep) => {
          if (currentStep === 2) {
            this.activeTabName = 'installation'
          }
        }
      }
    }
  },
  computed: {
    ...mapState('node', [
      'nodeList',
      'nodeForm'
    ])
  },
  methods: {
    onTabClick (tab) {
      this.$st.sendEv('节点详情', '切换标签', tab.name)
    },
    onNodeChange (id) {
      this.$router.push(`/nodes/${id}`)
      this.$st.sendEv('节点详情', '切换节点')
    }
  },
  created () {
    // get list of nodes
    this.$store.dispatch('node/getNodeList')

    // get node basic info
    this.$store.dispatch('node/getNodeData', this.$route.params.id)

    // get node task list
    this.$store.dispatch('node/getTaskList', this.$route.params.id)
  },
  mounted () {
    if (!this.$utils.tour.isFinishedTour('node-detail')) {
      this.$tours['node-detail'].start()
    }
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

  .label {
    width: 100px;
    text-align: right;
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
