<template>
  <div class="dialog-view">
    <el-dialog
      class="deploy-dialog"
      :title="title"
      :visible.sync="dialogVisible"
      width="40%">
      <!--message-->
      <label>{{message}}</label>

      <!--selection for node-->
      <el-select v-if="type === 'node'" v-model="activeSpider._id.$oid">
        <el-option v-for="op in spiderList" :key="op._id.$oid" :value="op._id.$oid" :label="op.name"></el-option>
      </el-select>

      <!--selection for spider-->
      <el-select v-else-if="type === 'spider'" v-model="activeNode._id">
        <el-option v-for="op in nodeList" :key="op._id" :value="op._id" :label="op.name"></el-option>
      </el-select>

      <!--action buttons-->
      <span slot="footer" class="dialog-footer">
        <el-button @click="onCancel">Cancel</el-button>
        <el-button type="danger" @click="onConfirm">Confirm</el-button>
      </span>
    </el-dialog>

  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'DialogView',
  computed: {
    ...mapState('spider', [
      'spiderList',
      'spiderForm'
    ]),
    ...mapState('node', [
      'nodeList'
    ]),
    ...mapState('dialogView', [
      'dialogType'
    ]),
    type () {
      if (this.dialogType === 'nodeDeploy') {
        return 'node'
      } else if (this.dialogType === 'nodeRun') {
        return 'node'
      } else if (this.dialogType === 'spiderDeploy') {
        return 'spider'
      } else if (this.dialogType === 'spiderRun') {
        return 'spider'
      } else {
        return ''
      }
    },
    activeNode: {
      get () {
        return this.$store.state.spider.activeNode
      },
      set () {
        this.$store.commit('spider/SET_ACTIVE_NODE')
      }
    },
    activeSpider: {
      get () {
        return this.$store.state.node.activeSpider
      },
      set () {
        this.$store.commit('node/SET_ACTIVE_SPIDER')
      }
    },
    dialogVisible: {
      get () {
        return this.$store.state.dialogView.dialogVisible
      },
      set (value) {
        this.$store.commit('dialogView/SET_DIALOG_VISIBLE', value)
      }
    },
    title () {
      if (this.dialogType === 'nodeDeploy') {
        return 'Deploy'
      } else if (this.dialogType === 'nodeRun') {
        return 'Run'
      } else if (this.dialogType === 'spiderDeploy') {
        return 'Deploy'
      } else if (this.dialogType === 'spiderRun') {
        return 'Run'
      } else {
        return ''
      }
    },
    message () {
      if (this.dialogType === 'nodeDeploy') {
        return 'Please select spider you would like to deploy'
      } else if (this.dialogType === 'nodeRun') {
        return 'Please select spider you would like to run'
      } else if (this.dialogType === 'spiderDeploy') {
        return 'Please select node you would like to deploy'
      } else if (this.dialogType === 'spiderRun') {
        return 'Please select node you would like to run'
      } else {
        return ''
      }
    }
  },
  methods: {
    onCancel () {
      this.$store.commit('dialogView/SET_DIALOG_VISIBLE', false)
    },
    onConfirm () {
      if (this.dialogType === 'nodeDeploy') {
      } else if (this.dialogType === 'nodeRun') {
      } else if (this.dialogType === 'spiderDeploy') {
        this.$store.dispatch('spider/deploySpider', {
          id: this.spiderForm._id.$oid,
          nodeId: this.activeNode._id
        })
          .then(() => {
            this.$message.success(`Spider "${this.spiderForm.name}" has been deployed on node "${this.activeNode._id}" successfully`)
          })
          .finally(() => {
            // get spider deploys
            this.$store.dispatch('spider/getDeployList', this.$route.params.id)

            // close dialog
            this.$store.commit('dialogView/SET_DIALOG_VISIBLE', false)
          })
      } else if (this.dialogType === 'spiderRun') {
        this.$store.dispatch('spider/crawlSpider', {
          id: this.spiderForm._id.$oid,
          nodeId: this.activeNode._id
        })
          .then(() => {
            this.$message.success(`Spider "${this.spiderForm.name}" started to run on node "${this.activeNode._id}"`)
          })
          .finally(() => {
            // get spider tasks
            setTimeout(() => {
              this.$store.dispatch('spider/getTaskList', this.$route.params.id)
            }, 500)

            // close dialog
            this.$store.commit('dialogView/SET_DIALOG_VISIBLE', false)
          })
      } else {
      }
    }
  },
  mounted () {
    if (!this.spiderList || !this.spiderList.length) this.$store.dispatch('spider/getSpiderList')
    if (!this.nodeList || !this.nodeList.length) this.$store.dispatch('node/getNodeList')
  }
}
</script>

<style scoped>

</style>
