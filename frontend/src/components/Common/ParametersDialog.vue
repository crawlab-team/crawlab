<template>
  <el-dialog
    :title="$t('Parameters')"
    :visible="visible"
    :before-close="beforeClose"
    class="parameters-dialog"
    width="720px"
  >
    <div class="action-wrapper">
      <el-button
        type="primary"
        size="small"
        @click="onAdd"
      >
        {{$t('Add')}}
      </el-button>
    </div>
    <el-table
      :data="paramData"
      border
    >
      <el-table-column
        :label="$t('Parameter Type')"
        width="100px"
      >
        <template slot-scope="scope">
          <el-select v-model="scope.row.type" size="small">
            <el-option
              :label="$t('Spider')"
              value="spider"
            />
            <el-option
              :label="$t('Setting')"
              value="setting"
            />
            <el-option
              :label="$t('Other')"
              value="other"
            />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('Parameter Name')"
        width="240px"
      >
        <template slot-scope="scope">
          <el-autocomplete
            v-if="scope.row.type === 'setting'"
            v-model="scope.row.name"
            size="small"
            suffix-icon="el-icon-edit"
            :fetch-suggestions="querySearch"
          />
          <el-input
            v-else-if="scope.row.type === 'spider'"
            v-model="scope.row.name"
            size="small"
            suffix-icon="el-icon-edit"
          />
          <div v-else style="text-align: center">
            N/A
          </div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('Parameter Value')"
      >
        <template slot-scope="scope">
          <el-input v-model="scope.row.value" size="small" suffix-icon="el-icon-edit"/>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('Action')"
        width="60px"
        align="center"
      >
        <template slot-scope="scope">
          <div class="action-btn-wrapper">
            <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.$index)" circle/>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <template slot="footer">
      <el-button type="plain" size="small" @click="$emit('close')">{{$t('Cancel')}}</el-button>
      <el-button type="primary" size="small" @click="onConfirm">
        {{$t('Confirm')}}
      </el-button>
    </template>
  </el-dialog>
</template>

<script>
export default {
  name: 'ParametersDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    param: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      paramData: []
    }
  },
  watch: {
    visible (value) {
      if (value) this.initParamData()
    }
  },
  methods: {
    beforeClose () {
      this.$emit('close')
    },
    initParamData () {
      const mArr = this.param.match(/((?:-[-a-zA-Z0-9] )?(?:\w+=)?\w+)/g)
      if (!mArr) {
        this.paramData = []
        this.paramData.push({ type: 'spider', name: '', value: '' })
        return
      }
      this.paramData = []
      mArr.forEach(s => {
        s = s.trim()
        let d = {}
        const arr = s.split(' ')
        if (arr.length === 1) {
          d.type = 'other'
          d.value = s
        } else {
          const arr2 = arr[1].split('=')
          d.name = arr2[0]
          d.value = arr2[1]
          if (arr[0] === '-a') {
            d.type = 'spider'
          } else if (arr[0] === '-s') {
            d.type = 'setting'
          } else {
            d.type = 'other'
            d.value = s
          }
        }
        this.paramData.push(d)
      })
      if (this.paramData.length === 0) {
        this.paramData.push({ type: 'spider', name: '', value: '' })
      }
    },
    onConfirm () {
      const param = this.paramData
        .filter(d => d.value)
        .map(d => {
          let s = ''
          if (d.type === 'setting') {
            s = `-s ${d.name}=${d.value}`
          } else if (d.type === 'spider') {
            s = `-a ${d.name}=${d.value}`
          } else if (d.type === 'other') {
            s = d.value
          }
          return s
        })
        .filter(s => !!s)
        .join(' ')
      this.$emit('confirm', param)
    },
    onRemove (index) {
      this.paramData.splice(index, 1)
    },
    onAdd () {
      this.paramData.push({ type: 'spider', name: '', value: '' })
    },
    querySearch (queryString, cb) {
      let data = this.$utils.scrapy.settingParamNames
      if (!queryString) {
        return cb(data.map(s => {
          return { value: s, label: s }
        }))
      }
      data = data
        .filter(s => s.match(new RegExp(queryString, 'i')))
        .sort((a, b) => a < b ? 1 : -1)
      cb(data.map(s => {
        return { value: s, label: s }
      }))
    }
  }
}
</script>

<style scoped>
  .parameters-dialog >>> .el-table td,
  .parameters-dialog >>> .el-table td .cell {
    padding: 0;
    margin: 0;
  }

  .parameters-dialog >>> .el-table td .cell .el-autocomplete {
    width: 100%;
  }

  .parameters-dialog >>> .el-table td .cell .el-input__inner {
    border: none;
  }

  .parameters-dialog .action-wrapper {
    margin-bottom: 10px;
    text-align: right;
  }

  .parameters-dialog .action-btn-wrapper {
  }
</style>
