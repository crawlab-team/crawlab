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

    <!--config detail-->
    <el-row>
      <el-form label-width="150px" ref="form" :model="spiderForm">
        <el-col :span="11" :offset="1">
          <el-form-item :label="$t('Crawl Type')">
            <el-button-group>
              <el-button v-for="type in crawlTypeList"
                         :key="type.value"
                         :type="type.value === spiderForm.crawl_type ? 'primary' : ''"
                         @click="onSelectCrawlType(type.value)">
                {{$t(type.label)}}
              </el-button>
            </el-button-group>
          </el-form-item>
          <el-form-item :label="$t('Start URL')" prop="start_url" required>
            <el-input v-model="spiderForm.start_url" :placeholder="$t('Start URL')"></el-input>
          </el-form-item>
          <el-form-item :label="$t('Obey robots.txt')">
            <el-switch v-model="spiderForm.obey_robots_txt" :placeholder="$t('Obey robots.txt')"></el-switch>
          </el-form-item>
          <!--<el-form-item :label="$t('URL Pattern')">-->
          <!--<el-input v-model="spiderForm.url_pattern" :placeholder="$t('URL Pattern')"></el-input>-->
          <!--</el-form-item>-->
        </el-col>
        <el-col :span="11" :offset="1">
          <el-form-item :label="$t('Item Selector')"
                        v-if="['list','list-detail'].includes(spiderForm.crawl_type)">
            <el-select style="width: 35%;margin-right: 10px;"
                       v-model="spiderForm.item_selector_type"
                       :placeholder="$t('Item Selector Type')">
              <el-option value="xpath" :label="$t('XPath')"></el-option>
              <el-option value="css" :label="$t('CSS')"></el-option>
            </el-select>
            <el-input style="width: calc(65% - 10px);"
                      v-model="spiderForm.item_selector"
                      :placeholder="$t('Item Selector')">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('Pagination Selector')"
                        v-if="['list','list-detail'].includes(spiderForm.crawl_type)">
            <el-select style="width: 35%;margin-right: 10px;"
                       v-model="spiderForm.pagination_selector_type"
                       :placeholder="$t('Pagination Selector Type')">
              <el-option value="xpath" :label="$t('XPath')"></el-option>
              <el-option value="css" :label="$t('CSS')"></el-option>
            </el-select>
            <el-input style="width: calc(65% - 10px);"
                      v-model="spiderForm.pagination_selector"
                      :placeholder="$t('Pagination Selector')">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('Item Threshold')"
                        v-if="['list','list-detail'].includes(spiderForm.crawl_type)">
            <el-input-number v-model="spiderForm.item_threshold"/>
          </el-form-item>
        </el-col>
      </el-form>
    </el-row>
    <!--./config detail-->

    <!--button group-->
    <el-row class="button-group-container">
      <div class="button-group">
        <el-button type="danger" @click="onCrawl">{{$t('Run')}}</el-button>
        <el-button type="primary" @click="onExtractFields" v-loading="extractFieldsLoading">{{$t('Extract Fields')}}
        </el-button>
        <el-button type="warning" @click="onPreview" v-loading="previewLoading">{{$t('Preview')}}</el-button>
        <el-button type="success" @click="onSave" v-loading="saveLoading">{{$t('Save')}}</el-button>
      </div>
    </el-row>
    <!--./button group-->

    <!--list field list-->
    <el-row v-if="['list','list-detail'].includes(spiderForm.crawl_type)"
            class="list-fields-container">
      <fields-table-view
        type="list"
        title="List Page Fields"
        :fields="spiderForm.fields"
      />
    </el-row>
    <!--./list field list-->

    <!--detail field list-->
    <el-row v-if="['detail','list-detail'].includes(spiderForm.crawl_type)"
            class="detail-fields-container"
            style="margin-top: 10px;">
      <fields-table-view
        type="detail"
        title="Detail Page Fields"
        :fields="spiderForm.detail_fields"
      />
    </el-row>
    <!--./detail field list-->
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import FieldsTableView from '../TableView/FieldsTableView'

export default {
  name: 'ConfigList',
  components: { FieldsTableView },
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
      columnsDict: {}
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
                    this.$set(this.spiderForm, 'item_selector_type', 'css')
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
  }
}
</script>

<style scoped>

  .button-group-container {
    margin-top: 10px;
    border-bottom: 1px dashed #dcdfe6;
    padding-bottom: 20px;
  }

  .button-group {
    text-align: right;
  }

  .list-fields-container {
    margin-top: 20px;
    border-bottom: 1px dashed #dcdfe6;
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
</style>
