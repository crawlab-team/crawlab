<template>
  <div class="info-view">
    <el-row>
      <el-form
        ref="nodeForm"
        label-width="150px"
        :model="nodeForm"
        class="node-form"
        label-position="right"
      >
        <el-form-item :label="$t('Node Name')">
          <el-input v-model="nodeForm.name" :placeholder="$t('Node Name')" :disabled="isView" />
        </el-form-item>
        <el-form-item :label="$t('Node IP')" prop="ip" required>
          <el-input v-model="nodeForm.ip" :placeholder="$t('Node IP')" disabled />
        </el-form-item>
        <el-form-item :label="$t('Node MAC')" prop="ip" required>
          <el-input v-model="nodeForm.mac" :placeholder="$t('Node MAC')" disabled />
        </el-form-item>
        <el-form-item :label="$t('Description')">
          <el-input v-model="nodeForm.description" type="textarea" :placeholder="$t('Description')" :disabled="isView" />
        </el-form-item>
      </el-form>
    </el-row>
    <el-row v-if="!isView" class="button-container">
      <el-button size="small" type="success" @click="onSave">{{ $t('Save') }}</el-button>
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
      onSave() {
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
