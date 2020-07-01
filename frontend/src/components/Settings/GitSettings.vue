<template>
  <el-tabs
    v-model="activeTabName"
    class="git-settings"
    @change="onChangeTab"
  >
    <el-tab-pane :label="$t('Settings')" name="settings">
      <el-form
        ref="git-settings-form"
        class="git-settings-form"
        label-width="150px"
        :model="spiderForm"
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
          />
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
            @blur="onGitUrlChange"
          />
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
            @blur="onGitUrlChange"
          />
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
        <el-form-item
          v-if="spiderForm.git_sync_error"
          :label="$t('Error Message')"
          prop="git_git_sync_error"
        >
          <el-alert
            type="error"
            :closable="false"
          >
            {{ spiderForm.git_sync_error }}
          </el-alert>
        </el-form-item>
        <el-form-item
          v-if="sshPublicKey"
          :label="$t('SSH Public Key')"
        >
          <el-alert
            type="info"
            :closable="false"
          >
            {{ sshPublicKey }}
          </el-alert>
          <span class="copy" @click="copySshPublicKey">
            <i class="el-icon-copy-document" />
            {{ $t('Copy') }}
          </span>
          <input v-show="true" id="ssh-public-key" v-model="sshPublicKey">
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
          {{ $t('Reset') }}
        </el-button>
        <el-button
          size="small"
          type="danger"
          :icon="isGitSyncLoading ? 'el-icon-loading' : 'el-icon-refresh'"
          :disabled="!spiderForm.git_url || isGitSyncLoading"
          @click="onSync"
        >
          {{ $t('Sync') }}
        </el-button>
        <el-button size="small" type="success" icon="el-icon-check" @click="onSave">
          {{ $t('Save') }}
        </el-button>
      </div>
    </el-tab-pane>
    <el-tab-pane label="Log" name="log">
      <el-timeline
        class="log"
      >
        <el-timeline-item
          v-for="c in commits"
          :key="c.hash"
          :timestamp="c.ts"
          :type="getCommitType(c)"
        >
          <div class="commit">
            <div class="row">
              <div class="message">
                {{ c.message }}
              </div>
              <div class="author">
                {{ c.author }} ({{ c.email }})
              </div>
            </div>
            <div class="row" style="margin-top: 10px">
              <div class="tags">
                <el-tag
                  v-if="c.is_head"
                  type="primary"
                  size="mini"
                >
                  <i class="fa fa-tag" />
                  HEAD
                </el-tag>
                <el-tag
                  v-for="b in c.branches"
                  :key="b.name"
                  :type="b.label === 'master' ? 'danger' : 'warning'"
                  size="mini"
                >
                  <i class="fa fa-tag" />
                  {{ b.label }}
                </el-tag>
                <el-tag
                  v-for="b in c.remote_branches"
                  :key="b.name"
                  type="info"
                  size="mini"
                >
                  <i class="fa fa-tag" />
                  {{ b.label }}
                </el-tag>
                <el-tag
                  v-for="t in c.tags"
                  :key="t.name"
                  type="success"
                  size="mini"
                >
                  <i class="fa fa-tag" />
                  {{ t.label }}
                </el-tag>
              </div>
              <div class="actions">
                <el-button
                  v-if="!c.is_head"
                  type="danger"
                  :icon="isGitCheckoutLoading ? 'el-icon-loading' : 'el-icon-position'"
                  size="mini"
                  @click="checkout(c)"
                >
                  Checkout
                </el-button>
              </div>
            </div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </el-tab-pane>
  </el-tabs>
</template>

