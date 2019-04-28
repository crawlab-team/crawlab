<template>
  <div class="spider-stats">
    <!--overall stats-->
    <el-row>
      <div class="metric-list">
        <metric-card :label="$t('30-Day Tasks')"
                     icon="fa fa-play"
                     :value="overviewStats.task_count"
                     type="danger"/>
        <metric-card :label="$t('30-Day Results')"
                     icon="fa fa-table"
                     :value="overviewStats.result_count"
                     type="primary"/>
        <metric-card :label="$t('Success Rate')"
                     icon="fa fa-check"
                     :value="getPercentStr(overviewStats.success_rate)"
                     type="success"/>
        <metric-card :label="$t('Avg Duration (sec)')"
                     icon="fa fa-hourglass"
                     :value="getDecimal(overviewStats.avg_duration)"
                     type="warning"/>
      </div>
    </el-row>
    <!--./overall stats-->

    <el-row>
      <el-col :span="8">
        <el-card class="chart-wrapper" style="margin-right: 20px;">
          <h4>{{$t('Tasks by Status')}}</h4>
          <div id="task-pie" class="chart"></div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card class="chart-wrapper">
          <h4>{{$t('Daily Tasks')}}</h4>
          <div id="task-line" class="chart"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row>
      <el-col :span="8">
        <el-card class="chart-wrapper" style="margin-right: 20px;">
          <h4>{{$t('Tasks by Status')}}</h4>
          <el-table class="table"></el-table>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card class="chart-wrapper">
          <h4>{{$t('Tasks by Status')}}</h4>
          <div id="duration-line" class="chart"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import MetricCard from './MetricCard'
import echarts from 'echarts'

export default {
  name: 'SpiderStats',
  components: { MetricCard },
  methods: {
    renderTaskPie () {
      const chart = echarts.init(this.$el.querySelector('#task-pie'))
      const option = {
        series: [{
          name: '',
          type: 'pie',
          radius: ['50%', '70%'],
          data: this.statusStats.map(d => {
            return {
              name: this.$t(d.name),
              value: d.value
            }
          })
        }]
      }
      chart.setOption(option)
    },

    renderTaskLine () {
      const chart = echarts.init(this.$el.querySelector('#task-line'))
      const option = {
        grid: {
          top: 20,
          bottom: 40
        },
        xAxis: {
          type: 'category',
          data: this.dailyStats.map(d => d.date)
        },
        yAxis: {
          type: 'value'
        },
        series: [{
          type: 'line',
          data: this.dailyStats.map(d => d.count),
          areaStyle: {},
          smooth: true
        }],
        tooltip: {
          trigger: 'axis',
          show: true
        }
      }
      chart.setOption(option)
    },

    render () {
      this.renderTaskPie()
      this.renderTaskLine()
    },

    update () {
      this.$store.dispatch('spider/getSpiderStats')
        .then(() => {
          this.render()
        })
    },

    getPercentStr (value) {
      if (value === undefined) return 'NA'
      return (value * 100).toFixed(2) + '%'
    },

    getDecimal (value) {
      if (value === undefined) return 'NA'
      return value.toFixed(2)
    }
  },
  computed: {
    ...mapState('spider', [
      'overviewStats',
      'statusStats',
      'dailyStats'
    ])
  },
  mounted () {
  }
}
</script>

<style scoped>
  .metric-list {
    display: flex;
  }

  .metric-list .metric-card {
    flex-basis: 25%;
  }

  .chart-wrapper {
    margin-top: 20px;
  }

  .chart {
    width: 100%;
    height: 240px;
  }

  .table {
    height: 240px;
  }

  h4 {
    display: inline-block;
    margin: 0
  }
</style>
