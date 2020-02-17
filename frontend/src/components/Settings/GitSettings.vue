<template>
  <div class="git-settings">
    <h3 class="title">{{$t('Git Settings')}}</h3>
    <el-form
      class="git-settings-form"
      label-width="150px"
      :model="spiderForm"
      ref="git-settings-form"
    >
      <el-form-item
        :label="$t('Git URL')"
        prop="git_url"
        required
      >
        <el-input
          v-model="spiderForm.git_url"
          :placeholder="$t('Git URL')"
          @blur="onGitUrlChange"
        >
        </el-input>
      </el-form-item>
      <el-form-item
        :label="$t('Has Credential')"
        prop="git_has_credential"
      >
        <el-switch
          v-model="spiderForm.git_has_credential"
          size="small"
          active-color="#67C23A"
        />
      </el-form-item>
      <el-form-item
        v-if="spiderForm.git_has_credential"
        :label="$t('Git Username')"
        prop="git_username"
      >
        <el-input
          v-model="spiderForm.git_username"
          :placeholder="$t('Git Username')"
        >
        </el-input>
      </el-form-item>
      <el-form-item
        v-if="spiderForm.git_has_credential"
        :label="$t('Git Password')"
        prop="git_password"
      >
        <el-input
          v-model="spiderForm.git_password"
          :placeholder="$t('Git Password')"
          type="password"
        >
        </el-input>
      </el-form-item>
      <el-form-item
        :label="$t('Git Branch')"
        prop="git_branch"
        required
      >
        <el-select
          v-model="spiderForm.git_branch"
          :placeholder="$t('Git Branch')"
          :disabled="!spiderForm.git_url || isGitBranchesLoading"
        >
          <el-option
            v-for="op in gitBranches"
            :key="op"
            :value="op"
            :label="op"
          />
        </el-select>
      </el-form-item>
      <el-form-item
        :label="$t('Auto Sync')"
        prop="git_auto_sync"
      >
        <el-switch
          v-model="spiderForm.git_auto_sync"
          size="small"
          active-color="#67C23A"
        />
      </el-form-item>
      <el-form-item
        v-if="spiderForm.git_auto_sync"
        :label="$t('Sync Frequency')"
        prop="git_sync_frequency"
        required
      >
        <el-select
          v-model="spiderForm.git_sync_frequency"
          :placeholder="$t('Sync Frequency')"
        >
          <el-option
            v-for="op in syncFrequencies"
            :key="op.value"
            :value="op.value"
            :label="op.label"
          />
        </el-select>
      </el-form-item>
    </el-form>
    <div class="action-wrapper">
      <el-button
        size="small"
        type="warning"
        :disabled="isGitResetLoading"
        :icon="isGitResetLoading ? 'el-icon-loading' : 'el-icon-refresh-left'"
        @click="onReset"
      >
        {{$t('Reset')}}
      </el-button>
      <el-button
        size="small"
        type="danger"
        :icon="isGitSyncLoading ? 'el-icon-loading' : 'el-icon-refresh'"
        :disabled="!spiderForm.git_url || isGitSyncLoading"
        @click="onSync"
      >
        {{$t('Sync')}}
      </el-button>
      <el-button size="small" type="success" @click="onSave" icon="el-icon-check">
        {{$t('Save')}}
      </el-button>
    </div>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'GitSettings',
  data () {
    return {
      gitBranches: [],
      isGitBranchesLoading: false,
      isGitSyncLoading: false,
      isGitResetLoading: false,
      syncFrequencies: [
        { label: '1m', value: '0 * * * * *' },
        { label: '5m', value: '0 0/5 * * * *' },
        { label: '15m', value: '0 0/15 * * * *' },
        { label: '30m', value: '0 0/30 * * * *' },
        { label: '1h', value: '0 0 * * * *' },
        { label: '6h', value: '0 0 0/6 * * *' },
        { label: '12h', value: '0 0 0/12 * * *' },
        { label: '1d', value: '0 0 0 0 * *' }
      ]
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm'
    ])
  },
  methods: {
    onSave () {
      this.$refs['git-settings-form'].validate(async valid => {
        if (!valid) return
        const res = await this.$store.dispatch('spider/editSpider')
        if (!res.data.error) {
          this.$message.success(this.$t('Spider info has been saved successfully'))
        }
      })
      this.$st.sendEv('爬虫详情', 'Git 设置', '保存')
    },
    async onGitUrlChange () {
      if (!this.spiderForm.git_url) return
      this.isGitBranchesLoading = true
      const res = await this.$request.get('/git/branches', { url: this.spiderForm.git_url })
      this.gitBranches = res.data.data
      if (!this.spiderForm.git_branch && this.gitBranches.length > 0) {
        this.$set(this.spiderForm, 'git_branch', this.gitBranches[0])
      }
      this.isGitBranchesLoading = false
    },
    async onSync () {
      this.isGitSyncLoading = true
      try {
        const res = await this.$request.post(`/spiders/${this.spiderForm._id}/git/sync`)
        if (!res.data.error) {
          this.$message.success(this.$t('Git has been synchronized successfully'))
        }
      } finally {
        this.isGitSyncLoading = false
      }
    },
    onReset () {
      this.$confirm(
        this.$t('This would delete all files of the spider. Are you sure to continue?'),
        this.$t('Notification'),
        {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        })
        .then(async () => {
          this.isGitResetLoading = true
          try {
            const res = await this.$request.post(`/spiders/${this.spiderForm._id}/git/reset`)
            if (!res.data.error) {
              this.$message.success(this.$t('Git has been reset successfully'))
            }
          } finally {
            this.isGitResetLoading = false
          }
        })
    }
  },
  async created () {
    if (this.spiderForm.git_url) {
      this.onGitUrlChange()
    }
  }
}
</script>

<style scoped>
  .git-settings {
    color: #606266;
  }

  .git-settings .git-settings-form {
    width: 640px;
  }

  .git-settings .title {
    border-bottom: 1px solid #DCDFE6;
    padding-bottom: 15px;
  }

  .git-settings .action-wrapper {
    width: 640px;
    text-align: right;
    margin-top: 10px;
  }
</style>
