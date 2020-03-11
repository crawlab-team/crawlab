<template>
  <div class="crawl-confirm-dialog-wrapper">
    <parameters-dialog
      :visible="isParametersVisible"
      :param="form.param"
      @confirm="onParametersConfirm"
      @close="isParametersVisible = false"
    />
    <el-dialog
      :title="$t('Notification')"
      :visible="visible"
      class="crawl-confirm-dialog"
      width="580px"
      :before-close="beforeClose"
    >
      <div style="margin-bottom: 20px;">{{$t('Are you sure to run this spider?')}}</div>
      <el-form label-width="140px" :model="form" ref="form">
        <el-form-item :label="$t('Run Type')" prop="runType" required inline-message>
          <el-select v-model="form.runType" :placeholder="$t('Run Type')">
            <el-option value="all-nodes" :label="$t('All Nodes')"/>
            <el-option value="selected-nodes" :label="$t('Selected Nodes')"/>
            <el-option value="random" :label="$t('Random')"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.runType === 'selected-nodes'" prop="nodeIds" :label="$t('Node')" required
                      inline-message>
          <el-select v-model="form.nodeIds" :placeholder="$t('Node')" multiple clearable>
            <el-option
              v-for="op in nodeList"
              :key="op._id"
              :value="op._id"
              :disabled="isNodeDisabled(op)"
              :label="op.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="spiderForm.is_scrapy && !multiple" :label="$t('Scrapy Spider')" prop="spider" required inline-message>
          <el-select v-model="form.spider" :placeholder="$t('Scrapy Spider')" :disabled="isLoading">
            <el-option
              v-for="s in spiderForm.spider_names"
              :key="s"
              :label="s"
              :value="s"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="spiderForm.is_scrapy || (multiple && scrapySpiders.length > 0)"
          :label="$t('Scrapy Log Level')"
          prop="scrapy_log_level"
          required
          inline-message
        >
          <el-select v-model="form.scrapy_log_level" :placeholder="$t('Scrapy Log Level')">
            <el-option value="INFO" label="INFO"/>
            <el-option value="DEBUG" label="DEBUG"/>
            <el-option value="WARN" label="WARN"/>
            <el-option value="ERROR" label="ERROR"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="spiderForm.type === 'customized'" :label="$t('Parameters')" prop="param" inline-message>
          <template v-if="spiderForm.is_scrapy && !multiple">
            <el-input v-model="form.param" :placeholder="$t('Parameters')" class="param-input"/>
            <el-button type="primary" icon="el-icon-edit" class="param-btn" @click="onOpenParameters"/>
          </template>
          <template v-else>
            <el-input v-model="form.param" :placeholder="$t('Parameters')"></el-input>
          </template>
        </el-form-item>
        <el-form-item class="disclaimer-wrapper">
          <div>
            <el-checkbox v-model="isAllowDisclaimer"/>
            <span style="margin-left: 5px">我已阅读并同意 <a href="javascript:"
                                                      @click="onClickDisclaimer">《免责声明》</a> 所有内容</span>
          </div>
          <div v-if="!spiderForm.is_long_task && !multiple">
            <el-checkbox v-model="isRedirect"/>
            <span style="margin-left: 5px">跳转到任务详情页</span>
          </div>
        </el-form-item>
      </el-form>
      <template slot="footer">
        <el-button type="plain" size="small" @click="$emit('close')">{{$t('Cancel')}}</el-button>
        <el-button type="primary" size="small" @click="onConfirm" :disabled="isConfirmDisabled">
          {{$t('Confirm')}}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import ParametersDialog from './ParametersDialog'

