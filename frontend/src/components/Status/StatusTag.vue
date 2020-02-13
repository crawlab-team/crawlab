<template>
  <el-tag :type="type" class="status-tag">
    <i :class="icon"></i>
    {{$t(label)}}
  </el-tag>
</template>

<script>
export default {
  name: 'StatusTag',
  props: {
    status: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      statusDict: {
        pending: { label: 'Pending', type: 'primary' },
        running: { label: 'Running', type: 'warning' },
        finished: { label: 'Finished', type: 'success' },
        error: { label: 'Error', type: 'danger' },
        cancelled: { label: 'Cancelled', type: 'info' }
      }
    }
  },
  computed: {
    type () {
      const s = this.statusDict[this.status]
      if (s) {
        return s.type
      }
      return ''
    },
    label () {
      const s = this.statusDict[this.status]
      if (s) {
        return s.label
      }
      return 'NA'
    },
    icon () {
      if (this.status === 'finished') {
        return 'el-icon-check'
      } else if (this.status === 'running') {
        return 'el-icon-loading'
      } else if (this.status === 'error') {
        return 'el-icon-error'
      }
      return ''
    }
  }
}
</script>

<style scoped>

</style>
