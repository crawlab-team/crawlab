<template>
  <div class="node-installation">
    <el-form inline>
      <el-form-item>
        <el-autocomplete size="small" clearable @clear="onSearch"
          v-if="activeLang.executable_name === 'python'"
          v-model="depName"
          style="width: 240px"
          :placeholder="$t('Search Dependencies')"
          :fetchSuggestions="fetchAllDepList"
          :minlength="2"
          @select="onSearch"
        ></el-autocomplete>
        <el-input
          v-else
          v-model="depName"
          style="width: 240px"
          :placeholder="$t('Search Dependencies')"
        />
      </el-form-item>
      <el-form-item>
        <el-button size="small"
          icon="el-icon-search"
          type="success"
          @click="onSearch"
        >
          {{$t('Search')}}
        </el-button>
      </el-form-item>
      <el-form-item>
        <el-checkbox v-model="isShowInstalled" :label="$t('Show installed')" @change="onIsShowInstalledChange"/>
      </el-form-item>
    </el-form>
    <el-tabs v-model="activeTab" @tab-click="onTabChange">
      <el-tab-pane v-for="lang in langList" :key="lang.name" :label="lang.name" :name="lang.executable_name"/>
    </el-tabs>
    <template v-if="activeLang.installed">
      <el-table
        height="calc(100vh - 280px)"
        :data="computedDepList"
        :empty-text="depName ? $t('No Data') : $t('Please search dependencies')"
        v-loading="loading"
        border
      >
        <el-table-column
          :label="$t('Name')"
          prop="name"
          width="180"
        />
        <el-table-column
          :label="!isShowInstalled ? $t('Latest Version') : $t('Version')"
          prop="version"
          width="100"
        />
        <el-table-column
          v-if="!isShowInstalled"
          :label="$t('Description')"
          prop="description"
        />
        <el-table-column
          :label="$t('Action')"
        >
          <template slot-scope="scope">
            <el-button
              v-if="!scope.row.installed"
              v-loading="getDepLoading(scope.row)"
              :disabled="getDepLoading(scope.row)"
              size="mini"
              type="primary"
              @click="onClickInstallDep(scope.row)"
            >
              {{$t('Install')}}
            </el-button>
            <el-button
              v-else
              v-loading="getDepLoading(scope.row)"
              :disabled="getDepLoading(scope.row)"
              size="mini"
              type="danger"
              @click="onClickUninstallDep(scope.row)"
            >
              {{$t('Uninstall')}}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>
    <template v-else>
      <div class="install-wrapper">
        <h3>{{activeLang.name + $t(' is not installed, do you want to install it?')}}</h3>
        <el-button
          v-loading="isLoadingInstallLang"
          :disabled="isLoadingInstallLang"
          type="primary"
          style="width: 240px;font-weight: bolder;font-size: 18px"
          @click="onClickInstallLang"
        >
          {{$t('Install')}}
        </el-button>
      </div>
    </template>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'NodeInstallation',
  data () {
    return {
      activeTab: '',
      langList: [],
      depName: '',
      depList: [],
      loading: false,
      isShowInstalled: false,
      installedDepList: [],
      depLoadingDict: {},
      isLoadingInstallLang: false
    }
  },
  computed: {
    ...mapState('node', [
      'nodeForm'
    ]),
    activeLang () {
      for (let i = 0; i < this.langList.length; i++) {
        if (this.langList[i].executable_name === this.activeTab) {
          return this.langList[i]
        }
      }
      return {}
    },
    activeLangName () {
      return this.activeLang.executable_name
    },
    computedDepList () {
      if (this.isShowInstalled) {
        return this.installedDepList
      } else {
        return this.depList
      }
    }
  },
  methods: {
    async getDepList () {
      this.loading = true
      const res = await this.$request.get(`/nodes/${this.nodeForm._id}/deps`, {
        lang: this.activeLang.executable_name,
        dep_name: this.depName
      })
      this.loading = false
      this.depList = res.data.data

      if (this.activeLangName === 'python') {
        // 排序
        this.depList = this.depList.sort((a, b) => a.name > b.name ? 1 : -1)

        // 异步获取python附加信息
        this.depList.map(async dep => {
          const resp = await this.$request.get(`/system/deps/${this.activeLang.executable_name}/${dep.name}/json`)
          if (resp) {
            dep.version = resp.data.data.version
            dep.description = resp.data.data.description
          }
        })
      }
    },
    async getInstalledDepList () {
      this.loading = true
      const res = await this.$request.get(`/nodes/${this.nodeForm._id}/deps/installed`, {
        lang: this.activeLang.executable_name
      })
      this.loading = false
      this.installedDepList = res.data.data
    },
    async fetchAllDepList (queryString, callback) {
      const res = await this.$request.get(`/system/deps/${this.activeLang.executable_name}`, {
        dep_name: queryString
      })
      callback(res.data.data ? res.data.data.map(d => {
        return { value: d, label: d }
      }) : [])
    },
    onSearch () {
      this.isShowInstalled = false
      this.getDepList()
      this.$st.sendEv('节点详情', '安装', '搜索依赖')
    },
    onIsShowInstalledChange (val) {
      if (val) {
        this.getInstalledDepList()
      } else {
        this.depName = ''
        this.depList = []
      }
      this.$st.sendEv('节点详情', '安装', '点击查看已安装')
    },
    async onClickInstallDep (dep) {
      const name = dep.name
      this.$set(this.depLoadingDict, name, true)
      const arr = this.$route.path.split('/')
      const id = arr[arr.length - 1]
      const data = await this.$request.post(`/nodes/${id}/deps/install`, {
        lang: this.activeLang.executable_name,
        dep_name: name
      })
      if (!data || data.error) {
        this.$notify.error({
          title: this.$t('Installing dependency failed'),
          message: this.$t('The dependency installation is unsuccessful: ') + name
        })
      } else {
        this.$notify.success({
          title: this.$t('Installing dependency successful'),
          message: this.$t('You have successfully installed a dependency: ') + name
        })
        dep.installed = true
      }
      this.$set(this.depLoadingDict, name, false)
      this.$st.sendEv('节点详情', '安装', '安装依赖')
    },
    async onClickUninstallDep (dep) {
      const name = dep.name
      this.$set(this.depLoadingDict, name, true)
      const arr = this.$route.path.split('/')
      const id = arr[arr.length - 1]
      const data = await this.$request.post(`/nodes/${id}/deps/uninstall`, {
        lang: this.activeLang.executable_name,
        dep_name: name
      })
      if (!data || data.error) {
        this.$notify.error({
          title: this.$t('Uninstalling dependency failed'),
          message: this.$t('The dependency uninstallation is unsuccessful: ') + name
        })
      } else {
        this.$notify.success({
          title: this.$t('Uninstalling dependency successful'),
          message: this.$t('You have successfully uninstalled a dependency: ') + name
        })
        dep.installed = false
      }
      this.$set(this.depLoadingDict, name, false)
      this.$st.sendEv('节点详情', '安装', '卸载依赖')
    },
    getDepLoading (dep) {
      const name = dep.name
      if (this.depLoadingDict[name] === undefined) {
        return false
      }
      return this.depLoadingDict[name]
    },
    async onClickInstallLang () {
      this.isLoadingInstallLang = true
      const res = await this.$request.post(`/nodes/${this.nodeForm._id}/langs/install`, {
        lang: this.activeLang.executable_name
      })
      if (!res || res.error) {
        this.$notify.error({
          title: this.$t('Installing language failed'),
          message: this.$t('The language installation is unsuccessful: ') + this.activeLang.name
        })
      } else {
        this.$notify.success({
          title: this.$t('Installing language successful'),
          message: this.$t('You have successfully installed a language: ') + this.activeLang.name
        })
      }
      this.isLoadingInstallLang = false
      this.$st.sendEv('节点详情', '安装', '安装语言')
    },
    onTabChange () {
      if (this.isShowInstalled) {
        this.getInstalledDepList()
      } else {
        this.depName = ''
        this.depList = []
      }
      this.$st.sendEv('节点详情', '安装', '切换标签')
    }
  },
  async created () {
    const arr = this.$route.path.split('/')
    const id = arr[arr.length - 1]
    const res = await this.$request.get(`/nodes/${id}/langs`)
    this.langList = res.data.data
    this.activeTab = this.langList[0].executable_name || ''
  }
}
</script>

<style scoped>
  .node-installation >>> .el-button .el-loading-spinner {
    margin-top: -13px;
    height: 28px;
  }

  .node-installation >>> .el-button .el-loading-spinner .circular {
    width: 28px;
    height: 28px;
  }
</style>
