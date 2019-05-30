<template>
  <div class="fields-table-view">
    <el-row class="button-group-container">
      <label class="title">{{$t(this.title)}}</label>
      <div class="button-group">
        <el-button type="primary" size="small" @click="addField" icon="el-icon-plus">{{$t('Add Field')}}</el-button>
      </div>
    </el-row>
    <el-row>
      <el-table :data="fields"
                class="table edit"
                :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
                border>
        <el-table-column v-if="type === 'list' && spiderForm.crawl_type === 'list-detail'"
                         :label="$t('Detail Page URL')"
                         align="center">
          <template slot-scope="scope">
            <el-checkbox v-model="scope.row.is_detail"
                         @change="onCheck(scope.row)">
            </el-checkbox>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Field Name')" width="200px">
          <template slot-scope="scope">
            <el-input v-model="scope.row.name" :placeholder="$t('Field Name')"
                      @change="onNameChange(scope.row)"></el-input>
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
        <el-table-column :label="$t('Action')" fixed="right" min-width="100px">
          <template slot-scope="scope">
            <div class="action-button-group">
              <el-button size="mini"
                         style="margin-left:10px"
                         icon="el-icon-delete"
                         type="danger"
                         @click="deleteField(scope.$index)">
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'FieldsTableView',
  props: {
    type: {
      type: String,
      default: 'list'
    },
    title: {
      type: String,
      default: ''
    },
    fields: {
      type: Array,
      default () {
        return []
      }
    }
  },
  computed: {
    ...mapState('spider', [
      'spiderForm'
    ])
  },
  methods: {
    addField () {
      this.fields.push({
        type: 'css',
        extract_type: 'text'
      })
      this.$st.sendEv('爬虫详情-配置', '添加字段')
    },
    deleteField (index) {
      this.fields.splice(index, 1)
      this.$st.sendEv('爬虫详情-配置', '删除字段')
    },
    onNameChange (row) {
      if (this.fields.filter(d => d.name === row.name).length > 1) {
        this.$message.error(this.$t(`Duplicated field names for ${row.name}`))
      }
      this.$st.sendEv('爬虫详情-配置', '更改字段')
    },
    onCheck (row) {
      this.fields.forEach(d => {
        if (row.name !== d.name) {
          this.$set(d, 'is_detail', false)
        }
      })
      this.$st.sendEv('爬虫详情-配置', '设置详情页URL')
    }
  }
}
</script>

<style scoped>
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

  .button-group-container {
    /*display: inline-block;*/
    /*width: 100%;*/
  }

  .button-group-container .title {
    float: left;
    line-height: 32px;
  }

  .button-group-container .button-group {
    float: right;
  }

  .action-button-group {
    display: flex;
    margin-left: 10px;
  }

  .action-button-group >>> .el-checkbox__label {
    font-size: 12px;
  }
</style>
