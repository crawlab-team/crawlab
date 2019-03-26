<template>
  <div class="info-view">
    <el-row>
      <el-form label-width="150px"
               :model="nodeForm"
               ref="nodeForm"
               class="node-form"
               label-position="right">
        <el-form-item :label="$t('Node Name')">
          <el-input v-model="nodeForm.name" :placeholder="$t('Node Name')" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Node IP')" prop="ip" required>
          <el-input v-model="nodeForm.ip" :placeholder="$t('Node IP')" :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Node Port')" prop="port" required>
          <el-input v-model="nodeForm.port" :placeholder="$t('Node Port')" :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Description')">
          <el-input type="textarea" v-model="nodeForm.description" :placeholder="$t('Description')" :disabled="isView">
          </el-input>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container" v-if="!isView">
      <el-button type="success" @click="onSave">{{$t('Save')}}</el-button>
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'

export default {
  name: 'NodeInfoView',
  props: {
    isView: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    ...mapState('node', [
      'nodeForm'
    ])
  },
  methods: {
    onSave () {
      this.$refs.nodeForm.validate(valid => {
        if (valid) {
          this.$store.dispatch('node/editNode')
            .then(() => {
              this.$message.success('Node has been saved successfully')
            })
        }
      })
    }
  }
}
</script>

<style scoped>
  .node-form {
    padding: 10px;
  }

  .button-container {
    padding: 0 10px;
    width: 100%;
    text-align: right;
  }
</style>
