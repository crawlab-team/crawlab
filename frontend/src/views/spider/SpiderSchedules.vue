<script>
import ScheduleList from '../schedule/ScheduleList'

export default {
  name: 'SpiderSchedules',
  extends: ScheduleList,
  computed: {
    isDisabledSpiderSchedule () {
      return true
    },
    spiderId () {
      const arr = this.$route.path.split('/')
      return arr[arr.length - 1]
    }
  },
  methods: {
    getNodeList () {
      this.$request.get('/nodes', {}).then(response => {
        this.nodeList = response.data.data.map(d => {
          d.systemInfo = {
            os: '',
            arch: '',
            num_cpu: '',
            executables: []
          }
          return d
        })
      })
    },
    getSpiderList () {
      this.$request.get('/spiders', {})
        .then(response => {
          this.spiderList = response.data.data.list || []
        })
    },
    onAdd () {
      this.isEdit = false
      this.dialogVisible = true
      this.$store.commit('schedule/SET_SCHEDULE_FORM', { node_ids: [], spider_id: this.spiderId })
      if (this.spiderForm.is_scrapy) {
        this.onSpiderChange(this.spiderForm._id)
      }
      this.$st.sendEv('定时任务', '添加定时任务')
    }
  },
  created () {
    const arr = this.$route.path.split('/')
    const id = arr[arr.length - 1]
    this.$store.dispatch(`spider/getScheduleList`, { id })

    // 节点列表
    this.getNodeList()

    // 爬虫列表
    this.getSpiderList()
  }
}
</script>
