<script>
import {
  mapState
} from 'vuex'
import TaskList from '../../views/task/TaskList'

export default {
  name: 'ScheduleTaskList',
  extends: TaskList,
  computed: {
    ...mapState('task', [
      'filter'
    ]),
    ...mapState('schedule', [
      'scheduleForm'
    ])
  },
  async created () {
    this.isFilterSpiderDisabled = true
    this.filter.spider_id = this.scheduleForm.spider_id
    await this.$store.dispatch('task/getTaskList')
  }
}
</script>
