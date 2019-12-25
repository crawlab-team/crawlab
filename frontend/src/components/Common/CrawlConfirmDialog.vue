<template>
  <el-dialog
    :title="$t('Notification')"
    :visible="visible"
    class="crawl-confirm-dialog"
    width="480px"
    :before-close="beforeClose"
  >
    <div style="margin-bottom: 20px;">{{$t('Are you sure to run this spider?')}}</div>
    <el-form label-width="80px" :model="form" ref="form">
      <el-form-item :label="$t('Run Type')" prop="runType" required inline-message>
        <el-select v-model="form.runType" :placeholder="$t('Run Type')">
          <el-option value="all-nodes" :label="$t('All Nodes')"/>
          <el-option value="selected-nodes" :label="$t('Selected Nodes')"/>
          <el-option value="random" :label="$t('Random')"/>
        </el-select>
      </el-form-item>
      <el-form-item v-if="form.runType === 'selected-nodes'" prop="nodeIds" :label="$t('Node')" required inline-message>
        <el-select v-model="form.nodeIds" :placeholder="$t('Node')" multiple clearable>
          <el-option
            v-for="op in nodeList"
            :key="op._id"
            :value="op._id"
            :disabled="op.status !== 'online'"
            :label="op.name"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('Parameters')" prop="param" inline-message>
        <el-input v-model="form.param" :placeholder="$t('Parameters')"></el-input>
      </el-form-item>
      <el-form-item class="disclaimer-wrapper">
        <el-checkbox v-model="isAllowDisclaimer"/>
        <span style="margin-left: 5px">我已阅读并同意 <a href="javascript:" @click="onClickDisclaimer">《免责声明》</a> 所有内容</span>
      </el-form-item>
    </el-form>
    <template slot="footer">
      <el-button type="plain" size="small" @click="$emit('close')">{{$t('Cancel')}}</el-button>
      <el-button type="primary" size="small" @click="onConfirm" :disabled="!isAllowDisclaimer">{{$t('Confirm')}}</el-button>
    </template>
  </el-dialog>
</template>

<script>
import request from '../../api/request'

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
      form: {
        runType: 'random',
        nodeIds: undefined,
        param: '',
        nodeList: []
      },
      isAllowDisclaimer: true
    }
  },
  methods: {
    beforeClose () {
      this.$emit('close')
    },
    onConfirm () {
      this.$refs['form'].validate(res => {
        if (!res) return

        this.$store.dispatch('spider/crawlSpider', {
          spiderId: this.spiderId,
          nodeIds: this.form.nodeIds,
          param: this.form.param,
          runType: this.form.runType
        })
          .then(() => {
            this.$message.success(this.$t('A task has been scheduled successfully'))
          })
        this.$emit('close')
        this.$st.sendEv('爬虫确认', '确认运行', this.form.runType)
      })
    },
    onClickDisclaimer () {
      this.$router.push('/disclaimer')
    }
  },
  created () {
    // 节点列表
    request.get('/nodes', {}).then(response => {
      this.nodeList = response.data.data.map(d => {
        d.systemInfo = {
          os: '',
          arch: '',
          num_cpu: '',
          executables: []
        }
        return d
      })
    })
  }
}
</script>

<style scoped>
  .crawl-confirm-dialog >>> .el-form .el-form-item {
    margin-bottom: 20px;
  }

  .crawl-confirm-dialog >>> .disclaimer-wrapper a {
    color: #409eff;
  }
</style>
