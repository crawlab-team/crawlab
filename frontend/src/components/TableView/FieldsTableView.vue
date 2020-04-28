<template>
  <div class="fields-table-view">
    <el-row>
      <el-table :data="fields"
                class="table edit"
                :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
                :cell-style="getCellClassStyle"
      >
        <el-table-column class-name="action" width="80px" align="right">
          <template slot-scope="scope">
            <i class="action-item el-icon-copy-document" @click="onCopyField(scope.row)"></i>
            <i class="action-item el-icon-remove-outline" @click="onRemoveField(scope.row)"></i>
            <i class="action-item el-icon-circle-plus-outline" @click="onAddField(scope.row)"></i>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Field Name')" width="150px">
          <template slot-scope="scope">
            <el-input v-model="scope.row.name"
                      :placeholder="$t('Field Name')"
                      suffix-icon="el-icon-edit"
                      @change="onNameChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column :label="$t('Selector Type')" width="150px" align="center" class-name="selector-type">
          <template slot-scope="scope">
            <span class="button-selector-item" @click="onClickSelectorType(scope.row, 'css')">
              <el-tag
                :class="scope.row.css ? 'active' : 'inactive'"
                type="success"
              >
                CSS
              </el-tag>
            </span>
            <span class="button-selector-item" @click="onClickSelectorType(scope.row, 'xpath')">
              <el-tag
                :class="scope.row.xpath ? 'active' : 'inactive'"
                type="primary"
              >
                XPath
              </el-tag>
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Selector')" width="200px">
          <template slot-scope="scope">
            <template v-if="scope.row.css">
              <el-input
                v-model="scope.row.css"
                :placeholder="$t('CSS / XPath')"
                suffix-icon="el-icon-edit"
              >
              </el-input>
            </template>
            <template v-else>
              <el-input
                v-model="scope.row.xpath"
                :placeholder="$t('CSS / XPath')"
                suffix-icon="el-icon-edit"
              >
              </el-input>
            </template>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Is Attribute')" width="150px" align="center">
          <template slot-scope="scope">
            <span class="button-selector-item" @click="onClickIsAttribute(scope.row, false)">
              <el-tag
                :class="!isShowAttr(scope.row) ? 'active' : 'inactive'"
                type="success"
              >
                {{$t('Text')}}
              </el-tag>
            </span>
            <span class="button-selector-item" @click="onClickIsAttribute(scope.row, true)">
              <el-tag
                :class="isShowAttr(scope.row) ? 'active' : 'inactive'"
                type="primary"
              >
                {{$t('Attribute')}}
              </el-tag>
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Attribute')" width="200px">
          <template slot-scope="scope">
            <template v-if="isShowAttr(scope.row)">
              <el-input
                v-model="scope.row.attr"
                :placeholder="$t('Attribute')"
                suffix-icon="el-icon-edit"
                @change="onAttrChange(scope.row)"
              />
            </template>
            <template v-else>
              <span style="margin-left: 15px; color: lightgrey">
                N/A
              </span>
            </template>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Next Stage')" width="250px">
          <template slot-scope="scope">
            <el-select
              v-model="scope.row.next_stage"
              :class="!scope.row.next_stage ? 'disabled' : ''"
              @change="onChangeNextStage(scope.row)"
            >
              <el-option :label="$t('No Next Stage')" value=""/>
              <el-option v-for="n in filteredStageNames" :key="n" :label="n" :value="n"/>
            </el-select>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Remark')" width="auto" min-width="120px">
          <template slot-scope="scope">
            <el-input v-model="scope.row.remark" :placeholder="$t('Remark')" suffix-icon="el-icon-edit"/>
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
    stage: {
      type: Object,
      default () {
        return {}
      }
    },
    stageNames: {
      type: Array,
      default () {
        return []
      }
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
    ]),
    filteredStageNames () {
      return this.stageNames.filter(n => n !== this.stage.name)
    }
  },
  methods: {
    onNameChange (row) {
      if (this.fields.filter(d => d.name === row.name).length > 1) {
        this.$message.error(this.$t(`Duplicated field names for ${row.name}`))
      }
      this.$st.sendEv('爬虫详情', '配置', '更改字段')
    },
    onClickSelectorType (row, selectorType) {
      this.$st.sendEv('爬虫详情', '配置', `点击字段选择器类别-${selectorType}`)
      if (selectorType === 'css') {
        if (row.xpath) this.$set(row, 'xpath', '')
        if (!row.css) this.$set(row, 'css', 'body')
      } else {
        if (row.css) this.$set(row, 'css', '')
        if (!row.xpath) this.$set(row, 'xpath', '//body')
      }
    },
    onClickIsAttribute (row, isAttribute) {
      this.$st.sendEv('爬虫详情', '配置', '设置字段属性')
      if (!isAttribute) {
        // 文本
        if (row.attr) this.$set(row, 'attr', '')
      } else {
        // 属性
        if (!row.attr) this.$set(row, 'attr', 'href')
      }
      this.$set(row, 'isAttrChange', false)
    },
    onCopyField (row) {
      for (let i = 0; i < this.fields.length; i++) {
        if (row.name === this.fields[i].name) {
          this.fields.splice(i, 0, JSON.parse(JSON.stringify(row)))
          break
        }
      }
    },
    onRemoveField (row) {
      this.$st.sendEv('爬虫详情', '配置', '删除字段')
      for (let i = 0; i < this.fields.length; i++) {
        if (row.name === this.fields[i].name) {
          this.fields.splice(i, 1)
          break
        }
      }
      if (this.fields.length === 0) {
        this.fields.push({
          xpath: '//body',
          next_stage: ''
        })
      }
    },
    onAddField (row) {
      this.$st.sendEv('爬虫详情', '配置', '添加字段')
      for (let i = 0; i < this.fields.length; i++) {
        if (row.name === this.fields[i].name) {
          this.fields.splice(i + 1, 0, {
            name: `field_${Math.floor(new Date().getTime()).toString()}`,
            xpath: '//body',
            next_stage: ''
          })
          break
        }
      }
    },
    getCellClassStyle ({ row, columnIndex }) {
      if (columnIndex === 1) {
        // 字段名称
        if (!row.name) {
          return {
            'border': '1px solid red'
          }
        }
      } else if (columnIndex === 3) {
        // 选择器
        if (!row.css && !row.xpath) {
          return {
            'border': '1px solid red'
          }
        }
      }
    },
    onChangeNextStage (row) {
      this.fields.forEach(f => {
        if (f.name !== row.name) {
          this.$set(f, 'next_stage', '')
        }
      })
    },
    onAttrChange (row) {
      this.$set(row, 'isAttrChange', !row.attr)
    },
    isShowAttr (row) {
      return (row.attr || row.isAttrChange)
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

  .el-table.edit >>> .button-selector-item {
    cursor: pointer;
    margin: 0 5px;
  }

  .el-table.edit >>> .el-tag.inactive {
    opacity: 0.5;
  }

  .el-table.edit >>> .action {
    background: none !important;
    border: none;
  }

  .el-table.edit >>> tr {
    border: none;
  }

  .el-table.edit >>> tr th {
    border-right: 1px solid rgb(220, 223, 230);
  }

  .el-table.edit >>> tr td:nth-child(2) {
    border-left: 1px solid rgb(220, 223, 230);
  }

  .el-table.edit >>> tr td {
    border-right: 1px solid rgb(220, 223, 230);
  }

  .el-table.edit::before {
    background: none;
  }

  .el-table.edit >>> .action-item {
    font-size: 14px;
    margin-right: 5px;
    cursor: pointer;
  }

  .el-table.edit >>> .action-item:last-child {
    margin-right: 10px;
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

  .el-table.edit >>> .el-select.disabled .el-input__inner {
    color: lightgrey;
  }
</style>
