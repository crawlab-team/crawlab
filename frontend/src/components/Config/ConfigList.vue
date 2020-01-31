<template>
  <div class="config-list">
    <!--tour-->
    <v-tour
      name="spider-detail-config"
      :steps="tourSteps"
      :callbacks="tourCallbacks"
      :options="$utils.tour.getOptions(true)"
    />
    <!--./tour-->

    <!--preview results-->
    <el-dialog :visible.sync="dialogVisible"
               :title="$t('Preview Results')"
               width="90%"
               :before-close="onDialogClose">
      <el-table class="table-header" :data="[{}]" :show-header="false">
        <el-table-column v-for="(f, index) in fields"
                         :key="f.name + '-' + index"
                         min-width="100px">
          <template>
            <el-input v-model="columnsDict[f.name]" size="mini" style="width: calc(100% - 15px)"></el-input>
            <a href="javascript:" style="margin-left: 2px;" @click="onDeleteField(index)">X</a>
            <!--<el-button size="mini" type="danger" icon="el-icon-delete" style="width:45px;margin-left:2px"></el-button>-->
          </template>
        </el-table-column>
      </el-table>
      <el-table :data="previewCrawlData"
                :show-header="false"
                border>
        <el-table-column v-for="(f, index) in fields"
                         :key="f.name + '-' + index"
                         :label="f.name"
                         min-width="100px">

          <template slot-scope="scope">
            {{getDisplayStr(scope.row[f.name])}}
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
    <!--./preview results-->

    <!--crawl confirm dialog-->
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="spiderForm._id"
      @close="crawlConfirmDialogVisible = false"
    />
    <!--./crawl confirm dialog-->

    <!--tabs-->
    <el-tabs :active-name="activeTab" @tab-click="onTabClick">
      <!--Stages-->
      <el-tab-pane name="stages" :label="$t('Stages')">
        <!--config detail-->
        <el-row>
          <el-form label-width="150px" ref="form" :model="spiderForm.config">
          </el-form>
        </el-row>
        <!--./config detail-->

        <el-row>
          <div class="top-wrapper">
            <ul class="list">
              <li class="item">
                <label>{{$t('Start URL')}}: </label>
                <el-input
                  id="start-url"
                  v-model="spiderForm.config.start_url"
                  :placeholder="$t('Start URL')"
                  :class="startUrlClass"
                />
              </li>
              <li class="item">
                <label>{{$t('Start Stage')}}: </label>
                <el-select
                  id="start-stage"
                  v-model="spiderForm.config.start_stage"
                  :placeholder="$t('Start Stage')"
                  :class="startStageClass"
                  @change="$st.sendEv('爬虫详情', '配置', '改变起始阶段')"
                >
                  <el-option
                    v-for="n in spiderForm.config.stages.map(s => s.name)"
                    :key="n"
                    :value="n"
                    :label="n"
                  />
                </el-select>
              </li>
              <li class="item">
                <label>{{$t('Engine')}}: </label>
                <el-select
                  v-model="spiderForm.config.engine"
                  :placeholder="$t('Start Stage')"
                  :class="startStageClass"
                  disabled
                >
                  <el-option
                    v-for="n in ['scrapy']"
                    :key="n"
                    :value="n"
                    :label="n"
                  />
                </el-select>
              </li>
              <li class="item">
                <label>{{$t('Selector Type')}}: </label>
                <div class="selector-type">
              <span class="selector-type-item" @click="onClickSelectorType('css')">
                <el-tag
                  :class="isCss ? 'active' : 'inactive'"
                  type="success"
                >
                  CSS
                </el-tag>
              </span>
                  <span class="selector-type-item" @click="onClickSelectorType('xpath')">
              <el-tag
                :class="isXpath ? 'active' : 'inactive'"
                type="primary"
              >
                XPath
              </el-tag>
            </span>
                </div>
              </li>
            </ul>
          </div>

          <div class="button-group-container">
            <div class="button-group">
              <el-button id="btn-run" type="danger" @click="onCrawl">{{$t('Run')}}</el-button>
              <!--              <el-button type="primary" @click="onExtractFields" v-loading="extractFieldsLoading">-->
              <!--                {{$t('ExtractFields')}}-->
              <!--              </el-button>-->
              <!--              <el-button type="warning" @click="onPreview" v-loading="previewLoading">{{$t('Preview')}}</el-button>-->
              <el-button id="btn-save" type="success" @click="onSave" v-loading="saveLoading">{{$t('Save')}}</el-button>
            </div>
          </div>
        </el-row>

        <el-collapse
          :value="activeNames"
        >
          <el-collapse-item
            v-for="(stage, index) in spiderForm.config.stages"
            :key="index"
            :name="stage.name"
          >
            <template slot="title">
              <ul class="stage-list">
                <!--actions-->
                <li class="stage-item actions" style="min-width: 80px; flex-basis: 80px; justify-content: flex-end"
                    @click="$event.stopPropagation()">
                  <i class="action-item el-icon-copy-document" @click="onCopyStage(stage)"></i>
                  <i class="action-item el-icon-remove-outline" @click="onRemoveStage(stage)"></i>
                  <i class="action-item el-icon-circle-plus-outline" @click="onAddStage(stage)"></i>
                </li>
                <!--./actions-->

                <!--stage name-->
                <li class="stage-item" style="min-width: 240px" @click="$event.stopPropagation()">
                  <label>{{$t('Stage Name')}}: </label>
                  <div v-if="!stage.isEdit" @click="onEditStage(stage)" class="text-wrapper">
                    <span class="text">
                      {{stage.name}}
                    </span>
                    <i class="el-icon-edit-outline"></i>
                  </div>
                  <el-input
                    v-else
                    :ref="`stage-name-${stage.name}`"
                    class="edit-text"
                    v-model="stage.name"
                    :placeholder="$t('Stage Name')"
                    @focus="onStageNameFocus($event)"
                    @blur="stage.isEdit = false"
                  />
                </li>
                <!--./stage name-->

                <!--list-->
                <li class="stage-item" style="min-width: 240px">
                  <label>{{$t('List')}}: </label>
                  <el-checkbox
                    style="text-align: left; flex-basis: 20px; margin-right: 5px"
                    :value="isList(stage)"
                    @change="onCheckIsList($event, stage)"
                  />
                  <el-popover v-model="stage.isListOpen" v-if="isList(stage)" placement="top" width="360">
                    <el-form label-width="120px">
                      <el-form-item :label="$t('Selector Type')">
                        <el-tag :class="!stage.list_xpath ? 'active' : 'inactive'" type="success"
                                @click="onSelectStageListType(stage, 'css')">CSS
                        </el-tag>
                        <el-tag :class="stage.list_xpath ? 'active' : 'inactive'" type="primary"
                                @click="onSelectStageListType(stage, 'xpath')">XPath
                        </el-tag>
                      </el-form-item>
                      <el-form-item :label="$t('Selector')" class="list-selector">
                        <el-input v-if="!stage.list_xpath" v-model="stage.list_css"/>
                        <el-input v-else v-model="stage.list_xpath"/>
                      </el-form-item>
                    </el-form>
                    <el-tag
                      v-if="!stage.list_xpath"
                      type="success"
                      slot="reference"
                      @click="onClickStageList($event, stage, 'css')"
                    >
                      <i v-if="!stage.isListOpen" class="el-icon-circle-plus-outline"></i>
                      <i v-else class="el-icon-remove-outline"></i>
                      CSS
                    </el-tag>
                    <el-tag
                      v-else
                      type="primary"
                      slot="reference"
                      @click="onClickStageList($event, stage, 'xpath')"
                    >
                      <i v-if="!stage.isListOpen" class="el-icon-circle-plus-outline"></i>
                      <i v-else class="el-icon-remove-outline"></i>
                      XPath
                    </el-tag>
                  </el-popover>
                </li>
                <!--./list-->

                <!--pagination-->
                <li class="stage-item" style="min-width: 240px">
                  <label>{{$t('Pagination')}}: </label>
                  <el-checkbox
                    style="text-align: left; flex-basis: 20px; margin-right: 5px"
                    :value="isPage(stage)"
                    @change="onCheckIsPage($event, stage)"
                    :disabled="!isList(stage)"
                  />
                  <el-popover v-model="stage.isPageOpen" v-if="isPage(stage)" placement="top" width="360">
                    <el-form label-width="120px">
                      <el-form-item :label="$t('Selector Type')">
                        <el-tag :class="!stage.page_xpath ? 'active' : 'inactive'" type="success"
                                @click="onSelectStagePageType(stage, 'css')">CSS
                        </el-tag>
                        <el-tag :class="stage.page_xpath ? 'active' : 'inactive'" type="primary"
                                @click="onSelectStagePageType(stage, 'xpath')">XPath
                        </el-tag>
                      </el-form-item>
                      <el-form-item :label="$t('Selector')" class="page-selector">
                        <el-input v-if="!stage.page_xpath" v-model="stage.page_css"/>
                        <el-input v-else v-model="stage.page_xpath"/>
                      </el-form-item>
                    </el-form>
                    <el-tag
                      v-if="!stage.page_xpath"
                      type="success"
                      slot="reference"
                      @click="onClickStagePage($event, stage, 'css')"
                    >
                      <i v-if="!stage.isPageOpen" class="el-icon-circle-plus-outline"></i>
                      <i v-else class="el-icon-remove-outline"></i>
                      CSS
                    </el-tag>
                    <el-tag
                      v-else
                      type="primary"
                      slot="reference"
                      @click="onClickStagePage($event, stage, 'xpath')"
                    >
                      <i v-if="!stage.isPageOpen" class="el-icon-circle-plus-outline"></i>
                      <i v-else class="el-icon-remove-outline"></i>
                      XPath
                    </el-tag>
                  </el-popover>
                </li>
                <!--./pagination-->

              </ul>
            </template>
            <fields-table-view
              type="list"
              title="List Page Fields"
              :fields="stage.fields"
              :stage="stage"
              :stage-names="spiderForm.config.stages.map(s => s.name)"
            />
          </el-collapse-item>
        </el-collapse>
      </el-tab-pane>
      <!--./Stages-->

      <!--Graph-->
      <el-tab-pane name="process" :label="$t('Process')">
        <div id="process-chart"></div>
      </el-tab-pane>
      <!--./Graph-->

      <!--Setting-->
      <el-tab-pane name="settings" :label="$t('Settings')">
        <div class="actions" style="text-align: right;margin-bottom: 10px">
          <el-button type="success" size="small" @click="onSave">
            {{$t('Save')}}
          </el-button>
        </div>
        <setting-fields-table-view
          type="list"
        />
      </el-tab-pane>
      <!--./Setting-->

      <!--Spiderfile-->
      <el-tab-pane name="spiderfile" label="Spiderfile">
        <div class="spiderfile-actions">
          <el-button type="primary" size="small" style="margin-right: 10px;" @click="onSpiderfileSave">
            <font-awesome-icon :icon="['fa', 'save']"/>
            {{$t('Save')}}
          </el-button>
        </div>
        <file-detail/>
      </el-tab-pane>
      <!--./Spiderfile-->
    </el-tabs>
    <!--./tabs-->
  </div>