export default {
  name: 'CrawlConfirmDialog',
  components: { ParametersDialog },
  props: {
    spiderId: {
      type: String,
      default: ''
    },
    spiders: {
      type: Array,
      default () {
        return []
      }
    },
    visible: {
      type: Boolean,
      default: false
    },
    multiple: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      form: {
        runType: 'random',
        nodeIds: undefined,
        spider: undefined,
        scrapy_log_level: 'INFO',
        param: '',
        nodeList: []
      },
      isAllowDisclaimer: true,
      isRedirect: true,
      isLoading: false,
      isParametersVisible: false,
      scrapySpidersNamesDict: {}
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm'
    ]),
    ...mapState('setting', [
      'setting'
    ]),
    isConfirmDisabled () {
      if (this.isLoading) return true
      if (!this.isAllowDisclaimer) return true
      return false
    },
    scrapySpiders () {
      return this.spiders.filter(d => d.type === 'customized' && d.is_scrapy)
    }
  },
  watch: {
    visible (value) {
      if (value) {
        this.onOpen()
      }
    }
  },
  methods: {
    beforeClose () {
      this.$emit('close')
    },
    beforeParameterClose () {
      this.isParametersVisible = false
    },
    async fetchScrapySpiderName (id) {
      const res = await this.$request.get(`/spiders/${id}/scrapy/spiders`)
      this.scrapySpidersNamesDict[id] = res.data.data
    },
    onConfirm () {
      this.$refs['form'].validate(async valid => {
        if (!valid) return

        // 请求响应
        let res

        if (!this.multiple) {
          // 运行单个爬虫

          // 参数
          let param = this.form.param

          // Scrapy爬虫特殊处理
          if (this.spiderForm.type === 'customized' && this.spiderForm.is_scrapy) {
            param = `${this.form.spider} --loglevel=${this.form.scrapy_log_level} ${this.form.param}`
          }

          // 发起请求
          res = await this.$store.dispatch('spider/crawlSpider', {
            spiderId: this.spiderId,
            nodeIds: this.form.nodeIds,
            param,
            runType: this.form.runType
          })
        } else {
          // 运行多个爬虫

          // 发起请求
          res = await this.$store.dispatch('spider/crawlSelectedSpiders', {
            nodeIds: this.form.nodeIds,
            runType: this.form.runType,
            taskParams: this.spiders.map(d => {
              // 参数
              let param = this.form.param

              // Scrapy爬虫特殊处理
              if (d.type === 'customized' && d.is_scrapy) {
                param = `${this.scrapySpidersNamesDict[d._id] ? this.scrapySpidersNamesDict[d._id][0] : ''} --loglevel=${this.form.scrapy_log_level} ${this.form.param}`
              }

              return {
                spider_id: d._id,
                param
              }
            })
          })
        }

        // 消息提示
        this.$message.success(this.$t('A task has been scheduled successfully'))

        this.$emit('close')
        if (this.multiple) {
          this.$st.sendEv('爬虫确认', '确认批量运行', this.form.runType)
        } else {
          this.$st.sendEv('爬虫确认', '确认运行', this.form.runType)
        }

        // 是否重定向
        if (
          this.isRedirect &&
          !this.spiderForm.is_long_task &&
          !this.multiple
        ) {
          // 返回任务id
          const id = res.data.data[0]
          this.$router.push('/tasks/' + id)
          this.$st.sendEv('爬虫确认', '跳转到任务详情')
        }

        this.$emit('confirm')
      })
    },
    onClickDisclaimer () {
      this.$router.push('/disclaimer')
    },
    async onOpen () {
      // 节点列表
      this.$request.get('/nodes', {}).then(response => {
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

      // 爬虫列表
      if (!this.multiple) {
        // 单个爬虫
        this.isLoading = true
        try {
          await this.$store.dispatch('spider/getSpiderData', this.spiderId)
          if (this.spiderForm.is_scrapy) {
            await this.$store.dispatch('spider/getSpiderScrapySpiders', this.spiderId)
            if (this.spiderForm.spider_names && this.spiderForm.spider_names.length > 0) {
              this.$set(this.form, 'spider', this.spiderForm.spider_names[0])
            }
          }
        } finally {
          this.isLoading = false
        }
      } else {
        // 多个爬虫
        this.isLoading = true
        try {
          // 遍历 Scrapy 爬虫列表
          await Promise.all(this.scrapySpiders.map(async d => {
            return this.fetchScrapySpiderName(d._id)
          }))
        } finally {
          this.isLoading = false
        }
      }
    },
    onOpenParameters () {
      this.isParametersVisible = true
    },
    onParametersConfirm (value) {
      this.form.param = value
      this.isParametersVisible = false
    },
    isNodeDisabled (node) {
      if (node.status !== 'online') return true
      if (node.is_master && this.setting.run_on_master === 'N') return true
      return false
    }
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

  .crawl-confirm-dialog >>> .param-input {
    width: calc(100% - 56px);
  }

  .crawl-confirm-dialog >>> .param-input .el-input__inner {
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
    border-right: none;
  }

  .crawl-confirm-dialog >>> .param-btn {
    width: 56px;
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
  }

</style>
