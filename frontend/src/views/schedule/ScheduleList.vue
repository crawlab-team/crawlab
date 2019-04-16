<template>
  <div class="app-container">
    <!--add popup-->
    <el-dialog
      :title="$t('Add Schedule')"
      :visible.sync="dialogVisible"
      width="60%"
      :before-close="onDialogClose">
      <el-form label-width="180px"
               :model="scheduleForm"
               :inline-message="true"
               ref="scheduleForm"
               label-position="right">
        <el-form-item :label="$t('Schedule Name')" prop="url" required>
          <el-input v-model="scheduleForm.name" :placeholder="$t('Schedule Name')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Cron')" prop="cron" required>
          <el-input v-model="scheduleForm.cron" :placeholder="$t('Cron')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Schedule Description')" prop="description">
          <el-input v-model="scheduleForm.description" :placeholder="$t('Schedule Description')"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="onCancel">{{$t('Cancel')}}</el-button>
        <el-button type="primary" @click="onAddSubmit">{{$t('Add')}}</el-button>
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
          <el-tooltip :content="$t('View')" placement="top">
            <el-button type="primary" icon="el-icon-search" size="mini" @click="onView(scope.row)"></el-button>
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
    return {
      columns: [
        { name: 'name', label: 'Name', width: '220' },
        { name: 'cron', label: 'Cron', width: '220' },
        { name: 'description', label: 'Description', width: 'auto' }
      ],
      dialogVisible: false
    }
  },
  computed: {
    ...mapState('schedule', [
      'scheduleList',
      'scheduleForm'
    ]),
    filteredTableData () {
      return this.scheduleList
    }
  },
  methods: {
    onDialogClose () {
    },
    onCancel () {
      this.dialogVisible = false
    },
    onAdd () {
      this.dialogVisible = true
    },
    onAddSubmit () {
    }
  },
  created () {
    this.$store.dispatch('schedule/getScheduleList')
  }
}
</script>

<style scoped>
  .filter .right {
    float: right;
  }
</style>
