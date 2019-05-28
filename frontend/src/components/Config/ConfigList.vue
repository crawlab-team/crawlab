<template>
  <div class="config-list">
    <!--preview results-->
    <el-dialog :visible.sync="dialogVisible"
               :title="$t('Preview Results')"
               width="90%"
               :before-close="onDialogClose">
      <el-table :data="previewCrawlData"
                :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
                border>
        <el-table-column v-for="(f, index) in fields"
                         :label="f.name"
                         :key="index"
                         min-width="100px">
          <template slot-scope="scope">
            {{scope.row[f.name]}}
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
    <!--./preview results-->

    <!--config detail-->
    <el-row>
      <el-col :span="11" :offset="1">
        <el-form label-width="150px">
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
          <el-form-item :label="$t('Start URL')">
            <el-input v-model="spiderForm.start_url" :placeholder="$t('Start URL')"></el-input>
          </el-form-item>
          <!--<el-form-item :label="$t('URL Pattern')">-->
            <!--<el-input v-model="spiderForm.url_pattern" :placeholder="$t('URL Pattern')"></el-input>-->
          <!--</el-form-item>-->
        </el-form>
      </el-col>
      <el-col :span="11" :offset="1">
        <el-form label-width="150px">
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
          <el-form-item :label="$t('Obey robots.txt')">
            <el-switch v-model="spiderForm.obey_robots_txt" :placeholder="$t('Obey robots.txt')"></el-switch>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
    <!--./config detail-->

    <!--button group-->
    <el-row class="button-group-container">
      <div class="button-group">
        <el-button type="danger" @click="onCrawl">{{$t('Run')}}</el-button>
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
      previewLoading: false,
      saveLoading: false,
      dialogVisible: false
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
    },
    onPreview () {
      this.onSave()
        .then(() => {
          this.previewLoading = true
          this.$store.dispatch('spider/getPreviewCrawlData')
            .then(() => {
              this.dialogVisible = true
            })
            .catch(() => {
              this.$message.error(this.$t('Something wrong happened'))
            })
            .finally(() => {
              this.previewLoading = false
            })
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
        })
    }
  },
  created () {
    // fields for list page
    if (!this.spiderForm.fields) {
      this.spiderForm.fields = []
      for (let i = 0; i < 3; i++) {
        this.spiderForm.fields.push({
          name: `field_${i + 1}`,
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
          name: `field_${i + 1}`,
          type: 'css',
          extract_type: 'text'
        })
      }
    }

    if (!this.spiderForm.crawl_type) this.$set(this.spiderForm, 'crawl_type', 'list')
    if (!this.spiderForm.start_url) this.$set(this.spiderForm, 'start_url', 'http://example.com')
    if (!this.spiderForm.item_selector_type) this.$set(this.spiderForm, 'item_selector_type', 'css')
    if (!this.spiderForm.pagination_selector_type) this.$set(this.spiderForm, 'pagination_selector_type', 'css')
    if (!this.spiderForm.obey_robots_txt) this.$set(this.spiderForm, 'obey_robots_txt', true)
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
</style>
