<template>
  <div class="spider-scrapy">
    <el-dialog
      :title="$t('Parameter Edit')"
      :visible="dialogVisible"
      class="setting-param-dialog"
      width="600px"
      :before-close="onCloseDialog"
    >
      <div class="action-wrapper" style="margin-bottom: 10px;text-align: right">
        <el-button
          type="primary"
          size="small"
          icon="el-icon-plus"
          @click="onActiveParamAdd"
        >
          {{$t('Add')}}
        </el-button>
      </div>
      <el-table
        :data="activeParamData"
      >
        <el-table-column
          v-if="activeParam.type === 'object'"
          :label="$t('Key')"
        >
          <template slot-scope="scope">
            <el-input v-model="scope.row.key" size="small"/>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Value')"
        >
          <template slot-scope="scope">
            <el-input
              v-if="activeParam.type === 'object'"
              v-model="scope.row.value"
              size="small"
              type="number"
              @change="() => scope.row.value = Number(scope.row.value)"
            />
            <el-input
              v-else-if="activeParam.type === 'array'"
              v-model="scope.row.value"
              size="small"
            />
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Action')"
          width="60px"
          align="center"
        >
          <template slot-scope="scope">
            <el-button
              type="danger"
              size="mini"
              icon="el-icon-delete"
              circle
              @click="onActiveParamRemove(scope.$index)"
            />
          </template>
        </el-table-column>
      </el-table>
      <template slot="footer">
        <el-button type="plain" size="small" @click="onCloseDialog">{{$t('Cancel')}}</el-button>
        <el-button type="primary" size="small" @click="onConfirm">
          {{$t('Confirm')}}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      :title="$t('Add Scrapy Spider')"
      :visible.sync="isAddSpiderVisible"
      width="480px"
    >
      <el-form
        :model="addSpiderForm"
        label-width="80px"
        ref="add-spider-form"
        inline-message
      >
        <el-form-item :label="$t('Name')" prop="name" required>
          <el-input v-model="addSpiderForm.name" :placeholder="$t('Name')"/>
        </el-form-item>
        <el-form-item :label="$t('Domain')" prop="domain" required>
          <el-input v-model="addSpiderForm.domain" :placeholder="$t('Domain')"/>
        </el-form-item>
      </el-form>
      <template slot="footer">
        <el-button type="plain" size="small" @click="isAddSpiderVisible = false">{{$t('Cancel')}}</el-button>
        <el-button
          type="primary"
          size="small"
          @click="onAddSpiderConfirm"
          :icon="isAddSpiderLoading ? 'el-icon-loading' : ''"
          :disabled="isAddSpiderLoading"
        >
          {{$t('Confirm')}}
        </el-button>
      </template>
    </el-dialog>

    <div class="spiders">
      <h3 class="title">{{$t('Scrapy Spiders')}}</h3>
      <div class="action-wrapper">
        <el-button
          type="primary"
          size="small"
          icon="el-icon-plus"
          @click="onAddSpider"
        >
          {{$t('Add Spider')}}
        </el-button>
      </div>
      <ul class="spider-list">
        <li
          v-for="s in spiderForm.spider_names"
          :key="s"
          class="spider-item"
        >
          <i class="el-icon-caret-right"></i>
          {{s}}
        </li>
      </ul>
    </div>
    <div class="settings">
      <h3 class="title">{{$t('Settings')}}</h3>
      <div class="top-action-wrapper">
        <el-button
          type="primary"
          size="small"
          icon="el-icon-plus"
          @click="onAdd"
        >
          {{$t('Add')}}
        </el-button>
        <el-button size="small" type="success" @click="onSave" icon="el-icon-check">
          {{$t('Save')}}
        </el-button>
      </div>
      <el-table
        :data="spiderScrapySettings"
        border
        :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
        max-height="calc(100vh - 240px)"
      >
        <el-table-column
          :label="$t('Variable Name')"
          width="240px"
        >
          <template slot-scope="scope">
            <el-autocomplete
              v-model="scope.row.key"
              size="small"
              suffix-icon="el-icon-edit"
              :fetch-suggestions="settingsKeysFetchSuggestions"
            />
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Variable Type')"
          width="120px"
        >
          <template slot-scope="scope">
            <el-select v-model="scope.row.type" size="small" @change="onParamTypeChange(scope.row)">
              <el-option value="string" :label="$t('String')"/>
              <el-option value="number" :label="$t('Number')"/>
              <el-option value="boolean" :label="$t('Boolean')"/>
              <el-option value="array" :label="$t('Array/List')"/>
              <el-option value="object" :label="$t('Object/Dict')"/>
            </el-select>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Variable Value')"
          width="calc(100% - 150px)"
        >
          <template slot-scope="scope">
            <el-input
              v-if="scope.row.type === 'string'"
              v-model="scope.row.value"
              size="small"
              suffix-icon="el-icon-edit"
            />
            <el-input
              v-else-if="scope.row.type === 'number'"
              type="number"
              v-model="scope.row.value"
              size="small"
              suffix-icon="el-icon-edit"
              @change="scope.row.value = Number(scope.row.value)"
            />
            <div
              v-else-if="scope.row.type === 'boolean'"
              style="margin-left: 10px"
            >
              <el-switch
                v-model="scope.row.value"
                size="small"
                active-color="#67C23A"
              />
            </div>
            <div
              v-else
              style="margin-left: 10px;font-size: 12px"
            >
              {{JSON.stringify(scope.row.value)}}
              <el-button
                type="warning"
                size="mini"
                icon="el-icon-edit"
                style="margin-left: 10px"
                @click="onEditParam(scope.row, scope.$index)"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Action')"
          width="60px"
          align="center"
        >
          <template slot-scope="scope">
            <el-button
              type="danger"
              size="mini"
              icon="el-icon-delete"
              circle
              @click="onRemove(scope.$index)"
            />
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'SpiderScrapy',
  computed: {
    ...mapState('spider', [
      'spiderForm',
      'spiderScrapySettings'
    ]),
    activeParamData () {
      if (this.activeParam.type === 'array') {
        return this.activeParam.value.map(s => {
          return { value: s }
        })
      } else if (this.activeParam.type === 'object') {
        return Object.keys(this.activeParam.value).map(key => {
          return {
            key,
            value: this.activeParam.value[key]
          }
        })
      }
      return []
    }
  },
  data () {
    return {
      dialogVisible: false,
      activeParam: {},
      activeParamIndex: undefined,
      isAddSpiderVisible: false,
      addSpiderForm: {
        name: '',
        domain: ''
      },
      isAddSpiderLoading: false
    }
  },
  methods: {
    onOpenDialog () {
      this.dialogVisible = true
    },
    onCloseDialog () {
      this.dialogVisible = false
    },
    onConfirm () {
      if (this.activeParam.type === 'array') {
        this.activeParam.value = this.activeParamData.map(d => d.value)
      } else if (this.activeParam.type === 'object') {
        const dict = {}
        this.activeParamData.forEach(d => {
          dict[d.key] = d.value
        })
        this.activeParam.value = dict
      }
      this.$set(this.spiderScrapySettings, this.activeParamIndex, JSON.parse(JSON.stringify(this.activeParam)))
      this.dialogVisible = false
    },
    onEditParam (row, index) {
      this.activeParam = JSON.parse(JSON.stringify(row))
      this.activeParamIndex = index
      this.onOpenDialog()
    },
    async onSave () {
      const res = await this.$store.dispatch('spider/saveSpiderScrapySettings', this.$route.params.id)
      if (!res.data.error) {
        this.$message.success(this.$t('Saved successfully'))
      }
    },
    onAdd () {
      const data = JSON.parse(JSON.stringify(this.spiderScrapySettings))
      data.push({
        key: '',
        value: '',
        type: 'string'
      })
      this.$store.commit('spider/SET_SPIDER_SCRAPY_SETTINGS', data)
    },
    onRemove (index) {
      const data = JSON.parse(JSON.stringify(this.spiderScrapySettings))
      data.splice(index, 1)
      this.$store.commit('spider/SET_SPIDER_SCRAPY_SETTINGS', data)
    },
    onActiveParamAdd () {
      if (this.activeParam.type === 'array') {
        this.activeParam.value.push('')
      } else if (this.activeParam.type === 'object') {
        if (!this.activeParam.value) {
          this.activeParam.value = {}
        }
        this.$set(this.activeParam.value, '', 999)
      }
    },
    onActiveParamRemove (index) {
      if (this.activeParam.type === 'array') {
        this.activeParam.value.splice(index, 1)
      } else if (this.activeParam.type === 'object') {
        const key = this.activeParamData[index].key
        const value = JSON.parse(JSON.stringify(this.activeParam.value))
        delete value[key]
        this.$set(this.activeParam, 'value', value)
      }
    },
    settingsKeysFetchSuggestions (queryString, cb) {
      const data = this.$utils.scrapy.settingParamNames
        .filter(s => {
          if (!queryString) return true
          return !!s.match(new RegExp(queryString, 'i'))
        })
        .map(s => {
          return {
            value: s,
            label: s
          }
        })
        .sort((a, b) => {
          return a > b ? -1 : 1
        })
      cb(data)
    },
    onParamTypeChange (row) {
      if (row.type === 'number') {
        row.value = Number(row.value)
      }
    },
    onAddSpiderConfirm () {
      this.$refs['add-spider-form'].validate(async valid => {
        if (!valid) return
        this.isAddSpiderLoading = true
        const res = await this.$store.dispatch('spider/addSpiderScrapySpider', {
          id: this.$route.params.id,
          form: this.addSpiderForm
        })
        if (!res.data.error) {
          this.$message.success(this.$t('Saved successfully'))
        }
        this.isAddSpiderVisible = false
        this.isAddSpiderLoading = false
        await this.$store.dispatch('spider/getSpiderScrapySpiders', this.$route.params.id)
      })
    },
    onAddSpider () {
      this.addSpiderForm = {
        name: '',
        domain: ''
      }
      this.isAddSpiderVisible = true
    }
  }
}
</script>

