<template>
  <el-tag :type="type">
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
    }
  }
}
</script>

<style scoped>

</style>
