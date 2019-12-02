<template>
  <div class="setting-list-table-view">
    <!--    <el-row class="button-group-container">-->
    <!--      <label class="title">{{$t(this.title)}}</label>-->
    <!--      <div class="button-group">-->
    <!--        <el-button type="primary" size="small" @click="addField" icon="el-icon-plus">{{$t('Add Field')}}</el-button>-->
    <!--      </div>-->
    <!--    </el-row>-->
    <el-row>
      <el-table :data="list"
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
        <el-table-column :label="$t('Name')" width="150px">
          <template slot-scope="scope">
            <el-input v-model="scope.row.name"
                      :placeholder="$t('Name')"
                      @change="onNameChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column :label="$t('Remark')" width="auto" min-width="120px">
          <template slot-scope="scope">
            <el-input v-model="scope.row.remark" :placeholder="$t('Remark')"/>
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
  name: 'SettingFieldsTableView',
  props: {
    type: {
      type: String,
      default: 'list'
    },
    title: {
      type: String,
      default: ''
    },
    stageNames: {
      type: Array,
      default () {
        return []
      }
    },
    list: {
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
      this.list.push({
        type: 'css',
        extract_type: 'text'
      })
      this.$st.sendEv('爬虫详情-配置', '添加字段')
    },
    deleteField (index) {
      this.list.splice(index, 1)
      this.$st.sendEv('爬虫详情-配置', '删除字段')
    },
    onNameChange (row) {
      if (this.list.filter(d => d.name === row.name).length > 1) {
        this.$message.error(this.$t(`Duplicated field names for ${row.name}`))
      }
      this.$st.sendEv('爬虫详情-配置', '更改字段')
    },
    onCheck (row) {
      this.list.forEach(d => {
        if (row.name !== d.name) {
          this.$set(d, 'is_detail', false)
        }
      })
      this.$st.sendEv('爬虫详情-配置', '设置详情页URL')
    },
    onClickSelectorType (row, selectorType) {
      if (selectorType === 'css') {
        if (row.xpath) this.$set(row, 'xpath', '')
        if (!row.css) this.$set(row, 'css', 'body')
      } else {
        if (row.css) this.$set(row, 'css', '')
        if (!row.xpath) this.$set(row, 'xpath', '//body')
      }
    },
    onClickIsAttribute (row, isAttribute) {
      if (!isAttribute) {
        // 文本
        if (row.attr) this.$set(row, 'attr', '')
      } else {
        // 属性
        if (!row.attr) this.$set(row, 'attr', 'href')
      }
    },
    onCopyField (row) {
      for (let i = 0; i < this.list.length; i++) {
        if (row.name === this.list[i].name) {
          this.list.splice(i, 0, JSON.parse(JSON.stringify(row)))
          break
        }
      }
    },
    onRemoveField (row) {
      for (let i = 0; i < this.list.length; i++) {
        if (row.name === this.list[i].name) {
          this.list.splice(i, 1)
          break
        }
      }
      if (this.list.length === 0) {
        this.list.push({
          css: 'body',
          next_stage: ''
        })
      }
    },
    onAddField (row) {
      for (let i = 0; i < this.list.length; i++) {
        if (row.name === this.list[i].name) {
          this.list.splice(i + 1, 0, {
            name: `field_${Math.floor(new Date().getTime()).toString()}`,
            css: 'body',
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
      this.list.forEach(f => {
        if (f.name !== row.name) {
          this.$set(f, 'next_stage', '')
        }
      })
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