<style scoped>
  .spider-scrapy {
    height: calc(100vh - 200px);
    color: #606266;
  }

  .spiders {
    float: left;
    display: inline-block;
    width: 240px;
    height: 100%;
    border: 1px solid #DCDFE6;
    border-radius: 3px;
    padding: 0 10px;
  }

  .spiders .title {
    border-bottom: 1px solid #DCDFE6;
    padding-bottom: 15px;
  }

  .spiders .action-wrapper {
    margin-bottom: 10px;
    text-align: right;
  }

  .spiders .spider-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .spiders .spider-list .spider-item {
    padding: 10px;
    cursor: pointer;
  }

  .spiders .spider-list .spider-item:hover {
    background: #F5F7FA;
  }

  .settings {
    margin-left: 20px;
    border: 1px solid #DCDFE6;
    float: left;
    width: calc(100% - 240px - 20px);
    height: 100%;
    border-radius: 3px;
    padding: 0 20px;
  }

  .settings .title {
    border-bottom: 1px solid #DCDFE6;
    padding-bottom: 15px;
  }

  .settings >>> .el-table td,
  .settings >>> .el-table td .cell {
    padding: 0;
    margin: 0;
  }

  .settings >>> .el-table td .cell .el-autocomplete {
    width: 100%;
  }

  .settings >>> .el-table td .cell .el-input__inner {
    border: none;
    font-size: 12px;
  }

  .settings >>> .action-wrapper {
    margin-top: 10px;
    text-align: right;
  }

  .settings >>> .top-action-wrapper {
    margin-bottom: 10px;
    text-align: right;
  }

  .settings >>> .top-action-wrapper .el-button {
    margin-left: 10px;
  }
</style>
