<template>
  <div class="navbar">
    <el-dialog
      :visible.sync="isLatestReleaseNoteVisible"
      :title="$t('Release') + ` ${latestRelease.name}`"
    >
      <el-tabs v-model="activeTabName">
        <el-tab-pane :label="$t('Release Note')" name="release-note">
          <div class="content markdown-body" v-html="latestReleaseNoteHtml">
          </div>
          <template slot="footer">
            <el-button type="primary" size="small" @click="isLatestReleaseNoteVisible = false">{{$t('Ok')}}</el-button>
          </template>
        </el-tab-pane>
        <el-tab-pane :label="$t('How to Upgrade')" name="how-to-upgrade">
          <div class="content markdown-body" v-html="howToUpgradeHtml">
          </div>
        </el-tab-pane>
      </el-tabs>
      <template slot="footer">
        <el-button type="primary" size="small" @click="isLatestReleaseNoteVisible = false">{{$t('Ok')}}</el-button>
      </template>
    </el-dialog>

    <hamburger :toggle-click="toggleSideBar" :is-active="sidebar.opened" class="hamburger-container"/>
    <breadcrumb class="breadcrumb"/>
    <el-dropdown class="avatar-container right" trigger="click">
      <span class="el-dropdown-link">
        {{username}}
        <i class="el-icon-arrow-down el-icon--right"></i>
      </span>
      <el-dropdown-menu slot="dropdown" class="user-dropdown">
        <el-dropdown-item>
          <span style="display:block;" @click="logout">{{$t('Logout')}}</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
    <el-dropdown class="lang-list right" trigger="click">
      <span class="el-dropdown-link">
        {{$t($store.getters['lang/lang'])}}
        <i class="el-icon-arrow-down el-icon--right"></i>
      </span>
      <el-dropdown-menu slot="dropdown">
        <el-dropdown-item @click.native="setLang('zh')">
          <span>中文</span>
        </el-dropdown-item>
        <el-dropdown-item @click.native="setLang('en')">
          <span>English</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
    <div class="documentation right">
      <a href="http://docs.crawlab.cn" target="_blank">
        <font-awesome-icon :icon="['far', 'question-circle']"/>
        <span style="margin-left: 5px;">{{$t('Documentation')}}</span>
      </a>
    </div>
    <div v-if="isUpgradable" class="upgrade right" @click="onClickUpgrade">
      <font-awesome-icon :icon="['fas', 'arrow-up']"/>
      <el-badge is-dot>
        <span style="margin-left: 5px;">{{$t('Upgrade')}}</span>
      </el-badge>
    </div>
    <el-popover
      class="wechat right"
      trigger="click"
    >
      <div style="margin-bottom: 5px">
        <label>{{$t('Add Wechat to join discussion group')}}</label>
      </div>
      <div>
        <img class="wechat-img" src="http://static-docs.crawlab.cn/wechat.jpg">
      </div>
      <div slot="reference">
        <i class="fa fa-wechat"></i>
      </div>
    </el-popover>
    <div class="github right">
      <!-- Place this tag where you want the button to render. -->
      <github-button
        href="https://github.com/crawlab-team/crawlab"
        data-color-scheme="no-preference: light; light: light; dark: dark;"
        data-size="large"
        data-show-count="true"
        :aria-label="$t('Star crawlab-team/crawlab on GitHub')"
        style="color: white"
      >
        Star
      </github-button>
    </div>
  </div>
</template>

<script>
import {
  mapState,
  mapGetters
} from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'
import GithubButton from 'vue-github-button'
import showdown from 'showdown'
import 'github-markdown-css/github-markdown.css'

