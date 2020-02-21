<template>
  <div class="spider-scrapy">
    <!--parameter edit-->
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
          @click="onSettingsActiveParamAdd"
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
              @click="onSettingsActiveParamRemove(scope.$index)"
            />
          </template>
        </el-table-column>
      </el-table>
      <template slot="footer">
        <el-button type="plain" size="small" @click="onCloseDialog">{{$t('Cancel')}}</el-button>
        <el-button type="primary" size="small" @click="onSettingsConfirm">
          {{$t('Confirm')}}
        </el-button>
      </template>
    </el-dialog>
    <!--./parameter edit-->

    <!--add scrapy spider-->
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
        <el-form-item :label="$t('Template')" prop="template" required>
          <el-select v-model="addSpiderForm.template" :placeholder="$t('Template')">
            <el-option value="basic" label="basic"/>
            <el-option value="crawl" label="crawl"/>
            <el-option value="csvfeed" label="csvfeed"/>
            <el-option value="xmlfeed" label="xmlfeed"/>
          </el-select>
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
    <!--./add scrapy spider-->

    <el-tabs
      v-model="activeTabName"
    >
      <!--settings-->
      <el-tab-pane :label="$t('Settings')" name="settings">
        <div class="settings">
          <div class="top-action-wrapper">
            <el-button
              type="primary"
              size="small"
              icon="el-icon-plus"
              @click="onSettingsAdd"
            >
              {{$t('Add')}}
            </el-button>
            <el-button size="small" type="success" @click="onSettingsSave" icon="el-icon-check">
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
                <el-select v-model="scope.row.type" size="small" @change="onSettingsParamTypeChange(scope.row)">
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
                    @click="onSettingsEditParam(scope.row, scope.$index)"
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
                  @click="onSettingsRemove(scope.$index)"
                />
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
      <!--./settings-->

      <!--spiders-->
      <el-tab-pane :label="$t('Spiders')" name="spiders">
        <div class="spiders">
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
      </el-tab-pane>
      <!--./spiders-->

      <!--items-->
      <el-tab-pane label="Items" name="items">
        <div class="items">
          <div class="action-wrapper">
            <el-button
              type="primary"
              size="small"
              icon="el-icon-plus"
              @click="onAddItem"
            >
              {{$t('Add Item')}}
            </el-button>
            <el-button size="small" type="success" @click="onItemsSave" icon="el-icon-check">
              {{$t('Save')}}
            </el-button>
          </div>
          <el-tree
            :data="spiderScrapyItemsConverted"
            default-expand-all
          >
          <span class="custom-tree-node" slot-scope="{ node, data }">
            <template v-if="data.level === 1">
              <span v-if="!node.isEdit" class="label" @click="onItemLabelEdit(node, data, $event)">
                {{ node.label }}
                <i class="el-icon-edit"></i>
              </span>
              <el-input
                v-else
                :ref="`el-input-${data.id}`"
                :placeholder="$t('Item Name')"
                v-model="data.label"
                size="mini"
                @blur="node.isEdit = false"
              />
              <span>
                <el-button
                  type="primary"
                  size="mini"
                  icon="el-icon-plus"
                  @click="onAddItemField(data, $event)">
                  {{$t('Add Field')}}
                </el-button>
                <el-button
                  type="danger"
                  size="mini"
                  icon="el-icon-delete"
                  @click="removeItem(data, $event)">
                  {{$t('Remove')}}
                </el-button>
              </span>
            </template>
            <template v-if="data.level === 2">
              <span v-if="!node.isEdit" class="label" @click="onItemLabelEdit(node, data, $event)">
                {{ node.label }}
                <i class="el-icon-edit"></i>
              </span>
              <el-input
                v-else
                :ref="`el-input-${data.id}`"
                :placeholder="$t('Field Name')"
                v-model="data.label"
                size="mini"
                @blur="node.isEdit = false"
              />
              <span>
                <el-button
                  type="danger"
                  size="mini"
                  icon="el-icon-delete"
                  @click="onRemoveItemField(node, data, $event)">
                  {{$t('Remove')}}
                </el-button>
              </span>
            </template>
          </span>
          </el-tree>
        </div>
      </el-tab-pane>
      <!--./items-->

      <!--pipelines-->
      <el-tab-pane label="Pipelines" name="pipelines">

      </el-tab-pane>
      <!--./pipelines-->
    </el-tabs>
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
      'spiderScrapySettings',
      'spiderScrapyItems'
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
    },
    spiderScrapyItemsConverted () {
      let id = 0
      return this.spiderScrapyItems.map(d => {
        d.id = id++
        d.label = d.name
        d.level = 1
        d.isEdit = false
        d.children = d.fields.map(f => {
          return {
            id: id++,
            label: f,
            level: 2,
            isEdit: false
          }
        })
        return d
      })
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
        domain: '',
        template: 'basic'
      },
      isAddSpiderLoading: false,
      activeTabName: 'settings'
    }
  },
  methods: {
    onOpenDialog () {
      this.dialogVisible = true
    },
    onCloseDialog () {
      this.dialogVisible = false
    },
    onSettingsConfirm () {
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
      this.$st('爬虫详情', 'Scrapy 设置', '确认编辑参数')
    },
    onSettingsEditParam (row, index) {
      this.activeParam = JSON.parse(JSON.stringify(row))
      this.activeParamIndex = index
      this.onOpenDialog()
      this.$st('爬虫详情', 'Scrapy 设置', '点击编辑参数')
    },
    async onSettingsSave () {
      const res = await this.$store.dispatch('spider/saveSpiderScrapySettings', this.$route.params.id)
      if (!res.data.error) {
        this.$message.success(this.$t('Saved successfully'))
      }
      this.$st('爬虫详情', 'Scrapy 设置', '保存设置')
    },
    onSettingsAdd () {
      const data = JSON.parse(JSON.stringify(this.spiderScrapySettings))
      data.push({
        key: '',
        value: '',
        type: 'string'
      })
      this.$store.commit('spider/SET_SPIDER_SCRAPY_SETTINGS', data)
      this.$st('爬虫详情', 'Scrapy 设置', '添加参数')
    },
    onSettingsRemove (index) {
      const data = JSON.parse(JSON.stringify(this.spiderScrapySettings))
      data.splice(index, 1)
      this.$store.commit('spider/SET_SPIDER_SCRAPY_SETTINGS', data)
      this.$st('爬虫详情', 'Scrapy 设置', '删除参数')
    },
    onSettingsActiveParamAdd () {
      if (this.activeParam.type === 'array') {
        this.activeParam.value.push('')
      } else if (this.activeParam.type === 'object') {
        if (!this.activeParam.value) {
          this.activeParam.value = {}
        }
        this.$set(this.activeParam.value, '', 999)
      }
      this.$st('爬虫详情', 'Scrapy 设置', '添加参数中参数')
    },
    onSettingsActiveParamRemove (index) {
      if (this.activeParam.type === 'array') {
        this.activeParam.value.splice(index, 1)
      } else if (this.activeParam.type === 'object') {
        const key = this.activeParamData[index].key
        const value = JSON.parse(JSON.stringify(this.activeParam.value))
        delete value[key]
        this.$set(this.activeParam, 'value', value)
      }
      this.$st('爬虫详情', 'Scrapy 设置', '删除参数中参数')
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
    onSettingsParamTypeChange (row) {
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
      this.$st('爬虫详情', 'Scrapy 设置', '确认添加爬虫')
    },
    onAddSpider () {
      this.addSpiderForm = {
        name: '',
        domain: '',
        template: 'basic'
      }
      this.isAddSpiderVisible = true
      this.$st('爬虫详情', 'Scrapy 设置', '添加爬虫')
    },
    onAddItem () {
      this.spiderScrapyItems.push({
        name: `Item_${+new Date()}`,
        fields: [
          `field_${+new Date()}`
        ]
      })
      this.$st('爬虫详情', 'Scrapy 设置', '添加Item')
    },
    removeItem (data, ev) {
      ev.stopPropagation()
      for (let i = 0; i < this.spiderScrapyItems.length; i++) {
        const item = this.spiderScrapyItems[i]
        if (item.name === data.label) {
          this.spiderScrapyItems.splice(i, 1)
          break
        }
      }
      this.$st('爬虫详情', 'Scrapy 设置', '删除Item')
    },
    onAddItemField (data, ev) {
      ev.stopPropagation()
      for (let i = 0; i < this.spiderScrapyItems.length; i++) {
        const item = this.spiderScrapyItems[i]
        if (item.name === data.label) {
          item.fields.push(`field_${+new Date()}`)
          break
        }
      }
      this.$st('爬虫详情', 'Scrapy 设置', '添加Items字段')
    },
    onRemoveItemField (node, data, ev) {
      ev.stopPropagation()
      for (let i = 0; i < this.spiderScrapyItems.length; i++) {
        const item = this.spiderScrapyItems[i]
        if (item.name === node.parent.label) {
          for (let j = 0; j < item.fields.length; j++) {
            const field = item.fields[j]
            if (field === data.label) {
              item.fields.splice(j, 1)
              break
            }
          }
        }
      }
      this.$st('爬虫详情', 'Scrapy 设置', '删除Items字段')
    },
    onItemLabelEdit (node, data, ev) {
      ev.stopPropagation()
      this.$set(node, 'isEdit', true)
      setTimeout(() => {
        this.$refs[`el-input-${data.id}`].focus()
      }, 0)
    },
    async onItemsSave () {
      const res = await this.$store.dispatch('spider/saveSpiderScrapyItems', this.$route.params.id)
      if (!res.data.error) {
        this.$message.success(this.$t('Saved successfully'))
      }
      this.$st('爬虫详情', 'Scrapy 设置', '保存Items')
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
    width: 100%;
    height: 100%;
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
    width: 100%;
    height: 100%;
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

  .items {
    width: 100%;
    height: 100%;
  }

  .items >>> .action-wrapper {
    text-align: right;
    padding-bottom: 10px;
    border-bottom: 1px solid #DCDFE6;
  }

  .items >>> .custom-tree-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    padding-right: 8px;
    min-height: 36px;
  }

  .items >>> .el-tree-node__content {
    height: auto;
  }

  .items >>> .custom-tree-node .label i.el-icon-edit {
    visibility: hidden;
  }

  .items >>> .custom-tree-node:hover .label i.el-icon-edit {
    visibility: visible;
  }

  .items >>> .custom-tree-node .el-input {
    width: 240px;
  }
</style>