</template>

<script>
import { mapState } from 'vuex'
import echarts from 'echarts'
import FieldsTableView from '../TableView/FieldsTableView'
import CrawlConfirmDialog from '../Common/CrawlConfirmDialog'

import 'codemirror/lib/codemirror.js'
import 'codemirror/mode/yaml/yaml.js'
import FileDetail from '../File/FileDetail'
import SettingFieldsTableView from '../TableView/SettingFieldsTableView'

export default {
  name: 'ConfigList',
  components: {
    SettingFieldsTableView,
    FileDetail,
    CrawlConfirmDialog,
    FieldsTableView
  },
  watch: {
    activeTab () {
      setTimeout(() => {
        // 渲染流程图
        this.renderProcessChart()

        // 获取Spiderfile
        this.getSpiderfile()

        // 获取config
        this.$store.dispatch('spider/getSpiderData', this.spiderForm._id)
      }, 0)
    }
  },
  data () {
    return {
      crawlTypeList: [
        { value: 'list', label: 'List Only' },
        { value: 'detail', label: 'Detail Only' },
        { value: 'list-detail', label: 'List + Detail' }
      ],
      extractFieldsLoading: false,
      previewLoading: false,
      saveLoading: false,
      dialogVisible: false,
      crawlConfirmDialogVisible: false,
      columnsDict: {},
      fieldColumns: [
        { name: 'name', label: 'Name' },
        { name: 'selector_type', label: 'Selector Type' },
        { name: 'selector', label: 'Selector' },
        { name: 'is_attr', label: 'Is Attribute' },
        { name: 'attr', label: 'Attribute' },
        { name: 'next_stage', label: 'Next Stage' }
      ],
      activeTab: 'stages',
      processChart: undefined,
      fileOptions: {
        mode: 'text/x-yaml',
        theme: 'darcula',
        styleActiveLine: true,
        lineNumbers: true,
        line: true,
        matchBrackets: true
      },
      tourSteps: [
        // stage
        {
          target: '.config-list .el-tabs__nav.is-top',
          content: this.$t('You can switch to each section of configurable spider.')
        },
        {
          target: '#start-url',
          content: this.$t('Here is the starting URL of the spider.')
        },
        {
          target: '#start-stage',
          content: this.$t('Here is the starting stage of the spider.<br><br>A <strong>Stage</strong> is basically a callback in the Scrapy spider.')
        },
        {
          target: '#btn-run',
          content: this.$t('You can run a spider task.<br><br>Spider will be automatically saved when clicking on this button.')
        },
        {
          target: '.stage-item.actions',
          content: this.$t('Add/duplicate/delete a stage.'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '.fields-table-view td.action',
          content: this.$t('Add/duplicate/delete an extract field in the stage.'),
          params: {
            placement: 'right'
          }
        },
        {
          target: '.stage-item:nth-child(3)',
          content: this.$t('You can decide whether this is a list page.<br><br>Click on the CSS/XPath tag to enter the selector expression for list items.<br>For example, "<code>ul > li</code>"'),
          params: {
            placement: 'top'
          }
        },
        {
          target: '.stage-item:nth-child(4)',
          content: this.$t('You can decide whether this is a list page with pagination.<br><br>Click on the CSS/XPath tag to enter the selector expression for the pagination.<br>For example, "<code>a.next</code>"'),
          params: {
            placement: 'top'
          }
        },
        {
          target: '.fields-table-view',
          content: this.$t('You should enter necessary information for all fields in the stage.'),
          params: {
            placement: 'top'
          }
        },
        {
          target: '.fields-table-view tr:nth-child(1) td:nth-child(7)',
          content: this.$t('If you have multiple stages, e.g. list page + detail page, you should select the next stage in the detail link\'s field.'),
          params: {
            placement: 'top'
          }
        },
        // process
        {
          target: '#tab-process',
          content: this.$t('You can view the<br> visualization of the stage<br> workflow.')
        },
        // settings
        {
          target: '#tab-settings',
          content: this.$t('You can add the settings here, which will be loaded in the Scrapy\'s <code>settings.py</code> file.<br><br>JSON and Array data are supported.')
        },
        // Spiderfile
        {
          target: '#tab-spiderfile',
          content: this.$t('You can edit the <code>Spiderfile</code> here.<br><br>For more information, please refer to the <a href="https://docs.crawlab.cn/Usage/Spider/ConfigurableSpider.html" target="_blank" style="color: #409EFF">Documentation (Chinese)</a>.')
        }
      ],
      tourCallbacks: {
        onStop: () => {
          this.$utils.tour.finishTour('spider-detail-config')
        },
        onPreviousStep: (currentStep) => {
          if (currentStep === 10) {
            this.activeTab = 'stages'
          } else if (currentStep === 11) {
            this.activeTab = 'process'
          } else if (currentStep === 12) {
            this.activeTab = 'settings'
          }
        },
        onNextStep: (currentStep) => {
          if (currentStep === 9) {
            this.activeTab = 'process'
          } else if (currentStep === 10) {
            this.activeTab = 'settings'
          } else if (currentStep === 11) {
            this.activeTab = 'spiderfile'
          }
        }
      }
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm',
      'previewCrawlData'
    ]),
    fields () {
      if (this.spiderForm.crawl_type === 'list') {
        return this.spiderForm.fields
      } else if (this.spiderForm.crawl_type === 'detail') {
        return this.spiderForm.detail_fields
      } else if (this.spiderForm.crawl_type === 'list-detail') {
        return this.spiderForm.fields.concat(this.spiderForm.detail_fields)
      } else {
        return []
      }
    },
    isCss () {
      let i = 0
      this.spiderForm.config.stages.forEach(stage => {
        stage.fields.forEach(field => {
          if (!field.css) i++
        })
      })
      return i === 0
    },
    isXpath () {
      let i = 0
      this.spiderForm.config.stages.forEach(stage => {
        stage.fields.forEach(field => {
          if (!field.xpath) i++
        })
      })
      return i === 0
    },
    activeNames () {
      return this.spiderForm.config.stages.map(d => d.name)
    },
    startUrlClass () {
      if (!this.spiderForm.config.start_url) {
        return 'invalid'
      } else if (!this.spiderForm.config.start_url.match(/^https?:\/\/.+|^\/\/.+/i)) {
        return 'invalid'
      }

      return ''
    },
    startStageClass () {
      if (!this.spiderForm.config.start_stage) {
        return 'invalid'
      } else if (!this.activeNames.includes(this.spiderForm.config.start_stage)) {
        return 'invalid'
      }
      return ''
    },
    stageNodes () {
      const startStage = this.spiderForm.config.stages[this.spiderForm.config.start_stage]
      const nodes = []
      const allStageNames = new Set()

      let i = 0
      let currentStage = startStage
      while (currentStage) {
        // 加入节点信息
        nodes.push({
          x: i++,
          y: 0,
          itemStyle: {
            color: '#409EFF'
          },
          ...currentStage
        })

        // 记录该节点
        allStageNames.add(currentStage.name)

        // 设置当前阶段为下一阶段
        currentStage = this.getNextStage(currentStage)

        if (currentStage && allStageNames.has(currentStage.name)) {
          currentStage = undefined
        }
      }

      // 加入剩余阶段
      i = 0
      const restStages = this.spiderForm.config.stages
        .filter(stage => !allStageNames.has(stage.name))
      restStages.forEach(stage => {
        // 加入节点信息
        nodes.push({
          x: i++,
          y: 1,
          itemStyle: {
            color: '#F56C6C'
          },
          ...stage
        })
      })

      return nodes
    },
    stageEdges () {
      const edges = []
      const stages = this.spiderForm.config.stages
      stages.forEach(stage => {
        for (let i = 0; i < stage.fields.length; i++) {
          const field = stage.fields[i]
          if (field.next_stage) {
            edges.push({
              source: stage.name,
              target: field.next_stage
            })
          }
        }
      })
      return edges
    }
  },
  methods: {
    onSelectCrawlType (value) {
      this.spiderForm.crawl_type = value
    },
    async onSave () {
      this.$st.sendEv('爬虫详情', '配置', '保存')
      this.saveLoading = true
      try {
        const res = await this.$store.dispatch('spider/postConfigSpiderConfig')
        if (!res.data.error) {
          this.$message.success(this.$t('Spider info has been saved successfully'))
          return true
        }
        return false
      } catch (e) {
        return false
      } finally {
        this.saveLoading = false
      }
    },
    onDialogClose () {
      this.dialogVisible = false
      this.fields.forEach(f => {
        f.name = this.columnsDict[f.name]
      })
    },
    onPreview () {
      this.$refs['form'].validate(res => {
        if (res) {
          this.onSave()
            .then(() => {
              this.previewLoading = true
              this.$store.dispatch('spider/getPreviewCrawlData')
                .then(() => {
                  this.fields.forEach(f => {
                    this.columnsDict[f.name] = f.name
                  })
                  this.dialogVisible = true
                })
                .catch(() => {
                  this.$message.error(this.$t('Something wrong happened'))
                })
                .finally(() => {
                  this.previewLoading = false
                })
              this.$st.sendEv('爬虫详情', '配置', '预览')
            })
        }
      })
    },
    async onCrawl () {
      this.$st.sendEv('爬虫详情', '配置', '点击运行')
      const res = await this.onSave()
      if (res) {
        this.crawlConfirmDialogVisible = true
      }
    },
    onExtractFields () {
      this.$refs['form'].validate(res => {
        if (res) {
          this.onSave()
            .then(() => {
              this.extractFieldsLoading = true
              this.$store.dispatch('spider/extractFields')
                .then(response => {
                  if (response.data.item_selector) {
                    this.$set(this.spiderForm, 'item_selector', response.data.item_selector)
                  }
                  if (response.data.item_selector_type) {
                    this.$set(this.spiderForm, 'item_selector_type', response.data.item_selector_type)
                  }

                  if (response.data.fields && response.data.fields.length) {
                    this.spiderForm.fields = response.data.fields
                  }

                  if (response.data.pagination_selector) {
                    this.spiderForm.pagination_selector = response.data.pagination_selector
                  }
                })
                .finally(() => {
                  this.extractFieldsLoading = false
                })
              this.$st.sendEv('爬虫详情', '配置', '提取字段')
            })
        }
      })
    },
    onDeleteField (index) {
      this.fields.splice(index, 1)
    },
    getDisplayStr (value) {
      if (!value) return value
      value = value.trim()
      if (value.length > 20) return value.substr(0, 20) + '...'
      return value
    },
    onClickSelectorType (selectorType) {
      this.$st.sendEv('爬虫详情', '配置', `点击阶段选择器类别-${selectorType}`)
      this.spiderForm.config.stages.forEach(stage => {
        // 列表
        if (selectorType === 'css') {
          if (stage.list_xpath) stage.list_xpath = ''
          if (!stage.list_css) stage.list_css = 'body'
        } else {
          if (stage.list_css) stage.list_css = ''
          if (!stage.list_xpath) stage.list_xpath = '//body'
        }

        // 分页
        if (selectorType === 'css') {
          if (stage.page_xpath) stage.page_xpath = ''
          if (!stage.page_css) stage.page_css = 'body'
        } else {
          if (stage.page_css) stage.page_css = ''
          if (!stage.page_xpath) stage.page_xpath = '//body'
        }

        // 字段
        stage.fields.forEach(field => {
          if (selectorType === 'css') {
            if (field.xpath) field.xpath = ''
            if (!field.css) field.css = 'body'
          } else {
            if (field.css) field.css = ''
            if (!field.xpath) field.xpath = '//body'
          }
        })
      })
    },
    onStageNameFocus (ev) {
      ev.stopPropagation()
    },
    onEditStage (stage) {
      this.$st.sendEv('爬虫详情', '配置', '更改阶段名称')
      this.$set(stage, 'isEdit', !stage.isEdit)
      setTimeout(() => {
        this.$refs[`stage-name-${stage.name}`][0].focus()
      }, 0)
    },
    onCopyStage (stage) {
      this.$st.sendEv('爬虫详情', '配置', '复制阶段')
      const stages = this.spiderForm.config.stages
      const ts = Math.floor(new Date().getTime()).toString()
      const newStage = JSON.parse(JSON.stringify(stage))
      newStage.name = `${stage.name}_copy_${ts}`
      for (let i = 0; i < stages.length; i++) {
        if (stage.name === stages[i].name) {
          stages.splice(i + 1, 0, newStage)
        }
      }
    },
    addStage (index) {
      const stages = this.spiderForm.config.stages
      const ts = Math.floor(new Date().getTime()).toString()
      const newStageName = `stage_${ts}`
      const newField = { name: `field_${ts}`, next_stage: '' }
      if (this.isCss) {
        newField['css'] = 'body'
      } else if (this.isXpath) {
        newField['xpath'] = '//body'
      } else {
        newField['xpath'] = '//body'
      }
      stages.splice(index + 1, 0, {
        name: newStageName,
        list_css: this.isCss ? 'body' : '',
        list_xpath: this.isXpath ? '//body' : '',
        page_css: '',
        page_xpath: '',
        fields: [newField]
      })
    },
    onRemoveStage (stage) {
      this.$st.sendEv('爬虫详情', '配置', '删除阶段')
      const stages = this.spiderForm.config.stages
      for (let i = 0; i < stages.length; i++) {
        if (stage.name === stages[i].name) {
          stages.splice(i, 1)
          break
        }
      }
      // 如果只剩一个stage，加入新的stage
      if (stages.length === 0) {
        this.addStage(0)
      }
      // 重置next_stage被设置为该stage的field
      stages.forEach(_stage => {
        _stage.fields.forEach(field => {
          if (field.next_stage === stage.name) {
            this.$set(field, 'next_stage', '')
          }
        })
      })
    },
    onAddStage (stage) {
      this.$st.sendEv('爬虫详情', '配置', '添加阶段')
      const stages = this.spiderForm.config.stages
      for (let i = 0; i < stages.length; i++) {
        if (stage.name === stages[i].name) {
          this.addStage(i)
          break
        }
      }
    },
    renderProcessChart () {
      const option = {
        title: {
          text: this.$t('Stage Process')
        },
        series: [
          {
            animation: false,
            type: 'graph',
            // layout: 'force',
            layout: 'none',
            symbolSize: 50,
            roam: true,
            label: {
              normal: {
                show: true
              }
            },
            edgeSymbol: ['circle', 'arrow'],
            edgeSymbolSize: [4, 10],
            edgeLabel: {
              normal: {
                textStyle: {
                  fontSize: 20
                }
              }
            },
            focusOneNodeAdjacency: true,
            force: {
              initLayout: 'force',
              repulsion: 100,
              gravity: 0.00001,
              edgeLength: 200
            },
            // draggable: true,
            data: this.stageNodes,
            links: this.stageEdges,
            lineStyle: {
              normal: {
                opacity: 0.9,
                width: 2,
                curveness: 0
              }
            }
          }
        ],
        tooltip: {
          // formatter: '{b0}: {c0}<br />{b1}: {c1}',
          formatter: (params) => {
            if (!params.data.fields) return

            let str = ''
            str += `<label>${this.$t('Stage')}: ${params.name}</label><br>`
            str += `<ul style="list-style: none; padding: 0; margin: 0;">`
            // 列表
            if (params.data.list_css || params.data.list_xpath) {
              str += `<li><span style="display: inline-block;min-width: 50px;font-weight: bolder;text-align: right;margin-right: 3px">${this.$t('List')}: </span>${params.data.list_css || params.data.list_xpath}</li>`
            }
            if (params.data.page_css || params.data.page_xpath) {
              str += `<li><span style="display: inline-block;min-width: 50px;font-weight: bolder;text-align: right;margin-right: 3px">${this.$t('Pagination')}: </span>${params.data.page_css || params.data.page_xpath}</li>`
            }
            str += `</ul><br>`

            // 字段
            str += `<label>${this.$t('Fields')}: </label><br>`
            str += '<ul style="list-style: none; padding: 0; margin: 0;">'
            for (let i = 0; i < params.data.fields.length; i++) {
              const f = params.data.fields[i]
              str += `
<li>
<span style="display: inline-block; min-width: 50px; font-weight: bolder; text-align: right">${f.name}: </span>
${f.css || f.xpath} ${f.attr ? ('(' + f.attr + ')') : ''} ${f.next_stage ? (' --> ' + '<span style="font-weight:bolder">' + f.next_stage + '</span>') : ''}
</li>
`
            }
            str += '</ul>'
            return str
          }
        }
      }
      const el = document.querySelector('#process-chart')
      this.processChart = echarts.init(el)
      this.processChart.setOption(option)
      this.processChart.resize()
    },
    onTabClick (tab) {
      this.activeTab = tab.name
    },
    update () {
      if (this.activeTab !== 'stages') return

      // 手动显示tab下划线
      const elBar = document.querySelector('.el-tabs__active-bar')
      const elStages = document.querySelector('#tab-stages')
      const totalWidth = Number(getComputedStyle(elStages).width.replace('px', ''))
      const paddingRight = Number(getComputedStyle(elStages).paddingRight.replace('px', ''))
      elBar.setAttribute('style', 'width:' + (totalWidth - paddingRight) + 'px')
    },
    getSpiderfile () {
      this.$store.commit('file/SET_FILE_CONTENT', '')
      this.$store.commit('file/SET_CURRENT_PATH', 'Spiderfile')
      this.$store.dispatch('file/getFileContent', { path: 'Spiderfile' })
    },
    async onSpiderfileSave () {
      try {
        await this.$store.dispatch('spider/saveConfigSpiderSpiderfile')
        this.$message.success(this.$t('Spiderfile saved successfully'))
      } catch (e) {
        this.$message.error('Something wrong happened')
      }
    },
    isList (stage) {
      return !!stage.is_list
    },
    onCheckIsList (value, stage) {
      stage.is_list = value
      if (value) {
        this.$st.sendEv('爬虫详情', '配置', '勾选列表页')
      } else {
        this.$st.sendEv('爬虫详情', '配置', '取消勾选列表页')
      }
    },
    onClickStageList ($event, stage, type) {
      $event.stopPropagation()
    },
    onSelectStageListType (stage, type) {
      if (type === 'css') {
        if (!stage.list_css) stage.list_css = 'body'
        stage.list_xpath = ''
      } else {
        if (!stage.list_xpath) stage.list_xpath = '//body'
        stage.list_css = ''
      }
    },
    isPage (stage) {
      return !!stage.page_css || !!stage.page_xpath
    },
    onCheckIsPage (value, stage) {
      if (value) {
        this.$st.sendEv('爬虫详情', '配置', '勾选分页')
        if (!stage.page_css && !stage.page_xpath) {
          stage.page_xpath = '//body'
        }
      } else {
        this.$st.sendEv('爬虫详情', '配置', '取消勾选分页')
        stage.page_css = ''
        stage.page_xpath = ''
      }
    },
    onClickStagePage ($event, stage, type) {
      $event.stopPropagation()
    },
    onSelectStagePageType (stage, type) {
      if (type === 'css') {
        if (!stage.page_css) stage.page_css = 'body'
        stage.page_xpath = ''
      } else {
        if (!stage.page_xpath) stage.page_xpath = '//body'
        stage.page_css = ''
      }
    },
    getNextStageField (stage) {
      return stage.fields.filter(f => !!f.next_stage)[0]
    },
    getNextStage (stage) {
      const nextStageField = this.getNextStageField(stage)
      if (!nextStageField) return
      return this.spiderForm.config.stages[nextStageField.next_stage]
    }
  },
  mounted () {
    this.activeNames = this.spiderForm.config.stages.map(stage => stage.name)
  }
}
</script>

