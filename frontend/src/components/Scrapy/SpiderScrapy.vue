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
          <div v-if="!spiderScrapySettings || !spiderScrapySettings.length" class="settings">
            <span class="empty-text">
              {{$t('No data available')}}
            </span>
            <template v-if="spiderScrapyErrors.settings">
              <label class="errors-label">{{$t('Errors')}}:</label>
              <el-alert type="error" v-html="getScrapyErrors('settings')"/>
            </template>
          </div>
          <div v-else class="settings">
            <div class="top-action-wrapper">
              <el-button
                type="primary"
                size="small"
                icon="el-icon-plus"
                @click="onSettingsAdd"
              >
                {{$t('Add Variable')}}
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
        </div>
      </el-tab-pane>
      <!--./settings-->

      <!--spiders-->
      <el-tab-pane :label="$t('Spiders')" name="spiders">
        <div class="spiders">
          <div v-if="!spiderForm.spider_names || !spiderForm.spider_names.length" class="spiders">
            <span class="empty-text error">
              {{$t('No data available. Please check whether your spiders are missing dependencies or no spiders created.')}}
            </span>
            <template v-if="spiderScrapyErrors.spiders">
              <label class="errors-label">{{$t('Errors')}}:</label>
              <el-alert type="error" v-html="getScrapyErrors('spiders')"/>
            </template>
          </div>
          <div v-else class="spiders">
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
            <ul class="list">
              <li
                v-for="s in spiderForm.spider_names"
                :key="s"
                class="item"
                @click="onClickSpider(s)"
              >
                <i class="el-icon-star-on"></i>
                {{s}}
                <i v-if="loadingDict[s]" class="el-icon-loading"></i>
              </li>
            </ul>
          </div>
        </div>
      </el-tab-pane>
      <!--./spiders-->

      <!--items-->
      <el-tab-pane label="Items" name="items">
        <div class="items">
          <div v-if="!spiderScrapyItems || !spiderScrapyItems.length" class="items">
            <span class="empty-text">
              {{$t('No data available')}}
            </span>
            <template v-if="spiderScrapyErrors.items">
              <label class="errors-label">{{$t('Errors')}}:</label>
              <el-alert type="error" v-html="getScrapyErrors('items')"/>
            </template>
          </div>
          <div v-else class="items">
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
              :data="spiderScrapyItems"
              default-expand-all
            >
              <span class="custom-tree-node" :class="`level-${data.level}`" slot-scope="{ node, data }">
                <template v-if="data.level === 1">
                  <span v-if="!node.isEdit" class="label" @click="onItemLabelEdit(node, data, $event)">
                    <i class="el-icon-star-on"></i>
                    {{ data.label }}
                    <i class="el-icon-edit"></i>
                  </span>
                  <el-input
                    v-else
                    :ref="`el-input-${data.id}`"
                    :placeholder="$t('Item Name')"
                    v-model="data.name"
                    size="mini"
                    @change="onItemChange(node, data, $event)"
                    @blur="$set(node, 'isEdit', false)"
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
                    <i class="el-icon-arrow-right"></i>
                    {{ node.label }}
                    <i class="el-icon-edit"></i>
                  </span>
                  <el-input
                    v-else
                    :ref="`el-input-${data.id}`"
                    :placeholder="$t('Field Name')"
                    v-model="data.name"
                    size="mini"
                    @change="onItemFieldChange(node, data, $event)"
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
        </div>
      </el-tab-pane>
      <!--./items-->

      <!--pipelines-->
      <el-tab-pane label="Pipelines" name="pipelines">
        <div v-if="!spiderScrapyPipelines || !spiderScrapyPipelines.length" class="pipelines">
          <span class="empty-text">
            {{$t('No data available')}}
          </span>
          <template v-if="spiderScrapyErrors.pipelines">
            <label class="errors-label">{{$t('Errors')}}:</label>
            <el-alert type="error" v-html="getScrapyErrors('pipelines')"/>
          </template>
        </div>
        <div class="pipelines">
          <ul class="list">
            <li
              v-for="s in spiderScrapyPipelines"
              :key="s"
              class="item"
              @click="$emit('click-pipeline')"
            >
              <i class="el-icon-star-on"></i>
              {{s}}
            </li>
          </ul>
        </div>
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
        'spiderScrapyItems',
        'spiderScrapyPipelines',
        'spiderScrapyErrors'
      ]),
      activeParamData() {
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
    data() {
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
        activeTabName: 'settings',
        loadingDict: {}
      }
    },
    methods: {
      onOpenDialog() {
        this.dialogVisible = true
      },
      onCloseDialog() {
        this.dialogVisible = false
      },
      onSettingsConfirm() {
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
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '确认编辑参数')
      },
      onSettingsEditParam(row, index) {
        this.activeParam = JSON.parse(JSON.stringify(row))
        this.activeParamIndex = index
        this.onOpenDialog()
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '点击编辑参数')
      },
      async onSettingsSave() {
        const res = await this.$store.dispatch('spider/saveSpiderScrapySettings', this.$route.params.id)
        if (!res.data.error) {
          this.$message.success(this.$t('Saved successfully'))
        }
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '保存设置')
      },
      onSettingsAdd() {
        const data = JSON.parse(JSON.stringify(this.spiderScrapySettings))
        data.push({
          key: '',
          value: '',
          type: 'string'
        })
        this.$store.commit('spider/SET_SPIDER_SCRAPY_SETTINGS', data)
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '添加参数')
      },
      onSettingsRemove(index) {
        const data = JSON.parse(JSON.stringify(this.spiderScrapySettings))
        data.splice(index, 1)
        this.$store.commit('spider/SET_SPIDER_SCRAPY_SETTINGS', data)
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '删除参数')
      },
      onSettingsActiveParamAdd() {
        if (this.activeParam.type === 'array') {
          this.activeParam.value.push('')
        } else if (this.activeParam.type === 'object') {
          if (!this.activeParam.value) {
            this.activeParam.value = {}
          }
          this.$set(this.activeParam.value, '', 999)
        }
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '添加参数中参数')
      },
      onSettingsActiveParamRemove(index) {
        if (this.activeParam.type === 'array') {
          this.activeParam.value.splice(index, 1)
        } else if (this.activeParam.type === 'object') {
          const key = this.activeParamData[index].key
          const value = JSON.parse(JSON.stringify(this.activeParam.value))
          delete value[key]
          this.$set(this.activeParam, 'value', value)
        }
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '删除参数中参数')
      },
      settingsKeysFetchSuggestions(queryString, cb) {
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
      onSettingsParamTypeChange(row) {
        if (row.type === 'number') {
          row.value = Number(row.value)
        }
      },
      onAddSpiderConfirm() {
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
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '确认添加爬虫')
      },
      onAddSpider() {
        this.addSpiderForm = {
          name: '',
          domain: '',
          template: 'basic'
        }
        this.isAddSpiderVisible = true
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '添加爬虫')
      },
      getMaxItemNodeId() {
        let max = 0
        this.spiderScrapyItems.forEach(d => {
          if (max < d.id) max = d.id
          d.children.forEach(f => {
            if (max < f.id) max = f.id
          })
        })
        return max
      },
      onAddItem() {
        const maxId = this.getMaxItemNodeId()
        this.spiderScrapyItems.push({
          id: maxId + 1,
          label: `Item_${+new Date()}`,
          level: 1,
          children: [
            {
              id: maxId + 2,
              level: 2,
              label: `field_${+new Date()}`
            }
          ]
        })
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '添加Item')
      },
      removeItem(data, ev) {
        ev.stopPropagation()
        for (let i = 0; i < this.spiderScrapyItems.length; i++) {
          const item = this.spiderScrapyItems[i]
          if (item.id === data.id) {
            this.spiderScrapyItems.splice(i, 1)
            break
          }
        }
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '删除Item')
      },
      onAddItemField(data, ev) {
        ev.stopPropagation()
        for (let i = 0; i < this.spiderScrapyItems.length; i++) {
          const item = this.spiderScrapyItems[i]
          if (item.id === data.id) {
            item.children.push({
              id: this.getMaxItemNodeId() + 1,
              level: 2,
              label: `field_${+new Date()}`
            })
            break
          }
        }
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '添加Items字段')
      },
      onRemoveItemField(node, data, ev) {
        ev.stopPropagation()
        for (let i = 0; i < this.spiderScrapyItems.length; i++) {
          const item = this.spiderScrapyItems[i]
          if (item.id === node.parent.data.id) {
            for (let j = 0; j < item.children.length; j++) {
              const field = item.children[j]
              if (field.id === data.id) {
                item.children.splice(j, 1)
                break
              }
            }
            break
          }
        }
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '删除Items字段')
      },
      onItemLabelEdit(node, data, ev) {
        ev.stopPropagation()
        this.$set(node, 'isEdit', true)
        this.$set(data, 'name', node.label)
        setTimeout(() => {
          this.$refs[`el-input-${data.id}`].focus()
        }, 0)
      },
      onItemChange(node, data, value) {
        for (let i = 0; i < this.spiderScrapyItems.length; i++) {
          const item = this.spiderScrapyItems[i]
          if (item.id === data.id) {
            item.label = value
            break
          }
        }
      },
      onItemFieldChange(node, data, value) {
        for (let i = 0; i < this.spiderScrapyItems.length; i++) {
          const item = this.spiderScrapyItems[i]
          if (item.id === node.parent.data.id) {
            for (let j = 0; j < item.children.length; j++) {
              const field = item.children[j]
              if (field.id === data.id) {
                item.children[j].label = value
                break
              }
            }
            break
          }
        }
      },
      async onItemsSave() {
        const res = await this.$store.dispatch('spider/saveSpiderScrapyItems', this.$route.params.id)
        if (!res.data.error) {
          this.$message.success(this.$t('Saved successfully'))
        }
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '保存Items')
      },
      async onClickSpider(spiderName) {
        if (this.loadingDict[spiderName]) return
        this.$set(this.loadingDict, spiderName, true)
        try {
          const res = await this.$store.dispatch('spider/getSpiderScrapySpiderFilepath', {
            id: this.$route.params.id,
            spiderName
          })
          this.$emit('click-spider', res.data.data)
        } finally {
          this.$set(this.loadingDict, spiderName, false)
        }
        this.$st.sendEv('爬虫详情', 'Scrapy 设置', '点击爬虫')
      },
      getScrapyErrors(type) {
        if (!this.spiderScrapyErrors || !this.spiderScrapyErrors[type] || (typeof this.spiderScrapyErrors[type] !== 'string')) return ''
        return this.$utils.html.htmlEscape(this.spiderScrapyErrors[type]).split('\n').join('<br/>')
      }
    }
  }
