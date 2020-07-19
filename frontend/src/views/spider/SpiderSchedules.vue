<script>
  import {
    mapState
  } from 'vuex'
  import ScheduleList from '../schedule/ScheduleList'

  export default {
    name: 'SpiderSchedules',
    extends: ScheduleList,
    computed: {
      ...mapState('spider', [
        'allSpiderList'
      ]),
      isDisabledSpiderSchedule() {
        return true
      },
      spiderId() {
        const arr = this.$route.path.split('/')
        return arr[arr.length - 1]
      }
    },
    created() {
      const arr = this.$route.path.split('/')
      const id = arr[arr.length - 1]
      this.$store.dispatch(`spider/getScheduleList`, { id })

      // 节点列表
      this.getNodeList()

      // 爬虫列表
      this.$store.dispatch('spider/getAllSpiderList')
    },
    methods: {
      getNodeList() {
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
      onAdd() {
        this.isEdit = false
        this.dialogVisible = true
        this.$store.commit('schedule/SET_SCHEDULE_FORM', { node_ids: [], spider_id: this.spiderId })
        if (this.spiderForm.is_scrapy) {
          this.onSpiderChange(this.spiderForm._id)
        }
        this.$st.sendEv('定时任务', '添加定时任务')
      }
    }
  }
</script>
