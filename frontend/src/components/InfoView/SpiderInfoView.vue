<template>
  <div class="info-view">
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="spiderForm._id"
      @close="crawlConfirmDialogVisible = false"
    />

    <el-row>
      <el-form label-width="150px"
               :model="spiderForm"
               ref="spiderForm"
               class="spider-form"
               label-position="right">
        <el-form-item :label="$t('Spider ID')">
          <el-input v-model="spiderForm._id" :placeholder="$t('Spider ID')" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Spider Name')">
          <el-input v-model="spiderForm.display_name" :placeholder="$t('Spider Name')" :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Source Folder')">
          <el-input v-model="spiderForm.src" :placeholder="$t('Source Folder')" disabled></el-input>
        </el-form-item>
        <el-form-item v-if="spiderForm.type === 'customized'" :label="$t('Execute Command')" prop="cmd" required :inline-message="true">
          <el-input v-model="spiderForm.cmd" :placeholder="$t('Execute Command')"
                    :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Results Collection')" prop="col" required :inline-message="true">
          <el-input v-model="spiderForm.col" :placeholder="$t('Results Collection')"
                    :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item v-if="false" :label="$t('Site')">
          <el-autocomplete v-model="spiderForm.site"
                           :placeholder="$t('Site')"
                           :fetch-suggestions="fetchSiteSuggestions"
                           clearable
                           :disabled="isView"
                           @select="onSiteSelect">
          </el-autocomplete>
        </el-form-item>
        <el-form-item :label="$t('Spider Type')">
          <el-select v-model="spiderForm.type" :placeholder="$t('Spider Type')" :disabled="true" clearable>
            <el-option value="configurable" :label="$t('Configurable')"></el-option>
            <el-option value="customized" :label="$t('Customized')"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('Remark')">
          <el-input v-model="spiderForm.remark"/>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container" v-if="!isView">
      <el-button size="small" v-if="isShowRun" type="danger" @click="onCrawl">{{$t('Run')}}</el-button>
      <el-button size="small" type="success" @click="onSave">{{$t('Save')}}</el-button>
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import CrawlConfirmDialog from '../Common/CrawlConfirmDialog'

export default {
  name: 'SpiderInfoView',
  components: { CrawlConfirmDialog },
  props: {
    isView: {
      default: false,
      type: Boolean
    }
  },
  data () {
    const cronValidator = (rule, value, callback) => {
      let patArr = []
      for (let i = 0; i < 6; i++) {
        patArr.push('[/*,0-9]+')
      }
      const pat = '^' + patArr.join(' ') + '$'
      if (this.spiderForm.cron_enabled) {
        if (!value) {
          callback(new Error('cron cannot be empty'))
        } else if (!value.match(pat)) {
          callback(new Error('cron format is invalid'))
        }
      }
      callback()
    }
    return {
      crawlConfirmDialogVisible: false,
      cmdRule: [
        { message: 'Execute Command should not be empty', required: true }
      ],
      cronRules: [
        { validator: cronValidator, trigger: 'blur' }
      ]
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm'
    ]),
    isShowRun () {
      return !!this.spiderForm.cmd
    }
  },
  methods: {
    onCrawl () {
      this.crawlConfirmDialogVisible = true
      this.$st.sendEv('爬虫详情-概览', '点击运行')
    },
    onSave () {
      this.$refs['spiderForm'].validate(res => {
        if (res) {
          this.$store.dispatch('spider/editSpider')
            .then(() => {
              this.$message.success(this.$t('Spider info has been saved successfully'))
            })
            .catch(error => {
              this.$message.error(error)
            })
        }
      })
      this.$st.sendEv('爬虫详情-概览', '保存')
    },
    fetchSiteSuggestions (keyword, callback) {
      this.$request.get('/sites', {
        keyword: keyword,
        page_num: 1,
        page_size: 100
      }).then(response => {
        const data = response.data.items.map(d => {
          d.value = `${d.name} | ${d.domain}`
          return d
        })
        callback(data)
      })
    },
    onSiteSelect (item) {
      this.spiderForm.site = item._id
    }
  }
}
</script>

<style scoped>
  .spider-form {
    padding: 10px;
  }

  .button-container {
    padding: 0 10px;
    width: 100%;
    text-align: right;
  }

  .el-autocomplete {
    width: 100%;
  }
</style>
