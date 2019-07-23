<template>
  <el-dialog
    :title="$t('Notification')"
    :visible="visible"
    width="480px"
    :before-close="beforeClose"
  >
    <div style="margin-bottom: 20px;">{{$t('Are you sure to run this spider?')}}</div>
    <el-form label-width="80px">
      <el-form-item :label="$t('Node')">
        <el-select v-model="nodeId">
          <el-option value="" :label="$t('All Nodes')"/>
          <el-option
            v-for="op in $store.state.node.nodeList"
            :key="op._id"
            :value="op._id"
            :disabled="op.status !== 'online'"
            :label="op.name"
          />
        </el-select>
      </el-form-item>
    </el-form>
    <template slot="footer">
      <el-button type="plain" size="small" @click="$emit('close')">{{$t('Cancel')}}</el-button>
      <el-button type="primary" size="small" @click="onConfirm">{{$t('Confirm')}}</el-button>
    </template>
  </el-dialog>
</template>

<script>
export default {
  name: 'CrawlConfirmDialog',
  props: {
    spiderId: {
      type: String,
      default: ''
    },
    visible: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      nodeId: ''
    }
  },
  methods: {
    beforeClose () {
      this.$emit('close')
    },
    onConfirm () {
      this.$store.dispatch('spider/crawlSpider', { id: this.spiderId, nodeId: this.nodeId })
        .then(() => {
          this.$message.success(this.$t('A task has been scheduled successfully'))
        })
      this.$emit('close')
      this.$st.sendEv('爬虫', '运行')
    }
  }
}
</script>

<style scoped>

</style>