</script>

<style scoped>
  .spider-scrapy {
    height: calc(100vh - 200px);
    color: #606266;
  }

  .spider-scrapy >>> .el-tabs__content {
    overflow: auto;
  }

  .spider-scrapy >>> .el-tab-pane {
    height: calc(100vh - 239px);
  }

  .settings {
    width: 100%;
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

  .spiders {
    width: 100%;
    height: auto;
    overflow: auto;
  }

  .spiders .action-wrapper {
    text-align: right;
    padding-bottom: 10px;
    border-bottom: 1px solid #DCDFE6;
  }

  .pipelines .list,
  .spiders .list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .pipelines .list .item,
  .spiders .list .item {
    font-size: 14px;
    padding: 10px;
    cursor: pointer;
  }

  .pipelines .list .item:hover,
  .spiders .list .item:hover {
    background: #F5F7FA;
  }

  .items {
    width: 100%;
    height: auto;
    overflow: auto;
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

  .items >>> .el-tree > .el-tree-node {
    border-bottom: 1px solid #e6e9f0;
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

  .empty-text {
    display: block;
    margin-bottom: 20px;
  }

  .empty-text.error {
    color: #f56c6c;
  }

  .errors-label {
    color: #f56c6c;
    display: block;
    margin-bottom: 10px;
  }
</style>

