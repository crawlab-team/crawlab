<template>
  <div class="app-container">
    <!--add popup-->
    <el-dialog
      :title="$t(dialogTitle)"
      :visible.sync="dialogVisible"
      width="60%"
      :before-close="onDialogClose">
      <el-form label-width="180px"
               :model="scheduleForm"
               :inline-message="true"
               ref="scheduleForm"
               label-position="right">
        <el-form-item :label="$t('Schedule Name')" prop="name" required>
          <el-input v-model="scheduleForm.name" :placeholder="$t('Schedule Name')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Spider')" prop="spider_id" required>
          <el-select v-model="scheduleForm.spider_id" filterable>
            <el-option v-for="op in spiderList" :key="op._id" :value="op._id" :label="op.name"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('Cron')" prop="cron" :rules="cronRules" required>
          <template slot="label">
            <el-tooltip :content="$t('Cron Format: [second] [minute] [hour] [day of month] [month] [day of week]')"
                        placement="top">
              <span>
                {{$t('Cron')}}
                <i class="fa fa-exclamation-circle"></i>
              </span>
            </el-tooltip>
          </template>
          <el-input v-model="scheduleForm.cron" :placeholder="$t('Cron')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Execute Command')" prop="params">
          <el-input v-model="spider.cmd"
                    :placeholder="$t('Execute Command')"
                    disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Parameters')" prop="params">
          <el-input v-model="scheduleForm.params"
                    :placeholder="$t('Parameters')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Schedule Description')" prop="description">
          <el-input v-model="scheduleForm.description" type="textarea"
                    :placeholder="$t('Schedule Description')"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="onCancel">{{$t('Cancel')}}</el-button>
        <el-button type="primary" @click="onAddSubmit">{{$t('Submit')}}</el-button>
      </span>
    </el-dialog>

    <!--filter-->
    <div class="filter">
      <div class="right">
        <el-button type="primary"
                   icon="el-icon-plus"
                   class="refresh"
                   @click="onAdd">
          {{$t('Add Schedule')}}
        </el-button>
      </div>
    </div>

    <!--table list-->
    <el-table :data="filteredTableData"
              class="table"
              :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
              border>
      <template v-for="col in columns">
        <el-table-column :key="col.name"
                         :property="col.name"
                         :label="$t(col.label)"
                         :sortable="col.sortable"
                         align="center"
                         :width="col.width">
        </el-table-column>
      </template>
      <el-table-column :label="$t('Action')" align="left" width="250">
        <template slot-scope="scope">
          <el-tooltip :content="$t('Edit')" placement="top">
            <el-button type="warning" icon="el-icon-edit" size="mini" @click="onEdit(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip :content="$t('Remove')" placement="top">
            <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip v-if="isShowRun(scope.row)" :content="$t('Run')" placement="top">
            <el-button type="success" icon="fa fa-bug" size="mini" @click="onCrawl(scope.row)"></el-button>
          </el-tooltip>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'ScheduleList',
  data () {
    const cronValidator = (rule, value, callback) => {
      let patArr = []
      for (let i = 0; i < 6; i++) {
        patArr.push('[/*,0-9]+')
      }
      const pat = '^' + patArr.join(' ') + '$'
      if (!value) {
        callback(new Error('cron cannot be empty'))
      } else if (!value.match(pat)) {
        callback(new Error('cron format is invalid'))
      }
      callback()
    }
    return {
      columns: [
        { name: 'name', label: 'Name', width: '220' },
        { name: 'cron', label: 'Cron', width: '220' },
        { name: 'description', label: 'Description', width: 'auto' }
      ],
      isEdit: false,
      dialogTitle: '',
      dialogVisible: false,
      cronRules: [
        { validator: cronValidator, trigger: 'blur' }
      ]
    }
  },
  computed: {
    ...mapState('schedule', [
      'scheduleList',
      'scheduleForm'
    ]),
    ...mapState('spider', [
      'spiderList'
    ]),
    filteredTableData () {
      return this.scheduleList
    },
    spider () {
      for (let i = 0; i < this.spiderList.length; i++) {
        if (this.spiderList[i]._id === this.scheduleForm.spider_id) {
          return this.spiderList[i]
        }
      }
      return {}
    }
  },
  methods: {
    onDialogClose () {
      this.dialogVisible = false
    },
    onCancel () {
      this.dialogVisible = false
    },
    onAdd () {
      this.isEdit = false
      this.dialogVisible = true
      this.$store.commit('schedule/SET_SCHEDULE_FORM', {})
    },
    onAddSubmit () {
      this.$refs.scheduleForm.validate(res => {
        if (res) {
          let action
          if (this.isEdit) {
            action = 'editSchedule'
          } else {
            action = 'addSchedule'
          }
          this.$store.dispatch('schedule/' + action, this.scheduleForm._id)
            .then(() => {
              this.dialogVisible = false
              setTimeout(() => {
                this.$store.dispatch('schedule/getScheduleList')
              }, 100)
            })
        }
      })
    },
    isShowRun () {
    },
    onEdit (row) {
      this.$store.commit('schedule/SET_SCHEDULE_FORM', row)
      this.dialogVisible = true
      this.isEdit = true
    },
    onRemove (row) {
      this.$store.dispatch('schedule/removeSchedule', row._id)
        .then(() => {
          setTimeout(() => {
            this.$store.dispatch('schedule/getScheduleList')
            this.$message.success(`Schedule "${row.name}" has been removed`)
          }, 100)
        })
    },
    onCrawl () {
    }
  },
  created () {
    this.$store.dispatch('schedule/getScheduleList')
    this.$store.dispatch('spider/getSpiderList')
  }
}
</script>

<style scoped>
  .filter .right {
    text-align: right;
  }

  .table {
    margin-top: 10px;
  }
</style>
