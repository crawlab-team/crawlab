<template>
  <div class="info-view">
    <el-row>
      <el-form label-width="150px"
               :model="spiderForm"
               ref="spiderForm"
               class="spider-form"
               label-position="right">
        <el-form-item label="Spider ID">
          <el-input v-model="spiderForm._id.$oid" placeholder="Spider ID" disabled></el-input>
        </el-form-item>
        <el-form-item label="Spider Name">
          <el-input v-model="spiderForm.name" placeholder="Spider Name" :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item label="Source Folder">
          <el-input v-model="spiderForm.src" placeholder="Source Folder" disabled></el-input>
        </el-form-item>
        <el-form-item label="Execute Command" prop="cmd" required :inline-message="true">
          <el-input v-model="spiderForm.cmd" placeholder="Execute Command"
                    :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item label="Results Collection">
          <el-input v-model="spiderForm.col" placeholder="Results Collection"
                    :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item label="Spider Type">
          <el-select v-model="spiderForm.type" placeholder="Select Spider Type" :disabled="isView" clearable>
            <el-option value="scrapy" label="Scrapy"></el-option>
            <el-option value="pyspider" label="PySpider"></el-option>
            <el-option value="webmagic" label="WebMagic"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="Language">
          <el-select v-model="spiderForm.lang" placeholder="Select Language" :disabled="isView" clearable>
            <el-option value="python" label="Python"></el-option>
            <el-option value="javascript" label="JavaScript"></el-option>
            <el-option value="java" label="Java"></el-option>
            <el-option value="go" label="Go"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="Schedule Enabled">
          <el-switch v-model="spiderForm.cron_enabled" :disabled="isView">
          </el-switch>
        </el-form-item>
        <el-form-item label="Schedule Cron" v-if="spiderForm.cron_enabled"
                      prop="cron"
                      :rules="cronRules"
                      :inline-message="true">
          <template slot="label">
            <el-tooltip content="Cron Format: [second] [minute] [hour] [day of month] [month] [day of week]"
                        placement="top">
              <span>
                Schedule Cron
                <i class="fa fa-exclamation-circle"></i>
              </span>
            </el-tooltip>
          </template>
          <el-input v-model="spiderForm.cron" placeholder="Schedule Cron"
                    :disabled="isView"></el-input>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container" v-if="!isView">
      <el-button type="danger" @click="onRun">Run</el-button>
      <el-button type="primary" @click="onDeploy">Deploy</el-button>
      <el-button type="success" @click="onSave">Save</el-button>
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
    ])
  },
  methods: {
    onRun () {
      const row = this.spiderForm
      this.$refs['spiderForm'].validate(res => {
        if (res) {
          this.$confirm('Are you sure to run this spider?', 'Notice', {
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel'
          })
            .then(() => {
              this.$store.dispatch('spider/crawlSpider', row._id.$oid)
                .then(() => {
                  this.$message.success(`Running spider "${row._id.$oid}" has been scheduled`)
                })
            })
        }
      })
    },
    onDeploy () {
      const row = this.spiderForm
      this.$refs['spiderForm'].validate(res => {
        if (res) {
          this.$confirm('Are you sure to deploy this spider?', 'Notice', {
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel'
          })
            .then(() => {
              this.$store.dispatch('spider/crawlSpider', row._id.$oid)
                .then(() => {
                  this.$message.success(`Spider "${row._id.$oid}" has been deployed`)
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
              this.$message.success('Spider info has been saved successfully')
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
