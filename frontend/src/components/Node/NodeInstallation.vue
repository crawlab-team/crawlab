<template>
  <div class="node-installation">
    <el-form inline>
      <el-form-item v-if="!isShowInstalled">
        <el-autocomplete
          v-model="depName"
          style="width: 240px"
          :placeholder="$t('Search Dependencies')"
          :fetchSuggestions="fetchAllDepList"
          minlength="2"
          @select="onSearch"
        />
      </el-form-item>
      <el-form-item>
        <el-button icon="el-icon-search" type="success" @click="onSearch">
          {{$t('Search')}}
        </el-button>
      </el-form-item>
      <el-form-item>
        <el-checkbox v-model="isShowInstalled" :label="$t('Show installed')" @change="onIsShowInstalledChange"/>
      </el-form-item>
    </el-form>
    <el-tabs v-model="activeTab">
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
          :label="$t('Latest Version')"
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
              size="mini"
              type="primary"
              @click="onClickInstallDep(scope.row)"
            >
              {{$t('Install')}}
            </el-button>
            <el-button
              v-else
              v-loading="getDepLoading(scope.row)"
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
          type="primary"
          style="width: 240px;font-weight: bolder;font-size: 18px"
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
      depLoadingDict: {}
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
      this.depList = res.data.data.sort((a, b) => a.name > b.name ? 1 : -1)
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
      const res = await this.$request.get('/system/deps', {
        lang: this.activeLang.executable_name,
        dep_name: queryString
      })
      callback(res.data.data ? res.data.data.map(d => {
        return { value: d, label: d }
      }) : [])
    },
    onSearch () {
      if (!this.isShowInstalled) {
        this.getDepList()
      } else {
        this.getInstalledDepList()
      }
    },
    onIsShowInstalledChange (val) {
      if (val) {
        this.getInstalledDepList()
      }
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
    },
    getDepLoading (dep) {
      const name = dep.name
      if (this.depLoadingDict[name] === undefined) {
        return false
      }
      return this.depLoadingDict[name]
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
  .install-wrapper >>> .el-button .el-loading-spinner {
    height: 100%;
  }
</style>
