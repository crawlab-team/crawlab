<template>
  <div class="config-list">
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
                  v-model="spiderForm.config.start_url"
                  :placeholder="$t('Start URL')"
                  :class="startUrlClass"
                />
              </li>
              <li class="item">
                <label>{{$t('Start Stage')}}: </label>
                <el-select
                  v-model="spiderForm.config.start_stage"
                  :placeholder="$t('Start Stage')"
                  :class="startStageClass"
                >
                  <el-option
                    v-for="n in Object.keys(spiderForm.config.stages)"
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
              <el-button type="danger" @click="onCrawl">{{$t('Run')}}</el-button>
              <el-button type="primary" @click="onExtractFields" v-loading="extractFieldsLoading">
                {{$t('ExtractFields')}}
              </el-button>
              <el-button type="warning" @click="onPreview" v-loading="previewLoading">{{$t('Preview')}}</el-button>
              <el-button type="success" @click="onSave" v-loading="saveLoading">{{$t('Save')}}</el-button>
            </div>
          </div>
        </el-row>

        <el-collapse
          :value="activeNames"
        >
          <el-collapse-item
            v-for="(stage, stageName) in spiderForm.config.stages"
            :key="stageName"
            :name="stageName"
          >
            <template slot="title">
              <ul class="stage-list">
                <li class="stage-item" style="flex-basis: 80px; justify-content: flex-end"
                    @click="$event.stopPropagation()">
                  <i class="action-item el-icon-copy-document" @click="onCopyStage(stage)"></i>
                  <i class="action-item el-icon-remove-outline" @click="onRemoveStage(stage)"></i>
                  <i class="action-item el-icon-circle-plus-outline" @click="onAddStage(stage)"></i>
                </li>
                <li class="stage-item" @click="$event.stopPropagation()">
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
              </ul>
            </template>
            <fields-table-view
              type="list"
              title="List Page Fields"
              :fields="stage.fields"
              :stage-names="Object.keys(spiderForm.config.stages)"
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
    </el-tabs>
    <!--./tabs-->
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import echarts from 'echarts'
import FieldsTableView from '../TableView/FieldsTableView'
import CrawlConfirmDialog from '../Common/CrawlConfirmDialog'

