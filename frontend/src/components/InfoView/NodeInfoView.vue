<template>
  <div class="info-view">
    <el-row>
      <el-form label-width="150px"
               :model="nodeForm"
               ref="nodeForm"
               class="node-form"
               label-position="right">
        <el-form-item :label="$t('Node Name')">
          <el-input v-model="nodeForm.name" :placeholder="$t('Node Name')" :disabled="isView"></el-input>
        </el-form-item>
        <el-form-item :label="$t('Node IP')" prop="ip" required>
          <el-input v-model="nodeForm.ip" :placeholder="$t('Node IP')" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Node MAC')" prop="ip" required>
          <el-input v-model="nodeForm.mac" :placeholder="$t('Node MAC')" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('Description')">
          <el-input type="textarea" v-model="nodeForm.description" :placeholder="$t('Description')" :disabled="isView">
          </el-input>
        </el-form-item>
      </el-form>
    </el-row>
    <el-row class="button-container" v-if="!isView">
      <el-button size="small" type="success" @click="onSave">{{$t('Save')}}</el-button>
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
              this.$message.success(this.$t('Node info has been saved successfully'))
            })
        }
      })
      this.$st.sendEv('节点详情', '概览', '保存')
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
