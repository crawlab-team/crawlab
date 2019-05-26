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
        <el-table-column v-for="(f, index) in spiderForm.fields"
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

    <el-row>
      <el-col :span="11" offset="1">
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
          <el-form-item :label="$t('Obey robots.txt')">
            <el-switch v-model="spiderForm.obey_robots_txt" :placeholder="$t('Obey robots.txt')"></el-switch>
          </el-form-item>
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
        </el-form>
      </el-col>
    </el-row>

    <!--button group-->
    <el-row style="margin-top: 10px">
      <div class="button-group-wrapper">
        <div class="button-group">
          <el-button type="primary" @click="addField" icon="el-icon-plus">{{$t('Add Field')}}</el-button>
        </div>
        <div class="button-group">
          <el-button type="danger" @click="onCrawl">{{$t('Run')}}</el-button>
          <el-button type="warning" @click="onPreview" v-loading="previewLoading">{{$t('Preview')}}</el-button>
          <el-button type="success" @click="onSave" v-loading="saveLoading">{{$t('Save')}}</el-button>
        </div>
      </div>
    </el-row>
    <!--./button group-->

    <!--field list-->
    <el-row style="margin-top: 10px;">
      <el-table :data="spiderForm.fields"
                class="table edit"
                :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
                border>
        <el-table-column :label="$t('Field Name')" width="200px">
          <template slot-scope="scope">
            <el-input v-model="scope.row.name" :placeholder="$t('Field Name')"></el-input>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Query Type')" width="200px">
          <template slot-scope="scope">
            <el-select v-model="scope.row.type" :placeholder="$t('Query Type')">
              <el-option value="css" :label="$t('CSS Selector')"></el-option>
              <el-option value="xpath" :label="$t('XPath')"></el-option>
            </el-select>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Query')" width="250px">
          <template slot-scope="scope">
            <el-input v-model="scope.row.query" :placeholder="$t('Query')"></el-input>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Extract Type')" width="120px">
          <template slot-scope="scope">
            <el-select v-model="scope.row.extract_type" :placeholder="$t('Extract Type')">
              <el-option value="text" :label="$t('Text')"></el-option>
              <el-option value="attribute" :label="$t('Attribute')"></el-option>
            </el-select>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Attribute')" width="250px">
          <template slot-scope="scope">
            <template v-if="scope.row.extract_type === 'attribute'">
              <el-input v-model="scope.row.attribute"
                        :placeholder="$t('Attribute')">
              </el-input>
            </template>
            <template v-else>
            </template>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Action')" fixed="right">
          <template slot-scope="scope">
            <div class="action-button-group">
              <el-button size="mini" icon="el-icon-delete" type="danger"
                         @click="deleteField(scope.$index)"></el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-row>
    <!--./field list-->
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'ConfigList',
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
    ])
  },
  methods: {
    addField () {
      this.spiderForm.fields.push({
        type: 'css',
        extract_type: 'text'
      })
    },
    deleteField (index) {
      this.spiderForm.fields.splice(index, 1)
    },
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
    if (!this.spiderForm.crawl_type) this.$set(this.spiderForm, 'crawl_type', 'list')
    if (!this.spiderForm.start_url) this.$set(this.spiderForm, 'start_url', 'http://example.com')
    if (!this.spiderForm.item_selector_type) this.$set(this.spiderForm, 'item_selector_type', 'css')
    if (!this.spiderForm.pagination_selector_type) this.$set(this.spiderForm, 'pagination_selector_type', 'css')
    if (!this.spiderForm.obey_robots_txt) this.$set(this.spiderForm, 'obey_robots_txt', true)
  }
}
</script>

<style scoped>
  .el-table {
  }

  .el-table.edit >>> .el-table__body td {
    padding: 0;
  }

  .el-table.edit >>> .el-table__body td .cell {
    padding: 0;
    font-size: 12px;
  }

  .el-table.edit >>> .el-input__inner:hover {
    text-decoration: underline;
  }

  .el-table.edit >>> .el-input__inner {
    height: 36px;
    border: none;
    border-radius: 0;
    font-size: 12px;
  }

  .el-table.edit >>> .el-select .el-input .el-select__caret {
    line-height: 36px;
  }

  .button-group-wrapper {
    display: flex;
    justify-content: space-between;
  }

  .button-group {
    text-align: right;
  }

  .action-button-group {
    margin-left: 10px;
  }
</style>
