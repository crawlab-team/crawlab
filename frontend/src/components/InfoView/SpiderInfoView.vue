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
        <el-form-item v-if="spiderForm.type === 'customized'" :label="$t('Execute Command')" prop="cmd" required
                      :inline-message="true">
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
          <el-input type="textarea" v-model="spiderForm.remark" :placeholder="$t('Remark')"/>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container" v-if="!isView">
      <el-upload
        v-if="spiderForm.type === 'customized'"
        :action="$request.baseUrl + `/spiders/${spiderForm._id}/upload`"
        :headers="{Authorization:token}"
        :on-progress="() => this.uploadLoading = true"
        :on-error="onUploadError"
        :on-success="onUploadSuccess"
        :file-list="fileList"
        style="display:inline-block;margin-right:10px"
      >
        <el-button size="normal" type="primary" icon="el-icon-upload" v-loading="uploadLoading">
          {{$t('Upload')}}
        </el-button>
      </el-upload>
      <el-button size="normal" v-if="isShowRun" type="danger" @click="onCrawl"
                 icon="el-icon-video-play">
        {{$t('Run')}}
      </el-button>
      <el-button size="normal" type="success" @click="onSave"
                 icon="el-icon-check">
        {{$t('Save')}}
      </el-button>
    </el-row>
  </div>
</template>

<script>
import {
  mapState,
  mapGetters
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
      uploadLoading: false,
      fileList: [],
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
    ...mapGetters('user', [
      'token'
    ]),
    isShowRun () {
      if (this.spiderForm.type === 'customized') {
        return !!this.spiderForm.cmd
      } else {
        return true
      }
    }
  },
  methods: {
    onCrawl () {
      this.crawlConfirmDialogVisible = true
      this.$st.sendEv('爬虫详情', '概览', '点击运行')
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
      this.$st.sendEv('爬虫详情', '概览', '保存')
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
    },
    onUploadSuccess () {
      this.$store.dispatch('file/getFileList', this.spiderForm.src)

      this.uploadLoading = false

      this.$message.success(this.$t('Uploaded spider files successfully'))
    },
    onUploadError () {
      this.uploadLoading = false
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

  .info-view >>> .el-upload-list {
    display: none;
  }
</style>
