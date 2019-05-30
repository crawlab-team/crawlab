<template>
  <div class="info-view">
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
          <el-input v-model="spiderForm.name" :placeholder="$t('Spider Name')" :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item v-if="isCustomized" :label="$t('Source Folder')">
          <el-input v-model="spiderForm.src" :placeholder="$t('Source Folder')" disabled></el-input>
        </el-form-item>
        <el-form-item v-if="isCustomized" :label="$t('Execute Command')" prop="cmd" required :inline-message="true">
          <el-input v-model="spiderForm.cmd" :placeholder="$t('Execute Command')"
                    :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Results Collection')">
          <el-input v-model="spiderForm.col" :placeholder="$t('Results Collection')"
                    :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Site')">
          <el-autocomplete v-model="spiderForm.site"
                           :placeholder="$t('Site')"
                           :fetch-suggestions="fetchSiteSuggestions"
                           clearable
                           @select="onSiteSelect">
          </el-autocomplete>
        </el-form-item>
        <el-form-item :label="$t('Spider Type')">
          <el-select v-model="spiderForm.type" :placeholder="$t('Spider Type')" :disabled="true" clearable>
            <el-option value="configurable" :label="$t('Configurable')"></el-option>
            <el-option value="customized" :label="$t('Customized')"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-if="isCustomized" :label="$t('Language')">
          <el-select v-model="spiderForm.lang" :placeholder="$t('Language')" :disabled="isView" clearable>
            <el-option value="python" label="Python"></el-option>
            <el-option value="javascript" label="JavaScript"></el-option>
            <el-option value="java" label="Java"></el-option>
            <el-option value="go" label="Go"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container" v-if="!isView">
      <el-button v-if="isShowRun" type="danger" @click="onCrawl">{{$t('Run')}}</el-button>
      <el-button v-if="isCustomized" type="primary" @click="onDeploy">{{$t('Deploy')}}</el-button>
      <el-button type="success" @click="onSave">{{$t('Save')}}</el-button>
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'SpiderInfoView',
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
      if (this.isCustomized) {
        // customized spider
        if (!this.spiderForm.deploy_ts) {
          return false
        }
        return !!this.spiderForm.cmd
      } else {
        // configurable spider
        return !!this.spiderForm.fields
      }
    },
    isCustomized () {
      return this.spiderForm.type === 'customized'
    }
  },
  methods: {
    onCrawl () {
      const row = this.spiderForm
      this.$refs['spiderForm'].validate(res => {
        if (res) {
          this.$confirm(this.$t('Are you sure to run this spider?'), this.$t('Notification'), {
            confirmButtonText: this.$t('Confirm'),
            cancelButtonText: this.$t('Cancel')
          })
            .then(() => {
              this.$store.dispatch('spider/crawlSpider', row._id)
                .then(() => {
                  this.$message.success(this.$t(`Spider task has been scheduled`))
                })
              this.$st.sendEv('爬虫详情-概览', '运行')
            })
        }
      })
    },
    onDeploy () {
      const row = this.spiderForm

      // save spider
      this.$store.dispatch('spider/editSpider', row._id)

      // validate fields
      this.$refs['spiderForm'].validate(res => {
        if (res) {
          this.$confirm(this.$t('Are you sure to deploy this spider?'), this.$t('Notification'), {
            confirmButtonText: this.$t('Confirm'),
            cancelButtonText: this.$t('Cancel')
          })
            .then(() => {
              this.$store.dispatch('spider/deploySpider', row._id)
                .then(() => {
                  this.$message.success(this.$t(`Spider has been deployed`))
                })
              this.$st.sendEv('爬虫详情-概览', '部署')
            })
        }
      })
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
