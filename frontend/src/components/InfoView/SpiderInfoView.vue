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
        <el-form-item :label="$t('Source Folder')">
          <el-input v-model="spiderForm.src" :placeholder="$t('Source Folder')" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Execute Command')" prop="cmd" required :inline-message="true">
          <el-input v-model="spiderForm.cmd" :placeholder="$t('Execute Command')"
                    :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Results Collection')">
          <el-input v-model="spiderForm.col" :placeholder="$t('Results Collection')"
                    :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Spider Type')">
          <el-select v-model="spiderForm.type" :placeholder="$t('Spider Type')" :disabled="isView" clearable>
            <el-option value="scrapy" label="Scrapy"></el-option>
            <el-option value="pyspider" label="PySpider"></el-option>
            <el-option value="webmagic" label="WebMagic"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('Language')">
          <el-select v-model="spiderForm.lang" :placeholder="$t('Language')" :disabled="isView" clearable>
            <el-option value="python" label="Python"></el-option>
            <el-option value="javascript" label="JavaScript"></el-option>
            <el-option value="java" label="Java"></el-option>
            <el-option value="go" label="Go"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('Schedule Enabled')">
          <el-switch v-model="spiderForm.cron_enabled" :disabled="isView">
          </el-switch>
        </el-form-item>
        <el-form-item :label="$t('Schedule Cron')" v-if="spiderForm.cron_enabled"
                      prop="cron"
                      :rules="cronRules"
                      :inline-message="true">
          <template slot="label">
            <el-tooltip :content="$t('Cron Format: [second] [minute] [hour] [day of month] [month] [day of week]')"
                        placement="top">
              <span>
                {{$t('Schedule Cron')}}
                <i class="fa fa-exclamation-circle"></i>
              </span>
            </el-tooltip>
          </template>
          <el-input v-model="spiderForm.cron" :placeholder="$t('Schedule Cron')"
                    :disabled="isView"></el-input>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container" v-if="!isView">
      <el-button v-if="isShowRun" type="danger" @click="onRun">{{$t('Run')}}</el-button>
      <el-button type="primary" @click="onDeploy">{{$t('Deploy')}}</el-button>
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
      if (!this.spiderForm.deploy_ts) {
        return false
      }
      if (!this.spiderForm.cmd) {
        return false
      }
      return true
    }
  },
  methods: {
    onRun () {
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
</style>