<script>
  import dayjs from 'dayjs'
  import {
    mapState
  } from 'vuex'

  export default {
    name: 'GitSettings',
    data() {
      return {
        gitBranches: [],
        isGitBranchesLoading: false,
        isGitSyncLoading: false,
        isGitResetLoading: false,
        isGitCheckoutLoading: false,
        syncFrequencies: [
          { label: '1m', value: '0 * * * * *' },
          { label: '5m', value: '0 0/5 * * * *' },
          { label: '15m', value: '0 0/15 * * * *' },
          { label: '30m', value: '0 0/30 * * * *' },
          { label: '1h', value: '0 0 * * * *' },
          { label: '6h', value: '0 0 0/6 * * *' },
          { label: '12h', value: '0 0 0/12 * * *' },
          { label: '1d', value: '0 0 0 0 * *' }
        ],
        sshPublicKey: '',
        activeTabName: 'settings',
        commits: []
      }
    },
    computed: {
      ...mapState('spider', [
        'spiderForm'
      ])
    },
    async created() {
      if (this.spiderForm.git_url) {
        this.onGitUrlChange()
      }
      this.getSshPublicKey()
      this.getCommits()
    },
    methods: {
      onSave() {
        this.$refs['git-settings-form'].validate(async valid => {
          if (!valid) return
          const res = await this.$store.dispatch('spider/editSpider')
          if (!res.data.error) {
            this.$message.success(this.$t('Spider info has been saved successfully'))
          }
        })
        this.$st.sendEv('爬虫详情', 'Git 设置', '保存')
      },
      async onGitUrlChange() {
        if (!this.spiderForm.git_url) return
        this.isGitBranchesLoading = true
        try {
          const res = await this.$request.get('/git/branches', { url: this.spiderForm.git_url })
          this.gitBranches = res.data.data
          if (!this.spiderForm.git_branch && this.gitBranches.length > 0) {
            this.$set(this.spiderForm, 'git_branch', this.gitBranches[0])
          }
        } finally {
          this.isGitBranchesLoading = false
        }
      },
      async onSync() {
        this.isGitSyncLoading = true
        try {
          const res = await this.$request.post(`/spiders/${this.spiderForm._id}/git/sync`)
          if (!res.data.error) {
            this.$message.success(this.$t('Git has been synchronized successfully'))
          }
        } finally {
          this.isGitSyncLoading = false
          await this.updateGit()
          await this.$store.dispatch('spider/getSpiderData', this.$route.params.id)
        }
        this.$st.sendEv('爬虫详情', 'Git 设置', '同步')
      },
      onReset() {
        this.$confirm(
          this.$t('This would delete all files of the spider. Are you sure to continue?'),
          this.$t('Notification'),
          {
            confirmButtonText: this.$t('Confirm'),
            cancelButtonText: this.$t('Cancel'),
            type: 'warning'
          })
          .then(async() => {
            this.isGitResetLoading = true
            try {
              const res = await this.$request.post(`/spiders/${this.spiderForm._id}/git/reset`)
              if (!res.data.error) {
                this.$message.success(this.$t('Git has been reset successfully'))
                this.$st.sendEv('爬虫详情', 'Git 设置', '确认重置')
              }
            } finally {
              this.isGitResetLoading = false
            // await this.updateGit()
            }
          })
        this.$st.sendEv('爬虫详情', 'Git 设置', '点击重置')
      },
      async getSshPublicKey() {
        const res = await this.$request.get('/git/public-key')
        this.sshPublicKey = res.data.data
      },
      copySshPublicKey() {
        const el = document.querySelector('#ssh-public-key')
        el.focus()
        el.setSelectionRange(0, this.sshPublicKey.length)
        document.execCommand('copy')
        this.$message.success(this.$t('SSH Public Key is copied to the clipboard'))
        this.$st.sendEv('爬虫详情', 'Git 设置', '拷贝 SSH 公钥')
      },
      async getCommits() {
        const res = await this.$request.get('/git/commits', { spider_id: this.spiderForm._id })
        this.commits = res.data.data.map(d => {
          d.ts = dayjs(d.ts).format('YYYY-MM-DD HH:mm:ss')
          return d
        })
      },
      async checkout(c) {
        this.isGitCheckoutLoading = true
        try {
          const res = await this.$request.post('/git/checkout', { spider_id: this.spiderForm._id, hash: c.hash })
          if (!res.data.error) {
            this.$message.success(this.$t('Checkout success'))
          }
        } finally {
          this.isGitCheckoutLoading = false
          await this.getCommits()
        }
        this.$st.sendEv('爬虫详情', 'Git Log', 'Checkout')
      },
      async updateGit() {
        this.getCommits()
      },
      getCommitType(c) {
        if (c.is_head) return 'primary'
        if (c.branches && c.branches.length) {
          if (c.branches.map(d => d.label).includes('master')) {
            return 'danger'
          } else {
            return 'warning'
          }
        }
        if (c.tags && c.tags.length) {
          return 'success'
        }
        if (c.remote_branches && c.remote_branches.length) {
          return 'info'
        }
      },
      onChangeTab() {
        this.$st.sendEv('爬虫详情', 'Git 切换标签', this.activeTabName)
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

  .git-settings .git-settings-form >>> .el-alert {
    padding: 0 5px;
    margin: 0;
  }

  .git-settings .git-settings-form >>> .el-alert__description {
    padding: 0;
    margin: 0;
    font-size: 14px;
    line-height: 24px;
  }

  .git-settings .git-settings-form .copy {
    display: inline;
    line-height: 14px;
    position: absolute;
    top: 5px;
    right: 5px;
    cursor: pointer;
  }

  .git-settings .git-settings-form .copy {
  }

  #ssh-public-key {
    position: absolute;
    z-index: -1;
    top: 0;
    left: 0;
    height: 0;
    /*visibility: hidden;*/
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

  .git-settings .log {
    height: calc(100vh - 280px);
    overflow: auto;
  }

  .git-settings .log .commit {
    border-top: 1px solid rgb(244, 244, 245);
    padding: 10px 0;
  }

  .git-settings .log .commit .row {
    display: flex;
    justify-content: space-between;
  }

  .git-settings .log .el-timeline-item {
    /*cursor: pointer;*/
  }

  .git-settings .log .commit .row .tags .el-tag {
    margin-right: 5px;
  }

  .git-settings .log .commit .row .actions {
    right: 0;
    bottom: 5px;
    position: absolute;
  }
</style>