export default {
  components: {
    Breadcrumb,
    Hamburger,
    GithubButton
  },
  data () {
    const converter = new showdown.Converter()
    return {
      isLatestReleaseNoteVisible: false,
      converter,
      activeTabName: 'release-note',
      howToUpgradeHtmlZh: `
### Docker 部署
\`\`\`bash
# 拉取最新镜像
docker pull tikazyq/crawlab:latest

# 删除容器
docker-compose down | true

# 启动容器
docker-compose up -d
\`\`\`

### 直接部署

1. 拉取最新 Github 代码
2. 重新构建前后端应用
3. 启动前后端应用
`,
      howToUpgradeHtmlEn: `
### Docker Deployment
\`\`\`bash
# pull the latest image
docker pull tikazyq/crawlab:latest

# delete containers
docker-compose down | true

# start containers
docker-compose up -d
\`\`\`

### Direct Deployment
1. Pull the latest Github repository
2. Build frontend and backend applications
3. Start frontend and backend applications
`
    }
  },
  computed: {
    ...mapState('version', [
      'latestRelease'
    ]),
    ...mapState('lang', [
      'lang'
    ]),
    ...mapGetters([
      'sidebar',
      'avatar'
    ]),
    username () {
      if (!this.$store.getters['user/userInfo']) return this.$t('User')
      if (!this.$store.getters['user/userInfo'].username) return this.$t('User')
      return this.$store.getters['user/userInfo'].username
    },
    isUpgradable () {
      if (!this.latestRelease.name) return false

      const currentVersion = sessionStorage.getItem('v')
      const latestVersion = this.latestRelease.name.replace('v', '')

      if (!latestVersion || !currentVersion) return false

      const currentVersionList = currentVersion.split('.')
      const latestVersionList = latestVersion.split('.')
      for (let i = 0; i < currentVersionList.length; i++) {
        let nc = Number(currentVersionList[i])
        let nl = Number(latestVersionList[i])
        if (isNaN(nl)) nl = 0
        if (nc < nl) return true
      }
      return false
    },
    latestReleaseNoteHtml () {
      if (!this.latestRelease.body) return ''
      const body = this.latestRelease.body
      return this.converter.makeHtml(body)
    },
    howToUpgradeHtml () {
      if (this.lang === 'zh') {
        return this.converter.makeHtml(this.howToUpgradeHtmlZh)
      } else if (this.lang === 'en') {
        return this.converter.makeHtml(this.howToUpgradeHtmlEn)
      } else {
        return ''
      }
    }
  },
  methods: {
    toggleSideBar () {
      this.$store.dispatch('ToggleSideBar')
    },
    logout () {
      this.$store.dispatch('user/logout')
      this.$store.dispatch('delAllViews')
      this.$router.push('/login')
      this.$st.sendEv('全局', '登出')
    },
    setLang (lang) {
      window.localStorage.setItem('lang', lang)
      this.$i18n.locale = lang
      this.$store.commit('lang/SET_LANG', lang)

      this.$st.sendEv('全局', '切换中英文', lang)
    },
    onClickUpgrade () {
      this.isLatestReleaseNoteVisible = true
      this.$st.sendEv('全局', '点击版本升级')
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
  .navbar {
    height: 50px;
    line-height: 50px;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, .12), 0 0 3px 0 rgba(0, 0, 0, .04);

    .hamburger-container {
      line-height: 58px;
      height: 50px;
      float: left;
      padding: 0 10px;
    }

    .lang-list {
      cursor: pointer;
      display: inline-block;
      margin-right: 35px;
      /*position: absolute;*/
      /*right: 80px;*/
    }

    .screenfull {
      position: absolute;
      right: 90px;
      top: 16px;
      color: red;
    }

    .avatar-container {
      cursor: pointer;
      height: 50px;
      display: inline-block;
      margin-right: 35px;
      /*position: absolute;*/
      /*right: 35px;*/
    }

    .documentation {
      margin-right: 35px;
      color: #606266;
      font-size: 14px;

      .span {
        margin-left: 5px;
      }
    }

    .github {
      height: 50px;
      margin-right: 35px;
      margin-top: -10px;
    }

    .upgrade {
      margin-top: 12.5px;
      line-height: 25px;
      cursor: pointer;
      font-size: 14px;
      color: #606266;
      margin-right: 35px;

      .span {
        margin-left: 5px;
      }
    }

    .wechat {
      color: #606266;
      margin-right: 35px;
      cursor: pointer;
    }

    .right {
      float: right
    }
  }
</style>
<style scoped>
  .navbar >>> .el-dialog__body {
    padding-top: 0;
  }
</style>
<style>
  .wechat-img {
    width: 240px;
  }
</style>
