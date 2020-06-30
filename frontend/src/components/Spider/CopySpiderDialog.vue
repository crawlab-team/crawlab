<template>
  <el-dialog
    ref="form"
    class="copy-spider-dialog"
    :title="$t('Copy Spider')"
    :visible="visible"
    width="580px"
    :before-close="onClose"
  >
    <el-form
      ref="form"
      label-width="160px"
      :model="form"
    >
      <el-form-item
        :label="$t('New Spider Name')"
        required
      >
        <el-input v-model="form.name" :placeholder="$t('New Spider Name')" />
      </el-form-item>
    </el-form>
    <template slot="footer">
      <el-button type="plain" size="small" @click="$emit('close')">{{ $t('Cancel') }}</el-button>
      <el-button
        type="primary"
        size="small"
        :icon="isLoading ? 'el-icon-loading' : ''"
        :disabled="isLoading"
        @click="onConfirm"
      >
        {{ $t('Confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script>
  export default {
    name: 'CopySpiderDialog',
    props: {
      spiderId: {
        type: String,
        default: ''
      },
      visible: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        form: {
          name: ''
        },
        isLoading: false
      }
    },
    methods: {
      onClose() {
        this.$emit('close')
      },
      onConfirm() {
        this.$refs['form'].validate(async valid => {
          if (!valid) return
          try {
            this.isLoading = true
            const res = await this.$request.post(`/spiders/${this.spiderId}/copy`, this.form)
            if (!res.data.error) {
              this.$message.success('Copied successfully')
            }
            this.$emit('confirm')
            this.$emit('close')
            this.$st.sendEv('爬虫复制', '确认提交')
          } finally {
            this.isLoading = false
          }
        })
      }
    }
  }
</script>

<style scoped>

</style>