<style scoped>

  .button-group-container {
    margin-top: 10px;
    /*border-bottom: 1px dashed #dcdfe6;*/
    padding-bottom: 20px;
  }

  .button-group {
    text-align: right;
  }

  .list-fields-container {
    margin-top: 20px;
    /*border-bottom: 1px dashed #dcdfe6;*/
    padding-bottom: 20px;
  }

  .detail-fields-container {
    margin-top: 20px;
  }

  .title {
    color: #606266;
    font-size: 14px;
  }

  .el-table.table-header >>> td {
    padding: 0;
  }

  .el-table.table-header >>> .cell {
    padding: 0;
  }

  .el-table.table-header >>> .el-input .el-input__inner {
    border-radius: 0;
  }

  .selector-type-item {
    margin: 0 5px;
    cursor: pointer;
    font-weight: bolder;
  }

  .el-tag {
    margin-right: 5px;
    font-weight: bolder;
    cursor: pointer;
  }

  .el-tag.inactive {
    opacity: 0.5;
  }

  .stage-list {
    width: 100%;
    /*width: calc(80px + 320px);*/
    display: flex;
    flex-wrap: wrap;
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .stage-list .stage-item {
    /*flex-basis: 320px;*/
    min-width: 120px;
    display: flex;
    align-items: center;
  }

  .stage-list .stage-item label {
    flex-basis: 90px;
    margin-right: 10px;
    justify-self: flex-end;
    text-align: right;
  }

  .stage-list .stage-item .el-input {
    flex-basis: calc(100% - 90px);
    height: 32px;
  }

  .stage-list .stage-item .el-input .el-input__inner {
    height: 32px;
    inline-size: 32px;
  }

  .stage-list .stage-item .action-item {
    cursor: pointer;
    width: 13px;
    margin-right: 5px;
  }

  .stage-list .stage-item .action-item:last-child {
    margin-right: 10px;
  }

  .stage-list .stage-item .text-wrapper {
    display: flex;
    align-items: center;
    max-width: calc(100% - 90px - 10px);
  }

  .stage-list .stage-item .text-wrapper .text {
    text-overflow: ellipsis;
    overflow: hidden;
  }

  .stage-list .stage-item .text-wrapper .text:hover {
    text-decoration: underline;
  }

  .stage-list .stage-item .text-wrapper i {
    margin-left: 5px;
  }

  .stage-list .stage-item >>> .edit-text {
    height: 32px;
    line-height: 32px;
  }

  .stage-list .stage-item >>> .edit-text .el-input__inner {
    height: 32px;
    line-height: 32px;
  }

  .top-wrapper {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .top-wrapper .list {
    list-style: none;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    padding: 0;
  }

  .top-wrapper .list .item {
    margin-bottom: 10px;
    display: flex;
    align-items: center;
  }

  .top-wrapper .list .item label {
    width: 100px;
    text-align: right;
    margin-right: 10px;
    font-size: 12px;
  }

  .top-wrapper .list .item label + * {
    width: 240px;
  }

  .invalid >>> .el-input__inner {
    border: 1px solid red !important;
  }

  #process-chart {
    width: 100%;
    height: 480px;
  }

  .config-list >>> .file-content {
    height: calc(100vh - 280px);
  }

  .spiderfile-actions {
    margin-bottom: 5px;
    text-align: right;
  }
</style>
