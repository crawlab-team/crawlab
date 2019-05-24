<template>
  <div class="config-list">
    <el-row>
      <div class="button-group">
        <el-button type="primary" @click="addEnv" icon="el-icon-plus">{{$t('Add Environment Variables')}}</el-button>
        <el-button type="success" @click="save">{{$t('Save')}}</el-button>
      </div>
    </el-row>
    <el-row>
      <el-table :data="spiderForm.fields">
        <el-table-column :label="$t('Field Name')">
          <template slot-scope="scope">
            <el-input v-model="scope.row.name" :placeholder="$t('Variable')"></el-input>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Extract Type')">
          <template slot-scope="scope">
            <el-input v-model="scope.row.type" :placeholder="$t('Value')"></el-input>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Query')">
          <template slot-scope="scope">
            <el-input v-model="scope.row.query" :placeholder="$t('Value')"></el-input>
          </template>
        </el-table-column>
        <el-table-column :label="$t('Action')">
          <template slot-scope="scope">
            <el-button size="mini" icon="el-icon-delete" type="danger" @click="deleteEnv(scope.$index)"></el-button>
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
  name: 'ConfigList',
  computed: {
    ...mapState('spider', [
      'spiderForm'
    ])
  }
}
</script>

<style scoped>

</style>
