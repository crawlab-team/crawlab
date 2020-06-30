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
    async created() {
      this.update()
    },
    methods: {
      update() {
        this.isFilterSpiderDisabled = true
        this.$set(this.filter, 'spider_id', this.scheduleForm.spider_id)
        this.filter.schedule_id = this.scheduleForm._id
        this.$store.dispatch('task/getTaskList')
      }
    }
  }
</script>
