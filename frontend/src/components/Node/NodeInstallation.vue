<template>
  <div class="node-installation">
    <el-form inline>
      <el-form-item>
        <el-input
          v-model="depName"
          style="width: 240px"
          :placeholder="$t('Search Dependencies')"
        />
      </el-form-item>
      <el-form-item>
        <el-button icon="el-icon-search" type="success" @click="getDepList">
          {{$t('Search')}}
        </el-button>
      </el-form-item>
    </el-form>
    <el-tabs v-model="activeTab">
      <el-tab-pane v-for="lang in langList" :key="lang.name" :label="lang.name" :name="lang.executable_name"/>
    </el-tabs>
    <template v-if="activeLang.installed">
      <el-table
        :data="depList"
        :empty-text="depName ? $t('No Data') : $t('Please search dependencies')"
        v-loading="loading"
      >
        <el-table-column
          :label="$t('Name')"
          prop="name"
        />
        <el-table-column
          :label="$t('Version')"
          prop="version"
        />
        <el-table-column
          :label="$t('Description')"
          prop="description"
        />
        <el-table-column
          :label="$t('Installed')"
          prop="installed"
        />
      </el-table>
    </template>
    <template v-else>
      <div class="install-wrapper">
        <h3>{{activeLang.name + $t(' is not installed, do you want to install it?')}}</h3>
        <el-button type="primary" style="width: 240px;font-weight: bolder;font-size: 18px">
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
      loading: false
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
    }
  },
  methods: {
    async getDepList () {
      const res = await this.$request.get(`/nodes/${this.nodeForm._id}/deps`, {
        lang: this.activeLang.executable_name,
        dep_name: this.depName
      })
      this.depList = res.data.data
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

</style>
