<template>
  <div class="node-installation-matrix">
    <div class="lang-table">
      <el-table
        class="table"
        :data="nodeList"
        :header-cell-style="{background:'rgb(48, 65, 86)',color:'white',height:'50px'}"
        border
      >
        <el-table-column
          :label="$t('Node')"
          width="240px"
          prop="name"
        />
        <el-table-column
          :label="$t('nodeList.type')"
          width="120px"
        >
          <template slot-scope="scope">
            <el-tag type="primary" v-if="scope.row.is_master">{{$t('Master')}}</el-tag>
            <el-tag type="warning" v-else>{{$t('Worker')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Status')"
          width="120px"
        >
          <template slot-scope="scope">
            <el-tag type="info" v-if="scope.row.status === 'offline'">{{$t('Offline')}}</el-tag>
            <el-tag type="success" v-else-if="scope.row.status === 'online'">{{$t('Online')}}</el-tag>
            <el-tag type="danger" v-else>{{$t('Unavailable')}}</el-tag>
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
              <span>{{scope.column.label}}</span>
              <el-button type="primary" size="mini" @click="onInstallAll(scope.column.name)">
                {{$t('Install')}}
              </el-button>
            </div>
          </template>
          <template slot-scope="scope">
            <template v-if="getLangInstallStatus(scope.row._id, l.name) === 'installed'">
              <el-tag type="success">
                <i class="el-icon-check"></i>
                {{$t('Installed')}}
              </el-tag>
            </template>
            <template v-else-if="getLangInstallStatus(scope.row._id, l.name) === 'installing'">
              <el-tag type="warning">
                <i class="el-icon-loading"></i>
                {{$t('Installing')}}
              </el-tag>
            </template>
            <template
              v-else-if="['installing-other', 'not-installed'].includes(getLangInstallStatus(scope.row._id, l.name))"
            >
              <div class="cell-with-action">
                <el-tag type="danger">
                  <i class="el-icon-error"></i>
                  {{$t('Not Installed')}}
                </el-tag>
                <el-button type="primary" size="mini" @click="onInstall(scope.row._id, scope.column.name)">
                  {{$t('Install')}}
                </el-button>
              </div>
            </template>
            <template v-else-if="getLangInstallStatus(scope.row._id, l.name) === 'na'">
              <el-tag type="info">
                <i class="el-icon-question"></i>
                {{$t('N/A')}}
              </el-tag>
            </template>
          </template>
        </el-table-column>
      </el-table>
    </div>
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
  data () {
    return {
      langs: [
        { label: 'Python', name: 'python' },
        { label: 'Node.js', name: 'node' },
        { label: 'Java', name: 'java' }
      ],
      dataDict: {},
      handle: undefined
    }
  },
  computed: {
    ...mapState('node', [
      'nodeList'
    ])
  },
  methods: {
    async getData () {
      await Promise.all(this.nodeList.map(async n => {
        const res = await this.$request.get(`/nodes/${n._id}/langs`)
        res.data.data.forEach(l => {
          const key = n._id + '|' + l.executable_name
          this.dataDict[key] = l
        })
      }))
    },
    getLang (nodeId, langName) {
      const key = nodeId + '|' + langName
      return this.dataDict[key]
    },
    getLangInstallStatus (nodeId, langName) {
      const lang = this.getLang(nodeId, langName)
      if (!lang || !lang.install_status) return 'na'
      return lang.install_status
    },
    async onInstall (nodeId, langName) {
      this.$request.post(`/nodes/${nodeId}/langs/install`, {
        lang: langName
      })
      setTimeout(() => {
        this.getData()
      }, 1000)
    },
    async onInstallAll (langName) {
      this.nodeList.map(async n => {
        return this.onInstall(n._id, langName)
      })
      setTimeout(() => {
        this.getData()
      }, 1000)
    }
  },
  async created () {
    await this.getData()

    this.handle = setInterval(() => {
      this.getData()
    }, 15000)
  },
  destroyed () {
    clearInterval(this.handle)
  }
}
</script>

<style scoped>
  .table {
    margin-top: 20px;
    border-radius: 5px;
  }

  .lang-table {
  }

  .lang-table >>> .el-table tr {
    cursor: initial;
  }

  .lang-table >>> .el-table .header-with-action,
  .lang-table >>> .el-table .cell-with-action {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
</style>
