<template>
  <div class="node-installation-matrix">
    <el-tabs v-model="activeTabName">
      <el-tab-pane :label="$t('Languages')" name="lang">
        <div class="lang-table">
          <el-table
            class="table"
            :data="nodeList"
            :header-cell-style="{background:'rgb(48, 65, 86)',color:'white',height:'50px'}"
            border
            @row-click="onLangTableRowClick"
          >
            <el-table-column
              :label="$t('Node')"
              width="240px"
              prop="name"
              fixed
            />
            <el-table-column
              :label="$t('nodeList.type')"
              width="120px"
              fixed
            >
              <template slot-scope="scope">
                <el-tag v-if="scope.row.is_master" type="primary">{{ $t('Master') }}</el-tag>
                <el-tag v-else type="warning">{{ $t('Worker') }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column
              :label="$t('Status')"
              width="120px"
              fixed
            >
              <template slot-scope="scope">
                <el-tag v-if="scope.row.status === 'offline'" type="info">{{ $t('Offline') }}</el-tag>
                <el-tag v-else-if="scope.row.status === 'online'" type="success">{{ $t('Online') }}</el-tag>
                <el-tag v-else type="danger">{{ $t('Unavailable') }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column
              v-for="l in langs"
              :key="l.name"
              :label="l.label"
              width="220px"
            >
              <template slot="header" slot-scope="scope">
                <div class="header-with-action">
                  <span>{{ scope.column.label }}</span>
                  <el-button type="primary" size="mini" @click="onInstallLangAll(scope.column.label, $event)">
                    {{ $t('Install') }}
                  </el-button>
                </div>
              </template>
              <template slot-scope="scope">
                <template v-if="getLangInstallStatus(scope.row._id, l.name) === 'installed'">
                  <el-tag type="success">
                    <i class="el-icon-check" />
                    {{ $t('Installed') }}
                  </el-tag>
                </template>
                <template v-else-if="getLangInstallStatus(scope.row._id, l.name) === 'installing'">
                  <el-tag type="warning">
                    <i class="el-icon-loading" />
                    {{ $t('Installing') }}
                  </el-tag>
                </template>
                <template
                  v-else-if="['installing-other', 'not-installed'].includes(getLangInstallStatus(scope.row._id, l.name))"
                >
                  <div class="cell-with-action">
                    <el-tag type="danger">
                      <i class="el-icon-error" />
                      {{ $t('Not Installed') }}
                    </el-tag>
                    <el-button type="primary" size="mini" @click="onInstallLang(scope.row._id, scope.column.label, $event)">
                      {{ $t('Install') }}
                    </el-button>
                  </div>
                </template>
                <template v-else-if="getLangInstallStatus(scope.row._id, l.name) === 'na'">
                  <el-tag type="info">
                    <i class="el-icon-question" />
                    {{ $t('N/A') }}
                  </el-tag>
                </template>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
      <el-tab-pane :label="$t('Dependencies')" name="dep">
        <el-form class="search-form" inline>
          <el-form-item>
            <el-input
              v-model="depName"
              style="width: 240px"
              :placeholder="$t('Search Dependencies')"
            />
          </el-form-item>
          <el-form-item>
            <el-button
              size="small"
              icon="el-icon-search"
              type="success"
              @click="onSearch"
            >
              {{ $t('Search') }}
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-checkbox v-model="isShowInstalled" :label="$t('Show installed')" @change="onIsShowInstalledChange" />
          </el-form-item>
        </el-form>
        <el-tabs v-model="activeLang">
          <el-tab-pane
            v-for="l in langsWithDeps"
            :key="l.name"
            :name="l.name"
            :label="l.label"
          />
        </el-tabs>
        <el-table
          v-loading="isDepsLoading"
          class="table"
          height="calc(100vh - 320px)"
          :data="computedDepsSet"
          :header-cell-style="{background:'rgb(48, 65, 86)',color:'white',height:'50px'}"
          border
        >
          <el-table-column
            :label="$t('Dependency')"
            prop="name"
            width="180px"
            fixed
          />
          <el-table-column
            v-if="false"
            :label="$t('Install on All Nodes')"
            width="120px"
            align="center"
            fixed
          >
            <template>
              <el-button
                size="mini"
                type="primary"
              >
                {{ $t('Install') }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column
            v-for="n in activeNodes"
            :key="n._id"
            :label="n.name"
            width="220px"
            align="center"
          >
            <template slot="header" slot-scope="scope">
              {{ scope.column.label }}
            </template>
            <template slot-scope="scope">
              <div
                v-if="getDepStatus(n, scope.row) === 'installed'"
                class="cell-with-action"
              >
                <el-tag type="success">
                  {{ $t('Installed') }}
                </el-tag>
                <el-button
                  size="mini"
                  type="danger"
                  @click="uninstallDep(n, scope.row)"
                >
                  {{ $t('Uninstall') }}
                </el-button>
              </div>
              <div
                v-else-if="getDepStatus(n, scope.row) === 'installing'"
                class="cell-with-action"
              >
                <el-tag type="warning">
                  <i class="el-icon-loading" />
                  {{ $t('Installing') }}
                </el-tag>
              </div>
              <div
                v-else-if="getDepStatus(n, scope.row) === 'uninstalled'"
                class="cell-with-action"
              >
                <el-tag type="danger">
                  {{ $t('Not Installed') }}
                </el-tag>
                <el-button
                  size="mini"
                  type="primary"
                  @click="installDep(n, scope.row)"
                >
                  {{ $t('Install') }}
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
  import { mapState } from 'vuex'

  export default {
    name: 'NodeInstallationMatrix',
    props: {
      activeTab: {
        type: String,
        default: ''
      }
    },
    data() {
      return {
        langs: [
          { label: 'Python', name: 'python', hasDeps: true },
          { label: 'Node.js', name: 'node', hasDeps: true },
          { label: 'Java', name: 'java', hasDeps: false },
          { label: '.Net Core', name: 'dotnet', hasDeps: false },
          { label: 'PHP', name: 'php', hasDeps: false }
        ],
        langsDataDict: {},
        handle: undefined,
        activeTabName: 'lang',
        depsDataDict: {},
        depsSet: new Set(),
        activeLang: 'python',
        isDepsLoading: false,
        depName: '',
        isShowInstalled: true,
        depList: []
      }
    },
    computed: {
      ...mapState('node', [
        'nodeList'
      ]),
      activeNodes() {
        return this.nodeList.filter(d => d.status === 'online')
      },
      computedDepsSet() {
        return Array.from(this.depsSet).map(d => {
          return {
            name: d
          }
        })
      },
      langsWithDeps() {
        return this.langs.filter(l => l.hasDeps)
      }
    },
    watch: {
      activeLang() {
        this.getDepsData()
      }
    },
    async created() {
      setTimeout(() => {
        this.getLangsData()
        this.getDepsData()
      }, 1000)

      this.handle = setInterval(() => {
        this.getLangsData()
      }, 10000)
    },
    destroyed() {
      clearInterval(this.handle)
    },
    methods: {
      async getLangsData() {
        await Promise.all(this.nodeList.map(async n => {
          if (n.status !== 'online') return
          const res = await this.$request.get(`/nodes/${n._id}/langs`)
          if (!res.data.data) return
          res.data.data.forEach(l => {
            const key = n._id + '|' + l.executable_name
            this.$set(this.langsDataDict, key, l)
          })
        }))
      },
      async getDepsData() {
        this.isDepsLoading = true
        this.depsDataDict = {}
        this.depsSet = new Set()
        const depsSet = new Set()
        await Promise.all(this.nodeList.map(async n => {
          if (n.status !== 'online') return
          const res = await this.$request.get(`/nodes/${n._id}/deps/installed`, { lang: this.activeLang })
          if (!res.data.data) return
          res.data.data.forEach(d => {
            depsSet.add(d.name)
            const key = n._id + '|' + d.name
            this.$set(this.depsDataDict, key, 'installed')
          })
        }))
        this.depsSet = depsSet
        this.isDepsLoading = false
      },
      getLang(nodeId, langName) {
        const key = nodeId + '|' + langName
        return this.langsDataDict[key]
      },
      getLangInstallStatus(nodeId, langName) {
        const lang = this.getLang(nodeId, langName)
        if (!lang || !lang.install_status) return 'na'
        return lang.install_status
      },
      getLangFromLabel(label) {
        for (let i = 0; i < this.langs.length; i++) {
          const lang = this.langs[i]
          if (lang.label === label) {
            return lang
          }
        }
      },
      async onInstallLang(nodeId, langLabel, ev) {
        if (ev) {
          ev.stopPropagation()
        }
        const lang = this.getLangFromLabel(langLabel)
        this.$request.post(`/nodes/${nodeId}/langs/install`, {
          lang: lang.name
        })
        const key = nodeId + '|' + lang.name
        this.$set(this.langsDataDict[key], 'install_status', 'installing')
        setTimeout(() => {
          this.getLangsData()
        }, 1000)
        this.$request.put('/actions', {
          type: 'install_lang'
        })
        this.$st.sendEv('节点列表', '安装', '安装语言')
      },
      async onInstallLangAll(langLabel, ev) {
        ev.stopPropagation()
        this.nodeList
          .filter(n => {
            if (n.status !== 'online') return false
            const lang = this.getLangFromLabel(langLabel)
            const key = n._id + '|' + lang.name
            if (!this.langsDataDict[key]) return false
            if (['installing', 'installed'].includes(this.langsDataDict[key].install_status)) return false
            return true
          })
          .forEach(n => {
            this.onInstallLang(n._id, langLabel, ev)
          })
        setTimeout(() => {
          this.getLangsData()
        }, 1000)
        this.$st.sendEv('节点列表', '安装', '安装语言-所有节点')
      },
      onLangTableRowClick(row) {
        this.$router.push(`/nodes/${row._id}`)
        this.$st.sendEv('节点列表', '安装', '查看节点详情')
      },
      getDepStatus(node, dep) {
        const key = node._id + '|' + dep.name
        if (!this.depsDataDict[key]) {
          return 'uninstalled'
        } else {
          return this.depsDataDict[key]
        }
      },
      async installDep(node, dep) {
        const key = node._id + '|' + dep.name
        this.$set(this.depsDataDict, key, 'installing')
        const data = await this.$request.post(`/nodes/${node._id}/deps/install`, {
          lang: this.activeLang,
          dep_name: dep.name
        })
        if (!data || data.error) {
          this.$notify.error({
            title: this.$t('Installing dependency failed'),
            message: this.$t('The dependency installation is unsuccessful: ') + dep.name
          })
          this.$set(this.depsDataDict, key, 'uninstalled')
        } else {
          this.$notify.success({
            title: this.$t('Installing dependency successful'),
            message: this.$t('You have successfully installed a dependency: ') + dep.name
          })
          this.$set(this.depsDataDict, key, 'installed')
        }
        this.$request.put('/actions', {
          type: 'install_dep'
        })
        this.$st.sendEv('节点列表', '安装', '安装依赖')
      },
      async uninstallDep(node, dep) {
        const key = node._id + '|' + dep.name
        this.$set(this.depsDataDict, key, 'installing')
        const data = await this.$request.post(`/nodes/${node._id}/deps/uninstall`, {
          lang: this.activeLang,
          dep_name: dep.name
        })
        if (!data || data.error) {
          this.$notify.error({
            title: this.$t('Uninstalling dependency failed'),
            message: this.$t('The dependency uninstallation is unsuccessful: ') + dep.name
          })
          this.$set(this.depsDataDict, key, 'installed')
        } else {
          this.$notify.success({
            title: this.$t('Uninstalling dependency successful'),
            message: this.$t('You have successfully uninstalled a dependency: ') + dep.name
          })
          this.$set(this.depsDataDict, key, 'uninstalled')
        }
        this.$st.sendEv('节点列表', '安装', '卸载依赖')
      },
      onSearch() {
        this.isShowInstalled = false
        this.getDepList()
        this.$st.sendEv('节点列表', '安装', '搜索依赖')
      },
      async getDepList() {
        const masterNode = this.nodeList.filter(n => n.is_master)[0]
        this.depsSet = []
        this.isDepsLoading = true
        const res = await this.$request.get(`/nodes/${masterNode._id}/deps`, {
          lang: this.activeLang,
          dep_name: this.depName
        })
        this.isDepsLoading = false
        this.depsSet = new Set(res.data.data.map(d => d.name))
      },
      onIsShowInstalledChange(val) {
        if (val) {
          this.getDepsData()
        } else {
          this.depsSet = []
        }
        this.$st.sendEv('节点列表', '安装', '点击查看已安装')
      }
    }
  }
</script>

<style scoped>
  .table {
    margin-top: 20px;
    border-radius: 5px;
  }

  .el-table tr {
    cursor: pointer;
  }

  .el-table .header-with-action,
  .el-table .cell-with-action {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
</style>