export default {
  name: 'ConfigList',
  components: {
    CrawlConfirmDialog,
    FieldsTableView
  },
  watch: {
    activeTab () {
      setTimeout(() => {
        this.renderProcessChart()
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
      processChart: undefined
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
      Object.values(this.spiderForm.config.stages).forEach(stage => {
        stage.fields.forEach(field => {
          if (!field.css) i++
        })
      })
      return i === 0
    },
    isXpath () {
      let i = 0
      Object.values(this.spiderForm.config.stages).forEach(stage => {
        stage.fields.forEach(field => {
          if (!field.xpath) i++
        })
      })
      return i === 0
    },
    activeNames () {
      return Object.values(this.spiderForm.config.stages).map(d => d.name)
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
      // const elChart = document.querySelector('#process-chart')
      // const totalWidth = Number(getComputedStyle(elChart).width.replace('px', ''))
      const stages = Object.values(this.spiderForm.config.stages)
      return stages.map((stage, i) => {
        return {
          name: stage.name,
          ...stage
        }
      })
    },
    stageEdges () {
      const edges = []
      const stages = Object.values(this.spiderForm.config.stages)
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
    onSave () {
      this.$st.sendEv('爬虫详情-配置', '保存')
      return new Promise((resolve, reject) => {
        this.saveLoading = true
        this.$store.dispatch('spider/updateSpiderFields')
          .then(() => {
            this.$store.dispatch('spider/editSpider')
              .then(() => {
                this.$message.success(this.$t('Spider info has been saved successfully'))
                resolve()
              })
              .catch(() => {
                this.$message.error(this.$t('Something wrong happened'))
                reject(new Error())
              })
              .finally(() => {
                this.saveLoading = false
              })
          })
          .then(() => {
            this.$store.dispatch('spider/updateSpiderDetailFields')
          })
          .catch(() => {
            this.$message.error(this.$t('Something wrong happened'))
            this.saveLoading = false
            reject(new Error())
          })
      })
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
              this.$st.sendEv('爬虫详情-配置', '预览')
            })
        }
      })
    },
    onCrawl () {
      this.$confirm(this.$t('Are you sure to run this spider?'), this.$t('Notification'), {
        confirmButtonText: this.$t('Confirm'),
        cancelButtonText: this.$t('Cancel')
      })
        .then(() => {
          this.$store.dispatch('spider/crawlSpider', this.spiderForm._id)
            .then(() => {
              this.$message.success(this.$t(`Spider task has been scheduled`))
            })
          this.$st.sendEv('爬虫详情-配置', '运行')
        })
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
              this.$st.sendEv('爬虫详情-配置', '提取字段')
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
      Object.values(this.spiderForm.config.stages).forEach(stage => {
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
      console.log(ev)
      ev.stopPropagation()
      // ev.preventDefault()
    },
    onEditStage (stage) {
      this.$set(stage, 'isEdit', !stage.isEdit)
      setTimeout(() => {
        this.$refs[`stage-name-${stage.name}`][0].focus()
      }, 0)
    },
    onCopyStage (stage) {
      const ts = Math.floor(new Date().getTime()).toString()
      const newStage = JSON.parse(JSON.stringify(stage))
      newStage.name = `stage_${ts}`
      this.$set(this.spiderForm.config.stages, newStage.name, newStage)
    },
    onRemoveStage (stage) {
      const stages = JSON.parse(JSON.stringify(this.spiderForm.config.stages))
      delete stages[stage.name]
      this.$set(this.spiderForm.config, 'stages', stages)
      if (Object.keys(stages).length === 0) {
        this.onAddStage()
      }
    },
    onAddStage (stage) {
      const stages = JSON.parse(JSON.stringify(this.spiderForm.config.stages))
      const ts = Math.floor(new Date().getTime()).toString()
      const newStageName = `stage_${ts}`
      const newField = { name: `field_${ts}`, next_stage: '' }
      if (this.isCss) {
        newField['css'] = 'body'
      } else if (this.isXpath) {
        newField['xpath'] = '//body'
      } else {
        newField['css'] = 'body'
      }
      stages[newStageName] = {
        name: newStageName,
        fields: [newField]
      }
      this.$set(this.spiderForm.config, 'stages', stages)
    },
    renderProcessChart () {
      const option = {
        title: {
          text: this.$t('Stage Process')
        },
        series: [
          {
            type: 'graph',
            layout: 'force',
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
              gravity: 0.01,
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
        ]
      }
      const el = document.querySelector('#process-chart')
      this.processChart = echarts.init(el)
      this.processChart.setOption(option)
      this.processChart.resize()
    },
    onTabClick (tab) {
      this.activeTab = tab.name
    }
  },
  created () {
    // fields for list page
    if (!this.spiderForm.fields) {
      this.spiderForm.fields = []
      for (let i = 0; i < 3; i++) {
        this.spiderForm.fields.push({
          name: 'field_' + (i + 1),
          type: 'css',
          extract_type: 'text'
        })
      }
    }

    // fields for detail page
    if (!this.spiderForm.detail_fields) {
      this.spiderForm.detail_fields = []
      for (let i = 0; i < 3; i++) {
        this.spiderForm.detail_fields.push({
          name: 'field_' + (i + 1),
          type: 'css',
          extract_type: 'text'
        })
      }
    }

    if (!this.spiderForm.crawl_type) this.$set(this.spiderForm, 'crawl_type', 'list')
    // if (!this.spiderForm.start_url) this.$set(this.spiderForm, 'start_url', 'http://example.com')
    if (!this.spiderForm.item_selector_type) this.$set(this.spiderForm, 'item_selector_type', 'css')
    if (!this.spiderForm.pagination_selector_type) this.$set(this.spiderForm, 'pagination_selector_type', 'css')
    if (this.spiderForm.obey_robots_txt == null) this.$set(this.spiderForm, 'obey_robots_txt', true)
    if (this.spiderForm.item_threshold == null) this.$set(this.spiderForm, 'item_threshold', 10)
  },
  mounted () {
    this.activeNames = Object.keys(this.spiderForm.config.stages)
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

  .selector-type-item > .el-tag.inactive {
    opacity: 0.5;
  }

  .stage-list {
    width: calc(80px + 320px);
    display: flex;
    flex-wrap: wrap;
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .stage-list .stage-item {
    flex-basis: 320px;
    width: 320px;
    display: flex;
    align-items: center;
  }

  .stage-list .stage-item label {
    flex-basis: 90px;
    margin-right: 10px;
    justify-self: flex-end;
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
</style>
